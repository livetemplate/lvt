package timeline

import (
	"embed"
	"html/template"

	"github.com/livetemplate/lvt/components/base"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the timeline template set.
func Templates() *base.TemplateSet {
	return base.WithFuncs(
		base.NewTemplateSet(templateFS, "templates/*.tmpl", "timeline"),
		template.FuncMap{
			"sub": func(a, b int) int { return a - b },
		},
	)
}
