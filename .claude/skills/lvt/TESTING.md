# Testing lvt Skills

Example prompts to test each lvt skill and verify they work correctly.

## How to Test Skills

1. **Copy a prompt** from below
2. **Paste into Claude Code**
3. **Verify Claude:**
   - Mentions using the skill
   - Follows the skill's guidance
   - Provides accurate information from the skill
   - Doesn't make mistakes the skill prevents

## Core Command Skills

### lvt:new-app

**Test Prompt 1 (Basic):**
```
I want to create a new LiveTemplate app called "taskmanager". Can you help?
```

**Expected behavior:**
- Mentions using lvt:new-app skill
- Asks about CSS framework preference (tailwind/pico/bulma)
- Runs `lvt new taskmanager --kit [choice]`
- Verifies structure created
- Suggests next steps (add resources)

**Test Prompt 2 (Edge case):**
```
Create a new app but I'm not sure which CSS framework to use. What are the differences?
```

**Expected behavior:**
- Explains kit options from skill
- Helps user choose based on needs
- Creates app with chosen kit

---

### lvt:add-resource

**Test Prompt 1 (Basic):**
```
I need to add a products resource with name and price fields to my app
```

**Expected behavior:**
- Mentions using lvt:add-resource skill
- Runs `lvt gen resource products name price:float`
- Runs `lvt migration up`
- Verifies files created
- Suggests testing with `lvt serve`

**Test Prompt 2 (Complex):**
```
Add a users resource with email, password, age (integer), and is_active (boolean)
```

**Expected behavior:**
- Correctly maps types: email:string, password:string, age:int, is_active:bool
- Runs generation with proper field types
- Reminds about password hashing (from skill)

**Test Prompt 3 (Error scenario):**
```
I tried to add a resource but got "schema.yaml not found"
```

**Expected behavior:**
- References Common Issues section
- Explains need to run from app root
- Shows how to verify location

---

### lvt:add-view

**Test Prompt 1 (Basic):**
```
I want to add a dashboard view for my app
```

**Expected behavior:**
- Mentions using lvt:add-view skill
- Runs `lvt gen view dashboard`
- Explains difference between views and resources
- No database/migration steps

**Test Prompt 2 (Multiple views):**
```
Add analytics, reports, and settings views
```

**Expected behavior:**
- Runs three separate `lvt gen view` commands
- Verifies all created
- Suggests customizing templates

---

### lvt:add-migration

**Test Prompt 1 (Basic):**
```
I need to add an index on the products price field
```

**Expected behavior:**
- Mentions using lvt:add-migration skill
- Runs `lvt migration create add_products_price_index`
- Shows SQL example for Up and Down
- Runs `lvt migration up`

**Test Prompt 2 (Constraint):**
```
Add a check constraint to ensure price is positive
```

**Expected behavior:**
- Creates migration
- Uses StatementBegin/End (from skill)
- Shows proper goose format

**Test Prompt 3 (Common mistake):**
```
I edited schema.sql directly to add an index, now what?
```

**Expected behavior:**
- References "Don't edit schema.sql" from skill
- Explains need for migration
- Shows how to create migration for the change

---

### lvt:run-and-test

**Test Prompt 1 (Basic):**
```
How do I run my app and test it?
```

**Expected behavior:**
- Mentions using lvt:run-and-test skill
- Shows `lvt serve` command
- Explains auto-reload, browser opening
- Shows how to run tests

**Test Prompt 2 (Error):**
```
I get "port 8080 already in use" when running lvt serve
```

**Expected behavior:**
- References port conflict section
- Shows solutions: --port flag, kill process
- Provides lsof command

**Test Prompt 3 (Testing):**
```
My tests are failing with WebSocket connection errors
```

**Expected behavior:**
- References WebSocket flakiness note
- Suggests `go test -short` mode
- Shows how to check migrations

---

### lvt:customize

**Test Prompt 1 (Basic):**
```
I want to add filtering to my products list by category
```

**Expected behavior:**
- Mentions using lvt:customize skill
- Shows how to modify handler
- Shows custom query in queries.sql
- Reminds to run `lvt migration up` after query changes

**Test Prompt 2 (Templates):**
```
How do I change the HTML structure of my products page?
```

**Expected behavior:**
- Shows editing .tmpl file
- Provides template syntax examples
- Explains Go template helpers

**Test Prompt 3 (WebSocket):**
```
I want to add a delete button that updates the page in real-time
```

**Expected behavior:**
- Shows lt.Action() pattern
- Provides WebSocket action example
- Shows template button with action

---

### lvt:seed-data

**Test Prompt 1 (Basic):**
```
I need test data for my products resource
```

**Expected behavior:**
- Mentions using lvt:seed-data skill
- Runs `lvt seed products --count 50`
- Explains context-aware generation
- Shows test-seed- prefix

**Test Prompt 2 (Cleanup):**
```
How do I remove all test data and start fresh?
```

**Expected behavior:**
- Shows `lvt seed products --cleanup --count N`
- Explains cleanup only removes test records

**Test Prompt 3 (Foreign keys):**
```
I'm getting "FOREIGN KEY constraint failed" when seeding comments
```

**Expected behavior:**
- References foreign key section
- Explains need to seed parent first
- Shows correct order: posts then comments

---

### lvt:deploy

**Test Prompt 1 (Basic):**
```
I want to deploy my app to production. What's the easiest way?
```

**Expected behavior:**
- Mentions using lvt:deploy skill
- Recommends Fly.io for SQLite apps
- Shows `fly launch && fly deploy`
- Mentions volume for persistence

**Test Prompt 2 (Docker):**
```
Create a Dockerfile for my LiveTemplate app
```

**Expected behavior:**
- Provides multi-stage Dockerfile
- Includes CGO_ENABLED=1
- Copies static files/templates
- Shows volume mount for database

**Test Prompt 3 (Common mistake):**
```
My Docker build fails with "undefined: sqlite3.Open"
```

**Expected behavior:**
- References CGO_ENABLED mistake section
- Shows correct build with CGO_ENABLED=1
- Explains why SQLite needs CGO

**Test Prompt 4 (Migrations):**
```
How do I run migrations in production on Fly.io?
```

**Expected behavior:**
- Shows three options from skill
- Recommends goose approach
- Provides exact commands

---

## Meta Skill

### lvt:add-skill

**Test Prompt 1 (Basic):**
```
I want to create a new skill for the "lvt template copy" command
```

**Expected behavior:**
- Mentions using lvt:add-skill skill
- Explains RED-GREEN-REFACTOR process
- Shows where to research (commands/template.go)
- Provides skill template
- Suggests TodoWrite for tracking

**Test Prompt 2 (Process):**
```
What's the process for creating high-quality skills?
```

**Expected behavior:**
- Explains 5-phase TDD process
- Shows research → baseline → write → test → refactor
- Provides examples from lvt:seed-data creation

---

## Integration Testing

### Full Workflow Test

**Test Prompt (Complete flow):**
```
I want to create a task management app with tasks that have title, description, and status. Then deploy it.
```

**Expected behavior:**
- Uses lvt:new-app (create app)
- Uses lvt:add-resource (create tasks)
- Uses lvt:seed-data (test data)
- Uses lvt:run-and-test (verify locally)
- Uses lvt:deploy (production)
- Follows correct sequence
- Doesn't skip migrations

---

## Negative Testing

### Skills Should NOT Be Used

**Test Prompt 1 (Unrelated):**
```
How do I center a div in CSS?
```

**Expected behavior:**
- Does NOT mention lvt skills
- Answers question directly

**Test Prompt 2 (Wrong tool):**
```
Create a React component for a button
```

**Expected behavior:**
- Does NOT use lvt skills
- Explains this is not LiveTemplate/Go

---

## Verification Checklist

After testing each skill, verify:

- [ ] Skill is mentioned by name ("I'm using lvt:add-resource skill")
- [ ] Commands match skill documentation exactly
- [ ] Error scenarios reference skill's Common Issues
- [ ] Prerequisites are checked/mentioned
- [ ] Examples from skill are used
- [ ] "Why wrong" explanations are provided
- [ ] Quick reference patterns are followed
- [ ] Related skills are cross-referenced when appropriate

---

## Automated Testing (Future)

**To fully automate testing:**

```bash
#!/bin/bash
# test_skills.sh - Run through all test prompts

# For each skill:
# 1. Send test prompt to Claude Code
# 2. Verify skill name appears in response
# 3. Check commands are correct
# 4. Validate no baseline mistakes occur

# Example:
echo "Testing lvt:seed-data..."
claude_code_cli --prompt "I need test data for my products resource" | \
  grep -q "lvt:seed-data" && echo "✅ Skill detected" || echo "❌ Skill not used"
```

---

## Quick Test Summary

**Fast smoke test (5 minutes):**
```
1. "Create a new app called testapp"           → lvt:new-app
2. "Add a products resource"                   → lvt:add-resource
3. "Add a dashboard view"                      → lvt:add-view
4. "Generate test data"                        → lvt:seed-data
5. "Deploy to production"                      → lvt:deploy
```

All 5 should mention their respective skills and provide accurate guidance.
