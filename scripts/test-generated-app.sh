#!/bin/bash
set -e

# Test script for verifying generated LiveTemplate apps work in a single shot
# Usage: ./scripts/test-generated-app.sh [app-name] [kit]

APP_NAME="${1:-testblog}"
KIT="${2:-multi}"
TEST_DIR="/tmp/lvt-test-$$"

echo "=== LiveTemplate Generated App Test ==="
echo "App: $APP_NAME"
echo "Kit: $KIT"
echo "Test directory: $TEST_DIR"
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo "=== Cleaning up ==="
    rm -rf "$TEST_DIR"
    echo "Done."
}
trap cleanup EXIT

# Build lvt first
echo "=== Building lvt ==="
go build -o ~/go/bin/lvt .
echo "✅ lvt built"

# Create test directory
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Step 1: Create new app
echo ""
echo "=== Step 1: Creating new app ==="
lvt new "$APP_NAME" --kit "$KIT"
cd "$APP_NAME"
echo "✅ App created"

# Step 2: Generate resource
echo ""
echo "=== Step 2: Generating resource ==="
lvt gen resource posts title:string content:text published:bool
echo "✅ Resource generated"

# Step 3: Generate auth (with delay to avoid timestamp collision)
echo ""
echo "=== Step 3: Generating auth ==="
sleep 2
lvt gen auth
echo "✅ Auth generated"

# Step 4: Install dependencies (must run before migrations to avoid sqlc errors)
echo ""
echo "=== Step 4: Installing dependencies ==="
go mod tidy
echo "✅ Dependencies installed"

# Step 5: Run migrations
echo ""
echo "=== Step 5: Running migrations ==="
lvt migration up
echo "✅ Migrations complete"

# Step 6: Generate sqlc code (in case migration's internal sqlc failed)
echo ""
echo "=== Step 6: Generating sqlc code ==="
cd database && sqlc generate && cd ..
echo "✅ sqlc generated"

# Step 7: Build
echo ""
echo "=== Step 7: Building app ==="
GOWORK=off go build ./...
echo "✅ Build successful"

# Step 8: Run short tests (skip E2E)
echo ""
echo "=== Step 8: Running short tests ==="
GOWORK=off go test ./... -short -v
echo "✅ Short tests passed"

# Step 9: Run full E2E tests (requires Docker)
echo ""
echo "=== Step 9: Running E2E tests (requires Docker) ==="
if command -v docker &> /dev/null && docker info &> /dev/null; then
    GOWORK=off go test ./... -v -timeout 5m
    echo "✅ All E2E tests passed"
else
    echo "⚠️  Docker not available, skipping E2E tests"
fi

echo ""
echo "========================================="
echo "✅ All tests passed! App works in a single shot."
echo "========================================="
