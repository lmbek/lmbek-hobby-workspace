# Workspace Controller Makefile
# Mirrors the steps in DEMO.md

.PHONY: help init sync validate up down doctor ssh-setup ssh

# Default workspace root if not set
export WORKSPACE_ROOT ?= ..

help:
	@echo "Workspace Controller - Available commands:"
	@echo "  make init      - [1] Bootstrap workspace (Pre-flight + Planning)"
	@echo "  make sync      - [2] Synchronize all repositories (Materialization + fetch/pull/hooks)"
	@echo "  make validate  - [3] Validate consistency and health"
	@echo "  make up        - [4] Start the system (docker-compose up)"
	@echo "  make down      - Stop the system (docker-compose down)"
	@echo "  make doctor    - [D] Diagnose environmental issues"
	@echo "  make ssh       - [S] Interactive SSH setup (alias for ssh-setup)"
	@echo "  make ssh-setup - [S] Interactive SSH setup (alias: ssh)"

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


ssh: ssh-setup

ssh-setup:
	@echo "==> Starting interactive SSH setup (ssh-setup)..."
	go run main.go ssh-setup
