package tooltip

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all tooltip template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the tooltip component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/tooltip"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(tooltip.Templates()),
//	)
//
// Available templates:
//   - "lvt:tooltip:default:v1" - Basic tooltip
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "tooltip")
}
