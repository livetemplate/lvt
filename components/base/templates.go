package base

import (
	"embed"
	"html/template"

	"github.com/livetemplate/livetemplate"
)

// TemplateSet is a type alias for livetemplate.TemplateSet, allowing component
// Templates() functions to return values directly compatible with
// livetemplate.WithComponentTemplates() without conversion.
type TemplateSet = livetemplate.TemplateSet

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
//	    return WithFuncs(NewTemplateSet(templateFS, "templates/*.tmpl", "dropdown"),
//	        template.FuncMap{
//	            "dropdownClass": func() string { return "dropdown" },
//	        })
//	}
func WithFuncs(ts *TemplateSet, funcs template.FuncMap) *TemplateSet {
	merged := make(template.FuncMap, len(ts.Funcs)+len(funcs))
	for k, v := range ts.Funcs {
		merged[k] = v
	}
	for k, v := range funcs {
		merged[k] = v
	}
	return &TemplateSet{
		FS:        ts.FS,
		Pattern:   ts.Pattern,
		Namespace: ts.Namespace,
		Funcs:     merged,
	}
}

// TemplateProvider is implemented by components that provide templates.
// The LiveTemplate framework uses this interface to collect all component templates.
type TemplateProvider interface {
	// Templates returns the component's template set for registration.
	Templates() *TemplateSet
}
