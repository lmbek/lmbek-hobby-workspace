# Controller Development Makefile

.PHONY: help test coverage

help:
	@echo "Workspace Controller Logic - Development commands:"
	@echo "  make test      - Run automated tests"
	@echo "  make coverage  - Generate test coverage report (HTML)"

test:
	@echo "==> Running tests..."
	go test -v ./...

coverage:
	@echo "==> Generating coverage reports..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "==> HTML report generated: coverage.html"

