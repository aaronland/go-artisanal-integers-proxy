prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/aaronland/go-artisanal-integers-proxy
	cp *.go src/github.com/aaronland/go-artisanal-integers-proxy/
	cp -r service src/github.com/aaronland/go-artisanal-integers-proxy/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:
	@GOPATH=$(shell pwd) go get "github.com/aaronland/go-artisanal-integers"
	@GOPATH=$(shell pwd) go get "github.com/aaronland/go-brooklynintegers-api"
	@GOPATH=$(shell pwd) go get "github.com/aaronland/go-missionintegers-api"
	@GOPATH=$(shell pwd) go get "github.com/aaronland/go-londonintegers-api"
	@GOPATH=$(shell pwd) go get "github.com/whosonfirst/go-whosonfirst-pool"
	@GOPATH=$(shell pwd) go get "github.com/whosonfirst/go-whosonfirst-log"
	rm -rf src/github.com/aaronland/go-brooklynintegers-api/vendor/github.com/aaronland/go-artisanal-integers
	rm -rf src/github.com/aaronland/go-missionintegers-api/vendor/github.com/aaronland/go-artisanal-integers
	rm -rf src/github.com/aaronland/go-londonintegers-api/vendor/github.com/aaronland/go-artisanal-integers

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt cmd/*.go
	go fmt service/*.go

bin:	self
	@GOPATH=$(shell pwd) go build -o bin/proxy-server cmd/proxy-server.go

