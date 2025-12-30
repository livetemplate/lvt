// Package cookie provides utilities for setting and clearing HTTP cookies
// with secure defaults appropriate for authentication and session management.
package cookie

import (
	"net/http"
	"strings"
)

// Set sets a cookie with secure defaults.
// Uses HttpOnly and SameSite=Lax for security.
//
// Example:
//
//	cookie.Set(w, "preference", "dark", 86400) // 1 day
func Set(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// SetSecure sets a cookie with strict security settings.
// Uses HttpOnly, Secure, and SameSite=Strict.
// Use this for sensitive cookies like session tokens.
//
// Example:
//
//	cookie.SetSecure(w, "session", token, 30*24*60*60) // 30 days
func SetSecure(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// Clear clears a cookie by setting MaxAge to -1.
//
// Example:
//
//	cookie.Clear(w, "session")
func Clear(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSecure clears a secure cookie by setting MaxAge to -1.
// Matches the attributes used by SetSecure.
//
// Example:
//
//	cookie.ClearSecure(w, "session")
func ClearSecure(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearLiveTemplateSession clears the LiveTemplate session cookie.
// This forces a fresh session state on the next page load, which is
// necessary after logout to ensure the home page shows the correct
// logged-out state.
//
// Example:
//
//	// In logout handler
//	cookie.Clear(w, "users_token")
//	cookie.ClearLiveTemplateSession(w)
//	http.Redirect(w, r, "/", http.StatusSeeOther)
func ClearLiveTemplateSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "livetemplate-id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// SetSession sets a session cookie with strict security settings.
// The maxAgeDays parameter specifies how long the session should last.
// Uses HttpOnly and SameSite=Strict for authentication cookies.
// The Secure flag is automatically set based on the request protocol.
//
// Example:
//
//	cookie.SetSession(w, r, "users_token", token, 30) // 30 days
func SetSession(w http.ResponseWriter, r *http.Request, name, value string, maxAgeDays int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAgeDays * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   IsSecure(r),
		SameSite: http.SameSiteStrictMode,
	})
}

// IsSecure returns true if the request is over HTTPS.
// It checks both the URL scheme and common proxy headers.
func IsSecure(r *http.Request) bool {
	// Check direct TLS connection
	if r.TLS != nil {
		return true
	}

	// Check URL scheme
	if r.URL.Scheme == "https" {
		return true
	}

	// Check common proxy headers
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}

	// Check if behind a load balancer
	if strings.HasPrefix(r.Header.Get("X-Forwarded-Ssl"), "on") {
		return true
	}

	return false
}

// Get retrieves a cookie value by name.
// Returns empty string if cookie doesn't exist.
//
// Example:
//
//	token := cookie.Get(r, "session")
//	if token == "" {
//	    // Not logged in
//	}
func Get(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}
