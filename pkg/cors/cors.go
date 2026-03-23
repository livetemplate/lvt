// Package cors provides CORS middleware for HTTP handlers.
package cors

import (
	"net/http"
	"strconv"
	"strings"
)

// Config holds CORS configuration.
type Config struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int // preflight cache duration in seconds
}

// Option configures CORS behavior.
type Option func(*Config)

// WithOrigins sets the allowed origins. Use "*" to allow all.
func WithOrigins(origins ...string) Option {
	return func(c *Config) { c.AllowOrigins = origins }
}

// WithMethods sets the allowed HTTP methods.
func WithMethods(methods ...string) Option {
	return func(c *Config) { c.AllowMethods = methods }
}

// WithHeaders sets the allowed request headers.
func WithHeaders(headers ...string) Option {
	return func(c *Config) { c.AllowHeaders = headers }
}

// WithCredentials enables credentials (cookies, auth headers).
func WithCredentials(allow bool) Option {
	return func(c *Config) { c.AllowCredentials = allow }
}

// WithMaxAge sets the preflight cache duration in seconds.
func WithMaxAge(seconds int) Option {
	return func(c *Config) { c.MaxAge = seconds }
}

// Middleware returns CORS middleware with the given options.
// Defaults: allow all origins, common methods, common headers, 1 hour max age.
func Middleware(opts ...Option) func(http.Handler) http.Handler {
	cfg := &Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:       3600,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	allowMethods := strings.Join(cfg.AllowMethods, ", ")
	allowHeaders := strings.Join(cfg.AllowHeaders, ", ")
	exposeHeaders := strings.Join(cfg.ExposeHeaders, ", ")
	maxAge := strconv.Itoa(cfg.MaxAge)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			if !cfg.isOriginAllowed(origin) {
				next.ServeHTTP(w, r)
				return
			}

			if len(cfg.AllowOrigins) == 1 && cfg.AllowOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if exposeHeaders != "" {
				w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
			}

			// Preflight
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", allowMethods)
				w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
				w.Header().Set("Access-Control-Max-Age", maxAge)
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (c *Config) isOriginAllowed(origin string) bool {
	for _, o := range c.AllowOrigins {
		if o == "*" || o == origin {
			return true
		}
	}
	return false
}
