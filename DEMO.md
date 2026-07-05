# Workspace Controller - Demo Guide

This guide will walk you through the core features of the **Workspace Controller**. We have moved from a monolithic setup to a flexible, decoupled architecture where services, infrastructure, and tools are independent, version-controlled entities.

> **Note:** In the examples below, `<cli>` refers to the command you are using (e.g., `./bin/workspace-controller.exe`, `go run main.go`, or `make`).

### Quick Start with Makefile
For a streamlined experience, use the included `Makefile`. **Note:** This requires a working SSH setup for GitHub.

```bash
make ssh-setup # Only if you haven't configured SSH yet
make init      # Step 1
make validate  # Step 2
make up        # Step 3
```

## 1. Initialization (`init`)
The `init` command is your primary entry point for setting up the workspace. It performs both **planning/analysis** and **materialization** (syncing repositories). It also enforces SSH connectivity.

**Run:**
```bash
# Set the root path for your workspace
$env:WORKSPACE_ROOT="." 
<cli> init
```

**What to look for:**
- **Pre-flight checks:** Automated verification of GitHub connectivity via SSH. If this fails, the controller will guide you through the `ssh-setup` process.
- **Planning & Analysis:** A detailed list of services and infrastructure derived from `repos.yaml`.
- **Materializing Workspace:** The controller clones missing repositories and updates existing ones using SSH.
- **Hooks in action:** Execution of the `post-sync` hook (e.g., a success message).

## 2. Validating Consistency (`validate`)
The `validate` command performs consistency checks across all components.

**Run:**
```bash
$env:WORKSPACE_ROOT="."
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
$env:WORKSPACE_ROOT="."
<cli> up
```

**What to look for:**
- The controller executes `docker-compose up -d` within the orchestrator directory.
- **Proxy & DNS:** Once started, access the central dashboard at `http://localhost`. Services are also available via `http://service1.localhost` and `http://service2.localhost` thanks to the integrated Traefik proxy.
- **Hooks in action:** Observe the execution of the `post-up` hook.

## 4. Diagnostics (`doctor`)
If you encounter environmental issues (e.g., SSH permission denied), use the built-in diagnostic tools.

**Run:**
```bash
<cli> doctor
# OR for advanced CLI features like interactive SSH setup
cd controller && go run main.go ssh
```

---

## Current Status
`workspace-controller` is now a professional orchestration tool:

- **Clean Root Structure:** The orchestration logic is tucked away in the `controller/` directory, keeping the project root focused on management via `Makefile` and `README.md`.
- **Simplified Entrypoint:** The `init` command combines planning and materialization into a single action.
- **Decoupled Architecture:** Infrastructure, platform, and tools are independent, version-controlled entities defined in `repos.yaml`.
- **Environment Variables:** Full support for dynamic paths via environment variables (e.g., `${WORKSPACE_ROOT}`) in the system definition.
- **Automated Lifecycle:** From initialization to health monitoring and shutdown.
