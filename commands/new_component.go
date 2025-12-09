package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

// NewComponent handles the `lvt new component` command.
func NewComponent(args []string) error {
	// Handle --help flag
	if ShowHelpIfRequested(args, printNewComponentHelp) {
		return nil
	}

	if len(args) < 1 {
		printNewComponentHelp()
		return nil
	}

	componentName := args[0]

	// Validate that component name doesn't look like a flag
	if err := ValidatePositionalArg(componentName, "component name"); err != nil {
		return err
	}

	// Validate component name
	if !isValidComponentName(componentName) {
		return fmt.Errorf("invalid component name: %s\nComponent names must be lowercase alphanumeric with optional hyphens", componentName)
	}

	// Default to components/ directory
	destDir := filepath.Join("components", componentName)

	// Parse options
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--dest", "-d":
			if i+1 < len(args) {
				destDir = args[i+1]
				i++
			}
		}
	}

	// Check if destination exists
	if _, err := os.Stat(destDir); err == nil {
		return fmt.Errorf("destination already exists: %s", destDir)
	}

	// Create directories
	templatesDir := filepath.Join(destDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	// Generate files
	data := componentData{
		Name:        componentName,
		PackageName: toPackageName(componentName),
		NamePascal:  toPascalCase(componentName),
		NameCamel:   toCamelCase(componentName),
	}

	files := []struct {
		path     string
		template string
	}{
		{filepath.Join(destDir, componentName+".go"), componentGoTemplate},
		{filepath.Join(destDir, "options.go"), optionsGoTemplate},
		{filepath.Join(destDir, "templates.go"), templatesGoTemplate},
		{filepath.Join(destDir, componentName+"_test.go"), testGoTemplate},
		{filepath.Join(templatesDir, "default.tmpl"), defaultTmplTemplate},
	}

	for _, f := range files {
		if err := writeTemplateFile(f.path, f.template, data); err != nil {
			return fmt.Errorf("failed to create %s: %v", f.path, err)
		}
	}

	fmt.Printf("✅ Created component scaffold: %s\n", componentName)
	fmt.Println()
	fmt.Println("Created files:")
	for _, f := range files {
		rel, _ := filepath.Rel(".", f.path)
		fmt.Printf("  - %s\n", rel)
	}
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Edit %s/%s.go to implement your component\n", destDir, componentName)
	fmt.Printf("  2. Edit %s/templates/default.tmpl for the HTML\n", destDir)
	fmt.Printf("  3. Run tests: go test ./%s/...\n", destDir)
	fmt.Println()
	fmt.Println("To contribute to the official library:")
	fmt.Println("  - Fork github.com/livetemplate/components")
	fmt.Println("  - Move your component to the fork")
	fmt.Println("  - Submit a pull request")

	return nil
}

type componentData struct {
	Name        string // Original name with hyphens (for templates)
	PackageName string // Go package name (no hyphens)
	NamePascal  string
	NameCamel   string
}

func isValidComponentName(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		if r == '-' {
			if i == 0 || i == len(name)-1 {
				return false
			}
			continue
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
		if unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func toPascalCase(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, "")
}

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) > 0 {
		return strings.ToLower(pascal[:1]) + pascal[1:]
	}
	return pascal
}

func toPackageName(s string) string {
	// Remove hyphens for Go package names
	return strings.ReplaceAll(s, "-", "")
}

func writeTemplateFile(path, tmplText string, data componentData) error {
	tmpl, err := template.New("file").Parse(tmplText)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

const componentGoTemplate = `// Package {{.PackageName}} provides {{.Name}} components for the LiveTemplate framework.
//
// # Available Templates
//
//   - lvt:{{.Name}}:default:v1 - Standard {{.Name}}
//
// # Basic Usage
//
//	c := {{.PackageName}}.New("my{{.NamePascal}}")
//
//	{{"{{"}}template "lvt:{{.Name}}:default:v1" .{{.NamePascal}}{{"}}"}}
package {{.PackageName}}

import (
	"github.com/livetemplate/components/base"
)

// {{.NamePascal}} represents a {{.Name}} component.
type {{.NamePascal}} struct {
	base.Base

	// Add your component fields here
	// Example:
	// Label string
	// Value string
}

// Option configures a {{.NamePascal}}.
type Option func(*{{.NamePascal}})

// New creates a new {{.NamePascal}} with the given ID and options.
func New(id string, opts ...Option) *{{.NamePascal}} {
	c := &{{.NamePascal}}{
		Base: base.NewBase(id, "{{.Name}}"),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Add your component methods here
// Example:
// func (c *{{.NamePascal}}) GetLabel() string {
//     return c.Label
// }
`

const optionsGoTemplate = `package {{.PackageName}}

// Add functional options here
// Example:
// func WithLabel(label string) Option {
//     return func(c *{{.NamePascal}}) {
//         c.Label = label
//     }
// }

// WithStyled sets the styled mode.
func WithStyled(styled bool) Option {
	return func(c *{{.NamePascal}}) {
		c.SetStyled(styled)
	}
}
`

const templatesGoTemplate = `package {{.PackageName}}

import (
	"embed"

	"github.com/livetemplate/components/base"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the {{.Name}} template set.
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "{{.Name}}")
}
`

const testGoTemplate = `package {{.PackageName}}

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("creates {{.Name}} with defaults", func(t *testing.T) {
		c := New("test")
		if c.ID() != "test" {
			t.Errorf("expected ID 'test', got %q", c.ID())
		}
		if c.Namespace() != "{{.Name}}" {
			t.Errorf("expected namespace '{{.Name}}', got %q", c.Namespace())
		}
	})
}

func TestWithStyled(t *testing.T) {
	c := New("test", WithStyled(true))
	if !c.IsStyled() {
		t.Error("expected styled")
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()
	if ts == nil {
		t.Fatal("Templates() returned nil")
	}
}
`

const defaultTmplTemplate = `{{"{{"}}define "lvt:{{.Name}}:default:v1"{{"}}"}}
{{"{{"}}if .IsStyled{{"}}"}}
{{"{{"}}/* Tailwind CSS styled version */{{"}}"}}
<div class="p-4 border rounded-lg" data-component="{{.Name}}" data-id="{{"{{"}}{{".ID"}}{{"}}"}}>
  {{"{{"}}/* Add your styled HTML here */{{"}}"}}
  <p>{{.NamePascal}} component</p>
</div>
{{"{{"}}else{{"}}"}}
{{"{{"}}/* Unstyled semantic HTML version */{{"}}"}}
<div data-component="{{.Name}}" data-id="{{"{{"}}{{".ID"}}{{"}}"}}>
  {{"{{"}}/* Add your unstyled HTML here */{{"}}"}}
  <p>{{.NamePascal}} component</p>
</div>
{{"{{"}}end{{"}}"}}
{{"{{"}}end{{"}}"}}
`
