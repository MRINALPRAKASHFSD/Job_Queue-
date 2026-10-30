// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package retry

import (
	"math"
	"math/rand"
	"time"
)

// Policy defines parameters for exponential retry calculation.
type Policy struct {
	BaseBackoff time.Duration
	MaxBackoff  time.Duration
	MaxAttempts int
	rng         *rand.Rand
}

// NewPolicy creates a retry policy. If rng is nil, a default random source is used.
func NewPolicy(base, max time.Duration, maxAttempts int, rng *rand.Rand) *Policy {
	if rng == nil {
		// Use local seeded RNG
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return &Policy{
		BaseBackoff: base,
		MaxBackoff:  max,
		MaxAttempts: maxAttempts,
		rng:         rng,
	}
}

// ShouldRetry returns true if current attempt is less than max allowed attempts.
func (p *Policy) ShouldRetry(attempt int) bool {
	return attempt < p.MaxAttempts
}

// FullJitter calculates exponential delay with uniform random jitter:
// sleep = min(MaxBackoff, BaseBackoff * 2^attempt)
// delay = random_uniform(0, sleep)
func (p *Policy) FullJitter(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	// Compute exponential factor 2^attempt, capped to prevent overflow
	multiplier := math.Pow(2, float64(attempt))
	maxSleepSec := p.BaseBackoff.Seconds() * multiplier

	if maxSleepSec > p.MaxBackoff.Seconds() || math.IsInf(maxSleepSec, 0) {
		maxSleepSec = p.MaxBackoff.Seconds()
	}

	maxSleepNanos := int64(maxSleepSec * float64(time.Second))
	if maxSleepNanos <= 0 {
		return p.BaseBackoff
	}

	jitteredNanos := p.rng.Int63n(maxSleepNanos)
	return time.Duration(jitteredNanos)
}

// DecorrelatedJitter calculates decorrelated jitter delay based on the previous delay:
// delay = min(MaxBackoff, random_uniform(BaseBackoff, prevDelay * 3))
func (p *Policy) DecorrelatedJitter(prevDelay time.Duration) time.Duration {
	minNanos := p.BaseBackoff.Nanoseconds()
	maxNanos := prevDelay.Nanoseconds() * 3
	if maxNanos < minNanos {
		maxNanos = minNanos * 3
	}

	capNanos := p.MaxBackoff.Nanoseconds()
	if maxNanos > capNanos {
		maxNanos = capNanos
	}

	delta := maxNanos - minNanos
	if delta <= 0 {
		return p.BaseBackoff
	}

	return time.Duration(minNanos + p.rng.Int63n(delta))
}
