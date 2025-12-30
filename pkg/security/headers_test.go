package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHeadersDefaults(t *testing.T) {
	handler := Headers()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	tests := []struct {
		header string
		want   string
	}{
		{"X-Frame-Options", "DENY"},
		{"X-Content-Type-Options", "nosniff"},
		{"X-XSS-Protection", "1; mode=block"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
	}

	for _, tt := range tests {
		got := rec.Header().Get(tt.header)
		if got != tt.want {
			t.Errorf("%s = %q, want %q", tt.header, got, tt.want)
		}
	}
}

func TestHeadersWithCSP(t *testing.T) {
	csp := "default-src 'self'; script-src 'self' 'unsafe-inline'"
	handler := Headers(WithCSP(csp))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("Content-Security-Policy")
	if got != csp {
		t.Errorf("Content-Security-Policy = %q, want %q", got, csp)
	}
}

func TestHeadersWithHSTS(t *testing.T) {
	handler := Headers(WithHSTS(365*24*time.Hour, true))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("Strict-Transport-Security")
	want := "max-age=31536000; includeSubDomains"
	if got != want {
		t.Errorf("Strict-Transport-Security = %q, want %q", got, want)
	}
}

func TestHeadersWithHSTSNoSubDomains(t *testing.T) {
	handler := Headers(WithHSTS(365*24*time.Hour, false))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("Strict-Transport-Security")
	want := "max-age=31536000"
	if got != want {
		t.Errorf("Strict-Transport-Security = %q, want %q", got, want)
	}
}

func TestHeadersWithFrameOptions(t *testing.T) {
	handler := Headers(WithFrameOptions("SAMEORIGIN"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("X-Frame-Options")
	if got != "SAMEORIGIN" {
		t.Errorf("X-Frame-Options = %q, want SAMEORIGIN", got)
	}
}

func TestHeadersWithReferrerPolicy(t *testing.T) {
	handler := Headers(WithReferrerPolicy("no-referrer"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("Referrer-Policy")
	if got != "no-referrer" {
		t.Errorf("Referrer-Policy = %q, want no-referrer", got)
	}
}

func TestHeadersMultipleOptions(t *testing.T) {
	handler := Headers(
		WithCSP("default-src 'self'"),
		WithHSTS(time.Hour, true),
		WithFrameOptions("SAMEORIGIN"),
		WithReferrerPolicy("no-referrer"),
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	tests := []struct {
		header string
		want   string
	}{
		{"Content-Security-Policy", "default-src 'self'"},
		{"Strict-Transport-Security", "max-age=3600; includeSubDomains"},
		{"X-Frame-Options", "SAMEORIGIN"},
		{"Referrer-Policy", "no-referrer"},
	}

	for _, tt := range tests {
		got := rec.Header().Get(tt.header)
		if got != tt.want {
			t.Errorf("%s = %q, want %q", tt.header, got, tt.want)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.FrameOptions != "DENY" {
		t.Errorf("FrameOptions = %q, want DENY", config.FrameOptions)
	}

	if !config.ContentTypeNoSniff {
		t.Error("ContentTypeNoSniff should be true by default")
	}

	if config.XSSProtection != "1; mode=block" {
		t.Errorf("XSSProtection = %q, want '1; mode=block'", config.XSSProtection)
	}

	if config.ReferrerPolicy != "strict-origin-when-cross-origin" {
		t.Errorf("ReferrerPolicy = %q, want strict-origin-when-cross-origin", config.ReferrerPolicy)
	}
}
