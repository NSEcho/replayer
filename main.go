package main

import (
	"flag"

	"github.com/lateralusd/replayer/config"
	"github.com/lateralusd/replayer/replay"
)

// obtain the method and path
// obtain the host header
// craft the url
// craft the response
// extract cookies
// extract headers
// extract body

var (
	count   = flag.Int("c", 1, "how many times to send request")
	timeout = flag.Int("t", 10, "timeout for request")
	proxy   = flag.String("p", "", "proxy for requests in format https?://ip:port")
	stdout  = flag.Bool("s", false, "print response on the stdout")
)

func main() {
	flag.Parse()
	r := replay.NewReplayer("req")
	r.Replay(&config.ReplayerConfig{
		Count:         *count,
		Timeout:       *timeout,
		Proxy:         *proxy,
		PrintOnStdout: *stdout,
	})
}
