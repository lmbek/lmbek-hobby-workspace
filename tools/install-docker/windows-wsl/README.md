# Install Docker on Windows (via WSL)

One-shot script that installs Docker Engine natively inside a WSL (Ubuntu) distribution on Windows.

## What It Does

1. Updates the package index.
2. Installs prerequisites (`ca-certificates`, `curl`, `gnupg`, `lsb-release`).
3. Adds Docker's official GPG key and APT repository.
4. Installs Docker Engine, CLI, containerd, Buildx, and the Compose plugin.
5. Adds your user to the `docker` group so you can run Docker without `sudo`.
6. Starts the Docker service.

## Usage

Open your **WSL Ubuntu** terminal and run:

```bash
bash install-docker.sh
```

The script requires `sudo` — you will be prompted for your password.

## Post-Install

- **Close and reopen your WSL terminal** (or run `newgrp docker`) for the group change to take effect.
- Verify the installation:

```bash
docker --version
docker compose version
```

### Auto-Start Docker on WSL Launch

Docker's service does not start automatically in WSL. Add this to your `~/.bashrc` or `~/.zshrc`:

```bash
if service docker status 2>&1 | grep -q "is not running"; then
    sudo service docker start > /dev/null 2>&1
fi
```

To skip the `sudo` password prompt for the auto-start, run `sudo visudo` and add:

```
%docker ALL=(ALL) NOPASSWD: /usr/sbin/service docker *
```

## Docker Desktop Is No Longer Needed

This script installs Docker Engine directly inside WSL — **Docker Desktop is redundant** after running it. Docker Desktop is a GUI wrapper that runs its own WSL backend, but with Docker Engine installed natively in your WSL distribution you get the same `docker` and `docker compose` CLI commands without the overhead, licensing costs, or resource consumption of Docker Desktop.

You can safely uninstall Docker Desktop from Windows after confirming that `docker run hello-world` works inside your WSL terminal.

## Requirements

- Windows 10/11 with WSL 2 and an Ubuntu distribution.
- Internet access (to download packages from Docker's APT repository).
