export GO111MODULE=on
export GOPROXY=direct
export GOSUMDB=off
export CGO_ENABLED=0

-include .env

start:
	go run -v ./cmd/start

proxyclient:
	go run -v ./cmd/proxyclient -port $(port)
