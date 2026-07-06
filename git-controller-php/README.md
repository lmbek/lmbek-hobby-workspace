# git-controller-php

PHP version of the workspace git-controller — a CLI tool that manages multi-repo workspaces using `system-definition.yaml`.

This is a direct port of the Go-based `git-controller` for PHP developers.

## Requirements

- PHP 8.1+
- Git
- Composer (for autoloading)

## Setup

```bash
cd git-controller-php
composer install
```

## Usage

```bash
php main.php [command]
```

### Workflow Commands

| Command    | Description                                                        |
|------------|--------------------------------------------------------------------|
| `init`     | Scaffold a new workspace (system-definition.yaml, Makefile, .gitignore) |
| `clone`    | Clone all repositories defined in system-definition.yaml           |
| `fetch`    | Fetch all remotes across all repositories                          |
| `pull`     | Pull latest changes across all repositories (clone if missing)     |
| `push`     | Push local commits across all repositories                         |
| `scaffold` | Initialise .git and set remote origin (no clone/fetch needed)      |
| `checkout` | Switch all repositories to their defined branch                    |
| `status`   | Show dashboard overview of all repository states                   |
| `update`   | Fetch, pull, and show status across all repositories               |
| `validate` | Validate repository consistency against the definition             |

### Setup Commands

| Command     | Description                                    |
|-------------|------------------------------------------------|
| `doctor`    | Diagnose environment (Git, PHP, SSH, Docker)   |
| `ssh-setup` | Interactive SSH key management (alias: `ssh`)  |
| `help`      | Show help                                      |

## Architecture

The PHP version mirrors the Go version's structure:

```
git-controller-php/
├── main.php                    # Entry point (equivalent to main.go)
├── composer.json               # Autoloading & dependencies
├── src/
│   ├── Commands/               # One class per command
│   │   ├── CheckoutCommand.php
│   │   ├── CloneCommand.php
│   │   ├── DoctorCommand.php
│   │   ├── FetchCommand.php
│   │   ├── PullCommand.php
│   │   ├── PushCommand.php
│   │   ├── ScaffoldCommand.php
│   │   ├── SshSetupCommand.php
│   │   ├── StatusCommand.php
│   │   ├── UpdateCommand.php
│   │   ├── ValidateCommand.php
│   │   └── WsInitCommand.php
│   ├── GitUtil/
│   │   └── GitUtil.php         # Git operations (clone, pull, push, etc.)
│   ├── SshUtil/
│   │   └── SshUtil.php         # SSH agent, connectivity, config parsing
│   ├── System/
│   │   ├── Component.php       # Repository component model
│   │   ├── Parser.php          # YAML parser for system-definition.yaml
│   │   ├── SystemDefinition.php # System definition model
│   │   └── Workspace.php       # Workspace path resolution
│   └── UI/
│       └── UI.php              # Coloured terminal output helpers
```

## Notes

- The built-in YAML parser handles the `system-definition.yaml` format without external dependencies. If the `yaml` PHP extension is available, it will be used instead.
- SSH enforcement is identical to the Go version — HTTP(S) repository URLs are rejected.
- All commands produce the same coloured terminal output as the Go version.
