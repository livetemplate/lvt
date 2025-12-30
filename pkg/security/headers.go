// Package security provides HTTP security middleware.
package security

import (
	"net/http"
	"time"
)

// HeadersConfig holds configuration for security headers.
type HeadersConfig struct {
	// CSP is the Content-Security-Policy header value.
	CSP string

	// HSTSMaxAge sets the Strict-Transport-Security max-age directive.
	// Set to 0 to disable HSTS.
	HSTSMaxAge time.Duration

	// HSTSIncludeSubDomains includes subdomains in HSTS.
	HSTSIncludeSubDomains bool

	// FrameOptions sets X-Frame-Options (e.g., "DENY", "SAMEORIGIN").
	FrameOptions string

	// ContentTypeNoSniff enables X-Content-Type-Options: nosniff.
	ContentTypeNoSniff bool

	// XSSProtection sets X-XSS-Protection header.
	XSSProtection string

	// ReferrerPolicy sets the Referrer-Policy header.
	ReferrerPolicy string
}

// Option configures HeadersConfig.
type Option func(*HeadersConfig)

// WithCSP sets the Content-Security-Policy header.
//
// Example:
//
//	security.Headers(security.WithCSP("default-src 'self'"))
func WithCSP(csp string) Option {
	return func(c *HeadersConfig) {
		c.CSP = csp
	}
}

// WithHSTS enables HTTP Strict Transport Security.
//
// Example:
//
//	security.Headers(security.WithHSTS(365 * 24 * time.Hour, true))
func WithHSTS(maxAge time.Duration, includeSubDomains bool) Option {
	return func(c *HeadersConfig) {
		c.HSTSMaxAge = maxAge
		c.HSTSIncludeSubDomains = includeSubDomains
	}
}

// WithFrameOptions sets X-Frame-Options header.
//
// Example:
//
//	security.Headers(security.WithFrameOptions("DENY"))
func WithFrameOptions(value string) Option {
	return func(c *HeadersConfig) {
		c.FrameOptions = value
	}
}

// WithContentTypeNoSniff enables X-Content-Type-Options: nosniff.
func WithContentTypeNoSniff() Option {
	return func(c *HeadersConfig) {
		c.ContentTypeNoSniff = true
	}
}

// WithXSSProtection sets X-XSS-Protection header.
//
// Example:
//
//	security.Headers(security.WithXSSProtection("1; mode=block"))
func WithXSSProtection(value string) Option {
	return func(c *HeadersConfig) {
		c.XSSProtection = value
	}
}

// WithReferrerPolicy sets Referrer-Policy header.
//
// Example:
//
//	security.Headers(security.WithReferrerPolicy("strict-origin-when-cross-origin"))
func WithReferrerPolicy(policy string) Option {
	return func(c *HeadersConfig) {
		c.ReferrerPolicy = policy
	}
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() HeadersConfig {
	return HeadersConfig{
		FrameOptions:       "DENY",
		ContentTypeNoSniff: true,
		XSSProtection:      "1; mode=block",
		ReferrerPolicy:     "strict-origin-when-cross-origin",
	}
}

// Headers returns middleware that adds security headers to responses.
//
// Example:
//
//	mux := http.NewServeMux()
//	handler := security.Headers(
//	    security.WithCSP("default-src 'self'"),
//	    security.WithHSTS(365*24*time.Hour, true),
//	)(mux)
func Headers(opts ...Option) func(http.Handler) http.Handler {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(&config)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()

			if config.CSP != "" {
				h.Set("Content-Security-Policy", config.CSP)
			}

			if config.HSTSMaxAge > 0 {
				value := "max-age=" + formatSeconds(config.HSTSMaxAge)
				if config.HSTSIncludeSubDomains {
					value += "; includeSubDomains"
				}
				h.Set("Strict-Transport-Security", value)
			}

			if config.FrameOptions != "" {
				h.Set("X-Frame-Options", config.FrameOptions)
			}

			if config.ContentTypeNoSniff {
				h.Set("X-Content-Type-Options", "nosniff")
			}

			if config.XSSProtection != "" {
				h.Set("X-XSS-Protection", config.XSSProtection)
			}

			if config.ReferrerPolicy != "" {
				h.Set("Referrer-Policy", config.ReferrerPolicy)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func formatSeconds(d time.Duration) string {
	seconds := int64(d.Seconds())
	return formatInt(seconds)
}

func formatInt(n int64) string {
	if n == 0 {
		return "0"
	}

	// Handle negative numbers
	negative := n < 0
	if negative {
		n = -n
	}

	// Build digits in reverse
	var digits [20]byte
	i := len(digits)
	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}

	if negative {
		i--
		digits[i] = '-'
	}

	return string(digits[i:])
}
