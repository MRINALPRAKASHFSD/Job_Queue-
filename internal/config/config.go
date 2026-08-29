// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Version and Commit are populated at build time via -ldflags.
var (
	Version = "0.1.0"
	Commit  = "none"
)

// Config encapsulates runtime configuration for the server, worker fleet, and database.
type Config struct {
	AppEnv     string
	AppDebug   bool
	Server     ServerConfig
	Database   DatabaseConfig
	Worker     WorkerConfig
	Retry      RetryConfig
	DeadLetter DeadLetterConfig
	Metrics    MetricsConfig
}

// ServerConfig defines API server settings.
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig defines PostgreSQL connection parameters.
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DSN returns the standard PostgreSQL connection string.
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode,
	)
}

// WorkerConfig defines worker concurrency and lease timeouts.
type WorkerConfig struct {
	Concurrency       int
	PollInterval      time.Duration
	LeaseDuration     time.Duration
	HeartbeatInterval time.Duration
	ShutdownTimeout   time.Duration
	Queues            []string
}

// RetryConfig defines exponential backoff policies.
type RetryConfig struct {
	BaseBackoff time.Duration
	MaxBackoff  time.Duration
	MaxAttempts int
	JitterMode  string
}

// DeadLetterConfig defines DLQ behavior.
type DeadLetterConfig struct {
	Enabled        bool
	MaxReplayLimit int
}

// MetricsConfig defines telemetry settings.
type MetricsConfig struct {
	Enabled bool
	Port    int
}

// LoadFromEnv loads configuration from environment variables with sensible defaults.
func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		AppDebug: getEnvBool("APP_DEBUG", true),
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres_password_placeholder"),
			Database:        getEnv("DB_NAME", "reliable_job_engine"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		Worker: WorkerConfig{
			Concurrency:       getEnvInt("WORKER_CONCURRENCY", 10),
			PollInterval:      getEnvDuration("WORKER_POLL_INTERVAL", 500*time.Millisecond),
			LeaseDuration:     getEnvDuration("WORKER_LEASE_DURATION", 30*time.Second),
			HeartbeatInterval: getEnvDuration("WORKER_HEARTBEAT_INTERVAL", 10*time.Second),
			ShutdownTimeout:   getEnvDuration("WORKER_SHUTDOWN_TIMEOUT", 20*time.Second),
			Queues:            parseList(getEnv("WORKER_QUEUES", "critical,high,default,low")),
		},
		Retry: RetryConfig{
			BaseBackoff: getEnvDuration("RETRY_BASE_BACKOFF", 1*time.Second),
			MaxBackoff:  getEnvDuration("RETRY_MAX_BACKOFF", 300*time.Second),
			MaxAttempts: getEnvInt("RETRY_MAX_ATTEMPTS", 5),
			JitterMode:  getEnv("RETRY_JITTER_MODE", "full"),
		},
		DeadLetter: DeadLetterConfig{
			Enabled:        getEnvBool("DLQ_ENABLED", true),
			MaxReplayLimit: getEnvInt("DLQ_MAX_REPLAY_LIMIT", 3),
		},
		Metrics: MetricsConfig{
			Enabled: getEnvBool("METRICS_ENABLED", true),
			Port:    getEnvInt("METRICS_PORT", 9090),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.ParseBool(strings.TrimSpace(val)); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := time.ParseDuration(strings.TrimSpace(val)); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func parseList(val string) []string {
	parts := strings.Split(val, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}
