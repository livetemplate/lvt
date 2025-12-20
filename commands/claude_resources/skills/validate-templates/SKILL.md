---
name: lvt-validate-templates
description: Validate LiveTemplate template files - check syntax, parse errors, execution issues, and common problems
category: core
version: 1.0.0
keywords: ["lvt", "livetemplate", "lt"]
---

# lvt:validate-templates

Validates LiveTemplate template files (*.tmpl) for syntax errors, parsing issues, execution problems, and common mistakes. Uses both html/template and LiveTemplate parsers to ensure templates work correctly.

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
- "Validate my template file"
- "Check if my template has errors"
- "Parse this template"
- "Is my template syntax correct?"
- "Debug template errors"

**Examples:**
- "Validate app/posts/posts.tmpl"
- "Check the dashboard template for errors"
- "Parse my custom template"
- "Why isn't my template working?"

## Quick Reference

```bash
# Validate a template file
lvt parse <template-file>

# Examples
lvt parse app/posts/posts.tmpl
lvt parse app/home/home.tmpl
lvt parse custom/template.tmpl
```

## What It Checks

**1. html/template parsing:**
- Basic Go template syntax
- Properly closed tags
- Valid actions ({{...}})
- No syntax errors

**2. LiveTemplate parsing:**
- LiveTemplate-specific features
- WebSocket actions
- Component references
- Framework helpers

**3. Template execution:**
- Can execute with sample data
- No runtime errors
- Generates valid HTML

**4. Common issues:**
- Unclosed tags
- Mismatched {{...}} pairs
- Invalid function calls
- Missing required fields
- Suspicious patterns

## Example Output

**Successful validation:**
```
Parsing template: app/posts/posts.tmpl
Template name: posts
============================================================

1. Testing html/template parsing...
   ✅ Successfully parsed with html/template

2. Defined templates:
   - posts
   - posts_table
   - posts_form

3. Testing LiveTemplate parsing...
   ✅ Successfully parsed with LiveTemplate

4. Testing template execution...
   ✅ Successfully executed (generated 2847 bytes of HTML)

5. Checking for common issues...
   ✅ No issues found

============================================================
✅ Template is valid!
```

**Failed validation:**
```
Parsing template: app/broken/broken.tmpl
Template name: broken
============================================================

1. Testing html/template parsing...
   ❌ Parse error: template: broken:15: unexpected "}" in operand

============================================================
❌ Template validation failed
```

## Checklist

- [ ] Extract template file path from user request
- [ ] Verify file exists
- [ ] Run: `lvt parse <template-file>`
- [ ] Check validation results
- [ ] If errors found, help user fix them
- [ ] If warnings, explain what they mean

## Common Issues Found

### Issue 1: Unclosed tags
**Example:** `{{range .Items` (missing }})
**Fix:** Close all template actions properly

### Issue 2: Invalid function calls
**Example:** `{{UnknownFunction}}`
**Fix:** Use built-in functions or kit helpers

### Issue 3: Missing end tags
**Example:** `{{range .Items}}` without `{{end}}`
**Fix:** Add matching `{{end}}` tag

### Issue 4: Type mismatches
**Example:** `{{.Count | add "text"}}`
**Fix:** Ensure function arguments match expected types

### Issue 5: Undefined fields
**Example:** `{{.NonExistentField}}`
**Warning:** Field might not exist in data structure

## Template Syntax Reminder

**Valid actions:**
```html
<!-- Variables -->
{{.FieldName}}
{{$var := .Value}}

<!-- Conditionals -->
{{if .Show}}...{{end}}
{{if .Show}}...{{else}}...{{end}}

<!-- Loops -->
{{range .Items}}
  {{.Name}}
{{end}}

<!-- With (context) -->
{{with .User}}
  {{.Name}}
{{end}}

<!-- Template inclusion -->
{{template "component" .}}

<!-- Functions/Pipes -->
{{.Title | uppercase}}
{{add .Count 1}}
```

## Use Cases

1. **Before deploying:** Validate all templates
2. **After customization:** Ensure changes didn't break syntax
3. **Debugging:** Find syntax errors quickly
4. **Learning:** Understand template structure
5. **CI/CD:** Automated template validation

## Notes

- Validates both html/template and LiveTemplate
- Execution test uses sample data (might not match real data)
- Execution errors might be OK if template expects specific data
- Can validate any .tmpl file, not just generated ones
- Catches most syntax errors before runtime
- Doesn't validate business logic or data correctness
- Fast validation (no server startup required)
