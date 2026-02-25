package knowledge

import "regexp"

// Pattern represents a known error pattern with its metadata and fix templates.
type Pattern struct {
	ID           string
	Name         string
	Description  string
	Confidence   float64
	Added        string
	FixCount     int
	SuccessRate  string
	UpstreamRepo string // empty for local patterns

	// Matching criteria
	ErrorPhase string         // "compilation", "runtime", "template", "generation"
	MessageRe  *regexp.Regexp // compiled from patterns.md
	ContextRe  *regexp.Regexp // optional, compiled from patterns.md

	// Fixes (can be multiple per pattern)
	Fixes []FixTemplate
}

// FixTemplate describes a find-and-replace operation on a file.
type FixTemplate struct {
	// File is a glob pattern for the target file, e.g. "*/handler.go.tmpl".
	// Supported formats: exact filename, or "*/filename" to match in any subdirectory.
	// Multi-level wildcards (e.g. "**/*.go") are not supported.
	File string
	FindPattern string
	Replace     string
	IsRegex     bool
}
