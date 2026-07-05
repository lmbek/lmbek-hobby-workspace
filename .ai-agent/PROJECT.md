# Project: Workspace Controller (Local Environment Orchestrator)

## Vision
The purpose of this project is to create a tool that allows defining and starting a complete distributed system locally with a single command. We want to move away from complex bash scripts towards a declarative approach, where the system state is described in a central configuration file.

A core focus is on **Cross-Platform Reliability**, ensuring that authentication (SSH) and repository management (Git) work seamlessly on Windows, Linux, and WSL without manual intervention.

## Core Concepts
- **Declarative Definition:** The entire system (services, versions, infrastructure) is described in `repos.yaml`.
- **Command-Based Lifecycle:** Clear separation between bootstrapping (`init`), materializing (`sync`), validating (`validate`), and running (`up`).
- **SSH & Git Synchronization:** The controller manages SSH agents and forces Git to use verified SSH clients to prevent authentication mismatches.
- **No Scripting:** Logic resides in the controller, not in scattered scripts. This ensures consistency across machines.
- **Root-Level Workspace:** All services are materialized directly in the project root to keep everything accessible while maintaining a clean organization.
- **Infrastructure as Code:** Infrastructure (databases, queues, etc.) configuration is managed in a dedicated, version-controlled repository.
- **Tools & Deployers:** Specialized tools (like Terraform-based deployers) handle the actual deployment of infrastructure.

## Structure
- `applications/`: Customer/business services.
- `proxy/`: Central workspace dashboard and reverse proxy.
- `orchestrator/`: Dedicated orchestration repositories.
- `infrastructure/`: Servers, networking, cloud, and Docker configurations.
- `platform/`: Developer tooling and observability.
- `tools/`: Local helper scripts and CLIs.
- `docs/`: Architecture and onboarding documentation.
- `repos.yaml`: The manifest defining all managed repositories.
- `controller/`: The Go source code for the **Workspace Controller** orchestrator.

---
*Refer to [../README.md](../README.md) for usage instructions and quick start. Update the README when project goals or structure change.*
