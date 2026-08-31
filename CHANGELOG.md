# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Phase 2: PostgreSQL transactional claiming engine with `SELECT ... FOR UPDATE SKIP LOCKED`.
- Phase 3: Worker execution daemon with background heartbeat lease renewal and fencing token validation.
- Phase 4: Exponential backoff with truncated jitter and automated Dead-Letter Queue (DLQ) containment.
- Phase 5: Transactional idempotency keys and deduplication storage.
- Phase 6: Queue statistics and Prometheus metrics export.

---

## [0.1.0] - 2026-11-01

### Added
- **Repository Architecture Foundation**:
  - Comprehensive architectural blueprint in [ARCHITECTURE.md](ARCHITECTURE.md) covering atomic queue leasing, fencing tokens, idempotency, and transactional outbox.
  - Distributed failure taxonomy and mitigation matrix in [FAILURE_MODES.md](FAILURE_MODES.md).
  - Multi-phase project development plan in [ROADMAP.md](ROADMAP.md).
- **Core Domain & Package Layout**:
  - `internal/queue`: Core job domain types, job states, priority levels, and queue payload definitions.
  - `internal/retry`: Exponential backoff algorithm with full jitter and decorrelated jitter, including comprehensive unit test coverage.
  - `internal/worker`: Worker lease management, heartbeat interfaces, and executor contracts.
  - `internal/idempotency`: Idempotency key models and deduplication interfaces.
  - `internal/dlq`: Dead-letter queue inspection and replay contracts.
  - `internal/config`: Type-safe configuration loader with environment variable bindings.
  - `internal/database`: PostgreSQL connection pool lifecycle and transactional abstractions.
  - `cmd/server`, `cmd/worker`, `cmd/cli`: Entrypoint skeletons for API inspection, worker daemon, and management CLI.
- **Database Schema**:
  - Baseline PostgreSQL migration (`000001_init_schema.up.sql` and `down.sql`) defining `jobs`, `job_executions`, `dead_letter_jobs`, and `idempotency_keys` tables with partial indexes for high-throughput locking.
- **CI/CD & DevSecOps Infrastructure**:
  - GitHub Actions CI workflow with Go toolchain verification, `gofmt`, `go vet`, race detector (`-race`), and PostgreSQL container integration testing.
  - GitHub CodeQL static analysis for Go security vulnerabilities.
  - Dependency review action and automated Dependabot configuration for Go modules, GitHub Actions, and Docker.
  - Semantic release workflow template.
- **Developer Experience**:
  - Fully containerized local environment using `docker-compose.yml` with health-checked PostgreSQL 16.
  - Multi-stage `Dockerfile` with non-root security context.
  - Ergonomic `Makefile` for local testing, formatting, linting, and container lifecycle.
  - Verified `.gitignore`, `.gitattributes`, `.editorconfig`, and `.dockerignore`.
  - Comprehensive `.env.example` with zero hardcoded credentials.
- **Governance & Open Source Standards**:
  - MIT License and comprehensive `NOTICE` third-party attribution.
  - `CONTRIBUTING.md` with Conventional Commits rules and quality gates.
  - `CODE_OF_CONDUCT.md` adopting Contributor Covenant v2.1.
  - `SECURITY.md` defining private vulnerability reporting and SLA.
  - GitHub issue templates (Bug report, Feature request, Documentation) and Pull Request template.
  - `CODEOWNERS` referencing repository maintainer `@MRINALPRAKASHFSD`.
  - Academic / software citation metadata in `CITATION.cff`.
