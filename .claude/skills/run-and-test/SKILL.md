---
name: lvt-run-and-test
description: Use when running and testing LiveTemplate applications - covers lvt serve for development, running generated tests, debugging issues, and verifying app works
---

# lvt:run-and-test

Run and test LiveTemplate applications during development.

## Overview

LiveTemplate apps can be run in two ways:
1. **Development server** (`lvt serve`) - Hot reload, browser auto-open, debugging
2. **Direct execution** (`go run`) - Production-like, no development features

Generated resources and views include E2E tests that verify HTTP endpoints and WebSocket connections.

## Running Your App

### Quick Start

```bash
# From project root
lvt serve

# Browser opens automatically at http://localhost:8080
# Server watches for changes and reloads
```

### Serve Options

```bash
# Custom port (if 8080 is busy)
lvt serve --port 3000

# Don't auto-open browser
lvt serve --no-browser

# Disable live reload
lvt serve --no-reload

# Specify directory
lvt serve --dir /path/to/app

# Force specific mode
lvt serve --mode app    # app, component, or kit
```

## Development Workflow

```bash
# 1. Create and migrate
lvt new myapp
cd myapp
lvt gen resource products name price:float
lvt migration up

# 2. Run development server
lvt serve

# 3. Open browser → http://localhost:8080/products
# 4. Edit code → Server auto-reloads
# 5. Test changes
```

## Running Tests

### Generated Tests

Each resource/view gets a test file:
- `internal/app/products/products_test.go`
- `internal/app/dashboard/dashboard_test.go`

Tests verify:
- HTTP endpoint responds
- WebSocket connection works
- Basic functionality

### Run All Tests

```bash
# From project root
go test ./...

# Or specific package
go test ./internal/app/products

# Verbose output
go test -v ./internal/app/products

# Skip slow WebSocket tests
go test -short ./...
```

### Before Running Tests

**Prerequisites:**
1. ✓ Migrations applied (`lvt migration up`)
2. ✓ Dependencies installed (`go mod tidy`)
3. ✓ No server running on test ports

```bash
# Complete test setup
lvt migration up
go mod tidy
go test ./...
```

## Common Issues

### ❌ Port Already in Use

```bash
# Error: listen tcp :8080: bind: address already in use

# Solutions:
# 1. Use different port
lvt serve --port 3000

# 2. Kill existing process
lsof -ti:8080 | xargs kill -9

# 3. Find and stop lvt serve
ps aux | grep "lvt serve"
kill <PID>
```

### ❌ Tests Fail: Server Won't Start

```bash
# Error in test: Failed to start server
# Or: Failed to connect to WebSocket

# Common causes:
# 1. Haven't run migrations
lvt migration up
go test ./...

# 2. Port conflict
# Tests use random free ports, but might conflict with running server
# Stop lvt serve before testing:
pkill -f "lvt serve"
go test ./...

# 3. Missing dependencies
go mod tidy
go test ./...

# 4. WebSocket tests are flaky
# Use -short mode for reliable, fast tests during development:
go test -short ./...
```

### ❌ 404 Not Found

```bash
# Browser shows 404 at http://localhost:8080/products

# Check:
# 1. Did migration run?
lvt migration status

# 2. Is route registered?
cat cmd/myapp/main.go | grep products

# 3. Restart server
# Ctrl+C then lvt serve
```

### ❌ Database Locked

```bash
# Error: database is locked

# Another process has app.db open:
# 1. Stop all running instances
pkill -f "go run"
pkill -f "lvt serve"

# 2. Close database connections
lsof app.db

# 3. Restart
lvt serve
```

## Test Structure

Generated tests follow this pattern:

```go
func TestProductsWebSocket(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping WebSocket test in short mode")
    }

    // 1. Get random free port
    serverPort, _ := getFreePort()

    // 2. Start app server
    go run cmd/myapp/main.go

    // 3. Wait for ready
    // Polls HTTP endpoint

    // 4. Test HTTP
    GET /products → 200 OK

    // 5. Test WebSocket
    ws://localhost:PORT/products

    // 6. Cleanup
    Kill server process
}
```

**Key points:**
- Tests start their own server instance
- Use random ports to avoid conflicts
- Skip in `-short` mode for faster feedback
- Clean up server process automatically

## Running Modes

LiveTemplate auto-detects mode from project structure:

| Mode | When Used | What It Does |
|------|-----------|--------------|
| **app** | Standard apps with cmd/ | Runs your application |
| **component** | Component development | Component preview server |
| **kit** | Kit development | Kit testing environment |

Force specific mode:
```bash
lvt serve --mode app
```

## Debugging Tips

### Check Server is Running

```bash
# Is anything on port 8080?
lsof -i:8080

# Is lvt serve running?
ps aux | grep "lvt serve"

# Test endpoint directly
curl http://localhost:8080/products
```

### View Server Logs

```bash
# lvt serve shows logs in terminal
lvt serve
# Watch for:
# - "Server starting on http://localhost:8080"
# - Database errors
# - Route registration
# - Template parsing errors
```

### Database State

```bash
# Check migrations applied
lvt migration status

# Inspect database directly
sqlite3 app.db
> .tables
> .schema products
> SELECT * FROM products;
> .quit
```

### Test Specific Handler

```bash
# Run single test
go test -v ./internal/app/products -run TestProductsWebSocket

# With more detail
go test -v ./internal/app/products -run TestProductsWebSocket 2>&1 | less
```

## Production Build

Development server is NOT for production:

```bash
# Build binary
go build -o myapp cmd/myapp/main.go

# Run binary
./myapp

# Or directly
go run cmd/myapp/main.go
```

For production deployment, see lvt:deploy skill.

## Quick Reference

**I want to...** | **Command**
---|---
Run development server | `lvt serve`
Change port | `lvt serve --port 3000`
Run all tests | `go test ./...`
Run fast tests only | `go test -short ./...`
Test one package | `go test ./internal/app/products`
Stop server | Ctrl+C or `pkill -f "lvt serve"`
Check what's on port | `lsof -i:8080`
View database | `sqlite3 app.db`
Check migration status | `lvt migration status`
Full reset | `pkill -f "lvt serve" && lvt migration up && go mod tidy`

## Typical Development Session

```bash
# Morning: Start fresh
cd myapp
lvt migration up
go mod tidy
lvt serve

# Develop features...
# Browser auto-reloads on changes

# Before committing
go test ./...         # All tests pass?
go build ./cmd/myapp  # Builds clean?

# End of day
Ctrl+C  # Stop server
```

## Remember

✓ `lvt serve` for development (hot reload, auto-browser)
✓ Run migrations before testing (`lvt migration up`)
✓ Tests start their own server on random ports
✓ Use `--port` if 8080 is busy
✓ Stop server before running tests to avoid conflicts
✓ `go test -short` to skip slow WebSocket tests

✗ Don't use `lvt serve` in production
✗ Don't forget `go mod tidy` after generating resources
✗ Don't run tests while `lvt serve` is running on same port
✗ Don't manually open database while server is running
