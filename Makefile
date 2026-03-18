.PHONY: generate lint fmt build test

generate:
	buf dep update
	buf generate
	go generate ./...

lint:
	golangci-lint run ./...

fmt:
	golangci-lint run --fix ./...

build:
	go build ./...

test:
	go test ./...