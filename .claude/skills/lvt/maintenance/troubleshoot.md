---
name: lvt:troubleshoot
description: Debug common LiveTemplate issues - build errors, migration problems, template errors, auth issues, deployment failures
category: maintenance
version: 1.0.0
---

# lvt:troubleshoot

Systematic debugging guide for common LiveTemplate issues. Helps diagnose and fix build errors, migration problems, template errors, authentication issues, and deployment failures.

## User Prompts

**When to use:**
- "My app won't build"
- "I'm getting an error"
- "Something's broken"
- "Migrations aren't working"
- "Templates are failing"
- "Authentication doesn't work"
- "Deployment failed"

**Examples:**
- "Help me fix this build error"
- "Why isn't my migration working?"
- "My template is throwing errors"
- "Users can't log in"
- "Server won't start"

## Troubleshooting Framework

### Step 1: Identify Issue Category

**Common categories:**
1. Build/compilation errors
2. Migration issues
3. Template errors
4. Authentication problems
5. WebSocket/reactivity issues
6. Deployment failures
7. Runtime errors

### Step 2: Gather Context

**Always collect:**
- Error messages (full output)
- What command was run
- What was expected vs. what happened
- Recent changes made
- Go version, OS, environment

### Step 3: Apply Solution Pattern

Each category has known solutions and patterns.

## Issue Categories

### 1. Build Errors

#### Error: "undefined: queries"

**Symptom:**
```
internal/app/posts/posts.go:25:15: undefined: queries
```

**Cause:** sqlc models not generated or outdated

**Solution:**
```bash
cd internal/database
sqlc generate
cd ../..
go mod tidy
```

**Prevention:** Run sqlc generate after any schema changes

---

#### Error: "package X is not in GOROOT"

**Symptom:**
```
package myapp/internal/app/auth is not in GOROOT
```

**Cause:** Missing dependency or incorrect module name

**Solution:**
```bash
# Check go.mod module name
cat go.mod | head -1

# Ensure it matches import paths
# If mismatch, update go.mod:
go mod edit -module <correct-name>
go mod tidy
```

**Prevention:** Use correct module name from project creation

---

#### Error: "CGO_ENABLED required"

**Symptom:**
```
# github.com/mattn/go-sqlite3
exec: "gcc": executable file not found in $PATH
```

**Cause:** SQLite requires CGO, gcc not available

**Solution:**
```bash
# macOS
xcode-select --install

# Linux (Ubuntu/Debian)
sudo apt-get install build-essential

# Linux (Alpine)
apk add gcc musl-dev

# Verify
CGO_ENABLED=1 go build
```

**Prevention:** Install build tools before using SQLite

---

#### Error: "imports cycle"

**Symptom:**
```
import cycle not allowed
package myapp/internal/app/posts
	imports myapp/internal/app/auth
	imports myapp/internal/app/posts
```

**Cause:** Circular dependencies between packages

**Solution:**
- Move shared code to internal/shared/
- Use interfaces to break cycles
- Restructure package dependencies

**Prevention:** Keep packages independent, share via interfaces

### 2. Migration Issues

#### Error: "migration already applied"

**Symptom:**
```
Error: version 20251104120000 already applied
```

**Cause:** Attempting to re-run applied migration

**Solution:**
```bash
# Check migration status
lvt migration status

# If stuck, rollback and retry
lvt migration down
lvt migration up
```

**Prevention:** Check status before running migrations

---

#### Error: "no such table"

**Symptom:**
```
Error: no such table: posts
```

**Cause:** Migrations not applied or database path incorrect

**Solution:**
```bash
# Verify database path
echo $DATABASE_PATH
# or check .env

# Apply migrations
lvt migration up

# Verify tables
sqlite3 dev.db ".tables"
```

**Prevention:** Always run migrations after schema changes

---

#### Error: "migration failed with syntax error"

**Symptom:**
```
Error 1: near "SERIAL": syntax error
```

**Cause:** PostgreSQL syntax in SQLite migration

**Solution:**
```sql
-- WRONG (PostgreSQL)
CREATE TABLE posts (
    id SERIAL PRIMARY KEY
);

-- RIGHT (SQLite)
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT
);
```

**Prevention:** Use SQLite syntax (INTEGER PRIMARY KEY)

---

#### Error: "foreign key constraint failed"

**Symptom:**
```
Error: FOREIGN KEY constraint failed
```

**Cause:** Foreign keys not enabled or referencing missing row

**Solution:**
```bash
# Check if foreign keys enabled
sqlite3 dev.db "PRAGMA foreign_keys;"
# Should return: 1

# If 0, enable in main.go:
# db.Exec("PRAGMA foreign_keys = ON")

# Check referenced table
sqlite3 dev.db "SELECT * FROM parent_table WHERE id = X;"
```

**Prevention:** Enable foreign keys, ensure parent rows exist

### 3. Template Errors

#### Error: "template: undefined function"

**Symptom:**
```
Error: template: posts.tmpl:15: function "formatDate" not defined
```

**Cause:** Template function not registered

**Solution:**
```go
// In handler, add custom funcs
funcMap := template.FuncMap{
    "formatDate": func(t time.Time) string {
        return t.Format("Jan 2, 2006")
    },
}
tmpl := template.New("posts.tmpl").Funcs(funcMap)
```

**Prevention:** Register all custom functions before parsing

---

#### Error: "template: unexpected EOF"

**Symptom:**
```
Error: template: posts.tmpl:45: unexpected EOF
```

**Cause:** Unclosed template action ({{ }})

**Solution:**
```bash
# Use lvt parse to find error
lvt parse internal/app/posts/posts.tmpl

# Common issues:
{{ if .Items }}  # Missing end
{{ range .Items  # Missing }}
```

**Prevention:** Use lvt parse before running app

---

#### Error: "can't evaluate field X"

**Symptom:**
```
Error: template: posts.tmpl:20: can't evaluate field Title in type *models.Post
```

**Cause:** Field doesn't exist or wrong type passed

**Solution:**
```bash
# Check database model
cat internal/database/models.go | grep "type Post"

# Verify field names match (case-sensitive)
# Title vs title

# Check what's passed to template
# Should be: tmpl.Execute(w, posts)
# Not: tmpl.Execute(w, wrongType)
```

**Prevention:** Match template field names to model fields exactly

### 4. Authentication Issues

#### Error: "invalid CSRF token"

**Symptom:**
- Form submissions fail with 403
- Error: "invalid CSRF token"

**Cause:** Missing or expired CSRF token

**Solution:**
```html
<!-- Add to all forms -->
<form method="POST">
    {{ csrfField }}
    <!-- form fields -->
</form>
```

**Prevention:** Always include {{ csrfField }} in POST forms

---

#### Error: "session not found"

**Symptom:**
- User logged in but shown as logged out
- Session disappears between requests

**Cause:** Session configuration issue or database problem

**Solution:**
```bash
# Check sessions table
sqlite3 dev.db "SELECT * FROM sessions;"

# Verify session middleware is used
# main.go should have:
# http.Handle("/", sessionMiddleware(handler))

# Check cookie settings
# Domain, Secure, HttpOnly, SameSite
```

**Prevention:** Test auth flows with browser dev tools open

---

#### Error: "password doesn't match"

**Symptom:**
- Correct password rejected
- Login always fails

**Cause:** Password not hashed correctly

**Solution:**
```bash
# Check password storage in database
sqlite3 dev.db "SELECT password_hash FROM users LIMIT 1;"
# Should be bcrypt hash (starts with $2a$ or $2b$)

# Verify bcrypt comparison
# golang.org/x/crypto/bcrypt
# bcrypt.CompareHashAndPassword(hashedPassword, password)
```

**Prevention:** Use bcrypt for password hashing, never plain text

---

#### Error: "email not sent"

**Symptom:**
- Magic link/reset emails not received
- No errors shown

**Cause:** Email sender not configured

**Solution:**
```bash
# Check email configuration
# In main.go:
emailSender := &email.ConsoleEmailSender{}  # For development
# or
emailSender := email.NewSMTPSender(config)  # For production

# For development, check console output
# For production, verify SMTP settings:
echo $SMTP_HOST
echo $SMTP_USER
# (Don't echo SMTP_PASS for security)
```

**Prevention:** Use ConsoleEmailSender for development

### 5. WebSocket/Reactivity Issues

#### Error: "WebSocket connection failed"

**Symptom:**
- Browser console: "WebSocket connection to 'ws://...' failed"
- No live updates

**Cause:** WebSocket endpoint not accessible or CORS issue

**Solution:**
```bash
# Check WebSocket endpoint
curl -i -N -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     http://localhost:3000/ws

# Should return 101 Switching Protocols

# Check main.go has WebSocket handler
# http.Handle("/ws", websocketHandler)

# For HTTPS, use wss:// not ws://
```

**Prevention:** Test WebSocket in browser dev tools Network tab

---

#### Error: "actions not triggering"

**Symptom:**
- Clicking buttons doesn't update UI
- No errors in console

**Cause:** Action not registered or wrong event binding

**Solution:**
```html
<!-- Check lf-action syntax -->
<button lf-action="delete" lf-target="post-{{.ID}}">Delete</button>

<!-- Verify action handler exists -->
<!-- In posts.go: -->
// case "delete":
//     return handleDelete(w, r, queries)
```

**Prevention:** Use browser dev tools to inspect WebSocket messages

---

#### Error: "template not found for partial"

**Symptom:**
```
Error: template "post-item" not found
```

**Cause:** Partial template not defined

**Solution:**
```html
<!-- Define partial in template -->
{{ define "post-item" }}
<div id="post-{{.ID}}">
    <!-- content -->
</div>
{{ end }}

<!-- Or extract to separate file and embed -->
```

**Prevention:** Define all partials used in WebSocket responses

### 6. Deployment Issues

#### Error: "port already in use"

**Symptom:**
```
Error: listen tcp :3000: bind: address already in use
```

**Cause:** Another process using the port

**Solution:**
```bash
# Find process using port
lsof -i :3000

# Kill process
kill -9 <PID>

# Or use different port
PORT=3001 lvt serve
```

**Prevention:** Stop previous server before starting new one

---

#### Error: "database locked"

**Symptom:**
```
Error: database is locked
```

**Cause:** Multiple connections writing simultaneously (SQLite limitation)

**Solution:**
```go
// Set connection pool for SQLite
db.SetMaxOpenConns(1)

// Or use WAL mode
db.Exec("PRAGMA journal_mode=WAL;")
```

**Prevention:** Configure SQLite for concurrent access

---

#### Error: "static files not found"

**Symptom:**
- CSS/JS not loading
- 404 for /static/ paths

**Cause:** Static file serving not configured

**Solution:**
```bash
# Verify static directory exists
ls static/

# Check main.go has file server
# http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

# For embedded FS:
// go:embed static
// var staticFS embed.FS
```

**Prevention:** Test static files before deployment

---

#### Error: "cannot find package in Docker"

**Symptom:**
```
Error: cannot find package "myapp/internal/app/posts"
```

**Cause:** Files not copied to Docker image

**Solution:**
```dockerfile
# Dockerfile must copy all source
COPY go.mod go.sum ./
RUN go mod download
COPY . .  # Copy entire project
```

**Prevention:** Test Docker build locally before deploying

### 7. Runtime Errors

#### Error: "panic: runtime error: invalid memory address"

**Symptom:**
```
panic: runtime error: invalid memory address or nil pointer dereference
```

**Cause:** Accessing nil pointer

**Solution:**
```go
// Always check for nil
if post == nil {
    http.Error(w, "Post not found", http.StatusNotFound)
    return
}

// Use sql.ErrNoRows check
post, err := queries.GetPost(ctx, id)
if err == sql.ErrNoRows {
    http.Error(w, "Not found", http.StatusNotFound)
    return
}
```

**Prevention:** Always check errors and nil values

---

#### Error: "too many open files"

**Symptom:**
```
Error: accept tcp [::]:3000: accept4: too many open files
```

**Cause:** File descriptor limit reached

**Solution:**
```bash
# Check limit
ulimit -n

# Increase limit (macOS/Linux)
ulimit -n 4096

# Or in code, close connections properly
defer conn.Close()
defer file.Close()
```

**Prevention:** Always close resources (defer)

## Debugging Checklist

When troubleshooting any issue:

- [ ] Read the full error message carefully
- [ ] Check what command was run
- [ ] Verify environment variables
- [ ] Check recent changes (git diff)
- [ ] Look for typos in code/config
- [ ] Test in isolation (minimal reproduction)
- [ ] Check logs (server + browser console)
- [ ] Verify database state
- [ ] Try the simplest solution first
- [ ] Document the solution for next time

## Diagnostic Commands

```bash
# Check app structure
ls -la

# Verify Go version
go version

# Check module
cat go.mod

# List database tables
sqlite3 dev.db ".tables"

# Check migration status
lvt migration status

# Validate template
lvt parse internal/app/posts/posts.tmpl

# Test build
go build -v

# Run tests
go test ./...

# Check running processes
ps aux | grep myapp

# Check port usage
lsof -i :3000

# View logs
tail -f /var/log/myapp.log
```

## Common Quick Fixes

**Problem:** "It was working before"
**Solution:** `git diff` to see what changed, revert if needed

**Problem:** "Weird caching behavior"
**Solution:** Clear browser cache, restart server

**Problem:** "Works locally but not in production"
**Solution:** Check environment variables, database path, file permissions

**Problem:** "Random intermittent errors"
**Solution:** Check for race conditions, add logging, use `go run -race`

**Problem:** "Nothing works after update"
**Solution:** `go mod tidy`, `sqlc generate`, restart server

## Prevention Best Practices

1. **Use lvt parse** before running app (catch template errors early)
2. **Run migrations** after every schema change
3. **Check logs** regularly (server + browser)
4. **Test auth flows** in incognito mode
5. **Use version control** (commit often, small changes)
6. **Read error messages** carefully (full output)
7. **Test locally** before deploying
8. **Keep dependencies updated** (`go mod tidy`)
9. **Use linters** (`golangci-lint run`)
10. **Monitor in production** (error tracking, logs)

## Getting Help

If troubleshooting doesn't resolve the issue:

1. **Search error message** - Exact error text often has solutions
2. **Check LiveTemplate docs** - May have specific guidance
3. **Simplify reproduction** - Minimal example that shows the issue
4. **Collect details:**
   - Full error output
   - Go version (`go version`)
   - OS/environment
   - Steps to reproduce
   - What was tried already
5. **Ask for help** with all details above

## Success Criteria

Issue is resolved when:
1. ✅ Error no longer occurs
2. ✅ Expected behavior works
3. ✅ Root cause identified
4. ✅ Solution documented
5. ✅ Prevention steps in place
6. ✅ Tests added (if applicable)

## Notes

- Most issues have simple solutions (typos, missed steps)
- Always check the basics first (migrations, sqlc, env vars)
- Use systematic debugging (isolate, test, verify)
- Document solutions for future reference
- Prevention is better than debugging
- When stuck, start fresh (new terminal, restart server)
- Browser dev tools are your friend (console, network, application)
- SQLite issues often fixed with WAL mode
- Template errors caught early with lvt parse
- Most deployment issues are environment/config related
