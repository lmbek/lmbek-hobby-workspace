# Project Roadmap

## Phase 1: Foundation & Architecture ✅
- [x] Initial project setup and vision definition.
- [x] Establishment of the `internal/` shared logic structure.
- [x] Unified CLI entrypoint (`main.go`) to replace scattered scripts.
- [x] Core data model for system definition (YAML parsing).

## Phase 2: Materialization & Sync ✅
- [x] Automated workspace directory management.
- [x] Git orchestration for cloning and versioning services.
- [x] Version control for infrastructure and tools repositories.
- [x] Security: Protection against placeholder URLs.

## Phase 3: Validation & Health ✅
- [x] Static validation: directory existence and Git state (dirty check, version match).
- [x] Runtime validation: Integrated HTTP and TCP health checks.
- [x] Consistent reporting of environment integrity.

## Phase 4: Lifecycle & Orchestration ✅
- [x] Implementation of `up` and `down` commands for environment control.
- [x] Support for environment variables in system definitions.
- [x] Hook system for post-sync and post-startup automation.
- [x] Dependency management between services (startup order).

## Phase 5: Decoupling & Specialization ✅
- [x] Infrastructure decoupling: Move configuration to specialized repositories.
- [x] Reorganization: Centralized `services/` directory within the workspace.
- [x] Introduction of the `tools/` ecosystem (e.g., Terraform-based `deployer`).
- [x] Modernized demo and documentation (Phase 5 Alignment).

## Phase 6: Reliability & Cross-Platform ✅
- [x] Windows/Linux/WSL compatibility for SSH and Git operations.
- [x] Robust SSH agent management (Start-Service on Windows, socket discovery on Linux).
- [x] Automated `known_hosts` management for GitHub.
- [x] Synchronization of SSH clients between pre-flight checks and Git operations (GIT_SSH_COMMAND).
- [x] Interactive `ssh-setup` tool for key management and consolidation.

## Phase 7: CLI Refinement & Optimization (Current State) ✅
- [x] Separation of `init` (bootstrap) and `sync` (materialize) commands.
- [x] Centralization of common logic into `gitutil` and `sshutil` packages.
- [x] Standardized CLI numbering [1] [2] [3] [4] [D] [S] for intuitive workflow.
- [x] Improved error handling and actionable diagnostics in `doctor` and `HandleGitError`.

## Phase 8: Quality Assurance & CI (Current State) *
- [x] Integrated unit tests for core logic (system, sshutil, validate).
- [x] Configured GitHub Actions CI for multi-platform build and test (Windows/Linux).
- [x] Added `make test` command for local verification.
- [x] Added `make workspace-controller-coverage` for HTML reports.

---
*This roadmap reflects the actual journey taken to build the workspace-controller. All phases are currently completed and verified.*
