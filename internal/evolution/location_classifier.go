package evolution

import (
	"strings"

	"github.com/livetemplate/lvt/internal/telemetry"
)

// LocationType classifies where a file originated in the project structure.
type LocationType string

// Location type constants for ErrorLocation.Type.
const (
	LocationComponent LocationType = "component"
	LocationKit       LocationType = "kit"
	LocationGenerated LocationType = "generated"
	LocationUnknown   LocationType = "unknown"
)

// ErrorLocation classifies where an error originated.
type ErrorLocation struct {
	Type      LocationType // one of Location* constants
	Component string       // normalized to lowercase; e.g. "modal" (only when Type == LocationComponent)
	Path      string       // original path as provided
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
	// Normalize path separators for cross-platform consistency
	normalized := strings.ReplaceAll(path, "\\", "/")
	loc := ErrorLocation{Path: normalized}
	lower := strings.ToLower(normalized)

	// Check kit paths first — kits may contain a components/ subdirectory
	// (e.g. internal/kits/system/multi/components/form.tmpl) which would
	// otherwise be misclassified as a top-level component.
	if strings.Contains(lower, "internal/kits/") {
		loc.Type = LocationKit
		return loc
	}

	// Check for top-level component paths: components/<name>/
	// Use segment-aware matching to avoid false positives from paths like
	// "custom_components/" — require "components/" at start or after "/".
	compIdx := -1
	if strings.HasPrefix(lower, "components/") {
		compIdx = 0
	} else if idx := strings.Index(lower, "/components/"); idx != -1 {
		compIdx = idx + 1 // skip the leading "/"
	}
	if compIdx != -1 {
		rest := lower[compIdx+len("components/"):]
		if slashIdx := strings.Index(rest, "/"); slashIdx > 0 {
			name := rest[:slashIdx]
			if !strings.Contains(name, "*") {
				loc.Type = LocationComponent
				loc.Component = name
				return loc
			}
		}
		// Bare component reference like "components/modal.go" or "components/modal.go.tmpl"
		// Strip .tmpl first so .go.tmpl → .go → (name)
		name := strings.TrimSuffix(rest, ".tmpl")
		name = strings.TrimSuffix(name, ".go")
		if name != "" && !strings.Contains(name, "*") {
			loc.Type = LocationComponent
			loc.Component = name
			return loc
		}
	}

	// Check for generated app code: must be at path start or after a separator
	// to avoid false positives from "webapp/", "myapp/", etc.
	if strings.HasPrefix(lower, "app/") || strings.Contains(lower, "/app/") {
		loc.Type = LocationGenerated
		return loc
	}

	loc.Type = LocationUnknown
	return loc
}
