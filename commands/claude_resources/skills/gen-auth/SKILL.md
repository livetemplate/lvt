---
name: lvt-gen-auth
description: Add authentication to an existing LiveTemplate app - generates session management, password auth, magic links, email confirmation, CSRF protection, and E2E tests
keywords: ["lvt", "livetemplate", "lt"]
category: core
version: 1.0.0
---

# lvt-gen-auth

Adds a complete authentication system to an existing LiveTemplate application. Generates password authentication, magic link authentication, session management, email confirmation, password reset, CSRF protection, route protection middleware, and comprehensive E2E tests.

## 🎯 ACTIVATION RULES

### Context Detection

This skill typically runs in **existing LiveTemplate projects** (.lvtrc exists).

**✅ Context Established By:**
1. **Project context** - `.lvtrc` exists (most common scenario)
2. **Agent context** - User is working with `lvt-assistant` agent
3. **Keyword context** - User mentions "lvt", "livetemplate", or "lt"

**Keyword matching** (case-insensitive): `lvt`, `livetemplate`, `lt`

### Trigger Patterns

**With Context:**
✅ "add authentication"
✅ "I need user login"
✅ "generate auth system"

**Without Context (needs keywords):**
✅ "add authentication to my lvt app"
✅ "use livetemplate to add user login"
❌ "add authentication" (no context, no keywords)

---

## User Prompts

This skill should activate when the user requests to add authentication:

**Explicit prompts:**
- "Add authentication to my app"
- "Generate auth for my app"
- "I need user authentication"
- "Add login and registration"
- "Set up auth system"

**Implicit prompts:**
- "Add user accounts"
- "I want users to sign in"
- "Add password login"
- "Set up user sessions"
- "Protect routes with authentication"
- "Add magic link login"

**Examples:**
- "Add authentication with password and magic links"
- "Generate auth system for my blog app"
- "I need user login for my app"
- "Set up authentication with email confirmation"
- "Add auth but skip magic links"

## Context Awareness

Before executing this skill, verify:

1. **In Project Directory:**
   - Check for `.lvtrc` file (confirms it's an lvt project)
   - Check for `go.mod` (confirms it's a Go project)
   - Check for `internal/database/` directory

2. **Dependencies Available:**
   - `lvt` binary is installed and accessible
   - Project was created with `lvt new`
   - Database initialized

3. **Not Already Exists:**
   - Check if `internal/app/auth/` already exists
   - Warn user if auth is already generated (will overwrite)

## Checklist

- [ ] **Step 1:** Verify we're in an lvt project directory
  - Check for `.lvtrc` file
  - Check for `go.mod` file
  - Check for `internal/database/` directory
  - If missing, inform user they need to create an app first (use lvt:new-app skill)

- [ ] **Step 2:** Validate prerequisites
  - Verify `lvt` command is available
  - Check current directory is project root
  - Verify database directory exists

- [ ] **Step 3:** Understand user requirements
  - Ask if they want password auth (default: yes)
  - Ask if they want magic link auth (default: yes)
  - Ask if they want email confirmation (default: yes)
  - Ask if they want password reset (default: yes)
  - Ask if they want sessions UI (default: yes)
  - Ask if they want CSRF protection (default: yes)
  - Ask for struct name (default: User)
  - Ask for table name (default: users)

- [ ] **Step 4:** Check for existing auth
  - Check if `internal/app/auth/` already exists
  - If exists, ask user if they want to:
    - Overwrite existing auth (warning: will lose customizations)
    - Cancel operation
    - Use different struct/table names

- [ ] **Step 5:** Build and run the `lvt gen auth` command
  - Format: `lvt gen auth [StructName] [table_name] [flags...]`
  - Default (all features): `lvt gen auth`
  - With flags: `lvt gen auth User users --no-magic-link --no-email-confirm`
  - Available flags:
    - `--no-password` - Skip password authentication
    - `--no-magic-link` - Skip magic link authentication
    - `--no-email-confirm` - Skip email confirmation
    - `--no-password-reset` - Skip password reset
    - `--no-sessions-ui` - Skip sessions management UI
    - `--no-csrf` - Skip CSRF protection

- [ ] **Step 6:** Verify auth generation succeeded
  - Check for success message from lvt
  - Verify files created:
    - `internal/app/auth/auth.go` (handler with all auth flows)
    - `internal/app/auth/auth.tmpl` (LiveTemplate UI)
    - `internal/app/auth/middleware.go` (route protection middleware)
    - `internal/app/auth/auth_e2e_test.go` (E2E tests)
    - `internal/shared/password/password.go` (bcrypt utilities)
    - `internal/shared/email/email.go` (email sender interface)
  - Verify files updated:
    - `internal/database/migrations/<timestamp>_create_auth_tables.sql`
    - `internal/database/queries.sql` (auth queries appended)

- [ ] **Step 7:** Run database migration
  - Execute: `lvt migration up`
  - Verify migration succeeded
  - Handle errors if migration fails

- [ ] **Step 8:** Generate sqlc models
  - Navigate to `internal/database/` directory
  - Run: `go run github.com/sqlc-dev/sqlc/cmd/sqlc generate`
  - Verify models generated successfully (User, Session types)
  - Return to project root

- [ ] **Step 9:** Run `go mod tidy`
  - Ensure all dependencies are up to date
  - Verify no errors

- [ ] **Step 10:** Verify app builds successfully
  - For multi/single kits: `go build ./cmd/<app>`
  - For simple kit: `go build`
  - If build fails, diagnose and fix issues

- [ ] **Step 11:** Guide user on wiring auth routes
  - Explain that routes are NOT auto-injected (by design)
  - Show example code for registering auth routes in main.go
  - Explain that user needs to decide which routes to expose
  - Provide clear code examples from internal/app/auth/auth.go comments

- [ ] **Step 12:** Guide user on configuring email sender
  - Explain that email interface needs implementation
  - Show example for console email (development)
  - Show example for SMTP email (production)
  - Point to internal/shared/email/email.go for interface

- [ ] **Step 13:** Provide user with success summary
  - List files created
  - List files updated
  - Show next steps (wire routes, configure email, run tests)
  - Provide example code for route wiring
  - Explain how to protect routes with middleware

## Authentication Features

### Password Authentication

**What it generates:**
- Registration with email + password
- Login with email + password
- Password hashing with bcrypt (cost 12)
- Password validation (min 8 chars)
- Logout

**Database:**
- `users` table with email, password_hash columns
- Unique index on email

**Routes:**
- POST /auth/register - User registration
- POST /auth/login - User login
- POST /auth/logout - User logout

### Magic Link Authentication

**What it generates:**
- Passwordless login via email
- Secure token generation (32 bytes, crypto/rand)
- Token expiration (15 minutes default)
- One-time use tokens

**Database:**
- `magic_link_tokens` table
- Tokens expire after use or timeout

**Routes:**
- POST /auth/magic-link/request - Request magic link
- GET /auth/magic-link/verify - Verify and login

### Email Confirmation

**What it generates:**
- Email verification tokens
- Confirmation email sending
- User status tracking (verified/unverified)
- Resend confirmation option

**Database:**
- `email_verification_tokens` table
- `users.email_verified` boolean column

**Routes:**
- POST /auth/confirm/send - Send confirmation email
- GET /auth/confirm/verify - Verify email

### Password Reset

**What it generates:**
- Password reset token generation
- Reset email sending
- Secure token validation
- Password update after reset

**Database:**
- `password_reset_tokens` table
- Tokens expire after 1 hour

**Routes:**
- POST /auth/password/reset/request - Request password reset
- POST /auth/password/reset/verify - Verify token and reset password

### Session Management

**What it generates:**
- Secure session tokens (32 bytes)
- Session storage in database
- Session expiration (30 days default)
- Active session tracking
- Session revocation

**Database:**
- `sessions` table with user_id, token, expires_at
- Index on token for fast lookup

**Middleware:**
- RequireAuth(next http.Handler) - Protect routes
- Session validation on each request
- Auto-cleanup of expired sessions

### CSRF Protection

**What it generates:**
- CSRF token generation per session
- Token validation middleware
- Form helper for templates
- Automatic token rotation

**How it works:**
- Token stored in session
- Must be included in forms
- Validated on state-changing requests

### Sessions UI (Optional)

**What it generates:**
- View active sessions
- Revoke individual sessions
- Revoke all sessions except current
- Session details (IP, user agent, last active)

**Routes:**
- GET /auth/sessions - List all sessions
- POST /auth/sessions/revoke - Revoke session

## Wiring Auth Routes

**IMPORTANT:** Routes are NOT auto-injected. You must manually wire them in main.go.

**Why manual wiring?**
- Gives you control over which auth features to expose
- Allows custom route prefixes
- Enables middleware customization
- Prevents accidental exposure of features you don't want

**Example wiring (all features):**

```go
// In main.go, after database initialization

import (
    "myapp/internal/app/auth"
    "myapp/internal/shared/email"
)

// Configure email sender (required for magic links, email confirm, password reset)
emailSender := &email.ConsoleEmailSender{} // For development
// emailSender := email.NewSMTPSender(...) // For production

// Create auth handler
authHandler := auth.New(queries, emailSender)

// Register all auth routes
http.Handle("/auth/register", authHandler.HandleRegister())
http.Handle("/auth/login", authHandler.HandleLogin())
http.Handle("/auth/logout", authHandler.HandleLogout())

// Magic links (if enabled)
http.Handle("/auth/magic-link/request", authHandler.HandleMagicLinkRequest())
http.Handle("/auth/magic-link/verify", authHandler.HandleMagicLinkVerify())

// Email confirmation (if enabled)
http.Handle("/auth/confirm/send", authHandler.HandleConfirmationSend())
http.Handle("/auth/confirm/verify", authHandler.HandleConfirmationVerify())

// Password reset (if enabled)
http.Handle("/auth/password/reset/request", authHandler.HandlePasswordResetRequest())
http.Handle("/auth/password/reset/verify", authHandler.HandlePasswordResetVerify())

// Sessions UI (if enabled)
http.Handle("/auth/sessions", authHandler.HandleSessions())
http.Handle("/auth/sessions/revoke", authHandler.HandleSessionRevoke())

// Protect routes with middleware
http.Handle("/dashboard", auth.RequireAuth(queries, dashboardHandler))
http.Handle("/profile", auth.RequireAuth(queries, profileHandler))
```

**Example wiring (password-only, no email features):**

```go
// If you used: lvt gen auth --no-magic-link --no-email-confirm --no-password-reset

authHandler := auth.New(queries, nil) // No email sender needed

http.Handle("/auth/register", authHandler.HandleRegister())
http.Handle("/auth/login", authHandler.HandleLogin())
http.Handle("/auth/logout", authHandler.HandleLogout())

// Protect routes
http.Handle("/dashboard", auth.RequireAuth(queries, dashboardHandler))
```

## Email Configuration

**Development (Console Sender):**

```go
import "myapp/internal/shared/email"

emailSender := &email.ConsoleEmailSender{} // Prints to console
```

**Production (SMTP):**

```go
import "myapp/internal/shared/email"

emailSender := email.NewSMTPSender(email.SMTPConfig{
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your-email@gmail.com",
    Password: "your-app-password",
    From:     "noreply@yourapp.com",
})
```

**Custom Implementation:**

Implement the `email.Sender` interface:

```go
type Sender interface {
    Send(to, subject, body string) error
}
```

## LiveTemplate v0.5.1+ HTTP APIs

The generated authentication code uses LiveTemplate v0.5.1+ ActionContext HTTP methods for cookie and redirect operations. These APIs make it easy to handle authentication flows directly from your `Change()` handler.

### Available HTTP Methods

**Setting Cookies:**
```go
func (s *AuthState) handleLogin(ctx *livetemplate.ActionContext) error {
    // ... authenticate user ...

    // Set session cookie
    err := ctx.SetCookie(&http.Cookie{
        Name:     "session_token",
        Value:    token,
        Path:     "/",
        MaxAge:   30 * 24 * 60 * 60, // 30 days
        HttpOnly: true,
        Secure:   true, // Use true in production
        SameSite: http.SameSiteStrictMode,
    })
    if err != nil {
        return fmt.Errorf("failed to set cookie: %w", err)
    }

    return nil
}
```

**Reading Cookies:**
```go
func (s *AuthState) checkSession(ctx *livetemplate.ActionContext) error {
    cookie, err := ctx.GetCookie("session_token")
    if err != nil {
        // Cookie not found or error
        return err
    }

    // Verify token
    // ...
}
```

**Deleting Cookies:**
```go
func (s *AuthState) handleLogout(ctx *livetemplate.ActionContext) error {
    // Delete session cookie
    ctx.DeleteCookie("session_token")

    return ctx.Redirect("/", http.StatusSeeOther)
}
```

**Redirecting:**
```go
func (s *AuthState) handleLogin(ctx *livetemplate.ActionContext) error {
    // ... authenticate and set cookie ...

    // Redirect to dashboard
    return ctx.Redirect("/dashboard", http.StatusSeeOther)
}
```

**Headers:**
```go
// Set response header
ctx.SetHeader("X-Custom-Header", "value")

// Read request header
userAgent := ctx.GetHeader("User-Agent")
```

**Checking HTTP Context:**
```go
// Check if we're in an HTTP request context (vs WebSocket)
if ctx.IsHTTP() {
    // Can safely use SetCookie, Redirect, etc.
}
```

### Security Best Practices

**Cookie Security:**
- Always set `HttpOnly: true` to prevent JavaScript access
- Set `Secure: true` in production (requires HTTPS)
- Use `SameSite: http.SameSiteStrictMode` to prevent CSRF
- Set appropriate `MaxAge` for session duration

**Redirect Validation:**
- ActionContext validates redirects automatically
- Only relative paths starting with `/` are allowed
- Prevents open redirect vulnerabilities

**Error Handling:**
- `ctx.SetCookie()` returns `ErrNoHTTPContext` if called from WebSocket
- `ctx.Redirect()` returns `ErrInvalidRedirectCode` for invalid status codes
- `ctx.Redirect()` returns `ErrInvalidRedirectURL` for non-relative URLs

## Common Issues and Fixes

### Issue 1: "failed to load project config"

**Why it's wrong:** Not in an lvt project directory

**Fix:**
```bash
# Check if you're in the right directory
ls .lvtrc  # Should exist

# If missing, create a new app first
lvt new myapp
cd myapp
lvt gen auth
```

### Issue 2: "at least one authentication method must be enabled"

**Why it's wrong:** Used both --no-password and --no-magic-link

**Fix:**
```bash
# ❌ Wrong - disables both auth methods
lvt gen auth --no-password --no-magic-link

# ✅ Correct - keep at least one method
lvt gen auth --no-magic-link  # Password only
lvt gen auth --no-password    # Magic link only
lvt gen auth                  # Both (recommended)
```

### Issue 3: Build fails with "User undefined"

**Why it's wrong:** Forgot to run `sqlc generate` after migration

**Fix:**
```bash
# Run migration first
lvt migration up

# Generate sqlc models
cd internal/database
go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
cd ../..

# Run go mod tidy
go mod tidy

# Build should work now
go build ./cmd/myapp
```

### Issue 4: Auth routes return 404

**Why it's wrong:** Forgot to wire auth routes in main.go

**Fix:**
```go
// In main.go, add auth route wiring (see "Wiring Auth Routes" section)
authHandler := auth.New(queries, emailSender)
http.Handle("/auth/login", authHandler.HandleLogin())
// ... etc
```

### Issue 5: "email sender not configured"

**Why it's wrong:** Passed `nil` for email sender when using email features

**Fix:**
```go
// ❌ Wrong - nil email sender with magic links enabled
authHandler := auth.New(queries, nil)

// ✅ Correct - provide email sender
emailSender := &email.ConsoleEmailSender{} // Development
authHandler := auth.New(queries, emailSender)
```

### Issue 6: CSRF token validation fails

**Why it's wrong:** Missing CSRF token in forms or using GET for state changes

**Fix:**
```html
<!-- In your templates, include CSRF token -->
<form method="POST" action="/auth/login">
    <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
    <!-- other fields -->
</form>
```

```go
// In your handler, pass CSRF token to template
data := struct {
    CSRFToken string
}{
    CSRFToken: session.CSRFToken,
}
```

### Issue 7: Sessions expire too quickly/slowly

**Why it's wrong:** Using default expiration (30 days)

**Fix:**
```go
// Modify session expiration in internal/app/auth/auth.go
// Change the constant:
const SessionDuration = 7 * 24 * time.Hour // 7 days instead of 30
```

## Running E2E Tests

**Prerequisites:**
- Docker installed and running (for headless Chrome)
- Auth system generated
- Migrations applied

**Run tests:**
```bash
# Run all auth E2E tests
go test ./internal/app/auth -run TestAuthE2E -v

# Run specific test
go test ./internal/app/auth -run TestAuthE2E/Registration -v
go test ./internal/app/auth -run TestAuthE2E/Login -v
go test ./internal/app/auth -run TestAuthE2E/MagicLink -v
```

**What the tests verify:**
- ✓ User registration flow
- ✓ Login with correct credentials
- ✓ Login rejection with wrong credentials
- ✓ Magic link request and verification
- ✓ Email confirmation flow
- ✓ Password reset flow
- ✓ Session management (active sessions, revocation)
- ✓ Route protection (middleware)
- ✓ CSRF protection

## Success Response

After successful auth generation, provide:

```
✅ Authentication system generated successfully!

📁 Generated files:
  - internal/app/auth/auth.go
  - internal/app/auth/auth.tmpl
  - internal/app/auth/middleware.go
  - internal/app/auth/auth_e2e_test.go
  - internal/shared/password/password.go
  - internal/shared/email/email.go
  - internal/database/migrations/<timestamp>_create_auth_tables.sql

📝 Files updated:
  - internal/database/queries.sql

⚠️  IMPORTANT: Auth routes are NOT auto-injected

🚀 Next steps:
  1. Run migrations:
     lvt migration up

  2. Generate sqlc code:
     cd internal/database && go run github.com/sqlc-dev/sqlc/cmd/sqlc generate && cd ../..

  3. Wire auth routes in main.go (see skill for examples)

  4. Configure email sender (see internal/shared/email/email.go)

  5. Run E2E tests:
     go test ./internal/app/auth -run TestAuthE2E -v

💡 Tip: Check internal/app/auth/auth.go for complete usage examples!
```

## Common User Scenarios

**Scenario 1: Full authentication with all features**
- User: "Add authentication to my app"
- Command: `lvt gen auth`
- Features: Password, magic links, email confirm, password reset, sessions UI, CSRF
- Next: Wire all routes, configure email sender

**Scenario 2: Simple password-only auth**
- User: "Add password login, skip magic links and email"
- Command: `lvt gen auth --no-magic-link --no-email-confirm --no-password-reset`
- Features: Password auth only, sessions, CSRF
- Next: Wire basic routes (register, login, logout), no email needed

**Scenario 3: Passwordless (magic link only)**
- User: "I want magic link authentication only"
- Command: `lvt gen auth --no-password`
- Features: Magic links, email confirm, sessions, CSRF
- Next: Wire magic link routes, configure email sender

**Scenario 4: Custom struct/table names**
- User: "Add auth for admins instead of users"
- Command: `lvt gen auth Admin admin_users`
- Features: All features but with Admin struct and admin_users table
- Next: Same wiring, but use Admin instead of User in code

**Scenario 5: Auth for multi-tenant app**
- User: "Add auth for accounts table"
- Command: `lvt gen auth Account accounts`
- Features: All features with Account struct
- Next: Wire routes, potentially add tenant isolation logic

## Validation Criteria

Auth generation is successful if:
1. ✅ All files created without errors
2. ✅ Migration file created with correct tables
3. ✅ Queries appended to queries.sql
4. ✅ `lvt migration up` succeeds
5. ✅ `sqlc generate` succeeds
6. ✅ `go build` succeeds
7. ✅ No compilation errors
8. ✅ E2E tests can be run (docker required)

## Advanced Customization

After generation, users can customize:

**1. Session duration:**
```go
// In internal/app/auth/auth.go
const SessionDuration = 7 * 24 * time.Hour // Change from 30 days
```

**2. Password requirements:**
```go
// In internal/shared/password/password.go
const MinPasswordLength = 12 // Change from 8
```

**3. Token expiration:**
```go
// In internal/app/auth/auth.go
const MagicLinkExpiration = 5 * time.Minute // Change from 15
const ResetTokenExpiration = 30 * time.Minute // Change from 1 hour
```

**4. Add custom user fields:**
```sql
-- In migration file before running lvt migration up
ALTER TABLE users ADD COLUMN display_name TEXT;
ALTER TABLE users ADD COLUMN avatar_url TEXT;
```

**5. Add OAuth providers:**
- Extend auth handler with OAuth endpoints
- Add OAuth provider tables
- Implement OAuth flow in auth.go

**6. Add 2FA/MFA:**
- Generate TOTP secrets
- Add verification step after login
- Store backup codes

**7. Add role-based access control (RBAC):**
- Add roles table
- Add user_roles junction table
- Create role-checking middleware
- Extend RequireAuth to check roles

## Notes

- Auth routes are intentionally NOT auto-injected for security and flexibility
- Email sender must be configured for magic links, email confirm, and password reset
- CSRF protection is enabled by default - include token in all forms
- Sessions are stored in database (not cookies) for security
- Bcrypt cost is 12 by default (good balance of security and performance)
- E2E tests require Docker for headless Chrome
- Generated code includes extensive comments and examples
- Middleware can be applied to individual routes or route groups
- Sessions auto-expire after inactivity
- One user can have multiple active sessions (different devices)
