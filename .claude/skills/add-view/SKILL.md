---
name: lvt-add-view
description: Add a view-only handler to an existing LiveTemplate app without database integration - perfect for static pages, dashboards, and UI-only components
category: core
version: 1.0.0
keywords: ["lvt", "livetemplate", "lt"]
---

# lvt:add-view

Adds a view-only handler to an existing LiveTemplate application. Unlike resources, views don't interact with the database - they're pure UI components perfect for static pages, dashboards, landing pages, or any content that doesn't need CRUD operations.

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

This skill should activate when the user requests to add a view-only handler:

**Explicit prompts:**
- "Add a [name] view to my app"
- "Create a view for [purpose]"
- "Generate a [name] page"
- "I need a static [name] page"

**Implicit prompts:**
- "Add a dashboard"
- "Create an about page"
- "I want a landing page"
- "Make a contact page"
- "Add a home page"

**Examples:**
- "Add a dashboard view"
- "Create an about page for my app"
- "I need a landing page"
- "Generate a contact view"
- "Make a stats dashboard"

## Context Awareness

Before executing this skill, verify:

1. **In Project Directory:**
   - Check for `.lvtrc` file (confirms it's an lvt project)
   - Check for `go.mod` (confirms it's a Go project)
   - Check for `app/` directory

2. **Dependencies Available:**
   - `lvt` binary is installed and accessible
   - Project was created with `lvt new`

3. **Not Already Exists:**
   - Check if view directory already exists in `app/[viewname]/`
   - Warn user if view name conflicts

## Checklist

- [ ] **Step 1:** Verify we're in an lvt project directory
  - Check for `.lvtrc` file
  - Check for `go.mod` file
  - Check for `app/` directory
  - If missing, inform user they need to create an app first (use lvt:new-app skill)

- [ ] **Step 2:** Validate prerequisites
  - Verify `lvt` command is available
  - Check current directory is project root

- [ ] **Step 3:** Extract view details from user request
  - View name (e.g., "dashboard", "about", "contact", "landing")
  - Purpose/description (helps understand what content to include)

- [ ] **Step 4:** Normalize view name
  - Convert to lowercase for directory/package name
  - Ensure it's a valid Go identifier (alphanumeric + underscore, no spaces)
  - Common names: dashboard, about, contact, landing, privacy, terms, faq, help

- [ ] **Step 5:** Check for naming conflicts
  - Check if `app/[viewname]/` already exists
  - If exists, ask user if they want to:
    - Overwrite existing view (warning: will lose customizations)
    - Choose a different name
    - Cancel operation

- [ ] **Step 6:** Build and run the `lvt gen view` command
  - Format: `lvt gen view <name>`
  - Example: `lvt gen view dashboard`

- [ ] **Step 7:** Verify view generation succeeded
  - Check for success message from lvt
  - Verify files created:
    - `app/<view>/<view>.go` (handler)
    - `app/<view>/<view>.tmpl` (template)
    - `app/<view>/<view>_test.go` (test)
  - Verify files updated:
    - `cmd/<app>/main.go` or `main.go` (route injected)
    - `app/home/home.tmpl` (view registered on home page)

- [ ] **Step 8:** Run `go mod tidy`
  - Ensure all dependencies are up to date
  - Verify no errors

- [ ] **Step 9:** Verify app builds successfully
  - For multi/single kits: `go build ./cmd/<app>`
  - For simple kit: `go build`
  - If build fails, diagnose and fix issues

- [ ] **Step 10:** Provide user with success summary
  - List files created
  - List files updated
  - Show the generated route
  - Suggest customization steps (edit handler, edit template)
  - Provide next steps (run app, test view)

## View Use Cases

Help users understand when to use views vs resources:

**Use Views For:**
- Static pages (about, contact, privacy, terms)
- Dashboards (displaying aggregated data without CRUD)
- Landing pages (marketing, feature showcases)
- Help/FAQ pages
- Custom UI components (charts, reports, analytics)
- Pages that read from API but don't use database

**Use Resources For:**
- CRUD operations (users, posts, products)
- Data management (create, read, update, delete)
- Database-backed features

## Error Handling

**If view directory already exists:**
1. Warn user about potential data loss
2. Ask for confirmation before proceeding
3. Suggest using a different name

**If not in an lvt project:**
1. Check for `.lvtrc` file
2. If missing, inform user they need to create an app first
3. Suggest using lvt:new-app skill

**If view name is invalid:**
1. Check for Go identifier validity (alphanumeric + underscore, no spaces)
2. Suggest corrections if needed
3. Avoid names that conflict with Go keywords

**If build fails after generation:**
1. Run `go mod tidy` again
2. Check for import errors
3. Check route injection didn't break main.go syntax

## Customization Guidance

After generating a view, guide users on common customizations:

**Handler Customization (`app/<view>/<view>.go`):**
- Add data to pass to template
- Fetch data from external APIs
- Add query parameters handling
- Add form handling (if needed)
- Add middleware (auth, logging, etc.)

**Template Customization (`app/<view>/<view>.tmpl`):**
- Update page title and content
- Add custom HTML/CSS
- Use CSS framework utilities (Tailwind, Bulma, Pico)
- Add interactive elements
- Include components from other views

**Test Customization (`app/<view>/<view>_test.go`):**
- Add E2E tests for specific functionality
- Test WebSocket interactions
- Test form submissions
- Verify page content

## Success Response

After successful view generation, provide:

```
✅ View '[name]' generated successfully!

📁 Files created:
  - app/[view]/[view].go
  - app/[view]/[view].tmpl
  - app/[view]/[view]_test.go

📝 Files updated:
  - cmd/[app]/main.go (route: /[view])
  - app/home/home.tmpl (view link added)

✅ App builds successfully

🚀 Next steps:
  1. Customize content: app/[view]/[view].tmpl
  2. Add logic: app/[view]/[view].go
  3. Start your app: lvt serve (or go run main.go)
  4. Visit: http://localhost:8080/[view]
  5. Run tests: go test ./app/[view]
```

## Common User Scenarios

**Scenario 1: Dashboard page**
- User: "Add a dashboard view"
- Generates: dashboard handler, template, test
- Route: /dashboard
- Purpose: Display analytics, charts, summaries

**Scenario 2: About page**
- User: "Create an about page"
- Generates: about handler, template, test
- Route: /about
- Purpose: Static company/app information

**Scenario 3: Landing page**
- User: "I need a landing page"
- Generates: landing handler, template, test
- Route: /landing
- Purpose: Marketing, feature showcase

**Scenario 4: Custom UI component**
- User: "Make a stats view that shows real-time metrics"
- Generates: stats handler, template, test
- Route: /stats
- Purpose: Real-time data display (WebSocket-powered)

## Validation Criteria

View generation is successful if:
1. ✅ All files created without errors
2. ✅ Route injected correctly in main.go
3. ✅ View registered on home page
4. ✅ `go build` succeeds
5. ✅ No compilation errors
6. ✅ App can be started and view accessed

## Notes

- View names are automatically lowercase for routes and packages
- Views don't have database migrations or sqlc queries
- Views use `Handler()` not `Handler(queries)`
- Views can still be interactive using WebSocket (LiveTemplate's reactivity)
- Views are perfect for pages that don't need persistent storage
- Multiple views can be added to the same app
- Views can include forms that POST to resources
