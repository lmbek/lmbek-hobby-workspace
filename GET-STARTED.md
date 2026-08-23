# 🚀 Quick Start Guide

Run these commands in your terminal to get started immediately.

---

### 1. Initialize & Start Local Environment

```bash
# Initialize all .env files across repositories
make init-repo-envs

# Start all local containers (web frontend, microservices, docs, proxy)
make up
```

*(Alternatively, run `make hotreload` for instant live-reloading when editing code or web pages).*

---

### 2. Test in Your Browser

- **Local Gateway / Proxy**: [http://localhost](http://localhost) (or [http://proxy.localhost](http://proxy.localhost))
- **Website Frontend**: [http://web.localhost](http://web.localhost)
- **Microservice 1 API**: [http://placeholder1.localhost](http://placeholder1.localhost)
- **Microservice 2 API**: [http://placeholder2.localhost](http://placeholder2.localhost)
- **Documentation Portal**: [http://docs.localhost](http://docs.localhost)
- **Traefik Proxy / Infra**: [http://infra.localhost](http://infra.localhost)

---

### 3. Check Workspace & Git Status

```bash
# Check status across all managed repositories
make status
```

---

### 4. Stop Local Environment

```bash
# Stop standard environment
make down

# (Or stop hotreload environment)
make hotreload-down
```

---

### 5. Deploy to Hetzner Cloud (When Ready)

> ⚠️ **Important Directory Requirement:** You **must** navigate to the `git-repositories/infrastructure/iac` directory before running any Terraform commands.

```bash
# 0. Install Terraform (if not already installed)
# Ubuntu / Linux (Snap): sudo snap install --classic terraform
# macOS (Homebrew):       brew install terraform
# Windows (Chocolatey):   choco install terraform

# 1. Navigate to the IaC directory (REQUIRED)
cd git-repositories/infrastructure/iac

# 2. Copy configuration template
cp terraform.tfvars.example terraform.tfvars

# 3. Open terraform.tfvars and paste your Hetzner Cloud API token:
# hcloud_token = "YOUR_HETZNER_API_TOKEN"

# 4. Initialize and apply (must be run from inside git-repositories/infrastructure/iac)
terraform init
terraform apply
```
