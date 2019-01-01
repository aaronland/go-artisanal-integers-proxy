package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronland/go-artisanal-integers"
	"github.com/aaronland/go-artisanal-integers-proxy/service"
	"github.com/aaronland/go-artisanal-integers/server"
	brooklyn_api "github.com/aaronland/go-brooklynintegers-api"
	london_api "github.com/aaronland/go-londonintegers-api"
	mission_api "github.com/aaronland/go-missionintegers-api"
	"github.com/whosonfirst/go-whosonfirst-pool"
	"net/url"
)

type ProxyServerArgs struct {
	Protocol string
	Host     string
	Port     int
}

type ProxyServiceArgs struct {
	BrooklynIntegers bool `json:"brooklyn_integers"`
	LondonIntegers   bool `json:"london_integers"`
	MissionIntegers  bool `json:"mission_integers"`
	MinCount         int  `json:"min_count"`
}

type ProxyServiceResponse struct {
	Integer int64 `json:"integer"`
}

type ProxyServiceLambdaFunc func(context.Context, ProxyServiceArgs) (*ProxyServiceResponse, error)

func NewProxyServiceWithPool(pl pool.LIFOPool, args ProxyServiceArgs) (artisanalinteger.Service, error) {

	opts, err := service.DefaultProxyServiceOptions()

	if err != nil {
		return nil, err
	}

	opts.Pool = pl
	opts.Minimum = args.MinCount

	clients := make([]artisanalinteger.Client, 0)

	if args.BrooklynIntegers {
		cl := brooklyn_api.NewAPIClient()
		clients = append(clients, cl)
	}

	if args.LondonIntegers {
		cl := london_api.NewAPIClient()
		clients = append(clients, cl)
	}

	if args.MissionIntegers {
		cl := mission_api.NewAPIClient()
		clients = append(clients, cl)
	}

	if len(clients) == 0 {
		return nil, errors.New("Insufficient clients")
	}

	return service.NewProxyService(opts, clients...)
}

func NewProxyServiceLambdaFunc(pl pool.LIFOPool) (ProxyServiceLambdaFunc, error) {

	f := func(ctx context.Context, args ProxyServiceArgs) (*ProxyServiceResponse, error) {

		svc, err := NewProxyServiceWithPool(pl, args)

		if err != nil {
			return nil, err
		}

		i, err := svc.NextInt()

		if err != nil {
			return nil, err
		}

		rsp := ProxyServiceResponse{
			Integer: i,
		}

		return &rsp, nil
	}

	return f, nil
}

func NewProxyServerWithService(svc artisanalinteger.Service, args ProxyServerArgs) (artisanalinteger.Server, error) {

	addr := fmt.Sprintf("%s://%s:%d", args.Protocol, args.Host, args.Port)
	u, err := url.Parse(addr)

	if err != nil {
		return nil, err
	}

	return server.NewArtisanalServer(args.Protocol, u)
}
