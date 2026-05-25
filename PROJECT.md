# Project: Workspace Controller (Local Environment Orchestrator)

## Vision
The purpose of this project is to create a tool that allows defining and starting a complete distributed system locally with a single command. We want to move away from complex bash scripts towards a declarative approach, where the system state is described in a central configuration file.

## Core Concepts
- **Declarative Definition:** The entire system (services, versions, infrastructure) is described in `system-definition.yaml`.
- **No Scripting:** Logic resides in the controller, not in scattered scripts. This ensures consistency across machines.
- **Isolated Workspace:** All services are materialized in a dedicated `workspace/` directory to avoid cluttering the rest of the system.
- **Infrastructure as Code:** Infrastructure (databases, queues, etc.) is managed centrally and automatically.

## Structure
- `workspace/`: (Project root) The main container for your local environment.
  - `services/`: Contains source code for services (e.g., `authentication-service/`, `user-service/`).
  - `infrastructure/`: Contains configuration for running shared services (Docker Compose, environment variables).
  - `tools/`: Additional tools and utilities.
- `workspace-controller/`: (Project root) The tool itself that orchestrates everything.

---
*Refer to [README.md](README.md) for usage instructions and quick start. Update the README when project goals or structure change.*
