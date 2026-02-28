package menu

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all menu template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the menu component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/menu"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(menu.Templates()),
//	)
//
// Available templates:
//   - "lvt:menu:default:v1" - Dropdown action menu
//   - "lvt:menu:context:v1" - Context/right-click menu
//   - "lvt:menu:nav:v1"     - Navigation menu
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "menu")
}
