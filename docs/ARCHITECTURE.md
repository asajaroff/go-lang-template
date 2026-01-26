# Architecture Documentation

This document provides a comprehensive overview of the Go microservice template architecture, including component interactions, data flow, and deployment topology.

## Table of Contents

- [System Overview](#system-overview)
- [High-Level Architecture](#high-level-architecture)
- [Component Architecture](#component-architecture)
- [Data Flow](#data-flow)
- [Deployment Architecture](#deployment-architecture)
- [Configuration Flow](#configuration-flow)
- [Health Check Flow](#health-check-flow)
- [Database Migration Flow](#database-migration-flow)

## System Overview

The Go microservice template is a production-ready foundation for building cloud-native microservices. It follows the 12-factor app methodology and provides essential features for enterprise deployments.

### Key Architectural Principles

1. **Separation of Concerns**: Clear separation between commands, internal logic, and external interfaces
2. **Dependency Injection**: Configurable dependencies for testing and flexibility
3. **Graceful Degradation**: Optional database connectivity with health checks
4. **Cloud-Native**: Designed for Kubernetes with proper health probes and logging
5. **Configuration Hierarchy**: Multiple configuration sources with clear precedence

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         External Clients                             │
│                    (Load Balancer, Users, K8s)                      │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             │ HTTP Requests
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Go Microservice Server                          │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │                     HTTP Server (port 8080)                    │  │
│  │                      (Graceful Shutdown)                       │  │
│  └───────────────────────────┬───────────────────────────────────┘  │
│                              │                                       │
│  ┌───────────────────────────┴───────────────────────────────────┐  │
│  │                       API Endpoints                            │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌───────────────────┐   │  │
│  │  │  /healthz/   │  │  /healthz/   │  │   /swagger/       │   │  │
│  │  │  liveness    │  │  readiness   │  │   index.html      │   │  │
│  │  └──────┬───────┘  └──────┬───────┘  └───────────────────┘   │  │
│  └─────────┼──────────────────┼─────────────────────────────────┘  │
│            │                  │                                     │
│  ┌─────────▼──────────────────▼─────────────────────────────────┐  │
│  │                    Internal Packages                          │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐   │  │
│  │  │  Health  │  │  Config  │  │  Logger  │  │  Database  │   │  │
│  │  │  Module  │  │  Module  │  │  Module  │  │   Module   │   │  │
│  │  └────┬─────┘  └────┬─────┘  └────┬─────┘  └─────┬──────┘   │  │
│  │       │             │             │              │           │  │
│  │       │             │             ▼              │           │  │
│  │       │             │      ┌─────────────┐       │           │  │
│  │       │             │      │   stdout/   │       │           │  │
│  │       │             │      │   stderr    │       │           │  │
│  │       │             │      └─────────────┘       │           │  │
│  │       │             │                            │           │  │
│  └───────┼─────────────┼────────────────────────────┼───────────┘  │
│          │             │                            │               │
└──────────┼─────────────┼────────────────────────────┼───────────────┘
           │             │                            │
           │             │                            │ SQL Connection
           │             │                            ▼
           │             │                  ┌──────────────────┐
           │             │                  │   PostgreSQL     │
           │             │                  │    Database      │
           │             │                  └──────────────────┘
           │             │
           │             ▼
           │    ┌─────────────────┐
           │    │  Configuration  │
           │    │    Sources      │
           │    │  • CLI Flags    │
           │    │  • Env Vars     │
           │    │  • YAML Files   │
           │    └─────────────────┘
           │
           ▼
    ┌──────────────┐
    │   Database   │
    │  Connection  │
    │    Status    │
    └──────────────┘
```

## Component Architecture

### Command Layer (`cmd/server/`)

The command layer provides the CLI interface and application entry point.

```
┌─────────────────────────────────────────────────────────┐
│                    cmd/server/                          │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  main.go                                          │ │
│  │  • Application entry point                        │ │
│  │  • Executes root command                          │ │
│  └─────────────────┬─────────────────────────────────┘ │
│                    │                                   │
│  ┌─────────────────▼─────────────────────────────────┐ │
│  │  root.go                                          │ │
│  │  • Root command definition (Cobra)                │ │
│  │  • Server initialization                          │ │
│  │  • Graceful shutdown handling                     │ │
│  │  • Configuration loading                          │ │
│  │  • HTTP endpoint registration                     │ │
│  └─────────────────┬─────────────────────────────────┘ │
│                    │                                   │
│  ┌─────────────────▼─────────────────────────────────┐ │
│  │  migrate.go                                       │ │
│  │  • Migration subcommands                          │ │
│  │  • Database version management                    │ │
│  │  • Migration execution                            │ │
│  └───────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### Internal Layer (`internal/`)

The internal layer contains private application logic.

```
┌─────────────────────────────────────────────────────────────┐
│                       internal/                             │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  config/                                              │ │
│  │  • Configuration struct definitions                   │ │
│  │  • Viper integration for file/environment variables/flag loading │ │
│  │  • Cobra flag binding                                 │ │
│  │  • Configuration validation                           │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  health/                                              │ │
│  │  • HTTP health check handlers                         │ │
│  │  • Checker interface for dependency validation        │ │
│  │  • Liveness probe (always healthy)                    │ │
│  │  • Readiness probe (checks dependencies)              │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  logger/                                              │ │
│  │  • Structured logging (log/slog)                      │ │
│  │  • JSON and text format support                       │ │
│  │  • Level-based routing (stdout/stderr)                │ │
│  │  • Configurable log levels                            │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  database/                                            │ │
│  │  • Database connection management                     │ │
│  │  • DSN construction                                   │ │
│  │  • Connection pooling (sql.DB)                        │ │
│  │  • Health check implementation                        │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  migrate/                                             │ │
│  │  • golang-migrate wrapper                             │ │
│  │  • Migration operations (up/down/steps)               │ │
│  │  • Version management                                 │ │
│  │  • Migration file handling                            │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

### Request Processing Flow

```
1. Client Request
   │
   ▼
2. HTTP Server (net/http)
   │
   ▼
3. Route Matching
   │
   ├─► /healthz/liveness  ──► Liveness Handler ──► 200 OK
   │
   ├─► /healthz/readiness ──► Readiness Handler
   │                           │
   │                           ├─► Check Database Connection
   │                           │   (if enabled)
   │                           │
   │                           └─► Return 200 OK or 503 Service Unavailable
   │
   └─► /swagger/* ──► Swagger UI ──► API Documentation
```

### Configuration Loading Flow

```
Application Startup
   │
   ▼
1. Parse CLI Flags (Cobra)
   │
   ▼
2. Load Configuration File (if specified)
   │   • Read YAML file
   │   • Parse configuration
   │
   ▼
3. Load Environment Variables
   │   • Read APP_* prefixed variables
   │
   ▼
4. Merge Configuration (Viper)
   │   • Priority: CLI Flags > Env Vars > Config File > Defaults
   │
   ▼
5. Validate Configuration
   │   • Check required fields
   │   • Validate values
   │
   ▼
6. Return Config Struct
```

## Deployment Architecture

### Docker Deployment

```
┌─────────────────────────────────────────────────────────────────┐
│                         Docker Host                             │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Go Microservice Container                              │   │
│  │  ┌───────────────────────────────────────────────────┐  │   │
│  │  │  Application (port 8080)                          │  │   │
│  │  └─────────────────┬─────────────────────────────────┘  │   │
│  │                    │                                     │   │
│  │  ┌─────────────────▼─────────────────────────────────┐  │   │
│  │  │  Environment Variables                            │  │   │
│  │  │  • APP_DB_HOST                                    │  │   │
│  │  │  • APP_DB_PASSWORD                                │  │   │
│  │  │  • APP_LOG_LEVEL                                  │  │   │
│  │  └───────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                 │
└──────────────────────────────┼─────────────────────────────────┘
                               │
                               │ TCP/IP
                               ▼
                     ┌──────────────────┐
                     │   PostgreSQL     │
                     │   Container      │
                     │   (port 5432)    │
                     └──────────────────┘
```

### Kubernetes Deployment

```
┌───────────────────────────────────────────────────────────────────┐
│                        Kubernetes Cluster                         │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                        Ingress                               │ │
│  │                    (External Access)                         │ │
│  └────────────────────────┬─────────────────────────────────────┘ │
│                           │                                       │
│  ┌────────────────────────▼─────────────────────────────────────┐ │
│  │                      Service (ClusterIP)                     │ │
│  │                Load Balances to Pods                         │ │
│  └────────────────────────┬─────────────────────────────────────┘ │
│                           │                                       │
│  ┌────────────────────────▼─────────────────────────────────────┐ │
│  │                     Deployment                               │ │
│  │                                                              │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │ │
│  │  │   Pod 1      │  │   Pod 2      │  │   Pod 3      │      │ │
│  │  │              │  │              │  │              │      │ │
│  │  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │      │ │
│  │  │ │   App    │ │  │ │   App    │ │  │ │   App    │ │      │ │
│  │  │ │Container │ │  │ │Container │ │  │ │Container │ │      │ │
│  │  │ └────┬─────┘ │  │ └────┬─────┘ │  │ └────┬─────┘ │      │ │
│  │  │      │       │  │      │       │  │      │       │      │ │
│  │  │      │ Probes│  │      │ Probes│  │      │ Probes│      │ │
│  │  │      ├───────┤  │      ├───────┤  │      ├───────┤      │ │
│  │  │      │Liveness  │      │Liveness  │      │Liveness      │ │
│  │  │      │Readiness │      │Readiness │      │Readiness     │ │
│  │  └──────┼───────┘  └──────┼───────┘  └──────┼───────┘      │ │
│  │         │                 │                 │              │ │
│  └─────────┼─────────────────┼─────────────────┼──────────────┘ │
│            │                 │                 │                │
│  ┌─────────▼─────────────────▼─────────────────▼──────────────┐ │
│  │                      ConfigMap                             │ │
│  │               (Non-sensitive config)                       │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                       Secret                               │ │
│  │            (Database password, API keys)                   │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │              PostgreSQL StatefulSet                        │ │
│  │                                                            │ │
│  │  ┌──────────────────────────────────────────────────────┐ │ │
│  │  │  PostgreSQL Pod                                      │ │ │
│  │  │  ┌────────────────────────────────────────────────┐  │ │ │
│  │  │  │  PostgreSQL Container                          │  │ │ │
│  │  │  └───────────────┬────────────────────────────────┘  │ │ │
│  │  │                  │                                   │ │ │
│  │  │  ┌───────────────▼────────────────────────────────┐  │ │ │
│  │  │  │  Persistent Volume Claim                       │  │ │ │
│  │  │  │  (Database storage)                            │  │ │ │
│  │  │  └────────────────────────────────────────────────┘  │ │ │
│  │  └──────────────────────────────────────────────────────┘ │ │
│  └────────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
```

## Configuration Flow

### Configuration Precedence Diagram

```
┌─────────────────────────────────────────────────┐
│          Configuration Sources                  │
│         (Order of Precedence)                   │
└─────────────────────────────────────────────────┘
         │
         │ 1. Highest Priority
         ▼
┌─────────────────────────────────────────────────┐
│          CLI Flags                              │
│  • --server-port 9090                           │
│  • --db-enabled                                 │
│  • --log-level debug                            │
└─────────────────┬───────────────────────────────┘
                  │
                  │ Overrides
                  ▼
┌─────────────────────────────────────────────────┐
│      Environment Variables                      │
│  • APP_SERVER_PORT=8080                         │
│  • APP_DB_ENABLED=true                          │
│  • APP_LOG_LEVEL=info                           │
└─────────────────┬───────────────────────────────┘
                  │
                  │ Overrides
                  ▼
┌─────────────────────────────────────────────────┐
│       Configuration File (YAML)                 │
│  server:                                        │
│    port: 8080                                   │
│  database:                                      │
│    enabled: false                               │
│  log:                                           │
│    level: warn                                  │
└─────────────────┬───────────────────────────────┘
                  │
                  │ Overrides
                  ▼
┌─────────────────────────────────────────────────┐
│         Default Values                          │
│  • server.port = 8080                           │
│  • database.enabled = false                     │
│  • log.level = "info"                           │
└─────────────────┬───────────────────────────────┘
                  │
                  │ 4. Lowest Priority
                  ▼
┌─────────────────────────────────────────────────┐
│         Final Configuration                     │
│        (Merged Result)                          │
└─────────────────────────────────────────────────┘
```

## Health Check Flow

### Readiness Probe Flow

```
┌─────────────────────────────────────────────────────────────┐
│          Kubernetes / Load Balancer                         │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ GET /healthz/readiness
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Readiness Handler                              │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ Check all registered checkers
                         ▼
              ┌──────────────────────┐
              │  Checker Interface   │
              └──────────┬───────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
┌────────────────┐ ┌────────────┐ ┌────────────┐
│   Database     │ │   Redis    │ │   Other    │
│   Checker      │ │   Checker  │ │   Checker  │
│  (if enabled)  │ │ (optional) │ │ (optional) │
└────────┬───────┘ └──────┬─────┘ └──────┬─────┘
         │                │              │
         │ db.Ping()      │ redis.Ping() │ custom check
         ▼                ▼              ▼
┌────────────────┐ ┌────────────┐ ┌────────────┐
│   healthy      │ │  healthy   │ │  healthy   │
└────────┬───────┘ └──────┬─────┘ └──────┬─────┘
         │                │              │
         └────────────────┴──────────────┘
                         │
                         │ All healthy?
                         ▼
                  ┌──────────────┐
                  │     Yes      │
                  └──────┬───────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  HTTP 200 OK                                                │
│  Content-Type: application/json                             │
│  {                                                          │
│    "status": "healthy",                                     │
│    "checks": {                                              │
│      "database": "healthy"                                  │
│    }                                                        │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘

                  ┌──────────────┐
                  │      No      │
                  └──────┬───────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  HTTP 503 Service Unavailable                              │
│  Content-Type: application/json                             │
│  {                                                          │
│    "status": "unhealthy",                                   │
│    "checks": {                                              │
│      "database": "connection failed: ..."                   │
│    }                                                        │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
```

## Database Migration Flow

### Migration Execution Flow

```
┌─────────────────────────────────────────────────┐
│  CLI: ./bin/server migrate up                   │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  1. Load Configuration                          │
│     • Database credentials                      │
│     • Migrations path                           │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  2. Connect to Database                         │
│     • Open SQL connection                       │
│     • Verify connectivity                       │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  3. Initialize Migrator (golang-migrate)        │
│     • Load migration files from disk            │
│     • Create schema_migrations table            │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  4. Check Current Version                       │
│     • Query schema_migrations table             │
│     • Determine pending migrations              │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  5. Execute Migrations                          │
│                                                 │
│  For each pending migration:                   │
│  ┌───────────────────────────────────────────┐ │
│  │  a. Begin transaction                     │ │
│  │  b. Execute .up.sql file                  │ │
│  │  c. Update schema_migrations version      │ │
│  │  d. Commit transaction                    │ │
│  └───────────────────────────────────────────┘ │
│                                                 │
│  If error occurs:                              │
│  ┌───────────────────────────────────────────┐ │
│  │  - Rollback transaction                   │ │
│  │  - Report error                           │ │
│  │  - Stop execution                         │ │
│  └───────────────────────────────────────────┘ │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  6. Report Results                              │
│     • Number of migrations applied              │
│     • Current version                           │
│     • Any errors encountered                    │
└─────────────────────────────────────────────────┘
```

### Migration Files Structure

```
migrations/
├── 000001_create_users_table.up.sql
│   • Creates users table
│   • Adds indexes
│   • Sets constraints
│
├── 000001_create_users_table.down.sql
│   • Drops indexes
│   • Drops users table
│   • Reverses all changes from .up.sql
│
├── 000002_add_email_index.up.sql
│   • Creates additional index on email
│
└── 000002_add_email_index.down.sql
    • Drops email index

Migration Version Tracking:
┌────────────────────────────────────────────┐
│       schema_migrations table              │
├────────────┬──────────────┬────────────────┤
│  version   │    dirty     │   sequence     │
├────────────┼──────────────┼────────────────┤
│     2      │    false     │      1         │
└────────────┴──────────────┴────────────────┘
```

## Graceful Shutdown Flow

```
┌─────────────────────────────────────────────────┐
│  OS Signal (SIGINT/SIGTERM)                     │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  1. Signal Handler Catches Signal               │
│     • Buffered channel receives signal          │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  2. Log Shutdown Initiation                     │
│     • "Shutting down server..."                 │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  3. Stop Accepting New Connections              │
│     • HTTP server stops listening               │
│     • Existing connections continue             │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  4. Wait for In-Flight Requests                 │
│     • Wait up to 30 seconds                     │
│     • Allow active requests to complete         │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  5. Close Database Connections                  │
│     • Close connection pool                     │
│     • Release resources                         │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  6. Exit Application                            │
│     • Clean exit with status code 0             │
└─────────────────────────────────────────────────┘
```

## Security Architecture

### Security Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    Security Layers                          │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  Layer 1: Network Security                            │ │
│  │  • Kubernetes Network Policies                        │ │
│  │  • Ingress TLS termination                            │ │
│  │  • Service mesh (optional)                            │ │
│  └───────────────────────────────────────────────────────┘ │
│                           │                                 │
│  ┌────────────────────────▼──────────────────────────────┐ │
│  │  Layer 2: Application Security                        │ │
│  │  • Secure configuration (no hardcoded secrets)        │ │
│  │  • Structured logging (no sensitive data)             │ │
│  │  • Graceful error handling                            │ │
│  └───────────────────────────────────────────────────────┘ │
│                           │                                 │
│  ┌────────────────────────▼──────────────────────────────┐ │
│  │  Layer 3: Data Security                               │ │
│  │  • Database SSL/TLS connections                       │ │
│  │  • Connection string encryption in secrets            │ │
│  │  • SQL injection prevention (prepared statements)     │ │
│  └───────────────────────────────────────────────────────┘ │
│                           │                                 │
│  ┌────────────────────────▼──────────────────────────────┐ │
│  │  Layer 4: Container Security                          │ │
│  │  • Non-root user execution                            │ │
│  │  • Read-only root filesystem                          │ │
│  │  • Minimal base image (distroless/scratch)            │ │
│  │  • Security scanning (trivy, grype)                   │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Observability Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  Observability Stack                        │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  Application                                          │ │
│  └────┬────────────────────┬────────────────────┬────────┘ │
│       │                    │                    │          │
│       │ Logs               │ Metrics            │ Traces   │
│       │                    │                    │          │
│       ▼                    ▼                    ▼          │
│  ┌──────────┐        ┌──────────┐        ┌──────────┐    │
│  │  stdout  │        │ Prometheus│       │  Jaeger  │    │
│  │  stderr  │        │ (future)  │       │ (future) │    │
│  └────┬─────┘        └──────────┘        └──────────┘    │
│       │                                                    │
│       ▼                                                    │
│  ┌──────────────────────────────────────────────────────┐ │
│  │  Log Aggregation (Future)                            │ │
│  │  • Fluentd / Fluent Bit                              │ │
│  │  • Elasticsearch                                     │ │
│  │  • Kibana                                            │ │
│  └──────────────────────────────────────────────────────┘ │
│                                                            │
│  ┌──────────────────────────────────────────────────────┐ │
│  │  Monitoring Dashboard (Future)                       │ │
│  │  • Grafana                                           │ │
│  │  • Prometheus Alerts                                 │ │
│  └──────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Technology Stack Summary

```
┌─────────────────────────────────────────────────────────────┐
│                    Technology Stack                         │
│                                                             │
│  Language & Runtime                                         │
│  ├─ Go 1.24.4+                                              │
│                                                             │
│  Core Libraries                                             │
│  ├─ net/http          (HTTP server)                         │
│  ├─ log/slog          (Structured logging)                  │
│  ├─ database/sql      (Database interface)                  │
│  └─ sql.DB            (Connection pooling)                  │
│                                                             │
│  Third-Party Dependencies                                   │
│  ├─ github.com/spf13/viper   (Configuration)                │
│  ├─ github.com/spf13/cobra   (CLI framework)                │
│  ├─ github.com/lib/pq        (PostgreSQL driver)            │
│  ├─ github.com/golang-migrate/migrate (DB migrations)       │
│  └─ github.com/swaggo/swag   (API documentation)            │
│                                                             │
│  Database                                                   │
│  └─ PostgreSQL 12+                                          │
│                                                             │
│  Container Runtime                                          │
│  └─ Docker                                                  │
│                                                             │
│  Orchestration                                              │
│  └─ Kubernetes 1.20+                                        │
│                                                             │
│  Development Tools                                          │
│  ├─ pre-commit        (Git hooks)                           │
│  ├─ golangci-lint     (Code linting)                        │
│  ├─ go test           (Testing)                             │
│  └─ swag              (Swagger generation)                  │
└─────────────────────────────────────────────────────────────┘
```

## Design Patterns and Best Practices

### Applied Design Patterns

1. **Dependency Injection**: Components receive dependencies through constructors
2. **Interface Segregation**: Small, focused interfaces (e.g., Checker interface)
3. **Factory Pattern**: Configuration and component creation functions
4. **Strategy Pattern**: Pluggable health checkers
5. **Singleton Pattern**: Logger and database connection instances
6. **Command Pattern**: CLI commands with Cobra

### Architectural Best Practices

1. **12-Factor App Compliance**
   - Configuration via environment
   - Stateless processes
   - Disposability (fast startup/shutdown)
   - Logs to stdout/stderr
   - Port binding for services

2. **Cloud-Native Principles**
   - Health checks for orchestration
   - Graceful shutdown handling
   - Configuration externalization
   - Horizontal scalability

3. **Code Organization**
   - Clear package boundaries
   - Internal package for privacy
   - Test files alongside source
   - Comprehensive test coverage

4. **Error Handling**
   - Explicit error returns
   - Context preservation
   - Structured error logging
   - No panic in production code

## Performance Considerations

### Connection Pooling

```
Database Connection Pool (sql.DB)
├── Max Open Connections: 25 (default)
├── Max Idle Connections: 25 (default)
├── Connection Lifetime: Unlimited
└── Connection Idle Time: Unlimited

Tuning recommendations:
- Set MaxOpenConns based on database capacity
- Set MaxIdleConns = MaxOpenConns for optimal reuse
- Set ConnMaxLifetime for connection rotation
- Monitor active/idle connection metrics
```

### HTTP Server Tuning

```
HTTP Server Configuration
├── Read Timeout: 15s (recommended)
├── Write Timeout: 15s (recommended)
├── Idle Timeout: 60s (recommended)
└── Max Header Bytes: 1MB (default)

Graceful Shutdown Timeout: 30s
```

## Future Enhancements

Potential architectural improvements:

1. **Metrics Collection**: Prometheus metrics endpoint
2. **Distributed Tracing**: OpenTelemetry integration
3. **Rate Limiting**: Per-endpoint request limiting
4. **Authentication**: JWT or OAuth2 middleware
5. **Authorization**: Role-based access control
6. **API Versioning**: URL-based or header-based versioning
7. **Caching Layer**: Redis integration
8. **Message Queue**: Async job processing
9. **Service Discovery**: Consul or etcd integration
10. **Circuit Breaker**: Resilience patterns for external dependencies

---

## References

- [12-Factor App](https://12factor.net/)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
