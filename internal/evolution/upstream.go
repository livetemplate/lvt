package evolution

import (
	"github.com/livetemplate/lvt/internal/evolution/knowledge"
	"github.com/livetemplate/lvt/internal/telemetry"
)

// UpstreamFix extends Fix with upstream repository targeting information.
type UpstreamFix struct {
	Fix
	UpstreamRepo string
	TargetBranch string // default "main"
}

// UpstreamProposer analyses errors against upstream patterns.
type UpstreamProposer struct {
	kb *knowledge.KnowledgeBase
}

// NewUpstreamProposer creates a new UpstreamProposer.
func NewUpstreamProposer(kb *knowledge.KnowledgeBase) *UpstreamProposer {
	return &UpstreamProposer{kb: kb}
}

// ProposeUpstreamFix returns an upstream fix if the pattern has UpstreamRepo set.
func (u *UpstreamProposer) ProposeUpstreamFix(pattern *knowledge.Pattern, err telemetry.GenerationError) *UpstreamFix {
	if pattern.UpstreamRepo == "" {
		return nil
	}
	if !pattern.Match(err) {
		return nil
	}
	if len(pattern.Fixes) == 0 {
		return nil
	}

	// TODO: v1 only uses the first fix. Multi-fix upstream proposals are not yet supported.
	fix := pattern.Fixes[0]
	return &UpstreamFix{
		Fix: Fix{
			ID:          generateFixID(),
			PatternID:   pattern.ID,
			TargetFile:  fix.File,
			FindPattern: fix.FindPattern,
			Replace:     fix.Replace,
			IsRegex:     fix.IsRegex,
			Confidence:  pattern.Confidence,
			Rationale:   pattern.Description,
			Source:      "knowledge_base",
		},
		UpstreamRepo: pattern.UpstreamRepo,
		TargetBranch: "main",
	}
}

// ListUpstreamPatterns returns only patterns with UpstreamRepo set.
func (u *UpstreamProposer) ListUpstreamPatterns() []*knowledge.Pattern {
	var upstream []*knowledge.Pattern
	for _, p := range u.kb.ListPatterns() {
		if p.UpstreamRepo != "" {
			upstream = append(upstream, p)
		}
	}
	return upstream
}
