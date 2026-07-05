# Controller Development Makefile

.PHONY: help test coverage version build

help:
	@echo "Workspace Controller Logic - Development commands:"
	@echo "  make test      - Run automated tests"
	@echo "  make coverage  - Generate test coverage report (HTML)"
	@echo "  make version   - Show version information"
	@echo "  make build     - Build the controller binary"

test:
	@echo "==> Running tests..."
	go test -v ./...

coverage:
	@echo "==> Generating coverage reports..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "==> HTML report generated: coverage.html"

version:
	@echo "==> Workspace Controller version:"
	go run main.go version

build:
	@echo "==> Building binary..."
	go build -o ../bin/workspace-controller.exe main.go
