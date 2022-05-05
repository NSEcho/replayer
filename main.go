package main

import (
	"flag"
	"fmt"
	"os"

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
	args := flag.Args()
	if len(args) == 0 || len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Please provide single filename")
		os.Exit(1)
	}
	r := replay.NewReplayer(args[0])
	err := r.Replay(&config.ReplayerConfig{
		Count:         *count,
		Timeout:       *timeout,
		Proxy:         *proxy,
		PrintOnStdout: *stdout,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error replaying request: %+v\n", err)
	}
}
