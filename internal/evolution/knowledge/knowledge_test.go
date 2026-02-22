package knowledge

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/livetemplate/lvt/internal/telemetry"
)

func findTestPatternsFile(t *testing.T) string {
	t.Helper()
	// From internal/evolution/knowledge/ go up 3 levels to repo root
	p, err := FindPatternsFile()
	if err != nil {
		t.Skipf("patterns.md not found: %v", err)
	}
	return p
}

func TestParseActualPatternsFile(t *testing.T) {
	path := findTestPatternsFile(t)
	patterns, err := Parse(path)
	if err != nil {
		t.Fatalf("parse patterns: %v", err)
	}

	// The real file has 13 patterns (10 local + 3 upstream)
	if len(patterns) < 13 {
		t.Errorf("expected at least 13 patterns, got %d", len(patterns))
	}

	// Check first pattern
	if patterns[0].ID != "editing-id-type" {
		t.Errorf("expected first pattern ID 'editing-id-type', got %q", patterns[0].ID)
	}
}

func TestParsePattern_Metadata(t *testing.T) {
	path := findTestPatternsFile(t)
	patterns, err := Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// Find editing-id-type
	var p *Pattern
	for _, pat := range patterns {
		if pat.ID == "editing-id-type" {
			p = pat
			break
		}
	}
	if p == nil {
		t.Fatal("pattern editing-id-type not found")
	}

	if p.Name != "EditingID Type Mismatch" {
		t.Errorf("name: got %q", p.Name)
	}
	if p.Confidence != 0.95 {
		t.Errorf("confidence: got %f", p.Confidence)
	}
	if p.Added != "2026-01-19" {
		t.Errorf("added: got %q", p.Added)
	}
	if p.FixCount != 0 {
		t.Errorf("fix count: got %d", p.FixCount)
	}
	if p.SuccessRate != "-" {
		t.Errorf("success rate: got %q", p.SuccessRate)
	}
}

func TestParsePattern_MultipleFixes(t *testing.T) {
	path := findTestPatternsFile(t)
	patterns, err := Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	var p *Pattern
	for _, pat := range patterns {
		if pat.ID == "modal-state-persistence" {
			p = pat
			break
		}
	}
	if p == nil {
		t.Fatal("pattern modal-state-persistence not found")
	}

	if len(p.Fixes) != 2 {
		t.Fatalf("expected 2 fixes, got %d", len(p.Fixes))
	}

	if p.Fixes[0].FindPattern != "IsAdding bool" {
		t.Errorf("fix 1 find: got %q", p.Fixes[0].FindPattern)
	}
	if p.Fixes[1].FindPattern != "EditingID string" {
		t.Errorf("fix 2 find: got %q", p.Fixes[1].FindPattern)
	}
}

func TestParsePattern_UpstreamRepo(t *testing.T) {
	path := findTestPatternsFile(t)
	patterns, err := Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	upstreamCount := 0
	for _, p := range patterns {
		if p.UpstreamRepo != "" {
			upstreamCount++
		}
	}

	if upstreamCount != 3 {
		t.Errorf("expected 3 upstream patterns, got %d", upstreamCount)
	}

	// Check specific upstream pattern
	var p *Pattern
	for _, pat := range patterns {
		if pat.ID == "morphdom-select-sync" {
			p = pat
			break
		}
	}
	if p == nil {
		t.Fatal("pattern morphdom-select-sync not found")
	}
	if p.UpstreamRepo != "github.com/livetemplate/client" {
		t.Errorf("upstream repo: got %q", p.UpstreamRepo)
	}
}

func TestParsePattern_ErrorRegex(t *testing.T) {
	path := findTestPatternsFile(t)
	patterns, err := Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	var p *Pattern
	for _, pat := range patterns {
		if pat.ID == "editing-id-type" {
			p = pat
			break
		}
	}
	if p == nil {
		t.Fatal("pattern editing-id-type not found")
	}

	if p.MessageRe == nil {
		t.Fatal("message regex not compiled")
	}
	if !p.MessageRe.MatchString("cannot convert X to type int") {
		t.Error("expected message regex to match 'cannot convert X to type int'")
	}
	if p.ContextRe == nil {
		t.Fatal("context regex not compiled")
	}
	if !p.ContextRe.MatchString("EditingID") {
		t.Error("expected context regex to match 'EditingID'")
	}
}

func TestMatch_CompilationError(t *testing.T) {
	path := findTestPatternsFile(t)
	kb, err := New(path)
	if err != nil {
		t.Fatalf("new kb: %v", err)
	}

	err2 := telemetry.GenerationError{
		Phase:   "compilation",
		Message: `cannot convert "abc" to type int`,
		Context: "EditingID",
	}

	matched := MatchAll(kb.ListPatterns(), err2)
	if len(matched) == 0 {
		t.Fatal("expected at least one match")
	}
	if matched[0].ID != "editing-id-type" {
		t.Errorf("expected editing-id-type, got %q", matched[0].ID)
	}
}

func TestMatch_RuntimeError(t *testing.T) {
	path := findTestPatternsFile(t)
	kb, err := New(path)
	if err != nil {
		t.Fatalf("new kb: %v", err)
	}

	err2 := telemetry.GenerationError{
		Phase:   "runtime",
		Message: "modal open after reload",
		Context: "IsEditing",
	}

	matched := MatchAll(kb.ListPatterns(), err2)
	if len(matched) == 0 {
		t.Fatal("expected at least one match for runtime error")
	}

	found := false
	for _, m := range matched {
		if m.ID == "modal-state-persistence" {
			found = true
		}
	}
	if !found {
		t.Error("expected modal-state-persistence to match")
	}
}

func TestMatch_PhaseMismatch(t *testing.T) {
	path := findTestPatternsFile(t)
	kb, err := New(path)
	if err != nil {
		t.Fatalf("new kb: %v", err)
	}

	// This error matches editing-id-type message but with wrong phase
	err2 := telemetry.GenerationError{
		Phase:   "runtime", // should be "compilation"
		Message: `cannot convert "abc" to type int`,
		Context: "EditingID",
	}

	matched := MatchAll(kb.ListPatterns(), err2)
	for _, m := range matched {
		if m.ID == "editing-id-type" {
			t.Error("editing-id-type should NOT match runtime phase")
		}
	}
}

func TestLookupFixes(t *testing.T) {
	path := findTestPatternsFile(t)
	kb, err := New(path)
	if err != nil {
		t.Fatalf("new kb: %v", err)
	}

	err2 := telemetry.GenerationError{
		Phase:   "compilation",
		Message: `cannot convert "abc" to type int`,
		Context: "EditingID",
	}

	fixes := kb.LookupFixes(err2)
	if len(fixes) == 0 {
		t.Fatal("expected at least one fix")
	}
	if fixes[0].File != "*/template.tmpl.tmpl" {
		t.Errorf("expected fix file '*/template.tmpl.tmpl', got %q", fixes[0].File)
	}
}

func TestKnowledgeBase_Synthetic(t *testing.T) {
	// Test with a synthetic patterns file to avoid depending on real file
	dir := t.TempDir()
	path := filepath.Join(dir, "patterns.md")
	content := `# Test Patterns

## Pattern: test-pattern

**Name:** Test Pattern
**Confidence:** 0.75
**Added:** 2026-01-01
**Fix Count:** 5
**Success Rate:** 80%

### Description

A test pattern for unit tests.

### Error Pattern

- **Phase:** compilation
- **Message Regex:** ` + "`test error \\d+`" + `
- **Context Regex:** ` + "`TestCtx`" + `

### Fix

- **File:** ` + "`*/test.go`" + `
- **Find:** ` + "`old code`" + `
- **Replace:** ` + "`new code`" + `
- **Is Regex:** false
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	kb, err := New(path)
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if kb.PatternCount() != 1 {
		t.Fatalf("expected 1 pattern, got %d", kb.PatternCount())
	}

	p, ok := kb.GetPattern("test-pattern")
	if !ok {
		t.Fatal("pattern not found")
	}
	if p.Name != "Test Pattern" {
		t.Errorf("name: %q", p.Name)
	}
	if p.Confidence != 0.75 {
		t.Errorf("confidence: %f", p.Confidence)
	}
	if p.FixCount != 5 {
		t.Errorf("fix count: %d", p.FixCount)
	}
	if p.SuccessRate != "80%" {
		t.Errorf("success rate: %q", p.SuccessRate)
	}
	if p.ErrorPhase != "compilation" {
		t.Errorf("phase: %q", p.ErrorPhase)
	}
	if !p.MessageRe.MatchString("test error 42") {
		t.Error("message regex should match 'test error 42'")
	}
	if len(p.Fixes) != 1 {
		t.Fatalf("expected 1 fix, got %d", len(p.Fixes))
	}
	if p.Fixes[0].FindPattern != "old code" {
		t.Errorf("fix find: %q", p.Fixes[0].FindPattern)
	}
}
