package token

import (
	"encoding/base64"
	"testing"
)

func TestGenerate(t *testing.T) {
	tok, err := Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if tok == "" {
		t.Fatal("Generate() returned empty string")
	}

	// Verify it's valid base64 URL encoding
	decoded, err := base64.URLEncoding.DecodeString(tok)
	if err != nil {
		t.Fatalf("Generate() returned invalid base64: %v", err)
	}

	// Should be 32 bytes when decoded
	if len(decoded) != 32 {
		t.Fatalf("Generate() decoded length = %d, want 32", len(decoded))
	}
}

func TestGenerateUniqueness(t *testing.T) {
	tokens := make(map[string]bool)
	count := 100

	for i := 0; i < count; i++ {
		tok, err := Generate()
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		if tokens[tok] {
			t.Fatalf("Generate() produced duplicate token on iteration %d", i)
		}
		tokens[tok] = true
	}
}

func TestGenerateLength(t *testing.T) {
	tok, err := Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// 32 bytes in base64 URL encoding = 44 characters (with padding)
	expectedLen := 44
	if len(tok) != expectedLen {
		t.Fatalf("Generate() length = %d, want %d", len(tok), expectedLen)
	}
}
