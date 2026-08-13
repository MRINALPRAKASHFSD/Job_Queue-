# Security Policy

The `reliable-job-engine` project takes software security, data isolation, and distributed system stability seriously. This document outlines our vulnerability disclosure policy and security engineering guidelines.

## Supported Versions

Only the latest minor release receives active security fixes and patches.

| Version | Supported          | Security Maintenance Status |
| ------- | ------------------ | --------------------------- |
| 0.1.x   | :white_check_mark: | Active development & fixes  |
| < 0.1.0 | :x:                | Unsupported                 |

---

## Reporting a Vulnerability

**DO NOT report security vulnerabilities through public GitHub issues, pull requests, discussions, or social channels.**

To report a vulnerability responsibly:

1. **GitHub Private Vulnerability Reporting (Preferred)**:
   Navigate to the repository's **Security** tab and click **"Report a vulnerability"** to open an encrypted advisory draft directly with the maintainers.
2. **Alternative Method**:
   If Private Vulnerability Reporting is unavailable in your view, initiate contact through the project maintainer's verified GitHub profile:
   [@MRINALPRAKASHFSD](https://github.com/MRINALPRAKASHFSD).

### Information to Include in Your Report

To help us triage and remediate the issue quickly, please provide:
- A clear description of the vulnerability and its potential impact (e.g., SQL injection, unauthorized job lease renewal, denial of service, secret leakage).
- Step-by-step reproduction instructions or a minimal Proof of Concept (PoC).
- Specific component(s) affected (`internal/queue`, `internal/api`, `migrations`, etc.).
- Proposed mitigation or patch if available.

### Disclosure Timeline and SLA

- **Initial Response & Acknowledgement**: Within 48 hours of receipt.
- **Triage & Severity Assessment**: Within 5 business days.
- **Fix Preparation & Patch Release**: Targeted within 14 business days, depending on complexity.
- **Public Disclosure**: Coordinated after a fix has been verified and released via a patch version.

---

## Security Engineering Principles

### 1. Database Security
- **Parameterized Queries**: All database interactions MUST strictly utilize parameterized SQL queries (`$1, $2, ...`) or prepared statements. Direct string concatenation or raw SQL interpolation is strictly forbidden.
- **Principle of Least Privilege**: Production worker processes must connect using dedicated database roles that only have `SELECT`, `INSERT`, `UPDATE`, `DELETE` privileges on the `jobs` and `job_executions` tables. Schema migration credentials (`DDL`) must remain segregated.
- **Transport Encryption**: Production database connections must enforce TLS (`sslmode=verify-full` or `sslmode=require`).

### 2. Secrets Management
- **Zero Secrets in Repository**: No credentials, API tokens, passwords, private keys (`*.pem`, `*.key`), or production `.env` files may ever be committed.
- Automated pre-commit checks and CI scans (e.g., CodeQL, Secret scanning) enforce secret exclusion.
- All secrets must be supplied at runtime via environment variables or external secret managers (e.g., HashiCorp Vault, AWS Secrets Manager, Kubernetes Secrets).

### 3. Sensitive Payload Handling
- Payloads stored in `jobs.payload` may contain customer data. The engine treats payloads as opaque JSON (`jsonb`).
- Applications requiring compliance (e.g., PCI-DSS, HIPAA, GDPR) should apply envelope encryption or client-side payload encryption prior to enqueuing.
- Worker logs must sanitize and mask sensitive keys (e.g., `password`, `token`, `secret`, `authorization`, `credit_card`) during tracing and error logging.

### 4. Dependency Vulnerability Management
- Automated dependency scanning is managed via GitHub Dependabot and Dependency Review action.
- Dependencies are verified against checksums in `go.sum` (`go mod verify`).
- Transitive dependencies are audited regularly using Go vulnerability scanning (`govulncheck`) and CodeQL.

### 5. API Authentication and Authorization
- Public-facing or intranet API endpoints exposing job inspection and dead-letter replay (`cmd/server`) must be shielded behind authentication mechanisms (e.g., mTLS, JWT, or internal gateway proxies).
- Fencing tokens and worker leases must ensure that only active leaseholders can acknowledge or extend job runs, preventing unauthorized lease stealing.
