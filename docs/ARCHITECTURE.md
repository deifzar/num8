# NuM8 Architecture Documentation

**Version:** 1.0  
**Last Updated:** January 2025  
**Document Type:** System Architecture

## Overview

NuM8 is a security scanning application that integrates with Nuclei (for vulnerability scanning) and Burp Suite (for web application testing). It uses a microservices architecture with RabbitMQ for message queuing and PostgreSQL for data persistence.

## High-Level Architecture

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
│  │ • Validation │  │ • Processing │  │ • Exchanges  │         │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘         │
│         │                 │                 │                 │
│         └─────────────────┼─────────────────┘                 │
│                           │                                   │
│  ┌──────────────┐  ┌──────┴───────┐  ┌──────────────┐         │
│  │   Model8     │  │   AMQPM8     │  │Notification8 │         │
│  │              │  │              │  │              │         │
│  │ • Entities   │  │ • Queues     │  │ • Channels   │         │
│  │ • DTOs       │  │ • Exchanges  │  │ • Templates  │         │
│  │ • Validation │  │ • Consumers  │  │ • Delivery   │         │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘         │
│         │                 │                 │                 │
│         └─────────────────┼─────────────────┘                 │
│                           │                                   │
│  ┌──────────────┐  ┌──────┴───────┐  ┌──────────────┐         │
│  │    DB8       │  │   Config     │  │    Log8      │         │
│  │              │  │              │  │              │         │
│  │ • Connection │  │ • Parser     │  │ • Structured │         │
│  │ • Queries    │  │ • Validation │  │ • Levels     │         │
│  │ • Migrations │  │ • Env Vars   │  │ • Rotation   │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────────────┘
          │                                              │
          ▼                                              ▼
┌─────────────────┐                            ┌─────────────────┐
│   PostgreSQL    │                            │    RabbitMQ     │
│                 │                            │                 │
│ • Schema: cptm8 │                            │ • Exchange:     │
│ • Tables:       │                            │   cptm8         │
│   - domains     │                            │ • Queues:       │
│   - endpoints   │                            │   - scan_queue  │
│   - issues      │                            │   - notif_queue │
│   - history     │                            │ • Routing Keys  │
└─────────────────┘                            └─────────────────┘
```

## Component Architecture

### 1. API Layer (`pkg/api8/`)

The API layer provides the HTTP interface for the application:

**Responsibilities:**
- HTTP request routing
- Middleware application (auth, logging, CORS)
- Request/response serialization
- Error handling and status codes

**Key Components:**
- `api8.go` - Main API server setup and configuration
- Route handlers for scan endpoints
- Middleware stack for cross-cutting concerns

**Endpoints:**
- `POST /scan` - Initiate general security scan
- `POST /domain/:id/scan` - Scan specific domain
- `POST /hostname/:id/scan` - Scan specific hostname
- `POST /endpoint/:id/scan` - Scan specific endpoint

### 2. Controller Layer (`pkg/controller8/`)

The controller layer orchestrates business logic and external tool integration:

**Responsibilities:**
- Business logic coordination
- External tool integration (Nuclei, Burp Suite)
- Scan workflow management
- Result processing and aggregation

**Key Components:**
- `controller8_numate.go` - Nuclei vulnerability scanning controller
- `controller8_burpmate.go` - Burp Suite integration controller

**Integration Points:**
- Nuclei templates and scanning engines
- Burp Suite REST API (ports 8080/8090)
- Result parsing and normalization

### 3. Orchestrator Layer (`pkg/orchestrator8/`)

The orchestrator manages message-driven workflows and inter-service communication:

**Responsibilities:**
- Workflow orchestration with manual message acknowledgment
- Message queue management with delivery tag tracking
- Service coordination via HTTP callbacks
- Asynchronous processing with scan completion tracking
- Message ACK/NACK handling based on scan outcomes

**Key Components:**
- `orchestrator8.go` - Main orchestration logic with ACK/NACK methods
- Queue management and routing
- Message publishing and consumption with delivery tag propagation
- Scan completion acknowledgment (`AckScanCompletion`, `NackScanMessage`)

**Message Flow with Manual ACK:**
```
RabbitMQ Message → Orchestrator (deliveryTag captured)
                       ↓
                HTTP Request (X-RabbitMQ-Delivery-Tag header)
                       ↓
                Controller (extract deliveryTag)
                       ↓
                Scan Execution (async goroutine)
                       ↓
           Scan Completes/Fails (defer function)
                       ↓
      AckScanCompletion(deliveryTag, scanCompleted)
                       ↓
            ACK (success) or NACK (failure, requeue)
```

**ACK/NACK Strategy:**
- **Successful scan completion:** ACK message (removed from queue)
- **Scan failure/panic/error:** NACK with requeue=true (retry scan)
- **Configuration/validation errors:** NACK with requeue=false (dead letter)
- **Handler errors:** NACK with requeue=true (temporary failure)
- **No handler found:** NACK with requeue=false (permanent failure)

### 4. Data Model Layer (`pkg/model8/`)

The data model layer defines application entities and data structures:

**Key Entities:**
- `endpoint8.go` - Target endpoint representations
- `historyissues.go` - Security findings and vulnerability data
- `notification8.go` - Notification message structures

**Data Flow:**
```
External Data → DTOs → Business Entities → Database Models
```

### 5. Database Layer (`pkg/db8/`)

The database layer manages data persistence and queries:

**Responsibilities:**
- Database connection management
- Query execution and optimization
- Transaction management
- Data migration support

**Database Schema:**
- **Database:** `cptm8`
- **Schema:** `public`
- **Key Tables:**
  - `domains` - Target domain information
  - `endpoints` - Scan targets and configurations
  - `history_issues` - Security findings and vulnerabilities
  - `notifications` - Alert and notification logs

### 6. Message Queue Layer (`pkg/amqpM8/`)

The message queue layer provides asynchronous communication infrastructure with connection pooling:

**Responsibilities:**
- AMQP connection pool management (min/max connections, health checks)
- Queue and exchange setup
- Message publishing and consuming with manual acknowledgment support
- Failure handling and retries with requeue logic
- Consumer health monitoring and recovery

**Queue Architecture:**
- **Exchange:** `cptm8` (topic exchange)
- **Queues:** Configured in `configuration.yaml`
- **Routing:** Topic-based routing with binding keys
- **Acknowledgment Mode:** Manual ACK (default: `autoACK: false`)
- **Message Reliability:** Delivery tags tracked for scan completion acknowledgment

**Connection Pool Features:**
- Configurable pool size (default: 2-10 connections)
- Automatic connection recycling and health checks
- Connection lifetime management (max 2h, idle 1h)
- Retry mechanism with exponential backoff (default: 10 retries, 2s delay)

### 7. Notification Layer (`pkg/notification8/`)

The notification layer handles alert delivery and external integrations:

**Responsibilities:**
- Multi-channel notification delivery
- Message formatting and templating
- Delivery status tracking
- Integration with external services

**Supported Channels:**
- Discord webhooks
- Email delivery
- Application-level notifications

## Data Flow Architecture

### 1. Scan Initiation Flow

```
Client Request → API8 → Controller8 → Orchestrator8 → Queue
                                   ↓
External Tool ← Worker ← Queue ← AMQPM8
                                   ↓
Results → DB8 → Notification8 → External Services
```

### 2. Message Processing Flow (with Manual ACK)

```
Publisher → Exchange (cptm8) → Routing Key → Queue → Consumer (captures deliveryTag)
                                                         ↓
                                              Handler (processes message)
                                                         ↓
                                    ┌────────────────────┴────────────────────┐
                                    ↓                                         ↓
                          Handler Success                           Handler Error
                                    ↓                                         ↓
                    Wait for Scan Completion                      NACK (requeue=true)
                                    ↓
                    ┌───────────────┴───────────────┐
                    ↓                               ↓
           Scan Completed                    Scan Failed/Panic
                    ↓                               ↓
            ACK (remove from queue)         NACK (requeue=true)
```

**Delivery Tag Propagation:**
1. Consumer receives message with `deliveryTag`
2. Orchestrator creates HTTP request with `X-RabbitMQ-Delivery-Tag` header
3. Controller extracts `deliveryTag` from request header
4. Controller passes `deliveryTag` to scan goroutine
5. Scan completion defer function ACKs/NACKs based on `scanCompleted` flag

### 3. Database Interaction Flow

```
Controller → DB8 → Connection Pool → PostgreSQL → Result → Model8
```

## Security Architecture

### Current Security Model
- **Authentication:** Not implemented (CRITICAL GAP)
- **Authorization:** Not implemented (CRITICAL GAP)
- **Input Validation:** Minimal (HIGH RISK)
- **Credential Management:** Plain text in config (CRITICAL RISK)

### Target Security Model
```
┌─────────────────┐
│   API Gateway   │
│                 │
│ • Authentication│
│ • Rate Limiting │
│ • Input Valid.  │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐
│  Auth Service   │
│                 │
│ • JWT Tokens    │
│ • RBAC          │
│ • Session Mgmt  │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐
│  Core Services  │
│                 │
│ • Authorized    │
│ • Validated     │
│ • Audited       │
└─────────────────┘
```

## Performance Architecture

### Current Performance Characteristics
- **Database:** Single connection per operation (BOTTLENECK)
- **HTTP Client:** No connection pooling (BOTTLENECK)
- **Processing:** Sequential scan processing (BOTTLENECK)
- **Memory:** Unbounded slice growth (RISK)

### Target Performance Model
```
┌─────────────────┐
│ Load Balancer   │
└─────────┬───────┘
          │
     ┌────┴────┐
     │         │
┌────▼───┐ ┌──▼────┐
│API Inst│ │API Inst│
│1       │ │2      │
└────┬───┘ └──┬────┘
     │        │
     └────┬───┘
          │
┌─────────▼───────┐
│ Connection Pool │
│                 │
│ • DB Pool       │
│ • HTTP Pool     │
│ • Worker Pool   │
└─────────────────┘
```

## Deployment Architecture

### Current Deployment
- Single binary deployment
- Manual configuration
- No container optimization
- No health checks

### Target Deployment Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                      │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Ingress   │  │  ConfigMap  │  │   Secrets   │        │
│  │             │  │             │  │             │        │
│  │ • SSL/TLS   │  │ • App Config│  │ • DB Creds  │        │
│  │ • Routing   │  │ • Features  │  │ • API Keys  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  NuM8 Pods  │  │  PostgreSQL │  │  RabbitMQ   │        │
│  │             │  │             │  │             │        │
│  │ • API       │  │ • Primary   │  │ • Cluster   │        │
│  │ • Workers   │  │ • Replica   │  │ • HA        │        │
│  │ • Health    │  │ • Backups   │  │ • Mirroring │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

## Monitoring and Observability Architecture

### Target Observability Stack
```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Prometheus    │  │     Grafana     │  │    Jaeger       │
│                 │  │                 │  │                 │
│ • Metrics       │  │ • Dashboards    │  │ • Tracing       │
│ • Alerting      │  │ • Alerts        │  │ • Dependencies  │
│ • Scraping      │  │ • Visualization │  │ • Performance   │
└─────────┬───────┘  └─────────┬───────┘  └─────────┬───────┘
          │                    │                    │
          └────────────────────┼────────────────────┘
                               │
                    ┌──────────▼──────────┐
                    │      NuM8 App       │
                    │                     │
                    │ • /metrics          │
                    │ • /health           │
                    │ • /debug/pprof      │
                    │ • Trace Headers     │
                    └─────────────────────┘
```

## Configuration Architecture

### Current Configuration
- YAML file with embedded secrets (CRITICAL ISSUE)
- No environment variable support
- No validation or hot reload

### Target Configuration Model
```
┌─────────────────┐
│  Environment    │
│                 │
│ • Variables     │
│ • Secrets       │
│ • Overrides     │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐
│ Config Service  │
│                 │
│ • Validation    │
│ • Hot Reload    │
│ • Versioning    │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐
│ Application     │
│                 │
│ • Typed Config  │
│ • Runtime       │
│ • Feature Flags │
└─────────────────┘
```

## Integration Points

### External Tool Integration
- **Nuclei:** Binary execution with template management
- **Burp Suite:** REST API integration (HTTP)
- **Discord:** Webhook integration for notifications
- **PostgreSQL:** Database persistence
- **RabbitMQ:** Message queue integration

### API Integration Points
- **Scan Execution:** Synchronous and asynchronous scanning
- **Result Retrieval:** Real-time and historical data access
- **Configuration:** Runtime configuration management
- **Monitoring:** Health checks and metrics endpoints

## Future Architecture Considerations

### Microservices Evolution
- Split monolith into focused services
- Service mesh implementation (Istio)
- Independent scaling and deployment
- Distributed configuration management

### Cloud-Native Features
- Auto-scaling based on queue depth
- Circuit breakers and bulkheads
- Distributed caching (Redis)
- Event-driven architecture

### Advanced Security
- Zero-trust network model
- Service-to-service authentication
- Encrypted inter-service communication
- Advanced threat detection

This architecture documentation provides a comprehensive view of the NuM8 system design, current state, and future evolution path. It serves as a reference for development, deployment, and operational decisions.