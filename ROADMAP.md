# Project Roadmap

## Phase 0: Foundations & Project Setup ✅
- [x] Vision and project definition (`PROJECT.md`).
- [x] Initial repository structure.
- [x] Go module initialization (`go mod init`).
- [x] Git configuration and initial project cleanup.

## Phase 1: Basic Structure & Refactoring ✅
- [x] Unification of CLI under one binary.
- [x] Establishment of `internal/` package structure.
- [x] Centralization of `system-definition.yaml` parsing.
- [x] Rename project to `workspace-controller`.

## Phase 2: Materialization (`sync`) ✅
- [x] Automatic creation of `workspace/` and `infrastructure/` directories.
- [x] Implementation of Git-cloning logic for services.
- [x] Logic to handle versions/tags (checkout).
- [x] Security check: Skip if Git URL is a placeholder (e.g., `git@company/...`).

## Phase 3: Validation & Consistency (`validate`) ✅
- [x] Check if local files match the definition.
- [x] Verification of Git versions on disk.
- [x] Identification of "dirty" repositories (local changes).

## Phase 4: Infrastructure & Orchestration (`start` v2) ✅
- [x] Automatic generation of `docker-compose.yaml` for infrastructure (Postgres, etc.).
- [x] Support for environment variables configuration.
- [x] Integration of local services into Docker Compose (running apps from `workspace/`).
- [x] Implementation of `up` / `down` life-cycle commands.

## Phase 5: Advanced Features & Polish ✅
- [x] Documentation and CLI polish (help commands).
- [ ] Health checks for synchronized services.
- [ ] Dependency management between services and infrastructure.
- [ ] Plugin/Hook system for custom setup steps (e.g. database migrations).

## Phase 6: Ecosystem & Developer Experience 🚀
- [ ] Interactive CLI mode (menus for service selection).
- [ ] Status dashboard (view logs and health of all services).
- [ ] Template system for new services.
- [ ] Integration with common CI/CD pipelines for validation.
- [ ] Multi-platform support verification (Windows, macOS, Linux).

---
*Status and usage are also tracked in [README.md](README.md). Keep both documents in sync.*
