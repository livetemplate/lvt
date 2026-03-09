# LVT CLI Roadmap: Generate Production Apps via LLM

## Progress Tracker

| Task | Status | Description |
|------|--------|-------------|
| T1 | DONE | Remove evolution system |
| T1a | DONE | Remove MCP server |
| T2 | DONE | Remove interactive gen mode |
| T3 | DONE | General dead code cleanup |
| T4 | DONE | Upgrade main.go.tmpl — structured logging (all 3 kits + fallback) |
| T5 | DONE | Upgrade main.go.tmpl — http.Server with timeouts (all 3 kits + fallback) |
| T6 | DONE | Upgrade main.go.tmpl — graceful shutdown (all 3 kits + fallback) |
| T7 | DONE | Upgrade main.go.tmpl — security headers (multi/single kits + fallback) |
| T8 | DONE | Upgrade main.go.tmpl — health endpoint (/health/live + /health/ready K8s-compatible) |
| T9 | DONE | Upgrade main.go.tmpl — recovery middleware (all 3 kits + fallback) |
| T10 | DONE | Upgrade main.go.tmpl — env var configuration (DATABASE_PATH, LOG_LEVEL, APP_ENV) |
| T11 | DONE | Clean up client library loading (CLIENT_LIB_PATH env var, removed 4-path hack) |
| T12 | DONE | Generate .env.example and .gitignore during lvt new |
| T13 | DONE | Verify auth integration with new template (all 5 modifications work) |

(Full plan in session transcript — this file tracks progress only)
