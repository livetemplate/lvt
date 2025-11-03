#!/bin/bash
#
# Setup local development with sibling core library
#
# This script configures the LVT repository to use a local checkout
# of the core livetemplate library instead of the published version.
#
# Usage:
#   ./scripts/setup-local-dev.sh          # Enable local development
#   ./scripts/setup-local-dev.sh --undo   # Revert to published version
#
# Requirements:
#   - Core library must be checked out at ../livetemplate/
#   - Both repos should be sibling directories under the same parent

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

CORE_PATH="../livetemplate"
MODULE_PATH="github.com/livetemplate/livetemplate"

# Check if we're undoing the setup
if [[ "$1" == "--undo" || "$1" == "undo" ]]; then
    echo "Reverting to published version..."
    go mod edit -dropreplace="$MODULE_PATH"
    go mod tidy
    echo -e "${GREEN}✓ Reverted to using published version of $MODULE_PATH${NC}"
    echo ""
    echo "To re-enable local development, run: ./scripts/setup-local-dev.sh"
    exit 0
fi

# Check if core library exists
if [ ! -d "$CORE_PATH" ]; then
    echo -e "${RED}✗ Core library not found at $CORE_PATH${NC}"
    echo ""
    echo "Expected directory structure:"
    echo "  parent/"
    echo "  ├── livetemplate/  (core library)"
    echo "  └── lvt/           (this repository)"
    echo ""
    echo "Please clone the core library:"
    echo "  git clone git@github.com:livetemplate/livetemplate.git $CORE_PATH"
    exit 1
fi

# Check if it's a valid Go module
if [ ! -f "$CORE_PATH/go.mod" ]; then
    echo -e "${RED}✗ $CORE_PATH does not appear to be a valid Go module${NC}"
    echo ""
    echo "Make sure you cloned the correct repository."
    exit 1
fi

# Verify it's the correct module
MODULE_DECL=$(head -n 1 "$CORE_PATH/go.mod" | awk '{print $2}')
if [ "$MODULE_DECL" != "$MODULE_PATH" ]; then
    echo -e "${YELLOW}⚠ Warning: Module path mismatch${NC}"
    echo "  Expected: $MODULE_PATH"
    echo "  Found:    $MODULE_DECL"
    echo ""
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo "Setting up local development..."
echo "  Core library: $CORE_PATH"
echo ""

# Add replace directive
go mod edit -replace="$MODULE_PATH=$CORE_PATH"
go mod tidy

echo -e "${GREEN}✓ Local development setup complete!${NC}"
echo ""
echo "The LVT repository now uses the local core library at:"
echo "  $CORE_PATH"
echo ""
echo "Changes you make to the core library will be immediately"
echo "reflected when building or testing LVT."
echo ""
echo "To revert to the published version, run:"
echo "  ./scripts/setup-local-dev.sh --undo"
