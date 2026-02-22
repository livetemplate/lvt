package tmplutil

import (
	"fmt"
	"strings"
	"testing"
)

func TestExtractLineNumber(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected int
	}{
		{"standard error", "template: test.tmpl:5: unexpected", 5},
		{"with column", "template: test.tmpl:12:22: function not defined", 12},
		{"no line number", "some other error", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractLineNumber(fmt.Errorf("%s", tt.errMsg))
			if got != tt.expected {
				t.Errorf("ExtractLineNumber(%q) = %d, want %d", tt.errMsg, got, tt.expected)
			}
		})
	}
}

func TestSourceContext_Arrow(t *testing.T) {
	content := "line1\nline2\nline3\nline4\nline5\nline6\nline7"

	ctx := SourceContext(content, 4, 2, "→ ")

	// Should include lines 2-6 (4 ± 2)
	if !strings.Contains(ctx, "line2") {
		t.Error("context should include line 2")
	}
	if !strings.Contains(ctx, "line6") {
		t.Error("context should include line 6")
	}
	// Line 4 should be marked with arrow
	if !strings.Contains(ctx, "→") {
		t.Error("context should mark error line with →")
	}
	// Should show line numbers
	if !strings.Contains(ctx, "4 |") {
		t.Error("context should show line numbers")
	}
}

func TestSourceContext_GreaterThan(t *testing.T) {
	content := "line1\nline2\nline3\nline4\nline5"

	ctx := SourceContext(content, 3, 1, "> ")

	if !strings.Contains(ctx, "> ") {
		t.Error("context should use > marker")
	}
	if strings.Contains(ctx, "→") {
		t.Error("context should not use → marker when > was requested")
	}
}

func TestSourceContext_FirstLine(t *testing.T) {
	content := "line1\nline2\nline3"

	ctx := SourceContext(content, 1, 2, "→ ")

	// Should include lines 1-3
	if !strings.Contains(ctx, "line1") {
		t.Error("context should include line 1")
	}
	if !strings.Contains(ctx, "→") {
		t.Error("context should mark error line")
	}
}

func TestSourceContext_OutOfRange(t *testing.T) {
	content := "line1\nline2"

	ctx := SourceContext(content, 10, 2, "→ ")
	if ctx != "" {
		t.Errorf("out of range should return empty, got: %s", ctx)
	}

	ctx = SourceContext(content, 0, 2, "→ ")
	if ctx != "" {
		t.Errorf("zero line should return empty, got: %s", ctx)
	}
}
