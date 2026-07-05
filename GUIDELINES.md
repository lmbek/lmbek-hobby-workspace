# Project MASTER Guidelines

This document serves as the MASTER guidelines for all repositories and projects within this workspace. All contributors, including human developers and AI agents, must follow these standards.

## Hierarchy of Guidelines
1.  **MASTER Guidelines** (This file, located at the workspace root)
2.  **Application/Project Guidelines** (Specific to sub-folders)

If a rule in a sub-project conflicts with the MASTER guidelines, the MASTER guidelines take precedence.

### Shared Guideline Files (must defer to MASTER)
To help humans and AI agents consistently discover the MASTER rules, the following common policy files exist at the workspace root and explicitly defer to this document:
- `AI_AGENTS.md`
- `AGENTS.md`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `SECURITY.md`

All of these files state that the MASTER guidelines in `GUIDELINES.md` are the first source of truth. Project-level policies come second.

#### Why these separate files exist
Some tools and AI agents prioritize certain filenames by default (for example, CONTRIBUTING.md). To ensure they all land on the same rules first, we keep those familiar filenames as thin pointers that explicitly defer to this MASTER document. This avoids duplication while improving discoverability. If you only want a single source of truth, read this file — the others simply link back here.

## Core Principles
*   **Decoupling**: Applications and infrastructure must be loosely coupled.
*   **Consistency**: Follow the established Go structure in `git-controller`.
*   **Automation**: Prioritize tools and scripts over manual configuration.
*   **Documentation**: All new features or changes must be documented in the relevant `README.md`.

## AI Agent Instructions
*   Always read `GUIDELINES.md` from the root before performing any task.
*   Respect the hierarchy of guidelines.
*   Do not hardcode configuration; use structured config objects (e.g., `WorkspaceConfig`).
*   When a repository contains any of the shared markdown files listed above, treat them as pointers back to this MASTER guideline first, then the local per-project docs after.

## Git Repositories
*   All external git repositories managed by the `git-controller` are placed in the `git-repositories/` directory (or specific category directories).
*   These directories must be ignored in the root `.gitignore`.
