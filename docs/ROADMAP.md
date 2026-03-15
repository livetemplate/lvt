# LVT CLI Feature Gap Roadmap

> A prioritized tracker of missing features for production-ready web applications, based on deep comparison with Rails 8, Phoenix LiveView, Laravel 12, Django, and Go frameworks (Buffalo, Beego).

## Executive Summary

LVT is already strong in code generation, authentication, real-time UI (WebSockets), database migrations (goose/sqlc), E2E testing (chromedp), deployment stacks, and developer experience (hot reload, AI agent integration, 20+ UI components).

However, **19 features** common across mature frameworks are missing or incomplete. These are organized into 4 milestones, from production blockers to developer experience improvements.

---

## Feature Comparison Matrix

| Feature                  | LVT  | Rails 8 | Phoenix | Laravel 12 | Django |
|--------------------------|------|---------|---------|------------|--------|
| Code Generation          | ✅   | ✅      | ✅      | ✅         | ✅     |
| Authentication           | ✅   | ✅      | ✅      | ✅         | ✅     |
| Real-time UI             | ✅   | ✅      | ✅      | ✅         | ⚠️     |
| Database Migrations      | ✅   | ✅      | ✅      | ✅         | ✅     |
| E2E Testing              | ✅   | ✅      | ✅      | ✅         | ✅     |
| Deployment Stacks        | ✅   | ✅      | ✅      | ✅         | ⚠️     |
| Hot Reload               | ✅   | ✅      | ✅      | ✅         | ⚠️     |
| Authorization / RBAC     | ❌   | ✅      | ⚠️     | ✅         | ✅     |
| Background Jobs          | ❌   | ✅      | ✅      | ✅         | ✅     |
| Email Sending            | ⚠️   | ✅      | ✅      | ✅         | ✅     |
| File Uploads & Storage   | ❌   | ✅      | ✅      | ✅         | ✅     |
| Form Validation          | ⚠️   | ✅      | ✅      | ✅         | ✅     |
| API / JSON Endpoints     | ❌   | ✅      | ✅      | ✅         | ✅     |
| Caching                  | ❌   | ✅      | ✅      | ✅         | ✅     |
| Admin Panel              | ❌   | ✅      | ❌      | ✅         | ✅     |
| Request Logging          | ⚠️   | ✅      | ✅      | ✅         | ✅     |
| Rate Limiting            | ❌   | ✅      | ⚠️ (Plug.Throttle) | ✅ | ✅ |
| Middleware Pipeline      | ⚠️   | ✅      | ✅      | ✅         | ✅     |
| i18n / Localization      | ❌   | ✅      | ✅      | ✅         | ✅     |
| Full-Text Search         | ❌   | ⚠️     | ❌      | ✅         | ✅     |
| WebSocket Channels       | ⚠️   | ✅      | ✅      | ✅         | ❌     |
| Scheduled Tasks          | ❌   | ✅      | ✅      | ✅         | ⚠️     |
| Asset Pipeline           | ❌   | ✅      | ✅      | ✅         | ⚠️     |
| DB Console / REPL        | ❌   | ✅      | ✅      | ✅         | ✅     |

**Legend**: ✅ Built-in / First-class | ⚠️ Partial / Basic | ❌ Missing

> **Package convention**: `pkg/` contains packages copied into generated applications (cookie, email, flash, password, token). `internal/` contains lvt CLI internals. Roadmap items follow this convention — new runtime packages go in `pkg/`, new CLI/generator packages go in `internal/`.

---

## Milestone 1: Unblock Production Use

**Goal**: Address the features that currently prevent lvt applications from being deployed to production with confidence.
**Estimated effort**: L (large) — 3 features, each requiring new packages, generator changes, and tests.

### 1.1 Background Job / Task Queue System

**Priority**: Critical — every real application needs async work (email delivery, image processing, report generation, webhook dispatch, data cleanup).

**What competitors offer**:
- Rails 8: Solid Queue (database-backed, zero external dependencies), Active Job unified API, job continuations (8.1)
- Laravel 12: Built-in queue with failover driver, batch processing, job chaining, retry logic
- Django: Celery (mature, Redis/RabbitMQ-backed)
- Phoenix: OTP processes, GenServer, Task.Supervisor — concurrency is first-class in Elixir
- Beego: Built-in task scheduler for cron jobs

**Acceptance Criteria**:
- [ ] Database-backed job queue table (SQLite-compatible, no Redis dependency)
- [ ] Worker pool using goroutines with configurable concurrency
- [ ] Retry logic with exponential backoff and max retry count
- [ ] Scheduled/delayed job execution (run at specific time)
- [ ] Job failure tracking with error messages and stack traces
- [ ] `lvt gen job <name>` command to scaffold new job types
- [ ] Job status API (pending, running, completed, failed)
- [ ] Graceful shutdown — finish in-progress jobs before exiting
- [ ] Dead letter queue for permanently failed jobs
- [ ] Integration with generated app's main.go (auto-start workers)

---

### 1.2 Production Email Sender

**Priority**: Critical — auth features (magic links, password reset) already generate emails but use `ConsoleEmailSender` (stdout only). These features are unusable in production.

**What competitors offer**:
- Rails: Action Mailer with SMTP, Mailgun, SendGrid, SES adapters + email preview in dev
- Laravel: Built-in Mail with drivers for SMTP, Mailgun, SES, Postmark, Resend
- Django: `django.core.mail` with SMTP backend + third-party providers
- Phoenix: Swoosh library with 10+ provider adapters

**Acceptance Criteria**:
- [ ] SMTP email sender implementation in `pkg/email/`
- [ ] Adapter interface for third-party providers (SendGrid, SES, Resend)
- [ ] At least one cloud provider adapter implemented (e.g., SendGrid)
- [ ] Email preview mode for development (render to browser instead of sending)
- [ ] HTML and plain-text email template support
- [ ] Configuration via environment variables (`SMTP_HOST`, `SMTP_PORT`, etc.)
- [ ] `lvt env generate` updated to detect email config requirements
- [ ] Connection pooling / reuse for SMTP connections
- [ ] Error handling with meaningful error messages on delivery failure
- [ ] Documentation with setup examples for common providers

---

### 1.3 Server-Side Form Validation

**Priority**: Critical — database constraints alone provide poor UX. Users need field-level validation with friendly error messages before data hits the database.

**What competitors offer**:
- Rails: Active Model Validations (presence, format, length, uniqueness, numericality, custom)
- Laravel: 90+ built-in validation rules, form requests, custom error messages
- Django: Form/Model validation, clean methods, field validators
- Phoenix: Ecto changesets with validate_required, validate_format, validate_length, etc.

**Acceptance Criteria**:
- [ ] Validation package in `pkg/validate/` with common rules (required, email, min/max length, numeric range, regex pattern, URL, in-list)
- [ ] Validation rules auto-generated from resource field types (e.g., `email:string` → email format rule, `age:integer` → numeric rule)
- [ ] Custom validator support (user-defined validation functions)
- [ ] Validation errors returned as structured data to templates
- [ ] Generated templates display inline field-level error messages
- [ ] Validation runs in handlers before database operations
- [ ] `--validate` flag or automatic validation in `lvt gen resource`
- [ ] Uniqueness validation with database check
- [ ] Cross-field validation support (e.g., password confirmation)
- [ ] Client-side validation attributes generated in HTML (e.g., `required`, `minlength`)

---

## Milestone 2: Feature Parity with Major Frameworks

**Goal**: Reach feature parity with the core capabilities that every major framework provides out of the box.
**Estimated effort**: XL (extra large) — 3 features with significant scope (file uploads, RBAC, API layer).

### 2.1 File Upload & Storage

**Priority**: High — almost every real application needs file uploads (profile photos, documents, attachments). This is table-stakes.

**What competitors offer**:
- Rails: Active Storage (local disk, S3, GCS, Azure) with image variants, direct uploads, content type validation
- Laravel: Filesystem abstraction (local, S3, FTP) with file validation and temporary URLs
- Django: `FileField`/`ImageField` on models, configurable storage backends, upload handlers
- Phoenix LiveView: Built-in uploads with progress bars, drag-and-drop, direct-to-cloud, chunk uploads

**Acceptance Criteria**:
- [ ] Storage package in `pkg/storage/` with `Store` interface
- [ ] Local disk storage adapter (configurable upload directory)
- [ ] S3-compatible storage adapter (works with AWS S3, MinIO, DigitalOcean Spaces)
- [ ] New field types in resource generation: `file` and `image`
- [ ] Multipart form handling in generated handlers
- [ ] File size limits and content type validation
- [ ] Generated UI with file input, drag-and-drop zone, and upload progress
- [ ] Image resizing/thumbnail generation (using Go's image packages)
- [ ] Secure file serving (prevent directory traversal, content type sniffing)
- [ ] Database migration includes file metadata columns (filename, content_type, size, storage_key)
- [ ] Cleanup of orphaned files on record deletion

---

### 2.2 Authorization / Role-Based Access Control

**Priority**: High — authentication without authorization is incomplete. Real apps need "who can do what."

**What competitors offer**:
- Rails: Pundit (policy objects), CanCanCan (ability-based authorization)
- Laravel: Gates & Policies (built-in), Spatie roles/permissions package
- Django: Built-in permission system with groups, object-level permissions, `@permission_required`
- Phoenix: Plugs and pipelines for scope-based authorization

**Acceptance Criteria**:
- [ ] Role system with configurable roles (at minimum: admin, user; extensible to custom roles)
- [ ] Policy package in `pkg/authz/` defining `Can(user, action, resource)` interface
- [ ] Authorization middleware for route-level protection
- [ ] Resource ownership checks (e.g., "only the creator can edit this record")
- [ ] `--with-authz` flag on `lvt gen resource` to generate authorization checks
- [ ] Role assignment during user registration and via admin
- [ ] Database migration for roles table and user_roles join table
- [ ] Generated handler code includes authorization checks before CRUD operations
- [ ] 403 Forbidden responses with appropriate error pages
- [ ] Authorization integrated with existing auth middleware

---

### 2.3 API / JSON Endpoints

**Priority**: High — modern applications need APIs for mobile clients, third-party integrations, SPAs, and internal microservices.

**Depends on**: Existing auth system for bearer token authentication. Should integrate with Rate Limiting (3.3) and Middleware Pipeline (3.4) once available.

**What competitors offer**:
- Rails: `respond_to` format blocks, API-only mode, ActiveModel Serializers, Jbuilder
- Laravel: API resources, API-only scaffolding, Sanctum/Passport for API auth
- Django: Django REST Framework (ViewSets, serializers, browsable API, throttling)
- Phoenix: JSON views, API pipelines, OpenAPI spec generation

**Acceptance Criteria**:
- [ ] `lvt gen api <resource> <fields...>` command generating JSON CRUD endpoints
- [ ] JSON serialization/deserialization for all generated field types
- [ ] RESTful route conventions (GET/POST/PUT/DELETE with proper status codes)
- [ ] API authentication via bearer tokens (integrated with existing auth system)
- [ ] Proper error responses in JSON format with error codes and messages
- [ ] Content-type negotiation (Accept header handling)
- [ ] API versioning support (URL prefix: `/api/v1/`)
- [ ] Generated API tests
- [ ] Pagination in JSON responses (page, per_page, total, links)
- [ ] CORS middleware for cross-origin requests

---

## Milestone 3: Production Hardening

**Goal**: Add the infrastructure features needed for reliable, observable, and secure production deployments.
**Estimated effort**: L (large) — 4 features, mostly middleware and infrastructure packages.

### 3.1 Request Logging Middleware

**Priority**: Medium-High — production applications need visibility into request traffic for debugging, performance monitoring, and audit trails.

**What competitors offer**:
- Rails: Tagged logging, Active Support instrumentation, request logging with timing
- Laravel: Logging channels (stack, daily, Slack, Sentry), Telescope for debugging
- Django: Python logging integration, django-debug-toolbar
- Beego: Built-in monitoring for performance and memory tracking

**Acceptance Criteria**:
- [ ] Request logging middleware in generated app (method, path, status code, response time, request ID)
- [ ] Structured JSON log output compatible with log aggregation tools (ELK, Datadog, etc.)
- [ ] Configurable log level (debug, info, warn, error)
- [ ] Request ID generation and propagation (X-Request-ID header)
- [ ] Sensitive data redaction (passwords, tokens, etc.)
- [ ] Access log format option (Apache/Nginx compatible)
- [ ] Slow request detection and warning (configurable threshold)
- [ ] Integration with existing `slog` structured logging

---

### 3.2 Caching Layer

**Priority**: Medium-High — database queries, rendered templates, and computed data all benefit from caching as applications scale.

**What competitors offer**:
- Rails: Solid Cache (database-backed), Russian doll caching, fragment caching, HTTP caching headers
- Laravel: Cache facade with file, database, Redis, Memcached drivers, tagged caching
- Django: Multi-level caching (per-site, per-view, template fragment, low-level API)
- Phoenix: ETS-based caching (built into Erlang VM), Cachex library

**Acceptance Criteria**:
- [ ] Cache package in `pkg/cache/` with `Cache` interface (Get, Set, Delete, Flush)
- [ ] In-memory cache adapter (with TTL and LRU eviction)
- [ ] SQLite-backed cache adapter (persistent across restarts)
- [ ] Template fragment caching in LiveTemplate engine
- [ ] Cache invalidation on database writes for related resources
- [ ] HTTP caching headers (ETag, Cache-Control, Last-Modified) for static resources
- [ ] Configurable TTL per cache key or group
- [ ] Cache statistics / hit rate logging

---

### 3.3 Rate Limiting

**Priority**: High — protects against brute force attacks on auth endpoints (login, password reset, magic links), abuse, and ensures fair resource usage. Arguably a production blocker given that auth is already generated.

**Depends on**: Auth middleware (existing). Should be applied to auth endpoints immediately and integrated into Middleware Pipeline (3.4) for general use.

**What competitors offer**:
- Laravel: Built-in RateLimiter facade with named limiters, per-route configuration
- Rails: Rack::Attack middleware (throttle, blocklist, safelist)
- Django: DRF throttling (user-based, IP-based, scoped)

**Acceptance Criteria**:
- [ ] Rate limiting middleware in `pkg/ratelimit/`
- [ ] Per-route rate limit configuration
- [ ] Per-IP and per-user rate limiting strategies
- [ ] Configurable time windows and request thresholds
- [ ] Returns HTTP 429 (Too Many Requests) with `Retry-After` header
- [ ] SQLite-backed rate limit storage (no Redis dependency)
- [ ] Automatic rate limiting on auth endpoints (login, registration, password reset)
- [ ] Allowlist/blocklist support for IPs

---

### 3.4 Composable Middleware Pipeline

**Priority**: Medium — real applications need different middleware stacks for different route groups (web vs API, public vs authenticated).

**What competitors offer**:
- Rails: `before_action`, `around_action`, controller concerns with scoping
- Laravel: Middleware groups (`web`, `api`), route middleware, global middleware, terminable middleware
- Phoenix: Plug pipelines per scope (`:browser`, `:api`), pipeline composition

**Acceptance Criteria**:
- [ ] Named middleware groups (e.g., `web`, `api`, `admin`)
- [ ] Per-route or per-group middleware assignment in router
- [ ] Middleware ordering control (before/after dependencies)
- [ ] `lvt new` generates default middleware groups
- [ ] `lvt gen resource` and `lvt gen api` assign appropriate middleware groups
- [ ] Easy custom middleware registration with clear interface
- [ ] Middleware short-circuiting (stop chain on auth failure, rate limit, etc.)

---

## Milestone 4: Developer Experience & Maturity

**Goal**: Quality-of-life features that improve developer productivity and bring lvt to the level of polish expected from mature frameworks.
**Estimated effort**: XL (extra large) — 7 features spanning admin UI, search, i18n, asset pipeline, and more.

### 4.1 Admin Panel Generator

**Priority**: Medium — every application needs data management. Auto-generated admin reduces boilerplate significantly.

**Depends on**: Authorization / RBAC (2.2) — admin panel requires role-based access control to restrict to admin users.

**What competitors offer**:
- Django: Auto-generated admin from models (Django's killer feature — register model, get full CRUD)
- Laravel: Nova, Filament, Backpack (rich admin panel ecosystems)
- Rails: ActiveAdmin, RailsAdmin, Administrate

**Acceptance Criteria**:
- [ ] `lvt gen admin` command that generates admin dashboard for all existing resources
- [ ] Auto-generated list/detail/create/edit views per resource
- [ ] User management interface (list users, toggle roles, confirm emails)
- [ ] Filtering and search across admin tables
- [ ] Admin-only middleware protection (admin role required)
- [ ] Responsive admin layout with navigation sidebar

---

### 4.2 Full-Text Search (SQLite FTS5)

**Priority**: Medium — search is expected in most applications. SQLite FTS5 is a natural fit for lvt's SQLite-first approach.

**What competitors offer**:
- Laravel: Scout with Algolia/Meilisearch drivers
- Django: Built-in PostgreSQL full-text search, django-haystack
- Rails: pg_search gem, Searchkick (Elasticsearch)

**Acceptance Criteria**:
- [ ] SQLite FTS5 virtual table generation for searchable resources
- [ ] `searchable` flag or `search:fts` field type in resource generation
- [ ] Auto-generated search input in resource list templates
- [ ] FTS5 query builder with match highlighting
- [ ] Search results ranked by relevance
- [ ] Trigger-based FTS index sync on insert/update/delete

---

### 4.3 WebSocket Channels / PubSub

**Priority**: Medium — lvt already has WebSocket support but lacks topic-based broadcasting for multi-user features (chat, notifications, collaborative editing).

**What competitors offer**:
- Phoenix: Channels with topics, presence tracking, distributed PubSub across nodes
- Rails: Action Cable with channels, subscriptions, broadcasting
- Laravel: Broadcasting with Echo, private/presence channels

**Acceptance Criteria**:
- [ ] Topic/channel abstraction on top of existing WebSocket manager
- [ ] Client can subscribe to specific channels (e.g., `room:123`, `user:456`)
- [ ] Server-side broadcast to channel (only subscribers receive messages)
- [ ] Presence tracking (who's currently viewing a page/channel)
- [ ] Private channels with authorization checks
- [ ] Updated client-side JavaScript library to support channel subscriptions

---

### 4.4 Scheduled Tasks / Cron

**Priority**: Medium — recurring tasks (cleanup, reports, data sync) need scheduling without external cron.

**What competitors offer**:
- Laravel: Task Scheduling (`schedule:run` — define schedules in code, run via single cron entry)
- Rails: whenever gem + Solid Queue recurring jobs
- Beego: Built-in task scheduler

**Acceptance Criteria**:
- [ ] Scheduler built on top of job queue (Milestone 1.1)
- [ ] `lvt gen task <name> --schedule "<cron>"` command to scaffold recurring tasks
- [ ] Cron expression support (e.g., `"0 * * * *"` for hourly)
- [ ] Named interval shortcuts (e.g., `@daily`, `@hourly`, `@every 5m`)
- [ ] Overlap prevention (don't run if previous instance still running)
- [ ] Task status tracking and last-run timestamp

---

### 4.5 Internationalization (i18n) / Localization

**Priority**: Low-Medium — essential for applications targeting multiple languages/regions.

**What competitors offer**:
- Rails: Built-in I18n with YAML locale files, pluralization, date/number formatting
- Laravel: Translation files with nested keys, pluralization, locale switching
- Django: Comprehensive i18n/l10n with gettext, timezone support, format localization
- Phoenix: Gettext-based translations with PO files

**Acceptance Criteria**:
- [ ] Translation key system (e.g., `t("resource.create.title")`)
- [ ] Locale files in YAML or JSON format
- [ ] Locale switching middleware (Accept-Language header, URL prefix, or cookie)
- [ ] Generated templates use translation keys instead of hardcoded strings
- [ ] Pluralization support
- [ ] Date, time, and number formatting per locale
- [ ] `lvt gen i18n` to extract hardcoded strings from existing templates

---

### 4.6 Asset Pipeline

**Priority**: Low-Medium — lvt currently relies on CDN-hosted CSS frameworks (Tailwind, Pico) with no local asset compilation step. A basic asset pipeline would enable custom CSS/JS bundling and optimized production builds.

**What competitors offer**:
- Rails: Propshaft (default in Rails 8), import maps for zero-build JavaScript
- Laravel: Vite integration for JS/CSS bundling, asset versioning
- Phoenix: esbuild integration, Tailwind plugin, asset digest for cache-busting
- Django: collectstatic, django-compressor, whitenoise for static serving

**Acceptance Criteria**:
- [ ] Static asset serving middleware with proper cache headers
- [ ] CSS framework integration (download Tailwind CSS for offline use)
- [ ] Asset fingerprinting/digest for cache-busting in production
- [ ] Optional esbuild or Tailwind CLI integration for local builds
- [ ] `lvt build` command to compile and minify assets for production
- [ ] Generated apps include static file handling in deployment config

---

### 4.7 Database Console

**Priority**: Low — quality-of-life improvement for development and debugging.

**What competitors offer**:
- Rails: `rails console` (Ruby REPL with app context), `rails dbconsole` (database CLI)
- Laravel: `php artisan tinker` (interactive REPL)
- Django: `manage.py shell` (Python REPL), `manage.py dbshell` (database CLI)
- Phoenix: `iex -S mix` (Elixir REPL with app context)

**Acceptance Criteria**:
- [ ] `lvt console` command opens interactive database shell
- [ ] Auto-detects database file location from app configuration
- [ ] SQLite: launches `sqlite3` with app's database
- [ ] PostgreSQL: launches `psql` with connection string from environment
- [ ] Loads schema context (shows tables on connect)
- [ ] History support for command recall

---

## Sources

- [Rails 8.0 Release Notes](https://guides.rubyonrails.org/8_0_release_notes.html)
- [Phoenix LiveView Documentation](https://hexdocs.pm/phoenix_live_view/Phoenix.LiveView.html)
- [Phoenix 1.8.0 Release Blog](https://www.phoenixframework.org/blog/phoenix-1-8-released)
- [Laravel 12 Release Notes](https://laravel.com/docs/12.x/releases)
- [Top Go Web Frameworks 2025](https://blog.logrocket.com/top-go-frameworks-2025/)
- [Go Web Frameworks Comparison](https://www.monocubed.com/blog/golang-web-frameworks/)
- [Django Channels Documentation](https://channels.readthedocs.io/en/stable/)
- [Django vs Next.js Comparison](https://stackshare.io/stackups/django-vs-next-js)
