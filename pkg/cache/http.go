package cache

import (
	"fmt"
	"net/http"
)

// HTTPCache returns middleware that sets Cache-Control headers.
// maxAge is the duration in seconds that the response may be cached.
func HTTPCache(maxAge int) func(http.Handler) http.Handler {
	value := fmt.Sprintf("public, max-age=%d", maxAge)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", value)
			next.ServeHTTP(w, r)
		})
	}
}

// NoCache returns middleware that sets headers to prevent caching.
func NoCache() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			next.ServeHTTP(w, r)
		})
	}
}
