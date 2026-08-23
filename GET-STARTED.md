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

- **Website Frontend**: [http://web.localhost](http://web.localhost) (or [http://localhost](http://localhost))
- **Microservice 1 API**: [http://placeholder1.localhost](http://placeholder1.localhost) (or [http://localhost/service1](http://localhost/service1))
- **Microservice 2 API**: [http://placeholder2.localhost](http://placeholder2.localhost) (or [http://localhost/service2](http://localhost/service2))
- **Documentation Portal**: [http://docs.localhost](http://docs.localhost) (or [http://localhost/docs](http://localhost/docs))
- **Traefik Proxy / Infra**: [http://infra.localhost](http://infra.localhost) (or [http://proxy.localhost](http://proxy.localhost))

---

### 3. Check Workspace & Git Status

```bash
# Check status across all managed repositories
make status
```

---

### 4. Stop Local Environment

```bash
make down
```

---

### 5. Deploy to Hetzner Cloud (When Ready)

```bash
cd git-repositories/infrastructure/iac

# Copy configuration template
cp terraform.tfvars.example terraform.tfvars

# Open terraform.tfvars and paste your Hetzner Cloud API token:
# hcloud_token = "YOUR_HETZNER_API_TOKEN"

# Provision the 2-node K3s Kubernetes cluster
terraform init
terraform apply
```
