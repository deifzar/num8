# NuM8

**NuM8** is a security scanning application that integrates with Nuclei and Burp Suite for comprehensive vulnerability assessment and web application testing. Built with a message-driven microservices architecture, NuM8 provides scalable, automated security scanning capabilities.

## Features

- **Nuclei Integration** - Automated vulnerability scanning using Nuclei v3.3.5 templates
- **Burp Suite Integration** - Web application security testing via Burp Suite REST API
- **Message-Driven Architecture** - RabbitMQ-based asynchronous processing with connection pooling
- **Manual Acknowledgment** - Reliable message processing with delivery tag tracking and ACK/NACK support
- **Multi-Channel Notifications** - Discord, email, and application-level alerting
- **Structured Logging** - Zerolog with automatic log rotation and configurable log levels
- **RESTful API** - Clean API design for scan orchestration and result retrieval
- **Docker Support** - Containerized deployment with environment variable configuration

## Technology Stack

- **Language**: Go 1.21.5+
- **Database**: PostgreSQL 14+ (database: `cptm8`, schema: `public`)
- **Message Queue**: RabbitMQ 3.8+ with connection pooling
- **HTTP Framework**: Gin
- **CLI Framework**: Cobra
- **Configuration**: Viper (with hot-reload support)
- **Logging**: Zerolog with Lumberjack rotation
- **Security Tools**: Nuclei v3.3.5, Burp Suite

## Quick Start

### Prerequisites

- Go 1.21.5 or higher
- PostgreSQL 14+
- RabbitMQ 3.8+
- Docker (optional)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/NuM8.git
cd NuM8

# Install dependencies
go mod download

# Build the application
go build -o num8 main.go
```

### Configuration

1. Copy the configuration template:
```bash
cp configs/configuration_template.yaml configs/configuration.yaml
```

2. Set environment variables for sensitive data:
```bash
export POSTGRESQL_PASSWORD="your_db_password"
export RABBITMQ_PASSWORD="your_rabbitmq_password"
export DISCORD_BOT_TOKEN="your_discord_token"
export DISCORD_WEBHOOK_TOKEN="your_webhook_token"
```

3. Update `configs/configuration.yaml` with your configuration:
```yaml
Database:
  location: "${POSTGRESQL_HOSTNAME}"
  password: "${POSTGRESQL_PASSWORD}"

RabbitMQ:
  location: "${RABBITMQ_HOSTNAME}"
  password: "${RABBITMQ_PASSWORD}"
```

### Running the Application

```bash
# Launch API service (default: 0.0.0.0:8003)
./num8 launch

# Launch with custom IP and port (port must be 8000-8999)
./num8 launch --ip 127.0.0.1 --port 8080

# Check version
./num8 version

# Get help
./num8 help
```

### Docker Deployment

```bash
# Build Docker image
docker build -t num8 -f dockerfile .

# Run container
docker run -p 8003:8003 \
  -e POSTGRESQL_HOSTNAME=db \
  -e POSTGRESQL_DB=cptm8 \
  -e POSTGRESQL_USERNAME=postgres \
  -e POSTGRESQL_PASSWORD=secret \
  -e RABBITMQ_HOSTNAME=rabbitmq \
  -e RABBITMQ_PASSWORD=secret \
  num8 launch
```

## API Endpoints

### Scan Endpoints

- `POST /scan` - Initiate general security scan
- `POST /domain/:id/scan` - Scan specific domain
- `POST /hostname/:id/scan` - Scan specific hostname
- `POST /endpoint/:id/scan` - Scan specific endpoint

### Example Request

```bash
curl -X POST http://localhost:8003/scan \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.com",
    "scan_type": "full"
  }'
```

## Architecture

NuM8 uses a microservices architecture with the following components:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client APIs   │    │  External Tools │    │  Notifications  │
│                 │    │                 │    │                 │
│ • REST API      │    │ • Nuclei        │    │ • Discord       │
│ • Web Interface │    │ • Burp Suite    │    │ • Email         │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                        NuM8 Core                               │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   API8       │  │ Controller8  │  │Orchestrator8 │         │
│  │              │  │              │  │              │         │
│  │ • Routes     │  │ • Numate     │  │ • Workflows  │         │
│  │ • Middleware │  │ • Burpmate   │  │ • Queues     │         │
│  │ • Validation │  │ • Processing │  │ • ACK/NACK   │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────────────┘
          │                                              │
          ▼                                              ▼
┌─────────────────┐                            ┌─────────────────┐
│   PostgreSQL    │                            │    RabbitMQ     │
│                 │                            │                 │
│ • Schema: cptm8 │                            │ • Pool: 2-10    │
│ • Tables:       │                            │ • Manual ACK    │
│   - domains     │                            │ • Retry: 10x    │
│   - endpoints   │                            │ • Health checks │
│   - issues      │                            │                 │
└─────────────────┘                            └─────────────────┘
```

### Message Flow with Manual Acknowledgment

NuM8 implements a robust message acknowledgment strategy:

1. RabbitMQ delivers message with `deliveryTag` to consumer
2. Orchestrator extracts `deliveryTag` and creates HTTP request with `X-RabbitMQ-Delivery-Tag` header
3. Controller extracts `deliveryTag` from request header
4. Controller passes `deliveryTag` to async scan goroutine
5. Scan executes with defer function tracking completion status
6. On completion, defer function calls `AckScanCompletion(deliveryTag, scanCompleted)`
7. Orchestrator ACKs (success) or NACKs (failure, requeue=true) the message

**ACK/NACK Strategy:**
- **Scan success**: ACK (message removed from queue)
- **Scan failure/panic**: NACK with requeue=true (retry scan)
- **Configuration/validation errors**: NACK with requeue=false (dead letter)
- **Handler errors**: NACK with requeue=true (temporary failure)

## Directory Structure

```
NuM8/
├── main.go                    # Entry point
├── cmd/                       # CLI commands (Cobra)
│   ├── launch.go             # Launch API server command
│   ├── version.go            # Version command
│   └── root.go               # Root command with graceful shutdown
├── pkg/                       # Core packages
│   ├── api8/                 # REST API layer (Gin)
│   ├── controller8/          # Business logic controllers
│   ├── model8/               # Data models (23 files)
│   ├── db8/                  # Database layer
│   ├── amqpM8/               # RabbitMQ connection pool
│   ├── orchestrator8/        # Message orchestration
│   ├── notification8/        # Notification system
│   ├── configparser/         # Configuration management
│   ├── log8/                 # Structured logging
│   ├── cleanup8/             # Temporary file cleanup
│   └── utils/                # Utility functions
├── configs/                   # Configuration files
├── docs/                      # Documentation
├── log/                       # Log files
└── tmp/                       # Temporary scan results
```

## Configuration

### Key Configuration Sections

```yaml
# Application environment
APP_ENV: DEV|TEST|PROD
LOG_LEVEL: 0-5 # (debug to panic)

# Service orchestration
ORCHESTRATORM8:
  Services: [num8, asmm8]
  Exchanges:
    - cptm8 (topic)
    - notification (topic)
  Queues:
    # Format: [consumer_name, queue_name, autoACK]
    - ["num8_consumer", "num8.scan.queue", "false"]  # Manual ACK

# Nuclei and Burp Suite integration
NUM8:
  BurpAPILocation: http://127.0.0.1:8090
  BurpProxyLocation: http://127.0.0.1:8080
  TemplateURLs: [...]
  TrustedDomains: [...]

# Database configuration
Database:
  location: "${POSTGRESQL_HOSTNAME}"
  port: 5432
  database: cptm8
  schema: public
  username: "${POSTGRESQL_USERNAME}"
  password: "${POSTGRESQL_PASSWORD}"

# RabbitMQ connection pool
RabbitMQ:
  location: "${RABBITMQ_HOSTNAME}"
  port: 5672
  username: "${RABBITMQ_USERNAME}"
  password: "${RABBITMQ_PASSWORD}"
  pool:
    max_connections: 10
    min_connections: 2
    max_idle_time: 1h
    max_lifetime: 2h
    health_check_period: 30m
    connection_timeout: 30s
    retry_attempts: 10
    retry_delay: 2s
```

### Environment Variable Support

All sensitive configuration values support environment variable substitution using `${VARIABLE_NAME}` syntax.

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

## Development

### Adding a New Endpoint

1. Define model in `pkg/model8/`
2. Add controller method in `pkg/controller8/`
3. Register route in `pkg/api8/api8.go`
4. Add validation logic
5. Update tests

### Adding a New RabbitMQ Consumer

1. Define message structure in `pkg/model8/`
2. Configure queue in `configs/configuration.yaml` with manual ACK: `[consumer_name, queue_name, "false"]`
3. Implement consumer in `pkg/orchestrator8/` with delivery tag propagation
4. Add handler logic in controller with ACK/NACK support
5. Test message flow and verify ACK/NACK behavior in logs

### Debugging Tips

```bash
# Check application logs
tail -f log/num8.log

# Track message acknowledgment flow
grep "deliveryTag" log/num8.log

# Check RabbitMQ queue status
rabbitmqctl list_queues name messages_ready messages_unacknowledged

# Monitor RabbitMQ management interface
open http://localhost:15672
```

## Documentation

Comprehensive documentation is available in the `docs/` directory:

- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - System architecture and component relationships
- [DEVELOPMENT.md](docs/DEVELOPMENT.md) - Development guidelines and best practices
- [SECURITY.md](docs/SECURITY.md) - Security considerations and recommendations
- [CODE_REVIEW.md](docs/CODE_REVIEW.md) - Code review checklist and findings
- [PERFORMANCE.md](docs/PERFORMANCE.md) - Performance analysis and optimization guide
- [TODO.md](docs/TODO.md) - Known issues and improvement roadmap
- [CLAUDE.md](CLAUDE.md) - AI assistant guidance for working with this codebase

## Security Considerations

**IMPORTANT**: This application has known security vulnerabilities. Review [SECURITY.md](docs/SECURITY.md) before deployment.

**Current Security Status:**
- **Authentication**: Not implemented - all endpoints publicly accessible
- **Authorization**: Not implemented - no RBAC
- **Input Validation**: Minimal - needs enhancement
- **Credential Management**: Supports environment variables (use them!)
- **Logging**: Structured logging with rotation
- **File Permissions**: Secure umask (0027) - Files: 640, Dirs: 750

**Recommendations for Production:**
- Always use environment variables for credentials
- Implement authentication middleware
- Add comprehensive input validation
- Enable TLS/SSL for API endpoints
- Implement rate limiting
- Add security headers
- Review and implement recommendations in [SECURITY.md](docs/SECURITY.md)

## Performance Considerations

**Current Limitations:**
- Sequential scan processing (opportunity for concurrency)
- RabbitMQ queue limits: max 2 messages (increase for production)
- No caching layer implemented
- Database connection pooling handled by standard library

**Optimization Opportunities:**
- Implement worker pools for concurrent scanning
- Increase RabbitMQ queue limits and connection pool
- Add Redis caching layer
- Optimize database queries with prepared statements
- Implement batch processing for large datasets

See [PERFORMANCE.md](docs/PERFORMANCE.md) for detailed analysis and recommendations.

## Logging

NuM8 uses structured logging with the following features:

- **Framework**: Zerolog with Lumberjack rotation
- **Log Levels**: trace(-1), debug(0), info(1), warn(2), error(3), fatal(4), panic(5)
- **Rotation**: Max 100MB per file, 3 backups
- **File Permissions**: 0640 (compatible with Vector log shipping)
- **Environment-Specific**:
  - DEV/TEST: Console output
  - PROD: File output with rotation

## Known Issues

See [TODO.md](docs/TODO.md) for a comprehensive list. Key issues:

- No authentication/authorization (CRITICAL)
- Limited input validation (HIGH)
- Sequential processing bottleneck (MEDIUM)
- Missing comprehensive tests (MEDIUM)
- No monitoring/metrics (MEDIUM)

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or contributions:

- **Issues**: Open an issue in the GitHub repository
- **Documentation**: See the `docs/` directory
- **Security**: Report security issues privately (see [SECURITY.md](docs/SECURITY.md))

## Acknowledgments

- **Nuclei** - ProjectDiscovery for the vulnerability scanning engine
- **Burp Suite** - PortSwigger for web application security testing
- **RabbitMQ** - For reliable message queuing
- **PostgreSQL** - For robust data persistence
- **Go Community** - For excellent libraries and frameworks

## Roadmap

See [TODO.md](docs/TODO.md) for the complete roadmap. Upcoming features:

- Authentication and authorization implementation
- Enhanced input validation
- Worker pool for concurrent scanning
- Redis caching layer
- Prometheus metrics and monitoring
- Comprehensive test coverage
- CI/CD pipeline integration
- Kubernetes deployment manifests

---

**Note**: This is a security scanning tool intended for authorized testing only. Ensure you have proper authorization before scanning any targets. Unauthorized scanning may be illegal.
