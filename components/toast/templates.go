package toast

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all toast template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the toast component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/toast"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(toast.Templates()),
//	)
//
// Available templates:
//   - "lvt:toast:container:v1" - Toast container with position
//   - "lvt:toast:message:v1"   - Individual toast message
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "toast")
}
