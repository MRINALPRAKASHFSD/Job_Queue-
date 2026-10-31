# syntax=docker/dockerfile:1.4
# Multi-stage Dockerfile for reliable-job-engine

# ==============================================================================
# Stage 1: Build binary
# ==============================================================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

WORKDIR /app

# Cache dependencies
COPY go.mod ./
# Copy go.sum if it exists
RUN if [ -f go.sum ]; then cp go.sum ./; fi
RUN go mod download || true

# Copy full source tree
COPY . .

# Build target binary (server, worker, or cli)
ARG TARGET_BIN=server
ARG VERSION=0.1.0
ARG COMMIT=unknown

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X 'github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/config.Version=${VERSION}' -X 'github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/config.Commit=${COMMIT}'" \
    -o /app/bin/app ./cmd/${TARGET_BIN}

# ==============================================================================
# Stage 2: Minimal runtime image
# ==============================================================================
FROM alpine:3.24 AS runner

# Create non-root unprivileged user and group
RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

# Install TLS certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /home/appuser

# Copy executable from builder
COPY --from=builder --chown=appuser:appgroup /app/bin/app /usr/local/bin/app

USER appuser:appgroup

EXPOSE 8080 9090

ENTRYPOINT ["/usr/local/bin/app"]
