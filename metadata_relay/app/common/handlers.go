package common

import (
	"encoding/json"
	"net/http"
)

// WriteJSONResponse writes a JSON response with HTTP 200 status.
func WriteJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteErrorResponse writes an error response with the specified status code.
func WriteErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
