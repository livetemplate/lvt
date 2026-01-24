# LVT + Components Integration Strategy

**Goal:** Use lvt as the primary driver for components evolution while maintaining component independence.

## The Bootstrap Problem

```
Components need usage to mature -> But they're too unstable to use -> So they don't get usage -> They stay unstable
```

**Solution:** lvt becomes the forcing function. By making lvt depend on and exercise components, we create:
1. Real-world usage patterns
2. Bug discovery through generation
3. Evolution feedback loop
4. Path to v1.0 stability

---

## Repository Strategy: Monorepo vs Multi-Repo

### Option A: Multi-Repo (Separate components repository)

```
github.com/livetemplate/lvt          # CLI + generator
github.com/livetemplate/components   # Reusable components (separate repo)
github.com/livetemplate/livetemplate # Core framework
```

**Pros:**
- Clear separation of concerns
- Independent release cycles
- Smaller repo size

**Cons:**
- Cross-repo PRs for related changes
- Version coordination overhead
- Slower feedback loop (fix -> release -> bump -> test)
- Evolution system needs multi-repo support

### Option B: Monorepo (Components inside lvt) - RECOMMENDED

```
github.com/livetemplate/lvt/
├── go.mod                           # github.com/livetemplate/lvt
├── main.go
├── commands/
├── internal/
│   ├── generator/
│   ├── kits/
│   └── evolution/
│
├── components/                      # Nested module
│   ├── go.mod                       # github.com/livetemplate/lvt/components
│   ├── modal/
│   │   ├── modal.go
│   │   ├── options.go
│   │   └── templates.go
│   ├── toast/
│   ├── dropdown/
│   ├── toggle/
│   └── styles/
│       ├── adapter.go
│       ├── tailwind/
│       └── unstyled/
│
└── (livetemplate/livetemplate remains separate - core framework)
```

**Pros:**
- Single repo for evolution system
- Atomic commits (component fix + lvt update together)
- Faster iteration cycle
- Simpler CI/CD
- One PR for related changes

**Cons:**
- Larger repo
- Need to ensure components stay independently importable

### Why Monorepo Works for Independence

Go's module system allows nested modules with independent import paths:

```go
// External app (not using lvt) can import components directly:
import "github.com/livetemplate/lvt/components/modal"

// lvt internally imports the same path:
import "github.com/livetemplate/lvt/components/modal"

// Generated apps also use the same import:
import "github.com/livetemplate/lvt/components/modal"
```

The nested `components/go.mod` makes it a separate module:

```go
// components/go.mod
module github.com/livetemplate/lvt/components

go 1.22

require github.com/livetemplate/livetemplate v0.8.0
// Note: NO dependency on github.com/livetemplate/lvt
```

### Migration Path from Existing Components Repo

```bash
# 1. Move components code into lvt
git subtree add --prefix=components \
    git@github.com:livetemplate/components.git main

# 2. Update components/go.mod with new module path
# Old: module github.com/livetemplate/components
# New: module github.com/livetemplate/lvt/components

# 3. Keep old repo as redirect (optional)
# github.com/livetemplate/components becomes a stub:
```

```go
// github.com/livetemplate/components (stub repo)
// go.mod
module github.com/livetemplate/components

// Deprecation notice + redirect
// All packages re-export from new location
```

```go
// github.com/livetemplate/components/modal/modal.go (stub)
package modal

import lvtmodal "github.com/livetemplate/lvt/components/modal"

// Re-export all types and functions
type State = lvtmodal.State
type ConfirmState = lvtmodal.ConfirmState
var New = lvtmodal.New
var NewConfirm = lvtmodal.NewConfirm
// ... etc
```

This provides backward compatibility for any existing users while consolidating development.

---

## Integration Architecture (Monorepo Version)

### Repository Structure

```
+-------------------------------------------------------------------------+
|                     github.com/livetemplate/lvt                          |
+-------------------------------------------------------------------------+
|                                                                         |
|  +-------------------------------------------------------------------+  |
|  |                     Evolution System                               |  |
|  |  - Telemetry capture                                               |  |
|  |  - Fix proposal (single repo - much simpler!)                      |  |
|  |  - All fixes in same PR                                            |  |
|  +-------------------------------------------------------------------+  |
|                                                                         |
|  +-------------------------------------------------------------------+  |
|  |                        components/                                 |  |
|  |  +----------+ +----------+ +----------+ +----------+              |  |
|  |  |  modal   | |  toast   | | dropdown | |  toggle  |  ...         |  |
|  |  +----------+ +----------+ +----------+ +----------+              |  |
|  |                                                                   |  |
|  |  +------------------------------------------------------+         |  |
|  |  |  styles/  (adapters: tailwind, bootstrap, etc)       |         |  |
|  |  +------------------------------------------------------+         |  |
|  |                                                                   |  |
|  |  go.mod: github.com/livetemplate/lvt/components                   |  |
|  |  (independent module - no lvt dependency)                         |  |
|  +-------------------------------------------------------------------+  |
|                                                                         |
|  +-------------------+  +-------------------+  +-------------------+    |
|  |   Kit: multi      |  |   Kit: single     |  |   Kit: simple     |    |
|  |  Uses components  |  |  Uses components  |  |  Uses components  |    |
|  +-------------------+  +-------------------+  +-------------------+    |
|                                                                         |
+-------------------------------------------------------------------------+
                                   |
                                   v
                    +-----------------------------+
                    |  livetemplate/livetemplate  |
                    |      (core framework)       |
                    |    (remains separate)       |
                    +-----------------------------+
```

### Key Principle: One-Way Dependency (Still Applies)

```
lvt (CLI) --> lvt/components --> livetemplate

lvt/components does NOT depend on lvt (CLI code)
livetemplate does NOT depend on lvt or lvt/components
```

This is enforced by the nested go.mod - components cannot import from parent.

### Independence Verification (CI)

```yaml
# .github/workflows/components-independence.yml
name: Verify Components Independence

on:
  push:
    paths:
      - 'components/**'
  pull_request:
    paths:
      - 'components/**'

jobs:
  check-independence:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check for lvt imports
        run: |
          cd components
          # Should not import anything from parent lvt module
          if grep -r '"github.com/livetemplate/lvt"' --include="*.go" .; then
            echo "ERROR: components must not import lvt CLI"
            exit 1
          fi
          if grep -r '"github.com/livetemplate/lvt/internal' --include="*.go" .; then
            echo "ERROR: components must not import lvt internals"
            exit 1
          fi

      - name: Build components standalone
        run: |
          cd components
          go build ./...

      - name: Test components standalone
        run: |
          cd components
          go test ./...

      - name: Verify go.mod independence
        run: |
          cd components
          # go.mod should only depend on livetemplate, not lvt
          if grep "github.com/livetemplate/lvt[^/]" go.mod; then
            echo "ERROR: components/go.mod must not depend on lvt"
            exit 1
          fi
```

---

## Benefits of Monorepo for Evolution System

### Simplified Fix Pipeline

**Multi-repo (complex):**
```
1. Detect bug in modal component
2. Create PR in components repo
3. Wait for components CI
4. Merge components PR
5. Release new components version
6. Create PR in lvt to bump version
7. Wait for lvt CI
8. Merge lvt PR
```

**Monorepo (simple):**
```
1. Detect bug in modal component
2. Create single PR fixing components/ and updating kits
3. Wait for CI
4. Merge
```

### Atomic Changes

```go
// Single commit can include:
// - Fix in components/modal/modal.go
// - Update in internal/kits/multi/templates/...
// - New test in e2e/modal_test.go

git commit -m "fix(modal): resolve state sync issue

- Fixed state persistence in modal component
- Updated multi kit template to use new API
- Added regression test"
```

### Simpler Evolution System

```go
// Evolution system only needs to work with one repo
type FixProposer struct {
    repo *git.Repository  // Just one repo!
}

func (p *FixProposer) ProposeFix(err GenerationError) *Fix {
    // Determine if fix is in components/ or internal/
    if isComponentError(err) {
        return &Fix{
            File: "components/modal/modal.go",
            // ...
        }
    }
    return &Fix{
        File: "internal/kits/multi/templates/...",
        // ...
    }
}

// No need for cross-repo coordination!
```

---

## Integration Levels

### Level 1: Generated Apps Use Components (Immediate)

Generated apps import and use components directly:

```go
// Generated main.go
import (
    "github.com/livetemplate/lvt/components/modal"
    "github.com/livetemplate/lvt/components/toast"
    "github.com/livetemplate/livetemplate"
)

type State struct {
    // App-specific state
    Items []Item

    // Component state (standardized)
    DeleteConfirm *modal.ConfirmState
    Toasts        *toast.ContainerState
}

func NewState() *State {
    return &State{
        DeleteConfirm: modal.NewConfirm("delete-confirm",
            modal.WithConfirmTitle("Delete Item"),
            modal.WithConfirmMessage("Are you sure?"),
        ),
        Toasts: toast.NewContainer("toasts",
            toast.WithPosition(toast.TopRight),
        ),
    }
}
```

**Benefits:**
- Immediate usage of components
- Standardized patterns across all generated apps
- Component bugs surface quickly

### Level 2: lvt Templates Reference Component Templates

lvt kit templates use component template names:

```html
<!-- internal/kits/system/multi/templates/resource/template.tmpl.tmpl -->

{{define "content"}}
<div class="resource-page">
    <!-- Use component's versioned template -->
    {{template "lvt:modal:confirm:v1" .DeleteConfirm}}
    {{template "lvt:toast:container:v1" .Toasts}}

    <!-- Kit-specific layout -->
    <div class="resource-list">
        {{range .Items}}
            <div class="item">{{.Name}}</div>
        {{end}}
    </div>
</div>
{{end}}
```

**Benefits:**
- Single source of truth for modal/toast/etc rendering
- Version pinning via template names
- Kit templates become simpler (just layout + component references)

### Level 3: lvt Internalizes Component Registration

lvt handles component template registration automatically:

```go
// internal/generator/app.go

func generateMain(config *Config) string {
    // Determine which components are used
    components := detectUsedComponents(config)

    return fmt.Sprintf(`
func main() {
    tmpl := livetemplate.NewTemplate()

    // Register component templates (auto-generated based on usage)
    %s

    // ... rest of app
}
`, generateComponentRegistration(components))
}

func generateComponentRegistration(components []string) string {
    var registrations []string
    for _, c := range components {
        registrations = append(registrations,
            fmt.Sprintf("livetemplate.RegisterComponentTemplates(tmpl, %s.Templates())", c))
    }
    return strings.Join(registrations, "\n    ")
}
```

---

## Component Adoption Priority

Based on git history analysis, prioritize components that address most common issues:

| Priority | Component | Addresses | Current Bug Rate |
|----------|-----------|-----------|------------------|
| 1 | modal | Modal state management (4+ fixes) | 40% of UI bugs |
| 2 | toast | User feedback standardization | N/A (new feature) |
| 3 | toggle | Form checkboxes | 10% of form bugs |
| 4 | dropdown | Select synchronization (2+ fixes) | 30% of form bugs |
| 5 | data-table | List rendering | 20% of display bugs |

---

## CSS Styling Architecture

### The Problem

Currently, styling is tightly coupled:
- Components ship with Tailwind classes baked in
- Kits have CSS helpers specific to their framework
- Changing styles requires template modifications
- No clean separation between structure and presentation

### Design Goals

1. **Swappable at generation time** - Choose Tailwind, Bootstrap, vanilla CSS, or custom
2. **Swappable post-generation** - Change styles without touching Go code
3. **Consistent between kits and components** - Same styling system for both
4. **Unstyled option** - Pure semantic HTML for custom styling

### Proposed Architecture: Style Adapters

```
+-------------------------------------------------------------------------+
|                         Style Adapter System                             |
+-------------------------------------------------------------------------+
|                                                                         |
|  Component/Kit Template                                                 |
|  +-------------------------------------------------------------------+  |
|  |  <button class="{{styles.Button.Primary}}">Save</button>          |  |
|  |  <div class="{{styles.Modal.Container}}">...</div>                |  |
|  +-------------------------------------------------------------------+  |
|                              |                                          |
|                              v                                          |
|  +-------------------------------------------------------------------+  |
|  |                    Style Adapter Interface                         |  |
|  |  type StyleAdapter interface {                                     |  |
|  |      Button() ButtonStyles                                         |  |
|  |      Modal() ModalStyles                                           |  |
|  |      Form() FormStyles                                             |  |
|  |      // ...                                                        |  |
|  |  }                                                                 |  |
|  +-------------------------------------------------------------------+  |
|                              |                                          |
|          +-------------------+-------------------+                      |
|          v                   v                   v                      |
|  +---------------+  +---------------+  +---------------+               |
|  |   Tailwind    |  |   Bootstrap   |  |   Unstyled    |               |
|  |    Adapter    |  |    Adapter    |  |    Adapter    |               |
|  +---------------+  +---------------+  +---------------+               |
|                                                                         |
+-------------------------------------------------------------------------+
```

### Style Adapter Interface

```go
// styles/adapter.go

package styles

// StyleAdapter provides CSS classes for components
type StyleAdapter interface {
    Name() string  // "tailwind", "bootstrap", "unstyled"

    // Component styles
    Button() ButtonStyles
    Modal() ModalStyles
    Form() FormStyles
    Table() TableStyles
    Toast() ToastStyles
    Dropdown() DropdownStyles

    // Layout styles
    Container() ContainerStyles
    Grid() GridStyles
    Flex() FlexStyles
}

type ButtonStyles struct {
    Base      string
    Primary   string
    Secondary string
    Danger    string
    Ghost     string
    Disabled  string
    Loading   string
    Sizes     SizeVariants
}

type ModalStyles struct {
    Overlay   string
    Container string
    Header    string
    Body      string
    Footer    string
    Close     string
    Animation string
}

type SizeVariants struct {
    XS string
    SM string
    MD string
    LG string
    XL string
}
```

### Tailwind Adapter Implementation

```go
// styles/tailwind/adapter.go

package tailwind

type Adapter struct{}

func (a *Adapter) Name() string { return "tailwind" }

func (a *Adapter) Button() styles.ButtonStyles {
    return styles.ButtonStyles{
        Base:      "inline-flex items-center justify-center rounded-md font-medium transition-colors focus-visible:outline-none focus-visible:ring-2",
        Primary:   "bg-blue-600 text-white hover:bg-blue-700",
        Secondary: "bg-gray-100 text-gray-900 hover:bg-gray-200",
        Danger:    "bg-red-600 text-white hover:bg-red-700",
        Ghost:     "hover:bg-gray-100",
        Disabled:  "opacity-50 cursor-not-allowed",
        Loading:   "opacity-70 cursor-wait",
        Sizes: styles.SizeVariants{
            XS: "h-7 px-2 text-xs",
            SM: "h-8 px-3 text-sm",
            MD: "h-9 px-4 text-sm",
            LG: "h-10 px-6 text-base",
            XL: "h-12 px-8 text-lg",
        },
    }
}

func (a *Adapter) Modal() styles.ModalStyles {
    return styles.ModalStyles{
        Overlay:   "fixed inset-0 bg-black/50 flex items-center justify-center z-50",
        Container: "bg-white rounded-lg shadow-xl max-w-md w-full mx-4",
        Header:    "px-6 py-4 border-b",
        Body:      "px-6 py-4",
        Footer:    "px-6 py-4 border-t flex justify-end gap-2",
        Close:     "absolute top-4 right-4 text-gray-400 hover:text-gray-600",
        Animation: "animate-fade-in",
    }
}
```

### Unstyled Adapter Implementation

```go
// styles/unstyled/adapter.go

package unstyled

type Adapter struct{}

func (a *Adapter) Name() string { return "unstyled" }

func (a *Adapter) Button() styles.ButtonStyles {
    return styles.ButtonStyles{
        Base:      "btn",
        Primary:   "btn-primary",
        Secondary: "btn-secondary",
        Danger:    "btn-danger",
        Ghost:     "btn-ghost",
        Disabled:  "btn-disabled",
        Loading:   "btn-loading",
        Sizes: styles.SizeVariants{
            XS: "btn-xs",
            SM: "btn-sm",
            MD: "btn-md",
            LG: "btn-lg",
            XL: "btn-xl",
        },
    }
}
```

### Component Template Integration

Components use style references, not hardcoded classes:

```html
<!-- components/modal/templates/confirm.tmpl -->
{{define "lvt:modal:confirm:v1"}}
<div class="{{.Styles.Modal.Overlay}}" {{if not .Open}}hidden{{end}}>
    <div class="{{.Styles.Modal.Container}}">
        <div class="{{.Styles.Modal.Header}}">
            <h3>{{.Title}}</h3>
        </div>
        <div class="{{.Styles.Modal.Body}}">
            {{.Message}}
        </div>
        <div class="{{.Styles.Modal.Footer}}">
            <button class="{{.Styles.Button.Secondary}} {{.Styles.Button.Sizes.MD}}"
                    lvt-click="{{.CancelAction}}">
                Cancel
            </button>
            <button class="{{.Styles.Button.Danger}} {{.Styles.Button.Sizes.MD}}"
                    lvt-click="{{.ConfirmAction}}">
                {{.ConfirmText}}
            </button>
        </div>
    </div>
</div>
{{end}}
```

### Generation Time Style Selection

```bash
# Choose style at generation time
lvt new myapp --styles=tailwind    # Default
lvt new myapp --styles=bootstrap
lvt new myapp --styles=unstyled
lvt new myapp --styles=custom:./mystyles

# Change styles post-generation
lvt styles set bootstrap
lvt styles eject  # Dumps current adapter for customization
```

### Custom Style Adapter

Users can create custom adapters:

```go
// mystyles/adapter.go
package mystyles

import "github.com/livetemplate/lvt/components/styles"

type Adapter struct{}

func (a *Adapter) Name() string { return "mystyles" }

func (a *Adapter) Button() styles.ButtonStyles {
    return styles.ButtonStyles{
        Base:    "my-btn",
        Primary: "my-btn--primary",
        // ... custom classes
    }
}
```

Register in app:

```go
// main.go
import "myapp/mystyles"

func main() {
    styles.Register(&mystyles.Adapter{})
    // ...
}
```

### Kit Template Integration

Kits use the same style adapter system:

```html
<!-- internal/kits/system/multi/templates/resource/template.tmpl.tmpl -->
{{define "content"}}
<div class="{{.Styles.Container.Default}}">
    <div class="{{.Styles.Flex.Between}} {{.Styles.Flex.ItemsCenter}}">
        <h1>{{.ResourceNamePlural}}</h1>
        <button class="{{.Styles.Button.Primary}} {{.Styles.Button.Sizes.MD}}"
                lvt-click="show_add">
            Add {{.ResourceName}}
        </button>
    </div>

    <!-- Table uses style adapter -->
    <table class="{{.Styles.Table.Container}}">
        <thead class="{{.Styles.Table.Header}}">
            <!-- ... -->
        </thead>
        <tbody class="{{.Styles.Table.Body}}">
            {{range .Items}}
            <tr class="{{$.Styles.Table.Row}}">
                <!-- ... -->
            </tr>
            {{end}}
        </tbody>
    </table>
</div>
{{end}}
```

### Style Inheritance Chain

```
1. Component Default Style
        |
        v
2. Kit Style Override (optional)
        |
        v
3. App Style Override (optional)
        |
        v
4. Per-Instance Override (optional)
```

```go
// Component with style override
modal.NewConfirm("delete",
    modal.WithStyles(customStyles),  // Override for this instance
)

// Or set app-wide
styles.SetDefault(&bootstrap.Adapter{})  // All components use Bootstrap
```

### CSS File Generation

For unstyled mode, generate a CSS scaffold:

```bash
lvt styles scaffold > styles.css
```

Generates:

```css
/* styles.css - Generated scaffold for unstyled mode */

/* Buttons */
.btn { /* Add your button styles */ }
.btn-primary { /* Add primary variant */ }
.btn-secondary { /* Add secondary variant */ }
.btn-danger { /* Add danger variant */ }

/* Modals */
.modal-overlay { /* Add overlay styles */ }
.modal-container { /* Add container styles */ }
/* ... */
```

### Benefits of Style Adapter System

1. **Clean separation** - Structure in templates, presentation in adapters
2. **Swappable** - Change framework without touching templates
3. **Customizable** - Create brand-specific adapters
4. **Testable** - Test components with unstyled adapter (no class noise)
5. **Consistent** - Same system for components and kits
6. **Evolvable** - Update styles without changing templates

### Evolution System Integration

Track style-related issues:

```go
type GenerationEvent struct {
    // ... existing fields ...

    StyleAdapter string  // Which adapter was used
    StyleErrors  []StyleError
}

type StyleError struct {
    Component string
    Property  string  // "Button.Primary", "Modal.Overlay"
    Issue     string  // "class conflict", "missing class"
}
```

Evolution system can propose style adapter fixes:

```go
func (p *FixProposer) proposeStyleFix(err StyleError) *Fix {
    return &Fix{
        File:     fmt.Sprintf("components/styles/%s/adapter.go", err.StyleAdapter),
        Section:  err.Property,
        // ... fix details
    }
}
```

---

## Conclusion

The **monorepo approach** (components inside lvt) is recommended because:

1. **Single feedback loop** - Evolution system only needs to work with one repo
2. **Atomic changes** - Component fix + kit update in same commit
3. **Faster iteration** - No cross-repo version coordination
4. **Independence maintained** - Nested go.mod ensures components stay standalone
5. **Simpler CI/CD** - One pipeline tests everything together

The key insight is that Go's module system allows us to have the best of both worlds: a single repo for development velocity, with independent importability for external users.

Components remain usable by anyone via:
```go
import "github.com/livetemplate/lvt/components/modal"
```

And CI ensures they never depend on lvt internals, preserving their standalone nature.
