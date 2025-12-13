package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

// GenerateToken create token from password
func GenerateToken(password string) string {
	// Create hash from password
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// ValidateToken check token
func ValidateToken(token string) bool {
	password := GetPassword()
	if password == "" {
		return true // if password is empty, auth is not required
	}
	
	// Check token with password
	expectedToken := GenerateToken(password)
	return token == expectedToken
}

// IsAuthRequired check if auth is required
func IsAuthRequired() bool {
	return GetPassword() != ""
}

// GetPassword get password from environment variable
func GetPassword() string {
	return os.Getenv("TODO_PASSWORD")
}