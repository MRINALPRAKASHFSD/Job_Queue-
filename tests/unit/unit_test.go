// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package unit_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/queue"
	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/retry"
)

func TestJob_StateTransitions(t *testing.T) {
	states := []queue.State{
		queue.StatePending,
		queue.StateLeased,
		queue.StateCompleted,
		queue.StateFailed,
		queue.StateDeadLetter,
	}

	for _, s := range states {
		if string(s) == "" {
			t.Errorf("empty state string encountered")
		}
	}
}

func TestJob_PriorityOrder(t *testing.T) {
	if queue.PriorityCritical <= queue.PriorityHigh {
		t.Errorf("expected Critical priority to be greater than High")
	}
	if queue.PriorityHigh <= queue.PriorityDefault {
		t.Errorf("expected High priority to be greater than Default")
	}
	if queue.PriorityDefault <= queue.PriorityLow {
		t.Errorf("expected Default priority to be greater than Low")
	}
}

func TestEnqueueRequest_JSONSerialization(t *testing.T) {
	req := queue.EnqueueRequest{
		Queue:          "orders",
		Priority:       queue.PriorityHigh,
		Payload:        json.RawMessage(`{"order_id":"123","amount":99.95}`),
		MaxAttempts:    5,
		ScheduledDelay: 10 * time.Second,
		IdempotencyKey: "order_req_123",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal enqueue request: %v", err)
	}

	var decoded queue.EnqueueRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal enqueue request: %v", err)
	}

	if decoded.Queue != "orders" || decoded.Priority != queue.PriorityHigh {
		t.Errorf("decoded request does not match original: %+v", decoded)
	}
}

func TestRetry_ExponentialPolicy(t *testing.T) {
	policy := retry.NewPolicy(100*time.Millisecond, 2*time.Second, 4, nil)

	if !policy.ShouldRetry(0) {
		t.Errorf("should retry at attempt 0")
	}
	if !policy.ShouldRetry(3) {
		t.Errorf("should retry at attempt 3")
	}
	if policy.ShouldRetry(4) {
		t.Errorf("should not retry at attempt 4 (max reached)")
	}
}
