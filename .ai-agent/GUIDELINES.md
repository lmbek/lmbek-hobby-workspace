# AI Agent Guidelines

This document contains guidelines and best practices for AI agents working on the **Workspace Controller** project.

## 1. Interaction Rules

- **Language:** All code (variables, functions, comments), commit messages, and documentation MUST be in **English**.
- **No Git Commands:** AI agents are strictly forbidden from executing `git` commands (e.g., `git init`, `git add`, `git commit`, `git tag`). These operations must be performed or requested by the human user.
- **Minimal Changes:** Follow the "Diff-Only Change Policy". Prefer minimal, focused changes over broad refactors unless explicitly requested.
- **Explicit Logic:** Avoid over-abstraction or hidden optimizations. Keep the system explicit and readable.
- **Cross-Platform Paths:** When executing Git commands on Windows, absolute paths in environment variables (like `GIT_SSH_COMMAND`) must use forward slashes (`/`) and be quoted if they contain spaces. Use `filepath.ToSlash()` for consistency.
- **SSH Mandatory:** This project strictly uses SSH for all Git operations. HTTPS URLs are forbidden in `repos.yaml` unless they are for internal/local paths. The controller must enforce SSH connectivity during `init`.
- **Branch Management:** The controller MUST ensure that repositories remain on the `main` branch (or a specific `feature` branch) and never enter a "detached HEAD" state. The `sync` command is configured to use `git fetch` and `git pull` on the current branch rather than `git checkout <tag>`. Version validation in `validate` will issue a warning if the branch/tag doesn't match the definition, but it will not block execution. AI agents must also respect this and never suggest commands that lead to a detached HEAD.

## 2. Project Architecture Principles

- **Declarative First:** The system state is defined by configuration (`repos.yaml`), not by imperative scripts.
- **Deterministic Execution:** Given the same configuration, the result must always be identical across different machines.
- **Command-Based Interface:** All interactions happen through explicit commands implemented in the controller.
- **Internal Boundaries:** Logic shared by multiple commands must reside in the `internal/` package.

## 3. Directory Structure

Agents should respect the following structure:
- `applications/`: Customer/business services.
- `proxy/`: Central workspace dashboard and reverse proxy.
- `orchestrator/`: Dedicated orchestration repositories.
- `infrastructure/`: Servers, networking, cloud, and Docker configurations.
- `platform/`: Developer tooling and observability.
- `tools/`: Local helper scripts and CLIs.
- `docs/`: Architecture and onboarding documentation.
- `repos.yaml`: The manifest defining all managed repositories.
- `controller/`: The Go source code for the **Workspace Controller** orchestrator.

## 4. Documentation

- Ensure `README.md` is always synchronized with the current system capabilities.
- Update `ROADMAP.md` as new phases are reached or planned.
- Reference this `guidelines.md` in handover documents for future AI sessions.
