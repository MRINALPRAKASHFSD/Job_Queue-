-- ==============================================================================
-- Migration: 000001_init_schema.up.sql
-- Description: Core tables and indexes for Reliable Job Engine
-- Copyright (c) 2026 Mrinal Prakash
-- SPDX-License-Identifier: MIT
-- ==============================================================================

-- Enable pgcrypto for gen_random_uuid() if not already available
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ------------------------------------------------------------------------------
-- 1. Main Jobs Table
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    queue VARCHAR(128) NOT NULL DEFAULT 'default',
    state VARCHAR(32) NOT NULL DEFAULT 'pending',
    priority INT NOT NULL DEFAULT 0,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    attempt INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    worker_id VARCHAR(128) DEFAULT NULL,
    fencing_token BIGINT NOT NULL DEFAULT 0,
    lease_expires_at TIMESTAMPTZ DEFAULT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT check_job_state CHECK (
        state IN ('pending', 'leased', 'completed', 'failed', 'dead_letter')
    ),
    CONSTRAINT check_attempts_positive CHECK (
        attempt >= 0 AND max_attempts > 0
    )
);

-- Partial index for atomic claim operations (SELECT ... FOR UPDATE SKIP LOCKED)
CREATE INDEX IF NOT EXISTS idx_jobs_claimable
    ON jobs (queue, priority DESC, scheduled_at ASC)
    WHERE state = 'pending';

-- Partial index for reaper / watchdog detecting expired worker leases
CREATE INDEX IF NOT EXISTS idx_jobs_expired_leases
    ON jobs (lease_expires_at ASC)
    WHERE state = 'leased';

-- General index for queue monitoring and inspection
CREATE INDEX IF NOT EXISTS idx_jobs_queue_state
    ON jobs (queue, state, created_at DESC);

-- ------------------------------------------------------------------------------
-- 2. Execution Audit Log (Append-only)
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_executions (
    id BIGSERIAL PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    worker_id VARCHAR(128) NOT NULL,
    attempt INT NOT NULL,
    fencing_token BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL,
    error_message TEXT DEFAULT NULL,
    duration_ms BIGINT NOT NULL DEFAULT 0,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_job_executions_job_id
    ON job_executions (job_id, executed_at DESC);

-- ------------------------------------------------------------------------------
-- 3. Dead-Letter Queue (DLQ) Quarantine Table
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS dead_letter_jobs (
    id UUID PRIMARY KEY,
    queue VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    total_attempts INT NOT NULL,
    last_error TEXT NOT NULL,
    last_worker_id VARCHAR(128) DEFAULT NULL,
    replay_count INT NOT NULL DEFAULT 0,
    quarantined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_replayed_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_dlq_queue_quarantined
    ON dead_letter_jobs (queue, quarantined_at DESC);

-- ------------------------------------------------------------------------------
-- 4. Idempotency Keys Table
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS idempotency_keys (
    key VARCHAR(255) PRIMARY KEY,
    job_id UUID NOT NULL,
    result_payload JSONB DEFAULT NULL,
    state VARCHAR(32) NOT NULL DEFAULT 'processing',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT check_idempotency_state CHECK (
        state IN ('processing', 'completed', 'failed')
    )
);

CREATE INDEX IF NOT EXISTS idx_idempotency_expires_at
    ON idempotency_keys (expires_at ASC);
