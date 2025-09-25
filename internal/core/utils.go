package core

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateSecureTokenBase64 generates a secure random token and returns it as base64 encoded string
func GenerateSecureTokenBase64() (string, error) {
	// Generate 32 random bytes (256 bits)
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode to base64 for safe transmission
	return base64.URLEncoding.EncodeToString(bytes), nil
}
