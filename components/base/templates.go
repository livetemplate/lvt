package base

import (
	"embed"
	"html/template"
)

// TemplateSet represents a collection of embedded templates for a component.
// Components create TemplateSet instances to expose their templates for registration
// with the LiveTemplate framework.
//
// Example:
//
//	package dropdown
//
//	import "embed"
//
//	//go:embed templates/*.tmpl
//	var templateFS embed.FS
//
//	func Templates() *TemplateSet {
//	    return &TemplateSet{
//	        FS:        templateFS,
//	        Pattern:   "templates/*.tmpl",
//	        Namespace: "dropdown",
//	    }
//	}
type TemplateSet struct {
	// FS is the embedded filesystem containing the template files.
	FS embed.FS

	// Pattern is the glob pattern for matching template files within FS.
	// Examples: "templates/*.tmpl", "*.tmpl", "**/*.tmpl"
	Pattern string

	// Namespace identifies the component type for this template set.
	// Used for documentation and debugging purposes.
	// Example: "dropdown" for templates like "lvt:dropdown:searchable:v1"
	Namespace string

	// Funcs provides additional template functions for this component.
	// These are merged with the base template functions.
	Funcs template.FuncMap
}

// NewTemplateSet creates a new TemplateSet with the given filesystem and pattern.
//
// Example:
//
//	//go:embed templates/*.tmpl
//	var templateFS embed.FS
//
//	func Templates() *TemplateSet {
//	    return NewTemplateSet(templateFS, "templates/*.tmpl", "dropdown")
//	}
func NewTemplateSet(fs embed.FS, pattern, namespace string) *TemplateSet {
	return &TemplateSet{
		FS:        fs,
		Pattern:   pattern,
		Namespace: namespace,
	}
}

// WithFuncs returns a copy of the TemplateSet with additional template functions.
//
// Example:
//
//	func Templates() *TemplateSet {
//	    return NewTemplateSet(templateFS, "templates/*.tmpl", "dropdown").
//	        WithFuncs(template.FuncMap{
//	            "dropdownClass": func() string { return "dropdown" },
//	        })
//	}
func (ts *TemplateSet) WithFuncs(funcs template.FuncMap) *TemplateSet {
	return &TemplateSet{
		FS:        ts.FS,
		Pattern:   ts.Pattern,
		Namespace: ts.Namespace,
		Funcs:     funcs,
	}
}

// TemplateProvider is implemented by components that provide templates.
// The LiveTemplate framework uses this interface to collect all component templates.
type TemplateProvider interface {
	// Templates returns the component's template set for registration.
	Templates() *TemplateSet
}
