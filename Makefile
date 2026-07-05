# Workspace Controller Makefile
# Mirrors the steps in DEMO.md

.PHONY: help init sync validate up down doctor ssh ssh-setup version v

# Default workspace root if not set
export WORKSPACE_ROOT ?= $(abspath .)

help:
	@echo "Workspace Controller - Available commands:"
	@echo "  make init      - [1] Bootstrap workspace (Pre-flight + Planning)"
	@echo "  make sync      - [2] Synchronize all repositories"
	@echo "  make validate  - [3] Validate consistency and health"
	@echo "  make up        - [4] Start the system"
	@echo "  make down      - Stop the system"
	@echo "  make doctor    - Diagnose environmental issues"
	@echo "  make ssh       - Interactive SSH setup (alias: ssh-setup)"
	@echo "  make version   - Show version info (alias: v)"

init:
	@echo "==> Bootstrapping workspace (init)..."
	cd controller && go run main.go init

sync:
	@echo "==> Synchronizing repositories (sync)..."
	cd controller && go run main.go sync

validate:
	@echo "==> Validating system consistency (validate)..."
	cd controller && go run main.go validate

up:
	@echo "==> Starting the system (up)..."
	cd controller && go run main.go up

down:
	@echo "==> Stopping the system (down)..."
	cd controller && go run main.go down

doctor:
	@echo "==> Running environment diagnostics (doctor)..."
	cd controller && go run main.go doctor

ssh:
	@echo "==> Running SSH setup tool..."
	cd controller && go run main.go ssh

ssh-setup: ssh

version:
	@echo "==> Workspace Controller version:"
	cd controller && go run main.go version

v: version
