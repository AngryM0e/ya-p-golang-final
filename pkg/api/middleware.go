package api

import (
	"net/http"

	"github.com/AngryM0e/ya-p-golang-final/pkg/auth"
)

// AuthMiddleware - middleware for auth
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if auth is required
		if !auth.IsAuthRequired() {
			next(w, r)
			return
		}
		
		// Get token from cookie
		var token string
		cookie, err := r.Cookie("token")
		if err == nil {
			token = cookie.Value
		}
		
		// Check token
		if !auth.ValidateToken(token) {
			writeJSONError(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	})
}