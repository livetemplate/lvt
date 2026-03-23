package cache

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPCache(t *testing.T) {
	handler := HTTPCache(3600)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	cc := rec.Header().Get("Cache-Control")
	if cc != "public, max-age=3600" {
		t.Errorf("Cache-Control = %q, want %q", cc, "public, max-age=3600")
	}
}

func TestNoCache(t *testing.T) {
	handler := NoCache()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	cc := rec.Header().Get("Cache-Control")
	if cc != "no-cache, no-store, must-revalidate" {
		t.Errorf("Cache-Control = %q, want %q", cc, "no-cache, no-store, must-revalidate")
	}

	pragma := rec.Header().Get("Pragma")
	if pragma != "no-cache" {
		t.Errorf("Pragma = %q, want %q", pragma, "no-cache")
	}

	expires := rec.Header().Get("Expires")
	if expires != "0" {
		t.Errorf("Expires = %q, want %q", expires, "0")
	}
}
