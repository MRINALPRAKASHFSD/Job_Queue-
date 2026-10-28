# Failure Modes and Recovery Analysis

Distributed background processing systems must operate with deterministic reliability in the presence of node crashes, network partitions, database restarts, and poison pill payloads.

This document details the failure modes anticipated in `reliable-job-engine`, their root causes, automated detection mechanisms, recovery procedures, and formal correctness guarantees.

---

## Failure Matrix Summary

| Failure Mode | Root Cause | Detection Mechanism | Recovery Procedure | Correctness Guarantee |
| :--- | :--- | :--- | :--- | :--- |
| **1. Worker Crash Before Execution** | Hardware fault, OOM kill, SIGKILL immediately after claiming job. | Lease expiration (`lease_expires_at < NOW()`). | Watchdog reaper resets job state to `pending`. | At-least-once delivery; no data loss. |
| **2. Worker Crash During Execution** | Handler panic, unhandled exception, machine shutdown midway through task. | Missed heartbeat; lease timeout exceeded. | Reaper releases lease back to pool; increment attempt. | At-least-once delivery; job resumed on another worker. |
| **3. Crash After Side-Effect Before Ack** | Worker succeeds in external API call (e.g. Stripe) but crashes before updating PostgreSQL. | Lease timeout; job reassigned to new worker. | New worker checks `idempotency_keys` table before executing side-effect. | Exactly-once business effect via idempotency keys. |
| **4. Database Restart / Failover** | PostgreSQL planned maintenance, crash, or Patroni/RDS primary failover. | Driver connection pool error; SQL query timeouts. | Workers enter exponential backoff connection retry; in-flight locks auto-expire. | ACID durability; committed jobs retained in WAL. |
| **5. Network Partition (Worker <-> DB)** | Switch failure, cloud VPC disruption, or firewall drop. | TCP keepalive failure; heartbeat queries return network timeout. | Worker terminates active job execution context; DB lease expires; other workers claim. | Fencing token prevents stale worker from committing after partition heals. |
| **6. Lease Expiration (GC Pause / CPU Spikes)** | Stop-the-world Go GC pause, CPU starvation, or handler taking longer than lease interval. | Database clock passes `lease_expires_at`; Reaper resets row. | Reaper bumps fencing token; new worker claims job. | Stale worker's eventual commit rejected due to mismatched fencing token. |
| **7. Duplicate Delivery** | Re-queueing after crash or transient network timeout on acknowledgement. | Worker receives job previously attempted (`attempt > 1`). | Worker executes inside transactional idempotency guard. | No duplicate side-effects. |
| **8. Retry Storms** | Downstream third-party API outage causing thousands of jobs to fail simultaneously. | Surge in error rate and `jobs_failed_total` metrics. | Exponential backoff with Full Jitter de-synchronizes subsequent retries. | Prevents thundering herd against downstream services. |
| **9. Poison Jobs (Corrupted Payload / Bug)** | Malformed JSON payload or bug causing unrecoverable panic on every attempt. | `attempt >= max_attempts` or explicit unrecoverable error. | Automatic isolation into `dead_letter_jobs` table; alerted via DLQ metric. | Queue unblocked; remaining valid jobs process normally. |
| **10. Stale / Zombie Worker** | Zombie worker awakens after long freeze and attempts to finish old job. | `UPDATE ... WHERE fencing_token = $expected` matches 0 rows. | Worker detects row count 0; aborts transaction and logs warning. | Zero split-brain data corruption. |
| **11. Dead-Letter Replay Loops** | Operator replays a DLQ job before fixing the underlying code bug. | Job fails again and increments DLQ replay counter. | Quarantined back to DLQ with audit trail if it fails maximum replay limit. | DLQ cannot cyclically exhaust queue worker capacity. |

---

## Detailed Failure Mode Analysis

### 1. Worker Crash Before Execution
- **Scenario**: Worker claims a batch of jobs via `SELECT ... FOR UPDATE SKIP LOCKED` and sets `state = 'leased'`. Before the worker begins dispatching the job handler, the process is killed (`SIGKILL` or OOM).
- **Detection**: The worker fails to send heartbeat updates. PostgreSQL evaluates `state = 'leased' AND lease_expires_at < NOW()`.
- **Recovery**: The background Reaper daemon identifies the expired lease and executes:
  ```sql
  UPDATE jobs
  SET state = 'pending',
      worker_id = NULL,
      lease_expires_at = NULL,
      fencing_token = fencing_token + 1,
      updated_at = NOW()
  WHERE state = 'leased'
    AND lease_expires_at < NOW();
  ```
- **Invariants**: The job is restored to `pending` and immediately becomes available to other healthy workers.

---

### 2. Worker Crash During Execution
- **Scenario**: A worker begins processing a computationally heavy or I/O-intensive task. Midway through, the host reboots or the process crashes.
- **Detection**: Heartbeat stops. Lease window expires.
- **Recovery**: The Reaper detects the expired lease. The attempt counter is incremented, and an execution audit record is logged with status `lease_timeout_abandoned`.
- **Invariants**: Work is never silently dropped. The job is scheduled for redelivery subject to `max_attempts`.

---

### 3. Crash After Side-Effect but Before Acknowledgement
- **Scenario**: The worker calls a payment processor (e.g., Stripe API), charges \$50, and receives HTTP 200. Before the worker can run `UPDATE jobs SET state = 'completed'`, the network interface drops or power fails.
- **Problem**: A new worker will eventually pick up the job and could charge the card a second time if not guarded.
- **Mitigation & Recovery**:
  1. The producer generates a deterministic `idempotency_key` (e.g., `payment_order_482910`).
  2. The job handler submits this `idempotency_key` as Stripe's `Idempotency-Key` HTTP header.
  3. If re-executed by Worker 2, Stripe returns the original cached response without re-charging.
  4. Furthermore, for database side-effects, the worker uses the `idempotency_keys` table inside an atomic transaction:
     ```sql
     INSERT INTO idempotency_keys (key, job_id, result_hash, created_at)
     VALUES ($1, $2, $3, NOW())
     ON CONFLICT (key) DO NOTHING;
     ```
- **Guarantee**: At-least-once physical invocation yields exactly-once logical side-effects.

---

### 4. Database Restart or Primary Failover
- **Scenario**: PostgreSQL undergoes a planned minor version restart, hardware fault, or automated failover in a high-availability cluster (e.g. AWS RDS Multi-AZ or Patroni).
- **Detection**: Go database driver (`pgx` or `database/sql`) returns connection refused, broken pipe, or context deadline exceeded.
- **Recovery**:
  - The worker pool pauses claiming loops and activates exponential backoff reconnection logic (1s, 2s, 4s, up to 30s).
  - Open transactions are automatically rolled back by PostgreSQL upon connection termination.
  - Active leases simply expire while the database recovers.
  - When the primary database returns, workers reconnect, and the Reaper recovers expired leases cleanly.
- **Durability Guarantee**: Jobs committed to WAL before the restart are fully preserved.

---

### 5. Network Partition Between Worker and Database
- **Scenario**: A network switch fails, severing the link between Worker Node A and PostgreSQL, while Worker Node B remains connected.
- **Detection**: Heartbeats from Worker A time out. PostgreSQL declares Worker A's lease expired.
- **Recovery**:
  1. Worker B reclaims the job and increments `fencing_token`.
  2. Worker A's database queries will fail immediately. Even if Worker A finishes local computation, any final commit attempt will fail because Worker A cannot reach PostgreSQL.
  3. When the partition heals, Worker A attempts to acknowledge using its old `fencing_token`. The query matches 0 rows and fails safely.
- **Guarantee**: Network partitions cannot produce divergent job states.

---

### 6. Lease Expiration due to Excessive CPU Load / GC Pauses
- **Scenario**: High memory allocation triggers a prolonged Go GC pause, or high host CPU load delays the worker's heartbeat goroutine past `lease_expires_at`.
- **Detection**: The database reaper identifies `lease_expires_at < NOW()`.
- **Recovery**: The reaper resets the job and increments `fencing_token = fencing_token + 1`. A second worker claims the job.
- **Stale Worker Defense**:
  When the frozen worker unfreezes, it attempts to complete the job:
  ```sql
  UPDATE jobs
  SET state = 'completed',
      completed_at = NOW()
  WHERE id = $1
    AND fencing_token = $2; -- Token from original claim
  ```
  Since the fencing token has been incremented, the statement affects `0` rows. The worker inspects `RowsAffected()`, realizes its lease was revoked, discards the result, and logs a warning.
- **Guarantee**: Zero duplicate state transitions or conflicting commits.

---

### 7. Retry Storms and Cascading Failures
- **Scenario**: A downstream dependency (e.g., an internal billing microservice or third-party web service) suffers an outage. 50,000 background jobs fail concurrently.
- **Risk**: If all 50,000 jobs retry at fixed intervals (e.g., exactly every 60 seconds), they will hammer the recovering service in synchronized spikes, preventing its recovery.
- **Mitigation**:
  1. The engine calculates retry intervals using **Full Jitter**:
     $$\text{delay} = \text{random\_uniform}(0, \min(300\text{s}, 1\text{s} \times 2^{\text{attempt}}))$$
  2. Retries are uniformly dispersed across the time spectrum, creating a smooth, flat request rate against downstream dependencies.

---

### 8. Poison Pill Jobs
- **Scenario**: A malformed payload triggers an unhandled memory panic or fatal logic error every time it is decoded.
- **Risk**: Without protection, the job could crash every worker in the fleet sequentially in an infinite loop.
- **Mitigation**:
  1. Worker execution loops wrap handler invocation in a deferred `recover()` block.
  2. If a panic occurs, the panic message and stack trace are captured as an error.
  3. The job increments `attempt`.
  4. Once `attempt >= max_attempts` (default: 5), the engine moves the row into `dead_letter_jobs` and marks `state = 'dead_letter'`.
  5. The active queue is unblocked, and operators receive an alert.

---

### 9. Dead-Letter Replay Loops
- **Scenario**: An operator uses the CLI to re-queue 1,000 dead-lettered jobs, but the underlying bug has not yet been deployed.
- **Mitigation**:
  - `dead_letter_jobs` tracks `replay_count`.
  - Replayed jobs retain their original lineage ID.
  - If a replayed job fails again after reaching `max_attempts`, it is returned to the DLQ, and `replay_count` is incremented.
  - Engines enforce a `max_replay_limit` (default: 3) to prevent infinite operator re-queue loops.
