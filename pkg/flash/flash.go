// Package flash provides utilities for one-time flash messages using cookies.
// Flash messages are commonly used to show feedback after form submissions
// or redirects (e.g., "Registration successful!" or "Invalid credentials").
package flash

import (
	"log"
	"net/http"
	"net/url"
)

const (
	// ErrorKey is the cookie name for error flash messages
	ErrorKey = "flash_error"
	// SuccessKey is the cookie name for success flash messages
	SuccessKey = "flash_success"
	// PendingKey marks that a flash message was just set (for one-time display)
	PendingKey = "flash_pending"

	// DefaultMaxAge is how long flash cookies live (30 seconds)
	// This allows time for slow redirects or network delays
	DefaultMaxAge = 30
)

// Set sets a flash message cookie with the given key.
// The message is URL-encoded to handle special characters.
//
// Example:
//
//	flash.Set(w, "info", "Your session will expire soon")
func Set(w http.ResponseWriter, key, message string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "flash_" + key,
		Value:    url.QueryEscape(message),
		Path:     "/",
		MaxAge:   DefaultMaxAge,
		HttpOnly: true,
	})
}

// Get reads and clears a flash message.
// Returns empty string if no flash message exists.
//
// Example:
//
//	if msg := flash.Get(r, w, "error"); msg != "" {
//	    // Display error message
//	}
func Get(r *http.Request, w http.ResponseWriter, key string) string {
	cookieName := "flash_" + key
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value == "" {
		return ""
	}

	// Clear the cookie after reading
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	msg, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		log.Printf("flash: failed to decode cookie value: %v", err)
		return ""
	}
	return msg
}

// Error sets an error flash message.
//
// Example:
//
//	flash.Error(w, "Invalid email or password")
func Error(w http.ResponseWriter, message string) {
	Set(w, "error", message)
}

// Success sets a success flash message.
//
// Example:
//
//	flash.Success(w, "Account created successfully!")
func Success(w http.ResponseWriter, message string) {
	Set(w, "success", message)
}

// GetError reads and clears the error flash message.
//
// Example:
//
//	if errMsg := flash.GetError(r, w); errMsg != "" {
//	    state.FlashError = errMsg
//	}
func GetError(r *http.Request, w http.ResponseWriter) string {
	return Get(r, w, "error")
}

// GetSuccess reads and clears the success flash message.
//
// Example:
//
//	if successMsg := flash.GetSuccess(r, w); successMsg != "" {
//	    state.FlashSuccess = successMsg
//	}
func GetSuccess(r *http.Request, w http.ResponseWriter) string {
	return Get(r, w, "success")
}

// SetPending marks that a flash message was just set.
// This is used to distinguish between a fresh flash and a page reload.
func SetPending(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     PendingKey,
		Value:    "1",
		Path:     "/",
		MaxAge:   DefaultMaxAge,
		HttpOnly: true,
	})
}

// ClearPending clears the pending marker.
func ClearPending(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   PendingKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

// IsPending checks if a flash message was just set (vs page reload).
func IsPending(r *http.Request) bool {
	_, err := r.Cookie(PendingKey)
	return err == nil
}

// RedirectWithError redirects to a URL with an error flash message.
//
// Example:
//
//	flash.RedirectWithError(w, r, "/auth", "Invalid credentials")
func RedirectWithError(w http.ResponseWriter, r *http.Request, redirectURL, message string) {
	Error(w, message)
	SetPending(w)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// RedirectWithSuccess redirects to a URL with a success flash message.
//
// Example:
//
//	flash.RedirectWithSuccess(w, r, "/", "Welcome back!")
func RedirectWithSuccess(w http.ResponseWriter, r *http.Request, redirectURL, message string) {
	Success(w, message)
	SetPending(w)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// Messages holds both error and success flash messages.
// Useful for passing to templates.
type Messages struct {
	Error   string
	Success string
}

// GetAll reads and clears both error and success flash messages.
//
// Example:
//
//	msgs := flash.GetAll(r, w)
//	state.FlashError = msgs.Error
//	state.FlashSuccess = msgs.Success
func GetAll(r *http.Request, w http.ResponseWriter) Messages {
	return Messages{
		Error:   GetError(r, w),
		Success: GetSuccess(r, w),
	}
}
