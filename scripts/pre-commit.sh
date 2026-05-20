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

# Step 1: Auto-format Go code before validation
echo "📝 Auto-formatting Go code..."
if go fmt ./...; then
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
if command -v golangci-lint >/dev/null 2>&1; then
    echo "🔍 Running golangci-lint..."
    if golangci-lint run --default=none --enable=errcheck,unused,staticcheck,ineffassign; then
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
if go test -v ./... -timeout=120s; then
    echo "✅ All tests passed"
else
    echo "❌ Tests failed - commit blocked"
    exit 1
fi

echo "✅ Pre-commit validation completed successfully"
