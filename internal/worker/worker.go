// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package worker

import (
	"context"
	"errors"
	"time"

	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/queue"
)

// Sentinel errors for worker execution
var (
	ErrExecutionPanicked = errors.New("job execution panicked")
	ErrLeaseLost         = errors.New("worker lease lost or fencing token revoked")
)

// HandlerFunc defines the business logic function signature executed by workers.
type HandlerFunc func(ctx context.Context, job *queue.Job) error

// Router maps queue names to handler functions.
type Router interface {
	Register(queueName string, handler HandlerFunc)
	Handler(queueName string) (HandlerFunc, bool)
}

// Pool defines the worker concurrency orchestrator.
type Pool interface {
	// Start starts the worker routines and lease management loops.
	Start(ctx context.Context) error

	// Stop gracefully drains in-flight jobs within timeout.
	Stop(ctx context.Context) error

	// ID returns the unique identifier for this worker node.
	ID() string
}

// ExecutionContext passes runtime lease renewal context to handler execution.
type ExecutionContext struct {
	Job          *queue.Job
	WorkerID     string
	FencingToken int64
	ClaimedAt    time.Time
}
