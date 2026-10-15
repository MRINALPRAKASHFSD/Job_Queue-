# Makefile for reliable-job-engine
# Copyright (c) 2026 Mrinal Prakash
# SPDX-License-Identifier: MIT

SHELL := /bin/bash
.DEFAULT_GOAL := help

# Build variables
BIN_DIR := bin
VERSION ?= 0.1.0
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w \
	-X 'github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/config.Version=$(VERSION)' \
	-X 'github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/config.Commit=$(GIT_COMMIT)'

# Database test connection
TEST_DB_URL ?= postgres://postgres:postgres_password_placeholder@127.0.0.1:5432/reliable_job_engine?sslmode=disable

.PHONY: help
help: ## Display available targets
	@echo "reliable-job-engine build targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: fmt lint test build ## Run formatting, linting, tests, and build all binaries

.PHONY: fmt
fmt: ## Format Go source code with gofmt
	@echo "==> Formatting code..."
	gofmt -s -w .

.PHONY: lint
lint: ## Run Go static analysis (vet)
	@echo "==> Running static analysis..."
	go vet ./...

.PHONY: test
test: test-unit test-race ## Run unit tests and race detection

.PHONY: test-unit
test-unit: ## Run unit tests
	@echo "==> Running unit tests..."
	go test -v ./internal/... ./tests/unit/...

.PHONY: test-race
test-race: ## Run tests with Go race detector
	@echo "==> Running tests with race detector..."
	go test -race -v ./internal/... ./tests/unit/...

.PHONY: test-integration
test-integration: ## Run PostgreSQL integration tests
	@echo "==> Running PostgreSQL integration tests..."
	DATABASE_URL="$(TEST_DB_URL)" go test -race -v -tags=integration ./tests/integration/...

.PHONY: build
build: build-server build-worker build-cli ## Build all binaries (server, worker, cli)

.PHONY: build-server
build-server: ## Build API server binary
	@echo "==> Building server binary..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/server ./cmd/server

.PHONY: build-worker
build-worker: ## Build worker daemon binary
	@echo "==> Building worker binary..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/worker ./cmd/worker

.PHONY: build-cli
build-cli: ## Build jobctl CLI binary
	@echo "==> Building jobctl CLI binary..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/jobctl ./cmd/cli

.PHONY: docker-up
docker-up: ## Start PostgreSQL and local development services in background
	@echo "==> Starting local infrastructure with Docker Compose..."
	docker compose up -d postgres

.PHONY: docker-down
docker-down: ## Stop all Docker Compose services and clean volumes
	@echo "==> Stopping Docker Compose services..."
	docker compose down

.PHONY: docker-build
docker-build: ## Build container images for server and worker
	@echo "==> Building Docker images..."
	docker compose build

.PHONY: migrate-up
migrate-up: ## Apply initial database schema migrations
	@echo "==> Applying migrations..."
	@if command -v psql >/dev/null 2>&1; then \
		psql "$(TEST_DB_URL)" -f migrations/000001_init_schema.up.sql; \
	else \
		docker exec -i reliable_job_engine_postgres psql -U postgres -d reliable_job_engine < migrations/000001_init_schema.up.sql; \
	fi

.PHONY: migrate-down
migrate-down: ## Rollback initial database schema migrations
	@echo "==> Rolling back migrations..."
	@if command -v psql >/dev/null 2>&1; then \
		psql "$(TEST_DB_URL)" -f migrations/000001_init_schema.down.sql; \
	else \
		docker exec -i reliable_job_engine_postgres psql -U postgres -d reliable_job_engine < migrations/000001_init_schema.down.sql; \
	fi

.PHONY: check-secrets
check-secrets: ## Scan repository for accidentally committed private keys or unencrypted secrets
	@echo "==> Checking for uncommitted secret files..."
	@git status --porcelain | grep -E '\.env$$|\.key$$|\.pem$$' && echo "ERROR: Potential secret file detected!" && exit 1 || echo "No secret files staged."

.PHONY: clean
clean: ## Remove compiled binaries and test coverage files
	@echo "==> Cleaning build artifacts..."
	rm -rf $(BIN_DIR) *.out *.test coverage.html
