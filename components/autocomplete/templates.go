package autocomplete

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all autocomplete template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the autocomplete component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/autocomplete"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(autocomplete.Templates()),
//	)
//
// Available templates:
//   - "lvt:autocomplete:default:v1" - Basic autocomplete
//   - "lvt:autocomplete:multi:v1"   - Multi-select autocomplete
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "autocomplete")
}
