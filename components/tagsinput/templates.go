package tagsinput

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all tagsinput template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the tagsinput component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/tagsinput"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(tagsinput.Templates()),
//	)
//
// Available templates:
//   - "lvt:tagsinput:default:v1" - Standard tags input
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "tagsinput")
}
