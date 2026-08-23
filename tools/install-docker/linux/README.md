# Install Docker on Linux

One-shot script that installs Docker Engine on a native Linux machine (Ubuntu/Debian). No WSL required.

## What It Does

1. Updates the package index.
2. Installs prerequisites (`ca-certificates`, `curl`, `gnupg`, `lsb-release`).
3. Adds Docker's official GPG key and APT repository.
4. Installs Docker Engine, CLI, containerd, Buildx, and the Compose plugin.
5. Adds your user to the `docker` group so you can run Docker without `sudo`.
6. Enables and starts the Docker service via `systemctl`.

## Usage

```bash
bash install-docker.sh
```

The script requires `sudo` — you will be prompted for your password.

## Post-Install

- **Log out and log back in** (or run `newgrp docker`) for the group change to take effect.
- Verify the installation:

```bash
docker --version
docker compose version
docker run hello-world
```

Docker will start automatically on boot thanks to `systemctl enable docker`.

## Requirements

- Ubuntu or Debian-based Linux distribution.
- Internet access (to download packages from Docker's APT repository).
