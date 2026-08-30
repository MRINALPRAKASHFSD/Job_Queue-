// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package retry

import (
	"math/rand"
	"testing"
	"time"
)

func TestPolicy_ShouldRetry(t *testing.T) {
	policy := NewPolicy(1*time.Second, 60*time.Second, 3, nil)

	if !policy.ShouldRetry(0) {
		t.Errorf("expected should retry for attempt 0")
	}
	if !policy.ShouldRetry(2) {
		t.Errorf("expected should retry for attempt 2")
	}
	if policy.ShouldRetry(3) {
		t.Errorf("expected should NOT retry for attempt 3 when max is 3")
	}
	if policy.ShouldRetry(5) {
		t.Errorf("expected should NOT retry for attempt 5")
	}
}

func TestPolicy_FullJitterBoundaries(t *testing.T) {
	// Deterministic seed for reproducible testing
	rng := rand.New(rand.NewSource(42))
	base := 1 * time.Second
	max := 10 * time.Second
	policy := NewPolicy(base, max, 5, rng)

	for attempt := 0; attempt < 10; attempt++ {
		for i := 0; i < 50; i++ {
			delay := policy.FullJitter(attempt)
			if delay < 0 {
				t.Fatalf("delay cannot be negative: %v", delay)
			}
			if delay > max {
				t.Fatalf("delay %v exceeded max backoff %v on attempt %d", delay, max, attempt)
			}
		}
	}
}

func TestPolicy_DecorrelatedJitterBoundaries(t *testing.T) {
	rng := rand.New(rand.NewSource(101))
	base := 500 * time.Millisecond
	max := 5 * time.Second
	policy := NewPolicy(base, max, 5, rng)

	prev := base
	for i := 0; i < 20; i++ {
		next := policy.DecorrelatedJitter(prev)
		if next < base {
			t.Fatalf("next delay %v less than base %v", next, base)
		}
		if next > max {
			t.Fatalf("next delay %v greater than max %v", next, max)
		}
		prev = next
	}
}
