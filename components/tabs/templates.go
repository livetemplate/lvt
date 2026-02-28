package tabs

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all tabs template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the tabs component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/tabs"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(tabs.Templates()),
//	)
//
// Available templates:
//   - "lvt:tabs:horizontal:v1" - Horizontal tabs (default)
//   - "lvt:tabs:vertical:v1"   - Vertical tabs (sidebar style)
//   - "lvt:tabs:pills:v1"      - Pill-style tabs
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "tabs")
}
