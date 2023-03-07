-include .env

start:
	go run -v ./cmd/start

proxyclient:
	go run -v ./cmd/proxyclient -port $(port)
