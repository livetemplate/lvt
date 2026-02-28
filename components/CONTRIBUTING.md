# Contributing to LiveTemplate Components

Thank you for your interest in contributing to the LiveTemplate Components Library!

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Submit a pull request

## Component Structure

Every component must follow this structure:

```
components/<name>/
  <name>.go           # Component struct, constructors, and methods
  options.go          # Functional options
  templates.go        # Template embedding with //go:embed
  <name>_test.go      # Comprehensive tests
  templates/
    default.tmpl      # Default template variant
    [variant].tmpl    # Additional variants as needed
```

## Naming Conventions

### Package Names
- Use lowercase with no hyphens
- For hyphenated component names, remove hyphens: `date-picker` -> `datepicker`

### Template Names
All templates must follow the pattern: `lvt:<component>:<variant>:v<version>`

```html
{{define "lvt:dropdown:default:v1"}}...{{end}}
{{define "lvt:dropdown:searchable:v1"}}...{{end}}
```

### Go Types
- Main struct: PascalCase matching component name (`Dropdown`, `DatePicker`)
- Option type: `Option` (always `func(*ComponentName)`)
- Constructor: `New(id string, opts ...Option)`

## Template Requirements

### Styled and Unstyled Variants
Every template must support both styled (Tailwind) and unstyled modes:

```html
{{define "lvt:mycomponent:default:v1"}}
{{if .IsStyled}}
{{/* Tailwind CSS styled version */}}
<div class="p-4 border rounded-lg bg-white shadow-sm" data-component="mycomponent" data-id="{{.ID}}">
  {{/* ... */}}
</div>
{{else}}
{{/* Unstyled semantic HTML version */}}
<div data-component="mycomponent" data-id="{{.ID}}">
  {{/* ... */}}
</div>
{{end}}
{{end}}
```

### Required Attributes
- `data-component="<name>"` - Identifies the component type
- `data-id="{{.ID}}"` - Unique instance identifier

### Use Generic lvt-* Attributes
Use existing generic attributes, NOT component-specific ones:

**Good:**
```html
lvt-click="toggle_{{.ID}}"
lvt-click-away="close_{{.ID}}"
lvt-input="search_{{.ID}}"
```

**Bad:**
```html
lvt-dropdown="..."      <!-- No component-specific attributes -->
lvt-calendar-nav="..."  <!-- Use generic lvt-click instead -->
```

### Accessibility
- Include appropriate ARIA attributes
- Support keyboard navigation where applicable
- Use semantic HTML elements

## Code Requirements

### Base Struct
Embed the base struct in every component:

```go
import "github.com/livetemplate/components/base"

type MyComponent struct {
    base.Base
    // Your fields here
}

func New(id string, opts ...Option) *MyComponent {
    c := &MyComponent{
        Base: base.NewBase(id, "mycomponent"),
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}
```

### Functional Options
Provide functional options for all customization:

```go
type Option func(*MyComponent)

func WithLabel(label string) Option {
    return func(c *MyComponent) {
        c.Label = label
    }
}

func WithStyled(styled bool) Option {
    return func(c *MyComponent) {
        c.SetStyled(styled)
    }
}
```

### Templates File
Use embedded templates with the correct pattern:

```go
package mycomponent

import (
    "embed"
    "github.com/livetemplate/components/base"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

func Templates() *base.TemplateSet {
    return base.NewTemplateSet(templateFS, "templates/*.tmpl", "mycomponent")
}
```

## Test Requirements

### Minimum Coverage
- Constructor creates component with defaults
- All options work correctly
- Template set is not nil
- Edge cases (empty values, nil inputs)

### Test Example

```go
func TestNew(t *testing.T) {
    t.Run("creates with defaults", func(t *testing.T) {
        c := New("test")
        if c.ID() != "test" {
            t.Errorf("expected ID 'test', got %q", c.ID())
        }
        if c.Namespace() != "mycomponent" {
            t.Errorf("expected namespace 'mycomponent', got %q", c.Namespace())
        }
    })
}

func TestOptions(t *testing.T) {
    t.Run("WithLabel", func(t *testing.T) {
        c := New("test", WithLabel("My Label"))
        if c.Label != "My Label" {
            t.Errorf("expected label 'My Label', got %q", c.Label)
        }
    })
}

func TestTemplates(t *testing.T) {
    ts := Templates()
    if ts == nil {
        t.Fatal("Templates() returned nil")
    }
}
```

## Pull Request Checklist

Before submitting a PR, ensure:

- [ ] Component follows the standard structure
- [ ] Package name is lowercase without hyphens
- [ ] Template names follow `lvt:<component>:<variant>:v1` convention
- [ ] Both styled (Tailwind) and unstyled variants work
- [ ] Uses generic `lvt-*` attributes (no component-specific ones)
- [ ] Embeds `base.Base` struct
- [ ] Provides functional options for all customization
- [ ] Includes comprehensive tests (all pass)
- [ ] Templates have proper ARIA attributes
- [ ] Code is formatted with `go fmt`
- [ ] No linting errors
- [ ] Updated `all.go` to include new component

## Updating all.go

After adding a component, update `components/all.go`:

```go
import (
    // ... existing imports ...
    "github.com/livetemplate/components/mycomponent"
)

func All() []*base.TemplateSet {
    return []*base.TemplateSet{
        // ... existing components ...
        mycomponent.Templates(),
    }
}
```

## Running Tests

```bash
# Run all component tests
go test ./...

# Run specific component tests
go test ./dropdown/...

# Run with verbose output
go test -v ./...
```

## Questions?

If you have questions about contributing:
1. Check existing components for examples
2. Open an issue for discussion
3. Ask in the LiveTemplate community

Thank you for contributing!
