package breadcrumbs

import (
	"embed"
	"html/template"

	"github.com/livetemplate/lvt/components/base"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the breadcrumbs template set.
func Templates() *base.TemplateSet {
	return base.WithFuncs(base.NewTemplateSet(templateFS, "templates/*.tmpl", "breadcrumbs"),
		template.FuncMap{
			"sub": func(a, b int) int { return a - b },
		})
}
