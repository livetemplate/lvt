package password

import (
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"simple password", "password123"},
		{"empty password", ""},
		{"unicode password", "p@ssw0rd!#$%^&*()"},
		{"max length password", "123456789012345678901234567890123456789012345678901234567890123456789012"}, // 72 bytes max
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := Hash(tt.password)
			if err != nil {
				t.Fatalf("Hash() error = %v", err)
			}

			if hash == "" {
				t.Fatal("Hash() returned empty string")
			}

			if hash == tt.password {
				t.Fatal("Hash() returned the original password (not hashed)")
			}

			if !Verify(tt.password, hash) {
				t.Fatal("Verify() returned false for correct password")
			}

			if Verify("wrong_password", hash) {
				t.Fatal("Verify() returned true for wrong password")
			}
		})
	}
}

func TestHashUniqueness(t *testing.T) {
	password := "testpassword"
	hash1, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	hash2, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("Hash() should produce unique hashes for the same password (due to salt)")
	}

	// Both should still verify correctly
	if !Verify(password, hash1) {
		t.Fatal("Verify() failed for hash1")
	}
	if !Verify(password, hash2) {
		t.Fatal("Verify() failed for hash2")
	}
}

func TestVerifyInvalidHash(t *testing.T) {
	// Should not panic with invalid hash formats
	if Verify("password", "not-a-valid-bcrypt-hash") {
		t.Fatal("Verify() should return false for invalid hash format")
	}

	if Verify("password", "") {
		t.Fatal("Verify() should return false for empty hash")
	}
}
