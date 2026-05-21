#!/bin/bash

# Pre-commit hook for LVT CLI
# Automatically formats code and runs validation

set -e

echo "🔄 Running pre-commit validation..."

# Step 0: Clean up generated test artifacts
if [ -d "commands/internal" ]; then
    echo "🧹 Cleaning up generated test artifacts..."
    rm -rf commands/internal/
fi

# All `go` invocations below run with GOWORK=off so the hook is portable
# across nested-worktree dev setups (e.g. .worktrees/<feature-branch>/),
# where the workspace's go.work doesn't include the worktree's local module
# and `go list ./...` otherwise reports "directory prefix . does not contain
# modules listed in go.work or their selected dependencies". Intentional
# policy: the hook always tests the module in isolation, ignoring any
# sibling-module replace directives a contributor may have configured in
# their workspace go.work. From the main checkout this is a no-op; from a
# worktree it's the gate that makes the hook self-contained.

# Step 1: Auto-format Go code before validation
echo "📝 Auto-formatting Go code..."
if GOWORK=off go fmt ./...; then
    echo "✅ Code formatting completed"

    # Add any formatted files to the commit
    FORMATTED_FILES=$(git diff --name-only)
    if [ -n "$FORMATTED_FILES" ]; then
        echo "📁 Adding formatted files to commit:"
        echo "$FORMATTED_FILES"
        git add $FORMATTED_FILES
    fi
else
    echo "❌ Code formatting failed"
    exit 1
fi

# Step 2: Run golangci-lint (if available)
# Requires golangci-lint >= 2.0 (`--default=none` replaced v1's `--disable-all`).
# If blocked here on v1, upgrade: https://golangci-lint.run/welcome/install/
# `gosimple` is intentionally omitted: in v2 its S1xxx checks live under
# `staticcheck`, so they still run; keeping `gosimple` enabled triggers a
# "linter not found" warning.
if command -v golangci-lint >/dev/null 2>&1; then
    echo "🔍 Running golangci-lint..."
    if GOWORK=off golangci-lint run --default=none --enable=errcheck,unused,staticcheck,ineffassign; then
        echo "✅ Linting passed"
    else
        echo "❌ Linting failed - commit blocked"
        echo "💡 Fix linting errors before committing"
        exit 1
    fi
else
    echo "⚠️  golangci-lint not installed (skipping)"
fi

# Step 3: Run all Go tests with timeout
echo "🧪 Running Go tests..."
if GOWORK=off go test -v ./... -timeout=120s; then
    echo "✅ All tests passed"
else
    echo "❌ Tests failed - commit blocked"
    exit 1
fi

echo "✅ Pre-commit validation completed successfully"
