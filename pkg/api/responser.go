package api

import (
	"encoding/json"
	"net/http"
)

// writeJSONError writes an error response in JSON format
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// writeJSONSuccess writes a success response in JSON format
func writeJSONSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-type", "application/json;  charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

// writeJSONEmpty writes an empty response in JSON format
func writeJSONEmpty(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
