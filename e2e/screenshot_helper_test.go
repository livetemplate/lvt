//go:build browser

// Package e2e provides screenshot capture utilities for the ui-polish skill.
// This file is meant to be copied to a test app's e2e directory and run there.
package e2e

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// ScreenshotManifest contains the paths to captured screenshots
type ScreenshotManifest struct {
	OutputDir      string            `json:"output_dir"`
	BaseURL        string            `json:"base_url"`
	Resource       string            `json:"resource"`
	Screenshots    map[string]string `json:"screenshots"`
	Errors         []string          `json:"errors,omitempty"`
	CaptureTime    time.Time         `json:"capture_time"`
	TotalCaptured  int               `json:"total_captured"`
	TotalRequested int               `json:"total_requested"`
}

// UIState represents a UI state to capture
type UIState struct {
	Name        string
	Description string
	Setup       func(ctx context.Context) error
}

var (
	outputDir   = flag.String("output", "/tmp/ui-polish-screenshots", "Output directory for screenshots")
	baseURL     = flag.String("url", "http://localhost:9999", "Base URL of the app")
	resourceArg = flag.String("resource", "posts", "Resource name to test")
)

// TestCaptureUIScreenshots captures screenshots of all UI states for the ui-polish skill
// Usage: go test -tags browser -run TestCaptureUIScreenshots -args -output /tmp/screenshots -url http://localhost:9999 -resource posts
func TestCaptureUIScreenshots(t *testing.T) {
	// Parse flags
	flag.Parse()

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Get Chrome context
	ctx, _, cleanup := GetPooledChrome(t)
	defer cleanup()

	// Extend timeout for screenshot capture
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	resource := *resourceArg
	resourceURL := fmt.Sprintf("%s/%s", *baseURL, resource)

	manifest := ScreenshotManifest{
		OutputDir:   *outputDir,
		BaseURL:     *baseURL,
		Resource:    resource,
		Screenshots: make(map[string]string),
		CaptureTime: time.Now(),
	}

	// Define UI states to capture
	states := []UIState{
		{
			Name:        "list_empty",
			Description: "Empty list view",
			Setup: func(ctx context.Context) error {
				return chromedp.Run(ctx,
					chromedp.Navigate(resourceURL),
					chromedp.WaitVisible(`body`, chromedp.ByQuery),
					chromedp.Sleep(500*time.Millisecond),
				)
			},
		},
		{
			Name:        "add_modal",
			Description: "Add modal open",
			Setup: func(ctx context.Context) error {
				// Look for common add button patterns
				return chromedp.Run(ctx,
					chromedp.Navigate(resourceURL),
					chromedp.WaitVisible(`body`, chromedp.ByQuery),
					chromedp.Sleep(300*time.Millisecond),
					// Try clicking add button using JavaScript
					chromedp.Evaluate(`
						(function() {
							const btn = document.querySelector('button[lf-action="showAdd"]') ||
							            document.querySelector('[data-action="add"]') ||
							            document.querySelector('.add-button') ||
							            document.querySelector('#add-button') ||
							            Array.from(document.querySelectorAll('button')).find(b => b.textContent.includes('Add'));
							if (btn) { btn.click(); return true; }
							return false;
						})()
					`, nil),
					chromedp.Sleep(500*time.Millisecond),
				)
			},
		},
		{
			Name:        "form_validation",
			Description: "Form with validation errors",
			Setup: func(ctx context.Context) error {
				return chromedp.Run(ctx,
					chromedp.Navigate(resourceURL),
					chromedp.WaitVisible(`body`, chromedp.ByQuery),
					chromedp.Sleep(300*time.Millisecond),
					// Open add modal
					chromedp.ActionFunc(func(ctx context.Context) error {
						return chromedp.Evaluate(`
							(function() {
								const btn = document.querySelector('button[lf-action="showAdd"]') ||
								            document.querySelector('[data-action="add"]') ||
								            Array.from(document.querySelectorAll('button')).find(b => b.textContent.includes('Add'));
								if (btn) { btn.click(); return true; }
								return false;
							})()
						`, nil).Do(ctx)
					}),
					chromedp.Sleep(300*time.Millisecond),
					// Try to submit empty form
					chromedp.ActionFunc(func(ctx context.Context) error {
						return chromedp.Evaluate(`
							(function() {
								const form = document.querySelector('form');
								if (form) {
									const submitBtn = form.querySelector('button[type="submit"]') ||
									                  form.querySelector('button:not([type="button"])');
									if (submitBtn) { submitBtn.click(); return true; }
								}
								return false;
							})()
						`, nil).Do(ctx)
					}),
					chromedp.Sleep(500*time.Millisecond),
				)
			},
		},
	}

	manifest.TotalRequested = len(states)

	// Capture each state
	for _, state := range states {
		t.Run(state.Name, func(t *testing.T) {
			t.Logf("Capturing state: %s (%s)", state.Name, state.Description)

			// Setup the state
			if err := state.Setup(ctx); err != nil {
				t.Logf("Warning: Setup for %s failed: %v", state.Name, err)
				manifest.Errors = append(manifest.Errors, fmt.Sprintf("%s: %v", state.Name, err))
				return
			}

			// Capture screenshot
			var screenshot []byte
			if err := chromedp.Run(ctx, chromedp.FullScreenshot(&screenshot, 90)); err != nil {
				t.Logf("Warning: Screenshot capture for %s failed: %v", state.Name, err)
				manifest.Errors = append(manifest.Errors, fmt.Sprintf("%s screenshot: %v", state.Name, err))
				return
			}

			// Save screenshot
			filename := fmt.Sprintf("%s.png", state.Name)
			filepath := filepath.Join(*outputDir, filename)
			if err := os.WriteFile(filepath, screenshot, 0644); err != nil {
				t.Logf("Warning: Failed to save screenshot %s: %v", state.Name, err)
				manifest.Errors = append(manifest.Errors, fmt.Sprintf("%s save: %v", state.Name, err))
				return
			}

			manifest.Screenshots[state.Name] = filepath
			manifest.TotalCaptured++
			t.Logf("Saved: %s (%d bytes)", filepath, len(screenshot))
		})
	}

	// Write manifest file
	manifestPath := filepath.Join(*outputDir, "manifest.json")
	manifestJSON, _ := json.MarshalIndent(manifest, "", "  ")
	if err := os.WriteFile(manifestPath, manifestJSON, 0644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	// Print manifest to stdout for the skill to capture
	fmt.Println("---MANIFEST_START---")
	fmt.Println(string(manifestJSON))
	fmt.Println("---MANIFEST_END---")

	t.Logf("Screenshot capture complete: %d/%d states captured", manifest.TotalCaptured, manifest.TotalRequested)
}

// TestCaptureResourceCRUD captures screenshots of full CRUD workflow
// This is a more comprehensive test that seeds data first
func TestCaptureResourceCRUD(t *testing.T) {
	flag.Parse()

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	ctx, _, cleanup := GetPooledChrome(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	resource := *resourceArg
	resourceURL := fmt.Sprintf("%s/%s", *baseURL, resource)

	manifest := ScreenshotManifest{
		OutputDir:   *outputDir,
		BaseURL:     *baseURL,
		Resource:    resource,
		Screenshots: make(map[string]string),
		CaptureTime: time.Now(),
	}

	// Helper to capture and save screenshot
	capture := func(name string) error {
		var screenshot []byte
		if err := chromedp.Run(ctx, chromedp.FullScreenshot(&screenshot, 90)); err != nil {
			return err
		}
		filename := fmt.Sprintf("%s.png", name)
		path := filepath.Join(*outputDir, filename)
		if err := os.WriteFile(path, screenshot, 0644); err != nil {
			return err
		}
		manifest.Screenshots[name] = path
		manifest.TotalCaptured++
		t.Logf("Captured: %s", path)
		return nil
	}

	// State 1: Navigate to list
	t.Log("State 1: Navigating to resource list")
	if err := chromedp.Run(ctx,
		chromedp.Navigate(resourceURL),
		chromedp.WaitVisible(`body`),
		chromedp.Sleep(500*time.Millisecond),
	); err != nil {
		t.Fatalf("Failed to navigate: %v", err)
	}
	capture("01_list_initial")
	manifest.TotalRequested++

	// State 2: Click Add button
	t.Log("State 2: Opening add modal")
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`
			(function() {
				const btn = document.querySelector('button[lf-action="showAdd"]') ||
				            Array.from(document.querySelectorAll('button')).find(b => b.textContent.includes('Add'));
				if (btn) { btn.click(); return true; }
				return false;
			})()
		`, nil),
		chromedp.Sleep(500*time.Millisecond),
	); err != nil {
		t.Logf("Warning: Could not open add modal: %v", err)
	}
	capture("02_add_modal")
	manifest.TotalRequested++

	// State 3: Fill form
	t.Log("State 3: Filling form")
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`
			(function() {
				const inputs = document.querySelectorAll('input[type="text"], textarea');
				inputs.forEach((input, i) => {
					if (input.name === 'title' || input.id === 'title') {
						input.value = 'Test Post Title';
					} else if (input.name === 'content' || input.id === 'content' || input.tagName === 'TEXTAREA') {
						input.value = 'This is test content for the post.';
					} else {
						input.value = 'Test Value ' + i;
					}
					input.dispatchEvent(new Event('input', { bubbles: true }));
				});
				return true;
			})()
		`, nil),
		chromedp.Sleep(300*time.Millisecond),
	); err != nil {
		t.Logf("Warning: Could not fill form: %v", err)
	}
	capture("03_form_filled")
	manifest.TotalRequested++

	// State 4: Submit form
	t.Log("State 4: Submitting form")
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`
			(function() {
				const form = document.querySelector('form');
				if (form) {
					const btn = form.querySelector('button[type="submit"]') ||
					            form.querySelector('button:not([type="button"])');
					if (btn) { btn.click(); return true; }
				}
				return false;
			})()
		`, nil),
		chromedp.Sleep(1*time.Second), // Wait for submission
	); err != nil {
		t.Logf("Warning: Could not submit form: %v", err)
	}
	capture("04_after_submit")
	manifest.TotalRequested++

	// State 5: Edit modal (if item was created)
	t.Log("State 5: Opening edit modal")
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`
			(function() {
				const editBtn = document.querySelector('button[lf-action="showEdit"]') ||
				                document.querySelector('[data-action="edit"]') ||
				                Array.from(document.querySelectorAll('button')).find(b => b.textContent.includes('Edit'));
				if (editBtn) { editBtn.click(); return true; }
				return false;
			})()
		`, nil),
		chromedp.Sleep(500*time.Millisecond),
	); err != nil {
		t.Logf("Warning: Could not open edit modal: %v", err)
	}
	capture("05_edit_modal")
	manifest.TotalRequested++

	// State 6: Delete confirmation
	t.Log("State 6: Delete confirmation")
	// Close edit modal first
	chromedp.Run(ctx,
		chromedp.Evaluate(`
			(function() {
				const closeBtn = document.querySelector('[lf-action="hideEdit"]') ||
				                 document.querySelector('.modal-close') ||
				                 document.querySelector('[data-dismiss="modal"]');
				if (closeBtn) { closeBtn.click(); }
				// Also try escape key
				document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
				return true;
			})()
		`, nil),
		chromedp.Sleep(300*time.Millisecond),
	)

	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`
			(function() {
				const delBtn = document.querySelector('button[lf-action="delete"]') ||
				               document.querySelector('[data-action="delete"]') ||
				               Array.from(document.querySelectorAll('button')).find(b => b.textContent.includes('Delete'));
				if (delBtn) { delBtn.click(); return true; }
				return false;
			})()
		`, nil),
		chromedp.Sleep(500*time.Millisecond),
	); err != nil {
		t.Logf("Warning: Could not trigger delete: %v", err)
	}
	capture("06_delete_confirm")
	manifest.TotalRequested++

	// Write manifest
	manifestPath := filepath.Join(*outputDir, "manifest.json")
	manifestJSON, _ := json.MarshalIndent(manifest, "", "  ")
	os.WriteFile(manifestPath, manifestJSON, 0644)

	fmt.Println("---MANIFEST_START---")
	fmt.Println(string(manifestJSON))
	fmt.Println("---MANIFEST_END---")

	t.Logf("CRUD workflow capture complete: %d/%d states", manifest.TotalCaptured, manifest.TotalRequested)
}
