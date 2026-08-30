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

# 3. Open terraform.tfvars and paste your Hetzner Cloud API token & SSH key:
# hcloud_token = "YOUR_HETZNER_API_TOKEN"
# ssh_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5..."
# (Optional) Restrict SSH to your IP for maximum security:
# allowed_ssh_ips = ["203.0.113.10/32"]

# 4. Initialize and apply (must be run from inside git-repositories/infrastructure/iac)
terraform init
terraform fmt -check
terraform validate
terraform plan
# Review the plan carefully, then apply only when ready.
terraform apply
```

---

### 6. Configure DNS & Verify from Your Local Environment

After `terraform apply`, get the public ingress address locally:

```bash
cd git-repositories/infrastructure/iac
terraform output -raw ingress_target_ip
```

Create this DNS `A` record at your DNS provider:

```text
@            A  <ingress_target_ip>
```

DNS providers usually use `@` as the record name; providers that require a fully
qualified name should use `lmbek.dk` instead. Do not add wildcard records: only
the web frontend is public.

Once DNS has propagated, verify the public routes from your local environment:

```bash
curl --fail --show-error --head https://lmbek.dk
curl --fail --show-error --head https://staging.lmbek.dk
```

The first certificates may take several minutes. Both public routes use trusted
Let's Encrypt certificates; local development remains HTTP-only. Do not test with
the server IP because HTTPS routing requires the hostname.

---

### 7. How to Update Infrastructure Provisioning

Whenever you edit `terraform.tfvars`, firewall rules, server locations, domain names,
or cloud-init:

```bash
cd git-repositories/infrastructure/iac

# Validate before creating a plan
terraform fmt -check
terraform validate

# Check what will change
terraform plan

# Apply only after reviewing the plan
terraform apply
```

Always review the plan before applying it. Changes to immutable server bootstrap
configuration, such as cloud-init, can require server replacement and new public IPs.
If the IP changes, update the single `@` DNS `A` record. Never SSH to production or
run manual server/Kubernetes commands; the repositories and Argo CD are the source
of truth.

---

### 8. Deploying an application change

Each service repository has GitHub Actions CI/CD. Use this workflow from your local
machine:

1. Run the service tests locally, then commit and push the service change to `main`.
2. Wait for GitHub Actions to test and publish the image to GHCR.
3. Read the published image's immutable `sha256` digest and update the matching image
   in `git-repositories/infrastructure/platform/overlays/staging/kustomization.yml`.
4. Open and merge the platform pull request. Argo CD deploys staging automatically.
5. Verify `https://staging.lmbek.dk`, then copy that exact digest to the production
   overlay, open and merge a second platform pull request.

Platform pull requests run Kustomize validation. Terraform pull requests run format
and configuration validation. Neither workflow applies Terraform or connects to
production.

---

### 9. How to Shut Down / Destroy Cloud Servers

When you want to stop billing or completely decommission the cloud servers:

```bash
cd git-repositories/infrastructure/iac

# Permanently destroy the server and firewall
terraform destroy
```
