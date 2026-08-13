# Getting Support

Thank you for using `reliable-job-engine`. We are committed to maintaining a high-quality, fault-tolerant background processing system.

## Documentation and Self-Service

Before opening a support request, please consult the project documentation:

- [README.md](README.md): System overview, quickstart, local development, and architecture fundamentals.
- [ARCHITECTURE.md](ARCHITECTURE.md): Deep-dive into database schema, atomic claiming (`SELECT FOR UPDATE SKIP LOCKED`), worker leases, fencing tokens, and idempotency guarantees.
- [FAILURE_MODES.md](FAILURE_MODES.md): Exhaustive failure matrix detailing crash recovery, network partitions, lease timeouts, and dead-letter handling.
- [ROADMAP.md](ROADMAP.md): Milestone progress and upcoming architectural capabilities.

---

## Where to Ask Questions

### 1. Questions, Architecture Discussions, and Ideas
If you have general questions on how to integrate `reliable-job-engine` into your Go microservices or how to design custom workers:
- Open a discussion topic in the **GitHub Discussions** area of this repository (or use the Question issue template).
- Provide context regarding your database environment (PostgreSQL version, connection pool configuration) and concurrency scale.

### 2. Bug Reports
If you suspect you have discovered a bug, race condition, or lease inconsistency:
- Review existing issues to ensure the bug has not already been reported.
- Open a new issue using our structured **[Bug Report Template](.github/ISSUE_TEMPLATE/bug_report.yml)**.
- Include a minimal reproducible example, reproduction steps, Go version, and PostgreSQL logs where applicable.

### 3. Feature Requests
To suggest architectural enhancements, new driver integrations, or observability endpoints:
- Open an issue using our structured **[Feature Request Template](.github/ISSUE_TEMPLATE/feature_request.yml)**.
- Explain the distributed systems use case, proposed API shape, and reliability implications.

### 4. Security Vulnerabilities
For security vulnerabilities, please **DO NOT** open a public issue. Follow the instructions outlined in our **[Security Policy](SECURITY.md)**.
