# Testing Guide

## Test Organization

- **Unit tests**: `internal/`, `commands/` - Fast, isolated logic tests
- **E2E tests**: `e2e/` - Browser automation, Docker integration tests
- **Golden tests**: `golden_test.go` - Template output validation

## Running Tests

### Quick Feedback (30 seconds)
```bash
make test-fast
```

Runs only unit tests. Use during development for fast feedback.

### Before Commit (3-4 minutes)
```bash
make test-commit
```

Runs unit + e2e tests. Run before creating commits.

### Full Suite (5-6 minutes)
```bash
make test-all
```

Includes deployment tests. Run before PRs.

### Specific Test
```bash
GOWORK=off go test -v -run=TestTutorialE2E ./e2e
```

## Test Optimization Architecture

The test suite uses several optimizations to achieve <6 minute execution:

1. **Shared Docker Base Image**: All e2e tests build on `lvt-base:latest` containing Go, sqlc, and dependencies. Each test builds only app-specific code (~10s vs ~60s).

2. **Chrome Container Pool**: 4 Chrome containers started once in TestMain. Tests borrow from pool instead of starting fresh containers.

3. **Parallel Execution**: Tests run with `-p 4` (4 packages concurrently) and `t.Parallel()` within packages.

4. **Optimized Timeouts**: Local development uses 10s WebSocket and 20s browser timeouts vs 30s/60s in CI.

## Performance Profile

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Total time | 10-15 min | 3-6 min | 50%+ faster |
| Docker builds | 12 min | 2 min | 83% faster |
| Chrome startups | 2.5 min | 1 min (once) | 60% faster |
| Memory peak | 4 GB | 2.2 GB | 45% less |

## Environment Variables

- `WEBSOCKET_TIMEOUT`: Override WebSocket ready timeout (default: 10s local, 30s CI)
- `BROWSER_TIMEOUT`: Override browser operation timeout (default: 20s local, 60s CI)
- `CI=true`: Enables CI-specific timeouts

## Cleanup

```bash
make test-clean
```

Removes all test Docker containers and images.

## Test Targets

- `make test-fast` - Unit tests only (~30s)
- `make test-commit` - Validation before commit (~3-4min)
- `make test-all` - Full suite including deployment (~5-6min)
- `make test-e2e` - E2E tests only
- `make test-unit` - Unit tests only (not short mode)
- `make test-clean` - Clean up Docker resources

## Troubleshooting

### Tests timeout
The test suite has a 5-minute timeout for `test-commit` and 10-minute timeout for `test-all`. If tests are timing out:

1. Run specific test packages to identify slow tests
2. Check Docker container status: `docker ps`
3. Clean up lingering resources: `make test-clean`
4. Monitor system resources during test run

### E2E tests fail
E2E tests require Docker and may fail if:
- Docker is not running
- Ports are already in use
- Previous test containers are still running (run `make test-clean`)

### Memory issues
E2E tests use Chrome containers which can consume memory. If you see memory-related failures:
- Close other applications
- Run `make test-clean` to remove old containers
- Run tests in smaller batches
