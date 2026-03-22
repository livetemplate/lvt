package authz

import (
	"net/http"
	"strings"
)

// RequireRole returns middleware that checks if the authenticated user
// has one of the allowed roles. The getUserRole function extracts the
// user's role from the request (typically reading from context set by
// auth.RequireAuth middleware).
//
// Returns 403 Forbidden if the user's role is not in the allowed list.
// Must be used after authentication middleware (user must be identified).
func RequireRole(getUserRole func(r *http.Request) string, roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := getUserRole(r)
			if role == "" || !allowed[role] {
				ServeForbidden(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ServeForbidden writes a 403 Forbidden response.
// Returns JSON for API requests (Accept: application/json),
// HTML for browser requests.
func ServeForbidden(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error":"forbidden","message":"You do not have permission to perform this action"}`))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`<!DOCTYPE html>
<html><head><title>403 Forbidden</title></head>
<body style="font-family: system-ui; max-width: 600px; margin: 4rem auto; text-align: center;">
<h1>403 Forbidden</h1>
<p>You do not have permission to perform this action.</p>
<a href="/">Go Home</a>
</body></html>`))
}
