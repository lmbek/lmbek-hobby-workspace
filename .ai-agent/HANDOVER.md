# Workspace Controller - Handover Document

## Context

This repository is a locally executable workspace controller designed to manage a multi-service development environment.

The goal is to provide a deterministic, reproducible way to:

- Describe a full system (services + infrastructure)
- Materialize it locally
- Validate its correctness
- Control its lifecycle through explicit commands

---

## Current Architecture

The system is structured around three main concepts:

### 1. system-definition.yaml (Source of Truth)

This file defines the entire system:

- Services (repositories + versions)
- Infrastructure dependencies
- System version

It is declarative and does NOT execute anything.

---

### 2. commands/ (User Interface Layer)

Each command is implemented in its own package:

#### init/
- Performs pre-flight checks (SSH agent status, GitHub connectivity).
- Automatically adds `github.com` to `known_hosts` on Windows.
- Reads `system-definition.yaml` and produces an execution plan.
- Guides user to run `sync` as the next step.

#### sync/
- Materializes system locally (clones/updates repos).
- Centralizes Git operations via `internal/gitutil`.
- Forces consistent SSH usage via `GIT_SSH_COMMAND`.
- Runs post-sync hooks.

#### validate/
- Validates system consistency (directory existence).
- Checks branch matches (defaulting to `main`).
- Performs real-time Health Checks (HTTP/TCP) on exposed host ports.
- Warns about local changes in infrastructure without blocking.

#### doctor/
- Diagnoses environmental issues (Git, SSH, Docker).
- Provides setup guidance for SSH keys and agents.
- Tests connectivity to GitHub with verbose diagnostic options.
- Detects structural SSH issues (e.g., identity file as directory).

#### ssh-setup / ssh
- Interactive tool for key generation and management.
- Configures `~/.ssh/config` with consolidation logic.
- (Windows only) Configures `git core.sshCommand` to use System OpenSSH.
- Automated cleanup and agent management.
- `ssh` is an alias for `ssh-setup`.

---

### 3. internal/ (Shared System Engine)
This folder contains shared logic used by all commands.

It is NOT a utility folder.

It contains:

- **system**: parsing logic, system model definitions, and CLI notes.
- **gitutil**: Git execution, environment isolation, and error handling.
- **sshutil**: Platform-specific SSH agent control, connectivity testing, and identity discovery.
- **validation**: Static check and health check (HTTP/TCP) implementation.

This enforces a strict boundary:
commands must not duplicate logic.

---

## Design Principles

### 1. Declarative First
The system is defined by configuration, not scripts.

### 2. Deterministic Execution
Given the same `system-definition.yaml`, the result must always be identical.

### 3. Command-Based Interface
All interactions happen through explicit, numbered commands:
- `[1] init`: Bootstrap environment.
- `[2] sync`: Materialize repositories.
- `[3] validate`: Check health and consistency.
- `[4] up`: Start containers.
- `[D] doctor`: Environmental diagnostics.
- `[S] ssh-setup`: Interactive SSH management.

### 4. Cross-Platform First
Code must handle Windows paths (slashes) and shell differences (PowerShell vs Sh) automatically.

---

## Installation & Build
To build the project, use:
```bash
go build -o bin/workspace-controller.exe main.go
```
The binary is ignored by Git. `go run main.go [command]` is the preferred way to execute during development. The `Makefile` provides a convenient wrapper.

---

## Current State of Implementation

Completed:

- Separated `init` and `sync` for granular control.
- Centralized `gitutil` and `sshutil` for robust orchestration.
- Fixed Windows health check connectivity by exposing ports in Compose.
- Automated SSH agent management and `known_hosts` setup.
- Comprehensive `doctor` and `ssh-setup` tools for troubleshooting.
- Cleaned up CLI output and standardized recommended actions.
- Full Windows/Linux/WSL compatibility verified.

---

## High-Level Goal

To provide a "one-click" setup experience for distributed local environments where all authentication and networking complexity is abstracted away by the controller.

---

## Key Constraints

1. **Language**: All code and documentation must be in **English**.
2. **No Git Commands**: AI agents are NOT allowed to suggest or execute Git commands.
3. **Cross-Platform Compatibility**: Use `filepath.ToSlash()` when passing absolute paths to Git on Windows.
4. **No Refactor Without Request**: Do not refactor unless explicitly asked.
5. **Diff-Only Change Policy**: Prefer minimal changes.
6. **Avoid Over-Abstraction**: Keep the system explicit and readable.
7. **AI Guidelines**: Refer to [guidelines.md](guidelines.md) for detailed AI interaction rules.

---
**Important**: The [../README.md](../README.md) serves as the primary entry point for new users and developers. Ensure it remains synchronized with the system's capabilities.
