#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
# Copyright (c) 2026 Mrinal Prakash

set -euo pipefail

echo "========================================================================"
echo "Reliable Job Engine - Local Development Setup"
echo "========================================================================"

# Check for .env file
if [ ! -f .env ]; then
  if [ -f .env.example ]; then
    echo "==> Creating .env from .env.example..."
    cp .env.example .env
    echo "==> Created .env (remember to review parameters if needed)."
  else
    echo "Error: .env.example not found!" >&2
    exit 1
  fi
else
  echo "==> Existing .env file found; preserving."
fi

# Check Docker
if command -v docker >/dev/null 2>&1; then
  echo "==> Docker is installed."
else
  echo "Warning: Docker is not installed or not in PATH. Local container services will not run." >&2
fi

# Check Docker Compose
if docker compose version >/dev/null 2>&1; then
  echo "==> Docker Compose is available."
else
  echo "Warning: Docker Compose is not available." >&2
fi

# Check Go
if command -v go >/dev/null 2>&1; then
  GO_VERSION=$(go version)
  echo "==> Go is installed: ${GO_VERSION}"
else
  echo "Warning: Go compiler is not currently in PATH. Please install Go 1.24+ to compile locally." >&2
fi

echo "========================================================================"
echo "Setup complete. Start local PostgreSQL via: make docker-up"
echo "========================================================================"
