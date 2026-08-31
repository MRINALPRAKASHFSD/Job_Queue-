// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/api"
	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting reliable-job-engine API server",
		"version", config.Version,
		"commit", config.Commit,
	)

	cfg, err := config.LoadFromEnv()
	if err != nil {
		slog.Error("Failed to load environment configuration", "error", err)
		os.Exit(1)
	}

	apiServer := api.NewServer(cfg)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      apiServer.Handler(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown channel
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("HTTP server listening", "address", addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-shutdownChan
	slog.Info("Received termination signal, starting graceful shutdown", "signal", sig.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Forced HTTP server shutdown with error", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped cleanly")
}
