#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
# Copyright (c) 2026 Mrinal Prakash

set -euo pipefail

ACTION="${1:-up}"
DB_URL="${DATABASE_URL:-postgres://postgres:postgres_password_placeholder@127.0.0.1:5432/reliable_job_engine?sslmode=disable}"

MIGRATION_DIR="migrations"

if [ ! -d "${MIGRATION_DIR}" ]; then
  echo "Error: Migration directory '${MIGRATION_DIR}' not found." >&2
  exit 1
fi

case "${ACTION}" in
  up)
    echo "==> Running schema migrations (UP)..."
    for file in "${MIGRATION_DIR}"/*.up.sql; do
      if [ -f "$file" ]; then
        echo "--> Applying $file..."
        if command -v psql >/dev/null 2>&1; then
          psql "${DB_URL}" -f "$file"
        elif docker ps | grep -q reliable_job_engine_postgres; then
          docker exec -i reliable_job_engine_postgres psql -U postgres -d reliable_job_engine < "$file"
        else
          echo "Error: Neither local psql nor running Docker container reliable_job_engine_postgres was found." >&2
          exit 1
        fi
      fi
    done
    echo "==> Migrations applied successfully."
    ;;
  down)
    echo "==> Rolling back schema migrations (DOWN)..."
    for file in $(ls -r "${MIGRATION_DIR}"/*.down.sql); do
      if [ -f "$file" ]; then
        echo "--> Applying rollback $file..."
        if command -v psql >/dev/null 2>&1; then
          psql "${DB_URL}" -f "$file"
        elif docker ps | grep -q reliable_job_engine_postgres; then
          docker exec -i reliable_job_engine_postgres psql -U postgres -d reliable_job_engine < "$file"
        else
          echo "Error: Neither local psql nor running Docker container reliable_job_engine_postgres was found." >&2
          exit 1
        fi
      fi
    done
    echo "==> Migrations rolled back successfully."
    ;;
  *)
    echo "Usage: $0 [up|down]" >&2
    exit 1
    ;;
esac
