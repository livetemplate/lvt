---
name: lvt-production-ready
description: Transform app to production-ready - adds authentication, deployment config, environment setup, and best practices
category: workflows
version: 1.0.0
keywords: ["lvt", "livetemplate", "lt"]
---

# lvt:production-ready

Transform a development app into a production-ready application with authentication, deployment configuration, environment variables, and production best practices.

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
- "Make my app production ready"
- "Add authentication and deployment"
- "Prepare my app for production"
- "I need to deploy this app"
- "Make this ready for real users"

**Examples:**
- "Make my blog production ready"
- "Add auth and prepare for deployment"
- "I want to deploy this to production"
- "Make this app ready for real users"

## Workflow Steps

This skill chains together:
1. **lvt:gen-auth** - Add authentication system
2. **lvt:deploy** - Add deployment configuration
3. Environment setup guidance
4. Production best practices checklist

### Step 1: Add Authentication

Use **lvt:gen-auth** skill:

**Ask user about auth requirements:**
- Need password auth? (default: yes)
- Need magic link auth? (default: yes)
- Need email confirmation? (default: yes)
- Need password reset? (default: yes)
- Need sessions UI? (default: yes)

**Generate auth:**
```bash
lvt gen auth

# Or with specific features:
lvt gen auth --no-magic-link  # Password only
lvt gen auth --no-password    # Magic link only
```

**Apply migrations:**
```bash
lvt migration up
cd internal/database && sqlc generate && cd ../..
go mod tidy
```

**Wire auth routes in main.go:**
```go
import (
    "myapp/internal/app/auth"
    "myapp/internal/shared/email"
)

// Email sender (production)
emailSender := email.NewSMTPSender(email.SMTPConfig{
    Host:     os.Getenv("SMTP_HOST"),
    Port:     587,
    Username: os.Getenv("SMTP_USER"),
    Password: os.Getenv("SMTP_PASS"),
    From:     os.Getenv("EMAIL_FROM"),
})

// Auth handler
authHandler := auth.New(queries, emailSender)

// Auth routes
http.Handle("/auth/register", authHandler.HandleRegister())
http.Handle("/auth/login", authHandler.HandleLogin())
http.Handle("/auth/logout", authHandler.HandleLogout())

// Protect existing routes
http.Handle("/posts", auth.RequireAuth(queries, posts.Handler(queries)))
http.Handle("/dashboard", auth.RequireAuth(queries, dashboard.Handler()))
```

### Step 2: Add Deployment Configuration

Use **lvt:deploy** skill:

**Ask user about deployment target:**
- Docker (for any platform)
- Fly.io (optimized for SQLite)
- Kubernetes (for scale)
- Traditional VPS

**Generate deployment files:**
```bash
lvt gen stack docker      # For Docker
lvt gen stack fly         # For Fly.io
lvt gen stack k8s         # For Kubernetes
```

**Verify deployment config created:**
- Dockerfile
- docker-compose.yml (if Docker)
- fly.toml (if Fly.io)
- k8s manifests (if Kubernetes)
- .env.example
- README.md with deployment instructions

### Step 3: Environment Variables Setup

Create `.env` file (never commit this):
```bash
# Database
DATABASE_PATH=/data/app.db

# Server
PORT=8080
HOST=0.0.0.0

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
EMAIL_FROM=noreply@yourapp.com

# Application
APP_ENV=production
APP_URL=https://yourapp.com

# Secrets
SESSION_SECRET=<generate-random-32-bytes>
CSRF_SECRET=<generate-random-32-bytes>
```

**Generate secrets:**
```bash
# Generate random secrets
openssl rand -hex 32  # For SESSION_SECRET
openssl rand -hex 32  # For CSRF_SECRET
```

### Step 4: Production Best Practices Checklist

#### Security ✅
- [ ] All routes using HTTPS
- [ ] CSRF protection enabled
- [ ] Session secrets are random and secure
- [ ] Password reset tokens expire
- [ ] Email confirmation required (if applicable)
- [ ] Rate limiting on auth endpoints
- [ ] SQL injection prevention (using sqlc ✓)
- [ ] XSS prevention (html/template escaping ✓)

#### Database ✅
- [ ] Migrations applied in correct order
- [ ] Database backups configured
- [ ] Database connection pooling set up
- [ ] Indexes on frequently queried columns
- [ ] Litestream configured (for SQLite)

#### Monitoring ✅
- [ ] Error logging configured
- [ ] Access logs enabled
- [ ] Health check endpoint added
- [ ] Metrics collection (optional)

#### Performance ✅
- [ ] Static assets cached
- [ ] Gzip compression enabled
- [ ] Database queries optimized
- [ ] Connection timeouts configured

#### Deployment ✅
- [ ] Environment variables documented
- [ ] Deployment instructions in README
- [ ] CI/CD pipeline set up (optional)
- [ ] Rollback strategy documented

### Step 5: Add Health Check Endpoint

Add to main.go:
```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteString(`{"status":"ok"}`)
})
```

### Step 6: Update README

Add production deployment section:
```markdown
## Production Deployment

### Prerequisites
- Docker installed (for Docker deployment)
- Environment variables configured (see .env.example)
- SMTP credentials for email

### Deploy to [Platform]
1. Copy .env.example to .env and fill in values
2. Generate secrets: `openssl rand -hex 32`
3. Build: `docker build -t myapp .`
4. Deploy: [platform-specific commands]

### Environment Variables
See .env.example for all required variables.

### Database
Run migrations before first deploy:
```bash
lvt migration up
```

### Email Configuration
Configure SMTP settings in .env for production email delivery.
```

## Quick Reference

### Full Production Setup (Docker + Auth)
```bash
# 1. Add authentication
lvt gen auth
lvt migration up
cd internal/database && sqlc generate && cd ../..
go mod tidy

# 2. Add deployment config
lvt gen stack docker

# 3. Create environment file
cp .env.example .env
# Edit .env with production values

# 4. Generate secrets
echo "SESSION_SECRET=$(openssl rand -hex 32)" >> .env
echo "CSRF_SECRET=$(openssl rand -hex 32)" >> .env

# 5. Build and test locally
docker-compose up --build

# 6. Deploy
docker-compose up -d
```

### Fly.io Deployment
```bash
# 1. Add auth
lvt gen auth
lvt migration up
cd internal/database && sqlc generate && cd ../..

# 2. Generate Fly.io config
lvt gen stack fly

# 3. Deploy
fly launch
fly deploy
```

## Checklist

- [ ] Verify app is in working state before starting
- [ ] Ask user about auth requirements
- [ ] Use lvt:gen-auth to add authentication
- [ ] Wire auth routes in main.go
- [ ] Protect existing routes with RequireAuth
- [ ] Ask user about deployment target
- [ ] Use lvt:deploy to generate deployment config
- [ ] Create .env.example with all variables
- [ ] Generate secure random secrets
- [ ] Add health check endpoint
- [ ] Review security checklist
- [ ] Update README with deployment instructions
- [ ] Test deployment locally (Docker)
- [ ] Guide user through actual deployment

## Platform-Specific Guides

### Docker
**Best for:** Any cloud provider, VPS, local hosting
**Pros:** Universal, portable, reproducible
**Setup:** Dockerfile + docker-compose.yml

**Deploy:**
```bash
docker-compose up -d
```

### Fly.io
**Best for:** SQLite apps, global edge deployment
**Pros:** Optimized for SQLite, automatic HTTPS, global CDN
**Setup:** fly.toml + litestream.yml

**Deploy:**
```bash
fly deploy
```

### Kubernetes
**Best for:** Large scale, enterprise deployments
**Pros:** Highly scalable, self-healing, enterprise-grade
**Setup:** Deployment, Service, Ingress manifests

**Deploy:**
```bash
kubectl apply -f k8s/
```

### VPS (Traditional)
**Best for:** Simple deployments, single server
**Pros:** Simple, cost-effective
**Setup:** systemd service + nginx

**Deploy:**
```bash
# Build binary
GOWORK=off CGO_ENABLED=1 go build -o myapp ./cmd/myapp
# Copy to server and run
./myapp
```

## Success Criteria

Production-ready setup is successful when:
1. ✅ Authentication system working
2. ✅ All auth routes wired and tested
3. ✅ Protected routes require authentication
4. ✅ Deployment config generated
5. ✅ Environment variables documented
6. ✅ Health check endpoint responding
7. ✅ App builds in production mode
8. ✅ Deployment tested locally
9. ✅ README has deployment instructions

## Common Post-Production Tasks

### Add Custom Domain
**Fly.io:**
```bash
fly certs add yourdomain.com
```

**Docker/VPS:**
Configure nginx or Caddy reverse proxy

### Set Up SSL/TLS
**Fly.io:** Automatic
**Docker:** Use Let's Encrypt + Caddy/nginx
**K8s:** Use cert-manager

### Configure Monitoring
- Add logging service (Sentry, LogRocket)
- Set up uptime monitoring (UptimeRobot)
- Configure alerts

### Database Backups
**SQLite + Litestream:**
```yaml
# litestream.yml already configured
# Backups to S3/B2/Azure automatically
```

**PostgreSQL:**
Set up pg_dump cron job

## Notes

- This is a meta-skill that chains multiple skills
- Always test locally before production deploy
- Keep secrets out of version control
- Document all environment variables
- Provide rollback instructions
- Test auth flows before deploying
- Verify email delivery in production
- Monitor error rates after deployment
