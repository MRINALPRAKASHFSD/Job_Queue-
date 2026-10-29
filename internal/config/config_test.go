// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadFromEnv_Defaults(t *testing.T) {
	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("expected no error loading config, got %v", err)
	}

	if cfg.AppEnv != "development" {
		t.Errorf("expected AppEnv 'development', got '%s'", cfg.AppEnv)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expected server port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Database.Database != "reliable_job_engine" {
		t.Errorf("expected db name 'reliable_job_engine', got '%s'", cfg.Database.Database)
	}

	if cfg.Retry.MaxAttempts != 5 {
		t.Errorf("expected max attempts 5, got %d", cfg.Retry.MaxAttempts)
	}
}

func TestLoadFromEnv_Overrides(t *testing.T) {
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("RETRY_MAX_ATTEMPTS", "10")
	os.Setenv("WORKER_LEASE_DURATION", "45s")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("RETRY_MAX_ATTEMPTS")
		os.Unsetenv("WORKER_LEASE_DURATION")
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 9000 {
		t.Errorf("expected server port 9000, got %d", cfg.Server.Port)
	}

	if cfg.Retry.MaxAttempts != 10 {
		t.Errorf("expected max attempts 10, got %d", cfg.Retry.MaxAttempts)
	}

	if cfg.Worker.LeaseDuration != 45*time.Second {
		t.Errorf("expected lease duration 45s, got %v", cfg.Worker.LeaseDuration)
	}
}
