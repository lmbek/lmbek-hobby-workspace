# Workspace Controller

The **Workspace Controller** is a declarative tool designed to manage complex, distributed local development environments. It allows you to define your entire system—services, versions, and infrastructure—in a single configuration file (`system-definition.yaml`) and orchestrate everything with simple commands.

## Vision
Move away from brittle bash scripts and manual setup. Use a single source of truth to materialize, validate, and run your local workspace.

## Key Features
- **Declarative Definition:** System state is described in YAML.
- **Automated Sync:** Clones and manages Git repositories automatically.
- **Infrastructure as Code:** Manages version-controlled infrastructure repositories.
- **Dependency Management:** Define service dependencies.
- **Hook System:** Run custom commands after sync or startup.
- **Consistency Checks:** Validates that your local environment matches the definition.
- **Health Monitoring:** Built-in support for HTTP and TCP health checks.

---

## Getting Started

### 1. Prerequisites
- [Go](https://golang.org/doc/install) (1.26 or later)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

### 2. Using the Makefile (Recommended)
You can use the provided `Makefile` for a convenient experience.

#### Core Workflow
```bash
make init      # [1] Bootstrap workspace
make sync      # [2] Synchronize all repositories
make validate  # [3] Validate consistency and health
make up        # [4] Start the system
make down      # Stop the system
```

#### Tooling & Diagnostics
```bash
make doctor    # [D] Diagnose environment issues
```

## Structure
The project is organized to keep the root clean:
- `git-controller/`: The Go source code for the orchestrator.
- `system-definition.yaml`: The manifest defining all managed repositories.
- `Makefile`: Root-level entry point.

---

## Documentation & Policies
- Start with the MASTER guidelines: [GUIDELINES.md](./GUIDELINES.md)
- Shared policy files at the workspace root — [AI_AGENTS.md](./AI_AGENTS.md), [AGENTS.md](./AGENTS.md), [CONTRIBUTING.md](./CONTRIBUTING.md), [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md), [SECURITY.md](./SECURITY.md) — all defer to the MASTER guidelines first, then to any application/project-specific docs.
- Sub-projects may include their own README and guidelines; always apply them after the MASTER rules.
