# Workspace Controller Makefile
# Mirrors the steps in DEMO.md

.PHONY: help init validate up down doctor ssh-setup ssh

# Default workspace root if not set
export WORKSPACE_ROOT ?= ..

help:
	@echo "Workspace Controller - Available commands:"
	@echo "  make init      - [1] Initialize workspace (Pre-flight + Planning + Materialization)"
	@echo "  make validate  - [2] Validate consistency and health"
	@echo "  make up        - [3] Start the system (docker-compose up)"
	@echo "  make down      - Stop the system (docker-compose down)"
	@echo "  make doctor    - [4] Diagnose environmental issues"
	@echo "  make ssh       - [4] Interactive SSH setup (alias for ssh-setup)"
	@echo "  make ssh-setup - [4] Interactive SSH setup"

init:
	@echo "==> Running initialization (init)..."
	go run main.go init

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
