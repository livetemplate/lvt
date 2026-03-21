package authz

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/livetemplate/lvt/pkg/cookie"
)

// TokenLookup is the interface for looking up a user ID from a session token.
// The generated sqlc Queries struct satisfies this via its GetUserToken method.
type TokenLookup interface {
	GetUserToken(ctx context.Context, arg interface{ GetToken() string; GetNow() time.Time }) (interface{ GetUserID() string }, error)
}

// CookieAuthenticator implements livetemplate.Authenticator by reading
// the auth session cookie to identify the user during WebSocket setup.
// This enables ctx.UserID() to return the real authenticated user ID
// in LiveTemplate controller actions.
type CookieAuthenticator struct {
	cookieName string
	lookupFn   func(ctx context.Context, token string) (userID string, err error)
}

// NewCookieAuthenticator creates an authenticator that reads the auth session
// cookie and looks up the user ID. The lookupFn should query the tokens table
// to resolve a session token to a user ID.
//
// Example:
//
//	authz.NewCookieAuthenticator("users_token", func(ctx context.Context, token string) (string, error) {
//	    row, err := queries.GetUserToken(ctx, models.GetUserTokenParams{Token: token, Now: time.Now()})
//	    if err != nil { return "", err }
//	    return row.UserID, nil
//	})
func NewCookieAuthenticator(cookieName string, lookupFn func(ctx context.Context, token string) (string, error)) *CookieAuthenticator {
	return &CookieAuthenticator{
		cookieName: cookieName,
		lookupFn:   lookupFn,
	}
}

// Identify returns the user ID from the auth session cookie.
// Returns "" for unauthenticated requests (no cookie or invalid token).
func (a *CookieAuthenticator) Identify(r *http.Request) (string, error) {
	token := cookie.Get(r, a.cookieName)
	if token == "" {
		return "", nil
	}
	userID, err := a.lookupFn(r.Context(), token)
	if err != nil {
		return "", nil // Invalid/expired token — treat as unauthenticated
	}
	return userID, nil
}

// GetSessionGroup returns the user ID as the session group so that
// all tabs for the same authenticated user share LiveTemplate state.
// Falls back to a browser-based cookie for unauthenticated users.
func (a *CookieAuthenticator) GetSessionGroup(r *http.Request, userID string) (string, error) {
	if userID != "" {
		return userID, nil
	}
	// Fallback for unauthenticated: use browser session cookie
	c, err := r.Cookie("livetemplate-id")
	if err == nil && c.Value != "" {
		return c.Value, nil
	}
	return fmt.Sprintf("anon-%d", time.Now().UnixNano()), nil
}
