# Workspace Controller - Demo Guide

This guide will walk you through the key features of the **Workspace Controller**. Since the default configuration uses placeholder Git URLs, this demo focuses on the orchestration and validation logic that you can experience immediately.

## 1. Planning and Infrastructure Generation
The `start` command reads your `system-definition.yaml` and prepares the environment.

**Run:**
```bash
go run main.go start
```

**What to look for:**
- A detailed "SYSTEM START PLAN" printed to your console.
- A new file created at `infrastructure/docker-compose.yaml`.
- Notice how `authentication-service` has `depends_on: [postgres]` in the generated YAML.

## 2. Materializing the Workspace
The `sync` command sets up your local folders and attempts to fetch source code.

**Run:**
```bash
go run main.go sync
```

**What to look for:**
- Creation of the `workspace/` and `infrastructure/` directories.
- The controller safely skips the placeholder URLs.
- **Hook in action:** You will see `[HOOK] echo "Sync complete! ..."` executed at the end.

## 3. Validating Consistency and Health
The `validate` command checks if your local state matches the definition and performs real-time health checks.

**Run:**
```bash
go run main.go validate
```

**What to look for:**
- `[MISSING]` status for services (since we haven't cloned real repos yet).
- `[UNHEALTHY]` status for Postgres (unless you have it running on port 5432).
- Helpful hints on how to fix the issues.

## 4. Lifecycle Management
You can control the entire stack with `up` and `down`.

**Try starting (requires Docker):**
```bash
go run main.go up
```
*Note: This will attempt to build services in the `workspace/` folder. Since they are currently empty/missing, this may fail until you add real services.*

**Cleaning up:**
```bash
go run main.go down
```

---

## Where are we so far?
We have built a **fully functional core** for a Local Environment Orchestrator. 

- **Phase 0-3 (Foundations & Sync):** We can parse complex YAML definitions and materialize a workspace from Git.
- **Phase 4 (Orchestration):** We automatically generate Docker Compose files, handling networking and environment variables.
- **Phase 5 (Advanced):** We have implemented:
    - **Health Checks:** Real-time TCP/HTTP monitoring.
    - **Dependency Management:** Order-aware container startup.
    - **Hooks:** Automation of post-setup tasks.

## What can we use it for?
1. **Onboarding:** A new developer clones the `workspace-controller` and runs `sync` + `up` to get everything running in minutes.
2. **Consistency:** Ensure every team member is running the exact same versions of 10+ microservices and infrastructure.
3. **Automation:** Replace 500-line bash scripts with a clean 30-line YAML file.
4. **Validation:** Quickly check if your local Git branches are out of sync with what the system expects.
