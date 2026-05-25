# AI Agent Guidelines

This document contains guidelines and best practices for AI agents working on the **Workspace Controller** project.

## 1. Interaction Rules

- **Language:** All code (variables, functions, comments), commit messages, and documentation MUST be in **English**.
- **No Git Commands:** AI agents are strictly forbidden from executing `git` commands (e.g., `git init`, `git add`, `git commit`, `git tag`). These operations must be performed or requested by the human user.
- **Minimal Changes:** Follow the "Diff-Only Change Policy". Prefer minimal, focused changes over broad refactors unless explicitly requested.
- **Explicit Logic:** Avoid over-abstraction or hidden optimizations. Keep the system explicit and readable.
- **Branch Management:** Always ensure that repositories remain on the `main` branch or a specific `feature` branch. Do NOT leave repositories in a "detached HEAD" state on a tag, even if the version matches. Versioning should be used for validation, but the working directory should stay on a branch.

## 2. Project Architecture Principles

- **Declarative First:** The system state is defined by configuration (`system-definition.yaml`), not by imperative scripts.
- **Deterministic Execution:** Given the same configuration, the result must always be identical across different machines.
- **Command-Based Interface:** All interactions happen through explicit commands implemented in the controller.
- **Internal Boundaries:** Logic shared by multiple commands must reside in the `internal/` package.

## 3. Directory Structure

Agents should respect the following structure:
- `workspace/`: The container for the local environment.
  - `services/`: Source code for services.
  - `infrastructure/`: Version-controlled infrastructure config.
  - `tools/`: Version-controlled development tools.
- `workspace-controller/`: The orchestration tool itself.

## 4. Documentation

- Ensure `README.md` is always synchronized with the current system capabilities.
- Update `ROADMAP.md` as new phases are reached or planned.
- Reference this `guidelines.md` in handover documents for future AI sessions.
