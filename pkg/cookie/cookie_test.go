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
	t.Run("HTTP request - not secure", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "http://localhost/login", nil)
		SetSession(w, r, "users_token", "abc123", 30)

		cookies := w.Result().Cookies()
		if len(cookies) != 1 {
			t.Fatalf("expected 1 cookie, got %d", len(cookies))
		}

		c := cookies[0]
		expectedMaxAge := 30 * 24 * 60 * 60
		if c.MaxAge != expectedMaxAge {
			t.Errorf("expected MaxAge %d, got %d", expectedMaxAge, c.MaxAge)
		}
		if c.Secure {
			t.Error("expected Secure to be false for HTTP")
		}
		if c.SameSite != http.SameSiteStrictMode {
			t.Errorf("expected SameSite %v, got %v", http.SameSiteStrictMode, c.SameSite)
		}
	})

	t.Run("HTTPS request via X-Forwarded-Proto", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "http://example.com/login", nil)
		r.Header.Set("X-Forwarded-Proto", "https")
		SetSession(w, r, "users_token", "abc123", 30)

		cookies := w.Result().Cookies()
		if len(cookies) != 1 {
			t.Fatalf("expected 1 cookie, got %d", len(cookies))
		}

		c := cookies[0]
		if !c.Secure {
			t.Error("expected Secure to be true for HTTPS")
		}
	})
}

func TestIsSecure(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*http.Request)
		expected bool
	}{
		{
			name:     "plain HTTP",
			setup:    func(r *http.Request) {},
			expected: false,
		},
		{
			name: "X-Forwarded-Proto https",
			setup: func(r *http.Request) {
				r.Header.Set("X-Forwarded-Proto", "https")
			},
			expected: true,
		},
		{
			name: "X-Forwarded-Ssl on",
			setup: func(r *http.Request) {
				r.Header.Set("X-Forwarded-Ssl", "on")
			},
			expected: true,
		},
		{
			name: "https URL scheme",
			setup: func(r *http.Request) {
				r.URL.Scheme = "https"
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			tt.setup(r)
			if got := IsSecure(r); got != tt.expected {
				t.Errorf("IsSecure() = %v, want %v", got, tt.expected)
			}
		})
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
