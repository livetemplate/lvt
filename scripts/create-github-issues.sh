#!/bin/bash
#
# Create GitHub Issues from Implementation Plan
#
# This script creates GitHub issues and milestones from the github-issues.json file.
#
# Prerequisites:
#   - gh CLI installed and authenticated (gh auth login)
#   - jq installed for JSON parsing
#
# Usage:
#   ./create-github-issues.sh                    # Create all issues
#   ./create-github-issues.sh --milestones-only  # Create only milestones
#   ./create-github-issues.sh --dry-run          # Show what would be created
#   ./create-github-issues.sh --milestone 1      # Create issues for milestone 1 only
#
# Note: Run from the lvt repo root directory

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_FILE="${SCRIPT_DIR}/github-issues.json"
REPO="livetemplate/lvt"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Parse arguments
DRY_RUN=false
MILESTONES_ONLY=false
SPECIFIC_MILESTONE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --milestones-only)
            MILESTONES_ONLY=true
            shift
            ;;
        --milestone)
            SPECIFIC_MILESTONE="$2"
            shift 2
            ;;
        --repo)
            REPO="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --dry-run           Show what would be created without actually creating"
            echo "  --milestones-only   Create only milestones, not issues"
            echo "  --milestone N       Create issues for milestone N only (1-6)"
            echo "  --repo OWNER/REPO   Target repository (default: livetemplate/lvt)"
            echo "  --help              Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Check prerequisites
if ! command -v gh &> /dev/null; then
    echo -e "${RED}Error: gh CLI is not installed${NC}"
    echo "Install it from: https://cli.github.com/"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is not installed${NC}"
    echo "Install it with: brew install jq (macOS) or apt install jq (Linux)"
    exit 1
fi

if ! gh auth status &> /dev/null; then
    echo -e "${RED}Error: gh CLI is not authenticated${NC}"
    echo "Run: gh auth login"
    exit 1
fi

if [[ ! -f "$JSON_FILE" ]]; then
    echo -e "${RED}Error: $JSON_FILE not found${NC}"
    exit 1
fi

echo -e "${BLUE}=== LVT GitHub Issues Creator ===${NC}"
echo -e "Repository: ${GREEN}${REPO}${NC}"
echo -e "JSON file: ${GREEN}${JSON_FILE}${NC}"
if [[ "$DRY_RUN" == "true" ]]; then
    echo -e "${YELLOW}DRY RUN MODE - no issues will be created${NC}"
fi
echo ""

# Create a mapping of issue IDs to GitHub issue numbers
# Using a temp file for bash 3.2 compatibility (macOS default)
ISSUE_MAP_FILE=$(mktemp)
trap "rm -f $ISSUE_MAP_FILE" EXIT

# Helper to set mapping: set_issue_map <id> <github_num>
set_issue_map() {
    echo "$1=$2" >> "$ISSUE_MAP_FILE"
}

# Helper to get mapping: get_issue_map <id>
get_issue_map() {
    grep "^$1=" "$ISSUE_MAP_FILE" 2>/dev/null | cut -d= -f2 | tail -1
}

# Function to create labels if they don't exist
create_labels() {
    echo -e "${BLUE}Creating labels...${NC}"

    local labels=(
        "testing|E2E and unit testing|0E8A16"
        "quick-win|Can be done quickly|FBCA04"
        "priority:critical|Critical priority|B60205"
        "priority:high|High priority|D93F0B"
        "priority:medium|Medium priority|FBCA04"
        "priority:low|Low priority|0E8A16"
        "bug|Something isn't working|D73A4A"
        "templates|Template-related changes|C5DEF5"
        "validation|Code validation|1D76DB"
        "generator|Code generator|5319E7"
        "infrastructure|Infrastructure changes|006B75"
        "mcp|MCP protocol related|0052CC"
        "evolution|Self-improvement system|9B59B6"
        "telemetry|Usage tracking|F9D0C4"
        "knowledge-base|Pattern knowledge base|FEF2C0"
        "cli|Command line interface|7057FF"
        "upstream|Upstream library changes|BFDADC"
        "components|UI components|C2E0C6"
        "architecture|Architectural changes|E99695"
        "monorepo|Monorepo structure|D4C5F9"
        "modal|Modal component|EDEDED"
        "toast|Toast component|EDEDED"
        "dropdown|Dropdown component|EDEDED"
        "styles|CSS/styling system|F9D0C4"
        "tailwind|Tailwind CSS|38BDF8"
        "unstyled|Unstyled/semantic CSS|EDEDED"
        "ci|Continuous integration|0E8A16"
        "kits|Kit templates|BFD4F2"
    )

    for label_info in "${labels[@]}"; do
        IFS='|' read -r name description color <<< "$label_info"
        if [[ "$DRY_RUN" == "true" ]]; then
            echo "  Would create label: $name"
        else
            gh label create "$name" --description "$description" --color "$color" --repo "$REPO" 2>/dev/null || true
        fi
    done

    echo -e "${GREEN}Labels ready${NC}"
    echo ""
}

# Function to create milestones
create_milestones() {
    echo -e "${BLUE}Creating milestones...${NC}"

    local milestone_count
    milestone_count=$(jq '.milestones | length' "$JSON_FILE")

    for ((i=0; i<milestone_count; i++)); do
        local title description
        title=$(jq -r ".milestones[$i].title" "$JSON_FILE")
        description=$(jq -r ".milestones[$i].description" "$JSON_FILE")

        if [[ "$DRY_RUN" == "true" ]]; then
            echo "  Would create milestone: $title"
        else
            # Check if milestone exists
            if gh api "repos/${REPO}/milestones" --jq ".[].title" 2>/dev/null | grep -q "^${title}$"; then
                echo -e "  ${YELLOW}Milestone exists: $title${NC}"
            else
                gh api "repos/${REPO}/milestones" \
                    --method POST \
                    --field title="$title" \
                    --field description="$description" \
                    --field state="open" > /dev/null
                echo -e "  ${GREEN}Created milestone: $title${NC}"
            fi
        fi
    done

    echo -e "${GREEN}Milestones ready${NC}"
    echo ""
}

# Function to get milestone number by title
get_milestone_number() {
    local title="$1"
    gh api "repos/${REPO}/milestones" --jq ".[] | select(.title == \"$title\") | .number" 2>/dev/null
}

# Function to create a single issue
create_issue() {
    local issue_index="$1"

    # Query JSON file directly (avoids bash variable escaping issues)
    local id title milestone body
    id=$(jq -r ".issues[$issue_index].id" "$JSON_FILE")
    title=$(jq -r ".issues[$issue_index].title" "$JSON_FILE")
    milestone=$(jq -r ".issues[$issue_index].milestone" "$JSON_FILE")
    body=$(jq -r ".issues[$issue_index].body" "$JSON_FILE")

    # Get labels as comma-separated string
    local labels
    labels=$(jq -r ".issues[$issue_index].labels | join(\",\")" "$JSON_FILE")

    # Build the title with issue ID prefix
    local full_title="[${id}] ${title}"

    if [[ "$DRY_RUN" == "true" ]]; then
        echo -e "  ${YELLOW}Would create:${NC} $full_title"
        echo "    Milestone: $milestone"
        echo "    Labels: $labels"
        return 0
    fi

    # Check if issue already exists (by title)
    local existing
    existing=$(gh issue list --repo "$REPO" --search "in:title [${id}]" --json number --jq '.[0].number' 2>/dev/null || echo "")

    if [[ -n "$existing" ]]; then
        echo -e "  ${YELLOW}Issue exists:${NC} #${existing} - ${full_title}"
        set_issue_map "$id" "$existing"
        return 0
    fi

    # Use a temp file for the body to handle special characters
    local body_file
    body_file=$(mktemp)
    echo "$body" > "$body_file"

    # Build gh issue create command
    local issue_url
    if [[ -n "$milestone" && "$milestone" != "null" ]]; then
        issue_url=$(gh issue create \
            --repo "$REPO" \
            --title "$full_title" \
            --body-file "$body_file" \
            ${labels:+--label "$labels"} \
            --milestone "$milestone" \
            2>&1)
    else
        issue_url=$(gh issue create \
            --repo "$REPO" \
            --title "$full_title" \
            --body-file "$body_file" \
            ${labels:+--label "$labels"} \
            2>&1)
    fi

    local issue_num
    issue_num=$(echo "$issue_url" | grep -oE '[0-9]+$' || echo "")

    rm -f "$body_file"

    if [[ -n "$issue_num" ]]; then
        set_issue_map "$id" "$issue_num"
        echo -e "  ${GREEN}Created:${NC} #${issue_num} - ${full_title}"
    else
        echo -e "  ${RED}Failed to create:${NC} ${full_title}"
        # Show error details if available
        if [[ -n "$issue_url" && "$issue_url" != *"github.com"* ]]; then
            echo -e "  ${RED}Error:${NC} $issue_url"
        fi
        return 1
    fi

    # Small delay to avoid rate limiting
    sleep 0.5
}

# Function to create all issues
create_issues() {
    echo -e "${BLUE}Creating issues...${NC}"

    local issue_count
    issue_count=$(jq '.issues | length' "$JSON_FILE")

    local created=0
    local skipped=0
    local failed=0

    for ((i=0; i<issue_count; i++)); do
        # Query JSON directly (avoid bash variable escaping issues)
        local id milestone_title
        id=$(jq -r ".issues[$i].id" "$JSON_FILE")
        milestone_title=$(jq -r ".issues[$i].milestone" "$JSON_FILE")

        # Filter by milestone if specified
        if [[ -n "$SPECIFIC_MILESTONE" ]]; then
            local milestone_num
            milestone_num=$(echo "$milestone_title" | grep -oE '^Milestone [0-9]+' | grep -oE '[0-9]+')
            if [[ "$milestone_num" != "$SPECIFIC_MILESTONE" ]]; then
                continue
            fi
        fi

        if create_issue "$i"; then
            ((created++))
        else
            ((failed++))
        fi
    done

    echo ""
    echo -e "${GREEN}Summary:${NC}"
    echo "  Created: $created"
    echo "  Failed: $failed"
}

# Function to update issue bodies with actual issue references
update_dependencies() {
    if [[ "$DRY_RUN" == "true" ]]; then
        echo -e "${YELLOW}Skipping dependency updates in dry run mode${NC}"
        return 0
    fi

    echo -e "${BLUE}Updating issue dependencies...${NC}"

    # For each issue, check if it has depends_on and update the body
    local issue_count
    issue_count=$(jq '.issues | length' "$JSON_FILE")

    for ((i=0; i<issue_count; i++)); do
        # Query JSON directly (avoid bash variable escaping issues)
        local id depends_on
        id=$(jq -r ".issues[$i].id" "$JSON_FILE")
        depends_on=$(jq -r ".issues[$i].depends_on // empty | @json" "$JSON_FILE")

        if [[ -z "$depends_on" || "$depends_on" == "null" || "$depends_on" == "[]" ]]; then
            continue
        fi

        local issue_num
        issue_num=$(get_issue_map "$id")
        if [[ -z "$issue_num" ]]; then
            continue
        fi

        # Build dependency links
        local dep_text="**Dependencies:**\n"
        local deps
        deps=$(jq -r ".issues[$i].depends_on[]" "$JSON_FILE")

        for dep in $deps; do
            if [[ "$dep" == milestone:* ]]; then
                dep_text+="- Milestone ${dep#milestone:} complete\n"
            else
                local dep_num
                dep_num=$(get_issue_map "$dep")
                if [[ -n "$dep_num" ]]; then
                    dep_text+="- #${dep_num}\n"
                else
                    dep_text+="- Issue ${dep}\n"
                fi
            fi
        done

        # This would need to update the issue body, which is complex
        # For now, just note that dependencies exist
        echo "  Issue #${issue_num} has dependencies"
    done
}

# Main execution
echo ""

# Create labels first
create_labels

# Create milestones
create_milestones

if [[ "$MILESTONES_ONLY" == "true" ]]; then
    echo -e "${GREEN}Done! Only milestones were created.${NC}"
    exit 0
fi

# Create issues
create_issues

echo ""
echo -e "${GREEN}=== Done! ===${NC}"

if [[ "$DRY_RUN" == "true" ]]; then
    echo ""
    echo "This was a dry run. To actually create the issues, run without --dry-run"
fi
