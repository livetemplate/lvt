package knowledge

import "github.com/livetemplate/lvt/internal/telemetry"

// Match checks if a GenerationError matches this pattern.
// Returns true when: phase matches (or pattern has empty phase) AND message regex
// matches AND context regex matches (or pattern has no context regex).
func (p *Pattern) Match(err telemetry.GenerationError) bool {
	// Phase must match if specified
	if p.ErrorPhase != "" && p.ErrorPhase != err.Phase {
		return false
	}

	// Message regex must match
	if p.MessageRe == nil {
		return false
	}
	if !p.MessageRe.MatchString(err.Message) {
		return false
	}

	// Context regex must match if specified
	if p.ContextRe != nil {
		// Check both the Context field and the Message for context matches
		if !p.ContextRe.MatchString(err.Context) && !p.ContextRe.MatchString(err.Message) {
			return false
		}
	}

	return true
}

// MatchAll returns all patterns that match the given error.
func MatchAll(patterns []*Pattern, err telemetry.GenerationError) []*Pattern {
	var matched []*Pattern
	for _, p := range patterns {
		if p.Match(err) {
			matched = append(matched, p)
		}
	}
	return matched
}
