# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 4.x     | ✅ Active support  |
| 3.x     | ⚠️ Security fixes only |
| < 3.0   | ❌ End of life     |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly:

1. **Do NOT open a public GitHub issue.**
2. Email the maintainers directly at (TODO: to be added later) with:
   - A description of the vulnerability.
   - Steps to reproduce.
   - Potential impact assessment.
   - Any suggested fix (optional).
3. You will receive an acknowledgement within **2 weeks**.
4. We aim to release a fix within **1 month** for critical issues.

## Security Best Practices

This workspace enforces the following security measures:

- **SSH-only Git access** — HTTP/HTTPS clone URLs are rejected by the git-controller.
- **No secrets in code** — Use environment variables or the deployment repo's secrets templates.
- **Dependency scanning** — Keep Go modules up to date and audit with `go mod verify`.
- **Signed commits** — Encouraged for all contributors (not yet enforced).

## Disclosure Policy

- We follow [coordinated vulnerability disclosure](https://en.wikipedia.org/wiki/Coordinated_vulnerability_disclosure).
- Credit will be given to reporters unless they prefer to remain anonymous.
