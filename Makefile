BINARY := fiken
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"
GOPATH := $(shell go env GOPATH)

.PHONY: build install test clean fmt lint

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

test:
	go test ./...

clean:
	rm -f $(BINARY)
	go clean

fmt:
	go fmt ./...

lint:
	golangci-lint run

.DEFAULT_GOAL := build
