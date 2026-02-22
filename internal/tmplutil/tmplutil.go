// Package tmplutil provides shared helpers for parsing and formatting
// Go html/template errors. It is intentionally small so that both the
// generator and the validation engine can import it without creating a
// circular dependency.
package tmplutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// LinePattern matches Go template parse error positions like
// "template: name:5:" or "template: name:5:22:".
var LinePattern = regexp.MustCompile(`template:.*?:(\d+)`)

// ExtractLineNumber pulls the line number from a Go template parse error.
// Returns 0 if no line number can be found.
func ExtractLineNumber(err error) int {
	m := LinePattern.FindStringSubmatch(err.Error())
	if len(m) < 2 {
		return 0
	}
	n, convErr := strconv.Atoi(m[1])
	if convErr != nil {
		return 0
	}
	return n
}

// SourceContext returns surrounding lines of source around lineNum, with
// line numbers and marker on the target line. marker is typically "→ " or
// "> " depending on the caller's display convention.
func SourceContext(content string, lineNum, surrounding int, marker string) string {
	lines := strings.Split(content, "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return ""
	}

	start := max(lineNum-surrounding-1, 0)
	end := min(lineNum+surrounding, len(lines))

	var b strings.Builder
	for i := start; i < end; i++ {
		m := "  "
		if i+1 == lineNum {
			m = marker
		}
		fmt.Fprintf(&b, "  %s%4d | %s\n", m, i+1, lines[i])
	}
	return strings.TrimRight(b.String(), "\n")
}
