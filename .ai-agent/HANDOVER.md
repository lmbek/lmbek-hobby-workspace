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
- Reads system-definition.yaml
- Produces an execution plan
- Materializes system locally (clones/updates repos)
- Runs post-sync hooks

#### validate/
- Validates system consistency
- Checks for missing or invalid configuration
- Verifies local Git state (dirty check)
- Performs real-time Health Checks (HTTP/TCP)
- Issues warnings for version mismatches (branch/tag vs definition) without blocking
- Ensures system integrity

#### doctor/
- Diagnoses environmental issues (Git, SSH, Docker)
- Provides setup guidance for SSH keys and agents
- Tests connectivity to GitHub

---

### 3. internal/ (Shared System Engine)

This folder contains shared logic used by all commands.

It is NOT a utility folder.

It contains:

- YAML parsing logic
- System model definitions
- Git orchestration logic
- Validation helpers
- Health Check logic

This enforces a strict boundary:
commands must not duplicate logic.

---

## Design Principles

### 1. Declarative First
The system is defined by configuration, not scripts.

### 2. Deterministic Execution
Given the same system-definition.yaml, the result must always be identical.

### 3. Command-Based Interface
All interactions happen through explicit commands:

- ./bin/workspace-controller init (OR go run main.go init)
- ./bin/workspace-controller validate (OR go run main.go validate)
- ./bin/workspace-controller up (OR go run main.go up)
- ./bin/workspace-controller down (OR go run main.go down)
- ./bin/workspace-controller doctor (OR go run main.go doctor)

### 4. No Hidden Behavior
No implicit scripts or side effects outside commands.

---

## Installation & Build
To build the project, use:
```bash
go build -o bin/workspace-controller.exe main.go
```
The binary is ignored by Git and should always reside in the `bin/` directory. However, for development and demo purposes, `go run main.go [command]` is the preferred way to execute the controller.

---

## Current State of Implementation

Completed:

- system-definition.yaml created (supporting Services, Infrastructure, and Tools)
- init command implemented (unified planning + materialization)
- validate command implemented (Git state + health checks)
- up/down commands implemented (Docker lifecycle orchestration)
- doctor command implemented (environmental diagnostics)
- Unified CLI entrypoint established
- Renamed project to workspace-controller
- Workspace reorganization (centralized services directory)
- Decoupled infrastructure and tools into specialized repositories
- Environment variable support in system definition (e.g., `${WORKSPACE_ROOT}`)

---

## High-Level Goal

To evolve this system into a fully reproducible local development environment where:

- Any developer can clone one repo
- Run a single command (`init`)
- Get a fully running distributed system locally (`up`)
- With identical behavior across machines

---

## Key Constraints

1. **Language**: All code (variables, functions, comments) and documentation must be in **English**.
2. **No Git Commands**: AI agents are NOT allowed to suggest or execute Git commands.
3. **No Refactor Without Request**: Do not refactor unless explicitly asked.
4. **No Hidden Optimization**: Keep logic explicit.
5. **Diff-Only Change Policy**: Prefer minimal changes.
6. **Avoid Over-Abstraction**: Keep the system explicit and readable.
7. **AI Guidelines**: Refer to [guidelines.md](guidelines.md) for detailed AI interaction rules.

---
**Important**: The [../README.md](../README.md) serves as the primary entry point for new users and developers. Ensure it remains synchronized with the system's capabilities.
