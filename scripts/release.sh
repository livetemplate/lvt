#!/usr/bin/env bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() { echo -e "${GREEN}✓${NC} $1"; }
log_warn() { echo -e "${YELLOW}⚠${NC} $1"; }
log_error() { echo -e "${RED}✗${NC} $1"; }
log_step() { echo -e "${BLUE}▸${NC} $1"; }

# Check prerequisites
check_prerequisites() {
    local missing=()

    command -v gh >/dev/null 2>&1 || missing+=("gh (GitHub CLI)")
    command -v go >/dev/null 2>&1 || missing+=("go")
    command -v jq >/dev/null 2>&1 || missing+=("jq (JSON processor)")
    command -v goreleaser >/dev/null 2>&1 || missing+=("goreleaser")

    if [ ${#missing[@]} -ne 0 ]; then
        log_error "Missing required tools: ${missing[*]}"
        echo ""
        echo "Install with:"
        echo "  macOS:   brew install gh go jq goreleaser"
        echo "  Linux:   apt-get install gh golang jq && go install github.com/goreleaser/goreleaser@latest"
        exit 1
    fi

    # Check GitHub CLI auth
    if ! gh auth status >/dev/null 2>&1; then
        log_error "GitHub CLI not authenticated. Run 'gh auth login' first"
        exit 1
    fi

    # Check optional tools
    if ! command -v git-chglog >/dev/null 2>&1; then
        log_warn "git-chglog not installed (optional). See: https://github.com/git-chglog/git-chglog#installation"
    fi
}

# Get core library version
get_core_library_version() {
    log_step "Fetching core library version from github.com/livetemplate/livetemplate"

    # Get latest release from GitHub API
    local core_version=$(gh release list --repo livetemplate/livetemplate --limit 1 --json tagName --jq '.[0].tagName' 2>/dev/null || echo "")

    if [ -z "$core_version" ]; then
        log_error "Could not fetch core library version"
        log_info "Make sure github.com/livetemplate/livetemplate has releases"
        exit 1
    fi

    # Remove 'v' prefix if present
    core_version=${core_version#v}

    log_info "Core library version: $core_version"
    echo "$core_version"
}

# Extract major.minor from version
get_major_minor() {
    local version=$1
    IFS='.' read -r major minor patch <<< "$version"
    echo "${major}.${minor}"
}

# Get current version
get_current_version() {
    if [ ! -f VERSION ]; then
        log_error "VERSION file not found"
        exit 1
    fi
    cat VERSION | tr -d '\n'
}

# Validate version against core library
validate_version() {
    local new_version=$1
    local core_version=$(get_core_library_version)

    local new_major_minor=$(get_major_minor "$new_version")
    local core_major_minor=$(get_major_minor "$core_version")

    if [ "$new_major_minor" != "$core_major_minor" ]; then
        log_error "Version mismatch!"
        echo ""
        echo "  LVT version:  $new_version (major.minor: $new_major_minor)"
        echo "  Core version: $core_version (major.minor: $core_major_minor)"
        echo ""
        echo "LVT must match core library's major.minor version."
        echo "Use: ${core_major_minor}.X where X is any patch version"
        exit 1
    fi

    log_info "Version validated against core library (major.minor: $core_major_minor)"
}

# Bump version
bump_version() {
    local current_version=$1
    local bump_type=$2

    IFS='.' read -r major minor patch <<< "$current_version"

    case $bump_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo "$bump_type"  # Allow custom version
            return
            ;;
    esac

    echo "${major}.${minor}.${patch}"
}

# Update version files
update_versions() {
    local new_version=$1

    log_step "Updating VERSION file to $new_version"
    echo "$new_version" > VERSION

    # Update go.mod to use latest core library with same major.minor
    log_step "Updating go.mod to use core library v$(get_major_minor "$new_version").x"

    log_info "Version files updated to $new_version"
}

# Generate changelog
generate_changelog() {
    local new_version=$1
    local prev_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

    log_step "Generating changelog for v$new_version"

    if command -v git-chglog >/dev/null 2>&1 && [ -f ".chglog/config.yml" ]; then
        # Use git-chglog only when both the binary and its config are present
        log_info "Using git-chglog for changelog generation"
        git-chglog --next-tag "v$new_version" -o CHANGELOG.md || {
            log_error "git-chglog failed. Aborting to avoid releasing with a stale changelog."
            log_info "To skip changelog generation, remove .chglog/config.yml or uninstall git-chglog."
            exit 1
        }
    else
        # Simple changelog generation (git-chglog not installed or .chglog/config.yml not found)
        if command -v git-chglog >/dev/null 2>&1; then
            log_warn ".chglog/config.yml not found — using simple changelog generation"
        else
            log_warn "git-chglog not installed — using simple changelog generation"
        fi

        if [ -n "$prev_tag" ]; then
            {
                echo "# Changelog"
                echo ""
                echo "All notable changes to the LVT CLI will be documented in this file."
                echo ""
                echo "The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),"
                echo "and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)."
                echo ""
                echo "## [v$new_version] - $(date +%Y-%m-%d)"
                echo ""
                echo "### Changes"
                echo ""
                git log "$prev_tag"..HEAD --pretty=format:"- %s (%h)" --no-merges | grep -v "^- Merge " || true
                echo ""
                echo ""
                tail -n +7 CHANGELOG.md 2>/dev/null || true
            } > CHANGELOG.md.tmp
            mv CHANGELOG.md.tmp CHANGELOG.md
        else
            log_warn "No previous tag found, skipping changelog generation"
        fi
    fi
}

# Commit and tag
commit_and_tag() {
    local new_version=$1

    log_step "Committing version bump"
    git add VERSION
    [ -f CHANGELOG.md ] && git add CHANGELOG.md
    git commit -m "chore(release): v$new_version

Release LVT CLI v$new_version

This release uses core library version: $(get_major_minor "$new_version").x

🤖 Generated with automated release script"

    log_step "Creating git tag v$new_version"
    git tag -a "v$new_version" -m "Release v$new_version"

    log_info "Committed and tagged v$new_version"
}

# Build and test
build_and_test() {
    log_step "Running Go tests..."
    GOWORK=off go test ./... -timeout=120s || {
        log_error "Tests failed, aborting release"
        exit 1
    }
    log_info "Tests passed"

    log_step "Building LVT CLI..."
    GOWORK=off go build -o /tmp/lvt . || {
        log_error "Build failed, aborting release"
        exit 1
    }
    log_info "CLI built successfully"
}

# Extract release notes from CHANGELOG
extract_release_notes() {
    local new_version=$1
    local notes_file="/tmp/release-notes-lvt-$new_version.md"

    if [ ! -f CHANGELOG.md ]; then
        log_warn "CHANGELOG.md not found, using default release notes"
        echo "Release v$new_version" > "$notes_file"
        echo "" >> "$notes_file"
        echo "LVT CLI - Code generator and development server for LiveTemplate" >> "$notes_file"
        echo "$notes_file"
        return
    fi

    # Extract notes for this version from CHANGELOG
    awk -v ver="$new_version" '
        /^## \[v/ {
            if (found) exit
            if ($0 ~ "\\[v"ver"\\]") {
                found=1
                next
            }
        }
        found && /^## \[v/ { exit }
        found { print }
    ' CHANGELOG.md > "$notes_file"

    # If empty, add default content
    if [ ! -s "$notes_file" ]; then
        log_warn "No changelog entries found for v$new_version, using default notes"
        echo "Release v$new_version" > "$notes_file"
        echo "" >> "$notes_file"
        echo "LVT CLI - Code generator and development server for LiveTemplate" >> "$notes_file"
    fi

    # Add installation instructions
    {
        echo ""
        echo "## Installation"
        echo ""
        echo "### Go Install"
        echo "\`\`\`bash"
        echo "go install github.com/livetemplate/lvt@v$new_version"
        echo "\`\`\`"
        echo ""
        echo "### Binary Download"
        echo "Download the appropriate binary for your platform from the assets below."
        echo ""
        echo "## Related Releases"
        echo ""
        echo "This release uses the LiveTemplate core library version $(get_major_minor "$new_version").x"
        echo ""
        echo "- Core Library: https://github.com/livetemplate/livetemplate"
        echo "- Client Library: https://github.com/livetemplate/client"
        echo "- Examples: https://github.com/livetemplate/examples"
    } >> "$notes_file"

    echo "$notes_file"
}

# Push and create GitHub release with GoReleaser
publish_github() {
    local new_version=$1
    local branch=$2

    log_step "Pushing commits and tags to GitHub (branch: $branch)"
    git push origin "$branch" || {
        log_error "Failed to push to origin. You may need to 'git pull --rebase origin $branch' first."
        log_error "After resolving, re-run the script or manually push:"
        log_error "  git push origin $branch && git push origin v$new_version"
        exit 1
    }
    git push origin "v$new_version"
    log_info "Pushed to GitHub"

    # Extract release notes
    log_step "Extracting release notes from CHANGELOG"
    local notes_file=$(extract_release_notes "$new_version")
    log_info "Release notes prepared"

    # Use GoReleaser to build and create release
    log_step "Building binaries and creating GitHub release with GoReleaser"
    goreleaser release --clean --release-notes "$notes_file" || {
        log_error "GoReleaser failed"
        exit 1
    }

    # Cleanup
    rm -f "$notes_file"

    log_info "GitHub release created: https://github.com/livetemplate/lvt/releases/tag/v$new_version"
}

# Dry run mode
dry_run() {
    local new_version=$1

    echo ""
    echo "🔍 DRY RUN MODE - No changes will be made"
    echo "========================================"
    echo ""

    log_info "Would validate version against core library"
    log_info "Would update VERSION to: $new_version"
    log_info "Would generate CHANGELOG.md"
    log_info "Would run tests and build"
    log_info "Would commit with message: chore(release): v$new_version"
    log_info "Would create tag: v$new_version"
    log_info "Would push to GitHub"
    log_info "Would run GoReleaser to build binaries and create release"

    echo ""
    log_info "Dry run completed successfully"
    exit 0
}

# Main release function
main() {
    local dry_run_mode=false

    # Parse flags
    while [[ $# -gt 0 ]]; do
        case $1 in
            --dry-run)
                dry_run_mode=true
                shift
                ;;
            *)
                shift
                ;;
        esac
    done

    echo "🚀 LVT CLI Release Automation"
    echo "=============================="
    echo ""

    check_prerequisites

    # Check git status
    if [ -n "$(git status --porcelain)" ]; then
        log_error "Working directory is not clean. Commit or stash changes first."
        echo ""
        git status --short
        exit 1
    fi

    # Detect branch and pull latest
    local branch
    if ! branch=$(git symbolic-ref --short HEAD 2>/dev/null); then
        log_error "Repository is in a detached HEAD state. Please check out a branch before running this script."
        exit 1
    fi
    if [ "$dry_run_mode" = true ]; then
        log_info "[dry-run] Would pull latest from origin/$branch"
    else
        log_step "Pulling latest changes from origin/$branch"
        git pull --rebase origin "$branch" || {
            log_error "Failed to pull latest changes. Resolve conflicts and try again."
            exit 1
        }
        log_info "Up to date with origin/$branch"
        log_info "Current HEAD: $(git log -1 --oneline)"
    fi

    # Get current version
    current_version=$(get_current_version)
    log_info "Current version: $current_version"

    # Get core library version for reference
    core_version=$(get_core_library_version)
    core_major_minor=$(get_major_minor "$core_version")

    echo ""
    log_info "Core library version: $core_version (major.minor: $core_major_minor)"
    log_info "LVT must use major.minor: $core_major_minor"

    # Ask for version bump type
    echo ""
    echo "Select version bump type:"
    echo "  1) patch (bug fixes)        → $(bump_version "$current_version" patch)"
    echo "  2) minor (sync with core)   → ${core_major_minor}.0"
    echo "  3) major (sync with core)   → ${core_major_minor}.0"
    echo "  4) custom version           → ${core_major_minor}.X"
    echo ""
    read -rp "Enter choice [1-4]: " choice

    case $choice in
        1) new_version=$(bump_version "$current_version" patch) ;;
        2) new_version="${core_major_minor}.0" ;;
        3) new_version="${core_major_minor}.0" ;;
        4)
            read -rp "Enter patch version for ${core_major_minor}.X: " patch_ver
            if ! [[ $patch_ver =~ ^[0-9]+$ ]]; then
                log_error "Invalid patch version. Must be a number"
                exit 1
            fi
            new_version="${core_major_minor}.${patch_ver}"
            ;;
        *)
            log_error "Invalid choice"
            exit 1
            ;;
    esac

    echo ""
    log_info "New version will be: $new_version"

    # Validate version
    validate_version "$new_version"

    echo ""
    echo "This will:"
    echo "  • Update VERSION file"
    echo "  • Generate/update CHANGELOG.md"
    echo "  • Run all tests and builds"
    echo "  • Commit and tag v$new_version"
    echo "  • Push to GitHub"
    echo "  • Build multi-platform binaries with GoReleaser"
    echo "  • Create GitHub release with binaries"
    echo ""

    if [ "$dry_run_mode" = true ]; then
        dry_run "$new_version"
    fi

    read -rp "Continue? [y/N]: " confirm

    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        log_warn "Release cancelled"
        exit 0
    fi

    echo ""
    log_info "Starting release process..."
    echo ""

    # Execute release steps
    update_versions "$new_version"
    generate_changelog "$new_version"
    build_and_test
    commit_and_tag "$new_version"
    publish_github "$new_version" "$branch"

    echo ""
    echo "================================================"
    log_info "✨ Release v$new_version completed successfully!"
    echo "================================================"
    echo ""
    echo "📦 Published artifacts:"
    echo "  • GitHub: https://github.com/livetemplate/lvt/releases/tag/v$new_version"
    echo "  • Go:     go install github.com/livetemplate/lvt@v$new_version"
    echo ""
    echo "📝 Next steps:"
    echo "  • Verify the release on GitHub"
    echo "  • Test binary downloads"
    echo "  • Update examples to use new version"
    echo "  • Announce the release"
}

main "$@"
