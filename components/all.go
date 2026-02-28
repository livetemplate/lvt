// Package components provides pre-built UI components for the LiveTemplate framework.
//
// This library provides reusable, styled components that can be registered with
// LiveTemplate applications using the WithComponentTemplates() option.
//
// # Quick Start
//
// Register all component templates in your main.go:
//
//	import "github.com/livetemplate/components"
//
//	func main() {
//	    tmpl := livetemplate.NewTemplates(
//	        livetemplate.WithComponentTemplates(components.All()),
//	        livetemplate.ParseGlob("internal/app/**/*.tmpl"),
//	    )
//	}
//
// Then use components in your templates:
//
//	{{template "lvt:dropdown:searchable:v1" .MyDropdown}}
//
// # Available Components
//
// See the component-specific packages for detailed documentation:
//
//   - dropdown: Dropdown menus (default, searchable, multi-select)
//   - tabs: Tab navigation and content panels
//   - modal: Modal dialogs and sheets
//   - toast: Toast notifications
//   - accordion: Collapsible sections
//   - datepicker: Date/time selection
//   - rating: Star rating
//   - And more...
//
// # Customization
//
// Components use Tailwind CSS by default. Override by defining the same template
// name in your project - project templates take precedence over library templates.
//
// # Contributing
//
// See the repository README for contribution guidelines.
package components

import (
	"github.com/livetemplate/components/accordion"
	"github.com/livetemplate/components/breadcrumbs"
	"github.com/livetemplate/components/autocomplete"
	"github.com/livetemplate/components/base"
	"github.com/livetemplate/components/datatable"
	"github.com/livetemplate/components/datepicker"
	"github.com/livetemplate/components/drawer"
	"github.com/livetemplate/components/dropdown"
	"github.com/livetemplate/components/menu"
	"github.com/livetemplate/components/modal"
	"github.com/livetemplate/components/popover"
	"github.com/livetemplate/components/progress"
	"github.com/livetemplate/components/rating"
	"github.com/livetemplate/components/skeleton"
	"github.com/livetemplate/components/tabs"
	"github.com/livetemplate/components/tagsinput"
	"github.com/livetemplate/components/timeline"
	"github.com/livetemplate/components/timepicker"
	"github.com/livetemplate/components/toast"
	"github.com/livetemplate/components/toggle"
	"github.com/livetemplate/components/tooltip"
)

// All returns all component template sets for registration with LiveTemplate.
//
// Usage:
//
//	tmpl := livetemplate.NewTemplates(
//	    livetemplate.WithComponentTemplates(components.All()),
//	)
//
// This registers all component templates, allowing you to use any component
// in your templates immediately.
func All() []*base.TemplateSet {
	return []*base.TemplateSet{
		accordion.Templates(),
		breadcrumbs.Templates(),
		autocomplete.Templates(),
		datatable.Templates(),
		datepicker.Templates(),
		drawer.Templates(),
		dropdown.Templates(),
		menu.Templates(),
		modal.Templates(),
		popover.Templates(),
		progress.Templates(),
		rating.Templates(),
		skeleton.Templates(),
		tabs.Templates(),
		tagsinput.Templates(),
		timeline.Templates(),
		timepicker.Templates(),
		toast.Templates(),
		toggle.Templates(),
		tooltip.Templates(),
	}
}

// Version returns the library version.
func Version() string {
	return "0.1.0"
}
