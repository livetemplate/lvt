package drawer

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all drawer template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the drawer component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/drawer"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(drawer.Templates()),
//	)
//
// Available templates:
//   - "lvt:drawer:default:v1" - Slide-out drawer/sidebar
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "drawer")
}
