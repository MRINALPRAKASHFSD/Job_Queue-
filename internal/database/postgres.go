// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/config"
)

// DB wraps standard database/sql connection pool with helper methods.
type DB struct {
	*sql.DB
}

// Open creates and verifies a PostgreSQL connection pool.
// Note: Driver registration (e.g. pgx or pq) occurs at program entrypoint.
func Open(cfg config.DatabaseConfig) (*DB, error) {
	// Standard database/sql open
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return &DB{DB: db}, nil
}

// PingWithTimeout validates database reachability within a specified timeout.
func (db *DB) PingWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return db.PingContext(ctx)
}

// WithinTransaction executes a closure inside an isolated SQL transaction.
// Automatically rolls back if an error or panic occurs, commits on success.
func (db *DB) WithinTransaction(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
