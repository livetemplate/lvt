package evolution

import (
	"strings"

	"github.com/livetemplate/lvt/internal/telemetry"
)

// ErrorLocation classifies where an error originated.
type ErrorLocation struct {
	Type      string // "component", "kit", "generated", "unknown"
	Component string // e.g. "modal" (only when Type == "component")
	Path      string // original path as provided
}

// ClassifyError determines if an error is in components/, internal/kits/,
// or generated app code.
func ClassifyError(err telemetry.GenerationError) ErrorLocation {
	return classifyPath(err.File)
}

// ClassifyFix determines the fix location category from a Fix's target file.
func ClassifyFix(fix Fix) ErrorLocation {
	return classifyPath(fix.TargetFile)
}

func classifyPath(path string) ErrorLocation {
	loc := ErrorLocation{Path: path}
	lower := strings.ToLower(path)

	// Check kit paths first — kits may contain a components/ subdirectory
	// (e.g. internal/kits/system/multi/components/form.tmpl) which would
	// otherwise be misclassified as a top-level component.
	if strings.Contains(lower, "internal/kits/") {
		loc.Type = "kit"
		return loc
	}

	// Check for top-level component paths: components/<name>/
	if idx := strings.Index(lower, "components/"); idx != -1 {
		rest := lower[idx+len("components/"):]
		if slashIdx := strings.Index(rest, "/"); slashIdx > 0 {
			loc.Type = "component"
			loc.Component = rest[:slashIdx]
			return loc
		}
		// Bare component reference like "components/modal.go"
		name := strings.TrimSuffix(rest, ".go")
		name = strings.TrimSuffix(name, ".tmpl")
		if name != "" {
			loc.Type = "component"
			loc.Component = name
			return loc
		}
	}

	// Check for generated app code: must be at path start or after a separator
	// to avoid false positives from "webapp/", "myapp/", etc.
	if strings.HasPrefix(lower, "app/") || strings.Contains(lower, "/app/") {
		loc.Type = "generated"
		return loc
	}

	loc.Type = "unknown"
	return loc
}
