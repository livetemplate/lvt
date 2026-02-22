package evolution

import "github.com/livetemplate/lvt/internal/validator"

// Fix represents a concrete fix to apply to a template file.
type Fix struct {
	ID          string
	PatternID   string
	TargetFile  string // glob pattern
	FindPattern string
	Replace     string
	IsRegex     bool
	Confidence  float64
	Rationale   string
	Source      string // "knowledge_base"
}

// Proposal groups fixes proposed for a specific generation event.
type Proposal struct {
	EventID string
	Fixes   []Fix
}

// TestResult captures the outcome of testing a fix in isolation.
type TestResult struct {
	FixID      string
	Success    bool
	Error      string
	Validation *validator.ValidationResult
}
