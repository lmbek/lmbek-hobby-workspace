# git-controller-ansible

Ansible-based port of the Go `git-controller`. Manages all workspace repositories
defined in `system-definition.yaml` using Ansible playbooks.

## Requirements

- **Ansible 2.14+** (with `ansible-playbook` on PATH)
- **Git** (on PATH)
- **SSH** configured for GitHub access

## Quick Start

```bash
# From the workspace root:
make -f Makefile-ansible status

# Or run playbooks directly:
cd git-controller-ansible
ansible-playbook playbooks/status.yml
```

## Available Commands

All commands mirror the Go and PHP versions:

| Command | Description |
|---------|-------------|
| `clone` | Clone all repositories (skip already cloned) |
| `fetch` | Fetch all remotes |
| `pull` | Pull updates (clone if missing) |
| `push` | Push local commits |
| `checkout` | Switch repos to their defined branch |
| `scaffold` | Init .git and set remote origin (no clone needed) |
| `status` | Show dashboard overview of all repository states |
| `update` | Fetch + pull + status in one go |
| `validate` | Verify repos match the system definition |
| `doctor` | Diagnose environment (Git, Ansible, SSH, Docker) |

## Project Structure

```
git-controller-ansible/
├── ansible.cfg              Ansible configuration
├── inventory.yml            Local inventory (localhost)
├── filter_plugins/
│   └── repo_filters.py     Custom filter to parse system-definition.yaml
├── playbooks/
│   ├── clone.yml
│   ├── fetch.yml
│   ├── pull.yml
│   ├── push.yml
│   ├── checkout.yml
│   ├── scaffold.yml
│   ├── status.yml
│   ├── update.yml
│   ├── validate.yml
│   └── doctor.yml
└── README.md
```
