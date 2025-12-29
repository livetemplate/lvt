// Package token provides secure random token generation for authentication.
package token

import (
	"crypto/rand"
	"encoding/base64"
)

// Token context constants for different use cases.
// These are stored in the database to identify the token's purpose.
const (
	// ContextSession is for session tokens (login)
	ContextSession = "session"
	// ContextMagic is for magic link tokens (passwordless login)
	ContextMagic = "magic"
	// ContextReset is for password reset tokens
	ContextReset = "reset"
	// ContextConfirm is for email confirmation tokens
	ContextConfirm = "confirm"
)

// Generate creates a cryptographically secure random token.
// Returns a 32-byte random value encoded as base64 URL-safe string.
//
// Example:
//
//	tok, err := token.Generate()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use tok as session token, magic link token, etc.
func Generate() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
