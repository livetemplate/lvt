package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// VisualIssue represents a visual problem found by the LLM reviewer.
type VisualIssue struct {
	Severity    string `json:"severity"`    // HIGH, MEDIUM, LOW
	Category    string `json:"category"`    // ALIGNMENT, HIERARCHY, SPACING, ERROR_STATE, LAYOUT, READABILITY
	Description string `json:"description"` // Human-readable description
}

const visualCheckPrompt = `You are a UI/UX reviewer for a web application using Pico CSS (https://picocss.com).
Analyze the screenshot at %s for visual issues.

Page description: %s

Check for these specific categories:
1. ALIGNMENT: Are input fields and buttons the same height when grouped together (inside fieldset)? Are elements properly aligned horizontally/vertically?
2. HIERARCHY: Are flash/notification messages visually distinct from surrounding content (colored background, border, or different styling)?
3. SPACING: Is spacing consistent between elements? No unexpected gaps or cramped areas?
4. ERROR_STATE: Do error messages use red/colored styling? Are invalid inputs visually highlighted?
5. LAYOUT: Is anything overflowing its container, overlapping other elements, or cut off?
6. READABILITY: Is all text readable with sufficient contrast?

Respond ONLY with a JSON array of issues found. Each issue must have severity (HIGH, MEDIUM, or LOW), category, and description.
Example: [{"severity":"HIGH","category":"ALIGNMENT","description":"Input and button have different heights in the form group"}]
If no issues found, respond with: []`

// ValidateScreenshotWithLLM captures a full-page screenshot and sends it to
// the Claude CLI for visual analysis. Skipped unless LVT_VISUAL_CHECK=true is set.
// Requires the 'claude' CLI to be installed and authenticated.
//
// Fails the test if any HIGH or CRITICAL severity issues are found.
// MEDIUM and LOW issues are logged as warnings.
func ValidateScreenshotWithLLM(t *testing.T, ctx context.Context, pageDescription string) {
	t.Helper()

	if os.Getenv("LVT_VISUAL_CHECK") != "true" {
		t.Skip("LLM visual check disabled — set LVT_VISUAL_CHECK=true to enable")
	}

	// Check claude CLI is available
	if _, err := exec.LookPath("claude"); err != nil {
		t.Skip("claude CLI not found — install Claude Code to enable LLM visual check")
	}

	// Capture screenshot
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf)); err != nil {
		t.Fatalf("Failed to capture screenshot: %v", err)
	}

	// Save to temp file
	screenshotPath := filepath.Join(t.TempDir(), "screenshot.png")
	if err := os.WriteFile(screenshotPath, buf, 0644); err != nil {
		t.Fatalf("Failed to write screenshot: %v", err)
	}

	t.Logf("Captured screenshot: %d bytes → %s", len(buf), screenshotPath)

	// Build prompt with screenshot path
	prompt := fmt.Sprintf(visualCheckPrompt, screenshotPath, pageDescription)

	// Run claude CLI in non-interactive mode
	cmd := exec.Command("claude", "-p", prompt, "--output-format", "text")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("claude CLI failed: %v\nOutput: %s", err, output)
	}

	responseText := strings.TrimSpace(string(output))
	t.Logf("LLM response: %s", responseText)

	// Parse JSON — strip markdown code fence if present
	jsonStr := responseText
	if strings.HasPrefix(jsonStr, "```") {
		lines := strings.Split(jsonStr, "\n")
		if len(lines) > 2 {
			jsonStr = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}
	jsonStr = strings.TrimSpace(jsonStr)

	var issues []VisualIssue
	if err := json.Unmarshal([]byte(jsonStr), &issues); err != nil {
		t.Logf("Warning: could not parse LLM response as JSON: %v", err)
		t.Logf("Raw response: %s", responseText)
		return
	}

	if len(issues) == 0 {
		t.Log("LLM visual check: no issues found")
		return
	}

	// Report issues by severity
	var highIssues []string
	for _, issue := range issues {
		msg := fmt.Sprintf("[%s] %s: %s", issue.Severity, issue.Category, issue.Description)
		switch strings.ToUpper(issue.Severity) {
		case "HIGH", "CRITICAL":
			highIssues = append(highIssues, msg)
			t.Errorf("Visual issue: %s", msg)
		default:
			t.Logf("Visual warning: %s", msg)
		}
	}

	if len(highIssues) > 0 {
		t.Errorf("Found %d high-severity visual issues", len(highIssues))
	}
}
