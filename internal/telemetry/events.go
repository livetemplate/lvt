package telemetry

import "time"

// GenerationEvent records a single code generation invocation.
type GenerationEvent struct {
	ID             string            `json:"id"`
	Timestamp      time.Time         `json:"timestamp"`
	Command        string            `json:"command"`                  // "gen resource", "gen view", "gen auth", "gen schema"
	Inputs         map[string]any    `json:"inputs"`                   // command args as key-value
	Kit            string            `json:"kit,omitempty"`            // e.g. "multi", "single"
	LvtVersion     string            `json:"lvt_version"`              // version of lvt used
	Success        bool              `json:"success"`                  // whether generation succeeded
	ValidationJSON string            `json:"validation,omitempty"`     // JSON of validator.ValidationResult
	Errors         []GenerationError `json:"errors,omitempty"`         // errors captured during generation
	DurationMs     int64             `json:"duration_ms"`              // wall-clock duration
	FilesGenerated []string          `json:"files_generated,omitempty"` // paths of generated files
}

// GenerationError records a single error captured during generation.
type GenerationError struct {
	Phase   string `json:"phase"`             // "generation", "compilation", "runtime", "template"
	File    string `json:"file,omitempty"`    // file where error occurred
	Line    int    `json:"line,omitempty"`    // line number
	Message string `json:"message"`           // error message
	Context string `json:"context,omitempty"` // surrounding code or context
}
