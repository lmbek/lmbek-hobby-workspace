# Workspace Controller - Demo Guide

This guide will walk you through the core features of the **Workspace Controller** as of **Phase 5: Decoupling & Specialization**. We have moved from a monolithic setup to a flexible, decoupled architecture where services, infrastructure, and tools are independent, version-controlled entities.

## 1. Planning (`start`)
The `start` command focuses on orchestration logic. It reads your `system-definition.yaml` and acknowledges infrastructure as a decoupled, version-controlled repository.

**Run:**
```bash
go run main.go start
```

**What to look for:**
- A "SYSTEM START PLAN" listing your services and infrastructure.
- Confirmation that infrastructure is managed within its own repository.

## 2. Materializing the Workspace (`sync`)
The `sync` command handles the materialization of services, infrastructure, and tools.

**Run:**
```bash
go run main.go sync
```

**What to look for:**
- Creation of the `workspace/` structure with `services/`, `infrastructure/`, and `tools/`.
- Each component being updated via `git fetch` and `git pull` while staying on the `main` branch.
- **Hooks in action:** Observe the `post-sync` hook execution (e.g., success message).

## 3. Specialized Tools: The Deployer
We've introduced specialized tools like the Terraform-based `deployer` in `workspace/tools/deployer`.

**Explore:**
Check out the files in `../workspace/tools/deployer/`. This shows how the system can scale to use industry-standard tools for infrastructure.

## 4. Validating Consistency (`validate`)
The `validate` command performs consistency checks across all components.

**Run:**
```bash
go run main.go validate
```

**What to look for:**
- Verification that all Git components match the defined versions.
- Local "dirty" state detection.
- **Health Checks:** Real-time HTTP/TCP monitoring of services.

## 5. Lifecycle Management (`up` / `down`)
Control your local environment via Docker Compose, targeting the decoupled infrastructure.

**Try starting:**
```bash
go run main.go up
```

**What to look for:**
- The controller executing `docker-compose up -d` within the infrastructure directory.
- **Hooks in action:** Observe the `post-up` hook execution.

---

## Current Status (Phase 5)
The `workspace-controller` is now a mature orchestration tool:

- **Decoupled Architecture:** Infrastructure and tools are first-class versioned citizens.
- **Specialized Tooling:** Integration with tools like Terraform.
- **Service-Centric:** All services are neatly organized in `workspace/services/`.
- **Automated Lifecycle:** From materialization to health monitoring and shutdown.

## Use Cases
1. **Deterministic Environments:** Every developer runs the exact same versions of code and infra.
2. **Simplified Onboarding:** Go from zero to a running system with `sync` and `up`.
3. **Infrastructure as Code:** Manage local infra with the same rigor as production.
