---
name: lvt:deploy
description: Use when deploying LiveTemplate applications to production - covers Docker containerization, Fly.io deployment, Kubernetes setup, database persistence, and production best practices
---

# lvt:deploy

Deploy LiveTemplate applications to production environments.

## Overview

LiveTemplate apps are standard Go binaries with SQLite databases, making them easy to deploy anywhere Go runs. This skill covers:

1. **Docker** - Containerize for any platform
2. **Fly.io** - Optimized for SQLite apps with built-in persistence
3. **Kubernetes** - For large-scale deployments
4. **Traditional VPS** - Simple binary deployment

## Prerequisites

**Before deployment:**
1. ✓ App builds successfully (`go build ./cmd/myapp`)
2. ✓ All migrations applied (`lvt migration status`)
3. ✓ Tests pass (`go test ./...`)
4. ✓ Production database prepared
5. ✓ Environment variables configured

## Docker Deployment

### 1. Create Dockerfile

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -o /myapp cmd/myapp/main.go

# Runtime stage
FROM alpine:latest

# Install SQLite (required for CGO)
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /myapp .

# Copy database directory (optional, for bundled data)
# COPY app.db .

# Copy migrations (needed if embedding migration runner)
COPY internal/database/migrations ./internal/database/migrations

# Copy template files (if not embedded in binary)
COPY internal/app ./internal/app

# Copy static files and client library
COPY static ./static
COPY client ./client

# Expose port
EXPOSE 8080

# Run binary
CMD ["./myapp"]
```

### 2. Create .dockerignore

```
# .dockerignore
app.db
*.log
.git
.env
tmp/
dist/
*.test
coverage.out
```

### 3. Build and Run

```bash
# Build image
docker build -t myapp:latest .

# Run container
docker run -p 8080:8080 \
  -v $(pwd)/data:/root/data \
  -e DATABASE_PATH=/root/data/app.db \
  myapp:latest

# Run with environment variables
docker run -p 8080:8080 \
  -v $(pwd)/data:/root/data \
  -e DATABASE_PATH=/root/data/app.db \
  -e PORT=8080 \
  -e ENV=production \
  myapp:latest
```

### 4. Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/root/data
    environment:
      - DATABASE_PATH=/root/data/app.db
      - PORT=8080
      - ENV=production
    restart: unless-stopped
```

```bash
# Start with compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## Fly.io Deployment

Fly.io is ideal for SQLite apps with built-in persistence and global distribution.

### 1. Install flyctl

```bash
# macOS
brew install flyctl

# Linux
curl -L https://fly.io/install.sh | sh

# Login
fly auth login
```

### 2. Create fly.toml

```toml
# fly.toml
app = "myapp"
primary_region = "sjc"

[build]
  builder = "paketobuildpacks/builder:base"

[env]
  PORT = "8080"

[[services]]
  http_checks = []
  internal_port = 8080
  processes = ["app"]
  protocol = "tcp"
  script_checks = []

  [services.concurrency]
    hard_limit = 25
    soft_limit = 20
    type = "connections"

  [[services.ports]]
    force_https = true
    handlers = ["http"]
    port = 80

  [[services.ports]]
    handlers = ["tls", "http"]
    port = 443

  [[services.tcp_checks]]
    grace_period = "1s"
    interval = "15s"
    restart_limit = 0
    timeout = "2s"

# SQLite persistence
[mounts]
  source = "myapp_data"
  destination = "/data"
```

### 3. Deploy to Fly.io

```bash
# Initialize app (first time only)
fly launch

# Create volume for SQLite database
fly volumes create myapp_data --region sjc --size 1

# Deploy
fly deploy

# Check status
fly status

# View logs
fly logs

# Open app
fly open

# SSH into instance
fly ssh console
```

### 4. Run Migrations on Fly.io

```bash
# Option 1: Run migrations before first deploy (in dev/CI)
lvt migration up
fly deploy

# Option 2: SSH and run goose
fly ssh console
goose -dir internal/database/migrations sqlite3 /data/app.db up

# Option 3: Auto-run migrations on deploy (add goose to Dockerfile)
# Add to fly.toml:
[deploy]
  release_command = "goose -dir internal/database/migrations sqlite3 /data/app.db up"
```

### 5. Scale on Fly.io

```bash
# Scale to multiple regions
fly regions add iad lhr syd

# Scale VM size
fly scale vm shared-cpu-1x --memory 512

# Scale instance count
fly scale count 2

# Auto-scale
fly autoscale set min=1 max=5
```

## Kubernetes Deployment

For large-scale production deployments.

### 1. Create Kubernetes Manifests

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: DATABASE_PATH
          value: "/data/app.db"
        volumeMounts:
        - name: data
          mountPath: /data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: myapp-pvc
```

**service.yaml:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  selector:
    app: myapp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

**pvc.yaml (for SQLite persistence):**
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: myapp-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```

### 2. Deploy to Kubernetes

```bash
# Apply manifests
kubectl apply -f pvc.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Check status
kubectl get pods
kubectl get services

# View logs
kubectl logs -f deployment/myapp

# Scale replicas
kubectl scale deployment myapp --replicas=5

# Update image
kubectl set image deployment/myapp myapp=myapp:v2
```

### 3. Important: SQLite + Kubernetes

**Warning:** SQLite doesn't support concurrent writes across multiple pods. For K8s:

**Option A: Single replica (simple)**
```yaml
spec:
  replicas: 1  # Only one pod writes to SQLite
```

**Option B: Read replicas (advanced)**
```yaml
# Use one writer pod + multiple read-only replicas
# Requires application-level read/write splitting
```

**Option C: Switch to PostgreSQL**
```bash
# For true horizontal scaling, migrate to PostgreSQL
# LiveTemplate's sqlc-based architecture makes this easier
```

## Traditional VPS Deployment

Simple deployment to DigitalOcean, Linode, AWS EC2, etc.

### 1. Build Binary

```bash
# Build for Linux (if developing on macOS)
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
  go build -o myapp cmd/myapp/main.go

# Or build on the server itself
ssh user@server
cd /opt/myapp
go build -o myapp cmd/myapp/main.go
```

### 2. Create systemd Service

```ini
# /etc/systemd/system/myapp.service
[Unit]
Description=MyApp LiveTemplate Application
After=network.target

[Service]
Type=simple
User=myapp
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/myapp
Restart=on-failure
RestartSec=5s

Environment="PORT=8080"
Environment="DATABASE_PATH=/var/lib/myapp/app.db"
Environment="ENV=production"

[Install]
WantedBy=multi-user.target
```

### 3. Deploy Steps

```bash
# 1. Copy files to server
scp -r . user@server:/opt/myapp/

# 2. SSH to server
ssh user@server

# 3. Setup
sudo useradd -r -s /bin/false myapp
sudo mkdir -p /var/lib/myapp
sudo chown myapp:myapp /var/lib/myapp
cd /opt/myapp

# 4. Install dependencies & build
go mod download
go build -o myapp cmd/myapp/main.go

# 5. Run migrations
# Option A: If lvt CLI is available
lvt migration up

# Option B: Use goose directly
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir internal/database/migrations sqlite3 /var/lib/myapp/app.db up

# 6. Enable systemd service
sudo systemctl enable myapp
sudo systemctl start myapp
sudo systemctl status myapp

# 7. View logs
sudo journalctl -u myapp -f
```

### 4. Nginx Reverse Proxy

```nginx
# /etc/nginx/sites-available/myapp
server {
    listen 80;
    server_name myapp.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/myapp /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# Add SSL with certbot
sudo certbot --nginx -d myapp.com
```

## Production Considerations

### 1. Database Backups

**SQLite backups:**
```bash
# Automated backup script
#!/bin/bash
# /opt/myapp/backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/var/backups/myapp"
DB_PATH="/var/lib/myapp/app.db"

mkdir -p $BACKUP_DIR

# Backup with SQLite's backup command
sqlite3 $DB_PATH ".backup '$BACKUP_DIR/app_$DATE.db'"

# Keep only last 7 days
find $BACKUP_DIR -name "app_*.db" -mtime +7 -delete

echo "Backup completed: app_$DATE.db"
```

**Cron job:**
```bash
# Run backup daily at 2 AM
0 2 * * * /opt/myapp/backup.sh
```

**Fly.io backup:**
```bash
# Snapshot volume
fly volumes snapshots create myapp_data

# List snapshots
fly volumes snapshots list myapp_data

# Restore from snapshot
fly volumes restore myapp_data <snapshot-id>
```

### 2. Environment Variables

**Create .env file (never commit to git):**
```bash
# .env
PORT=8080
DATABASE_PATH=/var/lib/myapp/app.db
ENV=production
SECRET_KEY=your-secret-key-here
```

**Load in app:**
```go
// cmd/myapp/main.go
package main

import (
    "os"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env in development
    if os.Getenv("ENV") != "production" {
        godotenv.Load()
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // ... rest of main
}
```

### 3. Migrations in Production

**IMPORTANT:** LiveTemplate migrations are run with `lvt migration up` (development) or `goose` (production).

**Before deployment - Run migrations in development:**
```bash
# Best practice: Run migrations before building
cd /path/to/app
lvt migration status
lvt migration up
go build ./cmd/myapp
```

**Option A: Include goose in Docker image:**
```dockerfile
# Add to Dockerfile before CMD
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations on container start
CMD goose -dir internal/database/migrations sqlite3 /data/app.db up && ./myapp
```

**Option B: Run migrations manually before deploy:**
```bash
# Docker
docker run --rm \
  -v $(pwd)/data:/data \
  myapp:latest \
  goose -dir internal/database/migrations sqlite3 /data/app.db up

# Fly.io
fly ssh console
goose -dir internal/database/migrations sqlite3 /data/app.db up

# VPS
ssh user@server
cd /opt/myapp
lvt migration up  # If lvt CLI is available
# OR
goose -dir internal/database/migrations sqlite3 /var/lib/myapp/app.db up
sudo systemctl restart myapp
```

**Option C: Auto-migrate on deploy (advanced):**
```go
// cmd/myapp/main.go
// Requires embedding goose or running migrations programmatically
import (
    "database/sql"
    "github.com/pressly/goose/v3"
)

func main() {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal(err)
    }

    // Run migrations
    if err := goose.Up(db, "internal/database/migrations"); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }

    // Start server
    // ...
}
```

### 4. Monitoring and Logs

**Structured logging:**
```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("Server starting", "port", port)
logger.Error("Database error", "error", err)
```

**Health check endpoint:**
```go
// Add to your routes
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    // Check database
    if err := db.Ping(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "unhealthy",
            "error": err.Error(),
        })
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
    })
})
```

### 5. Performance Tuning

**SQLite optimizations:**
```go
// In database initialization
db.Exec("PRAGMA journal_mode=WAL")  // Better concurrency
db.Exec("PRAGMA synchronous=NORMAL") // Faster writes
db.Exec("PRAGMA cache_size=-64000")  // 64MB cache
db.Exec("PRAGMA temp_store=MEMORY")  // In-memory temp tables
```

### 6. Static Files and Templates

**LiveTemplate serves:**
- Client library: `livetemplate-client.browser.js`
- Template files: `internal/app/**/*.tmpl`
- Static assets: CSS, images, etc.

**Option A: Embed in binary (recommended):**
```go
// cmd/myapp/main.go
import "embed"

//go:embed internal/app
var templates embed.FS

//go:embed static
var static embed.FS

//go:embed client
var client embed.FS

func main() {
    // Use embedded filesystems
    // ...
}
```

**Option B: Copy files in Docker:**
```dockerfile
# Already shown in Dockerfile above
COPY internal/app ./internal/app
COPY static ./static
COPY client ./client
```

**Verify static files are served:**
```bash
curl http://localhost:8080/static/livetemplate-client.browser.js
# Should return JavaScript file, not 404
```

## Common Deployment Mistakes

### ❌ Missing CGO_ENABLED for SQLite

```bash
# WRONG - SQLite won't work
CGO_ENABLED=0 go build ./cmd/myapp

# CORRECT
CGO_ENABLED=1 go build ./cmd/myapp
```

**Why wrong:** SQLite requires CGO. Without it, database operations fail.

### ❌ Not Persisting Database

```dockerfile
# WRONG - database lost on container restart
FROM alpine
COPY myapp .
CMD ["./myapp"]
```

```dockerfile
# CORRECT - mount volume for persistence
FROM alpine
COPY myapp .
VOLUME ["/data"]
ENV DATABASE_PATH=/data/app.db
CMD ["./myapp"]
```

**Why wrong:** Without volume, database is lost when container stops.

### ❌ Forgetting Migrations in Production

```bash
# WRONG - deploy without migrating
git push production main
# App crashes: "table not found"
```

```bash
# CORRECT
ssh production
./myapp migrate up  # Run migrations first
sudo systemctl restart myapp
```

**Why wrong:** Production database is out of sync with code.

### ❌ Multiple SQLite Writers in K8s

```yaml
# WRONG - SQLite corruption
spec:
  replicas: 5  # 5 pods writing to same SQLite file!
```

```yaml
# CORRECT - single writer
spec:
  replicas: 1  # SQLite handles one writer at a time
```

**Why wrong:** SQLite doesn't support concurrent writes from multiple processes.

### ❌ Hardcoded Paths and Ports

```go
// WRONG
db, err := sql.Open("sqlite3", "/Users/me/myapp/app.db")
server.ListenAndServe(":8080", nil)
```

```go
// CORRECT
dbPath := os.Getenv("DATABASE_PATH")
if dbPath == "" {
    dbPath = "./app.db"
}
db, err := sql.Open("sqlite3", dbPath)

port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
server.ListenAndServe(":" + port, nil)
```

**Why wrong:** Environment-specific paths break in production.

### ❌ Cross-Compiling with CGO

```bash
# WRONG - cross-compilation with CGO doesn't work simply
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build ./cmd/myapp
# Error: C compiler not found or wrong architecture
```

```bash
# CORRECT - use Docker for cross-platform builds
docker build -t myapp:latest .

# OR - build on target platform
ssh server
go build ./cmd/myapp

# OR - use cross-compilation tools (advanced)
# Install cross-compiler: brew install FiloSottile/musl-cross/musl-cross
CC=x86_64-linux-musl-gcc GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build ./cmd/myapp
```

**Why wrong:** SQLite requires CGO, which needs platform-specific C compilers. Docker multi-stage builds handle this automatically.

## Deployment Checklist

**Before deploying:**
- [ ] All tests pass (`go test ./...`)
- [ ] Builds successfully (`go build ./cmd/myapp`)
- [ ] Migrations applied (`lvt migration status`)
- [ ] Environment variables configured
- [ ] Database backup strategy in place
- [ ] Health check endpoint added
- [ ] Logs configured (structured JSON)
- [ ] .env file created (not committed)
- [ ] Dockerfile tested locally
- [ ] Static files/templates embedded or copied
- [ ] Client library accessible
- [ ] CGO enabled for SQLite (`CGO_ENABLED=1`)
- [ ] Reverse proxy configured (if needed)

**After deploying:**
- [ ] App accessible at production URL
- [ ] Health check returns 200
- [ ] Database queries work
- [ ] WebSocket connections work
- [ ] Static files load correctly
- [ ] Logs are being captured
- [ ] Backups running on schedule
- [ ] Monitoring alerts configured

## Quick Reference

**I want to...** | **Best Option** | **Command**
---|---|---
Deploy quickly | Fly.io | `fly launch && fly deploy`
Containerize | Docker | `docker build -t myapp .`
Use existing VPS | systemd | `systemctl enable myapp`
Auto-scale | Fly.io or K8s | `fly autoscale set min=1 max=5`
Global distribution | Fly.io | `fly regions add iad lhr`
Test locally | Docker Compose | `docker-compose up`
Backup database | Cron + SQLite | `sqlite3 app.db .backup`
Run migrations | goose | `goose -dir internal/database/migrations sqlite3 app.db up`

## Recommended: Fly.io for SQLite Apps

For most LiveTemplate apps, **Fly.io** is the best choice:

✓ Built-in SQLite persistence (volumes)
✓ Global distribution (multi-region)
✓ Auto-scaling
✓ Zero-downtime deploys
✓ Built-in SSL
✓ Easy rollbacks
✓ Affordable ($0-5/month for small apps)

```bash
# Complete Fly.io deployment
fly launch
fly volumes create myapp_data --size 1
fly deploy
fly open
```

## Remember

✓ Enable CGO for SQLite builds (`CGO_ENABLED=1`)
✓ Use volumes/mounts for database persistence
✓ Run migrations with `lvt migration up` (dev) or `goose` (production)
✓ Configure environment variables (never hardcode)
✓ Set up automated backups
✓ Add health check endpoint
✓ Test Docker builds locally before deploying
✓ SQLite = single writer (use replicas=1 in K8s)
✓ Embed or copy static files/templates/client library
✓ Use Docker for cross-platform builds (CGO cross-compilation is complex)

✗ Don't deploy without testing build
✗ Don't forget database persistence
✗ Don't skip migrations in production
✗ Don't assume `./myapp migrate up` exists (use `lvt` or `goose`)
✗ Don't use multiple SQLite writers in K8s
✗ Don't hardcode paths or ports
✗ Don't commit .env files or secrets
✗ Don't deploy without backup strategy
✗ Don't try simple cross-compilation with CGO (use Docker)
