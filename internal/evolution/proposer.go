package evolution

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/livetemplate/lvt/internal/evolution/knowledge"
	"github.com/livetemplate/lvt/internal/telemetry"
)

// Proposer analyses generation events and proposes fixes from the knowledge base.
type Proposer struct {
	kb *knowledge.KnowledgeBase
}

// NewProposer creates a new Proposer backed by the given knowledge base.
func NewProposer(kb *knowledge.KnowledgeBase) *Proposer {
	return &Proposer{kb: kb}
}

// ProposeFor analyses a generation event and returns fix proposals.
// In v1 this only uses the knowledge base (no LLM).
func (p *Proposer) ProposeFor(event *telemetry.GenerationEvent) (*Proposal, error) {
	if event == nil {
		return nil, fmt.Errorf("event is nil")
	}

	proposal := &Proposal{EventID: event.ID}

	// Track seen pattern+file combos for deduplication
	seen := make(map[string]bool)

	for _, genErr := range event.Errors {
		patterns := knowledge.MatchAll(p.kb.ListPatterns(), genErr)
		for _, pat := range patterns {
			for _, fix := range pat.Fixes {
				key := pat.ID + ":" + fix.File + ":" + fix.FindPattern + ":" + fix.Replace
				if seen[key] {
					continue
				}
				seen[key] = true

				proposal.Fixes = append(proposal.Fixes, Fix{
					ID:          generateFixID(),
					PatternID:   pat.ID,
					TargetFile:  fix.File,
					FindPattern: fix.FindPattern,
					Replace:     fix.Replace,
					IsRegex:     fix.IsRegex,
					Confidence:  pat.Confidence,
					Rationale:   fmt.Sprintf("Pattern %q (%s): %s", pat.Name, pat.ID, pat.Description),
					Source:      "knowledge_base",
				})
			}
		}
	}

	// Sort by confidence (highest first)
	sort.Slice(proposal.Fixes, func(i, j int) bool {
		return proposal.Fixes[i].Confidence > proposal.Fixes[j].Confidence
	})

	return proposal, nil
}

func generateFixID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "fix-" + hex.EncodeToString(b)
}
