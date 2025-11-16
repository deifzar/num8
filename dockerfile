# Build
FROM golang:1.25-alpine3.22 AS builder
RUN apk update && \ 
    apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' appuser
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with security flags
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o num8 .

# Release
FROM alpine:3.20
RUN apk upgrade --no-cache && \
    apk add --no-cache chromium ca-certificates tzdata && \
    apk --no-cache upgrade && \
    rm -rf /var/cache/apk/* && \
    adduser -D -g '' -s /bin/sh appuser
COPY --from=builder --chown=appuser:appuser /app/num8 /usr/local/bin/
# Security: Switch to non-root user
USER appuser
# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD num8 --help || exit 1
# Expose port (document the port used)
# EXPOSE 8000
CMD ["num8","help"]