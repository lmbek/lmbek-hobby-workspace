# Quick Start Guide

## Local environment

```bash
make init-repo-envs
make up
```

Open `http://web.localhost`, `http://placeholder1.localhost`, `http://placeholder2.localhost`, or `http://docs.localhost`. Use `make down` to stop the stack.

## One-time GitHub preparation

```bash
make github-setup
```

This idempotent command installs GitHub CLI from its official Debian/Ubuntu repository when needed. Unless `GH_TOKEN` is already set, it opens GitHub's browser authorization once; authentication cannot be delegated safely.

The command then enables Actions with read-only defaults in every repository, creates the platform's `production` environment restricted to `main`, runs and waits for any missing initial image builds, makes all four GHCR packages public, and installs required-CI `main` rulesets. The platform ruleset allows only the GitHub Actions integration to bypass it for the automated digest commit. Run it as a GitHub account with administration access to all repositories; the account must have a plan that supports repository rulesets.

The defaults target the `lmbek` account. Set `GITHUB_OWNER=another-owner` on both GitHub commands if the repositories were forked or transferred. Pull requests always use GitHub-hosted runners; only the trusted platform deployment job uses the production self-hosted runner.

## Provision production

```bash
cp git-repositories/infrastructure/iac/terraform.tfvars.example \
  git-repositories/infrastructure/iac/terraform.tfvars
# Set hcloud_token, ssh_public_key, allowed_ssh_ips, domain, and letsencrypt_email.
# Leave github_runner_token alone; the next command writes it without printing it.
chmod 600 git-repositories/infrastructure/iac/terraform.tfvars
make github-runner-token
cd git-repositories/infrastructure/iac
terraform init
terraform fmt -check -recursive
terraform validate
terraform plan -out=tfplan
terraform apply tfplan
terraform output -raw ingress_target_ip
```

`make github-runner-token` installs/authenticates GitHub CLI if necessary, requests the short-lived registration token, and updates only the ignored local `terraform.tfvars`. The token is single-use, expires after one hour, and must be regenerated immediately before any Terraform operation that replaces the server. No Kubernetes credential, PAT, cloud token, or SSH key is stored in GitHub Actions.

Create one DNS `A` record for the configured domain pointing to `ingress_target_ip`. Do not create a record for the Kubernetes API and do not open TCP 6443.

## First and subsequent deployments

Merge to `main` in the three service repositories and docs repository. Each hosted CI workflow runs tests/build checks, dependency scanning, a container build, Trivy, and then publishes `sha-<full-git-sha>` plus the moving `production` discovery tag.

The platform `Deploy production` workflow runs every five minutes or from **Actions -> Deploy production -> Run workflow**. It resolves each discovery tag to a digest, updates Git, applies the production Kustomization using namespace-only RBAC, waits for every rollout, and verifies HTTPS.

```bash
curl --fail --show-error --head https://lmbek.dk
```

Replace `lmbek.dk` if `domain` was changed in Terraform and `overlays/production/kustomization.yml`.

## Emergency operations

SSH is key-only, non-root, restricted by `allowed_ssh_ips`, and is never used by CI:

```bash
terraform output -raw admin_ssh_command
```

Rollback from an emergency administration session:

```bash
kubectl -n production rollout history deployment/web-frontend
kubectl -n production rollout undo deployment/web-frontend
kubectl -n production rollout status deployment/web-frontend --timeout=180s
```

Revert the corresponding platform digest commit afterward. To remove all cloud infrastructure, run `terraform destroy`; local persistent data is not recoverable unless separately backed up.