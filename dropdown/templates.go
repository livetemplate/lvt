package dropdown

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all dropdown template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the dropdown component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/dropdown"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(dropdown.Templates()),
//	)
//
// Available templates:
//   - "lvt:dropdown:default:v1"     - Basic single-select dropdown
//   - "lvt:dropdown:searchable:v1"  - Searchable dropdown with filter input
//   - "lvt:dropdown:multi:v1"       - Multi-select dropdown with checkboxes
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "dropdown")
}
