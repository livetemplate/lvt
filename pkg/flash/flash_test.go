package flash

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSet(t *testing.T) {
	w := httptest.NewRecorder()
	Set(w, "info", "Test message")

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	c := cookies[0]
	if c.Name != "flash_info" {
		t.Errorf("expected name 'flash_info', got '%s'", c.Name)
	}
	if c.Value != url.QueryEscape("Test message") {
		t.Errorf("expected URL-encoded value, got '%s'", c.Value)
	}
	if c.MaxAge != DefaultMaxAge {
		t.Errorf("expected MaxAge %d, got %d", DefaultMaxAge, c.MaxAge)
	}
}

func TestGet(t *testing.T) {
	// Test with flash present
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "flash_error",
		Value: url.QueryEscape("Something went wrong"),
	})

	w := httptest.NewRecorder()
	msg := Get(req, w, "error")

	if msg != "Something went wrong" {
		t.Errorf("expected 'Something went wrong', got '%s'", msg)
	}

	// Verify cookie was cleared
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie (clear), got %d", len(cookies))
	}
	if cookies[0].MaxAge != -1 {
		t.Errorf("expected MaxAge -1 (clear), got %d", cookies[0].MaxAge)
	}
}

func TestGetEmpty(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	msg := Get(req, w, "nonexistent")
	if msg != "" {
		t.Errorf("expected empty string, got '%s'", msg)
	}
}

func TestErrorAndSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, "Error message")
	Success(w, "Success message")

	cookies := w.Result().Cookies()
	if len(cookies) != 2 {
		t.Fatalf("expected 2 cookies, got %d", len(cookies))
	}

	// Find error cookie
	var errorCookie, successCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "flash_error" {
			errorCookie = c
		}
		if c.Name == "flash_success" {
			successCookie = c
		}
	}

	if errorCookie == nil {
		t.Error("expected flash_error cookie")
	}
	if successCookie == nil {
		t.Error("expected flash_success cookie")
	}
}

func TestGetErrorAndSuccess(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "flash_error",
		Value: url.QueryEscape("Error!"),
	})
	req.AddCookie(&http.Cookie{
		Name:  "flash_success",
		Value: url.QueryEscape("Success!"),
	})

	w := httptest.NewRecorder()

	errMsg := GetError(req, w)
	if errMsg != "Error!" {
		t.Errorf("expected 'Error!', got '%s'", errMsg)
	}

	successMsg := GetSuccess(req, w)
	if successMsg != "Success!" {
		t.Errorf("expected 'Success!', got '%s'", successMsg)
	}
}

func TestPending(t *testing.T) {
	// Test SetPending
	w := httptest.NewRecorder()
	SetPending(w)

	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == PendingKey && c.Value == "1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected flash_pending cookie to be set")
	}

	// Test IsPending
	req := httptest.NewRequest("GET", "/", nil)
	if IsPending(req) {
		t.Error("expected IsPending to be false without cookie")
	}

	req.AddCookie(&http.Cookie{Name: PendingKey, Value: "1"})
	if !IsPending(req) {
		t.Error("expected IsPending to be true with cookie")
	}
}

func TestRedirectWithError(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	RedirectWithError(w, req, "/auth", "Login failed")

	// Check redirect
	if w.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/auth" {
		t.Errorf("expected redirect to '/auth', got '%s'", loc)
	}

	// Check cookies
	cookies := w.Result().Cookies()
	hasError := false
	hasPending := false
	for _, c := range cookies {
		if c.Name == "flash_error" {
			hasError = true
		}
		if c.Name == PendingKey {
			hasPending = true
		}
	}
	if !hasError {
		t.Error("expected flash_error cookie")
	}
	if !hasPending {
		t.Error("expected flash_pending cookie")
	}
}

func TestGetAll(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "flash_error",
		Value: url.QueryEscape("Error message"),
	})
	req.AddCookie(&http.Cookie{
		Name:  "flash_success",
		Value: url.QueryEscape("Success message"),
	})

	w := httptest.NewRecorder()
	msgs := GetAll(req, w)

	if msgs.Error != "Error message" {
		t.Errorf("expected 'Error message', got '%s'", msgs.Error)
	}
	if msgs.Success != "Success message" {
		t.Errorf("expected 'Success message', got '%s'", msgs.Success)
	}
}

func TestSpecialCharacters(t *testing.T) {
	w := httptest.NewRecorder()
	Set(w, "test", "Hello & goodbye <script>")

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	w2 := httptest.NewRecorder()
	msg := Get(req, w2, "test")

	if msg != "Hello & goodbye <script>" {
		t.Errorf("expected special characters preserved, got '%s'", msg)
	}
}
