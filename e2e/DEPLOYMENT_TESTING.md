# Deployment Testing Guide

This guide explains how to use the deployment testing infrastructure for end-to-end testing of lvt applications.

## Overview

The deployment testing infrastructure provides:
- **Mock deployments**: Fast, offline testing with simulated Fly.io API
- **Real deployments**: Actual Fly.io deployments for confidence testing
- **Smoke tests**: Post-deployment verification (HTTP, WebSocket, templates)
- **Cleanup tracking**: Automatic resource cleanup to prevent cost leaks

## Quick Start

### Running Mock Tests

Mock tests don't require credentials and run quickly:

```bash
# Run all mock deployment tests
go test -v ./e2e -run TestDeployment.*Mock

# Run specific mock test
go test -v ./e2e -run TestMockDeploymentWorkflow

# Run in short mode (skips slow tests)
go test -v ./e2e -short
```

### Running Real Fly.io Tests

Real deployment tests require:
1. Fly.io account and API token
2. `flyctl` CLI installed
3. Environment variables set

```bash
# Install flyctl (macOS)
brew install flyctl

# Login to Fly.io
flyctl auth login

# Export API token
export FLY_API_TOKEN=$(flyctl auth token)

# Run real deployment tests
go test -v ./e2e -run TestRealFlyDeployment

# Run with explicit environment flag (for expensive tests)
RUN_FLY_DEPLOYMENT_TESTS=true go test -v ./e2e -run TestFlyDeploymentWithResources
```

## Architecture

### Core Components

1. **Credentials Management** (`testing/credentials.go`)
   - Loads credentials from environment
   - Validates provider requirements
   - Provides skip helpers for missing credentials

2. **App Name Generator** (`testing/naming.go`)
   - Generates unique Fly.io-compatible names
   - Format: `lvt-test-{random}-{timestamp}`
   - Validates against Fly.io constraints

3. **Deployment Harness** (`testing/deployment.go`)
   - Manages deployment lifecycle
   - Tracks cleanup functions
   - Provider abstraction (Fly.io, Docker, K8s)

4. **Smoke Tests** (`testing/smoke.go`)
   - HTTP root path verification
   - Health endpoint check
   - Static asset loading
   - WebSocket connection (browser-based)
   - Template rendering validation
   - Retry logic with exponential backoff

### Provider Implementations

#### Mock Provider (`testing/providers/mock.go`)

Fast, in-memory simulation of Fly.io API:
- Configurable delays (default: 2s deployment)
- State tracking for apps and volumes
- No actual network calls
- Perfect for CI/CD

```go
client := providers.NewMockFlyClient()
client.SimulateDelay = false // Instant execution
client.Launch(appName, "sjc")
client.Deploy(appName, appDir)
client.Destroy(appName)
```

#### Fly.io Provider (`testing/providers/fly.go`)

Real Fly.io deployments via flyctl CLI:
- Wraps all flyctl commands
- JSON parsing for structured data
- Volume management
- App lifecycle (launch, deploy, destroy)
- Wait for ready with timeout

```go
client := providers.NewFlyClient(apiToken, "personal")
client.Launch(appName, "sjc")
volumeID, _ := client.CreateVolume(appName, "sjc", 1)
client.Deploy(appName, appDir, "sjc")
client.WaitForAppReady(appName, 5*time.Minute)
url, _ := client.GetAppURL(appName)
client.Destroy(appName)
```

## Usage Examples

### Example 1: Basic Mock Deployment Test

```go
func TestMyFeature(t *testing.T) {
    // Create mock client
    client := providers.NewMockFlyClient()
    appName := lvttesting.GenerateTestAppName("test")

    // Deploy
    client.Launch(appName, "sjc")
    client.Deploy(appName, "/path/to/app")

    // Verify
    status, _ := client.Status(appName)
    if status.Status != "running" {
        t.Errorf("Expected running, got %s", status.Status)
    }

    // Cleanup
    client.Destroy(appName)
}
```

### Example 2: Full Deployment Harness

```go
func TestWithHarness(t *testing.T) {
    opts := &lvttesting.DeploymentOptions{
        Provider:  lvttesting.ProviderFly,
        Region:    "sjc",
        Kit:       "multi",
        Resources: []string{"posts title content"},
    }

    dt := lvttesting.SetupDeployment(t, opts)

    // Deploy (uses mock if no FLY_API_TOKEN, real otherwise)
    if err := dt.Deploy(); err != nil {
        t.Fatalf("Deploy failed: %v", err)
    }

    // Smoke tests
    suite, err := lvttesting.RunSmokeTests(dt.AppURL, nil)
    if err != nil || !suite.AllPassed() {
        t.Fatal("Smoke tests failed")
    }

    // Cleanup happens automatically via t.Cleanup()
}
```

### Example 3: Real Fly.io with Credentials

```go
func TestRealDeployment(t *testing.T) {
    // Skip if no credentials
    lvttesting.RequireFlyCredentials(t)

    // Check flyctl
    if err := providers.CheckFlyctlInstalled(); err != nil {
        t.Skip("flyctl not installed")
    }

    opts := &lvttesting.DeploymentOptions{
        Provider: lvttesting.ProviderFly,
        Region:   "sjc",
        Kit:      "simple",
    }

    dt := lvttesting.SetupDeployment(t, opts)

    if err := dt.Deploy(); err != nil {
        t.Fatalf("Deployment failed: %v", err)
    }

    t.Logf("Deployed at: %s", dt.AppURL)
}
```

## Environment Variables

### Required for Real Fly.io Tests

```bash
# Fly.io API token
export FLY_API_TOKEN="your-token-here"

# Get token from flyctl
export FLY_API_TOKEN=$(flyctl auth token)
```

### Optional

```bash
# AWS credentials (for Litestream S3 backup)
export AWS_ACCESS_KEY_ID="your-key"
export AWS_SECRET_ACCESS_KEY="your-secret"
export S3_BUCKET="your-bucket"
export S3_REGION="us-east-1"

# DigitalOcean (future support)
export DO_API_TOKEN="your-token"

# Enable expensive tests
export RUN_FLY_DEPLOYMENT_TESTS="true"
```

## Test Organization

```
e2e/
├── deployment_mock_test.go       # Mock infrastructure tests
├── deployment_fly_test.go        # Real Fly.io deployment tests
└── DEPLOYMENT_TESTING.md         # This guide

testing/
├── credentials.go                # Credential management
├── naming.go                     # App name generation
├── deployment.go                 # Deployment harness
├── smoke.go                      # Smoke test suite
└── providers/
    ├── mock.go                   # Mock Fly.io client
    └── fly.go                    # Real Fly.io client
```

## Smoke Tests

The smoke test suite runs 5 tests against deployed apps:

1. **HTTP Root Path**: Verifies `GET /` returns 200 OK
2. **Health Endpoint**: Checks `/health` returns OK/healthy
3. **Static Assets**: Verifies `/livetemplate-client.js` loads
4. **WebSocket**: Tests WebSocket connection (browser-based, optional)
5. **Template Rendering**: Checks for template errors in HTML

### Running Smoke Tests

```go
opts := lvttesting.DefaultSmokeTestOptions()
opts.Timeout = 5 * time.Minute
opts.RetryDelay = 2 * time.Second
opts.MaxRetries = 3
opts.SkipBrowser = true // Skip WebSocket test

suite, err := lvttesting.RunSmokeTests(appURL, opts)
if err != nil {
    log.Fatalf("Smoke tests failed: %v", err)
}

suite.PrintResults()
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Deployment Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install flyctl
        run: |
          curl -L https://fly.io/install.sh | sh
          echo "$HOME/.fly/bin" >> $GITHUB_PATH

      - name: Run Mock Tests
        run: go test -v ./e2e -run TestDeployment.*Mock

      - name: Run Real Fly.io Tests
        if: github.ref == 'refs/heads/main'
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
        run: go test -v ./e2e -run TestRealFlyDeployment
```

## Cleanup

Cleanup is automatic via `t.Cleanup()`:
- Apps destroyed
- Volumes deleted
- No manual cleanup needed
- Runs even if test fails

To verify cleanup:

```go
dt := lvttesting.SetupDeployment(t, opts)
dt.Deploy()

// Cleanup registered automatically
// Will run when test ends
```

## Troubleshooting

### flyctl not found

```bash
# Install flyctl
brew install flyctl  # macOS
curl -L https://fly.io/install.sh | sh  # Linux

# Verify
flyctl version
```

### Tests skip with "credentials not available"

```bash
# Check if token is set
echo $FLY_API_TOKEN

# Get token
flyctl auth login
export FLY_API_TOKEN=$(flyctl auth token)
```

### Deployment times out

- Increase timeout in smoke test options
- Check Fly.io status: https://status.flyio.net
- Verify region availability

### Cleanup fails

- Check Fly.io dashboard for orphaned apps
- Manually destroy: `flyctl apps destroy <app-name>`
- List apps: `flyctl apps list`

## Best Practices

1. **Use mock tests for CI**: Fast, free, no credentials needed
2. **Use real tests for releases**: Confidence before production
3. **Always use unique app names**: Prevents conflicts
4. **Set appropriate timeouts**: Mock: 30s, Real: 5min
5. **Skip browser tests in CI**: Requires headless Chrome
6. **Clean up manually if needed**: Check Fly.io dashboard

## Future Enhancements

- Docker local deployment support
- Kubernetes deployment testing
- DigitalOcean App Platform support
- Parallel deployment testing
- Cost tracking and limits
- Performance benchmarking
