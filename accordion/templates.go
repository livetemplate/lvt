package accordion

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all accordion template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the accordion component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/accordion"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(accordion.Templates()),
//	)
//
// Available templates:
//   - "lvt:accordion:default:v1" - Multi-open accordion
//   - "lvt:accordion:single:v1"  - Single-open accordion
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "accordion")
}
