# Workspace Controller

The **Workspace Controller** is a declarative tool designed to manage complex, distributed local development environments. It allows you to define your entire system—services, versions, and infrastructure—in a single configuration file and orchestrate everything with simple commands.

## Vision
Move away from brittle bash scripts and manual setup. Use a single source of truth (`system-definition.yaml`) to materialize, validate, and run your local workspace.

## Key Features
- **Declarative Definition:** System state is described in YAML.
- **Automated Sync:** Clones and manages Git repositories automatically.
- **Infrastructure as Code:** Manages version-controlled infrastructure repositories.
- **Deployer Tool:** Initial support for Terraform-based deployments in `tools/deployer`.
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
Clone this repository and build the binary into the `bin` directory (or run directly with Go):
```bash
git clone <this-repo-url>
cd workspace-controller
go build -o bin/workspace-controller.exe main.go
```

---

### Using the Makefile (Recommended)
You can use the provided `Makefile` for a more convenient experience. It automatically handles the `WORKSPACE_ROOT` environment variable.
```bash
make help      # Show available commands
make init      # [1] Initialize workspace
make validate  # [2] Validate consistency
make up        # [3] Start the system
make down      # Stop the system
make doctor    # [4] Diagnose issues
```

### Manual Usage & Commands
If you prefer running commands manually, ensure you set the `WORKSPACE_ROOT` environment variable.

> **Note:** In the examples below, `<cli>` refers to the command you are using (e.g., `./bin/workspace-controller.exe`, `go run main.go`, or `make`).

### 1. `init`
Initializes the workspace by performing an execution plan (analysis) and then materializing the system (cloning/updating repositories).
```bash
$env:WORKSPACE_ROOT=".."
<cli> init
```

### 2. `validate`
Verifies that your local workspace is consistent. It checks for missing directories, uncommitted changes, and service health (HTTP/TCP).
```bash
$env:WORKSPACE_ROOT=".."
<cli> validate
```

### 3. `up`
Starts the entire environment (services and infrastructure) using Docker Compose.
```bash
$env:WORKSPACE_ROOT=".."
<cli> up
```

### 4. `down`
Stops and removes all containers associated with the workspace.
```bash
$env:WORKSPACE_ROOT=".."
<cli> down
```

### 5. `doctor`
Diagnoses environmental issues (Git, SSH, Docker).
```bash
<cli> doctor
```

### 6. `ssh-setup`
Interactive tool to manage SSH keys and configuration.
```bash
<cli> ssh-setup
```

### 7. `help`
Shows the available commands and usage information.
```bash
<cli> help
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
- **Infrastructure:** A single version-controlled infrastructure repository.
- **Tools:** A single version-controlled development tools repository.

### Environment Variables
The controller supports environment variables to override default paths:
- `SERVICES_DIR`: Directory for services (default: `../workspace/services`)
- `INFRA_DIR`: Directory for infrastructure (default: `../workspace/infrastructure`)
- `TOOLS_DIR`: Directory for tools (default: `../workspace/tools`)

---

## Documentation
For more detailed information, see:
- [.ai-agent/ARCHITECTURE.md](.ai-agent/ARCHITECTURE.md) - Internal structure and design principles.
- [.ai-agent/PROJECT.md](.ai-agent/PROJECT.md) - Project vision and core concepts.
- [.ai-agent/ROADMAP.md](.ai-agent/ROADMAP.md) - Current progress and future plans.
- [DEMO.md](DEMO.md) - Hands-on guide to the controller's features.
- [.ai-agent/HANDOVER.md](.ai-agent/HANDOVER.md) - Context for developers taking over the project.
- [.ai-agent/guidelines.md](.ai-agent/guidelines.md) - The guidelines we are using in this project for ai agent

*Note: Always ensure that this README is updated when major features or changes are implemented.*
