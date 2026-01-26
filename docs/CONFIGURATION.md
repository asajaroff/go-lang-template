# Configuration Guide

This guide provides comprehensive documentation for all configuration options available in the Go Lang Template microservice.

## Table of Contents

- [Configuration Sources](#configuration-sources)
- [Priority Order](#priority-order)
- [Server Configuration](#server-configuration)
- [Database Configuration](#database-configuration)
- [Logging Configuration](#logging-configuration)
- [Configuration File Format](#configuration-file-format)
- [Environment Variables Reference](#environment-variables-reference)
- [CLI Flags Reference](#cli-flags-reference)
- [Configuration Examples](#configuration-examples)

## Configuration Sources

The application supports three configuration sources:

1. **CLI Flags** - Command-line arguments passed when starting the application
2. **Environment Variables** - System environment variables with `APP_` prefix
3. **Configuration File** - YAML file with structured configuration

## Priority Order

Configuration values are loaded with the following priority (highest to lowest):

1. CLI Flags (highest priority)
2. Environment Variables
3. Configuration File
4. Default Values (lowest priority)

This means that if a value is specified in multiple sources, the CLI flag will always take precedence.

## Server Configuration

Server configuration controls the HTTP server behavior.

### Available Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `server.host` | string | `0.0.0.0` | The host address the server listens on |
| `server.port` | int | `8080` | The port number the server listens on |

### Usage Examples

**CLI Flag:**
```bash
./server --server-host=127.0.0.1 --server-port=3000
```

**Environment Variable:**
```bash
export APP_SERVER_HOST=127.0.0.1
export APP_SERVER_PORT=3000
./server
```

**Config File:**
```yaml
server:
  host: 127.0.0.1
  port: 3000
```

## Database Configuration

Database configuration controls the PostgreSQL database connection. The database is optional and must be explicitly enabled.

### Available Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `database.enabled` | bool | `false` | Enable or disable database connection |
| `database.driver` | string | `postgres` | Database driver (currently only postgres supported) |
| `database.host` | string | `localhost` | Database server hostname |
| `database.port` | int | `5432` | Database server port |
| `database.name` | string | `app` | Database name |
| `database.user` | string | `postgres` | Database username |
| `database.password` | string | `""` | Database password |
| `database.sslmode` | string | `disable` | SSL mode for database connection |
| `database.migrationspath` | string | `./migrations` | Path to migration files |

### SSL Modes

The following SSL modes are supported for PostgreSQL connections:

- `disable` - No SSL encryption (default, suitable for local development)
- `require` - SSL encryption required, no certificate verification
- `verify-ca` - SSL encryption required, verify CA certificate
- `verify-full` - SSL encryption required, verify CA and hostname

### Usage Examples

**CLI Flag:**
```bash
./server \
  --db-enabled=true \
  --db-host=db.example.com \
  --db-port=5432 \
  --db-name=myapp \
  --db-user=appuser \
  --db-password=secretpassword \
  --db-sslmode=require \
  --db-migrations-path=/app/migrations
```

**Environment Variable:**
```bash
export APP_DATABASE_ENABLED=true
export APP_DATABASE_HOST=db.example.com
export APP_DATABASE_PORT=5432
export APP_DATABASE_NAME=myapp
export APP_DATABASE_USER=appuser
export APP_DATABASE_PASSWORD=secretpassword
export APP_DATABASE_SSLMODE=require
export APP_DATABASE_MIGRATIONSPATH=/app/migrations
./server
```

**Config File:**
```yaml
database:
  enabled: true
  driver: postgres
  host: db.example.com
  port: 5432
  name: myapp
  user: appuser
  password: secretpassword
  sslmode: require
  migrationspath: /app/migrations
```

## Logging Configuration

Logging configuration controls the application's logging behavior.

### Available Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `log.level` | string | `info` | Logging level (debug, info, warn, error) |
| `log.format` | string | `json` | Log output format (json, text) |

### Log Levels

The following log levels are supported (from most to least verbose):

- `debug` - Detailed debug information
- `info` - Informational messages (default)
- `warn` - Warning messages
- `error` - Error messages only

### Log Formats

- `json` - Structured JSON format (default, recommended for production)
- `text` - Human-readable text format (useful for local development)

### Log Output Routing

Logs are automatically routed based on severity:
- `INFO` and `DEBUG` logs → stdout
- `WARN` and `ERROR` logs → stderr

### Usage Examples

**CLI Flag:**
```bash
./server --log-level=debug --log-format=text
```

**Environment Variable:**
```bash
export APP_LOG_LEVEL=debug
export APP_LOG_FORMAT=text
./server
```

**Config File:**
```yaml
log:
  level: debug
  format: text
```

## Configuration File Format

The application looks for configuration files in the following locations (in order):

1. Path specified with `--config` flag
2. `./config.yaml` (current directory)
3. `./config/config.yaml` (config subdirectory)
4. `/etc/go-lang-template/config.yaml` (system-wide)

### Complete Configuration File Example

```yaml
# Server configuration
server:
  host: 0.0.0.0
  port: 8080

# Database configuration
database:
  enabled: true
  driver: postgres
  host: localhost
  port: 5432
  name: app
  user: postgres
  password: mypassword
  sslmode: disable
  migrationspath: ./migrations

# Logging configuration
log:
  level: info
  format: json
```

## Environment Variables Reference

All environment variables use the `APP_` prefix and follow a hierarchical naming structure with underscores.

### Complete List

```bash
# Server
APP_SERVER_HOST=0.0.0.0
APP_SERVER_PORT=8080

# Database
APP_DATABASE_ENABLED=false
APP_DATABASE_DRIVER=postgres
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_NAME=app
APP_DATABASE_USER=postgres
APP_DATABASE_PASSWORD=
APP_DATABASE_SSLMODE=disable
APP_DATABASE_MIGRATIONSPATH=./migrations

# Logging
APP_LOG_LEVEL=info
APP_LOG_FORMAT=json
```

## CLI Flags Reference

### Global Flags

```bash
--config string
    Config file path (default searches for config.yaml in standard locations)
```

### Server Flags

```bash
--server-host string
    Server host address (env: APP_SERVER_HOST) (default "0.0.0.0")

--server-port int
    Server port number (env: APP_SERVER_PORT) (default 8080)
```

### Database Flags

```bash
--db-enabled
    Enable database connection (env: APP_DATABASE_ENABLED) (default false)

--db-driver string
    Database driver (env: APP_DATABASE_DRIVER) (default "postgres")

--db-host string
    Database host (env: APP_DATABASE_HOST) (default "localhost")

--db-port int
    Database port (env: APP_DATABASE_PORT) (default 5432)

--db-name string
    Database name (env: APP_DATABASE_NAME) (default "app")

--db-user string
    Database user (env: APP_DATABASE_USER) (default "postgres")

--db-password string
    Database password (env: APP_DATABASE_PASSWORD)

--db-sslmode string
    Database SSL mode (env: APP_DATABASE_SSLMODE) (default "disable")

--db-migrations-path string
    Database migrations path (env: APP_DATABASE_MIGRATIONSPATH) (default "./migrations")
```

### Logging Flags

```bash
--log-level string
    Log level: debug, info, warn, error (env: APP_LOG_LEVEL) (default "info")

--log-format string
    Log format: json, text (env: APP_LOG_FORMAT) (default "json")
```

## Configuration Examples

### Development Environment

**Local development with text logging:**

```bash
./server --log-format=text --log-level=debug
```

Or with environment variables:

```bash
export APP_LOG_FORMAT=text
export APP_LOG_LEVEL=debug
./server
```

### Production Environment

**Production with database:**

```bash
./server \
  --db-enabled=true \
  --db-host=prod-db.example.com \
  --db-name=production_db \
  --db-user=prod_user \
  --db-password="${DB_PASSWORD}" \
  --db-sslmode=verify-full \
  --log-level=info \
  --log-format=json
```

Or use a configuration file:

```yaml
# config.yaml
server:
  host: 0.0.0.0
  port: 8080

database:
  enabled: true
  driver: postgres
  host: prod-db.example.com
  port: 5432
  name: production_db
  user: prod_user
  password: ${DB_PASSWORD}  # Can reference environment variables
  sslmode: verify-full
  migrationspath: /app/migrations

log:
  level: info
  format: json
```

Then run:

```bash
export DB_PASSWORD=secretpassword
./server --config=config.yaml
```

### Docker Environment

**Using environment variables with Docker:**

```bash
docker run -d \
  -e APP_SERVER_PORT=8080 \
  -e APP_DATABASE_ENABLED=true \
  -e APP_DATABASE_HOST=postgres \
  -e APP_DATABASE_PORT=5432 \
  -e APP_DATABASE_NAME=app \
  -e APP_DATABASE_USER=postgres \
  -e APP_DATABASE_PASSWORD=postgres \
  -e APP_LOG_LEVEL=info \
  -e APP_LOG_FORMAT=json \
  -p 8080:8080 \
  go-lang-template
```

### Kubernetes Environment

**Using ConfigMap and Secrets:**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_SERVER_HOST: "0.0.0.0"
  APP_SERVER_PORT: "8080"
  APP_DATABASE_ENABLED: "true"
  APP_DATABASE_HOST: "postgres-service"
  APP_DATABASE_PORT: "5432"
  APP_DATABASE_NAME: "app"
  APP_DATABASE_USER: "postgres"
  APP_DATABASE_SSLMODE: "require"
  APP_DATABASE_MIGRATIONSPATH: "/app/migrations"
  APP_LOG_LEVEL: "info"
  APP_LOG_FORMAT: "json"
---
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
stringData:
  APP_DATABASE_PASSWORD: "secretpassword"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-lang-template
  template:
    metadata:
      labels:
        app: go-lang-template
    spec:
      containers:
      - name: app
        image: go-lang-template:latest
        envFrom:
        - configMapRef:
            name: app-config
        - secretRef:
            name: app-secrets
        ports:
        - containerPort: 8080
```

## Security Considerations

### Sensitive Data

- **Never commit passwords or secrets** to version control
- Use environment variables or secrets management for sensitive data
- In Kubernetes, use Secrets instead of ConfigMaps for passwords
- Consider using external secrets management (e.g., HashiCorp Vault, AWS Secrets Manager)

### Database Security

- Always use SSL mode `require`, `verify-ca`, or `verify-full` in production
- Use strong, unique passwords for database users
- Limit database user permissions to only what's needed
- Consider using connection pooling and connection limits

### Network Security

- In production, consider binding to specific interfaces instead of `0.0.0.0`
- Use reverse proxies (nginx, Traefik) for TLS termination
- Implement network policies in Kubernetes to restrict traffic

## Troubleshooting

### Configuration Not Loading

1. Check configuration file path and permissions
2. Verify environment variable names (must use `APP_` prefix)
3. Check for typos in configuration keys
4. Enable debug logging: `--log-level=debug`

### Database Connection Issues

1. Verify database is accessible: `telnet db-host 5432`
2. Check database credentials
3. Verify SSL mode matches database server configuration
4. Check database logs for connection attempts
5. Ensure `database.enabled` is set to `true`

### Configuration Priority Issues

Remember the priority order:
1. CLI flags override everything
2. Environment variables override config file
3. Config file overrides defaults

If a value isn't being applied, check if it's being overridden by a higher-priority source.
