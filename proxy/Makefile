# Proxy Service Makefile

.PHONY: test build run

test:
	@echo "==> Running tests for proxy-service..."
	go test -v ./...

build:
	@echo "==> Building proxy-service..."
	go build -o proxy-service main.go

run:
	@echo "==> Running proxy-service..."
	go run main.go
