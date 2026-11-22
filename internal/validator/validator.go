package validator

import (
	"fmt"
	"strings"
)

// ValidationLevel represents the severity of a validation issue
type ValidationLevel string

const (
	LevelError   ValidationLevel = "error"
	LevelWarning ValidationLevel = "warning"
	LevelInfo    ValidationLevel = "info"
)

// ValidationIssue represents a single validation problem
type ValidationIssue struct {
	Level   ValidationLevel
	Message string
	File    string
	Line    int
	Hint    string
}

// ValidationResult holds the results of a validation check
type ValidationResult struct {
	Valid  bool
	Issues []ValidationIssue
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:  true,
		Issues: []ValidationIssue{},
	}
}

// AddError adds an error-level issue
func (r *ValidationResult) AddError(message string, file string, line int) {
	r.Valid = false
	r.Issues = append(r.Issues, ValidationIssue{
		Level:   LevelError,
		Message: message,
		File:    file,
		Line:    line,
	})
}

// AddErrorWithHint adds an error-level issue with a hint for context.
func (r *ValidationResult) AddErrorWithHint(message, file string, line int, hint string) {
	r.Valid = false
	r.Issues = append(r.Issues, ValidationIssue{
		Level:   LevelError,
		Message: message,
		File:    file,
		Line:    line,
		Hint:    hint,
	})
}

// AddWarning adds a warning-level issue
func (r *ValidationResult) AddWarning(message string, file string, line int) {
	r.Issues = append(r.Issues, ValidationIssue{
		Level:   LevelWarning,
		Message: message,
		File:    file,
		Line:    line,
	})
}

// AddInfo adds an info-level issue
func (r *ValidationResult) AddInfo(message string, file string, line int) {
	r.Issues = append(r.Issues, ValidationIssue{
		Level:   LevelInfo,
		Message: message,
		File:    file,
		Line:    line,
	})
}

// Merge combines another validation result into this one
func (r *ValidationResult) Merge(other *ValidationResult) {
	if !other.Valid {
		r.Valid = false
	}
	r.Issues = append(r.Issues, other.Issues...)
}

// HasErrors returns true if there are any error-level issues
func (r *ValidationResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Level == LevelError {
			return true
		}
	}
	return false
}

// ErrorCount returns the number of errors
func (r *ValidationResult) ErrorCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Level == LevelError {
			count++
		}
	}
	return count
}

// WarningCount returns the number of warnings
func (r *ValidationResult) WarningCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Level == LevelWarning {
			count++
		}
	}
	return count
}

// Format returns a formatted string representation
func (r *ValidationResult) Format() string {
	var sb strings.Builder

	if r.Valid {
		sb.WriteString("✅ Validation passed\n")
	} else {
		sb.WriteString("❌ Validation failed\n")
	}

	if len(r.Issues) == 0 {
		sb.WriteString("\nNo issues found\n")
		return sb.String()
	}

	fmt.Fprintf(&sb, "\nFound %d issue(s):\n", len(r.Issues))

	// Group by level
	errors := []ValidationIssue{}
	warnings := []ValidationIssue{}
	infos := []ValidationIssue{}

	for _, issue := range r.Issues {
		switch issue.Level {
		case LevelError:
			errors = append(errors, issue)
		case LevelWarning:
			warnings = append(warnings, issue)
		case LevelInfo:
			infos = append(infos, issue)
		}
	}

	// Print errors
	if len(errors) > 0 {
		fmt.Fprintf(&sb, "\n❌ Errors (%d):\n", len(errors))
		for _, issue := range errors {
			sb.WriteString(formatIssue(issue))
		}
	}

	// Print warnings
	if len(warnings) > 0 {
		fmt.Fprintf(&sb, "\n⚠️  Warnings (%d):\n", len(warnings))
		for _, issue := range warnings {
			sb.WriteString(formatIssue(issue))
		}
	}

	// Print infos
	if len(infos) > 0 {
		fmt.Fprintf(&sb, "\nℹ️  Info (%d):\n", len(infos))
		for _, issue := range infos {
			sb.WriteString(formatIssue(issue))
		}
	}

	return sb.String()
}

func formatIssue(issue ValidationIssue) string {
	var sb strings.Builder

	if issue.File != "" {
		if issue.Line > 0 {
			fmt.Fprintf(&sb, "  %s:%d: %s\n", issue.File, issue.Line, issue.Message)
		} else {
			fmt.Fprintf(&sb, "  %s: %s\n", issue.File, issue.Message)
		}
	} else {
		fmt.Fprintf(&sb, "  %s\n", issue.Message)
	}

	if issue.Hint != "" {
		fmt.Fprintf(&sb, "    💡 %s\n", issue.Hint)
	}

	return sb.String()
}
