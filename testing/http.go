// Package testing provides test utilities for LiveTemplate applications.
//
// HTTPTest provides a lightweight alternative to browser-based E2E tests.
// It starts the test server and makes HTTP requests directly, without requiring
// a browser. This is ideal for testing:
//   - Server-side rendering
//   - Form submission and validation
//   - CRUD operations
//   - API responses
//   - Template expression validation
//
// For tests that require browser JavaScript execution (WebSocket, focus management,
// animations), use the browser-based E2ETest instead.
package testing

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

// HTTPTest represents a configured HTTP test environment.
// It starts the test server and provides HTTP client methods for testing.
type HTTPTest struct {
	T          *testing.T
	Port       int
	Client     *http.Client
	DB         *sql.DB
	AppDir     string
	AppPath    string
	BaseURL    string
	ServerCmd  *exec.Cmd
	Server     *ServerLogger
	cleanupFns []func()
}

// HTTPSetupOptions configures the HTTP test environment.
type HTTPSetupOptions struct {
	// AppPath is the path to main.go (e.g., "./cmd/myapp/main.go")
	AppPath string

	// AppDir is the working directory for the app (defaults to directory of AppPath)
	AppDir string

	// Port is the server port (auto-allocated if 0)
	Port int

	// Timeout is the HTTP client timeout (default 10s)
	Timeout time.Duration

	// DB is an optional database connection for state verification
	DB *sql.DB
}

// SetupHTTP creates a new HTTP test environment.
// It starts the test server and returns an HTTPTest for making requests.
//
// Example:
//
//	test := lvttest.SetupHTTP(t, &lvttest.HTTPSetupOptions{
//	    AppPath: "../../cmd/myapp/main.go",
//	})
//	defer test.Cleanup()
//
//	resp := test.Get("/")
//	assert := lvttest.NewHTTPAssert(resp)
//	assert.StatusOK(t)
func SetupHTTP(t *testing.T, opts *HTTPSetupOptions) *HTTPTest {
	t.Helper()

	if opts == nil {
		opts = &HTTPSetupOptions{}
	}
	if opts.AppPath == "" {
		t.Fatal("AppPath is required in HTTPSetupOptions")
	}
	if opts.Timeout == 0 {
		opts.Timeout = 10 * time.Second
	}

	// Determine app directory
	appDir := opts.AppDir
	if appDir == "" {
		appDir = filepath.Dir(opts.AppPath)
	}

	// Allocate port
	port := opts.Port
	if port == 0 {
		var err error
		port, err = GetFreePort()
		if err != nil {
			t.Fatalf("Failed to allocate port: %v", err)
		}
	}

	// Create cookie jar for session persistence
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("Failed to create cookie jar: %v", err)
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: opts.Timeout,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects automatically - let tests handle them
			return http.ErrUseLastResponse
		},
	}

	// Start server
	serverLogger := NewServerLogger()
	serverLogger.Start()

	serverCmd := StartTestServer(t, opts.AppPath, port)

	test := &HTTPTest{
		T:         t,
		Port:      port,
		Client:    client,
		DB:        opts.DB,
		AppDir:    appDir,
		AppPath:   opts.AppPath,
		BaseURL:   fmt.Sprintf("http://localhost:%d", port),
		ServerCmd: serverCmd,
		Server:    serverLogger,
	}

	// Register cleanup
	t.Cleanup(test.Cleanup)

	return test
}

// Cleanup tears down the test environment.
// This is called automatically via t.Cleanup(), but can be called manually if needed.
func (h *HTTPTest) Cleanup() {
	// Run custom cleanup functions in reverse order
	for i := len(h.cleanupFns) - 1; i >= 0; i-- {
		h.cleanupFns[i]()
	}

	// Stop server logger
	if h.Server != nil {
		h.Server.Stop()
	}

	// Kill server process
	if h.ServerCmd != nil && h.ServerCmd.Process != nil {
		_ = h.ServerCmd.Process.Kill()
		_ = h.ServerCmd.Wait()
	}
}

// OnCleanup registers a function to be called during cleanup.
func (h *HTTPTest) OnCleanup(fn func()) {
	h.cleanupFns = append(h.cleanupFns, fn)
}

// URL returns the full URL for the given path.
func (h *HTTPTest) URL(path string) string {
	if path == "" || path == "/" {
		return h.BaseURL
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return h.BaseURL + path
}

// HTTPResponse wraps an HTTP response with helper methods.
type HTTPResponse struct {
	*http.Response
	Body    []byte
	bodyErr error
}

// readBody reads and caches the response body.
func (r *HTTPResponse) readBody() {
	if r.Body != nil || r.bodyErr != nil {
		return
	}
	r.Body, r.bodyErr = io.ReadAll(r.Response.Body)
	r.Response.Body.Close()
}

// String returns the response body as a string.
func (r *HTTPResponse) String() string {
	r.readBody()
	return string(r.Body)
}

// Get performs an HTTP GET request and returns the response.
func (h *HTTPTest) Get(path string) *HTTPResponse {
	h.T.Helper()

	resp, err := h.Client.Get(h.URL(path))
	if err != nil {
		h.T.Fatalf("GET %s failed: %v", path, err)
	}

	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	return httpResp
}

// GetWithHeaders performs an HTTP GET request with custom headers.
func (h *HTTPTest) GetWithHeaders(path string, headers map[string]string) *HTTPResponse {
	h.T.Helper()

	req, err := http.NewRequest("GET", h.URL(path), nil)
	if err != nil {
		h.T.Fatalf("Failed to create request: %v", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		h.T.Fatalf("GET %s failed: %v", path, err)
	}

	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	return httpResp
}

// PostForm submits form data via POST and returns the response.
func (h *HTTPTest) PostForm(path string, data url.Values) *HTTPResponse {
	h.T.Helper()

	resp, err := h.Client.PostForm(h.URL(path), data)
	if err != nil {
		h.T.Fatalf("POST %s failed: %v", path, err)
	}

	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	return httpResp
}

// PostJSON submits JSON data via POST and returns the response.
func (h *HTTPTest) PostJSON(path string, data interface{}) *HTTPResponse {
	h.T.Helper()

	body, err := json.Marshal(data)
	if err != nil {
		h.T.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", h.URL(path), bytes.NewReader(body))
	if err != nil {
		h.T.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		h.T.Fatalf("POST %s failed: %v", path, err)
	}

	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	return httpResp
}

// PostMultipart submits multipart form data (for file uploads) and returns the response.
func (h *HTTPTest) PostMultipart(path string, fields map[string]string, files map[string][]byte) *HTTPResponse {
	h.T.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			h.T.Fatalf("Failed to write field %s: %v", key, err)
		}
	}

	// Add files
	for fieldName, fileData := range files {
		part, err := writer.CreateFormFile(fieldName, fieldName)
		if err != nil {
			h.T.Fatalf("Failed to create form file %s: %v", fieldName, err)
		}
		if _, err := part.Write(fileData); err != nil {
			h.T.Fatalf("Failed to write file data for %s: %v", fieldName, err)
		}
	}

	if err := writer.Close(); err != nil {
		h.T.Fatalf("Failed to close multipart writer: %v", err)
	}

	req, err := http.NewRequest("POST", h.URL(path), &buf)
	if err != nil {
		h.T.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := h.Client.Do(req)
	if err != nil {
		h.T.Fatalf("POST %s failed: %v", path, err)
	}

	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	return httpResp
}

// Delete performs an HTTP DELETE request and returns the response.
func (h *HTTPTest) Delete(path string) *HTTPResponse {
	h.T.Helper()

	req, err := http.NewRequest("DELETE", h.URL(path), nil)
	if err != nil {
		h.T.Fatalf("Failed to create request: %v", err)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		h.T.Fatalf("DELETE %s failed: %v", path, err)
	}

	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	return httpResp
}

// FollowRedirect follows a redirect response and returns the new response.
// Returns nil if the response is not a redirect.
func (h *HTTPTest) FollowRedirect(resp *HTTPResponse) *HTTPResponse {
	h.T.Helper()

	location := resp.Header.Get("Location")
	if location == "" {
		return nil
	}

	return h.Get(location)
}

// FollowRedirects follows all redirects until a non-redirect response.
func (h *HTTPTest) FollowRedirects(resp *HTTPResponse) *HTTPResponse {
	h.T.Helper()

	for i := 0; i < 10; i++ { // Max 10 redirects
		location := resp.Header.Get("Location")
		if location == "" || resp.StatusCode < 300 || resp.StatusCode >= 400 {
			return resp
		}
		resp = h.Get(location)
	}

	h.T.Fatal("Too many redirects")
	return nil
}

// ExtractCSRFToken extracts a CSRF token from an HTML response.
// It looks for a hidden input field named "csrf_token" or "_csrf".
func (h *HTTPTest) ExtractCSRFToken(resp *HTTPResponse) string {
	h.T.Helper()

	doc, err := html.Parse(bytes.NewReader(resp.Body))
	if err != nil {
		return ""
	}

	var token string
	var findToken func(*html.Node)
	findToken = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			var name, value string
			var isHidden bool
			for _, attr := range n.Attr {
				switch attr.Key {
				case "name":
					name = attr.Val
				case "value":
					value = attr.Val
				case "type":
					isHidden = attr.Val == "hidden"
				}
			}
			if isHidden && (name == "csrf_token" || name == "_csrf" || name == "gorilla.csrf.Token") {
				token = value
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findToken(c)
		}
	}
	findToken(doc)

	return token
}

// SubmitForm extracts a form from the response, fills in values, and submits it.
// It automatically handles CSRF tokens.
func (h *HTTPTest) SubmitForm(resp *HTTPResponse, formSelector string, values map[string]string) *HTTPResponse {
	h.T.Helper()

	// Parse HTML to find form
	doc, err := html.Parse(bytes.NewReader(resp.Body))
	if err != nil {
		h.T.Fatalf("Failed to parse HTML: %v", err)
	}

	// Find form and extract action, method, and existing fields
	var formAction, formMethod string
	existingFields := make(map[string]string)

	var findForm func(*html.Node) *html.Node
	findForm = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "form" {
			// Check if this matches selector (simple implementation)
			for _, attr := range n.Attr {
				if attr.Key == "id" && formSelector == "#"+attr.Val {
					return n
				}
				if attr.Key == "action" {
					formAction = attr.Val
				}
				if attr.Key == "method" {
					formMethod = strings.ToUpper(attr.Val)
				}
			}
			// If no selector specified or matches first form
			if formSelector == "" || formSelector == "form" {
				return n
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found := findForm(c); found != nil {
				return found
			}
		}
		return nil
	}

	form := findForm(doc)
	if form == nil {
		h.T.Fatalf("Form not found: %s", formSelector)
	}

	// Extract form attributes
	for _, attr := range form.Attr {
		switch attr.Key {
		case "action":
			formAction = attr.Val
		case "method":
			formMethod = strings.ToUpper(attr.Val)
		}
	}

	if formMethod == "" {
		formMethod = "POST"
	}

	// Extract existing field values (hidden fields, CSRF, etc.)
	var extractFields func(*html.Node)
	extractFields = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			var name, value, inputType string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "name":
					name = attr.Val
				case "value":
					value = attr.Val
				case "type":
					inputType = attr.Val
				}
			}
			if name != "" && (inputType == "hidden" || inputType == "") {
				existingFields[name] = value
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractFields(c)
		}
	}
	extractFields(form)

	// Merge existing fields with provided values
	formData := url.Values{}
	for k, v := range existingFields {
		formData.Set(k, v)
	}
	for k, v := range values {
		formData.Set(k, v)
	}

	// Submit form
	if formAction == "" {
		formAction = resp.Request.URL.Path
	}

	return h.PostForm(formAction, formData)
}

// SetCookie sets a cookie for subsequent requests.
func (h *HTTPTest) SetCookie(name, value string) {
	h.T.Helper()

	u, err := url.Parse(h.BaseURL)
	if err != nil {
		h.T.Fatalf("Failed to parse base URL: %v", err)
	}

	h.Client.Jar.SetCookies(u, []*http.Cookie{
		{Name: name, Value: value},
	})
}

// GetCookie returns the value of a cookie, or empty string if not found.
func (h *HTTPTest) GetCookie(name string) string {
	u, err := url.Parse(h.BaseURL)
	if err != nil {
		return ""
	}

	for _, cookie := range h.Client.Jar.Cookies(u) {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

// ClearCookies removes all cookies.
func (h *HTTPTest) ClearCookies() {
	h.T.Helper()

	jar, err := cookiejar.New(nil)
	if err != nil {
		h.T.Fatalf("Failed to create cookie jar: %v", err)
	}
	h.Client.Jar = jar
}

// WaitForServer waits for the server to be ready by polling the health endpoint.
func (h *HTTPTest) WaitForServer(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	interval := 50 * time.Millisecond

	for time.Now().Before(deadline) {
		resp, err := h.Client.Get(h.URL("/health"))
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(interval)
	}

	return fmt.Errorf("server not ready after %v", timeout)
}

// SetDB sets the database connection for state verification.
func (h *HTTPTest) SetDB(db *sql.DB) {
	h.DB = db
}

// DBPath returns the path to the test database.
// This assumes TEST_MODE=1 is set, which uses :memory: by default.
func (h *HTTPTest) DBPath() string {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = ":memory:"
	}
	return dbPath
}

// templateErrorPatterns are patterns that indicate unflattened Go template expressions
var templateErrorPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\{\{\s*\.`),          // {{.Field}}
	regexp.MustCompile(`\{\{\s*if\s`),        // {{if ...}}
	regexp.MustCompile(`\{\{\s*range\s`),     // {{range ...}}
	regexp.MustCompile(`\{\{\s*end\s*\}\}`),  // {{end}}
	regexp.MustCompile(`\{\{\s*else\s*\}\}`), // {{else}}
	regexp.MustCompile(`\{\{\s*template\s`),  // {{template ...}}
	regexp.MustCompile(`\{\{\s*block\s`),     // {{block ...}}
	regexp.MustCompile(`\{\{\s*with\s`),      // {{with ...}}
	regexp.MustCompile(`\{\{\s*define\s`),    // {{define ...}}
	regexp.MustCompile(`\[\[\s*\.`),          // [[.Field]] (alternate delimiters)
	regexp.MustCompile(`\[\[\s*if\s`),        // [[if ...]]
	regexp.MustCompile(`\[\[\s*range\s`),     // [[range ...]]
}

// HasTemplateErrors checks if the response body contains unflattened template expressions.
func (r *HTTPResponse) HasTemplateErrors() bool {
	r.readBody()
	body := string(r.Body)

	for _, pattern := range templateErrorPatterns {
		if pattern.MatchString(body) {
			return true
		}
	}
	return false
}

// FindTemplateErrors returns all template error matches found in the response.
func (r *HTTPResponse) FindTemplateErrors() []string {
	r.readBody()
	body := string(r.Body)

	var errors []string
	for _, pattern := range templateErrorPatterns {
		matches := pattern.FindAllString(body, -1)
		errors = append(errors, matches...)
	}
	return errors
}
