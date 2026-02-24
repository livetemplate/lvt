package evolution

import (
	"regexp"
	"testing"

	"github.com/livetemplate/lvt/internal/evolution/knowledge"
	"github.com/livetemplate/lvt/internal/telemetry"
)

func TestUpstreamProposer_ProposeUpstreamFix(t *testing.T) {
	pattern := &knowledge.Pattern{
		ID:           "test-upstream",
		Name:         "Test Upstream Pattern",
		Description:  "A test upstream pattern",
		Confidence:   0.85,
		ErrorPhase:   "runtime",
		MessageRe:    regexp.MustCompile(`some error`),
		UpstreamRepo: "github.com/example/lib",
		Fixes: []knowledge.FixTemplate{
			{File: "*/handler.go", FindPattern: "old", Replace: "new"},
		},
	}

	kb := knowledge.NewFromPatterns([]*knowledge.Pattern{pattern})
	up := NewUpstreamProposer(kb)

	err := telemetry.GenerationError{
		Phase:   "runtime",
		Message: "some error",
	}

	fix := up.ProposeUpstreamFix(pattern, err)
	if fix == nil {
		t.Fatal("expected upstream fix, got nil")
	}
	if fix.UpstreamRepo != "github.com/example/lib" {
		t.Errorf("expected upstream repo, got %q", fix.UpstreamRepo)
	}
	if fix.TargetBranch != "main" {
		t.Errorf("expected target branch 'main', got %q", fix.TargetBranch)
	}
	if fix.PatternID != "test-upstream" {
		t.Errorf("expected pattern ID 'test-upstream', got %q", fix.PatternID)
	}
}

func TestUpstreamProposer_NoUpstreamRepo(t *testing.T) {
	pattern := &knowledge.Pattern{
		ID:         "local-only",
		Name:       "Local Pattern",
		ErrorPhase: "compilation",
		Fixes: []knowledge.FixTemplate{
			{File: "*/handler.go", FindPattern: "old", Replace: "new"},
		},
	}

	kb := knowledge.NewFromPatterns([]*knowledge.Pattern{pattern})
	up := NewUpstreamProposer(kb)

	err := telemetry.GenerationError{Phase: "compilation", Message: "error"}
	fix := up.ProposeUpstreamFix(pattern, err)
	if fix != nil {
		t.Error("expected nil fix for pattern without upstream repo")
	}
}

func TestUpstreamProposer_ListUpstreamPatterns(t *testing.T) {
	patterns := []*knowledge.Pattern{
		{ID: "local", Name: "Local"},
		{ID: "upstream-1", Name: "Upstream 1", UpstreamRepo: "github.com/a/b"},
		{ID: "upstream-2", Name: "Upstream 2", UpstreamRepo: "github.com/c/d"},
	}

	kb := knowledge.NewFromPatterns(patterns)
	up := NewUpstreamProposer(kb)

	upstream := up.ListUpstreamPatterns()
	if len(upstream) != 2 {
		t.Errorf("expected 2 upstream patterns, got %d", len(upstream))
	}
}
