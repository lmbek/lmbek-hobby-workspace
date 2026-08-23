#!/bin/bash
# ============================================================
# Docker Engine installer for Windows via WSL (Ubuntu)
# Run this inside your WSL terminal:
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

echo "==> Starting Docker service..."
sudo service docker start

echo ""
echo "============================================================"
echo "  Docker installed successfully!"
echo "============================================================"
echo ""
echo "  Verify with:  docker --version && docker compose version"
echo ""
echo "  IMPORTANT: Close and reopen your WSL terminal (or run"
echo "  'newgrp docker') so you can use docker without sudo."
echo ""
echo "  To auto-start Docker on WSL launch, add this to ~/.bashrc:"
echo '  if service docker status 2>&1 | grep -q "is not running"; then'
echo '      sudo service docker start > /dev/null 2>&1'
echo '  fi'
echo "============================================================"
