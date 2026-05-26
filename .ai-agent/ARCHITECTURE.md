# Workspace Controller Architecture

## Overview

This repository is structured as a command-based workspace controller for local development environments.

## Structure

### commands/
User-facing operations. Each folder represents a command that can be executed.

- init → initializes workspace: performs planning/analysis and materializes system locally (clones/updates repos)
- validate → checks system consistency and health (static checks + warnings for version mismatch)
- up → starts the system (docker-compose up)
- down → stops the system (docker-compose down)
- doctor → diagnoses environmental issues
- ssh-setup → manages SSH keys interactively
- help → shows usage instructions

### internal/
Shared system logic used by all commands.

This folder is NOT a utility folder.

It is an architecture boundary that contains:
- system parsing
- runtime logic
- git orchestration logic
- validation rules

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
- `go run main.go validate`
- `go run main.go up`
- `go run main.go help`

Or via the compiled binary: `./bin/workspace-controller [command]`

*Note: For a general overview and getting started guide, see [../README.md](../README.md). Always keep the README updated.*