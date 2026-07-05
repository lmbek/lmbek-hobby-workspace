# Workspace Controller Makefile
# Thin wrapper around git-controller commands.

.PHONY: help init clone pull push checkout status validate doctor ssh ssh-setup version v

export WORKSPACE_ROOT ?= $(abspath .)

help:
	@echo "Workspace Controller - Available commands:"
	@echo ""
	@echo "  Workflow:"
	@echo "  make init      - Scaffold a new workspace"
	@echo "  make clone     - Clone all repositories (initial setup)"
	@echo "  make pull      - Pull updates across all repositories (clone if missing)"
	@echo "  make push      - Push local commits across all repositories"
	@echo "  make checkout  - Switch all repos to their defined branch"
	@echo "  make status    - Show dashboard overview of all repository states"
	@echo "  make validate  - Validate repository consistency"
	@echo ""
	@echo "  Setup:"
	@echo "  make doctor    - Diagnose environment (Git, Go, SSH, Docker)"
	@echo "  make ssh       - Interactive SSH setup (alias: ssh-setup)"
	@echo "  make version   - Show version (alias: v)"

init:
	cd git-controller && go run . init

clone:
	cd git-controller && go run . clone

pull:
	cd git-controller && go run . pull

push:
	cd git-controller && go run . push

checkout:
	cd git-controller && go run . checkout

status:
	cd git-controller && go run . status

validate:
	cd git-controller && go run . validate

doctor:
	cd git-controller && go run . doctor

ssh:
	cd git-controller && go run . ssh

ssh-setup: ssh

version:
	cd git-controller && go run . version

v: version
