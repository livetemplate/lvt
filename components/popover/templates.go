package popover

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all popover template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the popover component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/popover"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(popover.Templates()),
//	)
//
// Available templates:
//   - "lvt:popover:default:v1" - Rich content popover
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "popover")
}
