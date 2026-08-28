# Project Roadmap

This roadmap outlines the planned development phases for `reliable-job-engine`. Each phase builds upon the distributed systems guarantees established in preceding milestones.

---

## Milestone Overview

```
Phase 1: Repository Foundation ───────► [COMPLETED / BASELINE]
Phase 2: Durable Queue ───────────────► [NEXT IN PROGRESS]
Phase 3: Worker Reliability ──────────► [PLANNED]
Phase 4: Retry and DLQ ───────────────► [PLANNED]
Phase 5: Idempotency ─────────────────► [PLANNED]
Phase 6: Observability ───────────────► [PLANNED]
Phase 7: Advanced Features ───────────► [PLANNED]
```

---

## Detailed Milestone Breakdown

### Phase 1: Repository Foundation :white_check_mark:
**Target**: v0.1.0 (Current)
- [x] Professional repository foundation and licensing (MIT License, Notice, Authors, Citation).
- [x] System architecture blueprint ([ARCHITECTURE.md](ARCHITECTURE.md)) and failure analysis ([FAILURE_MODES.md](FAILURE_MODES.md)).
- [x] Security policies ([SECURITY.md](SECURITY.md)) and community code of conduct ([CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)).
- [x] Developer experience tooling: `Makefile`, `Dockerfile`, `docker-compose.yml`, `.env.example`, `.editorconfig`.
- [x] GitHub Actions CI workflow with Go vet, race detection (`-race`), and PostgreSQL container integration testing.
- [x] GitHub CodeQL security analysis, Dependabot update rules, and PR/Issue templates.
- [x] Core domain packages skeleton (`internal/queue`, `internal/retry`, `internal/worker`, `internal/idempotency`, `internal/dlq`, `internal/config`, `internal/database`).
- [x] Initial SQL schema migrations (`migrations/000001_init_schema.up.sql`).

---

### Phase 2: Durable Queue Storage Engine :hourglass_flowing_sand:
**Target**: v0.2.0
- [ ] Implement robust PostgreSQL connection pool wrapper using `pgx/v5` with connection health checks.
- [ ] Implement atomic job enqueueing API:
  - Immediate execution.
  - Delayed execution (`scheduled_at`).
  - Integer-based priority levels (Critical, High, Default, Low).
- [ ] Implement batch claiming engine using `SELECT ... FOR UPDATE SKIP LOCKED`:
  - Partial index optimization on `state` and `scheduled_at`.
  - Configurable batch size per claim invocation.
- [ ] Integration test suite for multi-worker claiming with zero contention and zero duplicate claims.

---

### Phase 3: Worker Reliability & Heartbeat Engine
**Target**: v0.3.0
- [ ] Concurrent worker execution pool with graceful shutdown (`SIGTERM`/`SIGINT`) and context cancellation.
- [ ] Time-bounded worker leases (`lease_expires_at`).
- [ ] Background heartbeat renewal daemon renewing leases at configurable fractions (e.g. $1/3$ of lease duration).
- [ ] Monotonically increasing fencing tokens on each lease claim and renewal to protect against zombie worker writes.
- [ ] Reaper / Watchdog daemon detecting expired leases and returning jobs to `pending` status.
- [ ] Panic recovery boundary preventing malformed jobs from crashing worker processes.

---

### Phase 4: Retry Strategy & Dead-Letter Queue (DLQ)
**Target**: v0.4.0
- [ ] Exponential backoff scheduler with Full Jitter and Decorrelated Jitter algorithms.
- [ ] Configurable base backoff, maximum backoff, and maximum retry attempts per queue.
- [ ] Quarantined dead-letter storage (`dead_letter_jobs` table) with structured failure traces.
- [ ] Dead-letter inspection and manual replay CLI commands (`jobctl dlq replay <job-id>`).
- [ ] Replay loop prevention with maximum replay thresholds.

---

### Phase 5: Idempotency & Transactional Guarantees
**Target**: v0.5.0
- [ ] Unique `idempotency_key` enforcement across queues.
- [ ] Transactional Outbox pattern helper utilities for producer microservices.
- [ ] Automatic duplicate payload short-circuiting: return cached terminal results if already processed.
- [ ] Side-effect coordination guards ensuring at-least-once deliveries map to exactly-once business results.

---

### Phase 6: Observability, Metrics & Management API
**Target**: v0.6.0
- [ ] Structured logging with Go `log/slog` containing contextual trace IDs, worker IDs, and execution latencies.
- [ ] Native Prometheus metrics endpoint (`/metrics`) exposing queue depths, throughput, and error rates.
- [ ] Real-time queue inspection HTTP server (`cmd/server`) with REST endpoints for queues and DLQs.
- [ ] Interactive CLI (`cmd/cli`) for operational administration (listing backlog, purging queues, replaying DLQ).

---

### Phase 7: Advanced Features & Ecosystem
**Target**: v1.0.0
- [ ] Directed Acyclic Graph (DAG) job workflows and sequential dependencies.
- [ ] High-throughput benchmarking suite comparing queue latencies under 1,000 to 50,000 concurrent jobs.
- [ ] Lightweight web administration dashboard for visual monitoring.
- [ ] Production deployment templates (Kubernetes manifests, Helm chart, Terraform module).
