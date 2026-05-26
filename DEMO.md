# Workspace Controller - Demo Guide

This guide will walk you through the core features of the **Workspace Controller**. We have moved from a monolithic setup to a flexible, decoupled architecture where services, infrastructure, and tools are independent, version-controlled entities.

> **Note:** In the examples below, `<cli>` refers to the command you are using (e.g., `./bin/workspace-controller.exe`, `go run main.go`, or `make`).

### Quick Start with Makefile
For a streamlined experience, use the included `Makefile`:
```bash
make init      # Step 1
make validate  # Step 2
make up        # Step 3
```

## 1. Initialization (`init`)
The `init` command is your primary entry point for setting up the workspace. It performs both **planning/analysis** and **materialization** (syncing repositories).

**Run:**
```bash
# Set the root path for your workspace
$env:WORKSPACE_ROOT=".." 
<cli> init
```

**What to look for:**
- **Pre-flight checks:** Automated verification of GitHub connectivity via SSH.
- **Planning & Analysis:** A detailed list of services and infrastructure derived from `system-definition.yaml`.
- **Materializing Workspace:** The controller clones missing repositories and updates existing ones.
- **Hooks in action:** Execution of the `post-sync` hook (e.g., a success message).

## 2. Validating Consistency (`validate`)
The `validate` command performs consistency checks across all components.

**Run:**
```bash
$env:WORKSPACE_ROOT=".."
<cli> validate
```

**What to look for:**
- Verification that all Git components match the defined versions.
- Detection of local "dirty" changes.
- **Health Checks:** Real-time monitoring of services via HTTP/TCP.

## 3. Lifecycle Management (`up` / `down`)
Control your local environment via Docker Compose through the controller.

**Try starting:**
```bash
$env:WORKSPACE_ROOT=".."
<cli> up
```

**What to look for:**
- The controller executes `docker-compose up -d` within the infrastructure directory.
- **Hooks in action:** Observe the execution of the `post-up` hook.

## 4. Diagnostics (`doctor` / `ssh-setup`)
If you encounter environmental issues (e.g., SSH permission denied), use the built-in diagnostic tools.

**Run:**
```bash
<cli> doctor
# OR for interactive SSH management
<cli> ssh-setup
```

---

## Current Status
`workspace-controller` is now a mature orchestration tool:

- **Simplified Entrypoint:** The `init` command combines planning and materialization into a single action.
- **Decoupled Architecture:** Infrastructure and tools are independent, version-controlled entities.
- **Environment Variables:** Full support for dynamic paths via environment variables (e.g., `${WORKSPACE_ROOT}`) in the system definition.
- **Automated Lifecycle:** From initialization to health monitoring and shutdown.
- **Native Go Execution:** Optimized for `<cli>` for seamless development experience.
