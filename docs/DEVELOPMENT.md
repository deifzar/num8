# NuM8 Development Guide

**Version:** 1.0  
**Last Updated:** January 2025  
**Target Audience:** Developers, Contributors, DevOps Engineers

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Project Structure](#project-structure)
3. [Coding Standards](#coding-standards)
4. [Development Workflow](#development-workflow)
5. [Testing Guidelines](#testing-guidelines)
6. [Database Development](#database-development)
7. [Security Guidelines](#security-guidelines)
8. [Performance Guidelines](#performance-guidelines)
9. [Debugging and Troubleshooting](#debugging-and-troubleshooting)
10. [Deployment Guidelines](#deployment-guidelines)

## Development Environment Setup

### Prerequisites

**Required Software:**
- Go 1.21.5+ ([installation guide](https://golang.org/doc/install))
- PostgreSQL 14+ ([installation guide](https://www.postgresql.org/docs/))
- RabbitMQ 3.8+ ([installation guide](https://www.rabbitmq.com/download.html))
- Docker and Docker Compose ([installation guide](https://docs.docker.com/get-docker/))
- Git ([installation guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git))

**External Tools:**
- Nuclei v3.3.5+ ([installation guide](https://github.com/projectdiscovery/nuclei))
- Burp Suite Professional (optional, for Burpmate integration)

### Local Environment Setup

1. **Clone the repository:**
```bash
git clone <repository-url>
cd NuM8
```

2. **Set up environment variables:**
```bash
cp .env.example .env
# Edit .env with your local configuration
```

3. **Start dependencies with Docker:**
```bash
docker-compose up -d postgres rabbitmq
```

4. **Initialize the database:**
```bash
# Create database and schema
createdb cptm8
psql -d cptm8 -f scripts/schema.sql
```

5. **Install Go dependencies:**
```bash
go mod download
go mod tidy
```

6. **Build and run the application:**
```bash
go build -o num8 main.go
./num8 launch --ip 127.0.0.1 --port 8080
```

### Development Tools

**Recommended IDE Setup:**
- **VS Code** with Go extension
- **GoLand** by JetBrains
- **Vim/Neovim** with vim-go plugin

**Essential Go Tools:**
```bash
# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/sast-scan@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
```

**Pre-commit Hooks Setup:**
```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Run hooks manually
pre-commit run --all-files
```

## Project Structure

### Package Organization

```
NuM8/
├── main.go                    # Application entry point
├── configuration.yaml         # Configuration file (DO NOT commit secrets)
├── .env.example              # Environment variables template
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── Dockerfile                # Container build instructions
├── docker-compose.yml        # Local development stack
│
├── pkg/                      # Application packages
│   ├── api8/                 # HTTP API layer
│   │   ├── api8.go          # API server setup
│   │   ├── routes.go        # Route definitions
│   │   └── middleware.go    # HTTP middleware
│   │
│   ├── controller8/          # Business logic controllers
│   │   ├── controller8_numate.go           # Nuclei integration with manual ACK
│   │   ├── controller8_numate_interface.go # Controller interface (RunNumate signature)
│   │   ├── controller8_burpmate.go         # Burp Suite integration
│   │   └── base_controller.go              # Common controller logic
│   │
│   ├── orchestrator8/        # Workflow orchestration
│   │   ├── orchestrator8.go           # Main orchestration with ACK/NACK methods
│   │   ├── orchestrator8_interface.go # Interface with ACK methods
│   │   └── workflow.go                # Workflow definitions
│   │
│   ├── model8/               # Data models and structures
│   │   ├── endpoint8.go      # Endpoint models
│   │   ├── historyissues.go  # Security issue models
│   │   └── notification8.go  # Notification models
│   │
│   ├── db8/                  # Database layer
│   │   ├── db8.go           # Database connection and operations
│   │   ├── migrations/       # Database migration scripts
│   │   └── queries/          # SQL query definitions
│   │
│   ├── amqpM8/              # Message queue operations with connection pooling
│   │   ├── connection_pool.go  # Connection pool implementation
│   │   ├── pool_manager.go     # Global pool manager (singleton)
│   │   ├── pooled_amqp.go      # Pooled AMQP operations with manual ACK
│   │   ├── initialization.go   # Pool initialization from config
│   │   └── shared_state.go     # Shared state management
│   │
│   ├── notification8/        # Notification system
│   │   ├── notification8.go  # Notification logic
│   │   ├── channels/         # Notification channels
│   │   └── templates/        # Message templates
│   │
│   ├── configparser/         # Configuration management
│   │   ├── configparser.go   # Configuration parsing
│   │   └── validation.go     # Configuration validation
│   │
│   └── log8/                 # Logging utilities
│       ├── log8.go          # Logging setup and configuration
│       └── formatters.go     # Log formatters
│
├── tests/                    # Test files
│   ├── unit/                # Unit tests
│   ├── integration/         # Integration tests
│   ├── api/                 # API tests
│   └── fixtures/            # Test data and fixtures
│
├── scripts/                  # Development and deployment scripts
│   ├── build.sh             # Build scripts
│   ├── test.sh              # Test scripts
│   ├── migrate.sh           # Database migration scripts
│   └── deploy.sh            # Deployment scripts
│
├── docs/                     # Documentation
│   ├── API.md               # API documentation
│   ├── ARCHITECTURE.md      # System architecture
│   ├── DEVELOPMENT.md       # This file
│   ├── PERFORMANCE.md       # Performance analysis
│   └── TODO.md              # Issues and improvements
│
└── deployments/              # Deployment configurations
    ├── docker/              # Docker configurations
    ├── kubernetes/          # Kubernetes manifests
    └── helm/                # Helm charts
```

### Package Naming Conventions

- **Package names:** Use descriptive names with the `8` suffix (e.g., `api8`, `db8`)
- **File names:** Use snake_case for multi-word names
- **Interface names:** Use descriptive names ending with interface behavior (e.g., `Scanner`, `Publisher`)
- **Struct names:** Use PascalCase with descriptive names

## Coding Standards

### Go Style Guide

Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html) guidelines.

**Key Principles:**
1. **Simplicity:** Write simple, clear, and readable code
2. **Consistency:** Follow established patterns in the codebase
3. **Performance:** Write efficient code, but prioritize readability
4. **Error Handling:** Handle all errors explicitly and appropriately

### Code Formatting

**Automatic Formatting:**
```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .

# Run linter
golangci-lint run
```

**Pre-commit Configuration (.pre-commit-config.yaml):**
```yaml
repos:
  - repo: local
    hooks:
      - id: go-fmt
        name: go fmt
        entry: go fmt
        language: system
        files: \.go$
      
      - id: go-imports
        name: go imports
        entry: goimports
        language: system
        files: \.go$
        args: [-w]
      
      - id: golangci-lint
        name: golangci-lint
        entry: golangci-lint run
        language: system
        files: \.go$
```

### Error Handling Standards

**DO:**
```go
// Return errors explicitly
func (d *DB8) GetEndpoint(id int) (*model8.Endpoint8, error) {
    var endpoint model8.Endpoint8
    err := d.db.QueryRow("SELECT * FROM endpoints WHERE id = $1", id).Scan(&endpoint)
    if err != nil {
        return nil, fmt.Errorf("failed to get endpoint %d: %w", id, err)
    }
    return &endpoint, nil
}

// Use context for cancellation
func (s *Scanner) ScanWithContext(ctx context.Context, target string) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        return s.performScan(target)
    }
}
```

**DON'T:**
```go
// Don't use log.Fatal in business logic
func (d *DB8) Connect() {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err) // DON'T DO THIS
    }
}

// Don't ignore errors
func (s *Scanner) Scan(target string) {
    result, _ := s.performScan(target) // DON'T IGNORE ERRORS
    s.saveResult(result)
}
```

### Logging Standards

**Use Structured Logging:**
```go
import "github.com/rs/zerolog/log"

// Good logging
log.Info().
    Str("endpoint", endpoint.URL).
    Int("scan_id", scanID).
    Dur("duration", time.Since(start)).
    Msg("Scan completed successfully")

// Error logging with context
log.Error().
    Err(err).
    Str("function", "ProcessScan").
    Int("endpoint_id", endpointID).
    Msg("Failed to process scan")
```

### Interface Design

**Define Minimal Interfaces:**
```go
// Good - focused interface
type Scanner interface {
    Scan(ctx context.Context, target string) (*ScanResult, error)
}

// Good - testable interface
type DatabaseManager interface {
    GetEndpoint(id int) (*model8.Endpoint8, error)
    SaveScanResult(*ScanResult) error
}

// Avoid - too broad interface
type MegaInterface interface {
    Scan(string) error
    SaveResult(*ScanResult) error
    SendNotification(*Notification) error
    ParseConfig() error
    // ... many more methods
}
```

## Development Workflow

### Git Workflow

**Branch Naming Convention:**
- `feature/description` - New features
- `bugfix/description` - Bug fixes
- `hotfix/description` - Critical fixes
- `refactor/description` - Code refactoring

**Commit Message Format:**
```
type(scope): brief description

Longer description if needed

- List of changes
- Important notes
- Breaking changes

Fixes #123
```

**Example Commit Messages:**
```
feat(api): add authentication middleware

- Implement JWT token validation
- Add role-based access control
- Update API documentation

Breaking change: All endpoints now require authentication

Fixes #45

fix(db): resolve connection pool exhaustion

- Implement proper connection pooling
- Add connection health checks
- Fix goroutine leaks in database layer

Closes #67

refactor(orchestrator): improve error handling

- Replace log.Fatal with proper error returns
- Add context-based cancellation
- Improve test coverage

No breaking changes
```

### Code Review Process

**Before Submitting PR:**
1. Run all tests: `go test ./...`
2. Run linter: `golangci-lint run`
3. Check test coverage: `go test -cover ./...`
4. Update documentation if needed
5. Add tests for new functionality

**PR Checklist:**
- [ ] Code follows style guidelines
- [ ] Tests added and passing
- [ ] Documentation updated
- [ ] No security vulnerabilities introduced
- [ ] Performance impact considered
- [ ] Breaking changes documented

### Development Tasks

**Daily Development Routine:**
```bash
# Pull latest changes
git pull origin develop

# Run tests
go test ./...

# Check for issues
golangci-lint run

# Build application
go build -o num8 main.go

# Run application
./num8 launch --ip 127.0.0.1 --port 8080
```

## Testing Guidelines

### Test Structure

**Test File Organization:**
```go
package controller8

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)

func TestController8Numate_ProcessScan(t *testing.T) {
    // Test cases
    tests := []struct {
        name     string
        input    *ScanRequest
        expected *ScanResult
        wantErr  bool
    }{
        {
            name: "successful scan",
            input: &ScanRequest{
                Target: "example.com",
                Type:   "vuln",
            },
            expected: &ScanResult{
                Status: "completed",
                Issues: 5,
            },
            wantErr: false,
        },
        {
            name: "invalid target",
            input: &ScanRequest{
                Target: "",
                Type:   "vuln",
            },
            expected: nil,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            controller := NewController8Numate()
            
            // Execute
            result, err := controller.ProcessScan(context.Background(), tt.input)
            
            // Assert
            if tt.wantErr {
                require.Error(t, err)
                assert.Nil(t, result)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.expected.Status, result.Status)
                assert.Equal(t, tt.expected.Issues, result.Issues)
            }
        })
    }
}
```

### Test Categories

**Unit Tests:**
- Test individual functions and methods
- Use mocks for external dependencies
- Fast execution (< 1 second per test)

**Integration Tests:**
- Test component interactions
- Use test database and message queues
- Moderate execution time (< 10 seconds per test)

**API Tests:**
- Test HTTP endpoints end-to-end
- Use test environment
- Include authentication and authorization tests

**Performance Tests:**
- Benchmark critical paths
- Load testing for scalability
- Resource usage validation

### Test Data Management

**Test Fixtures:**
```go
// tests/fixtures/endpoints.go
package fixtures

var TestEndpoints = []model8.Endpoint8{
    {
        ID:       1,
        URL:      "https://example.com",
        Protocol: "https",
        Port:     443,
    },
    {
        ID:       2,
        URL:      "http://test.local",
        Protocol: "http",
        Port:     80,
    },
}
```

**Database Test Helpers:**
```go
// tests/helpers/database.go
package helpers

func SetupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("postgres", "postgres://test:test@localhost/test_cptm8")
    require.NoError(t, err)
    
    // Run migrations
    err = RunMigrations(db)
    require.NoError(t, err)
    
    t.Cleanup(func() {
        db.Close()
    })
    
    return db
}
```

## Database Development

### Schema Management

**Migration Files:**
```sql
-- migrations/001_initial_schema.up.sql
CREATE TABLE IF NOT EXISTS domains (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS endpoints (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER REFERENCES domains(id),
    url TEXT NOT NULL,
    protocol VARCHAR(10) NOT NULL,
    port INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_endpoints_domain_id ON endpoints(domain_id);
CREATE INDEX idx_endpoints_url ON endpoints(url);
```

**Migration Commands:**
```bash
# Apply migrations
./scripts/migrate.sh up

# Rollback migrations
./scripts/migrate.sh down

# Check migration status
./scripts/migrate.sh status
```

### Query Optimization

**Use Prepared Statements:**
```go
// Good - prepared statement
const getEndpointQuery = `
    SELECT id, url, protocol, port, created_at 
    FROM endpoints 
    WHERE domain_id = $1 AND active = true`

func (d *DB8) GetEndpointsByDomain(domainID int) ([]model8.Endpoint8, error) {
    stmt, err := d.db.Prepare(getEndpointQuery)
    if err != nil {
        return nil, fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()
    
    rows, err := stmt.Query(domainID)
    // ... rest of implementation
}
```

**Database Connection Best Practices:**
```go
func (d *DB8) InitConnectionPool() error {
    db, err := sql.Open("postgres", d.connectionString)
    if err != nil {
        return err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(2 * time.Minute)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }
    
    d.db = db
    return nil
}
```

## Security Guidelines

### Secure Coding Practices

**Input Validation:**
```go
import "github.com/go-playground/validator/v10"

type ScanRequest struct {
    Target   string `json:"target" validate:"required,url"`
    ScanType string `json:"scan_type" validate:"required,oneof=vuln spider proxy"`
    Options  map[string]interface{} `json:"options" validate:"dive"`
}

func (s *ScanRequest) Validate() error {
    validate := validator.New()
    return validate.Struct(s)
}
```

**SQL Injection Prevention:**
```go
// Always use parameterized queries
func (d *DB8) GetEndpointByURL(url string) (*model8.Endpoint8, error) {
    query := `SELECT id, url, protocol, port FROM endpoints WHERE url = $1`
    var endpoint model8.Endpoint8
    err := d.db.QueryRow(query, url).Scan(
        &endpoint.ID, &endpoint.URL, &endpoint.Protocol, &endpoint.Port)
    if err != nil {
        return nil, err
    }
    return &endpoint, nil
}
```

**Secret Management:**
```go
// Use environment variables for secrets
func LoadConfig() (*Config, error) {
    config := &Config{
        DatabaseURL:    os.Getenv("DATABASE_URL"),
        RabbitMQURL:    os.Getenv("RABBITMQ_URL"),
        DiscordToken:   os.Getenv("DISCORD_TOKEN"),
        JWTSecret:      os.Getenv("JWT_SECRET"),
    }
    
    // Validate required secrets
    if config.DatabaseURL == "" {
        return nil, errors.New("DATABASE_URL is required")
    }
    
    return config, nil
}
```

### Authentication Implementation

**JWT Middleware:**
```go
func JWTAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractToken(r)
        if token == "" {
            http.Error(w, "missing token", http.StatusUnauthorized)
            return
        }
        
        claims, err := validateToken(token)
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## Performance Guidelines

### Optimization Patterns

**Worker Pool Pattern:**
```go
type WorkerPool struct {
    jobs    chan Job
    results chan Result
    workers int
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        jobs:    make(chan Job, workers*2),
        results: make(chan Result, workers*2),
        workers: workers,
    }
}

func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        go wp.worker(ctx)
    }
}

func (wp *WorkerPool) worker(ctx context.Context) {
    for {
        select {
        case job := <-wp.jobs:
            result := job.Process()
            wp.results <- result
        case <-ctx.Done():
            return
        }
    }
}
```

**HTTP Client Optimization:**
```go
var httpClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxConnsPerHost:     30,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        DisableKeepAlives:   false,
    },
}
```

### Memory Management

**Batch Processing:**
```go
func (p *Processor) ProcessLargeDataset(data []Item) error {
    const batchSize = 100
    
    for i := 0; i < len(data); i += batchSize {
        end := i + batchSize
        if end > len(data) {
            end = len(data)
        }
        
        batch := data[i:end]
        if err := p.processBatch(batch); err != nil {
            return fmt.Errorf("failed to process batch %d-%d: %w", i, end, err)
        }
        
        // Optional: trigger GC for large datasets
        if i%1000 == 0 {
            runtime.GC()
        }
    }
    
    return nil
}
```

## Debugging and Troubleshooting

### Debugging Tools

**pprof Integration:**
```go
import (
    _ "net/http/pprof"
    "net/http"
)

func init() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

**Debug Commands:**
```bash
# Memory profiling
go tool pprof http://localhost:6060/debug/pprof/heap

# CPU profiling
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Goroutine analysis
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Check for race conditions
go run -race main.go
```

### Logging for Debugging

**Structured Debug Logging:**
```go
log.Debug().
    Str("function", "ProcessScan").
    Str("target", target).
    Int("scan_id", scanID).
    Msg("Starting scan process")

log.Debug().
    Str("query", query).
    Interface("params", params).
    Msg("Executing database query")
```

### Common Issues and Solutions

**Database Connection Issues:**
```bash
# Check connection pool status
SELECT count(*) FROM pg_stat_activity WHERE datname = 'cptm8';

# Monitor connection usage
watch -n 1 "netstat -an | grep :5432 | wc -l"
```

**RabbitMQ Issues:**
```bash
# Check queue status
rabbitmqctl list_queues

# Monitor message flow
rabbitmqctl list_exchanges

# Check connection status
rabbitmqctl list_connections

# Check unacknowledged messages
rabbitmqctl list_queues name messages_ready messages_unacknowledged

# Inspect connection pool health
# Check logs for: "Consumer health", "Connection pool", "ACK/NACK"
```

## RabbitMQ Manual Acknowledgment Patterns

### Overview

NuM8 uses **manual message acknowledgment** (manual ACK) for RabbitMQ message processing to ensure message reliability and prevent message loss during scan operations.

**Key Benefits:**
- Messages are only removed from queue after successful scan completion
- Failed scans can be automatically retried via NACK with requeue
- Prevents message loss during application crashes or panics
- Provides fine-grained control over message lifecycle

### Configuration

**Set autoACK in configuration.yaml:**
```yaml
ORCHESTRATORM8:
  num8:
    Queue:
      - "num8_consumer"    # Consumer name
      - "num8.scan.queue"  # Queue name
      - "false"            # autoACK (false = manual ACK mode)
```

**Default behavior:** If autoACK parsing fails, defaults to `false` (manual ACK).

### Implementation Pattern

**1. Message Handler (Orchestrator Layer):**

The orchestrator receives RabbitMQ messages and extracts the `deliveryTag`:

```go
func (o *Orchestrator8) createHandleAPICallByServiceWithConnection(service string) {
    handler := func(msg amqp.Delivery) error {
        // Extract routing key and build HTTP request
        requestURL := fmt.Sprintf("%s/%s", serviceURL, endpoint)

        req, err := http.NewRequest("POST", requestURL, bytes.NewReader(msg.Body))
        if err != nil {
            return err
        }

        // CRITICAL: Pass deliveryTag via HTTP header
        req.Header.Set("X-RabbitMQ-Delivery-Tag", fmt.Sprintf("%d", msg.DeliveryTag))
        req.Header.Set("X-RabbitMQ-Consumer-Tag", msg.ConsumerTag)

        resp, err := client.Do(req)
        if err != nil {
            return err // Handler error triggers NACK with requeue
        }
        defer resp.Body.Close()

        return nil // Success - wait for scan completion to ACK
    }

    // Register handler
    conn.RegisterHandler(queueName, handler)
}
```

**2. Controller Layer (Extract Delivery Tag):**

The controller extracts the `deliveryTag` from the HTTP request header:

```go
func (m *Controller8Numate) NumateScan(c *gin.Context) {
    // Extract delivery tag from request header
    deliveryTagStr := c.GetHeader("X-RabbitMQ-Delivery-Tag")
    var deliveryTag uint64
    if deliveryTagStr != "" {
        if tag, err := strconv.ParseUint(deliveryTagStr, 10, 64); err == nil {
            deliveryTag = tag
            log8.BaseLogger.Debug().Msgf("Scan triggered via RabbitMQ (deliveryTag: %d)", deliveryTag)
        }
    }

    // Perform validation and setup...

    // Launch async scan with deliveryTag
    go m.RunNumate(true, endpoints, options, outputFile, deliveryTag)

    c.JSON(http.StatusOK, gin.H{"msg": "Scan started"})
}
```

**3. Scan Execution with Defer ACK/NACK:**

The scan execution tracks completion and acknowledges the message:

```go
func (m *Controller8Numate) RunNumate(fullscan bool, endpoints []model8.Endpoint8,
    options model8.Model8Options8Interface, outputFileName string, deliveryTag uint64) {

    var scanCompleted bool = false

    defer func() {
        if r := recover(); r != nil {
            log8.BaseLogger.Error().Msgf("Panic during scan: %v", r)
            scanCompleted = false // Mark as failed
        }

        // ACK or NACK based on scan completion
        if deliveryTag > 0 {
            ackErr := m.Orch.AckScanCompletion(deliveryTag, scanCompleted)
            if ackErr != nil {
                log8.BaseLogger.Error().Msgf("Failed to ACK/NACK message (deliveryTag: %d): %v",
                    deliveryTag, ackErr)
            }
        }
    }()

    // Perform scan...
    err := performNucleiScan(endpoints, options, outputFileName)
    if err != nil {
        log8.BaseLogger.Error().Msgf("Scan failed: %v", err)
        scanCompleted = false
        return
    }

    // Mark as completed on success
    scanCompleted = true
}
```

**4. Orchestrator ACK/NACK Methods:**

```go
// AckScanCompletion acknowledges or rejects a message based on scan outcome
func (o *Orchestrator8) AckScanCompletion(deliveryTag uint64, scanCompleted bool) error {
    return amqpM8.WithPooledConnection(func(conn amqpM8.PooledAmqpInterface) error {
        ch := conn.GetChannel()
        if ch == nil {
            return fmt.Errorf("channel is nil, cannot acknowledge message")
        }

        if !scanCompleted {
            // Scan failed - NACK with requeue for retry
            log8.BaseLogger.Warn().Msgf("Scan incomplete (deliveryTag: %d) - NACK with requeue", deliveryTag)
            return ch.Nack(deliveryTag, false, true) // requeue=true
        }

        // Scan succeeded - ACK to remove from queue
        log8.BaseLogger.Info().Msgf("Scan completed (deliveryTag: %d) - ACK", deliveryTag)
        return ch.Ack(deliveryTag, false)
    })
}

// NackScanMessage rejects a message (for permanent failures)
func (o *Orchestrator8) NackScanMessage(deliveryTag uint64, requeue bool) error {
    return amqpM8.WithPooledConnection(func(conn amqpM8.PooledAmqpInterface) error {
        ch := conn.GetChannel()
        if ch == nil {
            return fmt.Errorf("channel is nil, cannot nack message")
        }

        log8.BaseLogger.Warn().Msgf("Rejecting message (deliveryTag: %d, requeue: %v)",
            deliveryTag, requeue)
        return ch.Nack(deliveryTag, false, requeue)
    })
}
```

### ACK/NACK Decision Matrix

| Scenario | Action | Requeue | Rationale |
|----------|--------|---------|-----------|
| Scan completed successfully | ACK | N/A | Message processed, remove from queue |
| Scan failed (runtime error) | NACK | true | Temporary failure, retry scan |
| Scan panicked | NACK | true | Unexpected error, retry scan |
| Configuration error | NACK | false | Permanent failure, send to DLQ |
| Validation error | NACK | false | Invalid input, don't retry |
| No handler found | NACK | false | Configuration issue, don't retry |
| Handler error (HTTP) | NACK | true | Network issue, retry |
| Database connection error | NACK | true | Temporary DB issue, retry |

### Consumer Handler Behavior

**Handler Success (No Immediate ACK):**
```go
if err := handler(msg); err != nil {
    // Handler failed - NACK immediately with requeue
    if !autoACK {
        log8.BaseLogger.Warn().Msgf("Handler failed, NACKing (deliveryTag: %d, requeue: true)",
            msg.DeliveryTag)
        if nackErr := msg.Nack(false, true); nackErr != nil {
            log8.BaseLogger.Error().Msgf("Failed to NACK message: %v", nackErr)
        }
    }
} else {
    // Handler succeeded - DON'T ACK yet!
    // Wait for scan completion (defer function will ACK)
    log8.BaseLogger.Debug().Msgf("Handler succeeded, waiting for scan completion to ACK (deliveryTag: %d)",
        msg.DeliveryTag)
}
```

### Best Practices

**DO:**
1. Always extract `deliveryTag` from HTTP headers in RabbitMQ-triggered requests
2. Pass `deliveryTag` to async goroutines that perform long-running operations
3. Use defer functions to ensure ACK/NACK is called even on panic
4. Set `scanCompleted` flag only after all operations succeed
5. Use `NACK(requeue=true)` for temporary failures (network, DB)
6. Use `NACK(requeue=false)` for permanent failures (validation, config)
7. Log all ACK/NACK operations with delivery tags for troubleshooting

**DON'T:**
1. Don't ACK immediately after handler success in manual ACK mode
2. Don't ignore `deliveryTag` - it's critical for message tracking
3. Don't use `log.Fatal()` in scan operations - it prevents defer ACK/NACK
4. Don't requeue messages indefinitely - implement retry limits
5. Don't ACK failed scans - use NACK with requeue instead

### Monitoring and Troubleshooting

**Check for stuck messages:**
```bash
# List queues with unacknowledged messages
rabbitmqctl list_queues name messages_ready messages_unacknowledged consumers

# Expected output:
# num8.scan.queue  0  1  1  (1 message being processed)
# num8.scan.queue  0  0  1  (idle, no messages)
```

**Debug log patterns:**
```bash
# Grep for ACK/NACK operations
grep "deliveryTag" log/num8.log

# Expected patterns:
# "Scan triggered via RabbitMQ (deliveryTag: 123)"
# "Handler succeeded, waiting for scan completion to ACK (deliveryTag: 123)"
# "Scan completed successfully (deliveryTag: 123) - sending ACK"
# "Scan incomplete (deliveryTag: 123) - sending NACK with requeue"
```

**Common issues:**
- **Messages stuck in "unacknowledged" state:** Check if scans are hanging or if ACK/NACK is not being called
- **Messages requeuing infinitely:** Implement retry limits or dead letter queue
- **Channel is nil errors:** Connection pool issue - check health checks
- **Duplicate message processing:** Ensure idempotency in scan operations

### Migration from Auto-ACK to Manual ACK

**Before (Auto-ACK):**
```yaml
Queue:
  - "consumer_name"
  - "queue_name"
  - "true"  # autoACK = true (messages auto-acknowledged)
```

**After (Manual ACK):**
```yaml
Queue:
  - "consumer_name"
  - "queue_name"
  - "false"  # autoACK = false (manual acknowledgment)
```

**Code changes required:**
1. Update controller signatures to accept `deliveryTag uint64` parameter
2. Add defer function in scan execution to ACK/NACK messages
3. Implement `AckScanCompletion` and `NackScanMessage` methods in orchestrator
4. Add `X-RabbitMQ-Delivery-Tag` header propagation in HTTP handlers
5. Update interface definitions to include new method signatures

## Deployment Guidelines

### Local Development Deployment

**Docker Compose Setup:**
```yaml
version: '3.8'
services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: cptm8
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  rabbitmq:
    image: rabbitmq:3-management
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: password
    ports:
      - "5672:5672"
      - "15672:15672"

  num8:
    build: .
    environment:
      DATABASE_URL: postgres://postgres:password@postgres:5432/cptm8
      RABBITMQ_URL: amqp://admin:password@rabbitmq:5672/
    ports:
      - "8003:8003"
    depends_on:
      - postgres
      - rabbitmq

volumes:
  postgres_data:
```

### Production Deployment

**Environment Variables:**
```bash
# Database configuration
export DATABASE_URL="postgres://user:pass@localhost:5432/cptm8"
export DATABASE_MAX_CONNECTIONS="25"
export DATABASE_MAX_IDLE="25"

# RabbitMQ configuration
export RABBITMQ_URL="amqp://user:pass@localhost:5672/"

# Security configuration
export JWT_SECRET="your-secret-key"
export API_KEY="your-api-key"

# Application configuration
export LOG_LEVEL="info"
export PORT="8003"
export ENVIRONMENT="production"
```

**Health Checks:**
```go
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now(),
        "checks": map[string]interface{}{
            "database": checkDatabaseHealth(),
            "rabbitmq": checkRabbitMQHealth(),
            "external_tools": checkExternalToolsHealth(),
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}
```

This development guide provides comprehensive guidelines for contributing to and maintaining the NuM8 project. It should be updated as the project evolves and new practices are adopted.