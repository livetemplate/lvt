package modal

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all modal template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the modal component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/modal"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(modal.Templates()),
//	)
//
// Available templates:
//   - "lvt:modal:default:v1" - Standard modal dialog
//   - "lvt:modal:confirm:v1" - Confirmation dialog
//   - "lvt:modal:sheet:v1" - Slide-in sheet panel
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "modal")
}
