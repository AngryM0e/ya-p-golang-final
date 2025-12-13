package api

import (
	"encoding/json"
	"net/http"

	"github.com/AngryM0e/ya-p-golang-final/pkg/auth"
)

type SignInRequest struct{
	Password string `json:"password"`
}

// SignInResponse - response for signin
type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// SignInHandler handler for /api/signin
func SignInHandler (w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse Request
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return // ДОБАВИТЬ return
	}

	// Getting password from variable environment
	expectedPassword := auth.GetPassword()

	// Check password
	if req.Password != expectedPassword {
		response := SignInResponse{Error: "Неверный пароль"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate token
	token := auth.GenerateToken(req.Password)

	// Set cookie
	http.SetCookie(w, &http.Cookie {
		Name: "token",
		Value: token,
		Path: "/",
		MaxAge: 8 * 60 * 60, // 8 часов
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteStrictMode,
	})

	// Return token
	response := SignInResponse{Token: token}
	writeJSONSuccess(w, response)
}