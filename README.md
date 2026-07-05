# Workspace Controller

The **Workspace Controller** is a declarative tool designed to manage complex, distributed local development environments. It allows you to define your entire system—services, versions, and infrastructure—in a single configuration file and orchestrate everything with simple commands.

## Vision
Move away from brittle bash scripts and manual setup. Use a single source of truth (`repos.yaml`) to materialize, validate, and run your local workspace.

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
- [Go](https://golang.org/doc/install) (1.26 or later)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

### Installation
**Note:** This project uses SSH for all Git operations. Ensure your SSH keys are configured and added to your GitHub account before proceeding.

Clone this repository and build the binary into the `bin` directory (or run directly with Go):
```bash
git clone git@github.com:Lars-Bek/LMBEK-HOBBY-WORKSPACE.git
cd LMBEK-HOBBY-WORKSPACE
go build -o bin/workspace-controller.exe ./controller/main.go
```

---

### Using the Makefile (Recommended)
You can use the provided `Makefile` for a more convenient experience. It automatically handles the `WORKSPACE_ROOT` environment variable.

#### SSH Setup (First Time Only)
If you haven't configured SSH for GitHub on this machine:
```bash
make ssh-setup   # Follow the interactive prompts to generate/add keys
```

#### Core Workflow
```bash
make help      # Show available commands
make init      # [1] Bootstrap workspace (Pre-flight + Planning)
make sync      # [2] Synchronize all repositories (Materialization)
make validate  # [3] Validate consistency and health
make up        # [4] Start the system
make down      # Stop the system
```

#### Tooling & Diagnostics
```bash
make doctor    # [D] Diagnose environment issues
```

#### Development & Testing
Commands for development, testing, and service management are located within their respective directories:
- **Controller Logic:** `cd controller && make help`
- **Services:** Each service in `applications/` has its own `Makefile` (e.g., `make test`, `make build`).

### Manual Usage & Commands
If you prefer running commands manually, ensure you set the `WORKSPACE_ROOT` environment variable.

> **Note:** In the examples below, `<cli>` refers to the command you are using (e.g., `./bin/workspace-controller.exe`, `go run main.go`, or `make`).

### 1. `init`
Bootstraps the workspace by performing pre-flight checks (SSH connectivity, agent status) and creating an execution plan.
```bash
$env:WORKSPACE_ROOT="."
<cli> init
```

### 2. `sync`
Materializes the system by cloning or updating all repositories defined in the system definition. It also runs post-sync hooks.
```bash
$env:WORKSPACE_ROOT="."
<cli> sync
```

### 3. `validate`
Verifies that your local workspace is consistent. It checks for missing directories, branch matches (defaulting to `main`), and service health (HTTP/TCP).
```bash
$env:WORKSPACE_ROOT="."
<cli> validate
```

### 4. `up`
Starts the entire environment (services and infrastructure) using Docker Compose.
```bash
$env:WORKSPACE_ROOT="."
<cli> up
```

### 5. `down`
Stops and removes all containers associated with the workspace.
```bash
$env:WORKSPACE_ROOT="."
<cli> down
```

### 6. `doctor` [D]
Diagnoses environmental issues (Git, SSH, Docker) and provides automated fixes for common SSH agent and key problems. This is an integrated diagnostic tool in the workspace-controller CLI.
```bash
<cli> doctor
```

### 7. `ssh-setup` / `ssh` [S]
Interactive tool to manage SSH keys, configure `~/.ssh/config`, and ensure Git is using the correct SSH client (especially important on Windows). Both `ssh-setup` and the `ssh` alias are integrated commands of the workspace-controller CLI.
```bash
<cli> ssh-setup
# or
<cli> ssh
```

### 8. `help`
Shows the available commands and usage information.
```bash
<cli> help
```

---

## Structure
The project is organized to keep the root clean while having all managed components readily available:
- `applications/`: Customer/business services (e.g., Auth, User).
- `proxy/`: Central workspace dashboard and reverse proxy.
- `orchestrator/`: Dedicated orchestration repositories (e.g., Docker Compose, Kubernetes).
- `infrastructure/`: Cloud resources, servers, and networking.
- `platform/`: Developer tooling and observability.
- `tools/`: Local helper scripts and CLIs.
- `docs/`: Architecture and onboarding documentation.
- `repos.yaml`: The manifest defining all managed repositories.
- `controller/`: The Go source code for the **Workspace Controller** orchestrator.
- `Makefile`: Root-level entry point for managing the entire workspace.
- `README.md`: Project documentation.

The system is defined in `repos.yaml`. You can specify:
- **Applications:** Git repository URLs, versions (tags/branches), and environment variables.
- **Infrastructure:** Multiple version-controlled infrastructure repositories.
- **Platform:** Developer tooling, deployment systems, and observability repositories.
- **Tools:** Local helper scripts and CLI repositories.
- **Docs:** Architecture and onboarding documentation repositories.

### Environment Variables
The controller supports environment variables to override default paths and control behavior:
- `DEBUG`: Set to `true` to enable verbose debug logging (default: `false`).
- `WORKSPACE_ROOT`: The absolute path to the workspace root.
- `SERVICES_DIR`: Directory for applications (default: `./applications`).
- `INFRA_DIR`: Directory for infrastructure (default: `./infrastructure`).
- `PLATFORM_DIR`: Directory for platform (default: `./platform`).
- `TOOLS_DIR`: Directory for tools (default: `./tools`).
- `DOCS_DIR`: Directory for documentation (default: `./docs`).

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
