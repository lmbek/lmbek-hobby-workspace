#!/bin/bash
# ============================================================
# Docker Engine installer for Linux (Ubuntu/Debian)
# Run directly on your Linux machine:
#   bash install-docker.sh
# ============================================================
set -euo pipefail

echo "==> Updating package index..."
sudo apt update -y

echo "==> Installing prerequisites..."
sudo apt install -y ca-certificates curl gnupg lsb-release

echo "==> Adding Docker GPG key..."
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo "==> Adding Docker repository..."
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

echo "==> Installing Docker Engine..."
sudo apt update -y
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

echo "==> Adding $USER to docker group..."
sudo usermod -aG docker "$USER"

echo "==> Enabling and starting Docker service..."
sudo systemctl enable docker
sudo systemctl start docker

echo ""
echo "============================================================"
echo "  Docker installed successfully!"
echo "============================================================"
echo ""
echo "  Verify with:  docker --version && docker compose version"
echo ""
echo "  IMPORTANT: Log out and log back in (or run 'newgrp docker')"
echo "  so you can use docker without sudo."
echo "============================================================"
