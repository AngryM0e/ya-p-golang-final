package tests

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

var Port = 7540
var DBFile = getTestDBPath()
var FullNextDate = true
var Search = false
var Token = generateTestToken()

// getTestDBPath return path to test database
func getTestDBPath() string {
	// Try getting path from environment variable
	if envPath := os.Getenv("TODO_DBFILE"); envPath != "" {
		return resolveTestPath(envPath)
	}
	
	// By default, use scheduler.db in the tests directory
	return resolveTestPath("../scheduler.db")
}

// resolveTestPath converts relative paths to absolute paths relative to the tests directory
func resolveTestPath(relativePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	
	// getting the absolute path to the tests directory
	testsDir, err := filepath.Abs(".")
	if err != nil {
		return relativePath
	}
	
	// Combine paths
	fullPath := filepath.Join(testsDir, relativePath)
	
	// Normalize path to avoid problems with . and .. in the path
	cleanPath, err := filepath.Abs(fullPath)
	if err != nil {
		return fullPath
	}
	
	return cleanPath
}

func generateTestToken() string {
	// Getting password from environment variable
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return "" // If the password is not set, return an empty token
	}

	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}