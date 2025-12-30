// Package password provides secure password hashing and verification using bcrypt.
package password

import "golang.org/x/crypto/bcrypt"

// Hash generates a bcrypt hash from a plain-text password.
// Uses bcrypt.DefaultCost (currently 10) for a good balance of security and performance.
//
// Example:
//
//	hash, err := password.Hash("mysecretpassword")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Store hash in database
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Verify checks if a plain-text password matches a bcrypt hash.
// Returns true if the password matches, false otherwise.
//
// Example:
//
//	if password.Verify("mysecretpassword", storedHash) {
//	    // Password is correct
//	}
func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
