// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package dlq

import (
	"context"
	"encoding/json"
	"time"
)

// DeadLetterJob represents a quarantined job that has exhausted retries.
type DeadLetterJob struct {
	ID             string          `json:"id"`
	Queue          string          `json:"queue"`
	Payload        json.RawMessage `json:"payload"`
	TotalAttempts  int             `json:"total_attempts"`
	LastError      string          `json:"last_error"`
	LastWorkerID   *string         `json:"last_worker_id,omitempty"`
	ReplayCount    int             `json:"replay_count"`
	QuarantinedAt  time.Time       `json:"quarantined_at"`
	LastReplayedAt *time.Time      `json:"last_replayed_at,omitempty"`
}

// Filter provides query parameters for inspecting quarantined items.
type Filter struct {
	Queue  string
	Limit  int
	Offset int
}

// Manager defines operational contracts for Dead-Letter Queue administration.
type Manager interface {
	// List returns paginated dead letter jobs matching the filter.
	List(ctx context.Context, filter Filter) ([]DeadLetterJob, error)

	// Get retrieves a specific dead-letter job by ID.
	Get(ctx context.Context, id string) (*DeadLetterJob, error)

	// Replay moves a dead-letter job back to pending status in its original queue.
	Replay(ctx context.Context, id string) error

	// Purge permanently removes jobs quarantined prior to a specified timestamp.
	Purge(ctx context.Context, olderThan time.Time) (int64, error)
}
