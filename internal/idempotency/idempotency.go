// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrDuplicateExecution = errors.New("operation with this idempotency key is already running or completed")
	ErrKeyNotFound        = errors.New("idempotency key not found")
)

// State represents the execution phase associated with an idempotency key.
type State string

const (
	StateProcessing State = "processing"
	StateCompleted  State = "completed"
	StateFailed     State = "failed"
)

// Record tracks execution deduplication metadata.
type Record struct {
	Key           string          `json:"key"`
	JobID         string          `json:"job_id"`
	ResultPayload json.RawMessage `json:"result_payload,omitempty"`
	State         State           `json:"state"`
	CreatedAt     time.Time       `json:"created_at"`
	ExpiresAt     time.Time       `json:"expires_at"`
}

// Store defines deduplication key operations.
type Store interface {
	// Acquire claims an idempotency key. Returns ErrDuplicateExecution if active or completed.
	Acquire(ctx context.Context, key string, jobID string, ttl time.Duration) error

	// Complete stores the final execution result and marks key as completed.
	Complete(ctx context.Context, key string, result json.RawMessage) error

	// Release removes or fails an in-progress idempotency key when execution errors.
	Release(ctx context.Context, key string) error
}
