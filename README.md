# NuM8 - Nuclei and BurpSuite Mate

<div align="center">

**Production-grade Go microservice for automated security vulnerability scanning and web application testing**

[![Go Version](https://img.shields.io/badge/Go-1.21.5+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?logo=docker)](dockerfile)
[![Status](https://img.shields.io/badge/status-active%20development-yellow)](https://github.com/yourusername/num8)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)

[Features](#features) • [Quick Start](#quick-start) • [Documentation](#documentation) • [Architecture](#architecture) • [API](#api-reference)

</div>

---

## Overview

NuM8 (Nuclei and BurpSuite Mate) is a production-grade Go microservice designed for automated security vulnerability scanning and web application testing. It provides a robust REST API for orchestrating security scans using industry-standard tools like Nuclei and Burp Suite.

**Built for:**
- Security researchers and penetration testers
- Red teams conducting security assessments
- Bug bounty hunters identifying vulnerabilities
- Security operations teams managing vulnerability assessments
- DevSecOps teams integrating security into CI/CD pipelines

### Key Features

- **Dual Scanning Modes**: Nuclei-based vulnerability scanning and Burp Suite integration for web app testing
- **Asynchronous Processing**: RabbitMQ-based message queuing with advanced connection pooling
- **Manual Acknowledgment**: Reliable message processing with delivery tag tracking and smart ACK/NACK logic
- **Persistent Storage**: PostgreSQL integration with security findings history and tracking
- **Production Ready**: Docker containerization, structured logging, and graceful shutdown
- **Scalable Architecture**: Interface-based design with dependency injection
- **Multi-Channel Notifications**: Email and custom notification channels

### Stats at a Glance

```
📁 46 Go source files        🔧 11 major components
📊 4,980 lines of code       🐳 Multi-stage Docker build
📦 11 packages               🔌 REST API with 4+ endpoints
🛠️  2 security tools         💾 PostgreSQL + RabbitMQ
```

---

## Quick Start

### Prerequisites

- **Go** 1.21.5 or higher
- **PostgreSQL** 14+ (for data persistence)
- **RabbitMQ** 3.8+ (for message queuing)
- **Burp Suite** (optional, for web application testing)
- **Docker** (optional, for containerized deployment)

### Installation

#### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/NuM8.git
cd NuM8

# Install Go dependencies
go mod download

# Build the binary
go build -o num8 main.go

# Run the service
./num8 launch --ip 0.0.0.0 --port 8003
```

#### Option 2: Docker

```bash
# Build the Docker image
docker build -t num8:latest -f dockerfile .

# Run the container
docker run -d \
  -p 8003:8003 \
  -e POSTGRESQL_HOSTNAME=your-db-host \
  -e POSTGRESQL_PASSWORD=your-db-password \
  -e RABBITMQ_HOSTNAME=your-rabbitmq-host \
  -e RABBITMQ_PASSWORD=your-rabbitmq-password \
  --name num8 \
  num8:latest
```

### Configuration

1. Copy the example configuration:
```bash
cp configs/configuration_template.yaml configs/configuration.yaml
```

2. Edit `configs/configuration.yaml` with your settings:
```yaml
APP_ENV: PROD
LOG_LEVEL: "1"  # 0=debug, 1=info, 2=warn, 3=error

NUM8:
  BurpAPILocation: http://127.0.0.1:8090
  BurpProxyLocation: http://127.0.0.1:8080
  NucleiTemplateURLs:
    - https://github.com/projectdiscovery/nuclei-templates
  TrustedDomains:
    - example.com

Database:
  location: localhost
  port: 5432
  database: cptm8
  schema: public
  username: cpt_dbuser
  password: ${POSTGRESQL_PASSWORD}  # Use environment variables for secrets

RabbitMQ:
  location: localhost
  port: 5672
  username: ${RABBITMQ_USERNAME}
  password: ${RABBITMQ_PASSWORD}
  pool:
    max_connections: 10
    min_connections: 2
    retry_attempts: 10

ORCHESTRATORM8:
  Queues:
    # Format: [consumer_name, queue_name, autoACK]
    - ["num8_consumer", "num8.scan.queue", "false"]  # Manual ACK mode
  Exchanges:
    - name: cptm8
      type: topic
    - name: notification
      type: topic
```

3. Set environment variables for sensitive data:
```bash
export POSTGRESQL_PASSWORD="your-secure-password"
export RABBITMQ_USERNAME="your-rabbitmq-user"
export RABBITMQ_PASSWORD="your-rabbitmq-password"
```

### Database Setup

```sql
-- Create database
CREATE DATABASE cptm8;

-- Create tables (example schema)
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    companyname VARCHAR(255),
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE endpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain_id UUID REFERENCES domains(id),
    url VARCHAR(2048) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE history_issues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    endpoint_id UUID REFERENCES endpoints(id),
    severity VARCHAR(50),
    title VARCHAR(500),
    description TEXT,
    status VARCHAR(10),  -- U (Unverified), V (Verified), I (In Progress), FP (False Positive), R (Resolved)
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP
);

CREATE INDEX idx_issues_endpoint ON history_issues(endpoint_id);
CREATE INDEX idx_issues_status ON history_issues(status);
CREATE INDEX idx_issues_severity ON history_issues(severity);
```

---

## Usage

### Basic Workflow

1. **Trigger a general security scan:**
```bash
curl -X POST http://localhost:8003/scan \
  -H "Content-Type: application/json" \
  -d '{
    "target": "https://example.com",
    "templateFilters": ["cve", "vulnerabilities"]
  }'
```

2. **Scan a specific domain:**
```bash
curl -X POST http://localhost:8003/domain/{domain-uuid}/scan \
  -H "Content-Type: application/json"
```

3. **Scan a specific endpoint:**
```bash
curl -X POST http://localhost:8003/endpoint/{endpoint-uuid}/scan \
  -H "Content-Type: application/json" \
  -d '{
    "scanType": "nuclei",
    "severity": ["critical", "high", "medium"]
  }'
```

4. **Retrieve security findings:**
```bash
curl http://localhost:8003/endpoint/{endpoint-uuid}/issues
```

### Scanning Modes

#### Nuclei Vulnerability Scanning
- Template-based vulnerability detection
- Supports custom template filters (CVEs, technologies, severity levels)
- Automatic template updates from GitHub
- Results stored with severity classification
- False positive tracking and management

#### Burp Suite Integration
- Sitemap crawling and analysis
- Passive security checks
- Active scanning capabilities (configurable)
- Integration with Burp Suite Professional/Enterprise
- HTTP status code and content-type filtering

---

## API Reference

### Scan Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/scan` | Initiate general security scan |
| `POST` | `/domain/:id/scan` | Scan specific domain |
| `POST` | `/hostname/:id/scan` | Scan specific hostname |
| `POST` | `/endpoint/:id/scan` | Scan specific endpoint |

### Health Check Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Liveness probe (always 200 OK) |
| `GET` | `/ready` | Readiness probe (checks DB + RabbitMQ) |

**Example Response:**
```json
{
  "status": "ready",
  "database": "connected",
  "rabbitmq": "connected",
  "timestamp": "2025-11-16T12:34:56Z"
}
```

---

## Architecture

### High-Level Overview

```
┌─────────────┐
│   Client    │───┐
│   (HTTP)    │   │
└─────────────┘   │
                  │    ┌─────────────┐    ┌─────────────┐
┌─────────────┐   └───▶│  NuM8 API   │───▶│  Database   │
│  RabbitMQ   │───────▶│  (Gin/Go)   │    │ (PostgreSQL)│
│  (Message   │        └──────┬──────┘    └─────────────┘
│   Queue)    │               │
└─────────────┘               │
                              │
                ┌─────────────┼─────────────┐
                ▼             ▼             ▼
        ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
        │   Nuclei    │ │ Burp Suite  │ │Notification │
        │Vulnerability│ │   Web App   │ │  Channels   │
        │   Scanner   │ │   Testing   │ │             │
        └─────────────┘ └─────────────┘ └─────────────┘
```

### Scanning Workflow

**Complete Scan Pipeline:**

```
1. SCAN INITIATION
   ├─→ REST API Request (HTTP POST)
   └─→ RabbitMQ Message (with deliveryTag) → Controller

2. VULNERABILITY SCANNING (Nuclei)
   ├─→ Template selection and filtering
   ├─→ Target validation and preparation
   ├─→ Nuclei execution with custom templates
   ├─→ Result parsing and normalization
   └─→ Database storage (history_issues table)

3. WEB APPLICATION TESTING (Burp Suite)
   ├─→ Sitemap retrieval via Burp API
   ├─→ Passive security analysis
   ├─→ Active scanning (if configured)
   └─→ Findings aggregation and storage

4. NOTIFICATION & ORCHESTRATION
   ├─→ Security findings published to notification queue
   ├─→ Email alerts for critical findings
   ├─→ Downstream service notification via RabbitMQ
   └─→ Message ACK/NACK based on scan completion

5. MANUAL ACKNOWLEDGMENT FLOW
   └─→ Scan completion → ACK (success) or NACK+requeue (failure)
```

### Package Structure

```
num8/
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Base command setup with graceful shutdown
│   ├── launch.go          # API service launcher (port 8000-8999)
│   └── version.go         # Version information
├── pkg/                    # 11 packages, 46 Go files
│   ├── amqpM8/            # RabbitMQ connection pooling (5 files)
│   ├── api8/              # HTTP API routes and initialization
│   ├── cleanup8/          # Temporary file cleanup utilities
│   ├── configparser/      # Configuration management (Viper)
│   ├── controller8/       # Business logic controllers (Numate, Burpmate)
│   ├── db8/               # Database access layer (PostgreSQL)
│   ├── log8/              # Structured logging (zerolog + lumberjack)
│   ├── model8/            # Data models and domain entities (23 files)
│   ├── notification8/     # Notification system (multi-channel)
│   ├── orchestrator8/     # Service orchestration and message routing
│   └── utils/             # Utility functions (IP validation, jq filtering)
├── configs/               # Configuration files (YAML)
├── docs/                  # Comprehensive documentation
├── log/                   # Application logs
├── tmp/                   # Temporary scan results
└── main.go                # Application entry point (sets umask, initializes logger)
```

### Key Components

- **[API Layer](pkg/api8/)** - Gin-based REST API
- **[Controllers](pkg/controller8/)** - Business logic (Numate, Burpmate)
- **[Database Layer](pkg/db8/)** - PostgreSQL repository pattern
- **[Message Queue](pkg/amqpM8/)** - RabbitMQ with connection pooling
- **[Orchestrator](pkg/orchestrator8/)** - Service coordination and message routing
- **[Notification System](pkg/notification8/)** - Multi-channel alerting
- **[Configuration](pkg/configparser/)** - Viper-based config management
- **[Logging](pkg/log8/)** - Zerolog structured logging with rotation

For detailed architecture documentation, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

---

## Advanced Features

### RabbitMQ Integration

**Advanced Connection Pooling:**
- Configurable pool size (default: 2-10 connections)
- Automatic connection recovery on failures
- Periodic health checks (30-minute intervals)
- Manual message acknowledgment with smart ACK/NACK logic
- Delivery tag tracking for message lifecycle management

**Message Flow:**
```
RabbitMQ → Consumer → Extract deliveryTag → HTTP Request (X-RabbitMQ-Delivery-Tag header)
    ↓
RunNumate()/RunBurpmate() → Extract from header → Async scan execution
    ↓
Defer ACK/NACK → Success: ACK | Failure/Panic: NACK+requeue
```

**ACK/NACK Strategy:**
- **Scan completed successfully**: ACK (message removed from queue)
- **Scan failed/panic**: NACK with requeue=true (retry scan)
- **Configuration errors**: NACK with requeue=false (dead letter)
- **Handler errors**: NACK with requeue=true (temporary failure)

### Error Handling

**Resilient Design:**
- Panic recovery with defer blocks in scan functions
- Automatic requeue on scan failures
- Error notifications via RabbitMQ
- Graceful degradation on tool failures
- Context-based cancellation for long-running operations

### Nuclei Template Management

- Automatic template updates from GitHub repositories
- Custom template filtering (severity, tags, CVEs, technologies)
- Template caching for improved performance
- Trusted domain whitelisting for false positive reduction
- Support for custom template repositories

### Burp Suite Integration

- Sitemap crawling with configurable filters
- HTTP status code filtering (configurable ranges)
- Content-type filtering (HTML, JSON, XML, etc.)
- Passive and active scanning modes
- Result normalization and deduplication
- Integration with Burp Suite REST API (ports 8080/8090)

---

## Documentation

Comprehensive documentation is available in the [docs/](docs/) directory:

| Document | Description |
|----------|-------------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Detailed system architecture and component relationships |
| [DEVELOPMENT.md](docs/DEVELOPMENT.md) | Development setup and guidelines |
| [SECURITY.md](docs/SECURITY.md) | Security best practices and vulnerability assessment |
| [PERFORMANCE.md](docs/PERFORMANCE.md) | Performance optimization and analysis |
| [TODO.md](docs/TODO.md) | Known issues and roadmap |
| [CODE_REVIEW.md](docs/CODE_REVIEW.md) | Code review checklist |

---

## External Tools

NuM8 integrates with industry-standard security testing tools:

| Tool | Version | Purpose | Integration |
|------|---------|---------|-------------|
| [Burp Suite](https://portswigger.net/burp) | Professional/Enterprise | Web application security testing | REST API |

---

## Performance

**Typical Performance Metrics:**

- **Nuclei Scan**: 1-10 minutes (depends on template count and target complexity)
- **Burp Suite Scan**: 5-60 minutes (depends on sitemap size and scan depth)
- **Database Operations**: Optimized with connection pooling and batch insertions
- **Message Processing**: Manual ACK ensures reliable delivery and retry logic
- **Concurrent Scans**: Goroutine-based parallel execution

**Resource Requirements:**

- **CPU**: 2+ cores recommended
- **Memory**: 2 GB minimum, 4 GB recommended
- **Storage**: 10 GB for application + logs + temporary scan results
- **Network**: Stable internet connection for Nuclei template updates

For optimization tips, see [docs/PERFORMANCE.md](docs/PERFORMANCE.md).

---

## Security Considerations

### Current Limitations

- **No authentication** on API endpoints (planned for v2.0)
- Database credentials in configuration file (**use environment variables**)
- Limited input validation (basic Gin binding only)
- **CRITICAL**: Exposed credentials in sample configs must be rotated

### Recommendations

1. **Use environment variables** for all secrets and credentials
2. **Deploy behind API gateway** with authentication (JWT, OAuth2)
3. **Enable TLS/SSL** for production deployments
4. **Implement rate limiting** to prevent abuse and DoS
5. **Run as non-root user** in containers (already configured)
6. **Rotate all exposed credentials** before production deployment
7. **Review and implement** recommendations in [docs/SECURITY.md](docs/SECURITY.md)

See [docs/SECURITY.md](docs/SECURITY.md) for comprehensive security guidelines.

---

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Follow Go best practices and existing code style
4. Add tests for new functionality
5. Update documentation as needed
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed development guidelines.

---

## Roadmap

### Version 1.x (Current)
- [x] Core vulnerability scanning functionality (Nuclei)
- [x] Burp Suite integration
- [x] PostgreSQL persistence
- [x] RabbitMQ message queuing with connection pooling
- [x] Docker containerization
- [x] Manual ACK/NACK message reliability
- [x] Multi-channel notifications

### Version 2.0 (Planned)
- [ ] JWT-based authentication
- [ ] Rate limiting and request throttling
- [ ] Unit and integration tests (target: 80% coverage)
- [ ] Kubernetes deployment manifests
- [ ] Prometheus metrics integration
- [ ] GraphQL API option
- [ ] Web dashboard for visualization
- [ ] Enhanced Burp Suite active scanning
- [ ] Custom Nuclei template repository management

See [docs/TODO.md](docs/TODO.md) for the complete roadmap and known issues.

---

## Troubleshooting

### Common Issues

**1. Database connection failures**
```bash
# Check PostgreSQL is running
systemctl status postgresql

# Verify connection settings in configs/configuration.yaml
# Ensure database and tables are created
psql -U postgres -d cptm8 -c "\dt"
```

**2. RabbitMQ connection errors**
```bash
# Check RabbitMQ status
systemctl status rabbitmq-server

# Verify credentials and port in configuration
# Check exchange and queue creation
rabbitmqctl list_exchanges
rabbitmqctl list_queues
```

**3. Nuclei not found**
```bash
# Ensure Nuclei is in PATH
which nuclei

# Reinstall if needed
go install github.com/projectdiscovery/nuclei/v3/cmd/nuclei@v3.3.5

# Update templates
nuclei -update-templates
```

**4. Burp Suite connection errors**
```bash
# Ensure Burp Suite is running
# Check REST API is enabled in Burp Suite settings
# Verify ports 8080 (proxy) and 8090 (REST API) are accessible
curl http://127.0.0.1:8090/v1/
```

**5. Permission errors**
```bash
# Ensure log directory is writable
chmod 755 log/

# Check file permissions for config files
chmod 640 configs/configuration.yaml

# Verify umask is set correctly (0027)
umask
```

**6. Message stuck in RabbitMQ queue**
```bash
# Check unacknowledged messages
rabbitmqctl list_queues name messages_ready messages_unacknowledged

# Track delivery tag flow in logs
grep "deliveryTag" log/num8.log

# Look for ACK/NACK patterns
grep -E "ACK|NACK" log/num8.log
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Acknowledgments

- [ProjectDiscovery](https://github.com/projectdiscovery) for excellent security tools (Nuclei)
- [PortSwigger](https://portswigger.net/) for Burp Suite web security testing platform
- [Gin Web Framework](https://github.com/gin-gonic/gin) for the HTTP router
- [Cobra CLI](https://github.com/spf13/cobra) for CLI framework
- [Viper](https://github.com/spf13/viper) for configuration management
- [Zerolog](https://github.com/rs/zerolog) for structured logging
- [RabbitMQ](https://www.rabbitmq.com/) for reliable message queuing
- [PostgreSQL](https://www.postgresql.org/) for robust data persistence

---

## Contact

For questions, issues, or feature requests, please open an issue on GitHub.

**Project Link:** [https://github.com/yourusername/NuM8](https://github.com/yourusername/NuM8)

---

<div align="center">

**Built with ❤️ for the security community**

**⚠️ Note**: This is a security scanning tool intended for **authorized testing only**. Ensure you have proper authorization before scanning any targets. Unauthorized scanning may be illegal.

[⬆ Back to Top](#num8---nuclei-and-burpsuite-mate)

</div>
