# LVT + Components Integration Strategy

**Goal:** Use lvt as the primary driver for components evolution while maintaining component independence.

## The Bootstrap Problem

```
Components need usage to mature → But they're too unstable to use → So they don't get usage → They stay unstable
```

**Solution:** lvt becomes the forcing function. By making lvt depend on and exercise components, we create:
1. Real-world usage patterns
2. Bug discovery through generation
3. Evolution feedback loop
4. Path to v1.0 stability

## Integration Architecture

### Dependency Graph

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           USER / LLM                                     │
└─────────────────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              lvt CLI                                     │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Evolution System                             │   │
│  │  - Telemetry capture                                             │   │
│  │  - Fix proposal (for lvt AND components)                         │   │
│  │  - Multi-repo PR creation                                        │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐  │
│  │   Kit: multi      │  │   Kit: single     │  │   Kit: simple     │  │
│  │                   │  │                   │  │                   │  │
│  │ Uses components:  │  │ Uses components:  │  │ Uses components:  │
│  │ - modal           │  │ - modal           │  │ - (minimal)       │
│  │ - toast           │  │ - toast           │  │                   │
│  │ - dropdown        │  │ - toggle          │  │                   │
│  │ - data-table      │  │                   │  │                   │
│  └───────────────────┘  └───────────────────┘  └───────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
                                   │
                    ┌──────────────┴──────────────┐
                    ▼                              ▼
┌─────────────────────────────┐    ┌─────────────────────────────┐
│   livetemplate/components   │    │    livetemplate/livetemplate │
│                             │    │                              │
│  - Remains independent      │    │  - Core framework            │
│  - No lvt dependency        │    │  - Components depend on this │
│  - Can be used standalone   │    │                              │
└─────────────────────────────┘    └─────────────────────────────┘
```

### Key Principle: One-Way Dependency

```
lvt → components → livetemplate

components does NOT depend on lvt
livetemplate does NOT depend on components or lvt
```

This ensures components remain independently usable.

## Integration Levels

### Level 1: Generated Apps Use Components (Immediate)

Generated apps import and use components directly:

```go
// Generated main.go
import (
    "github.com/livetemplate/components/modal"
    "github.com/livetemplate/components/toast"
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

## Evolution System Integration

### Cross-Repo Telemetry

When a generated app fails, capture which component was involved:

```go
type GenerationError struct {
    Phase       string
    File        string
    Message     string

    // NEW: Component attribution
    Component   string  // e.g., "modal", "toast", ""
    ComponentVersion string
}

func attributeErrorToComponent(err error, files []string) string {
    // Analyze error and generated files to determine if a component caused it
    for _, f := range files {
        content := readFile(f)
        if strings.Contains(content, "lvt:modal") && isModalRelated(err) {
            return "modal"
        }
        // ... other components
    }
    return ""
}
```

### Cross-Repo Fix Proposal

When evolution system identifies a component bug:

```go
func (p *FixProposer) proposeComponentFix(err GenerationError) *CrossRepoFix {
    if err.Component == "" {
        return nil  // Not a component issue
    }

    return &CrossRepoFix{
        PrimaryRepo: "livetemplate/components",
        PrimaryFix: Fix{
            File:    fmt.Sprintf("%s/templates.go", err.Component),
            // ... fix details
        },

        // Also update lvt to use fixed version
        SecondaryRepo: "livetemplate/lvt",
        SecondaryFix: Fix{
            File: "go.mod",
            Find: fmt.Sprintf("github.com/livetemplate/components v%s", currentVersion),
            Replace: fmt.Sprintf("github.com/livetemplate/components v%s", newVersion),
        },

        // Link the PRs
        LinkedPRs: true,
    }
}
```

### Component Quality Dashboard

Track component health through lvt usage:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     Component Health Dashboard                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Component    │ Usage (30d) │ Success Rate │ Common Errors │ Status    │
│  ─────────────┼─────────────┼──────────────┼───────────────┼────────── │
│  modal        │    847      │    94.2%     │ State sync    │ ⚠️ Needs  │
│  toast        │    623      │    98.1%     │ Position      │ ✅ Stable │
│  dropdown     │    412      │    87.3%     │ Selection     │ 🔴 Broken │
│  data-table   │    156      │    91.4%     │ Pagination    │ ⚠️ Needs  │
│  toggle       │    534      │    99.2%     │ (none)        │ ✅ Stable │
│                                                                         │
│  Recommended actions:                                                   │
│  1. Fix dropdown selection bug (3 PRs proposed)                        │
│  2. Investigate modal state sync issue                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Keeping Components Independent

### Design Constraints

1. **No lvt imports in components:**
   ```go
   // components/modal/modal.go
   import (
       "github.com/livetemplate/livetemplate"  // OK - core framework
       // NO: "github.com/livetemplate/lvt/..."  // NEVER
   )
   ```

2. **Components work without lvt:**
   ```go
   // Standalone usage (no lvt involved)
   import "github.com/livetemplate/components/modal"

   state := modal.NewConfirm("my-modal", modal.WithTitle("Hello"))
   tmpl.RegisterComponentTemplates(modal.Templates())
   ```

3. **Components have own test suite:**
   ```
   components/
   ├── modal/
   │   ├── modal.go
   │   ├── modal_test.go      # Unit tests
   │   └── modal_e2e_test.go  # Standalone E2E tests
   ```

4. **Components have own examples:**
   ```
   components/
   ├── examples/
   │   ├── modal-basic/       # Works without lvt
   │   ├── toast-positions/
   │   └── full-app/          # Shows all components together
   ```

### Independence Verification

Add CI check to components repo:

```yaml
# .github/workflows/independence.yml
name: Verify Independence

on: [push, pull_request]

jobs:
  no-lvt-dependency:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check for lvt imports
        run: |
          if grep -r "github.com/livetemplate/lvt" --include="*.go" .; then
            echo "ERROR: Components must not import lvt"
            exit 1
          fi

      - name: Build without lvt
        run: go build ./...

      - name: Test without lvt
        run: go test ./...
```

## Implementation Phases

### Phase 1: Add Components Dependency (Week 1)

```bash
# In lvt repo
go get github.com/livetemplate/components@latest
```

Update `go.mod`:
```go
require (
    github.com/livetemplate/livetemplate v0.8.0
    github.com/livetemplate/components v0.1.0  // NEW
)
```

### Phase 2: Integrate Modal Component (Week 2)

Replace kit modal templates with component usage:

**Before (3 separate implementations):**
```
internal/kits/system/multi/templates/components/modal.tmpl
internal/kits/system/single/templates/components/modal.tmpl
internal/kits/system/simple/templates/components/modal.tmpl  (missing)
```

**After (single component):**
```go
// Generated handler.go
import "github.com/livetemplate/components/modal"

type State struct {
    EditModal   *modal.State
    DeleteModal *modal.ConfirmState
}
```

```html
<!-- Generated template.tmpl -->
{{template "lvt:modal:default:v1" .EditModal}}
{{template "lvt:modal:confirm:v1" .DeleteModal}}
```

### Phase 3: Integrate Toast Component (Week 3)

Add standardized notifications to all generated apps:

```go
type State struct {
    Toasts *toast.ContainerState
}

func (c *Controller) Create(ctx context.Context, s *State) *State {
    // ... create logic ...
    s.Toasts.Add(toast.Success("Item created successfully"))
    return s
}
```

### Phase 4: Evolution System for Components (Week 4)

Extend telemetry to track component usage and failures:

```go
type GenerationEvent struct {
    // ... existing fields ...

    ComponentsUsed    []ComponentUsage
    ComponentErrors   []ComponentError
}

type ComponentUsage struct {
    Name    string  // "modal", "toast"
    Version string  // "v1"
    Count   int     // How many instances
}

type ComponentError struct {
    Component string
    Version   string
    Error     string
    Context   string
}
```

### Phase 5: Cross-Repo Fix Pipeline (Week 5-6)

```go
type CrossRepoFixPipeline struct {
    lvtRepo        *GitRepo
    componentsRepo *GitRepo
}

func (p *CrossRepoFixPipeline) ProposeLinkedFixes(event *GenerationEvent) error {
    // 1. Identify if issue is in components
    componentIssues := filterComponentIssues(event.Errors)

    for _, issue := range componentIssues {
        // 2. Propose fix to components repo
        componentFix := p.proposeComponentFix(issue)
        componentPR := p.componentsRepo.CreatePR(componentFix)

        // 3. Propose version bump to lvt repo
        versionBump := p.proposeVersionBump(componentPR)
        lvtPR := p.lvtRepo.CreatePR(versionBump)

        // 4. Link the PRs
        p.linkPRs(componentPR, lvtPR)
    }

    return nil
}
```

## Component Adoption Priority

Based on git history analysis, prioritize components that address most common issues:

| Priority | Component | Addresses | Current Bug Rate |
|----------|-----------|-----------|------------------|
| 1 | modal | Modal state management (4+ fixes) | 40% of UI bugs |
| 2 | toast | User feedback standardization | N/A (new feature) |
| 3 | toggle | Form checkboxes | 10% of form bugs |
| 4 | dropdown | Select synchronization (2+ fixes) | 30% of form bugs |
| 5 | data-table | List rendering | 20% of display bugs |

## Success Metrics

### For lvt

| Metric | Current | Target (3mo) | Target (6mo) |
|--------|---------|--------------|--------------|
| Generation success rate | ~75% | 90% | 95% |
| Compilation success rate | Unknown | 95% | 99% |
| Template files maintained | 15+ | 5 | 3 |

### For Components

| Metric | Current | Target (3mo) | Target (6mo) |
|--------|---------|--------------|--------------|
| Components at v1 | 0 | 5 | 15 |
| Test coverage | Unknown | 80% | 95% |
| lvt-driven bug fixes | 0 | 20 | 50 |
| Standalone users | 1 (examples) | 10 | 50 |

## Risk Mitigation

### Risk: Component Breaking Change Breaks lvt

**Mitigation:** Version pinning in templates

```html
{{template "lvt:modal:confirm:v1" .Modal}}  <!-- Pinned to v1 -->
```

When component releases v2, lvt continues using v1 until explicitly upgraded.

### Risk: Components Become lvt-Specific

**Mitigation:** Independence tests + standalone examples

```yaml
# Components CI must pass:
- Build without lvt
- Test without lvt
- Example apps work without lvt
```

### Risk: Two Repos Slow Development

**Mitigation:** Automated cross-repo tooling

```bash
# Single command to test change across both repos
lvt dev test-with-components --component-branch=fix-modal-bug

# Runs:
# 1. Build components from branch
# 2. Build lvt with local components
# 3. Run full lvt test suite
# 4. Report results
```

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
┌─────────────────────────────────────────────────────────────────────────┐
│                         Style Adapter System                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Component/Kit Template                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  <button class="{{styles.Button.Primary}}">Save</button>        │   │
│  │  <div class="{{styles.Modal.Container}}">...</div>              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Style Adapter Interface                       │   │
│  │  type StyleAdapter interface {                                   │   │
│  │      Button() ButtonStyles                                       │   │
│  │      Modal() ModalStyles                                         │   │
│  │      Form() FormStyles                                           │   │
│  │      // ...                                                      │   │
│  │  }                                                               │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│          ┌───────────────────┼───────────────────┐                     │
│          ▼                   ▼                   ▼                      │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐              │
│  │   Tailwind    │  │   Bootstrap   │  │   Unstyled    │              │
│  │    Adapter    │  │    Adapter    │  │    Adapter    │              │
│  └───────────────┘  └───────────────┘  └───────────────┘              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
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

import "github.com/livetemplate/components/styles"

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
        ↓
2. Kit Style Override (optional)
        ↓
3. App Style Override (optional)
        ↓
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
        Repo:     "livetemplate/components",
        File:     fmt.Sprintf("styles/%s/adapter.go", err.StyleAdapter),
        Section:  err.Property,
        // ... fix details
    }
}
```

## Conclusion

Tight integration with components is the right path forward because:

1. **lvt needs stable components** - Current template drift proves we need single sources of truth
2. **Components need real usage** - lvt provides thousands of generated apps as test cases
3. **Evolution system can span both** - Unified feedback loop improves everything
4. **Independence is maintainable** - One-way dependency + CI checks ensure components stay standalone

The key insight is that lvt becomes the **forcing function** for component maturity, while the evolution system provides the **feedback loop** that drives improvement in both.
