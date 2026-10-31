// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting reliable-job-engine worker daemon",
		"version", config.Version,
		"commit", config.Commit,
	)

	cfg, err := config.LoadFromEnv()
	if err != nil {
		slog.Error("Failed to load environment configuration", "error", err)
		os.Exit(1)
	}

	hostname, _ := os.Hostname()
	workerID := "worker-" + hostname

	slog.Info("Worker daemon initialized (Foundation Phase 1)",
		"worker_id", workerID,
		"concurrency", cfg.Worker.Concurrency,
		"queues", cfg.Worker.Queues,
		"lease_duration", cfg.Worker.LeaseDuration.String(),
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("Worker daemon running. Awaiting job scheduler activation (Phase 3)...")

	<-ctx.Done()
	slog.Info("Shutting down worker daemon gracefully, draining tasks...")

	// Draining simulation
	time.Sleep(500 * time.Millisecond)
	slog.Info("Worker daemon drained and stopped cleanly")
}
