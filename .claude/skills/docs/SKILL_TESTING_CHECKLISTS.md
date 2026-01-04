# LVT Claude Code Skills - Testing Checklists

This document provides detailed manual testing checklists for validating Claude Code skills.

---

## Table of Contents

1. [General Testing Checklist](#general-testing-checklist)
2. [Core Skills Checklists](#core-skills-checklists)
3. [Workflow Skills Checklists](#workflow-skills-checklists)
4. [Maintenance Skills Checklists](#maintenance-skills-checklists)
5. [Rating Guidelines](#rating-guidelines)

---

## General Testing Checklist

**Use this for every skill test:**

### Before Starting
- [ ] Clean test environment (fresh session directory)
- [ ] Browser DevTools open (Console tab visible)
- [ ] Terminal visible for server logs
- [ ] Timer started (track testing duration)

### After Skill Execution
- [ ] Automated validation passed
- [ ] All expected files generated
- [ ] No error messages in terminal
- [ ] Git status clean (if applicable)

### During Manual Testing
- [ ] Server starts without errors
- [ ] Page loads in browser
- [ ] No JavaScript console errors
- [ ] No browser console warnings
- [ ] WebSocket connection established
- [ ] UI elements render correctly
- [ ] Interactions work as expected

### After Testing
- [ ] Stop server cleanly
- [ ] Check for orphaned processes
- [ ] Document any issues found
- [ ] Rate experience (1-5 stars)
- [ ] Update SESSION.md with results

---

## Core Skills Checklists

### Skill: lvt:new-app

**Test Scenarios:**
1. Basic app with multi kit (Tailwind)
2. App with single kit (SPA mode)
3. App with simple kit (Pico CSS)
4. App with custom module name
5. Error: Invalid app name

**Manual Testing Checklist:**

#### Application Structure
- [ ] Directory created with app name
- [ ] `.lvtrc` file exists with module name
- [ ] `go.mod` file has correct module name
- [ ] `cmd/{app}/main.go` exists
- [ ] `app/` directory exists (handlers/templates)
- [ ] `database/` directory exists (schema, queries, models)
- [ ] `shared/` directory exists
- [ ] `web/assets/` directory exists

#### Build & Dependencies
- [ ] `go build` succeeds without errors
- [ ] `go mod tidy` completes successfully
- [ ] All dependencies resolved

#### Development Server
- [ ] Server starts on specified port
- [ ] No error messages during startup
- [ ] Server logs show "Listening on..."

#### Homepage
- [ ] Navigate to `http://localhost:8080`
- [ ] Page loads within 2 seconds
- [ ] Title displays correctly
- [ ] Layout renders properly
- [ ] CSS framework styles applied

#### DevTools Console
- [ ] No JavaScript errors (red messages)
- [ ] No JavaScript warnings
- [ ] WebSocket connection successful
- [ ] No 404 errors for assets

#### Network Tab
- [ ] HTML page loads (200 status)
- [ ] CSS loads (200 status)
- [ ] JavaScript loads (200 status)
- [ ] WebSocket upgrade successful (101 status)

#### Server Logs
- [ ] No error messages
- [ ] HTTP requests logged
- [ ] WebSocket connection logged

#### Rating Questions
- Was the app generated correctly? (Yes/No)
- Did everything work on first try? (Yes/No)
- How would you rate the experience? (1-5 stars)
- Would you use this in a real project? (Yes/No/Maybe)

---

### Skill: lvt:add-resource

**Test Scenarios:**
1. Simple resource (3 fields, auto-infer types)
2. Complex resource (10+ fields, mixed types)
3. Resource with explicit types
4. Resource with foreign key
5. Resource with multiple FKs
6. Error: Duplicate resource name
7. Error: Invalid field name

**Manual Testing Checklist:**

#### Code Generation
- [ ] Handler file created: `app/{resource}/{resource}.go`
- [ ] Template file created: `app/{resource}/{resource}.tmpl`
- [ ] Test file created: `app/{resource}/{resource}_test.go`
- [ ] WebSocket test created: `app/{resource}/{resource}_ws_test.go`

#### Database
- [ ] Migration file created with timestamp
- [ ] Schema includes all fields with correct types
- [ ] Foreign key constraints correct (if applicable)
- [ ] Queries added to `queries.sql`
- [ ] `lvt migration up` succeeded
- [ ] `sqlc generate` succeeded

#### Build & Tests
- [ ] `go build` succeeds
- [ ] `go test ./app/{resource}` passes
- [ ] No compilation errors

#### List View
- [ ] Navigate to `/{resource}` route
- [ ] Page loads without errors
- [ ] Shows empty state (if no records)
- [ ] OR shows existing records (if seeded)
- [ ] Table headers match field names
- [ ] Pagination controls visible (if applicable)
- [ ] Search bar visible (if applicable)
- [ ] "New {Resource}" button visible
- [ ] "New {Resource}" button clickable

#### Create Operation
- [ ] Click "New {Resource}" button
- [ ] Modal/page opens (depending on edit mode)
- [ ] Form displays all fields
- [ ] Field types correct (text/number/checkbox/date)
- [ ] Fill in all required fields
- [ ] Click "Submit" or "Create"
- [ ] Form submits without errors
- [ ] New record appears in list
- [ ] Success message shown (if applicable)

#### Read/Detail View
- [ ] Click on record in list
- [ ] Detail view/modal opens
- [ ] All fields display correct values
- [ ] Layout looks clean

#### Update Operation
- [ ] Click "Edit" button on record
- [ ] Edit form opens with current values
- [ ] Modify one or more fields
- [ ] Click "Save" or "Update"
- [ ] Changes persist
- [ ] Updated record shows new values in list
- [ ] No errors in console or server logs

#### Delete Operation
- [ ] Click "Delete" button on record
- [ ] Confirmation prompt appears (if applicable)
- [ ] Confirm deletion
- [ ] Record removed from list
- [ ] No errors in console or server logs
- [ ] Database record actually deleted (check with another operation)

#### Search (if implemented)
- [ ] Type in search box
- [ ] Results filter in real-time
- [ ] Correct records shown
- [ ] Clear search returns all records

#### Pagination (if implemented)
- [ ] Navigate to page 2
- [ ] Different records shown
- [ ] Navigate back to page 1
- [ ] Original records shown

#### Sort (if implemented)
- [ ] Click column header
- [ ] Records sort ascending
- [ ] Click again
- [ ] Records sort descending

#### Console & Logs
- [ ] No JavaScript errors during any operation
- [ ] WebSocket messages sent/received
- [ ] Server logs show each operation
- [ ] No SQL errors

#### Rating Questions
- Did CRUD operations all work? (Yes/No)
- Any broken features? (Yes/No - specify)
- How would you rate the UX? (1-5 stars)
- Would you deploy this to production? (Yes/No/With changes)

---

### Skill: lvt:add-view

**Manual Testing Checklist:**

#### Code Generation
- [ ] Handler file created: `app/{view}/{view}.go`
- [ ] Template file created: `app/{view}/{view}.tmpl`
- [ ] Route added to main.go

#### View Rendering
- [ ] Navigate to `/{view}` route
- [ ] Page loads without errors
- [ ] Content renders correctly
- [ ] Styling applied (CSS framework)
- [ ] Layout consistent with rest of app

#### Console & Logs
- [ ] No JavaScript errors
- [ ] WebSocket connects (if using LiveTemplate features)
- [ ] No server errors

#### Rating Questions
- Did the view generate correctly? (Yes/No)
- Does the styling match the app theme? (Yes/No)
- How would you rate the output? (1-5 stars)

---

### Skill: lvt:add-auth

**Manual Testing Checklist:**

#### Database
- [ ] Auth migrations created
- [ ] Users table exists
- [ ] Tokens table exists (if magic-link enabled)
- [ ] Sessions table exists (if applicable)
- [ ] Queries added to queries.sql

#### Code Generation
- [ ] Password utilities created: `shared/password/`
- [ ] Email utilities created: `shared/email/`
- [ ] Auth queries available

#### Integration (Phase 2 when available)
- [ ] Login handler exists
- [ ] Signup handler exists
- [ ] Logout handler exists
- [ ] Session middleware exists
- [ ] CSRF protection enabled

#### Manual Wiring (Phase 1)
- [ ] Instructions provided for wiring handlers
- [ ] Examples given for middleware usage
- [ ] Next steps clearly explained

#### Rating Questions
- Were the instructions clear? (Yes/No)
- Could you complete the setup? (Yes/No/Partially)
- How would you rate the guidance? (1-5 stars)

---

### Skill: lvt:deploy

**Manual Testing Checklist:**

#### Files Generated
- [ ] `deploy/` directory created
- [ ] Provider-specific configs present
- [ ] Dockerfile exists (if applicable)
- [ ] `.lvtstack` tracking file created

#### Configuration Review
- [ ] Environment variables documented
- [ ] Database configuration correct
- [ ] Backup configuration (if enabled)
- [ ] Redis configuration (if enabled)
- [ ] Storage configuration (if enabled)

#### Validation
- [ ] Run `lvt stack validate`
- [ ] No validation errors
- [ ] Run `lvt stack info`
- [ ] Information displayed correctly

#### Documentation
- [ ] Deployment instructions provided
- [ ] Required secrets documented
- [ ] Next steps clear

#### Rating Questions
- Are the configs production-ready? (Yes/No)
- Would you feel confident deploying this? (Yes/No)
- What's missing? (Free text)
- How would you rate this? (1-5 stars)

---

### Skill: lvt:dev

**Manual Testing Checklist:**

#### Server Startup
- [ ] Server starts on specified port
- [ ] No startup errors
- [ ] Browser opens automatically (if not --no-browser)
- [ ] Page loads in browser

#### Hot Reload (if applicable)
- [ ] Make change to template
- [ ] Page reloads automatically
- [ ] Change visible in browser

#### Monitoring
- [ ] Server logs visible
- [ ] HTTP requests logged
- [ ] Errors highlighted

#### Shutdown
- [ ] Ctrl+C stops server cleanly
- [ ] No orphaned processes

#### Rating Questions
- Did the dev server work correctly? (Yes/No)
- Was hot reload working? (Yes/No/N/A)
- How would you rate the DX? (1-5 stars)

---

### Skill: lvt:test

**Manual Testing Checklist:**

#### Test Execution
- [ ] Tests run without hanging
- [ ] Output is readable
- [ ] Pass/fail clearly indicated
- [ ] Duration shown

#### Results Parsing
- [ ] Failed tests highlighted
- [ ] Error messages shown
- [ ] File paths with line numbers
- [ ] Suggestions for fixes (if available)

#### Rating Questions
- Were the results clear? (Yes/No)
- Were failures easy to understand? (Yes/No)
- How would you rate this? (1-5 stars)

---

## Workflow Skills Checklists

### Skill: lvt:quickstart

**Manual Testing Checklist:**

#### End-to-End Flow
- [ ] User prompted for app type
- [ ] App generated successfully
- [ ] Example resource added
- [ ] Migrations run automatically
- [ ] Dev server starts
- [ ] Browser opens to app
- [ ] CRUD operations work

#### Guidance
- [ ] Each step explained
- [ ] Progress visible
- [ ] Errors handled gracefully
- [ ] Next steps suggested

#### Final State
- [ ] Working app running
- [ ] Tests passing
- [ ] User can immediately interact
- [ ] Clear path forward

#### Timing
- [ ] Total time < 2 minutes
- [ ] Each step completes quickly
- [ ] No long delays

#### Rating Questions
- Did you get a working app quickly? (Yes/No)
- Was the process smooth? (Yes/No)
- Would you recommend this to others? (Yes/No)
- How would you rate this experience? (1-5 stars)

---

### Skill: lvt:production-ready

**Manual Testing Checklist:**

#### Analysis Phase
- [ ] Current app state analyzed correctly
- [ ] Missing features identified
- [ ] Recommendations make sense

#### Auth Addition (if needed)
- [ ] User prompted for auth options
- [ ] Auth setup completes
- [ ] Migrations applied

#### Deployment Setup
- [ ] User prompted for provider
- [ ] Stack generated successfully
- [ ] Configs look correct

#### Environment Setup
- [ ] `.env` template generated
- [ ] All required vars documented
- [ ] Secrets clearly marked

#### Final Checklist
- [ ] Security features verified
- [ ] Tests all passing
- [ ] Deployment instructions provided
- [ ] Production-ready checklist shown

#### Rating Questions
- Do you feel confident deploying? (Yes/No)
- What's still missing? (Free text)
- How would you rate this? (1-5 stars)

---

### Skill: lvt:add-related-resources

**Manual Testing Checklist:**

#### Domain Detection
- [ ] Domain identified correctly
- [ ] Related resources suggested
- [ ] Relationships explained

#### Resource Generation
- [ ] Resources generated in correct order
- [ ] Foreign keys set up correctly
- [ ] Migrations applied successfully
- [ ] All resources working

#### Relationship Verification
- [ ] Can create parent records
- [ ] Can create child records referencing parent
- [ ] Foreign key constraints enforced
- [ ] Cascade behavior correct (if applicable)

#### Final Schema
- [ ] Schema makes sense
- [ ] Relationships clear
- [ ] No circular dependencies

#### Rating Questions
- Were the suggestions helpful? (Yes/No)
- Do the relationships work correctly? (Yes/No)
- How would you rate this? (1-5 stars)

---

## Maintenance Skills Checklists

### Skill: lvt:analyze

**Manual Testing Checklist:**

#### Analysis Accuracy
- [ ] Kit detected correctly
- [ ] Module name correct
- [ ] Resources listed accurately
- [ ] Field types correct
- [ ] Foreign keys identified
- [ ] Auth status correct
- [ ] Deployment status correct

#### Output Quality
- [ ] Information well-formatted
- [ ] Easy to read
- [ ] Actionable insights
- [ ] Clear next steps

#### Rating Questions
- Was the analysis accurate? (Yes/No)
- Was it helpful? (Yes/No)
- How would you rate this? (1-5 stars)

---

### Skill: lvt:suggest

**Manual Testing Checklist:**

#### Pattern Recognition
- [ ] Domain identified correctly
- [ ] Suggestions make sense
- [ ] Relationships explained
- [ ] Prioritized appropriately

#### Suggestion Quality
- [ ] Relevant to current app
- [ ] Not suggesting duplicates
- [ ] Explains reasoning
- [ ] Offers to generate

#### Rating Questions
- Were suggestions helpful? (Yes/No)
- Would you implement them? (Yes/No/Some)
- How would you rate this? (1-5 stars)

---

### Skill: lvt:troubleshoot

**Manual Testing Checklist:**

#### Problem Detection
- [ ] Pending migrations found
- [ ] Missing dependencies found
- [ ] Port conflicts found
- [ ] Database issues found
- [ ] Test failures identified

#### Diagnosis Quality
- [ ] Root cause identified
- [ ] Clear explanation
- [ ] Specific fix suggested
- [ ] Commands to fix provided

#### Fix Application
- [ ] Fixes actually work
- [ ] Problem resolved
- [ ] No new issues introduced

#### Rating Questions
- Did it find the issue? (Yes/No)
- Did the fix work? (Yes/No)
- How would you rate this? (1-5 stars)

---

## Rating Guidelines

### 5 Stars ⭐⭐⭐⭐⭐
- Everything worked perfectly
- No errors encountered
- Great user experience
- Would use in production immediately
- No improvements needed

### 4 Stars ⭐⭐⭐⭐
- Worked well with minor issues
- Small UX improvements possible
- Would use in production after review
- Clear path to 5 stars

### 3 Stars ⭐⭐⭐
- Worked but rough edges
- Some confusion during use
- Requires improvements before production
- Multiple issues to address

### 2 Stars ⭐⭐
- Basic functionality works
- Significant issues encountered
- Needs major improvements
- Not ready for production

### 1 Star ⭐
- Doesn't work as expected
- Critical bugs present
- Unusable in current state
- Requires complete rework

---

## Testing Best Practices

### Do's
- ✅ Follow checklist completely
- ✅ Take screenshots of issues
- ✅ Copy error messages exactly
- ✅ Test in fresh environment
- ✅ Document everything
- ✅ Be honest in ratings

### Don'ts
- ❌ Skip checklist items
- ❌ Test in production environment
- ❌ Ignore console errors
- ❌ Give generous ratings without justification
- ❌ Batch testing without clean environment
- ❌ Assume it works without verifying

---

## Issue Reporting Template

When you find an issue, document it like this:

```markdown
## Issue Found During Testing

**Skill:** lvt:skill-name
**Test Scenario:** Specific scenario
**Step:** Checklist item that failed

**Expected:**
What should have happened

**Actual:**
What actually happened

**Error Messages:**
```
Copy exact error messages here
```

**Screenshots:**
[Attach if applicable]

**Console Logs:**
```
Copy browser console output
```

**Server Logs:**
```
Copy server terminal output
```

**Severity:** P0 / P1 / P2 / P3

**Rating Impact:**
How this affects the rating (e.g., "Blocks 5 star rating")
```

---

## Checklist for Checklist Testing

Meta, but important:

- [ ] Checklist is complete
- [ ] Steps are clear
- [ ] Order makes sense
- [ ] Nothing is ambiguous
- [ ] Can be followed by someone unfamiliar with skill
- [ ] Covers happy path
- [ ] Covers error cases
- [ ] Rating guidelines clear

---

**Remember:** The goal of manual testing is to catch issues that automated tests miss - primarily UX, clarity, and real-world usability problems. Be thorough, be critical, and be honest in your assessments!
