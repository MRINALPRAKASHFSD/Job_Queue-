# Contributing to Reliable Job Engine

Thank you for your interest in contributing to `reliable-job-engine`. We are building a mission-critical, fault-tolerant background processing engine with strict distributed systems guarantees.

To maintain code quality, deterministic concurrency, and production reliability, all contributions must adhere to the standards outlined in this guide.

---

## Code of Conduct

All contributors and maintainers are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md). Please report any unacceptable behavior according to the instructions in that document.

---

## Getting Started: Fork and Clone Workflow

1. **Fork the Repository**:
   Fork `https://github.com/MRINALPRAKASHFSD/Job_Queue-` to your personal GitHub account.
2. **Clone Locally**:
   ```bash
   git clone https://github.com/<your-username>/Job_Queue-.git
   cd Job_Queue-
   ```
3. **Configure Upstream Remote**:
   ```bash
   git remote add upstream https://github.com/MRINALPRAKASHFSD/Job_Queue-.git
   git fetch upstream
   ```
4. **Create a Dedicated Feature Branch**:
   Never develop directly on `main`. Keep your local `main` synchronized with `upstream/main`.

---

## Branch Naming Conventions

Branch names should be descriptive, lowercase, hyphen-separated, and prefixed with the category of change:

- `feat/<short-description>`: New functionality or architectural capability.
- `fix/<short-description>`: Bug fix or concurrency edge-case resolution.
- `perf/<short-description>`: Performance optimization or query tuning.
- `refactor/<short-description>`: Code restructuring without behavioral change.
- `test/<short-description>`: New test cases or test harness enhancements.
- `docs/<short-description>`: Documentation additions or corrections.
- `chore/<short-description>`: Build tooling, dependencies, or repository maintenance.
- `security/<short-description>`: Security-related hardening or dependency updates.

**Examples:**
- `feat/atomic-job-claiming`
- `fix/stale-lease-acknowledgement`
- `test/exponential-backoff-jitter`
- `docs/fencing-tokens-architecture`

---

## Commit Message Conventions

This repository strictly enforces **Conventional Commits (v1.0.0)**. Commit messages must follow this structure:

```
<type>(<scope>): <short imperative summary>

[optional body explaining motivation, distributed systems trade-offs, and invariants]

[optional footer(s), e.g. Fixes #123]
```

### Allowed Types

- `feat`: A new feature for users or workers.
- `fix`: A bug fix or race condition remediation.
- `docs`: Documentation-only changes.
- `test`: Adding missing tests or correcting existing tests.
- `refactor`: A code change that neither fixes a bug nor adds a feature.
- `perf`: A code change that improves performance or reduces database lock contention.
- `chore`: Changes to build scripts, auxiliary tools, or non-source code.
- `ci`: Changes to CI/CD workflows and automation scripts.
- `security`: Security patches, vulnerability mitigations, or permission hardening.

### Common Scopes

- `queue`: Core job enqueuing, priority scheduling, and table locks.
- `worker`: Worker pool, leasing, execution loop, and heartbeating.
- `retry`: Backoff calculations, jitter algorithms, and attempt counters.
- `dlq`: Dead-letter queue isolation, failure logging, and manual replay.
- `idempotency`: Deduplication keys, transaction wrappers, and state fences.
- `inspector`: CLI and HTTP inspection endpoints.
- `db`: PostgreSQL migrations, connection pool, and query definitions.
- `config`: Environment and configuration loading.

### Verified Examples

- `feat(queue): implement atomic job claiming`
- `fix(worker): prevent stale lease acknowledgement`
- `test(retry): cover exponential backoff`
- `docs(architecture): document crash recovery`
- `perf(queue): add partial index on state and scheduled_at`

---

## Pull Request Requirements

Before opening a Pull Request:

1. **Link Related Issues**: Every non-trivial PR must reference an existing GitHub issue (e.g., `Closes #42` or `Addresses #15`).
2. **Atomic Changes**: Keep PRs focused on a single logical change. Large multi-faceted changes should be broken into sequential, reviewable PRs.
3. **Fill the PR Template**: Complete all sections of the [Pull Request Template](.github/PULL_REQUEST_TEMPLATE.md).
4. **Mandatory Tests**:
   - Every bug fix must include a regression test proving the failure and verifying the resolution.
   - Every new feature must include both unit tests and integration tests against real PostgreSQL.
5. **No Secrets or Extraneous Files**: Ensure `.env`, IDE files, build binaries, or debug logs are never included.

---

## Code Quality and Engineering Standards

### 1. Code Formatting
All Go source code must be formatted using standard tooling:
```bash
make fmt
# or manually:
gofmt -s -w .
```

### 2. Linting and Static Analysis
Code must pass standard Go static analysis without warnings:
```bash
make lint
# runs: go vet ./...
```
If `golangci-lint` is installed locally, run:
```bash
golangci-lint run ./...
```

### 3. Testing Requirements
We require rigorous testing, particularly for concurrency and lease mechanics:
- **Unit Tests**:
  ```bash
  make test-unit
  ```
- **Race Condition Detection**:
  All tests must run with the Go race detector enabled:
  ```bash
  go test -race ./...
  ```
- **PostgreSQL Integration Tests**:
  Mock databases cannot accurately simulate row-level locking (`FOR UPDATE SKIP LOCKED`), transactional deadlocks, or lease boundaries. All queue claiming and worker leasing tests must run against an active PostgreSQL instance:
  ```bash
  make test-integration
  ```

### 4. Documentation Requirements
- Public functions, types, and interfaces in Go packages must include meaningful doc comments adhering to Go standard conventions.
- If a change alters system guarantees, failure handling, or schema definitions, the corresponding documentation in `README.md`, `ARCHITECTURE.md`, or `FAILURE_MODES.md` must be updated within the same pull request.

---

## Security Guidelines

- Never hardcode credentials, connection strings, or encryption keys.
- Always use parameterized SQL queries (`$1, $2, ...`) without exception.
- If you discover a security vulnerability, do NOT open a public PR or issue. Follow the instructions in [SECURITY.md](SECURITY.md).

---

## Review Process

1. **Automated CI Validation**:
   Once opened, your PR will trigger GitHub Actions workflows:
   - Linting & Vet checks
   - Unit tests with `-race`
   - Real PostgreSQL integration tests
   - CodeQL security analysis
   - Dependency vulnerability review
2. **Architectural Review**:
   A repository maintainer will review the code for distributed systems correctness, locking behavior, edge-case safety, and code style.
3. **Revisions**:
   Address review feedback by pushing new commits to your branch. Do not force-push over active reviews unless requested.
4. **Merging**:
   Once approved and all CI checks pass, the PR will be merged using "Squash and merge" or "Rebase and merge" to maintain a clean linear commit history.
