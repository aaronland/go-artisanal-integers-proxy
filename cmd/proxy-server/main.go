package main

import (
	"flag"
	"github.com/aaronland/go-artisanal-integers-proxy"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-pool"
	"io"
	"os"
)

func main() {

	var protocol = flag.String("protocol", "http", "The protocol to use for the proxy server.")
	var host = flag.String("host", "localhost", "Host to listen on.")
	var port = flag.Int("port", 8080, "Port to listen on.")
	var min = flag.Int("min", 5, "The minimum number of artisanal integers to keep on hand at all times.")
	var loglevel = flag.String("loglevel", "info", "Log level.")

	var brooklyn_integers = flag.Bool("brooklyn-integers", false, "Use Brooklyn Integers as an artisanal integer source.")
	var london_integers = flag.Bool("london-integers", false, "Use London Integers as an artisanal integer source.")
	var mission_integers = flag.Bool("mission-integers", false, "Use Mission Integers as an artisanal integer source.")

	flag.Parse()

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger("[proxy-server]")
	logger.AddLogger(writer, *loglevel)

	pl, err := pool.NewMemLIFOPool()

	if err != nil {
		logger.Fatal(err)
	}

	svc_args := proxy.ProxyServiceArgs{
		BrooklynIntegers: *brooklyn_integers,
		LondonIntegers:   *london_integers,
		MissionIntegers:  *mission_integers,
		MinCount:         *min,
		// Logger: logger,
	}

	svc, err := proxy.NewProxyServiceWithPool(pl, svc_args)

	svr_args := proxy.ProxyServerArgs{
		Protocol: *protocol,
		Host:     *host,
		Port:     *port,
		// Logger: logger,
	}

	svr, err := proxy.NewProxyServerWithService(svc, svr_args)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Status("Listening for requests on %s", svr.Address())

	err = svr.ListenAndServe(svc)

	if err != nil {
		logger.Fatal(err)
	}

	os.Exit(0)
}
