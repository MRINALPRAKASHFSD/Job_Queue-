-- ==============================================================================
-- Migration: 000001_init_schema.down.sql
-- Description: Drop core tables and indexes for Reliable Job Engine
-- Copyright (c) 2026 Mrinal Prakash
-- SPDX-License-Identifier: MIT
-- ==============================================================================

DROP TABLE IF EXISTS idempotency_keys CASCADE;
DROP TABLE IF EXISTS dead_letter_jobs CASCADE;
DROP TABLE IF EXISTS job_executions CASCADE;
DROP TABLE IF EXISTS jobs CASCADE;
