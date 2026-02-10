package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// GenerateToken generates a cryptographically secure random token.
// Returns 32 bytes encoded as base64url (RawURLEncoding).
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashToken returns the SHA256 hash of the raw token as a hex string.
func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return fmt.Sprintf("%x", hash)
}
