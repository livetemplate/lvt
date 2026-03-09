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
| T14 | DONE | E2e tests for production template features (chromedp) |
| T15 | DONE | Template edge case testing (4 combos: bare, views-only, resource, auth-only) |
| T16 | N/A | No bugs found in template testing — all e2e tests pass |
| T16a | DONE | Modal + toast verified via TestCompleteWorkflow_BlogApp (all subtests pass) |
| T16b | PENDING | Dropdown in select field (requires separate test) |
| T16c | PENDING | Assess remaining component browser test needs |
| T16d | N/A | No component bugs found yet |
| T17 | DONE | Audit core skills — fixed health endpoint references |
| T18 | DONE | Audit deployment skills — fixed health endpoint references |
| T19 | DONE | Audit remaining skills — no stale references found |
| T20 | N/A | MCP server removed in T1a |
| T21 | DONE | Automated skill pipeline test (full LLM workflow, chromedp verified) |

### Remaining tasks (Phases 5-7)

| Task | Status | Description |
|------|--------|-------------|
| T22-T27 | PENDING | Phase 5: Expose LiveTemplate library features (health checkers, Redis, broadcasting) |
| T28-T30 | PENDING | Phase 6: End-to-end LLM workflow test (Docker, Fly.io, flakiness fixes) |
| T31-T36 | PENDING | Phase 7: Documentation and release prep |

(Full plan in main plan file — this file tracks progress only)
