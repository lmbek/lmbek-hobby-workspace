# LMBEK Hobby Workspace

A simple, robust multi-repo cloud platform and developer workspace. Manage all microservices and websites locally with Docker Compose hot-reloading (`make up`), provision 2 cheap Hetzner Cloud servers with Terraform in minutes, and let ArgoCD GitOps automatically deploy every push to `main` to **Staging** and every GitHub Release to **Production**.

---

## 🧭 3-Tier Lifecycle Overview

| Environment | Deploy / Up Command | Teardown / Down Command | Target URL & Ingress | Description |
|---|---|---|---|---|
| **1. Local** | `make up` or `make hotreload` | `make down` | `http://localhost` | Zero-cloud-cost local development with Traefik routing matching production |
| **2. Staging** | `git push origin main` | `kubectl delete -k git-repositories/infrastructure/platform/overlays/staging` | `https://staging.<domain>` | Continuous pre-production environment auto-deployed on push to `main` with automated TLS |
| **3. Production** | Create GitHub Release (`v*.*.*`) | `cd git-repositories/infrastructure/iac && terraform destroy` | `https://<domain>` | Live stable production workloads with automated Let's Encrypt SSL and edge proxying |

---

## 🌐 Proxies, Load Balancers & Best Practice URLs

All web traffic passes through hardened reverse proxies with automatic SSL/TLS encryption:

| URL Endpoint | Service Routed | Environment | Ingress / Proxy Layer | TLS Security |
|---|---|---|---|---|
| `https://example.com` (or `https://web.example.com`) | Web Frontend Website | `production` | Traefik + Hetzner LB (optional) | Let's Encrypt Prod SSL (Auto-renew) |
| `https://placeholder1.example.com` | Microservice 1 API | `production` | Traefik Ingress | Let's Encrypt Prod SSL (Auto-renew) |
| `https://placeholder2.example.com` | Microservice 2 API | `production` | Traefik Ingress | Let's Encrypt Prod SSL (Auto-renew) |
| `https://docs.example.com` | Docs Portal | `production` | Traefik Ingress | Let's Encrypt Prod SSL (Auto-renew) |
| `https://staging.example.com` | Web Frontend Website | `staging` | Traefik Ingress | Automated Staging TLS Certificate |
| `https://placeholder1.staging.example.com` | Microservice 1 API | `staging` | Traefik Ingress | Automated Staging TLS Certificate |
| `https://placeholder2.staging.example.com` | Microservice 2 API | `staging` | Traefik Ingress | Automated Staging TLS Certificate |
| `https://docs.staging.example.com` | Docs Portal | `staging` | Traefik Ingress | Automated Staging TLS Certificate |

### 🔒 Security & Private Network Architecture:
- **Zero Kubernetes API Exposure**: Port 6443 (K3s API) is strictly bound to the private network (`10.0.1.0/24`) and blocked from the public internet by cloud firewall rules.
- **Edge Security Headers**: Traefik enforces HTTP-to-HTTPS redirect, HSTS (HTTP Strict Transport Security), XSS filters, and rate limiting out-of-the-box.
- **Automated Certificate Lifecycle**: cert-manager continuously manages and renews Let's Encrypt SSL certificates via ACME HTTP-01 solvers.
- **High Availability Edge Proxy**: Optional Hetzner Load Balancer (`lb11`) distributes incoming HTTP (80) and HTTPS (443) traffic across master and worker nodes.

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

To spin up your 2 cheap cloud servers (~€3.79/mo each) on Hetzner Cloud:

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
ssh_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5..." # (Optional) Your public SSH key
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
*That's it!* Terraform automatically provisions:
- **Server 1 (`k3s-master` / `10.0.1.10`)**: Installs K3s control plane, Traefik Ingress, and ArgoCD GitOps operator.
- **Server 2 (`k3s-worker` / `10.0.1.11`)**: Automatically connects to Server 1 over the private cloud network (`10.0.1.0/24`).
- **Cloud Firewall**: Pre-configured for SSH (22), HTTP (80), HTTPS (443), and intra-cluster communication.

### Step 4: Updating Server Provisioning
Whenever you modify variables, firewall rules, or domain settings:
```bash
cd git-repositories/infrastructure/iac
terraform plan
terraform apply
```

### Step 5: Shutting Down / Destroying Cloud Servers
When you want to stop billing or tear down the cloud servers completely:
```bash
cd git-repositories/infrastructure/iac
terraform destroy
```

---

## 🚀 How Releases Work (Pure GitOps — No SSH Needed)

After bootstrapping the servers once with Terraform, **code takes over completely**:

### 1. Staging Release Flow (Automated on push to `main`):
```
Developer pushes code to main
        │
        ▼
GitHub Actions builds container image & tags staging-latest / staging-<sha>
        │
        ▼
ArgoCD detects new image & auto-syncs to Staging namespace (https://staging.<your-ip>)
```

### 2. Production Release Flow (Automated on GitHub Release):
```
Developer creates a GitHub Release (e.g. v1.0.0)
        │
        ▼
GitHub Actions builds container image & tags latest, v1.0.0, <sha>
        │
        ▼
ArgoCD detects new release & auto-syncs to Production namespace (https://<your-ip>)
```

### Day-to-Day Development Workflow:
1. **Local Development**: Edit website HTML/CSS/JS or backend services and test instantly with `make up` or `make hotreload`.
2. **Deploy to Staging**: Push your changes to the `main` branch:
   ```bash
   git add .
   git commit -m "Update homepage UI and features"
   git push origin main
   ```
   *GitHub Actions automatically builds the `staging-latest` image and ArgoCD deploys it to the Staging namespace.*
3. **Promote to Production**: When staging validation passes, create a release on GitHub (e.g. `v1.0.0`):
   - GitHub Actions automatically builds the release container image (`latest` + `v1.0.0`).
   - ArgoCD auto-syncs the release into the live Production namespace with zero downtime!

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
│   │   ├── iac/                 # Terraform modules for Hetzner Cloud 2-node cluster
│   │   ├── servers/             # K3s master/worker cloud-init & Go platform CLI
│   │   └── platform/            # Kubernetes base manifests, overlays, and ArgoCD CRDs
│   ├── deployment/              # Environment release configurations
│   ├── observability/           # Grafana dashboards & Prometheus alerts
│   ├── docs/                    # Architecture records (ADRs) & documentation portal
│   └── repo-definition.yaml     # Single source of truth for managed repositories
├── Makefile                     # Root developer interface
├── .gitignore                   # Workspace gitignore rules (safeguards secrets & keys)
└── README.md                    # This guide
```
