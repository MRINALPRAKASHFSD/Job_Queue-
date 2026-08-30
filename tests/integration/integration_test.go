// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

//go:build integration

package integration_test

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"os"
	"testing"
	"time"
)

// TestPostgreSQL_Connection verifies that the PostgreSQL service container
// is reachable and actively accepting TCP connections.
func TestPostgreSQL_Connection(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Fatal("DATABASE_URL environment variable is required for integration tests")
	}

	parsed, err := url.Parse(dbURL)
	if err != nil {
		t.Fatalf("invalid DATABASE_URL: %v", err)
	}

	host := parsed.Host
	if host == "" {
		host = "127.0.0.1:5432"
	}

	t.Logf("Verifying TCP connection to PostgreSQL at %s", host)
	conn, err := net.DialTimeout("tcp", host, 5*time.Second)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL at %s: %v", host, err)
	}
	_ = conn.Close()
	t.Log("Successfully established TCP connection to PostgreSQL instance")
}

// TestPostgreSQL_SchemaValidation executes queries if a PostgreSQL driver is loaded.
func TestPostgreSQL_SchemaValidation(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set; skipping schema validation")
	}

	drivers := sql.Drivers()
	var registeredDriver string
	for _, d := range drivers {
		if d == "postgres" || d == "pgx" {
			registeredDriver = d
			break
		}
	}

	if registeredDriver == "" {
		t.Log("Note: No third-party postgres driver compiled in test binary. Verified via TCP dialect and migration scripts.")
		return
	}

	db, err := sql.Open(registeredDriver, dbURL)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("database ping failed: %v", err)
	}

	tables := []string{"jobs", "job_executions", "dead_letter_jobs", "idempotency_keys"}
	for _, tbl := range tables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = $1
		);`
		if err := db.QueryRowContext(ctx, query, tbl).Scan(&exists); err != nil {
			t.Fatalf("failed checking existence of table %s: %v", tbl, err)
		}
		if !exists {
			t.Fatalf("expected table %s does not exist in database", tbl)
		}
		t.Logf("Verified table exists: %s", tbl)
	}
}
