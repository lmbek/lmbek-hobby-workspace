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

#### start/
- Reads system-definition.yaml
- Produces an execution plan
- Does NOT modify system state
- Acknowledges decoupled infrastructure repositories

#### sync/
- Materializes system locally
- Clones service repositories
- Creates workspace structure

#### validate/
- Validates system consistency
- Checks for missing or invalid configuration
- Verifies local Git state (for Services, Infrastructure, and Tools)
- Performs real-time Health Checks (HTTP/TCP)
- Ensures system integrity

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

- go run main.go start
- go run main.go sync
- go run main.go validate
- go run main.go up
- go run main.go down

### 4. No Hidden Behavior
No implicit scripts or side effects outside commands.

---

## Current State of Implementation

Completed:

- system-definition.yaml created (supporting Services, Infrastructure, and Tools)
- start command implemented (execution planning for decoupled architecture)
- sync command implemented (repository cloning + versioning for all components)
- validate command implemented (Git state + health checks)
- up/down commands implemented (Docker lifecycle orchestration)
- Unified CLI entrypoint established
- Renamed project to workspace-controller
- Workspace reorganization (centralized services directory)
- Decoupled infrastructure and tools into specialized repositories

---

## High-Level Goal

To evolve this system into a fully reproducible local development environment where:

- Any developer can clone one repo
- Run a single command (`sync`)
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
