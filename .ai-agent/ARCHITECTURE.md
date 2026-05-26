# Workspace Controller Architecture

## Overview

This repository is structured as a command-based workspace controller for local development environments.

## Structure

### commands/
User-facing operations, separated into core workspace lifecycle and supporting tools.

#### Core Workflow:
- init → bootstraps workspace: performs planning/analysis and pre-flight checks (SSH, agent)
- sync → materializes system locally: clones/updates repos and runs hooks
- validate → checks system consistency and health (static checks + warnings for version mismatch)
- up → starts the system (docker-compose up)
- down → stops the system (docker-compose down)

#### Tooling & Diagnostics (Integrated CLI Commands):
- doctor → [D] diagnoses environmental issues
- test → [T] runs automated tests (make workspace-controller-test)
- coverage → [C] generates HTML coverage report
- ssh-setup / ssh → [S] manages SSH keys interactively and configures Git SSH
- help → shows usage instructions

Note: The controller prioritizes safety and will not automatically rename or delete user configuration files. Manual intervention is requested for critical structural issues.

### internal/
Shared system logic used by all commands.

This folder is NOT a utility folder.

It is an architecture boundary that contains:
- system: parsing, models, and CLI helpers
- gitutil: centralized Git orchestration, error handling, and environment isolation (GIT_SSH_COMMAND)
- sshutil: platform-aware SSH agent management and connectivity diagnostics
- validation: rules and health check logic (HTTP/TCP)

## Why "internal"?

This is a Go language convention.

It enforces that:
- shared logic cannot be imported from outside this repository
- commands must use internal APIs instead of duplicating logic

## Rule

If logic is reused by more than one command, it belongs in internal/.

---

### Command Interface

All interactions happen through the `main.go` entrypoint using Go directly:

- `go run main.go init`
- `go run main.go sync`
- `go run main.go validate`
- `go run main.go up`
- `go run main.go help`

Or via the compiled binary: `./bin/workspace-controller [command]`
Or via the provided `Makefile`: `make sync`

*Note: For a general overview and getting started guide, see [../README.md](../README.md). Always keep the README updated.*