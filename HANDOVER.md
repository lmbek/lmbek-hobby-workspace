# System Controller - Handover Document

## Context

This repository is a locally executable system controller designed to manage a multi-service development environment.

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

Each folder represents a command:

#### start/
- Reads system-definition.yaml
- Produces an execution plan
- Does NOT modify system state

#### sync/
- Materializes system locally
- Clones service repositories
- Creates workspace structure

#### validate/
- Validates system consistency
- Checks for missing or invalid configuration
- Ensures system integrity

---

### 3. internal/ (Shared System Engine)

This folder contains shared logic used by all commands.

It is NOT a utility folder.

It contains:

- YAML parsing logic
- System model definitions
- Git orchestration logic
- Validation rules

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

- system-controller start
- system-controller sync
- system-controller validate

### 4. No Hidden Behavior
No implicit scripts or side effects outside commands.

---

## Current State of Implementation

Completed:

- system-definition.yaml created
- start command implemented (execution planning)
- sync command implemented (repository cloning)
- validate command implemented (basic structure validation)

Partially completed:

- internal/ architecture introduced (shared logic layer not fully extracted yet)
- consistent CLI interface not yet unified (still folder-based entrypoints)

---

## Next Intended Step

The next architectural step is to:

1. Extract shared logic into internal/ properly
2. Remove duplicated YAML parsing across commands
3. Introduce consistent command entrypoint pattern
4. Improve deterministic execution guarantees
5. Strengthen validation rules

---

## High-Level Goal

To evolve this system into a fully reproducible local development environment where:

- Any developer can clone one repo
- Run a single command
- Get a fully running distributed system locally
- With identical behavior across machines

---

## Key Constraint

Avoid over-abstraction.

Keep the system:

- explicit
- readable
- command-driven
- deterministic

## Documentation

For detailed documentation on each command and internal architecture, refer to the respective files in the repository.



ALSO:
“no refactor without request”
“no hidden optimization”
“diff-only change policy”
