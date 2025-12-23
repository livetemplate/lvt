#!/bin/bash
# Test lvt-plan skill behavior
# Usage: ./scripts/test-skill.sh

set -e

TEST_DIR="/tmp/lvt-skill-test-$$"
PROMPT="create a blog app with authentication using livetemplate"

echo "=== Testing lvt-plan skill ==="
echo "Test directory: $TEST_DIR"
echo "Prompt: $PROMPT"
echo ""

# Create test directory
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Install agent
echo "Installing agent..."
lvt install-agent --llm claude --force > /dev/null 2>&1

# Show what skill contains
echo "=== Skill content (first 30 lines) ==="
head -30 .claude/skills/plan/SKILL.md
echo ""

# Show settings instructions
echo "=== Settings instructions ==="
cat .claude/settings.json | python3 -c "import json,sys; d=json.load(sys.stdin); print(d.get('instructions','')[:800])"
echo ""

# Run Claude with the prompt and capture response
echo "=== Running Claude ==="
echo "Sending prompt: $PROMPT"
echo ""

# Use claude CLI with --print to get single response
# --no-input prevents interactive mode
RESPONSE=$(echo "$PROMPT" | timeout 60 claude --print --dangerously-skip-permissions 2>/dev/null || echo "TIMEOUT or ERROR")

echo "=== Response ==="
echo "$RESPONSE"
echo ""

# Check for failure patterns
echo "=== Analysis ==="
if echo "$RESPONSE" | grep -qi "question 1\|what would you like\|let me gather\|let me walk you through"; then
    echo "❌ FAIL: Response contains questionnaire pattern"
    echo ""
    echo "Found patterns:"
    echo "$RESPONSE" | grep -i "question\|what would you\|let me gather\|let me walk" | head -5
else
    echo "✅ Might be OK - no obvious questionnaire patterns"
fi

if echo "$RESPONSE" | grep -qi "Plan for your\|Ready to create\|Commands I'll run"; then
    echo "✅ PASS: Response contains plan table pattern"
else
    echo "❌ FAIL: Response missing plan table pattern"
fi

# Cleanup
cd /
rm -rf "$TEST_DIR"

echo ""
echo "=== Done ==="
