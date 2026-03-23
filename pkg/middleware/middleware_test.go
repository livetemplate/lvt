package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain_Order(t *testing.T) {
	var order []string

	mw := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+"-before")
				next.ServeHTTP(w, r)
				order = append(order, name+"-after")
			})
		}
	}

	chain := Chain(mw("first"), mw("second"), mw("third"))
	handler := chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	expected := []string{
		"first-before", "second-before", "third-before",
		"handler",
		"third-after", "second-after", "first-after",
	}

	if len(order) != len(expected) {
		t.Fatalf("got %d calls, want %d: %v", len(order), len(expected), order)
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d] = %q, want %q", i, order[i], v)
		}
	}
}

func TestChain_Empty(t *testing.T) {
	chain := Chain()
	called := false
	handler := chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("handler should be called with empty chain")
	}
}

func TestChain_Single(t *testing.T) {
	headerSet := false
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "true")
			headerSet = true
			next.ServeHTTP(w, r)
		})
	}

	chain := Chain(mw)
	handler := chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !headerSet {
		t.Error("middleware should have been called")
	}
	if rec.Header().Get("X-Test") != "true" {
		t.Error("middleware should have set header")
	}
}

func TestChain_Compose(t *testing.T) {
	addHeader := func(key, value string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(key, value)
				next.ServeHTTP(w, r)
			})
		}
	}

	// Compose two chains
	base := Chain(addHeader("X-Base", "true"))
	extended := Chain(base, addHeader("X-Extra", "true"))

	handler := extended(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-Base") != "true" {
		t.Error("base middleware should set X-Base")
	}
	if rec.Header().Get("X-Extra") != "true" {
		t.Error("extended middleware should set X-Extra")
	}
}
