SHELL := /bin/bash

APP ?= api
BIN := ./bin

.PHONY: all
all: tools fmt lint test build

.PHONY: tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: fmt
fmt:
	go fmt ./...
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: build
build:
	mkdir -p $(BIN)
	go build -trimpath -ldflags="-s -w -X main.version=$$(git rev-parse --short HEAD 2>/dev/null || echo dev)" \
		-o $(BIN)/$(APP) ./cmd/$(APP)

.PHONY: run
run: build
	$(BIN)/$(APP)

.PHONY: clean
clean:
	rm -rf $(BIN)

.PHONY: ci
ci: all
