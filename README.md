# LMBEK Hobby Workspace

A multi-repo workspace that uses a single CLI to manage all your Git repositories in one place. Instead of cloning and updating each repo manually, the **git-controller** reads a central config file (`system-definition.yaml`) and handles everything for you.

## Table of Contents

- [What Is This?](#what-is-this)
- [Architecture Overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Available Commands](#available-commands)
- [Project Structure](#project-structure)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [Security](#security)
- [Guidelines](#guidelines)

---

## What Is This?

This repository is a **workspace root**. It doesn't contain application code itself — instead, it contains:

1. A CLI tool (`git-controller/`) that clones, pulls, and pushes all your project repositories.
2. A definition file (`system-definition.yaml`) that lists which repos to manage and where they go.
3. A `git-repositories/` directory (gitignored) where all managed repos live on disk.

When you run the CLI, it reads the definition file and ensures every listed repository is cloned, up to date, and on the correct branch.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     Workspace Root                              │
│  system-definition.yaml ──► git-controller (CLI)                │
│                                  │                              │
│                    ┌─────────────┼─────────────┐                │
│                    ▼             ▼             ▼                │
│              applications   orchestrator   infrastructure       │
│              ┌──────────┐   ┌──────────┐   ┌──────────┐        │
│              │ service1  │   │ compose  │   │ terraform│        │
│              │ service2  │   │ manifests│   │ servers  │        │
│              └──────────┘   └──────────┘   └──────────┘        │
│                    │             │             │                │
│                    ▼             ▼             ▼                │
│              deployment    observability     tools              │
│              ┌──────────┐  ┌──────────┐   ┌──────────┐        │
│              │ dev/      │  │ grafana  │   │ scripts  │        │
│              │ staging/  │  │ prometheus│  │ utilities│        │
│              │ prod/     │  │ alerts   │   └──────────┘        │
│              └──────────┘  └──────────┘                        │
│                    │             │                              │
│                    ▼             ▼                              │
│                        docs                                    │
│              ┌─────────────────────────┐                       │
│              │ ADRs, runbooks          │                       │
│              └─────────────────────────┘                       │
└─────────────────────────────────────────────────────────────────┘
```

**Key principles:**

- **One repo per service** — each microservice has its own repository, Dockerfile, and CI/CD pipeline.
- **Orchestrator wires services** — uses pre-built images, never builds from source.
- **Deployment is environment config** — folder-per-environment (dev/staging/prod) in a single repo.
- **Observability is separate** — dashboards, alerts, and log pipelines live outside application code.
- **Everything is YAML-driven** — `system-definition.yaml` is the single source of truth.

---

## Prerequisites

Before you begin, make sure you have:

- **Git** — any recent version
- **Go 1.26+** — needed to build and run the git-controller CLI
- **SSH access** — your SSH key must be authorized for the repositories listed in `system-definition.yaml`

### Setting up Go on Windows (PowerShell)

If you have multiple Go versions installed, make sure the correct one is first in PATH:

```powershell
$env:GOTOOLCHAIN = "local"
$env:PATH = "C:\Users\larsz\sdk\go1.26.4\bin;" + $env:PATH
go version   # should print go1.26.4
```

---

## Getting Started

Open a terminal at the workspace root and run these commands in order:

### 1. Check your environment

```
make doctor
```

This checks that Git, Go, SSH, and Docker are available and correctly configured. Fix anything it reports before continuing.

### 2. Set up SSH (if needed)

```
make ssh
```

An interactive wizard that helps you generate an SSH key and add it to your Git host.

### 3. Clone all repositories

```
make clone
```

This reads `system-definition.yaml` and clones every listed repository into `git-repositories/`.

### 4. Keep everything up to date

```
make pull
```

Pulls the latest changes across all repositories. If any repo is missing, it clones it first.

### 5. Push local changes

```
make push
```

Pushes local commits across all cloned repositories.

### 6. Switch branches

```
make checkout
```

Switches all repositories to the branch defined in `system-definition.yaml`. Override with `BRANCH=feature-x make checkout`.

### 7. Verify everything is correct

```
make validate
```

Checks that every repository matches the definition (correct branch, clean state).

---

## Available Commands

All commands can be run via `make` from the workspace root, or directly with `go run .` from the `git-controller/` directory.

| Command          | Description                                              |
|------------------|----------------------------------------------------------|
| `make init`      | Scaffold a new workspace (definition, Makefile, .gitignore) |
| `make clone`     | Clone all repositories (first-time setup)                |
| `make pull`      | Pull updates across all repositories (clone if missing)  |
| `make push`      | Push local commits across all repositories               |
| `make checkout`  | Switch all repos to their defined branch                 |
| `make status`    | Show dashboard overview of all repository states         |
| `make validate`  | Verify repos match the system definition                 |
| `make doctor`    | Diagnose environment issues (Git, Go, SSH, Docker)       |
| `make ssh`       | Interactive SSH key setup wizard                         |

---

## Project Structure

```
LMBEK-HOBBY-WORKSPACE/
├── git-controller/          Go CLI that manages repositories
├── git-repositories/        All managed repos live here (gitignored)
│   ├── deployment/          Environment configs and manifests (folders per env)
│   ├── applications/        Microservices (one repo per service)
│   ├── orchestrator/        Docker Compose / K8s manifests to run the stack
│   ├── infrastructure/      Terraform / IaC and server provisioning
│   ├── observability/       Monitoring, dashboards, alerts, log pipelines
│   ├── tools/               General-purpose utilities and scripts
│   └── docs/                Architecture docs, API docs, manual runbooks
├── system-definition.yaml   Defines which repos/branches to manage
├── Makefile                 Convenience targets for the CLI
├── .env.example             Documented environment variables
├── .github/                 Issue templates, PR template, CI workflows
├── AGENTS.md               AI agent guidelines (defers to GUIDELINES.md)
├── GUIDELINES.md            Master guidelines for all projects
├── CONTRIBUTING.md          How to contribute (workflow, PR process, commits)
├── GIT-REPOS.md             Guide to all git repository categories and why they exist
├── SECURITY.md              Vulnerability reporting policy
└── README.md                You are here
```

### Enterprise Repository Layout

| Category         | Purpose                                                    | Example Repo                        |
|------------------|------------------------------------------------------------|-------------------------------------|
| `deployment`     | Per-environment values, secrets templates, and promotion config (folder-per-env) | `lmbek-hobby-deployment` |
| `applications`   | Independent microservices, each with its own Dockerfile and CI/CD | `lmbek-hobby-placeholder1-service`  |
| `orchestrator`   | Central Docker Compose / K8s manifests that wire services using pre-built images | `lmbek-hobby-orchestrator`          |
| `infrastructure` | Terraform modules for cloud resources and server provisioning (Ansible, Packer, cloud-init) | `lmbek-hobby-infrastructure`, `lmbek-hobby-servers` |
| `observability`  | Grafana dashboards, Prometheus rules, alert definitions, log pipelines | `lmbek-hobby-observability` |
| `tools`          | General-purpose utilities, scripts, and helper tooling     | `lmbek-hobby-tools`                 |
| `docs`           | Architecture Decision Records, API docs, manual runbooks | `lmbek-hobby-docs` |

After cloning, start all services via the **orchestrator** repo:

```
cd git-repositories/orchestrator && docker compose up -d
```

### Where Do Compose Files Live?

| Compose file | Location | Purpose |
|---|---|---|
| Per-service `docker-compose.yml` | Each application repo | Run that single service in isolation during development (with its own DB, cache, etc.) |
| Central `docker-compose.yml` | Orchestrator repo | Wire **all** services together using pre-built images — the full-stack local environment |
| Environment overrides | Deployment repo | Per-environment values (dev/staging/prod folders) consumed by CI/CD or GitOps tooling |

Each service repo owns its own `Dockerfile`. The orchestrator references images built from those Dockerfiles — it does **not** build from source.

---

## Troubleshooting

- **Files from `git-repositories/` showing in git status?**
  They were committed before the ignore rule existed. Untrack them once:
  ```
  git rm -r --cached git-repositories
  git commit -m "Stop tracking git-repositories"
  ```

- **Wrong Go version?**
  Place the correct Go bin directory first in PATH and set `GOTOOLCHAIN=local`. See [Setting up Go on Windows](#setting-up-go-on-windows-powershell).

- **SSH permission errors?**
  Run `make doctor` to see what's wrong, then `make ssh` to fix it.

---

## Contributing

We welcome contributions! Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for development workflow, PR process, and commit message conventions.

---

## Security

For reporting vulnerabilities, see [SECURITY.md](./SECURITY.md). Do not open public issues for security concerns.

---

## Guidelines

All projects in this workspace follow the master guidelines first, then any project-specific guidelines after.
See [GUIDELINES.md](./GUIDELINES.md).
