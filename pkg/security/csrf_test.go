package security

import (
	"net/http/httptest"
	"testing"
)

func TestValidateOrigin(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		origin   string
		referer  string
		expected bool
	}{
		{
			name:     "matching origin",
			host:     "example.com",
			origin:   "https://example.com",
			expected: true,
		},
		{
			name:     "matching origin with port",
			host:     "example.com:8080",
			origin:   "https://example.com:8080",
			expected: true,
		},
		{
			name:     "mismatched origin",
			host:     "example.com",
			origin:   "https://evil.com",
			expected: false,
		},
		{
			name:     "subdomain attack",
			host:     "example.com",
			origin:   "https://evil.example.com",
			expected: false,
		},
		{
			name:     "localhost http",
			host:     "localhost:8080",
			origin:   "http://localhost:8080",
			expected: true,
		},
		{
			name:     "localhost to 127.0.0.1",
			host:     "localhost:8080",
			origin:   "http://127.0.0.1:8080",
			expected: true,
		},
		{
			name:     "127.0.0.1 to localhost",
			host:     "127.0.0.1:8080",
			origin:   "http://localhost:8080",
			expected: true,
		},
		{
			name:     "no origin header",
			host:     "example.com",
			origin:   "",
			expected: false,
		},
		{
			name:     "referer fallback",
			host:     "example.com",
			origin:   "",
			referer:  "https://example.com/login",
			expected: true,
		},
		{
			name:     "referer mismatch",
			host:     "example.com",
			origin:   "",
			referer:  "https://evil.com/page",
			expected: false,
		},
		{
			name:     "invalid origin URL",
			host:     "example.com",
			origin:   "://invalid",
			expected: false,
		},
		{
			name:     "origin contains host but different domain",
			host:     "example.com",
			origin:   "https://notexample.com",
			expected: false,
		},
		{
			name:     "different ports same host",
			host:     "localhost:8080",
			origin:   "http://localhost:3000",
			expected: true, // Same host, different port is OK for localhost
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "http://"+tt.host+"/login", nil)
			req.Host = tt.host
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.referer != "" {
				req.Header.Set("Referer", tt.referer)
			}

			result := ValidateOrigin(req)
			if result != tt.expected {
				t.Errorf("ValidateOrigin() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateOriginAllowEmpty(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		origin   string
		referer  string
		expected bool
	}{
		{
			name:     "no headers allowed",
			host:     "example.com",
			origin:   "",
			referer:  "",
			expected: true,
		},
		{
			name:     "valid origin",
			host:     "example.com",
			origin:   "https://example.com",
			expected: true,
		},
		{
			name:     "invalid origin",
			host:     "example.com",
			origin:   "https://evil.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "http://"+tt.host+"/api/login", nil)
			req.Host = tt.host
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.referer != "" {
				req.Header.Set("Referer", tt.referer)
			}

			result := ValidateOriginAllowEmpty(req)
			if result != tt.expected {
				t.Errorf("ValidateOriginAllowEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStripPort(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"example.com:8080", "example.com"},
		{"example.com", "example.com"},
		{"localhost:3000", "localhost"},
		{"127.0.0.1:8080", "127.0.0.1"},
		{"[::1]:8080", "[::1]"},
		{"[::1]", "[::1]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := stripPort(tt.input)
			if result != tt.expected {
				t.Errorf("stripPort(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"localhost", true},
		{"LOCALHOST", true},
		{"127.0.0.1", true},
		{"::1", true},
		{"[::1]", true},
		{"example.com", false},
		{"localhost.evil.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isLocalhost(tt.input)
			if result != tt.expected {
				t.Errorf("isLocalhost(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
