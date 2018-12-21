package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-artisanal-integers"
	"github.com/aaronland/go-artisanal-integers-proxy/service"
	"github.com/aaronland/go-artisanal-integers/server"
	brooklyn_api "github.com/aaronland/go-brooklynintegers-api"
	london_api "github.com/aaronland/go-londonintegers-api"
	mission_api "github.com/aaronland/go-missionintegers-api"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-pool"
	"io"
	"net/url"
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

	logger := log.NewWOFLogger("[proxy-server] ")
	logger.AddLogger(writer, *loglevel)

	// set up one or more clients to proxy integers from

	clients := make([]artisanalinteger.Client, 0)

	if *brooklyn_integers {
		cl := brooklyn_api.NewAPIClient()
		clients = append(clients, cl)
	}

	if *london_integers {
		cl := london_api.NewAPIClient()
		clients = append(clients, cl)
	}

	if *mission_integers {
		cl := mission_api.NewAPIClient()
		clients = append(clients, cl)
	}

	if len(clients) == 0 {
		logger.Fatal("Insufficient clients")
	}

	// set up a local pool for proxied integers

	// this needs to be tweaked to keep a not-just-in-memory copy of the
	// pool so that we can use this in offline-mode (20181206/thisisaaronland)

	pl, err := pool.NewMemLIFOPool()

	if err != nil {
		logger.Fatal(err)
	}

	// set up the proxy service

	opts, err := service.DefaultProxyServiceOptions()

	if err != nil {
		logger.Fatal(err)
	}

	opts.Logger = logger
	opts.Pool = pl
	opts.Minimum = *min

	pr, err := service.NewProxyService(opts, clients...)

	if err != nil {
		logger.Fatal(err)
	}

	// set up the actual server endpoint

	addr := fmt.Sprintf("%s://%s:%d", *protocol, *host, *port)
	u, err := url.Parse(addr)

	if err != nil {
		logger.Fatal(err)
	}

	svr, err := server.NewArtisanalServer(*protocol, u)

	if err != nil {
		logger.Fatal(err)
	}

	// go!

	logger.Status("Listening for requests on %s", svr.Address())

	err = svr.ListenAndServe(pr)

	if err != nil {
		logger.Fatal(err)
	}

	os.Exit(0)
}
