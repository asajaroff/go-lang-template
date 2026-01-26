# Production Deployment Guide

This guide covers deploying the Go microservice template to production environments using Docker, Docker Compose, and Kubernetes.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Build and Release](#build-and-release)
- [Docker Deployment](#docker-deployment)
- [Docker Compose Deployment](#docker-compose-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Production Configuration](#production-configuration)
- [Security Best Practices](#security-best-practices)
- [Monitoring and Logging](#monitoring-and-logging)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Tools

- Docker 20.10 or later
- Docker Compose 2.0 or later (for local/staging)
- kubectl 1.24 or later (for Kubernetes)
- Access to a container registry (Docker Hub, ECR, GCR, etc.)
- PostgreSQL 13 or later (production database)

### Environment Requirements

- Linux-based production environment (Ubuntu 20.04+, RHEL 8+, etc.)
- Minimum 2 CPU cores and 2GB RAM per instance
- Persistent storage for PostgreSQL data
- Network access to database and required services

## Build and Release

### 1. Version Your Release

Follow semantic versioning for releases:

```bash
# Tag your release
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### 2. Build Production Binary

```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/server ./cmd/server

# Verify binary
./bin/server --version
```

### 3. Build Docker Image

Create a production Dockerfile if not present:

```dockerfile
# Multi-stage build for minimal image size
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server ./cmd/server

# Final stage - minimal runtime image
FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz/liveness || exit 1

# Run application
ENTRYPOINT ["/app/server"]
```

Build and tag the image:

```bash
# Build image
docker build -t your-registry/go-lang-template:v1.0.0 .

# Tag as latest
docker tag your-registry/go-lang-template:v1.0.0 your-registry/go-lang-template:latest

# Push to registry
docker push your-registry/go-lang-template:v1.0.0
docker push your-registry/go-lang-template:latest
```

## Docker Deployment

### Standalone Container

Run the application as a standalone container:

```bash
docker run -d \
  --name go-microservice \
  --restart unless-stopped \
  -p 8080:8080 \
  -e APP_SERVER_HOST=0.0.0.0 \
  -e APP_SERVER_PORT=8080 \
  -e APP_DB_ENABLED=true \
  -e APP_DB_HOST=postgres-prod.example.com \
  -e APP_DB_PORT=5432 \
  -e APP_DB_NAME=production_db \
  -e APP_DB_USER=app_user \
  -e APP_DB_PASSWORD=secure_password_here \
  -e APP_DB_SSLMODE=require \
  -e APP_LOG_LEVEL=info \
  -e APP_LOG_FORMAT=json \
  your-registry/go-lang-template:v1.0.0
```

### Using Secrets

For sensitive data, use Docker secrets or mount files:

```bash
# Create secret file
echo "secure_password" > db_password.txt

# Run with secret mounted
docker run -d \
  --name go-microservice \
  -p 8080:8080 \
  -e APP_DB_PASSWORD_FILE=/run/secrets/db_password \
  -v $(pwd)/db_password.txt:/run/secrets/db_password:ro \
  your-registry/go-lang-template:v1.0.0
```

## Docker Compose Deployment

### Production Docker Compose

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  app:
    image: your-registry/go-lang-template:v1.0.0
    container_name: go-microservice
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      APP_SERVER_HOST: 0.0.0.0
      APP_SERVER_PORT: 8080
      APP_DB_ENABLED: "true"
      APP_DB_HOST: postgres
      APP_DB_PORT: 5432
      APP_DB_NAME: production_db
      APP_DB_USER: app_user
      APP_DB_PASSWORD: ${DB_PASSWORD}
      APP_DB_SSLMODE: require
      APP_LOG_LEVEL: info
      APP_LOG_FORMAT: json
    depends_on:
      postgres:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/healthz/liveness"]
      interval: 30s
      timeout: 3s
      start_period: 10s
      retries: 3
    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    container_name: postgres-db
    restart: unless-stopped
    environment:
      POSTGRES_DB: production_db
      POSTGRES_USER: app_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app_user -d production_db"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network

volumes:
  postgres-data:
    driver: local

networks:
  app-network:
    driver: bridge
```

Deploy with Docker Compose:

```bash
# Set database password
export DB_PASSWORD="secure_password_here"

# Deploy
docker-compose -f docker-compose.prod.yml up -d

# Run migrations
docker exec go-microservice /app/server migrate up

# Check logs
docker-compose -f docker-compose.prod.yml logs -f app

# Check health
curl http://localhost:8080/healthz/readiness
```

## Kubernetes Deployment

### Namespace Setup

Create a dedicated namespace:

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: go-microservice
```

```bash
kubectl apply -f namespace.yaml
```

### Secrets Management

Create secrets for sensitive data:

```yaml
# secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: go-microservice
type: Opaque
stringData:
  db-password: "secure_password_here"
  db-user: "app_user"
```

```bash
kubectl apply -f secrets.yaml
```

### ConfigMap

Create a ConfigMap for non-sensitive configuration:

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: go-microservice
data:
  server-host: "0.0.0.0"
  server-port: "8080"
  db-host: "postgres-service"
  db-port: "5432"
  db-name: "production_db"
  db-sslmode: "require"
  log-level: "info"
  log-format: "json"
```

```bash
kubectl apply -f configmap.yaml
```

### PostgreSQL Deployment

Deploy PostgreSQL with persistent storage:

```yaml
# postgres-deployment.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: go-microservice
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: go-microservice
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-name
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: db-user
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: db-password
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        livenessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - app_user
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - app_user
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
  namespace: go-microservice
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
  type: ClusterIP
```

```bash
kubectl apply -f postgres-deployment.yaml
```

### Application Deployment

Deploy the application:

```yaml
# app-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-microservice
  namespace: go-microservice
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-microservice
  template:
    metadata:
      labels:
        app: go-microservice
    spec:
      containers:
      - name: app
        image: your-registry/go-lang-template:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: APP_SERVER_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: server-host
        - name: APP_SERVER_PORT
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: server-port
        - name: APP_DB_ENABLED
          value: "true"
        - name: APP_DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-host
        - name: APP_DB_PORT
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-port
        - name: APP_DB_NAME
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-name
        - name: APP_DB_USER
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: db-user
        - name: APP_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: db-password
        - name: APP_DB_SSLMODE
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-sslmode
        - name: APP_LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: log-level
        - name: APP_LOG_FORMAT
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: log-format
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /healthz/liveness
            port: 8080
          initialDelaySeconds: 15
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /healthz/readiness
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
---
apiVersion: v1
kind: Service
metadata:
  name: go-microservice-service
  namespace: go-microservice
spec:
  selector:
    app: go-microservice
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

```bash
kubectl apply -f app-deployment.yaml
```

### Ingress Configuration

Setup Ingress for external access:

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-microservice-ingress
  namespace: go-microservice
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - api.example.com
    secretName: go-microservice-tls
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-microservice-service
            port:
              number: 80
```

```bash
kubectl apply -f ingress.yaml
```

### Run Database Migrations in Kubernetes

Create a Kubernetes Job for migrations:

```yaml
# migration-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: db-migration
  namespace: go-microservice
spec:
  template:
    spec:
      containers:
      - name: migrate
        image: your-registry/go-lang-template:v1.0.0
        command: ["/app/server", "migrate", "up"]
        env:
        - name: APP_DB_ENABLED
          value: "true"
        - name: APP_DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-host
        - name: APP_DB_PORT
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-port
        - name: APP_DB_NAME
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-name
        - name: APP_DB_USER
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: db-user
        - name: APP_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: db-password
        - name: APP_DB_SSLMODE
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db-sslmode
      restartPolicy: OnFailure
  backoffLimit: 3
```

```bash
kubectl apply -f migration-job.yaml
kubectl logs -n go-microservice -l job-name=db-migration
```

### Horizontal Pod Autoscaling

Enable autoscaling based on CPU/memory:

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: go-microservice-hpa
  namespace: go-microservice
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: go-microservice
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

```bash
kubectl apply -f hpa.yaml
```

## Production Configuration

### Environment-Specific Settings

#### Production
```yaml
server:
  host: 0.0.0.0
  port: 8080

database:
  enabled: true
  driver: postgres
  host: production-db.example.com
  port: 5432
  name: production_db
  user: app_user
  password: ${DB_PASSWORD}  # From secrets
  sslmode: require

log:
  level: info
  format: json

migrations:
  path: ./migrations
```

#### Staging
```yaml
server:
  host: 0.0.0.0
  port: 8080

database:
  enabled: true
  driver: postgres
  host: staging-db.example.com
  port: 5432
  name: staging_db
  user: app_user
  password: ${DB_PASSWORD}
  sslmode: require

log:
  level: debug
  format: json

migrations:
  path: ./migrations
```

## Security Best Practices

### 1. Container Security

- Run as non-root user
- Use read-only root filesystem
- Drop all capabilities
- Scan images for vulnerabilities regularly
- Use minimal base images (Alpine)

### 2. Secret Management

- Never commit secrets to version control
- Use Kubernetes Secrets or external secret managers (AWS Secrets Manager, HashiCorp Vault)
- Rotate secrets regularly
- Use different credentials per environment

### 3. Network Security

- Enable SSL/TLS for database connections (`sslmode: require`)
- Use network policies in Kubernetes to restrict traffic
- Enable HTTPS for all external endpoints
- Use private networks for internal services

### 4. Database Security

- Use strong passwords
- Limit database user permissions
- Enable SSL connections
- Regular backups
- Encrypt data at rest

### 5. Application Security

- Keep dependencies updated
- Enable rate limiting
- Implement proper authentication and authorization
- Validate all inputs
- Use security headers

## Monitoring and Logging

### Health Check Monitoring

Monitor the health endpoints:

```bash
# Liveness check
curl http://api.example.com/healthz/liveness

# Readiness check
curl http://api.example.com/healthz/readiness
```

### Log Aggregation

The application outputs structured JSON logs. Integrate with:

- **ELK Stack**: Elasticsearch, Logstash, Kibana
- **Loki**: With Grafana for visualization
- **CloudWatch**: For AWS deployments
- **Stackdriver**: For GCP deployments

### Metrics Collection

Integrate with monitoring systems:

- **Prometheus**: For metrics collection
- **Grafana**: For visualization
- **DataDog**: Comprehensive monitoring
- **New Relic**: Application performance monitoring

### Sample Prometheus ServiceMonitor

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: go-microservice
  namespace: go-microservice
spec:
  selector:
    matchLabels:
      app: go-microservice
  endpoints:
  - port: http
    path: /metrics
    interval: 30s
```

## Troubleshooting

### Common Issues

#### 1. Application Not Starting

**Check logs:**
```bash
# Docker
docker logs go-microservice

# Kubernetes
kubectl logs -n go-microservice -l app=go-microservice
```

**Common causes:**
- Database connection failure
- Missing environment variables
- Port already in use
- Permission issues

#### 2. Database Connection Errors

**Verify connectivity:**
```bash
# Test database connection
psql -h production-db.example.com -U app_user -d production_db

# Check from container
docker exec -it go-microservice /app/server --db-enabled --db-host=production-db.example.com
```

**Check:**
- Database host is reachable
- Credentials are correct
- SSL mode matches database configuration
- Network policies allow connection

#### 3. Readiness Probe Failing

**Debug readiness issues:**
```bash
# Check readiness endpoint
curl -v http://localhost:8080/healthz/readiness

# Check database connectivity
kubectl exec -it -n go-microservice deployment/go-microservice -- /app/server migrate version
```

**Common causes:**
- Database not ready
- Database credentials incorrect
- Network connectivity issues
- Application initialization taking too long

#### 4. Migration Failures

**Check migration status:**
```bash
# View current version
./bin/server migrate version

# View migration logs
kubectl logs -n go-microservice -l job-name=db-migration

# Force version if needed (use with caution)
./bin/server migrate force <version>
```

#### 5. Memory/CPU Issues

**Monitor resource usage:**
```bash
# Docker
docker stats go-microservice

# Kubernetes
kubectl top pods -n go-microservice
```

**Adjust resources:**
- Increase memory/CPU limits
- Add more replicas
- Enable autoscaling
- Optimize application code

### Debug Mode

Enable debug logging:

```bash
# Set log level to debug
export APP_LOG_LEVEL=debug

# Or via CLI
./bin/server --log-level=debug
```

### Support Checklist

When reporting issues, provide:

1. Application logs
2. Environment configuration
3. Database connection status
4. Resource usage metrics
5. Recent changes or deployments
6. Error messages and stack traces

## Rollback Procedures

### Docker Rollback

```bash
# Stop current version
docker stop go-microservice

# Start previous version
docker run -d --name go-microservice your-registry/go-lang-template:v0.9.0
```

### Kubernetes Rollback

```bash
# View rollout history
kubectl rollout history deployment/go-microservice -n go-microservice

# Rollback to previous version
kubectl rollout undo deployment/go-microservice -n go-microservice

# Rollback to specific revision
kubectl rollout undo deployment/go-microservice -n go-microservice --to-revision=2
```

### Database Migration Rollback

```bash
# Rollback one migration
./bin/server migrate down

# Rollback N migrations
./bin/server migrate steps -N
```

## Maintenance

### Regular Tasks

1. **Weekly:**
   - Review application logs for errors
   - Check resource usage trends
   - Verify backup integrity

2. **Monthly:**
   - Update dependencies
   - Review and rotate secrets
   - Analyze performance metrics
   - Database maintenance (VACUUM, ANALYZE)

3. **Quarterly:**
   - Security audit
   - Load testing
   - Disaster recovery drill
   - Review and update documentation

### Backup Procedures

**Database backups:**

```bash
# PostgreSQL backup
pg_dump -h production-db.example.com -U app_user -d production_db > backup_$(date +%Y%m%d).sql

# Automated backups in Kubernetes
kubectl create cronjob db-backup --image=postgres:15-alpine \
  --schedule="0 2 * * *" -- pg_dump -h postgres-service -U app_user production_db > /backups/backup.sql
```

### Upgrade Procedures

1. Test in staging environment
2. Back up production database
3. Deploy new version
4. Run database migrations
5. Verify health checks
6. Monitor for issues
7. Rollback if needed

## Conclusion

This deployment guide covers the essential steps for deploying the Go microservice template to production. Always test thoroughly in staging environments before deploying to production, and maintain proper monitoring and backup procedures.

For additional help:
- Review application logs
- Check [CONFIGURATION.md](CONFIGURATION.md) for configuration details
- See [README.md](../README.md) for development information
- Open an issue on GitHub for support
