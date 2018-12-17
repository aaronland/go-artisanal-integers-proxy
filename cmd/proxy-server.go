package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/server"
	"github.com/aaronland/go-brooklynintegers-api"
	"github.com/aaronland/go-brooklynintegers-proxy/service"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-pool"
	"io"
	"net/url"
	"os"
)

// this needs to be tweaked to keep a not-just-in-memory copy of the
// pool so that we can use this in offline-mode (20181206/thisisaaronland)

func main() {

	var host = flag.String("host", "localhost", "Host to listen on")
	var port = flag.Int("port", 8080, "Port to listen on")
	var min = flag.Int("min", 5, "The minimum number of artisanal integers to keep on hand at all times")
	var loglevel = flag.String("loglevel", "info", "Log level")

	flag.Parse()

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger("[big-integer] ")
	logger.AddLogger(writer, *loglevel)

	cl := api.NewAPIClient()

	pl, err := pool.NewMemLIFOPool()

	if err != nil {
		logger.Fatal(err)
	}

	opts, err := service.DefaultProxyServiceOptions()

	if err != nil {
		logger.Fatal(err)
	}

	opts.Logger = logger
	opts.Pool = pl
	opts.Minimum = *min

	pr, err := service.NewProxyService(opts, cl)

	if err != nil {
		logger.Fatal(err)
	}

	addr := fmt.Sprintf("http://%s:%d", *host, *port)
	u, err := url.Parse(addr)

	if err != nil {
		logger.Fatal(err)
	}

	svr, err := server.NewHTTPServer(u)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Status("Listening for requests on %s", svr.Address())

	err = svr.ListenAndServe(pr)

	if err != nil {
		logger.Fatal(err)
	}

	os.Exit(0)
}
