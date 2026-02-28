package skeleton

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all skeleton template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the skeleton component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/skeleton"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(skeleton.Templates()),
//	)
//
// Available templates:
//   - "lvt:skeleton:default:v1" - Basic skeleton line
//   - "lvt:skeleton:avatar:v1" - Circular avatar placeholder
//   - "lvt:skeleton:card:v1" - Card placeholder
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "skeleton")
}
