// Package security provides HTTP security middleware and utilities.
package security

import (
	"net/http"
	"net/url"
	"strings"
)

// ValidateOrigin checks if the request's Origin or Referer header
// matches the expected host. This provides CSRF protection for
// form submissions.
//
// Returns true if the request origin is valid, false otherwise.
//
// The validation:
// - Parses the Origin header (or Referer as fallback) as a URL
// - Compares the parsed host against the request's Host header
// - Allows localhost/127.0.0.1 for development
//
// Example:
//
//	func HandleLogin(w http.ResponseWriter, r *http.Request) {
//	    if !security.ValidateOrigin(r) {
//	        http.Error(w, "Invalid request origin", http.StatusForbidden)
//	        return
//	    }
//	    // Process login...
//	}
func ValidateOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = r.Header.Get("Referer")
	}

	// If neither header is present, reject the request
	// This prevents CSRF attacks that strip these headers
	if origin == "" {
		return false
	}

	originURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	// Get the expected host from the request
	expectedHost := r.Host

	// Strip port from expected host for comparison if needed
	expectedHostNoPort := stripPort(expectedHost)
	originHostNoPort := stripPort(originURL.Host)

	// Strict host comparison
	if originHostNoPort == expectedHostNoPort {
		return true
	}

	// Allow localhost variants for development
	if isLocalhost(originHostNoPort) && isLocalhost(expectedHostNoPort) {
		return true
	}

	return false
}

// ValidateOriginAllowEmpty is like ValidateOrigin but allows requests
// without Origin/Referer headers. This is useful for endpoints that
// need to accept requests from non-browser clients while still
// protecting against CSRF from browsers.
//
// Example:
//
//	func HandleAPILogin(w http.ResponseWriter, r *http.Request) {
//	    if !security.ValidateOriginAllowEmpty(r) {
//	        http.Error(w, "Invalid request origin", http.StatusForbidden)
//	        return
//	    }
//	    // Process login...
//	}
func ValidateOriginAllowEmpty(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	referer := r.Header.Get("Referer")

	// Allow if neither header is present (e.g., curl, API clients)
	if origin == "" && referer == "" {
		return true
	}

	return ValidateOrigin(r)
}

// stripPort removes the port from a host string.
func stripPort(host string) string {
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		// Check if it's IPv6 (has [ ])
		if strings.Contains(host, "]") {
			// IPv6 with port: [::1]:8080
			if bracketIdx := strings.LastIndex(host, "]"); bracketIdx < idx {
				return host[:idx]
			}
			return host
		}
		return host[:idx]
	}
	return host
}

// isLocalhost checks if the host is a localhost variant.
func isLocalhost(host string) bool {
	host = strings.ToLower(host)
	return host == "localhost" ||
		host == "127.0.0.1" ||
		host == "::1" ||
		host == "[::1]"
}
