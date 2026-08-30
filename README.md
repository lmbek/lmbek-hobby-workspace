# LMBEK Hobby Workspace

A multi-repository developer workspace with Docker Compose locally and a small GitOps
deployment on one Hetzner server. Terraform owns the server and firewall; Argo CD owns
the Kubernetes workloads.

---

## 🧭 3-Tier Lifecycle Overview

| Environment | Deploy / Up Command | Teardown / Down Command | Target URL & Ingress | Description |
|---|---|---|---|---|
| **1. Local** | `make up` or `make hotreload` | `make down` | `http://localhost` | Zero-cloud-cost local development with Traefik routing matching production |
| **2. Staging** | Record a staging digest in Git | Git revert | `https://staging.lmbek.dk` | Pre-production frontend with trusted automatic TLS |
| **3. Production** | Promote image digests in Git | Git revert | `https://lmbek.dk` | Immutable images with trusted automatic TLS |

---

## 🌐 Public URLs

Only the web frontend is publicly exposed in staging and production. Traffic passes
through Traefik with automatic SSL/TLS encryption:

| URL Endpoint | Service Routed | Environment | Ingress / Proxy Layer | TLS Security |
|---|---|---|---|---|
| `https://lmbek.dk` | Web Frontend Website | `production` | Traefik | Let's Encrypt SSL (Auto-renew) |
| `https://staging.lmbek.dk` | Web Frontend Website | `staging` | Traefik | Let's Encrypt SSL (Auto-renew) |

The placeholder services and documentation remain internal Kubernetes services and
have no public Ingress routes.

### 🔒 Security:
- **No public Kubernetes API**: The Hetzner firewall only exposes ports 22, 80, and 443.
- **Edge Security Headers**: Traefik enforces HTTP-to-HTTPS redirect, HSTS (HTTP Strict Transport Security), XSS filters, and rate limiting out-of-the-box.
- **Automated Certificate Lifecycle**: cert-manager continuously manages and renews Let's Encrypt SSL certificates via ACME HTTP-01 solvers.

---

## ⚡ Quick Start (Local Environment)

Get the entire platform running on your local machine in 4 simple commands:

```bash
# 1. Clone all repositories and create local workspaces
make clone

# 2. Automatically initialize all .env files from .env.example across all directories
make init-repo-envs

# 3. Start all services, websites, proxy, and docs locally
make up

# (Optional) Run with live hot-reloading and source mounting instead
make hotreload
```

Open in your browser:
- **Local Gateway / Proxy**: [http://localhost](http://localhost) (or [http://proxy.localhost](http://proxy.localhost))
- **Web Frontend**: [http://web.localhost](http://web.localhost)
- **Service 1 API**: [http://placeholder1.localhost](http://placeholder1.localhost)
- **Service 2 API**: [http://placeholder2.localhost](http://placeholder2.localhost)
- **Documentation**: [http://docs.localhost](http://docs.localhost)
- **Traefik Proxy / Infra**: [http://infra.localhost](http://infra.localhost)

To stop the local stack:
```bash
make down
```

---

## 🔑 Where to Put Your Hetzner API Key (One-Time Server Setup)

To provision the single K3s server on Hetzner Cloud:

### Step 1: Obtain Your Hetzner API Token
1. Go to [Hetzner Cloud Console](https://console.hetzner.cloud).
2. Select your Project &rarr; **Security** &rarr; **API Tokens** &rarr; click **Generate API Token** (Read & Write permissions).

### Step 2: Configure Your API Key
Navigate to `git-repositories/infrastructure/iac`:
```bash
cd git-repositories/infrastructure/iac
cp terraform.tfvars.example terraform.tfvars
```
Open `terraform.tfvars` and paste your token:
```hcl
hcloud_token   = "YOUR_HETZNER_API_TOKEN_HERE"
ssh_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5..." # Required emergency access key
server_type    = "cx23"                               # 2 vCPU, 4 GB RAM (~€3.79/mo)
server_location = "fsn1"                              # Falkenstein (or nbg1, hel1)
```
*(Alternatively, you can simply run `export HCLOUD_TOKEN="your-token"` in your terminal.)*

> **Security Guarantee:** `terraform.tfvars` is listed in `.gitignore` and will **never** be committed to Git.

### Step 3: Run Terraform Apply
*(Make sure Terraform is installed: `sudo snap install --classic terraform` on Ubuntu, or `brew install terraform` on macOS).*

> ⚠️ **Directory Requirement:** Run `terraform init` and `terraform apply` directly inside `git-repositories/infrastructure/iac`:

```bash
cd git-repositories/infrastructure/iac
terraform init
terraform apply
```
Terraform provisions one K3s server and one cloud firewall. Immutable cloud-init installs
K3s, cert-manager, and Argo CD; no server-side setup commands are required.

### Step 4: Updating Server Provisioning
Whenever you modify variables, firewall rules, domain settings, or cloud-init:
```bash
cd git-repositories/infrastructure/iac
terraform fmt -check
terraform validate
terraform plan
# Apply only after reviewing the plan.
terraform apply
```

Cloud-init changes can replace the server and change its public IP. Update the DNS
record after such a replacement. Do not SSH to production or run manual `kubectl`
commands; all operational changes belong in these repositories.

### Step 5: Shutting Down / Destroying Cloud Servers
When you want to stop billing or tear down the cloud servers completely:
```bash
cd git-repositories/infrastructure/iac
terraform destroy
```

---

## 🚀 How Releases Work (Pure GitOps — No SSH Needed)

After bootstrapping the server once with Terraform, **Git takes over completely**:

### 1. Staging Release Flow:
```
Developer pushes code to main
        │
        ▼
GitHub Actions builds container image & tags staging-latest / staging-<sha>
        │
        ▼
Developer records the immutable image digest in the staging overlay and merges it
```

### 2. Production Promotion Flow:
```
Developer verifies a staging image digest
        │
        ▼
Developer tests staging, then updates the digest in overlays/production/kustomization.yml
        │
        ▼
ArgoCD detects the Git change and syncs the Production namespace
```

### Day-to-Day Development Workflow:
1. **Local Development**: Edit website HTML/CSS/JS or backend services and test instantly with `make up` or `make hotreload`.
2. **Deploy to Staging**: Push service changes to the service repository `main` branch:
   ```bash
   git add .
   git commit -m "Update homepage UI and features"
   git push origin main
   ```
   Wait for GitHub Actions to pass and publish the GHCR image. Copy its `sha256` digest
   into `git-repositories/infrastructure/platform/overlays/staging/kustomization.yml`,
   then open and merge a platform pull request.
3. **Promote to Production**: Verify `https://staging.lmbek.dk`, copy the same tested
   digest into `overlays/production/kustomization.yml`, and merge that platform change.
   Argo CD deploys both merged changes automatically; no SSH or server-side command is needed.

---

## 🛠️ Available Workspace Commands

All commands can be executed directly from the workspace root:

| Command | Alias | Description |
|---|---|---|
| `make clone` | | Clones all repositories and initializes workspace directories |
| `make init-repo-envs` | `make envs` | Automatically creates `.env` from `.env.example` across all repos |
| `make up` | | Starts the entire local stack (web, services, docs, proxy) |
| `make hotreload` | | Starts local stack with live volume mounts for instant hot reloading |
| `make down` | | Stops all running local containers |
| `make down-v` | | Stops local containers and removes persistent volumes |
| `make status` | | Shows Git branch and status dashboard for all repositories |
| `make sync` | | Safely pulls upstream changes across all repositories (`--ff-only`) |
| `make fetch` | | Fetches all remotes across all repositories |
| `make checkout` | | Switches branches (`make checkout BRANCH=main` or `go run . checkout main`) |
| `make doctor` | | Verifies local development prerequisites (Git, Go, Docker, SSH) |
| `make ssh-helper` | `make ssh` | Interactive SSH key configuration helper |
| `make ps` | | Lists active container statuses and exposed ports |
| `make logs` | | Tails unified container logs |
| `make help` | | Shows available commands and quick help |

---

## 📁 Project Structure

```text
.
├── git-controller/              # Go CLI managing multi-repo workflows & operations
├── local-orchestrator/          # Unified Docker Compose & live hot-reload setup
│   ├── docker-compose.yml       # Root compose definition with modular includes
│   ├── compose/                 # Stack modules (websites, services, docs, proxy, observability)
│   ├── hotreload/               # Live source mounting overlays for instant UI/code feedback
│   └── proxy/                   # Traefik routing configuration mirroring production paths
├── tools/                       # Workspace development tools and helper scripts
├── git-repositories/            # Managed workspace repositories (gitignored)
│   ├── services/
│   │   ├── web-frontend/        # Responsive website frontend (Go + HTML5/CSS)
│   │   ├── placeholder1-service/# Microservice 1 (Port 8082)
│   │   └── placeholder2-service/# Microservice 2 (Port 8081)
│   ├── infrastructure/
│   │   ├── iac/                 # Terraform for one Hetzner server and firewall
│   │   ├── servers/             # Immutable K3s cloud-init
│   │   └── platform/            # Kubernetes base manifests, overlays, and ArgoCD CRDs
│   ├── deployment/              # Environment release configurations
│   ├── observability/           # Grafana dashboards & Prometheus alerts
│   ├── docs/                    # Architecture records (ADRs) & documentation portal
│   └── repo-definition.yaml     # Single source of truth for managed repositories
├── Makefile                     # Root developer interface
├── .gitignore                   # Workspace gitignore rules (safeguards secrets & keys)
└── README.md                    # This guide
```
