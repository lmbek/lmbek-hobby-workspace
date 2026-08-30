# Quick Start Guide

## Local environment

```bash
make init-repo-envs
make up
```

Open `http://web.localhost`, `http://placeholder1.localhost`, `http://placeholder2.localhost`, or `http://docs.localhost`. Use `make down` to stop the stack.

## One-time GitHub preparation

1. In every service and docs repository, enable GitHub Actions and set **Settings -> Actions -> General -> Workflow permissions** to read repository contents. The workflow's explicit `packages: write` permission publishes to GHCR.
2. In the platform repository, permit the workflow's explicit `contents: write` and `packages: read` permissions. If `main` is protected, allow the GitHub Actions app to bypass only for the automated digest commit.
3. Create a `production` environment in the platform repository and restrict deployment branches to `main`. Required reviewers are optional; enabling them intentionally pauses every scheduled deployment for approval.
4. Make all four GHCR packages public, or grant the platform repository Actions read access under each package's **Manage Actions access** page.
5. Require each repository's CI workflow in the `main` branch ruleset. Never route pull-request jobs to the self-hosted runner.

Authenticate GitHub CLI and create the short-lived runner registration token immediately before provisioning:

```bash
gh auth login
gh api --method POST \
  repos/lmbek/lmbek-hobby-platform/actions/runners/registration-token \
  --jq .token
```

The token is single-use and expires after one hour. Put it only in your local `terraform.tfvars`; no Kubernetes credential, PAT, cloud token, or SSH key is stored in GitHub Actions.

## Provision production

```bash
cd git-repositories/infrastructure/iac
cp terraform.tfvars.example terraform.tfvars
# Set hcloud_token, ssh_public_key, allowed_ssh_ips, domain,
# letsencrypt_email, and github_runner_token.
chmod 600 terraform.tfvars
terraform init
terraform fmt -check -recursive
terraform validate
terraform plan -out=tfplan
terraform apply tfplan
terraform output -raw ingress_target_ip
```

Create one DNS `A` record for the configured domain pointing to `ingress_target_ip`. Do not create a record for the Kubernetes API and do not open TCP 6443.

## First and subsequent deployments

Push `main` in the three service repositories and docs repository. Each hosted CI workflow runs tests/build checks, dependency scanning, a container build, Trivy, and then publishes `sha-<full-git-sha>` plus the moving `production` discovery tag.

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