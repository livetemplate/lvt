package rating

import (
	"embed"

	"github.com/livetemplate/components/base"
)

// templateFS contains all rating template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the rating component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/rating"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(rating.Templates()),
//	)
//
// Available templates:
//   - "lvt:rating:default:v1"  - Interactive star rating
//   - "lvt:rating:readonly:v1" - Read-only rating display
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "rating")
}
