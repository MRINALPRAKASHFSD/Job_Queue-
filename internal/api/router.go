// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/config"
)

// Server provides HTTP REST endpoints for queue inspection, health checks, and metrics.
type Server struct {
	cfg *config.Config
	mux *http.ServeMux
}

// HealthResponse represents standard health check response.
type HealthResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Commit    string    `json:"commit"`
	Timestamp time.Time `json:"timestamp"`
}

// NewServer initializes the HTTP router and registers standard endpoints.
func NewServer(cfg *config.Config) *Server {
	s := &Server{
		cfg: cfg,
		mux: http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

// Handler returns the underlying http.Handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthCheck)
	s.mux.HandleFunc("GET /readyz", s.handleReadiness)
	s.mux.HandleFunc("GET /api/v1/version", s.handleVersion)

	// API placeholder endpoints (to be wired to Engine/Inspector in subsequent phases)
	s.mux.HandleFunc("POST /api/v1/jobs", s.handleEnqueuePlaceholder)
	s.mux.HandleFunc("GET /api/v1/queues", s.handleQueuesOverviewPlaceholder)
	s.mux.HandleFunc("GET /api/v1/dlq", s.handleDLQListPlaceholder)
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, HealthResponse{
		Status:    "pass",
		Version:   config.Version,
		Commit:    config.Commit,
		Timestamp: time.Now().UTC(),
	})
}

func (s *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	// In production, verify database ping
	respondJSON(w, http.StatusOK, map[string]string{
		"readiness": "ready",
	})
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"version": config.Version,
		"commit":  config.Commit,
	})
}

func (s *Server) handleEnqueuePlaceholder(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusNotImplemented, map[string]string{
		"error":   "not_implemented",
		"message": "Job enqueuing API planned for Phase 2",
	})
}

func (s *Server) handleQueuesOverviewPlaceholder(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusNotImplemented, map[string]string{
		"error":   "not_implemented",
		"message": "Queue overview inspection API planned for Phase 6",
	})
}

func (s *Server) handleDLQListPlaceholder(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusNotImplemented, map[string]string{
		"error":   "not_implemented",
		"message": "Dead-letter queue management API planned for Phase 4",
	})
}

func respondJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
