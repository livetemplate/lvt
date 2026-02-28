package datepicker

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all datepicker template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the datepicker component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/datepicker"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(datepicker.Templates()),
//	)
//
// Available templates:
//   - "lvt:datepicker:single:v1" - Single date picker
//   - "lvt:datepicker:range:v1"  - Date range picker
//   - "lvt:datepicker:inline:v1" - Inline calendar (always visible)
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "datepicker")
}
