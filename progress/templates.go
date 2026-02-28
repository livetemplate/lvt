package progress

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all progress template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the progress component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/progress"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(progress.Templates()),
//	)
//
// Available templates:
//   - "lvt:progress:default:v1" - Linear progress bar
//   - "lvt:progress:circular:v1" - Circular progress indicator
//   - "lvt:progress:spinner:v1" - Loading spinner
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "progress")
}
