package replay

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lateralusd/replayer/config"
)

var (
	methodPathRe = regexp.MustCompile(`^([A-Z]+)\s(\/.*?)\s`)
	hostRe       = regexp.MustCompile(`Host:\s(.*)`)
	headersRe    = regexp.MustCompile(`([A-Za-z0-9-]+):\s(.*)`)
)

func NewReplayer(filename string) *Replayer {
	return &Replayer{
		filename: filename,
	}
}

type Replayer struct {
	filename string
}

type reqData struct {
	method  string
	query   string
	host    string
	headers map[string]string
}

func (r *Replayer) Replay(cfg *config.ReplayerConfig) error {
	f, err := os.Open(r.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	methodQuery := extractSingleData(data, methodPathRe)
	host := extractSingleData(data, hostRe)
	headers := extractAllData(data, headersRe)

	rData := reqData{
		method:  methodQuery[0],
		query:   methodQuery[1],
		host:    host[0],
		headers: make(map[string]string),
	}

	for k, v := range headers {
		rData.headers[k] = v
	}

	u, err := craftURL(rData.host, rData.query)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(rData.method, u.String(), nil)
	if err != nil {
		return err
	}

	for k, v := range rData.headers {
		if k != "Cookie" && k != "Host" {
			req.Header.Add(k, v)
		}
		if k == "Cookie" {
			cookies := extractCookies(v)
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
		}
	}

	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	if cfg.Proxy != "" {
		fmt.Println("Not")
		proxyURL, err := url.Parse(cfg.Proxy)
		if err != nil {
			return err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if cfg.PrintHeaders {
		fmt.Printf("%s %d %s\n",
			resp.Proto,
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
		for key, val := range resp.Header {
			fmt.Printf("%s: %s\n", key, strings.Join(val, " "))
		}
		fmt.Printf("\n\n")
	}

	if cfg.PrintOnStdout {
		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}
		io.Copy(os.Stdout, reader)
	}

	return nil
}

func extractCookies(cookies string) []*http.Cookie {
	var res []*http.Cookie
	splitted := strings.Split(cookies, ";")
	for _, spl := range splitted {
		cookie := strings.Split(spl, "=")
		if len(cookie) >= 2 {
			res = append(res, &http.Cookie{
				Name:  cookie[0],
				Value: strings.Join(cookie[1:], "="),
			})
		}
	}
	return res
}

func craftURL(host, query string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("http://%s%s", host, query))
}

func extractSingleData(content []byte, re *regexp.Regexp) []string {
	matches := re.FindStringSubmatch(string(content))
	if len(matches) == 0 {
		return []string{}
	}
	return matches[1:]
}

func extractAllData(content []byte, re *regexp.Regexp) map[string]string {
	matches := re.FindAllStringSubmatch(string(content), -1)
	data := make(map[string]string)
	for _, match := range matches {
		if len(match) > 1 {
			data[match[1]] = match[2]
		}
	}
	return data
}
