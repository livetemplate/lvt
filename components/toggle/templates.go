package toggle

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all toggle template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the toggle component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/toggle"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(toggle.Templates()),
//	)
//
// Available templates:
//   - "lvt:toggle:default:v1" - Switch toggle
//   - "lvt:toggle:checkbox:v1" - Styled checkbox
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "toggle")
}
