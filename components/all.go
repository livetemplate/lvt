// Package components provides pre-built UI components for the LiveTemplate framework.
//
// This library provides reusable, styled components that can be registered with
// LiveTemplate applications using the WithComponentTemplates() option.
//
// # Quick Start
//
// Register all component templates in your main.go:
//
//	import "github.com/livetemplate/lvt/components"
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
	"github.com/livetemplate/lvt/components/accordion"
	"github.com/livetemplate/lvt/components/autocomplete"
	"github.com/livetemplate/lvt/components/breadcrumbs"
	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/datatable"
	"github.com/livetemplate/lvt/components/datepicker"
	"github.com/livetemplate/lvt/components/drawer"
	"github.com/livetemplate/lvt/components/dropdown"
	"github.com/livetemplate/lvt/components/menu"
	"github.com/livetemplate/lvt/components/modal"
	"github.com/livetemplate/lvt/components/popover"
	"github.com/livetemplate/lvt/components/progress"
	"github.com/livetemplate/lvt/components/rating"
	"github.com/livetemplate/lvt/components/skeleton"
	"github.com/livetemplate/lvt/components/tabs"
	"github.com/livetemplate/lvt/components/tagsinput"
	"github.com/livetemplate/lvt/components/timeline"
	"github.com/livetemplate/lvt/components/timepicker"
	"github.com/livetemplate/lvt/components/toast"
	"github.com/livetemplate/lvt/components/toggle"
	"github.com/livetemplate/lvt/components/tooltip"
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
		autocomplete.Templates(),
		breadcrumbs.Templates(),
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
