package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSet(t *testing.T) {
	w := httptest.NewRecorder()
	Set(w, "test", "value", 3600)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	c := cookies[0]
	if c.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", c.Name)
	}
	if c.Value != "value" {
		t.Errorf("expected value 'value', got '%s'", c.Value)
	}
	if c.MaxAge != 3600 {
		t.Errorf("expected MaxAge 3600, got %d", c.MaxAge)
	}
	if !c.HttpOnly {
		t.Error("expected HttpOnly to be true")
	}
	if c.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite Lax, got %v", c.SameSite)
	}
}

func TestSetSecure(t *testing.T) {
	w := httptest.NewRecorder()
	SetSecure(w, "session", "token123", 86400)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	c := cookies[0]
	if c.Name != "session" {
		t.Errorf("expected name 'session', got '%s'", c.Name)
	}
	if !c.HttpOnly {
		t.Error("expected HttpOnly to be true")
	}
	if !c.Secure {
		t.Error("expected Secure to be true")
	}
	if c.SameSite != http.SameSiteStrictMode {
		t.Errorf("expected SameSite Strict, got %v", c.SameSite)
	}
}

func TestClear(t *testing.T) {
	w := httptest.NewRecorder()
	Clear(w, "test")

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	c := cookies[0]
	if c.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", c.Name)
	}
	if c.Value != "" {
		t.Errorf("expected empty value, got '%s'", c.Value)
	}
	if c.MaxAge != -1 {
		t.Errorf("expected MaxAge -1, got %d", c.MaxAge)
	}
}

func TestClearLiveTemplateSession(t *testing.T) {
	w := httptest.NewRecorder()
	ClearLiveTemplateSession(w)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	c := cookies[0]
	if c.Name != "livetemplate-id" {
		t.Errorf("expected name 'livetemplate-id', got '%s'", c.Name)
	}
	if c.MaxAge != -1 {
		t.Errorf("expected MaxAge -1, got %d", c.MaxAge)
	}
}

func TestSetSession(t *testing.T) {
	w := httptest.NewRecorder()
	SetSession(w, "users_token", "abc123", 30)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	c := cookies[0]
	expectedMaxAge := 30 * 24 * 60 * 60
	if c.MaxAge != expectedMaxAge {
		t.Errorf("expected MaxAge %d, got %d", expectedMaxAge, c.MaxAge)
	}
}

func TestGet(t *testing.T) {
	// Test with cookie present
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "test", Value: "myvalue"})

	val := Get(req, "test")
	if val != "myvalue" {
		t.Errorf("expected 'myvalue', got '%s'", val)
	}

	// Test with cookie not present
	req2 := httptest.NewRequest("GET", "/", nil)
	val2 := Get(req2, "nonexistent")
	if val2 != "" {
		t.Errorf("expected empty string, got '%s'", val2)
	}
}
