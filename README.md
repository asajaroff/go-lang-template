# Go Lang Template - Microservice Template

A production-ready Go microservice template with configuration management, health checks, database migrations, structured logging, and API documentation.

## Features

- **Configuration Management**: Load configs from files, environment variables, or CLI flags using Viper and Cobra
- **Health Checks**: Built-in liveness and readiness probes for Kubernetes deployments
- **Database Support**: Optional PostgreSQL backend with connection pooling
- **Database Migrations**: Version-controlled schema migrations using golang-migrate
- **Structured Logging**: JSON and text format logging with configurable levels
- **API Documentation**: Auto-generated Swagger/OpenAPI documentation
- **Development Tools**: Pre-commit hooks for code quality and conventional commits
- **Testing**: Comprehensive test coverage with integration tests

## Prerequisites

- Go 1.24.4 or later
- PostgreSQL (optional, only if using database features)
- Make (for using Makefile commands)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/asajaroff/go-lang-template.git
cd go-lang-template
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Install Development Tools

```bash
make install-tools
```

This installs:
- `swag` - Swagger documentation generator

### 4. Set Up Pre-commit Hooks (Optional but Recommended)

```bash
make setup-hooks
```

This installs git hooks that enforce:
- Conventional commit messages
- Code formatting
- Static analysis
- Running tests before commits

### 5. Build the Application

```bash
make build
```

The binary will be created in `bin/server`.

### 6. Run the Application

```bash
# Run with default settings (port 8080, no database)
./bin/server

# Or run directly with go
go run ./cmd/server
```

The server will start on `http://localhost:8080`.

## Configuration

The application can be configured using three methods (in order of precedence):

1. **CLI Flags** (highest priority)
2. **Environment Variables** (prefix: `APP_`)
3. **Configuration File** (YAML format)

### CLI Flags

```bash
./bin/server --help

Available flags:
  --config string           Config file path (default "config.yaml")
  --server-host string      Server host (default "localhost")
  --server-port int         Server port (default 8080)
  --db-enabled              Enable database connection (default false)
  --db-driver string        Database driver (default "postgres")
  --db-host string          Database host (default "localhost")
  --db-port int             Database port (default 5432)
  --db-name string          Database name (default "mydb")
  --db-user string          Database user (default "postgres")
  --db-password string      Database password
  --db-sslmode string       Database SSL mode (default "disable")
  --log-level string        Log level: debug, info, warn, error (default "info")
  --log-format string       Log format: json, text (default "json")
  --migrations-path string  Path to migration files (default "./migrations")
```

### Environment Variables

All flags can be set via environment variables with the `APP_` prefix:

```bash
export APP_SERVER_PORT=8080
export APP_DB_ENABLED=true
export APP_DB_HOST=localhost
export APP_DB_PORT=5432
export APP_DB_NAME=mydb
export APP_DB_USER=postgres
export APP_DB_PASSWORD=secretpassword
export APP_LOG_LEVEL=info
export APP_LOG_FORMAT=json
```

### Configuration File

Create a `config.yaml` file:

```yaml
server:
  host: localhost
  port: 8080

database:
  enabled: true
  driver: postgres
  host: localhost
  port: 5432
  name: mydb
  user: postgres
  password: secretpassword
  sslmode: disable

log:
  level: info
  format: json

migrations:
  path: ./migrations
```

Then run with:

```bash
./bin/server --config config.yaml
```

## Running with PostgreSQL

### 1. Start PostgreSQL (using Docker)

```bash
docker run --name postgres-dev \
  -e POSTGRES_PASSWORD=secretpassword \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -d postgres:latest
```

### 2. Run Database Migrations

```bash
./bin/server migrate up --db-enabled --db-password=secretpassword
```

Available migration commands:

```bash
# Apply all pending migrations
./bin/server migrate up

# Rollback the last migration
./bin/server migrate down

# Apply or rollback N migrations
./bin/server migrate steps 2    # Apply 2 migrations
./bin/server migrate steps -1   # Rollback 1 migration

# Check current migration version
./bin/server migrate version

# Force set migration version (use with caution)
./bin/server migrate force 1
```

### 3. Run the Server with Database

```bash
./bin/server --db-enabled --db-password=secretpassword
```

## API Endpoints

### Health Checks

- **Liveness Probe**: `GET /healthz/liveness`
  - Always returns `200 OK` if server is running
  - Used by Kubernetes to restart unhealthy pods

- **Readiness Probe**: `GET /healthz/readiness`
  - Returns `200 OK` if server is ready to accept traffic
  - Checks database connection if enabled
  - Used by Kubernetes to route traffic

### API Documentation

- **Swagger UI**: `GET /swagger/index.html`
  - Interactive API documentation
  - Available after running `make swagger`

## Development

### Project Structure

```
.
├── cmd/
│   └── server/          # Main application and CLI commands
│       ├── main.go      # Application entry point
│       ├── root.go      # Root CLI command
│       └── migrate.go   # Migration CLI commands
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── database/        # Database connection handling
│   ├── health/          # Health check handlers
│   ├── logger/          # Structured logging
│   └── migrate/         # Migration management
├── migrations/          # Database migration files
├── docs/                # Generated Swagger documentation
├── pkg/                 # Public libraries (if any)
└── scripts/             # Build and utility scripts
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

### Linting

```bash
# Run all linters
make lint

# Individual checks
go fmt ./...
go vet ./...
```

### Generate Swagger Documentation

After modifying API annotations:

```bash
make swagger
```

### Available Make Commands

```bash
make help          # Show available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests with coverage
make lint          # Run linters
make fmt           # Format code
make swagger       # Generate Swagger docs
make clean         # Clean build artifacts
make all           # Lint, test, and build
```

## Creating Database Migrations

### 1. Create Migration Files

Create two files in the `migrations/` directory:

- `<version>_<description>.up.sql` - Apply changes
- `<version>_<description>.down.sql` - Revert changes

Example:

```sql
-- 000002_add_email_index.up.sql
CREATE INDEX idx_users_email ON users(email);

-- 000002_add_email_index.down.sql
DROP INDEX IF EXISTS idx_users_email;
```

### 2. Apply Migrations

```bash
./bin/server migrate up
```

## Deployment

### Docker

Build a Docker image:

```bash
docker build -t go-lang-template:latest .
```

Run the container:

```bash
docker run -p 8080:8080 \
  -e APP_DB_ENABLED=true \
  -e APP_DB_HOST=postgres \
  -e APP_DB_PASSWORD=secretpassword \
  go-lang-template:latest
```

### Kubernetes

The health check endpoints are designed for Kubernetes:

```yaml
livenessProbe:
  httpGet:
    path: /healthz/liveness
    port: 8080
  initialDelaySeconds: 3
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /healthz/readiness
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on:
- Commit message conventions (Conventional Commits)
- Development workflow
- Code style guidelines
- Testing requirements

## License

Apache 2.0

## Support

For issues and questions:
- Open an issue on GitHub
- Check the Swagger documentation at `/swagger/index.html`
- Review the [CONTRIBUTING.md](CONTRIBUTING.md) guide
