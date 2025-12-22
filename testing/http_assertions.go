package testing

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// HTTPAssert provides assertion helpers for HTTP responses.
type HTTPAssert struct {
	Response *HTTPResponse
}

// NewHTTPAssert creates an assertion helper for the given HTTP response.
//
// Example:
//
//	resp := test.Get("/")
//	assert := lvttest.NewHTTPAssert(resp)
//	assert.StatusOK(t)
//	assert.Contains(t, "Welcome")
func NewHTTPAssert(resp *HTTPResponse) *HTTPAssert {
	return &HTTPAssert{Response: resp}
}

// StatusOK asserts that the response has a 200 OK status.
func (a *HTTPAssert) StatusOK(t *testing.T) {
	t.Helper()
	if a.Response.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", a.Response.StatusCode)
	}
}

// StatusCode asserts that the response has the expected status code.
func (a *HTTPAssert) StatusCode(t *testing.T, expected int) {
	t.Helper()
	if a.Response.StatusCode != expected {
		t.Errorf("expected status %d, got %d", expected, a.Response.StatusCode)
	}
}

// StatusRedirect asserts that the response is a redirect (3xx).
func (a *HTTPAssert) StatusRedirect(t *testing.T) {
	t.Helper()
	if a.Response.StatusCode < 300 || a.Response.StatusCode >= 400 {
		t.Errorf("expected redirect (3xx), got %d", a.Response.StatusCode)
	}
}

// RedirectTo asserts that the response redirects to the expected location.
func (a *HTTPAssert) RedirectTo(t *testing.T, expectedLocation string) {
	t.Helper()
	a.StatusRedirect(t)
	location := a.Response.Header.Get("Location")
	if location != expectedLocation {
		t.Errorf("expected redirect to %q, got %q", expectedLocation, location)
	}
}

// StatusNotFound asserts that the response has a 404 status.
func (a *HTTPAssert) StatusNotFound(t *testing.T) {
	t.Helper()
	if a.Response.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", a.Response.StatusCode)
	}
}

// StatusUnauthorized asserts that the response has a 401 status.
func (a *HTTPAssert) StatusUnauthorized(t *testing.T) {
	t.Helper()
	if a.Response.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", a.Response.StatusCode)
	}
}

// StatusForbidden asserts that the response has a 403 status.
func (a *HTTPAssert) StatusForbidden(t *testing.T) {
	t.Helper()
	if a.Response.StatusCode != 403 {
		t.Errorf("expected status 403, got %d", a.Response.StatusCode)
	}
}

// StatusBadRequest asserts that the response has a 400 status.
func (a *HTTPAssert) StatusBadRequest(t *testing.T) {
	t.Helper()
	if a.Response.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", a.Response.StatusCode)
	}
}

// StatusServerError asserts that the response has a 5xx status.
func (a *HTTPAssert) StatusServerError(t *testing.T) {
	t.Helper()
	if a.Response.StatusCode < 500 || a.Response.StatusCode >= 600 {
		t.Errorf("expected server error (5xx), got %d", a.Response.StatusCode)
	}
}

// Contains asserts that the response body contains the expected text.
func (a *HTTPAssert) Contains(t *testing.T, text string) {
	t.Helper()
	a.Response.readBody()
	if !strings.Contains(string(a.Response.Body), text) {
		t.Errorf("response does not contain %q\n\nBody:\n%s", text, truncateBody(a.Response.Body))
	}
}

// NotContains asserts that the response body does NOT contain the text.
func (a *HTTPAssert) NotContains(t *testing.T, text string) {
	t.Helper()
	a.Response.readBody()
	if strings.Contains(string(a.Response.Body), text) {
		t.Errorf("response should not contain %q but it does", text)
	}
}

// ContainsAll asserts that the response body contains all expected texts.
func (a *HTTPAssert) ContainsAll(t *testing.T, texts ...string) {
	t.Helper()
	a.Response.readBody()
	body := string(a.Response.Body)
	var missing []string
	for _, text := range texts {
		if !strings.Contains(body, text) {
			missing = append(missing, text)
		}
	}
	if len(missing) > 0 {
		t.Errorf("response missing: %v", missing)
	}
}

// Matches asserts that the response body matches the regular expression.
func (a *HTTPAssert) Matches(t *testing.T, pattern string) {
	t.Helper()
	a.Response.readBody()
	re, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("invalid pattern %q: %v", pattern, err)
	}
	if !re.Match(a.Response.Body) {
		t.Errorf("response does not match pattern %q\n\nBody:\n%s", pattern, truncateBody(a.Response.Body))
	}
}

// NoTemplateErrors asserts that the response has no unflattened template expressions.
// This catches bugs where {{.Field}}, {{if}}, {{range}}, etc. appear in the output.
func (a *HTTPAssert) NoTemplateErrors(t *testing.T) {
	t.Helper()
	if a.Response.HasTemplateErrors() {
		errors := a.Response.FindTemplateErrors()
		t.Errorf("found unflattened template expressions: %v\n\nBody:\n%s", errors, truncateBody(a.Response.Body))
	}
}

// HasElement asserts that at least one element matches the CSS selector.
// Note: This uses a simple implementation - for complex selectors, use a full CSS selector library.
func (a *HTTPAssert) HasElement(t *testing.T, selector string) {
	t.Helper()
	a.Response.readBody()

	if !a.hasElement(selector) {
		t.Errorf("element %q not found in response", selector)
	}
}

// HasNoElement asserts that no elements match the CSS selector.
func (a *HTTPAssert) HasNoElement(t *testing.T, selector string) {
	t.Helper()
	a.Response.readBody()

	if a.hasElement(selector) {
		t.Errorf("element %q should not exist but it does", selector)
	}
}

// ElementCount asserts that exactly n elements match the selector.
func (a *HTTPAssert) ElementCount(t *testing.T, selector string, expected int) {
	t.Helper()
	a.Response.readBody()

	count := a.countElements(selector)
	if count != expected {
		t.Errorf("expected %d elements matching %q, got %d", expected, selector, count)
	}
}

// ElementText asserts that an element has the expected text content.
func (a *HTTPAssert) ElementText(t *testing.T, selector string, expectedText string) {
	t.Helper()
	a.Response.readBody()

	text := a.getElementText(selector)
	text = strings.TrimSpace(text)
	expectedText = strings.TrimSpace(expectedText)

	if text != expectedText {
		t.Errorf("element %q has text %q, expected %q", selector, text, expectedText)
	}
}

// ElementTextContains asserts that an element's text contains the expected substring.
func (a *HTTPAssert) ElementTextContains(t *testing.T, selector string, expectedSubstring string) {
	t.Helper()
	a.Response.readBody()

	text := a.getElementText(selector)
	if !strings.Contains(text, expectedSubstring) {
		t.Errorf("element %q text %q does not contain %q", selector, text, expectedSubstring)
	}
}

// HasCSRFToken asserts that the response contains a CSRF token field.
func (a *HTTPAssert) HasCSRFToken(t *testing.T) {
	t.Helper()
	a.Response.readBody()

	// Look for common CSRF token field names
	patterns := []string{
		`name="csrf_token"`,
		`name="_csrf"`,
		`name="gorilla.csrf.Token"`,
	}

	body := string(a.Response.Body)
	for _, pattern := range patterns {
		if strings.Contains(body, pattern) {
			return
		}
	}

	t.Error("CSRF token field not found in response")
}

// FormFieldValue asserts that a form field has the expected value.
func (a *HTTPAssert) FormFieldValue(t *testing.T, fieldName string, expectedValue string) {
	t.Helper()
	a.Response.readBody()

	value := a.getInputValue(fieldName)
	if value != expectedValue {
		t.Errorf("field %q has value %q, expected %q", fieldName, value, expectedValue)
	}
}

// HasFormField asserts that a form field with the given name exists.
func (a *HTTPAssert) HasFormField(t *testing.T, fieldName string) {
	t.Helper()
	a.Response.readBody()

	if !a.hasInputField(fieldName) {
		t.Errorf("form field %q not found", fieldName)
	}
}

// ContentType asserts that the response has the expected Content-Type.
func (a *HTTPAssert) ContentType(t *testing.T, expected string) {
	t.Helper()
	ct := a.Response.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, expected) {
		t.Errorf("expected Content-Type %q, got %q", expected, ct)
	}
}

// ContentTypeHTML asserts that the response is HTML.
func (a *HTTPAssert) ContentTypeHTML(t *testing.T) {
	t.Helper()
	a.ContentType(t, "text/html")
}

// ContentTypeJSON asserts that the response is JSON.
func (a *HTTPAssert) ContentTypeJSON(t *testing.T) {
	t.Helper()
	a.ContentType(t, "application/json")
}

// Header asserts that the response has the expected header value.
func (a *HTTPAssert) Header(t *testing.T, name, expected string) {
	t.Helper()
	actual := a.Response.Header.Get(name)
	if actual != expected {
		t.Errorf("expected header %q=%q, got %q", name, expected, actual)
	}
}

// HasHeader asserts that the response has the specified header.
func (a *HTTPAssert) HasHeader(t *testing.T, name string) {
	t.Helper()
	if a.Response.Header.Get(name) == "" {
		t.Errorf("expected header %q to be present", name)
	}
}

// TableRowCount asserts that a table has the expected number of rows.
func (a *HTTPAssert) TableRowCount(t *testing.T, expected int) {
	t.Helper()
	a.ElementCount(t, "tbody tr", expected)
}

// helper functions

func truncateBody(body []byte) string {
	const maxLen = 1000
	s := string(body)
	if len(s) > maxLen {
		return s[:maxLen] + "... (truncated)"
	}
	return s
}

// hasElement checks if an element matching the selector exists.
// Simple implementation supporting: tag, .class, #id, tag.class, tag#id
func (a *HTTPAssert) hasElement(selector string) bool {
	doc, err := html.Parse(bytes.NewReader(a.Response.Body))
	if err != nil {
		return false
	}

	return a.findElement(doc, selector) != nil
}

// countElements counts elements matching the selector.
func (a *HTTPAssert) countElements(selector string) int {
	doc, err := html.Parse(bytes.NewReader(a.Response.Body))
	if err != nil {
		return 0
	}

	return a.findAllElements(doc, selector)
}

// getElementText gets the text content of the first matching element.
func (a *HTTPAssert) getElementText(selector string) string {
	doc, err := html.Parse(bytes.NewReader(a.Response.Body))
	if err != nil {
		return ""
	}

	el := a.findElement(doc, selector)
	if el == nil {
		return ""
	}

	return extractText(el)
}

// getInputValue gets the value of an input field by name.
func (a *HTTPAssert) getInputValue(fieldName string) string {
	doc, err := html.Parse(bytes.NewReader(a.Response.Body))
	if err != nil {
		return ""
	}

	var value string
	var findInput func(*html.Node)
	findInput = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "input" || n.Data == "textarea" || n.Data == "select") {
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == fieldName {
					// For input, get value attribute
					if n.Data == "input" {
						for _, a := range n.Attr {
							if a.Key == "value" {
								value = a.Val
								return
							}
						}
					}
					// For textarea, get text content
					if n.Data == "textarea" {
						value = extractText(n)
						return
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findInput(c)
		}
	}
	findInput(doc)

	return value
}

// hasInputField checks if an input field with the given name exists.
func (a *HTTPAssert) hasInputField(fieldName string) bool {
	doc, err := html.Parse(bytes.NewReader(a.Response.Body))
	if err != nil {
		return false
	}

	var found bool
	var findInput func(*html.Node)
	findInput = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "input" || n.Data == "textarea" || n.Data == "select") {
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == fieldName {
					found = true
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found {
				return
			}
			findInput(c)
		}
	}
	findInput(doc)

	return found
}

// findElement finds the first element matching the selector.
func (a *HTTPAssert) findElement(n *html.Node, selector string) *html.Node {
	tag, id, class := parseSelector(selector)

	var find func(*html.Node) *html.Node
	find = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode {
			if matchesSelector(n, tag, id, class) {
				return n
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found := find(c); found != nil {
				return found
			}
		}
		return nil
	}

	return find(n)
}

// findAllElements counts all elements matching the selector.
func (a *HTTPAssert) findAllElements(n *html.Node, selector string) int {
	tag, id, class := parseSelector(selector)

	var count int
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if matchesSelector(n, tag, id, class) {
				count++
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}

	find(n)
	return count
}

// parseSelector parses a simple CSS selector into tag, id, and class.
// Supports: tag, .class, #id, tag.class, tag#id, tag.class1.class2
func parseSelector(selector string) (tag, id, class string) {
	// Handle descendant selectors (e.g., "tbody tr")
	parts := strings.Fields(selector)
	if len(parts) > 1 {
		// For multi-part selectors, only match the last part
		selector = parts[len(parts)-1]
	}

	// Parse the selector
	if strings.HasPrefix(selector, "#") {
		id = selector[1:]
	} else if strings.HasPrefix(selector, ".") {
		class = selector[1:]
	} else if strings.Contains(selector, "#") {
		idx := strings.Index(selector, "#")
		tag = selector[:idx]
		id = selector[idx+1:]
	} else if strings.Contains(selector, ".") {
		idx := strings.Index(selector, ".")
		tag = selector[:idx]
		class = selector[idx+1:]
	} else {
		tag = selector
	}

	return
}

// matchesSelector checks if a node matches the tag, id, and class.
func matchesSelector(n *html.Node, tag, id, class string) bool {
	if tag != "" && n.Data != tag {
		return false
	}

	if id != "" {
		hasID := false
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				hasID = true
				break
			}
		}
		if !hasID {
			return false
		}
	}

	if class != "" {
		hasClass := false
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				classes := strings.Fields(attr.Val)
				for _, c := range classes {
					if c == class {
						hasClass = true
						break
					}
				}
			}
		}
		if !hasClass {
			return false
		}
	}

	return true
}

// extractText extracts all text content from a node.
func extractText(n *html.Node) string {
	var buf strings.Builder
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(n)
	return buf.String()
}

// DatabaseAssert provides database assertion helpers.
type DatabaseAssert struct {
	test *HTTPTest
}

// NewDatabaseAssert creates a database assertion helper.
func NewDatabaseAssert(test *HTTPTest) *DatabaseAssert {
	return &DatabaseAssert{test: test}
}

// RecordExists asserts that a record exists in the table matching the conditions.
func (d *DatabaseAssert) RecordExists(t *testing.T, table string, conditions map[string]interface{}) {
	t.Helper()

	if d.test.DB == nil {
		t.Fatal("database connection not set - use HTTPSetupOptions.DB or HTTPTest.SetDB()")
	}

	query, args := buildSelectQuery(table, conditions)
	var count int
	err := d.test.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		t.Fatalf("database query failed: %v", err)
	}

	if count == 0 {
		t.Errorf("no record found in %s matching %v", table, conditions)
	}
}

// RecordNotExists asserts that no record exists in the table matching the conditions.
func (d *DatabaseAssert) RecordNotExists(t *testing.T, table string, conditions map[string]interface{}) {
	t.Helper()

	if d.test.DB == nil {
		t.Fatal("database connection not set - use HTTPSetupOptions.DB or HTTPTest.SetDB()")
	}

	query, args := buildSelectQuery(table, conditions)
	var count int
	err := d.test.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		t.Fatalf("database query failed: %v", err)
	}

	if count > 0 {
		t.Errorf("expected no record in %s matching %v, but found %d", table, conditions, count)
	}
}

// RecordCount asserts that the table has exactly the expected number of records.
func (d *DatabaseAssert) RecordCount(t *testing.T, table string, expected int) {
	t.Helper()

	if d.test.DB == nil {
		t.Fatal("database connection not set - use HTTPSetupOptions.DB or HTTPTest.SetDB()")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	var count int
	err := d.test.DB.QueryRow(query).Scan(&count)
	if err != nil {
		t.Fatalf("database query failed: %v", err)
	}

	if count != expected {
		t.Errorf("expected %d records in %s, got %d", expected, table, count)
	}
}

// RecordDeleted asserts that a record with the given ID was deleted.
func (d *DatabaseAssert) RecordDeleted(t *testing.T, table string, id interface{}) {
	t.Helper()
	d.RecordNotExists(t, table, map[string]interface{}{"id": id})
}

func buildSelectQuery(table string, conditions map[string]interface{}) (string, []interface{}) {
	var whereClauses []string
	var args []interface{}
	i := 1

	for col, val := range conditions {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	where := ""
	if len(whereClauses) > 0 {
		where = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return fmt.Sprintf("SELECT COUNT(*) FROM %s%s", table, where), args
}
