# LMBEK Hobby Workspace

This workspace combines several small repositories into one local development checkout and one minimal production platform. Production is a single Ubuntu 24.04 server running K3s and bundled Traefik; Terraform creates the Hetzner server and firewall, while GitHub Actions tests and publishes immutable GHCR images and deploys them through an outbound-only self-hosted runner.

## Architecture

```text
Developer -> GitHub -> hosted CI -> GHCR
                                 |
                                 v
                    outbound-only Actions runner
                                 |
                     localhost K3s API (RBAC)
                                 |
Internet -> 80/443 -> Traefik -> production namespace
```

There is no Argo CD, cert-manager, service mesh, dashboard, monitoring stack, public Kubernetes API, or CI SSH deployment. Traefik obtains and renews Let's Encrypt certificates itself. Only TCP 22 from configured trusted CIDRs and TCP 80/443 from the internet reach the server.

## Local development

```bash
make clone
make init-repo-envs
make up
```

Useful local routes are `http://web.localhost`, `http://placeholder1.localhost`, `http://placeholder2.localhost`, and `http://docs.localhost`. Stop the stack with `make down`; use `make hotreload` and `make hotreload-down` for live reload.

## Repository map

| Repository | Responsibility |
|---|---|
| `git-repositories/infrastructure/iac` | Hetzner server, firewall, secure cloud-init, K3s, Traefik ACME, deployment runner and RBAC bootstrap |
| `git-repositories/infrastructure/platform` | Production Kustomize manifests and scheduled CD workflow |
| `git-repositories/services/*` | Go service source, non-root images, tests, dependency scanning and GHCR publishing |
| `git-repositories/docs` | LikeC4 documentation and unprivileged documentation image |
| `local-orchestrator` | Docker Compose-based local development only |

## Production quick start

The authoritative end-to-end setup and GitHub configuration guide is in `git-repositories/infrastructure/iac/README.md`. In summary:

```bash
cd git-repositories/infrastructure/iac
cp terraform.tfvars.example terraform.tfvars
# Fill all placeholders, including a trusted SSH CIDR and one-time runner token.
terraform init
terraform fmt -check -recursive
terraform validate
terraform plan -out=tfplan
terraform apply tfplan
terraform output -raw ingress_target_ip
```

Point the configured domain's DNS `A` record to that output. Push each application repository's `main` branch once to publish its `production` tag; the platform workflow resolves those tags to immutable digests, records them in Git, rolls out the manifests, and verifies `https://<domain>/healthz`.

## Security and operations

- Runtime and CI secrets are never committed. `terraform.tfvars`, Terraform state, kubeconfigs, private keys, and `.env` files are ignored.
- The deployment runner uses a service-account token that can manage only selected resources in `production`; it cannot read Secrets or other namespaces.
- Pull requests never execute on the production runner. Hosted runners perform untrusted validation and image builds.
- Containers run as non-root, drop all capabilities, use seccomp, declare resource requests/limits, and have readiness/liveness probes.
- Internal services use `ClusterIP`; a default-deny NetworkPolicy permits only required namespace, DNS, ingress, and web egress traffic.
- Images are blocked on fixable `HIGH`/`CRITICAL` Trivy findings and dependency scanners. Exceptions should be reviewed and documented rather than disabling the gate.

Rollback the most recent application revision with:

```bash
kubectl -n production rollout undo deployment/<name>
kubectl -n production rollout status deployment/<name> --timeout=180s
```

Then revert the digest commit in the platform repository so Git again matches the cluster. Infrastructure and manifests are reproducible, but they are not backups: local K3s volumes and application databases require an independent off-server backup strategy.

Destroy the server and firewall from the IaC directory with `terraform destroy`. Destruction also destroys local K3s data.