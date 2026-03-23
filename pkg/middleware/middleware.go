// Package middleware provides composable middleware helpers for HTTP handlers.
// It enables named middleware groups for different route scopes (web, api, admin).
package middleware

import "net/http"

// Middleware is the standard HTTP middleware signature.
type Middleware = func(http.Handler) http.Handler

// Chain composes multiple middlewares into a single middleware.
// Middlewares are applied in order: the first argument is the outermost layer.
//
// Example:
//
//	global := middleware.Chain(rateLimiter, securityHeaders, recovery, logging)
//	handler := global(mux) // rateLimiter(securityHeaders(recovery(logging(mux))))
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
