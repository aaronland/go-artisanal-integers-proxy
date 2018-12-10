package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-brooklynintegers-api"
	"github.com/aaronland/go-brooklynintegers-proxy"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-pool"
	"io"
	"net/http"
	"os"
	"strconv"
)

// this needs to be tweaked to keep a not-just-in-memory copy of the
// pool so that we can use this in offline-mode (20181206/thisisaaronland)

func main() {

	var port = flag.Int("port", 8080, "Port to listen")
	var min = flag.Int("min", 5, "The minimum number of Brooklyn Integers to keep on hand at all times")
	var loglevel = flag.String("loglevel", "info", "Log level")
	var cors = flag.Bool("cors", false, "Enable CORS headers")

	flag.Parse()

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger("[big-integer] ")
	logger.AddLogger(writer, *loglevel)

	cl := api.NewAPIClient()

	pl, err := pool.NewMemLIFOPool()

	if err != nil {
		logger.Fatal(err)
	}

	pr := proxy.NewProxy(cl, pl, int64(*min), logger)
	pr.Init()

	handler := func(rsp http.ResponseWriter, r *http.Request) {

		i, err := pr.Integer()

		if err != nil {
			msg := fmt.Sprintf("Failed to retrieve integer because %v", err)
			http.Error(rsp, msg, http.StatusBadRequest)
		}

		if *cors {
			rsp.Header().Set("Access-Control-Allow-Origin", "*")
			return
		}

		io.WriteString(rsp, strconv.FormatInt(i, 10))
	}

	http.HandleFunc("/", handler)

	str_port := ":" + strconv.Itoa(*port)
	err = http.ListenAndServe(str_port, nil)

	if err != nil {
		logger.Fatal("Failed to start server, because %v\n", err)
	}

}
