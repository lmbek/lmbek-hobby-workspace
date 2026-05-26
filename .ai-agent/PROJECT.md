# Project: Workspace Controller (Local Environment Orchestrator)

## Vision
The purpose of this project is to create a tool that allows defining and starting a complete distributed system locally with a single command. We want to move away from complex bash scripts towards a declarative approach, where the system state is described in a central configuration file.

A core focus is on **Cross-Platform Reliability**, ensuring that authentication (SSH) and repository management (Git) work seamlessly on Windows, Linux, and WSL without manual intervention.

## Core Concepts
- **Declarative Definition:** The entire system (services, versions, infrastructure) is described in `system-definition.yaml`.
- **Command-Based Lifecycle:** Clear separation between bootstrapping (`init`), materializing (`sync`), validating (`validate`), and running (`up`).
- **SSH & Git Synchronization:** The controller manages SSH agents and forces Git to use verified SSH clients to prevent authentication mismatches.
- **No Scripting:** Logic resides in the controller, not in scattered scripts. This ensures consistency across machines.
- **Isolated Workspace:** All services are materialized in a dedicated `workspace/` directory to avoid cluttering the rest of the system.
- **Infrastructure as Code:** Infrastructure (databases, queues, etc.) configuration is managed in a dedicated, version-controlled repository.
- **Tools & Deployers:** Specialized tools (like Terraform-based deployers) handle the actual deployment of infrastructure.

## Structure
- `workspace/`: (Project root) The main container for your local environment.
  - `services/`: Contains source code for services (e.g., `authentication-service/`, `user-service/`).
  - `infrastructure/`: Contains version-controlled configuration for shared services (Docker Compose, environment variables).
  - `tools/`: Additional version-controlled development tools and utilities.
- `workspace-controller/`: (Project root) The tool itself that orchestrates everything.

---
*Refer to [../README.md](../README.md) for usage instructions and quick start. Update the README when project goals or structure change.*
