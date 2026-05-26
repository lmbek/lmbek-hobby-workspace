# Workspace Controller Makefile
# Mirrors the steps in DEMO.md

.PHONY: help init sync validate up down doctor ssh-setup ssh workspace-controller-test workspace-controller-help version

# Default workspace root if not set
export WORKSPACE_ROOT ?= ..

help:
	@echo "Workspace Controller - Available commands:"
	@echo "  Core Workflow:"
	@echo "    make init      - [1] Bootstrap workspace (Pre-flight + Planning)"
	@echo "    make sync      - [2] Synchronize all repositories (Materialization + fetch/pull/hooks)"
	@echo "    make validate  - [3] Validate consistency and health"
	@echo "    make up        - [4] Start the system (docker-compose up)"
	@echo "    make down      - Stop the system (docker-compose down)"
	@echo ""
	@echo "  Tooling & Diagnostics:"
	@echo "    make doctor    - [D] Diagnose environmental issues"
	@echo "    make workspace-controller-test - [T] Run automated tests"
	@echo "    make workspace-controller-coverage - [C] Generate test coverage report (HTML)"
	@echo "    make workspace-controller-help - Show CLI help"
	@echo "    make version   - [V] Show version information"
	@echo "    make ssh       - [S] Interactive SSH setup (alias for ssh-setup)"
	@echo "    make ssh-setup - [S] Interactive SSH setup (alias: ssh)"

init:
	@echo "==> Bootstrapping workspace (init)..."
	go run main.go init

sync:
	@echo "==> Synchronizing repositories (sync)..."
	go run main.go sync

validate:
	@echo "==> Validating system consistency (validate)..."
	go run main.go validate

up:
	@echo "==> Starting the system (up)..."
	go run main.go up

down:
	@echo "==> Stopping the system (down)..."
	go run main.go down

doctor:
	@echo "==> Running environment diagnostics (doctor)..."
	go run main.go doctor

workspace-controller-test:
	@echo "==> Running tests..."
	go test -v ./...

workspace-controller-help:
	@echo "==> Showing CLI help:"
	go run main.go help

workspace-controller-coverage:
	@echo "==> Generating coverage reports..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "==> HTML report generated: coverage.html"


ssh:
	@echo "==> Starting interactive SSH setup (ssh)..."
	go run main.go ssh

ssh-setup:
	@echo "==> Starting interactive SSH setup (ssh-setup)..."
	go run main.go ssh-setup

version:
	@echo "==> Workspace Controller version:"
	go run main.go version
