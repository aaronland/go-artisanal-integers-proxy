.PHONY=tools
tools:
	go build -o bin/proxy-server cmd/proxy-server/main.go

.PHONY=fmt
fmt:
	go fmt **/*.go

