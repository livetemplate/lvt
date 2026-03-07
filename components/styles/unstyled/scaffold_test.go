package unstyled

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateCSS(t *testing.T) {
	var buf bytes.Buffer
	err := GenerateCSS(&buf)
	if err != nil {
		t.Fatalf("GenerateCSS returned error: %v", err)
	}

	css := buf.String()

	// Should contain the header
	if !strings.Contains(css, "LiveTemplate Unstyled CSS Scaffold") {
		t.Error("missing scaffold header")
	}

	// Should contain utility classes
	if !strings.Contains(css, ".lvt-sr-only") {
		t.Error("missing .lvt-sr-only utility class")
	}

	// Should contain section headers for all component groups
	expectedSections := []string{
		"Accordion", "Autocomplete", "Breadcrumbs", "Data Table",
		"Date Picker", "Drawer", "Dropdown", "Menu", "Modal",
		"Confirm Modal", "Sheet", "Popover", "Progress",
		"Circular Progress", "Spinner", "Rating", "Skeleton",
		"Avatar Skeleton", "Card Skeleton", "Tabs", "Tags Input",
		"Timeline", "Timeline Item", "Time Picker", "Toast",
		"Toggle", "Checkbox", "Tooltip",
	}
	for _, sec := range expectedSections {
		if !strings.Contains(css, "/* === "+sec+" === */") {
			t.Errorf("missing section header for %q", sec)
		}
	}

	// Should contain essential layout rules for overlay components
	if !strings.Contains(css, "position: fixed") {
		t.Error("missing essential layout rules (position: fixed)")
	}

	// Should contain modal overlay rule
	if !strings.Contains(css, ".lvt-modal__overlay") {
		t.Error("missing .lvt-modal__overlay class")
	}

	// Should contain toast container
	if !strings.Contains(css, ".lvt-toast-container") {
		t.Error("missing .lvt-toast-container class")
	}
}

func TestGenerateCSSContainsAllClasses(t *testing.T) {
	var buf bytes.Buffer
	if err := GenerateCSS(&buf); err != nil {
		t.Fatalf("GenerateCSS returned error: %v", err)
	}

	css := buf.String()
	allNames := AllClassNames()

	// Filter out lvt-sr-only since it's in the utility section, not a component
	for _, name := range allNames {
		if name == "lvt-sr-only" {
			continue
		}
		selector := "." + name
		if !strings.Contains(css, selector) {
			t.Errorf("CSS output missing class %q", name)
		}
	}
}

func TestEssentialRulesApplied(t *testing.T) {
	var buf bytes.Buffer
	if err := GenerateCSS(&buf); err != nil {
		t.Fatalf("GenerateCSS returned error: %v", err)
	}

	css := buf.String()

	// Check that overlay components have essential rules
	essentialClasses := []struct {
		class    string
		contains string
	}{
		{"lvt-modal", "z-index: 50"},
		{"lvt-modal__overlay", "rgba(0,0,0,0.5)"},
		{"lvt-drawer", "z-index: 40"},
		{"lvt-sheet__panel", "position: fixed"},
		{"lvt-toast-container", "pointer-events: none"},
	}

	for _, ec := range essentialClasses {
		// Find the CSS block for this class
		idx := strings.Index(css, "."+ec.class+" {")
		if idx < 0 {
			t.Errorf("missing CSS block for .%s", ec.class)
			continue
		}
		// Get the block content (find closing brace)
		blockEnd := strings.Index(css[idx:], "}")
		if blockEnd < 0 {
			t.Errorf("unclosed CSS block for .%s", ec.class)
			continue
		}
		block := css[idx : idx+blockEnd]
		if !strings.Contains(block, ec.contains) {
			t.Errorf(".%s block should contain %q, got: %s", ec.class, ec.contains, block)
		}
	}
}
