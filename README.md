# LiveTemplate CLI Generator (`lvt`)

A Phoenix-inspired code generator for LiveTemplate applications with CRUD functionality and interactive TUI wizards.

## Installation

```bash
go install github.com/livetemplate/lvt@latest
```

Or download pre-built binaries from the [releases page](https://github.com/livetemplate/lvt/releases).

Or build from source:

```bash
git clone https://github.com/livetemplate/lvt
cd lvt
go build -o lvt .
```

## Related Projects

- **[LiveTemplate Core](https://github.com/livetemplate/livetemplate)** - Go library for server-side rendering
- **[Client Library](https://github.com/livetemplate/client)** - TypeScript client for browsers
- **[Examples](https://github.com/livetemplate/examples)** - Example applications

## Version Synchronization

LVT follows the LiveTemplate core library's major.minor version:
- Core: `v0.1.5` → LVT: `v0.1.x` (any patch version)
- Core: `v0.2.0` → LVT: `v0.2.0` (must match major.minor)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

All pull requests require passing CI checks including tests, linting, and code formatting.

## Support

- **Issues**: [GitHub Issues](https://github.com/livetemplate/lvt/issues)
- **Discussions**: [GitHub Discussions](https://github.com/livetemplate/lvt/discussions)

## AI-Assisted Development

`lvt` provides AI assistance through multiple approaches, supporting all major AI assistants:

### Quick Start

```bash
# List available AI agents
lvt install-agent --list

# Install agent for your AI assistant
lvt install-agent --llm <type>    # claude, copilot, cursor, aider, generic

```

### Supported AI Assistants

| AI Assistant | Installation | Best For |
|--------------|-------------|----------|
| **Claude Code** | `lvt install-agent --llm claude` | Full workflows with 20+ skills |
| **GitHub Copilot** | `lvt install-agent --llm copilot` | In-editor suggestions |
| **Cursor** | `lvt install-agent --llm cursor` | Rule-based development |
| **Aider** | `lvt install-agent --llm aider` | CLI-driven development |
| **Generic/Other** | `lvt install-agent --llm generic` | Custom LLMs via CLI |

### Claude Code (Recommended)

Full-featured agent with skills and workflows:

```bash
# Install
lvt install-agent --llm claude

# Upgrade
lvt install-agent --upgrade

# Start Claude Code
claude
```

**Features:**
- 20+ skills for lvt commands
- Project management agent
- Guided workflows
- Best practices enforcement

**Try asking:**
- "Add a posts resource with title and content"
- "Generate authentication system"
- "Create a quickstart blog app"

### GitHub Copilot

Instructions-based integration:

```bash
# Install
lvt install-agent --llm copilot

# Open in VS Code - Copilot automatically understands LiveTemplate
```

### Cursor

Rule-based development patterns:

```bash
# Install
lvt install-agent --llm cursor

# Open project in Cursor - rules apply to *.go files automatically
```

### Aider

CLI configuration:

```bash
# Install
lvt install-agent --llm aider

# Start Aider - configuration loads automatically
aider
```

### Upgrading Agents

```bash
# Upgrade any agent type
lvt install-agent --llm <type> --upgrade
```

This preserves your custom settings while updating the agent files.

### Complete Setup Guide

For detailed setup instructions for each AI assistant, see:

- **[docs/AGENT_SETUP.md](docs/AGENT_SETUP.md)** - Complete setup guide for all AI assistants
- **[docs/WORKFLOWS.md](docs/WORKFLOWS.md)** - Common development workflows
- **[docs/AGENT_USAGE_GUIDE.md](docs/AGENT_USAGE_GUIDE.md)** - Claude Code usage examples

## Quick Start

You can use `lvt` in two modes: **Interactive** (TUI wizards) or **Direct** (CLI arguments).

**Important:** Create apps **outside** of existing Go module directories. If you create an app inside another Go module (e.g., for testing), you'll need to use `GOWORK=off` when running commands:
```bash
GOWORK=off go run cmd/myapp/main.go
```

### Interactive Mode (Recommended for New Users)

```bash
# Launch interactive app creator
lvt new

# Launch interactive resource builder
lvt gen

# Launch interactive view creator
lvt gen view
```

### Direct Mode

### 1. Create a New App

```bash
lvt new myapp
cd myapp
```

This generates:
- Complete Go project structure
- Database layer with sqlc integration
- go.mod with Go 1.24+ tools directive
- README with next steps

### 2. Generate a CRUD Resource

```bash
# With explicit types
lvt gen users name:string email:string age:int

# With inferred types (NEW!)
lvt gen products name price quantity enabled created_at
# → Infers: name:string price:float quantity:int enabled:bool created_at:time
```

This generates:
- `app/users/users.go` - Full CRUD handler
- `app/users/users.tmpl` - Tailwind CSS UI
- `app/users/users_ws_test.go` - WebSocket tests
- `app/users/users_test.go` - Chromedp E2E tests
- Database schema and queries (appended)

### 3. Run Migrations

```bash
lvt migration up  # Runs pending migrations and auto-generates database code
```

This automatically:
- Applies pending database migrations
- Runs `sqlc generate` to create Go database code
- Updates your query interfaces

### 4. Wire Up Routes

Add to `cmd/myapp/main.go`:

```go
import "myapp/app/users"

// In main():
http.Handle("/users", users.Handler(queries))
```

### 5. Run the App

```bash
go run cmd/myapp/main.go
```

Open http://localhost:8080/users

## Tutorial: Building a Blog System

Let's build a complete blog system with posts, comments, and categories to demonstrate lvt's capabilities.

### Step 1: Create the Blog App

```bash
lvt new myblog
cd myblog
```

This creates your project structure with database setup, main.go, and configuration. Dependencies are automatically installed via `go get ./...`.

### Step 2: Generate Resources

```bash
lvt gen posts title content:string published:bool
lvt gen categories name description
lvt gen comments post_id:references:posts author text
```

This generates for each resource:
- ✅ `app/{resource}/{resource}.go` - CRUD handler with LiveTemplate integration
- ✅ `app/{resource}/{resource}.tmpl` - Component-based template with Tailwind CSS
- ✅ `app/{resource}/{resource}_test.go` - E2E tests with chromedp
- ✅ Database migration file with unique timestamps
- ✅ SQL queries appended to `database/queries.sql`

For the `comments` resource with `post_id:references:posts`:
- ✅ Creates `post_id` field as TEXT (matching posts.id type)
- ✅ Adds foreign key constraint: `FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE`
- ✅ Creates index on `post_id` for query performance
- ✅ No manual migration needed!

### Step 3: Run Migrations

```bash
lvt migration up
```

This command:
- ✅ Runs all pending database migrations
- ✅ Automatically generates Go database code with sqlc
- ✅ Creates type-safe query interfaces

You'll see output like:
```
Running pending migrations...
OK   20240315120000_create_posts.sql
OK   20240315120001_create_categories.sql
OK   20240315120002_create_comments.sql
Generating database code with sqlc...
✅ Database code generated successfully!
✅ Migrations complete!
```

### Step 4: Resolve Dependencies

```bash
go mod tidy
```

This resolves all internal package imports created by the generated code. Required before running the app.

### Step 5: Wire Up Routes

The routes are auto-injected, but verify in `cmd/myblog/main.go`:

```go
import (
    "myblog/app/posts"
    "myblog/app/categories"
    "myblog/app/comments"
)

func main() {
    // ... database setup ...

    // Routes (auto-injected)
    http.Handle("/posts", posts.Handler(queries))
    http.Handle("/categories", categories.Handler(queries))
    http.Handle("/comments", comments.Handler(queries))

    // Start server
    http.ListenAndServe(":8080", nil)
}
```

### Step 6: Run the Blog

```bash
go run cmd/myblog/main.go
```

Visit:
- http://localhost:8080/posts - Create and manage blog posts
- http://localhost:8080/categories - Organize posts by category
- http://localhost:8080/comments - View all comments

**Note:** Visiting http://localhost:8080/ will show a 404 since no root handler exists. You can add a homepage next.

### Step 7: Add a Custom View for the Homepage (Optional)

```bash
lvt gen view home
go mod tidy
go run cmd/myblog/main.go
```

This creates a view-only handler (no database operations). Edit `app/home/home.tmpl` to create your landing page, then visit http://localhost:8080/home.

### Step 8: Test the Application

```bash
# Run all tests (E2E + WebSocket)
go test ./...

# Run specific resource tests
go test ./app/posts -v
```

### Customization Ideas

**1. Generate resources (CSS framework determined by kit):**
```bash
# Resources use the CSS framework from your chosen kit
# Multi and single kits use Tailwind CSS
# Simple kit uses no CSS framework (semantic HTML)

lvt gen tags name

# To use a different CSS framework, create your app with a different kit
lvt new myapp --kit simple  # Uses no CSS (semantic HTML)
cd myapp
lvt gen authors name bio    # Will use semantic HTML
```

**2. Use Type Inference:**
```bash
# Field types are inferred from names
lvt gen articles title content published_at author email price

# Infers: title=string, content=string, published_at=time,
#         author=string, email=string, price=float
```

**3. Create Custom Templates:**
```bash
# Copy templates to customize
lvt template copy all

# Edit templates in .lvt/templates/
# Your customizations apply to all new resources
```

**4. Define Relationships with `references`:**
```bash
# Basic reference (ON DELETE CASCADE - default)
lvt gen comments post_id:references:posts author text

# Custom ON DELETE behavior
lvt gen audit_logs user_id:references:users:set_null action:string
  # Makes user_id nullable, sets NULL when user deleted

# Multiple references
lvt gen likes user_id:references:users post_id:references:posts

# Restrict deletion (prevent deleting parent if children exist)
lvt gen invoices customer_id:references:customers:restrict amount:float
```

**5. Add More Features:**
```bash
# Tags for posts
lvt gen tags name color:string

# Post-tag relationship (many-to-many with references)
lvt gen post_tags post_id:references:posts tag_id:references:tags

# User accounts
lvt gen users username email password_hash:string

# Post reactions with proper relationships
lvt gen reactions post_id:references:posts user_id:references:users type:string
```

### Project Structure

After completing the tutorial, your project looks like:

```
myblog/
├── cmd/myblog/main.go
├── internal/
│   ├── app/
│   │   ├── posts/
│   │   │   ├── posts.go
│   │   │   ├── posts.tmpl
│   │   │   ├── posts_test.go
│   │   │   └── posts_ws_test.go
│   │   ├── categories/
│   │   ├── comments/
│   │   └── home/
│   └── database/
│       ├── db.go
│       ├── migrations/
│       │   ├── 20240315120000_create_posts.sql
│       │   ├── 20240315120001_create_categories.sql
│       │   └── ...
│       ├── queries.sql
│       └── models/          # Generated by sqlc
│           ├── db.go
│           ├── models.go
│           └── queries.sql.go
└── go.mod
```

### Next Steps

1. **Add Authentication** - Integrate session management
2. **Rich Text Editor** - Add markdown or WYSIWYG editor to post content
3. **Image Uploads** - Add image upload functionality
4. **Search** - Implement full-text search across posts
5. **RSS Feed** - Generate RSS feed from posts
6. **Admin Dashboard** - Create `lvt gen view admin`
7. **API Endpoints** - Add JSON API alongside HTML views

### Tips

- **Start simple** - Begin with core resources, add features incrementally
- **Use migrations** - Always use `lvt migration create` for schema changes
- **Test continuously** - Run `go test ./...` after each change
- **Customize templates** - Copy and modify templates to match your design
- **Component mode** - Use `--mode single` for SPA-style applications

## Commands

### `lvt new <app-name>`

Creates a new LiveTemplate application with:

```
myapp/
├── cmd/myapp/main.go           # Application entry point
├── go.mod                      # With //go:tool directive
├── internal/
│   ├── app/                    # Handlers and templates
│   ├── database/
│   │   ├── db.go              # Connection & migrations
│   │   ├── schema.sql         # Database schema
│   │   ├── queries.sql        # SQL queries (sqlc)
│   │   ├── sqlc.yaml          # sqlc configuration
│   │   └── models/            # Generated code
│   └── shared/                # Shared utilities
├── web/assets/                # Static assets
└── README.md
```

### `lvt gen <resource> <field:type>...`

Generates a full CRUD resource with database integration.

**Example:**
```bash
lvt gen posts title:string content:string published:bool views:int
```

**Generated Files:**
- Handler with State struct, Change() method, Init() method
- Bulma CSS template with:
  - Create form with validation
  - List view with search, sort, pagination
  - Delete functionality
  - Real-time WebSocket updates
- WebSocket unit tests
- Chromedp E2E tests
- Database schema and queries

**Features:**
- ✅ CRUD operations (Create, Read, Update, Delete)
- ✅ Search across string fields
- ✅ Sorting by fields
- ✅ Pagination
- ✅ Real-time updates via WebSocket
- ✅ Form validation
- ✅ Statistics/counts
- ✅ Bulma CSS styling
- ✅ Comprehensive tests
- ✅ **Auto-injected routes** - Automatically adds route and import to `main.go`

### `lvt gen view <name>`

Generates a view-only handler without database integration (like the counter example).

**Example:**
```bash
lvt gen view dashboard
```

**Generates:**
- `app/dashboard/dashboard.go` - View handler with state management
- `app/dashboard/dashboard.tmpl` - Bulma CSS template
- `app/dashboard/dashboard_ws_test.go` - WebSocket tests
- `app/dashboard/dashboard_test.go` - Chromedp E2E tests

**Features:**
- ✅ State management
- ✅ Real-time updates via WebSocket
- ✅ Bulma CSS styling
- ✅ Comprehensive tests
- ✅ No database dependencies
- ✅ **Auto-injected routes** - Automatically adds route and import to `main.go`

### `lvt gen auth`

Generates a complete authentication system similar to Phoenix's `mix phx.gen.auth`.

**Example:**
```bash
# Generate with default settings (password + magic-link auth)
lvt gen auth

# Generate with only password authentication
lvt gen auth --no-magic-link

# Generate with only magic-link authentication
lvt gen auth --no-password

# Disable email confirmation
lvt gen auth --no-email-confirm

# Disable CSRF protection
lvt gen auth --no-csrf
```

**Flags:**
- `--no-password` - Disable password authentication
- `--no-magic-link` - Disable magic-link authentication
- `--no-email-confirm` - Disable email confirmation flow
- `--no-password-reset` - Disable password reset functionality
- `--no-sessions-ui` - Disable session management UI
- `--no-csrf` - Disable CSRF protection middleware

**Note:** At least one authentication method (password or magic-link) must be enabled.

**Generates:**
- `shared/password/password.go` - Password hashing utilities (bcrypt)
- `shared/email/email.go` - Email sender interface with console logger
- `database/migrations/YYYYMMDDHHMMSS_create_auth_tables.sql` - Auth tables migration
- Auth queries appended to `database/queries.sql`

**Features:**
- ✅ Password authentication with bcrypt hashing
- ✅ Magic-link email authentication
- ✅ Email confirmation flow
- ✅ Password reset functionality
- ✅ Session management
- ✅ CSRF protection with gorilla/csrf
- ✅ Auto-updates `go.mod` dependencies
- ✅ EmailSender interface (console logger + SMTP/Mailgun examples)
- ✅ Case-insensitive email matching
- ✅ Configurable features via flags

**Database Tables:**
- `users` - User accounts with email and optional hashed password
- `user_tokens` - Tokens for magic links, email confirmation, password reset

**Next Steps After Generation:**
```bash
# 1. Run migrations
lvt migration up

# 2. Generate sqlc code
sqlc generate

# 3. Update main.go to register auth handler
# (Implementation in Phase 2)
```

## Router Auto-Update

When you generate a resource or view, `lvt` automatically:

1. **Adds the import** to your `cmd/*/main.go`:
   ```go
   import (
       "yourapp/app/users"  // ← Auto-added
   )
   ```

2. **Injects the route** after the TODO comment:
   ```go
   // TODO: Add routes here
   http.Handle("/users", users.Handler(queries))  // ← Auto-added
   ```

3. **Maintains idempotency** - Running the same command twice won't duplicate routes

This eliminates the manual step of wiring up routes, making the development workflow smoother. Routes are inserted in the order you generate them, right after the TODO marker.

## Type Mappings

| CLI Type | Go Type      | SQL Type   |
|----------|--------------|------------|
| string   | string       | TEXT       |
| int      | int64        | INTEGER    |
| bool     | bool         | BOOLEAN    |
| float    | float64      | REAL       |
| time     | time.Time    | DATETIME   |

**Aliases:**
- `str`, `text` → `string`
- `integer` → `int`
- `boolean` → `bool`
- `float64`, `decimal` → `float`
- `datetime`, `timestamp` → `time`

## Smart Type Inference (🆕 Phase 1)

The CLI includes an intelligent type inference system that automatically suggests types based on field names:

### How It Works

When using the type inference system, you can omit explicit types and let the system infer them:

```go
// In ui.InferType("email") → returns "string"
// In ui.InferType("age") → returns "int"
// In ui.InferType("price") → returns "float"
// In ui.InferType("enabled") → returns "bool"
// In ui.InferType("created_at") → returns "time"
```

### Inference Rules

**String fields** (default for unknown):
- Exact: `name`, `title`, `description`, `email`, `username`, `url`, `slug`, `address`, etc.
- Contains: `*email*`, `*url*`

**Integer fields:**
- Exact: `age`, `count`, `quantity`, `views`, `likes`, `score`, `rank`, `year`
- Suffix: `*_count`, `*_number`, `*_index`

**Float fields:**
- Exact: `price`, `amount`, `rating`, `latitude`, `longitude`
- Suffix/Contains: `*_price`, `*_amount`, `*_rate`, `*price*`, `*amount*`

**Boolean fields:**
- Exact: `enabled`, `active`, `published`, `verified`, `approved`, `deleted`
- Prefix: `is_*`, `has_*`, `can_*`

**Time fields:**
- Exact: `created_at`, `updated_at`, `deleted_at`, `published_at`
- Suffix: `*_at`, `*_date`, `*_time`

### Usage

The inference system is available via the `ui` package:

```go
import "github.com/livetemplate/lvt/internal/ui"

// Infer type from field name
fieldType := ui.InferType("email")  // → "string"

// Parse field input (with or without type)
name, typ := ui.ParseFieldInput("email")      // → "email", "string" (inferred)
name, typ := ui.ParseFieldInput("age:float")  // → "age", "float" (explicit override)
```

### Future Enhancement

In upcoming phases, this will power:
- Interactive field builders that suggest types as you type
- Direct mode support: `lvt gen users name email age` (without explicit types)
- Smart defaults that reduce typing

## Project Layout

The generated app follows idiomatic Go conventions:

- **`cmd/`** - Application entry points
- **`app/`** - Handlers and templates (co-located!)
- **`database/`** - Database layer with sqlc
- **`shared/`** - Shared utilities
- **`web/assets/`** - Static assets

**Key Design Decision:** Templates live next to their handlers for easy discovery.

## Generated Handler Structure

```go
package users

type State struct {
    Queries        *models.Queries
    Users          []User
    SearchQuery    string
    SortBy         string
    CurrentPage    int
    PageSize       int
    TotalPages     int
    // ...
}

// Action methods - automatically dispatched based on action name
func (s *State) Add(ctx *livetemplate.ActionContext) error {
    // Create user
    return nil
}

func (s *State) Update(ctx *livetemplate.ActionContext) error {
    // Update user
    return nil
}

func (s *State) Delete(ctx *livetemplate.ActionContext) error {
    // Delete user
    return nil
}

func (s *State) Search(ctx *livetemplate.ActionContext) error {
    // Search users
    return nil
}

func (s *State) Init() error {
    // Load initial data
    return nil
}

func Handler(queries *models.Queries) http.Handler {
    tmpl := livetemplate.New("users")
    state := &State{Queries: queries, PageSize: 10}
    return tmpl.Handle(state)
}
```

## Testing

The project includes comprehensive testing infrastructure at multiple levels.

### Make Targets (Recommended)

Use these convenient make targets for different testing workflows:

```bash
make test-fast     # Unit tests only (~30s)
make test-commit   # Before committing (~3-4min)
make test-all      # Full suite (~5-6min)
make test-clean    # Clean Docker resources
```

See [Testing Guide](docs/testing.md) for detailed documentation on test optimization and architecture.

### Quick Start

```bash
# Run all tests (fast mode - skips deployment tests)
go test ./... -short

# Run all tests (including slower e2e tests)
go test ./...

# Run specific package tests
go test ./internal/generator -v

# Run tests with coverage
go test ./... -cover
```

### Test Types

#### 1. Unit Tests
Fast tests for individual packages and functions:

```bash
# Internal packages
go test ./internal/config ./internal/generator ./internal/parser -v

# Commands package
go test ./commands -v
```

**Duration**: <5 seconds

#### 2. WebSocket Tests (`*_ws_test.go`)

Fast unit tests for WebSocket protocol and state changes in generated resources:

```bash
go test ./app/users -run WebSocket
```

**Features**:
- Test server startup with dynamic ports
- WebSocket connection testing
- CRUD action testing
- Server log capture for debugging

**Duration**: 2-5 seconds per resource

#### 3. E2E Browser Tests (`*_test.go`)

Full browser testing with real user interactions for generated resources:

```bash
go test ./app/users -run E2E
```

**Features**:
- Docker Chrome container
- Real browser interactions (clicks, typing, forms)
- Visual verification
- Screenshot capture
- Console log access

**Duration**: 20-60 seconds per resource

#### 4. Deployment Tests (Advanced)

Comprehensive deployment testing infrastructure for testing real deployments:

```bash
# Mock deployment tests (fast, no credentials needed)
go test ./e2e -run TestDeploymentInfrastructure_Mock -v

# Docker deployment tests (requires Docker)
RUN_DOCKER_DEPLOYMENT_TESTS=true go test ./e2e -run TestDockerDeployment -v

# Fly.io deployment tests (requires credentials)
export FLY_API_TOKEN="your_token"
RUN_FLY_DEPLOYMENT_TESTS=true go test ./e2e -run TestRealFlyDeployment -v
```

**Features**:
- Mock, Docker, and Fly.io deployment testing
- Automatic cleanup and resource management
- Smoke tests (HTTP, health, WebSocket, templates)
- Credential-based access control

**Duration**: 2 minutes (mock) to 15 minutes (real deployments)

### Test Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `RUN_DOCKER_DEPLOYMENT_TESTS` | Enable Docker deployment tests | `false` |
| `RUN_FLY_DEPLOYMENT_TESTS` | Enable Fly.io deployment tests | `false` |
| `FLY_API_TOKEN` | Fly.io API token for real deployments | - |

### Continuous Integration

Tests run automatically on every pull request via GitHub Actions:

- ✅ Code formatting validation
- ✅ Unit tests (all internal packages)
- ✅ Commands tests
- ✅ E2E tests (short mode)
- ✅ Mock deployment tests

**On-demand/scheduled deployment testing** available via manual workflow dispatch or weekly schedule.

For detailed CI/CD documentation, see:
- [CI Deployment Testing Guide](e2e/CI_DEPLOYMENT_TESTING.md)
- [Deployment Testing Documentation](e2e/DEPLOYMENT_TESTING.md)

### Skip Slow Tests

Use `-short` flag to skip slow tests (deployment tests, long-running e2e tests):

```bash
go test -short ./...
```

### Test Documentation

For comprehensive testing documentation, see:
- **[Deployment Testing](e2e/DEPLOYMENT_TESTING.md)** - Complete deployment testing guide
- **[CI/CD Testing](e2e/CI_DEPLOYMENT_TESTING.md)** - CI/CD workflows and setup
- **[Deployment Plan](e2e/DEPLOYMENT_TESTING_PLAN_UPDATE.md)** - Implementation progress and status

## Go 1.24+ Tools Support

Generated `go.mod` includes:

```go
//go:tool github.com/sqlc-dev/sqlc/cmd/sqlc
```

Run migrations (automatically runs sqlc):
```bash
lvt migration up
```

## CSS Framework

All generated templates use [Bulma CSS](https://bulma.io/) by default:

```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@1.0.4/css/bulma.min.css">
```

Components used:
- `.section`, `.container` - Layout
- `.box` - Content containers
- `.table` - Data tables
- `.button`, `.input`, `.select` - Form controls
- `.pagination` - Pagination controls

## Development Workflow

1. **Create app:** `lvt new myapp`
2. **Generate resources:** `lvt gen users name:string email:string`
3. **Run migrations:** `lvt migration up` (auto-generates DB code)
4. **Wire routes** in `main.go`
5. **Run tests:** `go test ./...`
6. **Run app:** `go run cmd/myapp/main.go`

## Examples

### Blog App

```bash
lvt new myblog
cd myblog

# Generate posts resource
lvt gen posts title:string content:string published:bool

# Generate comments resource
lvt gen comments post_id:string author:string text:string

# Run migrations (auto-generates DB code)
lvt migration up

# Run
go run cmd/myblog/main.go
```

### E-commerce

```bash
lvt new mystore
cd mystore

lvt gen products name:string price:float stock:int
lvt gen customers name:string email:string
lvt gen orders customer_id:string total:float

lvt migration up  # Runs migrations and generates DB code
go run cmd/mystore/main.go
```

## Architecture

### Template System

The generator uses custom delimiters (`[[`, `]]`) to avoid conflicts with Go template syntax:

- **Generator templates:** `[[.ResourceName]]` - Replaced during generation
- **Output templates:** `{{.Title}}` - Used at runtime by LiveTemplate

### Embedded Templates

All templates are embedded using `embed.FS` for easy distribution.

### Code Generation Strategy

1. Parse field definitions (`name:type`)
2. Map types to Go and SQL types
3. Render templates with resource data
4. Generate handler, template, tests
5. Append to database files

## Testing the Generator

### Run All Tests

```bash
go test ./cmd/lvt -v
```

### Test Layers

1. **Parser Tests** (`cmd/lvt/internal/parser/fields_test.go`)
   - Field parsing and validation
   - Type mapping correctness
   - 13 comprehensive tests

2. **Golden File Tests** (`cmd/lvt/golden_test.go`)
   - Regression testing for generated code
   - Validates handler and template output
   - Update with: `UPDATE_GOLDEN=1 go test ./cmd/lvt -run Golden`

3. **Integration Tests** (`cmd/lvt/integration_test.go`)
   - Go syntax validation
   - File structure validation
   - Generation pipeline testing

4. **Smoke Test** (`scripts/test_cli_smoke.sh`)
   - End-to-end CLI workflow
   - App creation and resource generation
   - File structure verification

## Roadmap

- [x] ~~`lvt gen view` - View-only handlers~~ ✅ Complete
- [x] ~~Router auto-update~~ ✅ Complete
- [x] ~~Bubbletea interactive UI~~ ✅ Complete (Phase 1-3)
  - [x] Dependencies & infrastructure
  - [x] Smart type inference system (50+ patterns)
  - [x] UI styling framework (Lipgloss)
  - [x] Interactive app creation wizard
  - [x] Interactive resource builder
  - [x] Interactive view builder
  - [x] Mode detection (auto-switch based on args)
  - [x] Type inference in direct mode
  - [x] ~~Enhanced validation & help system (Phase 4)~~ ✅ Complete
    - [x] Real-time Go identifier validation
    - [x] SQL reserved word warnings (25+ keywords)
    - [x] Help overlay with `?` key in all wizards
    - [x] Color-coded feedback (✓✗⚠)
    - [x] All 3 wizards enhanced
- [x] ~~Migration commands~~ ✅ Complete
  - [x] Goose integration with minimal wrapper (~410 lines)
  - [x] Auto-generate migrations from `lvt gen resource`
  - [x] Commands: `up`, `down`, `status`, `create <name>`
  - [x] Timestamped migration files with Up/Down sections
  - [x] Schema versioning and rollback support
- [x] ~~Custom template support~~ ✅ Complete
  - [x] Cascading template lookup (project → user → embedded)
  - [x] `lvt template copy` command for easy customization
  - [x] Project templates in `.lvt/templates/` (version-controlled)
  - [x] User-wide templates in `~/.config/lvt/templates/`
  - [x] Selective override (only customize what you need)
  - [x] Zero breaking changes (~250 lines total)
- [x] ~~Multiple CSS frameworks~~ ✅ Complete
  - [x] Tailwind CSS v4 (default)
  - [x] Bulma 1.0.4
  - [x] Pico CSS v2
  - [x] None (pure HTML)
  - [x] CSS framework determined by kit (multi/single use Tailwind, simple uses Pico)
  - [x] 57 CSS helper functions for framework abstraction
  - [x] Conditional template rendering (single source of truth)
  - [x] Semantic HTML support for Pico CSS (<main>, <article>)
  - [x] Zero breaking changes (~550 lines total)
- [x] ~~`lvt gen auth` - Authentication system~~ ✅ Phase 1 Complete
  - [x] Password authentication (bcrypt)
  - [x] Magic-link email authentication
  - [x] Email confirmation flow
  - [x] Password reset functionality
  - [x] Session management tables
  - [x] CSRF protection (gorilla/csrf)
  - [x] Auto-dependency updates (go.mod)
  - [x] EmailSender interface with examples
  - [x] Configurable via flags
  - [ ] Auth handlers (Phase 2)
  - [ ] Custom authenticator (Phase 3)
  - [ ] Middleware templates (Phase 4)
- [ ] GraphQL support

## Contributing

See the main [LiveTemplate CLAUDE.md](../../CLAUDE.md) for development guidelines.

## License

Same as LiveTemplate project.
