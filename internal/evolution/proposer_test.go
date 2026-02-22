package evolution

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/livetemplate/lvt/internal/evolution/knowledge"
	"github.com/livetemplate/lvt/internal/telemetry"
)

func newTestKB(t *testing.T) *knowledge.KnowledgeBase {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "patterns.md")
	content := `# Test Patterns

## Pattern: test-compile-error

**Name:** Test Compile Error
**Confidence:** 0.90
**Added:** 2026-01-01
**Fix Count:** 0
**Success Rate:** -

### Description

A test compilation error pattern.

### Error Pattern

- **Phase:** compilation
- **Message Regex:** ` + "`cannot convert .* to type int`" + `
- **Context Regex:** ` + "`EditingID`" + `

### Fix

- **File:** ` + "`*/template.tmpl`" + `
- **Find:** ` + "`{{if ne .EditingID 0}}`" + `
- **Replace:** ` + "`{{if ne .EditingID \"\"}}`" + `
- **Is Regex:** false

---

## Pattern: test-runtime-error

**Name:** Test Runtime Error
**Confidence:** 0.85
**Added:** 2026-01-01
**Fix Count:** 0
**Success Rate:** -

### Description

A test runtime error pattern.

### Error Pattern

- **Phase:** runtime
- **Message Regex:** ` + "`modal open after reload`" + `

### Fix

- **File:** ` + "`*/handler.go`" + `
- **Find:** ` + "`IsAdding bool`" + `
- **Replace:** ` + "`IsAdding bool (transient)`" + `
- **Is Regex:** false
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	kb, err := knowledge.New(path)
	if err != nil {
		t.Fatalf("new kb: %v", err)
	}
	return kb
}

func TestProposeFor_KnownPattern(t *testing.T) {
	kb := newTestKB(t)
	proposer := NewProposer(kb)

	event := &telemetry.GenerationEvent{
		ID:      "evt-1",
		Command: "gen resource",
		Errors: []telemetry.GenerationError{
			{
				Phase:   "compilation",
				Message: `cannot convert "abc" to type int`,
				Context: "EditingID",
			},
		},
	}

	proposal, err := proposer.ProposeFor(event)
	if err != nil {
		t.Fatalf("propose: %v", err)
	}
	if len(proposal.Fixes) == 0 {
		t.Fatal("expected at least one fix")
	}
	if proposal.Fixes[0].PatternID != "test-compile-error" {
		t.Errorf("expected pattern 'test-compile-error', got %q", proposal.Fixes[0].PatternID)
	}
	if proposal.Fixes[0].Source != "knowledge_base" {
		t.Errorf("expected source 'knowledge_base', got %q", proposal.Fixes[0].Source)
	}
}

func TestProposeFor_UnknownError(t *testing.T) {
	kb := newTestKB(t)
	proposer := NewProposer(kb)

	event := &telemetry.GenerationEvent{
		ID:      "evt-2",
		Command: "gen resource",
		Errors: []telemetry.GenerationError{
			{
				Phase:   "compilation",
				Message: "some unknown error that doesn't match anything",
			},
		},
	}

	proposal, err := proposer.ProposeFor(event)
	if err != nil {
		t.Fatalf("propose: %v", err)
	}
	if len(proposal.Fixes) != 0 {
		t.Errorf("expected 0 fixes for unknown error, got %d", len(proposal.Fixes))
	}
}

func TestProposeFor_MultipleErrors(t *testing.T) {
	kb := newTestKB(t)
	proposer := NewProposer(kb)

	event := &telemetry.GenerationEvent{
		ID:      "evt-3",
		Command: "gen resource",
		Errors: []telemetry.GenerationError{
			{
				Phase:   "compilation",
				Message: `cannot convert "x" to type int`,
				Context: "EditingID",
			},
			{
				Phase:   "runtime",
				Message: "modal open after reload",
			},
		},
	}

	proposal, err := proposer.ProposeFor(event)
	if err != nil {
		t.Fatalf("propose: %v", err)
	}
	if len(proposal.Fixes) != 2 {
		t.Fatalf("expected 2 fixes, got %d", len(proposal.Fixes))
	}

	// Should be sorted by confidence (0.90, 0.85)
	if proposal.Fixes[0].Confidence < proposal.Fixes[1].Confidence {
		t.Error("expected fixes sorted by confidence descending")
	}
}

func TestProposeFor_Deduplication(t *testing.T) {
	kb := newTestKB(t)
	proposer := NewProposer(kb)

	// Two identical errors should produce deduplicated fixes
	event := &telemetry.GenerationEvent{
		ID:      "evt-4",
		Command: "gen resource",
		Errors: []telemetry.GenerationError{
			{
				Phase:   "compilation",
				Message: `cannot convert "a" to type int`,
				Context: "EditingID",
			},
			{
				Phase:   "compilation",
				Message: `cannot convert "b" to type int`,
				Context: "EditingID",
			},
		},
	}

	proposal, err := proposer.ProposeFor(event)
	if err != nil {
		t.Fatalf("propose: %v", err)
	}
	// Same pattern+file should be deduplicated
	if len(proposal.Fixes) != 1 {
		t.Errorf("expected 1 deduplicated fix, got %d", len(proposal.Fixes))
	}
}
