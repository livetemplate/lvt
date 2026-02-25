package knowledge

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Parse reads evolution/patterns.md and returns all patterns with compiled regexes.
func Parse(path string) ([]*Pattern, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open patterns file: %w", err)
	}
	defer f.Close()

	var patterns []*Pattern
	var current *Pattern
	var currentFix *FixTemplate
	section := "" // "description", "error", "fix"
	var descLines []string

	flushFix := func() {
		if current != nil && currentFix != nil && currentFix.File != "" {
			current.Fixes = append(current.Fixes, *currentFix)
			currentFix = nil
		}
	}

	flushPattern := func() {
		flushFix()
		if current != nil {
			if len(descLines) > 0 {
				current.Description = strings.TrimSpace(strings.Join(descLines, "\n"))
			}
			patterns = append(patterns, current)
			current = nil
			descLines = nil
			section = ""
		}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip HTML comments (template section at bottom)
		if strings.HasPrefix(strings.TrimSpace(line), "<!--") {
			flushPattern()
			break
		}

		// New pattern starts
		if strings.HasPrefix(line, "## Pattern: ") {
			flushPattern()
			id := strings.TrimSpace(strings.TrimPrefix(line, "## Pattern: "))
			current = &Pattern{ID: id}
			section = "metadata"
			continue
		}

		// Skip non-pattern H2 sections
		if strings.HasPrefix(line, "## ") && !strings.HasPrefix(line, "## Pattern: ") {
			// e.g. "## Upstream Patterns" — keep current context
			continue
		}

		if current == nil {
			continue
		}

		// Section headers
		if strings.HasPrefix(line, "### Description") {
			section = "description"
			continue
		}
		if strings.HasPrefix(line, "### Error Pattern") {
			section = "error"
			continue
		}
		if strings.HasPrefix(line, "### Fix") {
			flushFix()
			currentFix = &FixTemplate{}
			section = "fix"
			continue
		}

		trimmed := strings.TrimSpace(line)

		switch section {
		case "metadata":
			parseMetadataLine(current, trimmed)

		case "description":
			if trimmed != "" {
				descLines = append(descLines, trimmed)
			}

		case "error":
			if err := parseErrorLine(current, trimmed); err != nil {
				return nil, fmt.Errorf("pattern %q: %w", current.ID, err)
			}

		case "fix":
			if currentFix != nil {
				parseFixLine(currentFix, trimmed)
			}
		}
	}

	// Flush the last pattern
	flushPattern()

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan patterns: %w", err)
	}

	return patterns, nil
}

func parseMetadataLine(p *Pattern, line string) {
	switch {
	case strings.HasPrefix(line, "**Name:**"):
		p.Name = strings.TrimSpace(strings.TrimPrefix(line, "**Name:**"))
	case strings.HasPrefix(line, "**Confidence:**"):
		v := strings.TrimSpace(strings.TrimPrefix(line, "**Confidence:**"))
		p.Confidence, _ = strconv.ParseFloat(v, 64)
	case strings.HasPrefix(line, "**Added:**"):
		p.Added = strings.TrimSpace(strings.TrimPrefix(line, "**Added:**"))
	case strings.HasPrefix(line, "**Fix Count:**"):
		v := strings.TrimSpace(strings.TrimPrefix(line, "**Fix Count:**"))
		p.FixCount, _ = strconv.Atoi(v)
	case strings.HasPrefix(line, "**Success Rate:**"):
		p.SuccessRate = strings.TrimSpace(strings.TrimPrefix(line, "**Success Rate:**"))
	case strings.HasPrefix(line, "**Upstream Repo:**"):
		p.UpstreamRepo = strings.TrimSpace(strings.TrimPrefix(line, "**Upstream Repo:**"))
	}
}

func parseErrorLine(p *Pattern, line string) error {
	switch {
	case strings.HasPrefix(line, "- **Phase:**"):
		p.ErrorPhase = strings.TrimSpace(strings.TrimPrefix(line, "- **Phase:**"))
	case strings.HasPrefix(line, "- **Message Regex:**"):
		raw := extractBacktickValue(strings.TrimPrefix(line, "- **Message Regex:**"))
		re, err := regexp.Compile(raw)
		if err != nil {
			return fmt.Errorf("compile message regex %q: %w", raw, err)
		}
		p.MessageRe = re
	case strings.HasPrefix(line, "- **Context Regex:**"):
		raw := extractBacktickValue(strings.TrimPrefix(line, "- **Context Regex:**"))
		re, err := regexp.Compile(raw)
		if err != nil {
			return fmt.Errorf("compile context regex %q: %w", raw, err)
		}
		p.ContextRe = re
	}
	return nil
}

func parseFixLine(fix *FixTemplate, line string) {
	switch {
	case strings.HasPrefix(line, "- **File:**"):
		fix.File = extractBacktickValue(strings.TrimPrefix(line, "- **File:**"))
	case strings.HasPrefix(line, "- **Find:**"):
		fix.FindPattern = extractBacktickValue(strings.TrimPrefix(line, "- **Find:**"))
	case strings.HasPrefix(line, "- **Replace:**"):
		fix.Replace = extractBacktickValue(strings.TrimPrefix(line, "- **Replace:**"))
	case strings.HasPrefix(line, "- **Is Regex:**"):
		v := strings.TrimSpace(strings.TrimPrefix(line, "- **Is Regex:**"))
		fix.IsRegex = strings.ToLower(v) == "true"
	}
}

// extractBacktickValue extracts the value between backticks, or returns trimmed text.
func extractBacktickValue(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "`") && strings.HasSuffix(s, "`") {
		return s[1 : len(s)-1]
	}
	// Handle empty backticks (e.g., Replace: ``)
	if s == "``" {
		return ""
	}
	return s
}
