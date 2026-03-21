package authz

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireRole_Allowed(t *testing.T) {
	getRoleFromReq := func(r *http.Request) string { return "admin" }

	handler := RequireRole(getRoleFromReq, "admin", "editor")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRequireRole_Denied(t *testing.T) {
	getRoleFromReq := func(r *http.Request) string { return "user" }

	handler := RequireRole(getRoleFromReq, "admin")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestRequireRole_EmptyRole(t *testing.T) {
	getRoleFromReq := func(r *http.Request) string { return "" }

	handler := RequireRole(getRoleFromReq, "admin")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 for empty role, got %d", rec.Code)
	}
}

func TestServeForbidden_HTML(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()

	ServeForbidden(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected HTML content type, got %q", ct)
	}
}

func TestServeForbidden_JSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()

	ServeForbidden(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected JSON content type, got %q", ct)
	}
	body := rec.Body.String()
	if body == "" || body[0] != '{' {
		t.Errorf("expected JSON body, got %q", body)
	}
}
