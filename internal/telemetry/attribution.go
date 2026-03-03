package telemetry

import (
	"reflect"
	"strings"
)

// ComponentError links an error to a specific component.
type ComponentError struct {
	Component string `json:"component"`
	Phase     string `json:"phase"`
	Message   string `json:"message"`
	File      string `json:"file,omitempty"`
}

// AttributeErrors scans GenerationErrors and matches them to components
// by checking file paths (e.g. "components/modal/modal.go") and error
// messages (e.g. "modal.New: invalid size"). Only components listed in
// componentsUsed are considered.
func AttributeErrors(errors []GenerationError, componentsUsed []string) []ComponentError {
	if len(errors) == 0 || len(componentsUsed) == 0 {
		return nil
	}

	var result []ComponentError
	for _, genErr := range errors {
		comp := matchComponent(genErr, componentsUsed)
		if comp == "" {
			continue
		}
		result = append(result, ComponentError{
			Component: comp,
			Phase:     genErr.Phase,
			Message:   genErr.Message,
			File:      genErr.File,
		})
	}
	return result
}

// matchComponent checks whether a GenerationError relates to a specific component
// by scanning the file path and error message.
func matchComponent(genErr GenerationError, componentsUsed []string) string {
	fileLower := strings.ToLower(genErr.File)
	msgLower := strings.ToLower(genErr.Message)

	for _, comp := range componentsUsed {
		compLower := strings.ToLower(comp)

		// Check file path patterns: "components/modal/", "/modal/", "modal.go", "modal.tmpl"
		if strings.Contains(fileLower, "components/"+compLower+"/") ||
			strings.Contains(fileLower, "/"+compLower+"/"+compLower+".") ||
			strings.HasSuffix(fileLower, "/"+compLower+".go") ||
			strings.HasSuffix(fileLower, "/"+compLower+".tmpl") {
			return comp
		}

		// Check error message patterns: "modal.New:", "modal:", "modal "
		if strings.Contains(msgLower, compLower+".") ||
			strings.Contains(msgLower, compLower+":") ||
			strings.HasPrefix(msgLower, compLower+" ") {
			return comp
		}
	}
	return ""
}

// ComponentsFromUsage converts a ComponentUsage struct (or any struct with
// bool fields prefixed "Use") into a string slice of component names.
// For example, UseModal=true becomes "modal".
func ComponentsFromUsage(usage any) []string {
	if usage == nil {
		return nil
	}

	v := reflect.ValueOf(usage)
	// Dereference pointer if needed
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	var components []string
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() != reflect.Bool {
			continue
		}
		if !v.Field(i).Bool() {
			continue
		}
		// Convert "UseModal" → "modal", "UseToast" → "toast".
		// Only process fields with "Use" prefix — other bool fields are ignored.
		name := field.Name
		if !strings.HasPrefix(name, "Use") {
			continue
		}
		name = name[3:]
		components = append(components, strings.ToLower(name))
	}
	return components
}
