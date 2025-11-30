---
name: lvt-manage-env
description: Manage environment variables - set, unset, list, and validate configuration required for generated app features
category: core
version: 1.0.0
keywords: ["lvt", "livetemplate", "lt"]
---

# lvt:manage-env

Manage environment variables for your LiveTemplate application. Helps detect required configuration based on features, guide users to set values, validate configuration, and ensure apps are ready to run.

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
✅ Generic prompts related to this skill's purpose

**Without Context (needs keywords):**
✅ Must mention "lvt", "livetemplate", or "lt"
❌ Generic requests without keywords

---

## User Prompts

**When to use:**
- "Set up environment variables"
- "Configure my app environment"
- "What env vars do I need?"
- "Validate my environment"
- "Check if my app is configured correctly"
- "Add SMTP configuration"
- "Set up database connection"

**Examples:**
- "Configure environment for production"
- "What environment variables are missing?"
- "Set up email for my app"
- "Validate my .env file"
- "Check if I'm ready to deploy"

## Context Awareness

Before executing this skill, verify:

1. **In Project Directory:**
   - Check for `.lvtrc` file (confirms it's an lvt project)
   - Check for `go.mod` (confirms it's a Go project)

2. **Understand User Intent:**
   - Want to see what's configured? → Use `lvt env list`
   - Want to validate? → Use `lvt env validate`
   - Want to set values? → Use `lvt env set`
   - Want to remove values? → Use `lvt env unset`
   - Want to generate template? → Use `lvt env generate`

## Commands Available

### Generate .env.example Template
```bash
lvt env generate
```

**What it does:**
- Detects features in your app (auth, database, email, etc.)
- Generates .env.example with appropriate variables
- Includes helpful comments and security tips
- Feature-aware (only includes vars for features you use)

### Set Environment Variable
```bash
lvt env set KEY VALUE
```

**Examples:**
```bash
lvt env set APP_ENV development
lvt env set DATABASE_PATH ./app.db
lvt env set SESSION_SECRET $(openssl rand -hex 32)
lvt env set SMTP_HOST smtp.gmail.com
```

**What it does:**
- Sets or updates environment variable in .env
- Validates key format (UPPERCASE_SNAKE_CASE)
- Masks sensitive values in output
- Auto-adds .env to .gitignore
- Creates .env if it doesn't exist

### Unset Environment Variable
```bash
lvt env unset KEY
```

**Example:**
```bash
lvt env unset SMTP_HOST
```

**What it does:**
- Removes environment variable from .env
- Confirms the variable exists before removal

### List Environment Variables
```bash
lvt env list                # Masked values
lvt env list --show-values  # Show actual values
lvt env list --required-only # Only required vars
```

**What it does:**
- Shows all variables from .env
- Marks required variables with [REQUIRED]
- Masks sensitive values by default
- Feature-aware (knows what's required based on your app)

### Validate Environment
```bash
lvt env validate         # Check required vars are set
lvt env validate --strict # Also validate values
```

**What it does:**
- Checks all required variables are set
- Detects placeholder values that need replacement
- Provides helpful error messages with reasons
- Strict mode validates:
  - APP_ENV is valid (development/staging/production)
  - Secrets are strong enough (32+ chars)
  - EMAIL_PROVIDER is valid
  - SMTP vars present when needed
  - PORT is numeric

## Feature Detection

The commands automatically detect what features your app uses:

**Database:**
- Checks for `internal/database/schema.sql`
- Checks for `internal/database/migrations/`
- Requires: `DATABASE_PATH`

**Auth:**
- Checks for `internal/app/auth/` directory
- Requires: `SESSION_SECRET`, `CSRF_SECRET`

**Email (Auth with Email Features):**
- Checks auth files for magic link/confirmation/reset
- Requires: `EMAIL_PROVIDER`
- If `EMAIL_PROVIDER=smtp`, also requires:
  - `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`

**Server (Always):**
- Requires: `APP_ENV`

## Checklist

- [ ] **Step 1:** Understand what user wants to do
  - Configure environment? → Set variables
  - Check status? → List variables
  - Verify readiness? → Validate
  - Generate template? → Generate .env.example

- [ ] **Step 2:** Detect current state
  - Run `lvt env list` to see what's already set
  - Run `lvt env validate` to check for missing vars

- [ ] **Step 3:** Guide user based on validation results
  - If missing vars, explain what each one is for
  - Suggest commands to set them
  - Provide examples with actual values when safe

- [ ] **Step 4:** Help set values
  - For secrets: Suggest `openssl rand -hex 32`
  - For EMAIL_PROVIDER: Explain options (console, smtp)
  - For APP_ENV: Explain options (development, staging, production)
  - For DATABASE_PATH: Suggest appropriate path for environment

- [ ] **Step 5:** Validate again after changes
  - Run `lvt env validate`
  - For production: Use `lvt env validate --strict`

## Common Workflows

### Workflow 1: Initial Setup

User: "Set up environment variables"

**Steps:**
1. Run `lvt env generate` to create .env.example
2. Copy to .env: `cp .env.example .env`
3. Run `lvt env validate` to see what needs to be set
4. Guide user to set each required variable:
   ```bash
   lvt env set APP_ENV development
   lvt env set DATABASE_PATH ./dev.db
   lvt env set SESSION_SECRET $(openssl rand -hex 32)
   # ... etc
   ```
5. Run `lvt env validate` again to confirm all set

### Workflow 2: Check Configuration

User: "What environment variables are missing?"

**Steps:**
1. Run `lvt env list` to show current state
2. Run `lvt env validate` to identify missing vars
3. Explain what each missing var is for
4. Provide set commands

### Workflow 3: Production Setup

User: "Configure for production deployment"

**Steps:**
1. Run `lvt env validate --strict` to ensure all values are production-ready
2. Check for placeholder values
3. Ensure secrets are strong (32+ chars)
4. Verify APP_ENV=production
5. Check EMAIL_PROVIDER is set correctly
6. If SMTP, verify all SMTP credentials are set

### Workflow 4: Email Setup

User: "Set up email for my app"

**Steps:**
1. Check if email features are present (run detection)
2. If not present, explain they need `lvt gen auth` first
3. If present, ask for email provider:
   - Development: `lvt env set EMAIL_PROVIDER console`
   - Production: `lvt env set EMAIL_PROVIDER smtp`
4. If SMTP, set all required vars:
   ```bash
   lvt env set SMTP_HOST smtp.gmail.com
   lvt env set SMTP_PORT 587
   lvt env set SMTP_USER your-email@gmail.com
   lvt env set SMTP_PASS your-app-password
   ```
5. Validate with `lvt env validate --strict`

## Error Handling

**Error: "subcommand required"**
- Missing subcommand
- Show available commands
- Suggest common usage

**Error: "invalid key format"**
- Key must be UPPERCASE_SNAKE_CASE
- Provide corrected example

**Error: ".env file not found"**
- For list/unset/validate commands
- Suggest: `lvt env generate && cp .env.example .env`

**Error: "validation failed"**
- Shows specific missing or invalid variables
- Explains why each is required
- Provides fix commands

**Error: "not in a LiveTemplate app directory"**
- Not in project root
- Need to cd to project directory first

## Success Response Templates

### After Setting Variable

```
✅ Set SESSION_SECRET=a3f4****

Tip: This value is masked for security. Use --show-values to see actual values.
```

### After Validation Success

```
✅ All required environment variables are set

Your app is configured and ready to run!

Next steps:
  - Start dev server: lvt serve
  - Run tests: go test ./...
  - Deploy: see lvt:deploy skill
```

### After Validation with Issues

```
❌ Environment validation failed

Missing required variables:
  - SESSION_SECRET (session security - auth enabled)
  - EMAIL_PROVIDER (email functionality - auth with email features)

Set them with:
  lvt env set SESSION_SECRET $(openssl rand -hex 32)
  lvt env set EMAIL_PROVIDER console

Invalid placeholder values (need to be replaced):
  - CSRF_SECRET=change-me-to-random-32-byte-hex

Fix with:
  lvt env set CSRF_SECRET $(openssl rand -hex 32)

Fix these issues, then run 'lvt env validate' again
```

## Security Best Practices

Always remind users:

1. **Never commit .env** - It contains secrets
   - .env is auto-added to .gitignore
   - .env.example is safe to commit (no actual values)

2. **Generate strong secrets**
   ```bash
   openssl rand -hex 32
   ```

3. **Use different secrets per environment**
   - Don't reuse development secrets in production

4. **File permissions**
   - .env is created with 0600 (owner read/write only)

5. **Rotate secrets regularly**
   - Change SESSION_SECRET, CSRF_SECRET periodically in production

6. **Use environment-specific files**
   - .env.development (local dev)
   - .env.staging (staging server)
   - .env.production (production - use secrets manager)

## Validation Rules

### APP_ENV
- Must be: `development`, `staging`, or `production`
- Required: Always

### SESSION_SECRET / CSRF_SECRET
- Minimum length: 32 characters (strict mode)
- Cannot contain: Placeholder values
- Generate with: `openssl rand -hex 32`
- Required: If auth is enabled

### DATABASE_PATH
- Required: If database features present
- Common values:
  - Development: `./dev.db`
  - Production: `/data/app.db` or persistent volume

### EMAIL_PROVIDER
- Must be: `console` or `smtp`
- Required: If email features present (magic link, confirmation, reset)
- `console` = prints emails to terminal (dev only)
- `smtp` = sends real emails (requires SMTP_* vars)

### SMTP Variables
- Required when: `EMAIL_PROVIDER=smtp`
- Must set all:
  - `SMTP_HOST` (e.g., smtp.gmail.com)
  - `SMTP_PORT` (e.g., 587)
  - `SMTP_USER` (email address)
  - `SMTP_PASS` (app password)

### PORT
- Must be: Numeric value
- Default: 3000 or 8080
- Production: Usually set by platform (Fly.io, Railway, etc.)

## Example Interactions

### Example 1: Initial Setup

**User:** "Set up environment variables for my app"

**Skill:**
1. Runs: `lvt env generate`
2. Creates .env.example with detected features
3. Runs: `cp .env.example .env`
4. Runs: `lvt env validate`
5. Shows missing variables with explanations
6. Guides user to set each one:
   ```bash
   lvt env set APP_ENV development
   lvt env set DATABASE_PATH ./dev.db
   lvt env set SESSION_SECRET $(openssl rand -hex 32)
   lvt env set CSRF_SECRET $(openssl rand -hex 32)
   ```
7. Runs: `lvt env validate` again
8. Confirms: "✅ All required environment variables are set"

### Example 2: Check Status

**User:** "What environment variables are configured?"

**Skill:**
1. Runs: `lvt env list`
2. Shows all variables (masked)
3. Highlights required ones
4. Shows if any are missing

### Example 3: Production Validation

**User:** "Am I ready to deploy?"

**Skill:**
1. Runs: `lvt env validate --strict`
2. Checks all required vars
3. Validates values are production-ready
4. Checks for placeholder values
5. Verifies secrets are strong
6. Reports: Ready ✅ or issues to fix ❌

### Example 4: Email Configuration

**User:** "Add SMTP configuration"

**Skill:**
1. Asks: "Which SMTP provider? (Gmail, SendGrid, Mailgun, custom)"
2. Based on choice, provides specific guidance:

   **For Gmail:**
   ```bash
   lvt env set SMTP_HOST smtp.gmail.com
   lvt env set SMTP_PORT 587
   lvt env set SMTP_USER your-email@gmail.com
   lvt env set SMTP_PASS your-app-password
   lvt env set EMAIL_FROM noreply@yourdomain.com
   ```

   Explains: "For Gmail, you need an App Password (not your regular password)"
   Link: https://support.google.com/accounts/answer/185833

3. Runs: `lvt env validate --strict`
4. Suggests testing: "Try sending a test email with your app"

## Tips for Claude

- **Be proactive**: Run `lvt env validate` automatically to check status
- **Be specific**: Don't just say "set SESSION_SECRET", show the exact command with value generation
- **Be secure**: Always remind about security best practices
- **Be helpful**: Explain WHY each variable is needed, not just WHAT it is
- **Be environment-aware**: Different advice for development vs production

## Notes

- .env file is created with restricted permissions (0600)
- Sensitive values are automatically masked in output
- Commands are idempotent (safe to run multiple times)
- Feature detection is automatic (no config needed)
- .gitignore is updated automatically
- Works with existing .env files (doesn't overwrite without confirmation)
