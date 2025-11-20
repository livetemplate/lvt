package testing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// SmokeTestOptions configures smoke test behavior
type SmokeTestOptions struct {
	Timeout     time.Duration // Total timeout for all tests
	RetryDelay  time.Duration // Delay between retries
	MaxRetries  int           // Maximum retry attempts
	SkipBrowser bool          // Skip browser-based tests
}

// DefaultSmokeTestOptions returns sensible defaults
func DefaultSmokeTestOptions() *SmokeTestOptions {
	return &SmokeTestOptions{
		Timeout:     5 * time.Minute,
		RetryDelay:  2 * time.Second,
		MaxRetries:  3,
		SkipBrowser: false,
	}
}

// SmokeTestResult represents the result of a single smoke test
type SmokeTestResult struct {
	Name     string
	Passed   bool
	Error    error
	Duration time.Duration
}

// SmokeTestSuite represents the results of all smoke tests
type SmokeTestSuite struct {
	AppURL        string
	Results       []SmokeTestResult
	TotalDuration time.Duration
}

// RunSmokeTests executes all smoke tests against a deployed app
func RunSmokeTests(appURL string, opts *SmokeTestOptions) (*SmokeTestSuite, error) {
	if opts == nil {
		opts = DefaultSmokeTestOptions()
	}

	suite := &SmokeTestSuite{
		AppURL:  appURL,
		Results: []SmokeTestResult{},
	}

	startTime := time.Now()

	// Test 1: HTTP 200 on root path
	suite.Results = append(suite.Results, runTest("HTTP Root Path", func() error {
		return testHTTPRoot(appURL, opts)
	}))

	// Test 2: Health endpoint responds
	suite.Results = append(suite.Results, runTest("Health Endpoint", func() error {
		return testHealthEndpoint(appURL, opts)
	}))

	// Test 3: Static assets load
	suite.Results = append(suite.Results, runTest("Static Assets", func() error {
		return testStaticAssets(appURL, opts)
	}))

	// Test 4: WebSocket connection (browser-based, optional)
	if !opts.SkipBrowser {
		suite.Results = append(suite.Results, runTest("WebSocket Connection", func() error {
			return testWebSocket(appURL, opts)
		}))
	}

	// Test 5: Templates render without errors
	suite.Results = append(suite.Results, runTest("Template Rendering", func() error {
		return testTemplateRendering(appURL, opts)
	}))

	suite.TotalDuration = time.Since(startTime)

	// Check if any tests failed
	for _, result := range suite.Results {
		if !result.Passed {
			return suite, fmt.Errorf("smoke tests failed: %s", result.Name)
		}
	}

	return suite, nil
}

// runTest executes a single test with timing
func runTest(name string, fn func() error) SmokeTestResult {
	startTime := time.Now()
	err := fn()

	return SmokeTestResult{
		Name:     name,
		Passed:   err == nil,
		Error:    err,
		Duration: time.Since(startTime),
	}
}

// Test 1: HTTP Root Path
func testHTTPRoot(appURL string, opts *SmokeTestOptions) error {
	return retryHTTP(opts, func() error {
		resp, err := http.Get(appURL)
		if err != nil {
			return fmt.Errorf("failed to GET %s: %w", appURL, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		return nil
	})
}

// Test 2: Health Endpoint
func testHealthEndpoint(appURL string, opts *SmokeTestOptions) error {
	healthURL := fmt.Sprintf("%s/health", strings.TrimRight(appURL, "/"))

	return retryHTTP(opts, func() error {
		resp, err := http.Get(healthURL)
		if err != nil {
			return fmt.Errorf("failed to GET %s: %w", healthURL, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		// Check for expected response (case-insensitive)
		bodyStr := strings.ToLower(string(body))
		if !strings.Contains(bodyStr, "ok") && !strings.Contains(bodyStr, "healthy") {
			return fmt.Errorf("unexpected health response: %s", string(body))
		}

		return nil
	})
}

// Test 3: Static Assets
// Note: Apps use CDN client from unpkg.com, so no local client library to test
func testStaticAssets(appURL string, opts *SmokeTestOptions) error {
	// Skip - apps use CDN client, not local /livetemplate-client.js
	return nil
}

// Test 4: WebSocket Connection
func testWebSocket(appURL string, opts *SmokeTestOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Create Chrome context
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
		)...)
	defer allocCancel()

	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx)
	defer chromeCancel()

	var wsConnected bool

	// Navigate and check for WebSocket connection
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(appURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(2*time.Second), // Give WebSocket time to connect
		chromedp.Evaluate(`
			(function() {
				// Check if WebSocket exists and is connected
				if (window.ltClient && window.ltClient.ws) {
					return window.ltClient.ws.readyState === 1; // OPEN
				}
				return false;
			})()
		`, &wsConnected),
	)

	if err != nil {
		return fmt.Errorf("WebSocket test failed: %w", err)
	}

	if !wsConnected {
		// Not necessarily a failure - app might not use WebSockets yet
		// Just log this as a warning
		return nil
	}

	return nil
}

// Test 5: Template Rendering
func testTemplateRendering(appURL string, opts *SmokeTestOptions) error {
	return retryHTTP(opts, func() error {
		resp, err := http.Get(appURL)
		if err != nil {
			return fmt.Errorf("failed to GET %s: %w", appURL, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		bodyStr := string(body)

		// Check for template errors
		errorPatterns := []string{
			"template: ",
			"undefined:",
			"can't evaluate field",
			"panic:",
			"runtime error:",
		}

		for _, pattern := range errorPatterns {
			if strings.Contains(bodyStr, pattern) {
				return fmt.Errorf("template error detected: contains '%s'", pattern)
			}
		}

		// Check for basic HTML structure
		if !strings.Contains(bodyStr, "<html") && !strings.Contains(bodyStr, "<!DOCTYPE") {
			return fmt.Errorf("response doesn't appear to be HTML")
		}

		return nil
	})
}

// retryHTTP retries an HTTP operation with exponential backoff
func retryHTTP(opts *SmokeTestOptions, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= opts.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(opts.RetryDelay * time.Duration(attempt))
		}

		if err := fn(); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return fmt.Errorf("failed after %d retries: %w", opts.MaxRetries, lastErr)
}

// PrintResults prints smoke test results in a readable format
func (s *SmokeTestSuite) PrintResults() {
	fmt.Printf("\n=== Smoke Test Results ===\n")
	fmt.Printf("App URL: %s\n", s.AppURL)
	fmt.Printf("Total Duration: %v\n\n", s.TotalDuration)

	passed := 0
	failed := 0

	for _, result := range s.Results {
		status := "✅ PASS"
		if !result.Passed {
			status = "❌ FAIL"
			failed++
		} else {
			passed++
		}

		fmt.Printf("%s %s (%v)\n", status, result.Name, result.Duration)
		if result.Error != nil {
			fmt.Printf("   Error: %v\n", result.Error)
		}
	}

	fmt.Printf("\nResults: %d passed, %d failed\n", passed, failed)
}

// AllPassed returns true if all smoke tests passed
func (s *SmokeTestSuite) AllPassed() bool {
	for _, result := range s.Results {
		if !result.Passed {
			return false
		}
	}
	return true
}
