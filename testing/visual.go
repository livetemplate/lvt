package testing

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/chromedp/chromedp"
)

// VisualIssue represents a visual problem found by the LLM reviewer.
type VisualIssue struct {
	Severity    string `json:"severity"`    // HIGH, MEDIUM, LOW
	Category    string `json:"category"`    // ALIGNMENT, HIERARCHY, SPACING, ERROR_STATE, LAYOUT, READABILITY
	Description string `json:"description"` // Human-readable description
}

const visualCheckPrompt = `You are a UI/UX reviewer for a web application using Pico CSS (https://picocss.com).
Analyze this screenshot for visual issues.

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
// Claude for visual analysis. Skipped unless LVT_VISUAL_CHECK=true is set.
// Requires ANTHROPIC_API_KEY environment variable.
//
// Fails the test if any HIGH or CRITICAL severity issues are found.
// MEDIUM and LOW issues are logged as warnings.
func ValidateScreenshotWithLLM(t *testing.T, ctx context.Context, pageDescription string) {
	t.Helper()

	if os.Getenv("LVT_VISUAL_CHECK") != "true" {
		t.Skip("LLM visual check disabled — set LVT_VISUAL_CHECK=true to enable")
	}

	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("ANTHROPIC_API_KEY not set — required for LLM visual check")
	}

	// Capture screenshot
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf)); err != nil {
		t.Fatalf("Failed to capture screenshot: %v", err)
	}

	t.Logf("Captured screenshot: %d bytes", len(buf))

	// Encode as base64 for the API
	b64 := base64.StdEncoding.EncodeToString(buf)

	// Call Anthropic API
	client := anthropic.NewClient()
	prompt := fmt.Sprintf(visualCheckPrompt, pageDescription)

	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewImageBlockBase64("image/png", b64),
				anthropic.NewTextBlock(prompt),
			),
		},
	})
	if err != nil {
		t.Fatalf("Anthropic API call failed: %v", err)
	}

	// Extract text response
	var responseText string
	for _, block := range message.Content {
		if block.Type == "text" {
			responseText = block.Text
			break
		}
	}

	if responseText == "" {
		t.Fatal("Empty response from LLM")
	}

	t.Logf("LLM response: %s", responseText)

	// Parse JSON — strip markdown code fence if present
	jsonStr := strings.TrimSpace(responseText)
	if strings.HasPrefix(jsonStr, "```") {
		lines := strings.Split(jsonStr, "\n")
		// Remove first and last lines (code fence markers)
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
