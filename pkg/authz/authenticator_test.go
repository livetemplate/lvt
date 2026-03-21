package authz

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCookieAuthenticator_Identify_ValidToken(t *testing.T) {
	auth := NewCookieAuthenticator("users_token", func(ctx context.Context, token string) (string, error) {
		if token == "valid-token" {
			return "user-123", nil
		}
		return "", fmt.Errorf("not found")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "users_token", Value: "valid-token"})

	userID, err := auth.Identify(req)
	if err != nil {
		t.Fatalf("Identify() error = %v", err)
	}
	if userID != "user-123" {
		t.Errorf("Identify() = %q, want %q", userID, "user-123")
	}
}

func TestCookieAuthenticator_Identify_NoCookie(t *testing.T) {
	auth := NewCookieAuthenticator("users_token", func(ctx context.Context, token string) (string, error) {
		return "should-not-be-called", nil
	})

	req := httptest.NewRequest("GET", "/", nil)

	userID, err := auth.Identify(req)
	if err != nil {
		t.Fatalf("Identify() error = %v", err)
	}
	if userID != "" {
		t.Errorf("Identify() = %q, want empty for no cookie", userID)
	}
}

func TestCookieAuthenticator_Identify_InvalidToken(t *testing.T) {
	auth := NewCookieAuthenticator("users_token", func(ctx context.Context, token string) (string, error) {
		return "", fmt.Errorf("expired")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "users_token", Value: "expired-token"})

	userID, err := auth.Identify(req)
	if err != nil {
		t.Fatalf("Identify() error = %v", err)
	}
	if userID != "" {
		t.Errorf("Identify() = %q, want empty for invalid token", userID)
	}
}

func TestCookieAuthenticator_GetSessionGroup_Authenticated(t *testing.T) {
	auth := NewCookieAuthenticator("users_token", nil)

	req := httptest.NewRequest("GET", "/", nil)
	group, err := auth.GetSessionGroup(req, "user-123")
	if err != nil {
		t.Fatalf("GetSessionGroup() error = %v", err)
	}
	if group != "user-123" {
		t.Errorf("GetSessionGroup() = %q, want %q", group, "user-123")
	}
}

func TestCookieAuthenticator_GetSessionGroup_Anonymous(t *testing.T) {
	auth := NewCookieAuthenticator("users_token", nil)

	req := httptest.NewRequest("GET", "/", nil)
	group, err := auth.GetSessionGroup(req, "")
	if err != nil {
		t.Fatalf("GetSessionGroup() error = %v", err)
	}
	if group == "" {
		t.Error("GetSessionGroup() should return non-empty for anonymous")
	}
}

func TestCookieAuthenticator_GetSessionGroup_AnonymousWithCookie(t *testing.T) {
	auth := NewCookieAuthenticator("users_token", nil)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "livetemplate-id", Value: "browser-123"})

	group, err := auth.GetSessionGroup(req, "")
	if err != nil {
		t.Fatalf("GetSessionGroup() error = %v", err)
	}
	if group != "browser-123" {
		t.Errorf("GetSessionGroup() = %q, want %q", group, "browser-123")
	}
}
