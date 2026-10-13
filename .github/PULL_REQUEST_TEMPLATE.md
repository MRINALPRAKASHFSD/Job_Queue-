## Description

<!-- Provide a clear summary of the changes introduced by this pull request. Explain the problem solved and the architectural approach taken. -->

Closes #(issue number)

## Type of Change

<!-- Please select all options that apply: -->
- [ ] `feat`: New feature or capability
- [ ] `fix`: Bug fix or race condition resolution
- [ ] `perf`: Performance or query optimization
- [ ] `refactor`: Code refactoring without behavior change
- [ ] `test`: New tests or test infrastructure update
- [ ] `docs`: Documentation addition or correction
- [ ] `chore`: Build scripts, dependencies, or toolchain updates
- [ ] `security`: Security patch or hardening

## Architectural and Concurrency Checklist

<!-- For changes affecting queue claiming, worker leases, or database operations: -->
- [ ] Does this change affect database row locking (`SELECT FOR UPDATE SKIP LOCKED`)?
- [ ] Are worker leases, heartbeats, or fencing token invariants preserved?
- [ ] Are all database operations strictly parameterized to prevent SQL injection?
- [ ] Is error handling structured, with retryable vs non-retryable errors distinguished?
- [ ] Have failure modes and edge cases been evaluated against [FAILURE_MODES.md](FAILURE_MODES.md)?

## Testing Checklist

- [ ] Unit tests added / updated and passing (`make test-unit`)
- [ ] Race detector passed with no data races (`go test -race ./...`)
- [ ] PostgreSQL integration tests executed and passing (`make test-integration`)
- [ ] Code formatted with `gofmt` (`make fmt`)
- [ ] Go static analysis passes with no warnings (`make lint`)

## Security and Quality Assurance

- [ ] No secrets, credentials, `.env` files, or private keys are committed
- [ ] Documentation updated in `README.md`, `ARCHITECTURE.md`, or code comments where appropriate
- [ ] Conventional Commit message guidelines followed
