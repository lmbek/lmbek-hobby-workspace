# Local Orchestrator

Docker Compose orchestrator wiring together all workspace services for local development and live reload.

## Local Subdomains & Endpoints

When running `make up` or `make hotreload`, Traefik automatically routes requests across all services:

| Subdomain / URL | Service | Description |
|---|---|---|
| `http://localhost` / `http://gateway.localhost` | Local Gateway | Proxy homepage dashboard & service catalog |
| `http://web.localhost` | Web Frontend | Responsive HTML5/Go landing page |
| `http://placeholder1.localhost` | Placeholder 1 Service | Microservice 1 API (Port 8082) |
| `http://placeholder2.localhost` | Placeholder 2 Service | Microservice 2 API (Port 8081) |
| `http://docs.localhost` | Docs Portal | LikeC4 Architecture Documentation |
| `http://infra.localhost` | Traefik Proxy | Traefik reverse proxy dashboard & API |
| `http://grafana.localhost` | Grafana | Observability metrics & log dashboards |
| `http://prometheus.localhost` | Prometheus | Metrics scraper & TSDB |
| `http://alloy.localhost` | Grafana Alloy | OpenTelemetry / telemetry collector |
| `http://loki.localhost` | Grafana Loki | Centralized log aggregation |
| `http://cadvisor.localhost` | cAdvisor | Container resource metrics |

## Quick Start

```bash
make up         # Start all services in the background (build & run)
make down       # Stop all services
make hotreload  # Start all services with live source mounting & hot reload
make ps         # View status of running containers
make logs       # Tail logs across all services
```

Start/stop a single service:

```bash
make service-up   S=placeholder1-service
make service-down S=placeholder1-service
make service-logs S=placeholder1-service
```

## Structure

```
local-orchestrator/
├── docker-compose.yml          # Unified entry compose file with includes & networks
├── compose/                    # Modular compose definitions
│   ├── proxy.yml               # Traefik reverse proxy & homepage
│   ├── services.yml            # Application microservices (build definitions)
│   ├── docs.yml                # LikeC4 architecture docs service
│   └── observability.yml       # Observability stack include (Prometheus, Grafana, Alloy, Loki)
├── hotreload/                  # Live reload orchestrator
│   ├── docker-compose.yml      # Includes base docker-compose.yml and mount overlays
│   └── mounts/                 # Volume mounting definitions for hot reload
│       ├── services.yml        # Source mounts for Go microservices
│       └── docs.yml            # Source mounts for Documentation
├── proxy/                      # Traefik configuration & error pages
│   ├── traefik.yml             # Traefik static configuration
│   ├── dynamic.yml             # Traefik routing rules & services
│   └── homepage/               # Proxy landing and error pages (502, 503, 504)
├── Makefile                    # Orchestration commands
└── README.md
```

## Available Makefile Targets

| Target | Description |
|--------|-------------|
| `make up` | Start full local stack in background |
| `make down` | Stop full local stack |
| `make restart` | Restart all services |
| `make hotreload` | Start full stack with live volume mounts & hot reload (`make dev`) |
| `make down-hotreload` | Stop hotreload stack |
| `make ps` | Show container statuses |
| `make logs` | Tail logs for all services |
| `make service-up S=<name>` | Start a specific service |
| `make service-down S=<name>` | Stop a specific service |
| `make service-logs S=<name>` | Tail logs for a specific service |
| `make validate` | Validate standard and hotreload compose configurations |
| `make clean` | Stop containers and remove volumes/orphans |
| `make help` | Show target list and usage instructions |
