SHELL := /bin/bash

APP     ?= api
BIN     := ./bin
VERSION ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)

# Keep in sync with .github/workflows/ci.yml.
GOLANGCI_LINT_VERSION ?= v2.5.0
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: all
all: fmt lint test build

.PHONY: fmt
fmt:
	go fmt ./...
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run ./...

.PHONY: test
test:
	go test -race -shuffle=on ./...

.PHONY: build
build:
	mkdir -p $(BIN)
	go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" \
		-o $(BIN)/$(APP) ./cmd/$(APP)

.PHONY: run
run: build
	$(BIN)/$(APP)

.PHONY: clean
clean:
	rm -rf $(BIN) cover.out

.PHONY: ci
ci: all
