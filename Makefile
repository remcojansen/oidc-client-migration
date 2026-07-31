.PHONY: build install test run fmt lint clean help

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
	go build -o bin/ocm ./export

install:
	@GOBIN_DIR="$$(go env GOBIN)"; \
	if [ -z "$$GOBIN_DIR" ]; then GOBIN_DIR="$$(go env GOPATH)/bin"; fi; \
	go build -o "$$GOBIN_DIR/ocm" ./export

test:
	go test ./...
	terraform -chdir=import validate

run:
	go run ./export

fmt:
	go fmt ./...
	terraform fmt -recursive import/

lint:
	golangci-lint run ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/ coverage.out
