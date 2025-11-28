---
name: lvt-manage-kits
description: Manage CSS framework kits - list available kits, view kit info, validate kits, and create custom kits
category: core
version: 1.0.0
---

# lvt:manage-kits

Manage CSS framework kits in LiveTemplate. List system/local/community kits, view kit details, validate kit structure, customize existing kits, and create new kits.

## User Prompts

**When to use:**
- "What CSS frameworks are available?"
- "Show me all kits"
- "Tell me about the Tailwind kit"
- "Validate my custom kit"
- "Create a new kit for Bootstrap"

**Examples:**
- "List all kits"
- "Show info for multi kit"
- "Validate my custom-theme kit"
- "What kits can I use?"

## Quick Reference

```bash
# List all kits
lvt kits list
lvt kits list --filter system        # Only system kits
lvt kits list --filter local          # Only local kits
lvt kits list --format json           # JSON output

# View kit info
lvt kits info multi                   # Show multi kit details
lvt kits info single                  # Show single kit details

# Validate a kit
lvt kits validate custom-kit          # Validate kit structure

# Create new kit (advanced)
lvt kits create my-bootstrap-kit --based-on multi

# Customize existing kit
lvt kits customize multi --name my-custom-multi
```

## What It Does

### List Command (`lvt kits list`)

**Shows:**
- Kit name
- CSS framework (Tailwind, Bulma, Pico, etc.)
- Description
- Source (system, local, community)

**Filters:**
- `--filter system` - Built-in kits (multi, single, simple)
- `--filter local` - Your custom kits
- `--filter community` - Community-contributed kits
- `--search <term>` - Search by name or framework

**Formats:**
- `--format table` - Pretty table (default)
- `--format json` - Machine-readable JSON
- `--format simple` - Name-only list

### Info Command (`lvt kits info <kit>`)

**Shows:**
- Kit name and description
- CSS framework
- Components included
- Templates available
- Helper functions
- Configuration options

### Validate Command (`lvt kits validate <kit>`)

**Checks:**
- Kit manifest (kit.yaml) exists and valid
- Required components present
- Template syntax valid
- Helpers properly defined
- No missing dependencies

## System Kits

**multi** (Multi-page with Tailwind CSS)
- Full-featured multi-page applications
- Tailwind CSS framework
- All pagination modes
- Modal and page edit modes
- Component library included

**single** (Single-page with Tailwind CSS)
- SPA-style applications
- Tailwind CSS framework
- Client-side routing
- Component-based architecture

**simple** (Minimal with Pico CSS)
- Lightweight applications
- Pico CSS (classless)
- Minimal dependencies
- Perfect for prototypes

## Checklist

- [ ] Determine which kit command user wants (list/info/validate/create/customize)
- [ ] For list: Run with appropriate filters
- [ ] For info: Extract kit name, run `lvt kits info <kit>`
- [ ] For validate: Extract kit name, run `lvt kits validate <kit>`
- [ ] For create/customize: Guide through advanced workflow
- [ ] Explain output to user

## Common Issues

**Issue: "kit not found"**
- Kit name doesn't exist
- Fix: Run `lvt kits list` to see available kits

**Issue: "invalid manifest"**
- kit.yaml has syntax errors
- Fix: Validate YAML syntax, check required fields

**Issue: "missing components"**
- Kit missing required component files
- Fix: Compare with system kit structure, add missing files

## Example Output

**List (table format):**
```
Available kits (3):
Name     Framework   Description                          Source
multi    Tailwind    Multi-page apps with Tailwind CSS    system
single   Tailwind    SPA with Tailwind CSS                system
simple   Pico        Minimal apps with Pico CSS           system
```

**Info:**
```
Kit: multi
Framework: Tailwind CSS
Description: Full-featured multi-page applications

Components:
  - layout.tmpl (base page layout)
  - table.tmpl (data tables)
  - form.tmpl (input forms)
  - pagination.tmpl (infinite, load-more, prev-next, numbers)
  - toolbar.tmpl (actions, search, filters)
  - stats.tmpl (metrics display)

Templates:
  - resource/* (CRUD resources)
  - view/* (UI-only pages)
  - auth/* (authentication)
  - app/* (application base)

Helpers:
  - Icon(), Button(), Alert(), Badge()
  - Table(), Form(), Input(), Select()
```

## Use Cases

1. **Choose a kit:** Before creating app, see what's available
2. **Validate custom kit:** After creating custom kit, verify structure
3. **Explore features:** See what components each kit provides
4. **Create new kit:** Build kit for your preferred CSS framework
5. **Customize kit:** Modify existing kit for specific needs

## Creating Custom Kits (Advanced)

**Directory structure:**
```
~/.config/lvt/kits/my-custom-kit/
├── kit.yaml                    # Manifest
├── components/                 # Reusable components
│   ├── layout.tmpl
│   ├── table.tmpl
│   └── form.tmpl
└── templates/                  # Generation templates
    ├── resource/
    ├── view/
    └── app/
```

**Manifest (kit.yaml):**
```yaml
name: my-custom-kit
description: Custom kit with Bootstrap
cssFramework: Bootstrap
version: 1.0.0

components:
  - layout
  - table
  - form

helpers:
  - Icon
  - Button
```

## Notes

- System kits (multi, single, simple) cannot be modified
- Custom kits go in ~/.config/lvt/kits/
- Kit chosen during `lvt new` cannot be changed later
- Generated code uses kit's components and helpers
- Validating kits is recommended before use
- Community kits require internet connection
