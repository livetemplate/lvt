package timepicker

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all timepicker template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the timepicker component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/timepicker"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(timepicker.Templates()),
//	)
//
// Available templates:
//   - "lvt:timepicker:default:v1"  - Time picker
//   - "lvt:timepicker:duration:v1" - Duration picker
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "timepicker")
}
