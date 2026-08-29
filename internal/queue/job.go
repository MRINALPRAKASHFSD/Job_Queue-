// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package queue

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// Common queue errors
var (
	ErrJobNotFound     = errors.New("job not found")
	ErrLeaseExpired    = errors.New("job lease has expired or fencing token is invalid")
	ErrInvalidState    = errors.New("invalid job state transition")
	ErrPayloadTooLarge = errors.New("payload exceeds maximum allowable size")
)

// State represents the lifecycle stage of a background job.
type State string

const (
	StatePending    State = "pending"
	StateLeased     State = "leased"
	StateCompleted  State = "completed"
	StateFailed     State = "failed"
	StateDeadLetter State = "dead_letter"
)

// Priority defines standard job execution priorities.
type Priority int

const (
	PriorityLow      Priority = -10
	PriorityDefault  Priority = 0
	PriorityHigh     Priority = 10
	PriorityCritical Priority = 100
)

// Job represents a durable background task stored in PostgreSQL.
type Job struct {
	ID             string          `json:"id"`
	Queue          string          `json:"queue"`
	State          State           `json:"state"`
	Priority       Priority        `json:"priority"`
	Payload        json.RawMessage `json:"payload"`
	Attempt        int             `json:"attempt"`
	MaxAttempts    int             `json:"max_attempts"`
	WorkerID       *string         `json:"worker_id,omitempty"`
	FencingToken   int64           `json:"fencing_token"`
	LeaseExpiresAt *time.Time      `json:"lease_expires_at,omitempty"`
	ScheduledAt    time.Time       `json:"scheduled_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	CompletedAt    *time.Time      `json:"completed_at,omitempty"`
}

// EnqueueRequest encapsulates parameters to schedule a new job.
type EnqueueRequest struct {
	Queue          string          `json:"queue"`
	Priority       Priority        `json:"priority"`
	Payload        json.RawMessage `json:"payload"`
	MaxAttempts    int             `json:"max_attempts"`
	ScheduledDelay time.Duration   `json:"scheduled_delay,omitempty"`
	IdempotencyKey string          `json:"idempotency_key,omitempty"`
}

// ClaimedJob represents a job leased to a worker with fencing details.
type ClaimedJob struct {
	Job
	ClaimedAt time.Time `json:"claimed_at"`
}

// Engine defines the primary contract for durable queue operations.
type Engine interface {
	// Enqueue adds a job to the queue.
	Enqueue(ctx context.Context, req EnqueueRequest) (*Job, error)

	// ClaimBatch atomically claims up to limit pending jobs using SELECT FOR UPDATE SKIP LOCKED.
	ClaimBatch(ctx context.Context, queueName string, workerID string, limit int, leaseDuration time.Duration) ([]ClaimedJob, error)

	// Heartbeat extends the lease for an active job, verifying the worker and fencing token.
	Heartbeat(ctx context.Context, jobID string, workerID string, fencingToken int64, extension time.Duration) error

	// AcknowledgeSuccess marks a job as completed.
	AcknowledgeSuccess(ctx context.Context, jobID string, workerID string, fencingToken int64) error

	// AcknowledgeFailure increments the attempt, applies backoff, or moves to DLQ.
	AcknowledgeFailure(ctx context.Context, jobID string, workerID string, fencingToken int64, err error, nextAttemptAt time.Time) error
}
