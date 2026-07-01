.PHONY: build test run fmt lint clean help

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  install  - Install the binary to GOBIN"
	@echo "  test     - Run tests"
	@echo "  run      - Run with go run"
	@echo "  fmt      - Format code"
	@echo "  lint     - Run linter"
	@echo "  clean    - Remove build artifacts"
	@echo "  coverage - Run tests with coverage report"

build:
	go build

install:
	go install

test:
	go test ./...
	terraform validate --test-directory provisioning/

run:
	go run main.go

fmt:
	go fmt ./...
	terraform fmt --recursive provisioning/

lint:
	golangci-lint run ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/ coverage.out
