# Workspace Controller

The **Workspace Controller** is a declarative tool designed to manage complex, distributed local development environments. It allows you to define your entire system—services, versions, and infrastructure—in a single configuration file and orchestrate everything with simple commands.

## Vision
Move away from brittle bash scripts and manual setup. Use a single source of truth (`system-definition.yaml`) to materialize, validate, and run your local workspace.

## Key Features
- **Declarative Definition:** System state is described in YAML.
- **Automated Sync:** Clones and manages Git repositories automatically.
- **Infrastructure as Code:** Automatically generates Docker Compose configurations.
- **Dependency Management:** Define service dependencies (e.g., waiting for database).
- **Hook System:** Run custom commands after sync or startup.
- **Consistency Checks:** Validates that your local environment matches the definition.
- **Health Monitoring:** Built-in support for HTTP and TCP health checks during validation.

---

## Getting Started

### 1. Try the Demo
For a hands-on experience of all features, follow the [DEMO.md](DEMO.md) guide.

### 2. Prerequisites
- [Go](https://golang.org/doc/install) (1.20 or later)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

### Installation
Clone this repository and build the binary (or run directly with Go):
```bash
git clone <this-repo-url>
cd workspace-controller
```

---

## Usage & Commands

The controller is operated through a set of explicit commands:

### 1. `start`
Generates an execution plan based on `system/system-definition.yaml`. It also creates the `infrastructure/docker-compose.yaml` file.
```bash
go run main.go start
```

### 2. `sync`
Materializes the system. It creates the `workspace/` directory and clones all services defined in the configuration.
```bash
go run main.go sync
```

### 3. `validate`
Verifies that your local workspace is consistent with the definition. It checks for missing directories, version mismatches, and uncommitted changes.
```bash
go run main.go validate
```

### 4. `up`
Starts the entire environment (services and infrastructure) using Docker Compose.
```bash
go run main.go up
```

### 5. `down`
Stops and removes all containers associated with the workspace.
```bash
go run main.go down
```

### 6. `help`
Shows the available commands and usage information.
```bash
go run main.go help
```

---

## Structure
The project is organized to keep everything related to the workspace together:
- `workspace/`: (Root) The main container for your local environment.
  - `services/`: Contains source code for services (e.g., `authentication-service/`, `user-service/`).
  - `infrastructure/`: Configuration for shared services (Docker Compose).
  - `tools/`: Additional development tools and utilities.
- `workspace-controller/`: (Root) The tool itself that orchestrates everything.

The system is defined in `system/system-definition.yaml`. You can specify:
- **Services:** Git repository URLs, versions (tags/branches), and environment variables.
- **Infrastructure:** Docker images, versions, and now **version-controlled configuration** (Git).
- **Tools:** Development utilities managed via Git.

---

## Documentation
For more detailed information, see:
- [ARCHITECTURE.md](ARCHITECTURE.md) - Internal structure and design principles.
- [PROJECT.md](PROJECT.md) - Project vision and core concepts.
- [ROADMAP.md](ROADMAP.md) - Current progress and future plans.
- [HANDOVER.md](HANDOVER.md) - Context for developers taking over the project.

*Note: Always ensure that this README is updated when major features or changes are implemented.*
