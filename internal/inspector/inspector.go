// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package inspector

import (
	"context"
	"time"

	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/queue"
)

// QueueStats provides aggregated point-in-time metrics for a specific queue.
type QueueStats struct {
	QueueName      string    `json:"queue_name"`
	PendingCount   int64     `json:"pending_count"`
	LeasedCount    int64     `json:"leased_count"`
	CompletedCount int64     `json:"completed_count"`
	FailedCount    int64     `json:"failed_count"`
	DeadLetterCount int64    `json:"dead_letter_count"`
	OldestPending  *time.Time `json:"oldest_pending,omitempty"`
}

// ActiveLease describes a job currently claimed by a worker.
type ActiveLease struct {
	JobID          string    `json:"job_id"`
	Queue          string    `json:"queue"`
	WorkerID       string    `json:"worker_id"`
	FencingToken   int64     `json:"fencing_token"`
	Attempt        int       `json:"attempt"`
	LeaseExpiresAt time.Time `json:"lease_expires_at"`
	SecondsRemaining float64 `json:"seconds_remaining"`
}

// Inspector defines read-only observational queries for operators and dashboards.
type Inspector interface {
	// Overview returns aggregate statistics across all known queues.
	Overview(ctx context.Context) ([]QueueStats, error)

	// QueueDetail returns fine-grained metrics for a single queue.
	QueueDetail(ctx context.Context, queueName string) (*QueueStats, error)

	// ActiveLeases returns currently leased jobs nearing expiration or actively heartbeating.
	ActiveLeases(ctx context.Context, queueName string) ([]ActiveLease, error)

	// ListJobs returns paginated jobs filtered by state.
	ListJobs(ctx context.Context, queueName string, state queue.State, limit, offset int) ([]queue.Job, error)
}
