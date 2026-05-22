.PHONY: build test lint install clean

build:
	go build -o bin/doorloop-pp-cli ./cmd/doorloop-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/doorloop-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/doorloop-pp-mcp ./cmd/doorloop-pp-mcp

install-mcp:
	go install ./cmd/doorloop-pp-mcp

build-all: build build-mcp
