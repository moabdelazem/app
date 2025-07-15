package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// RespondJSON writes a JSON response with the given status code and payload
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			slog.Error("Failed to encode JSON response", "error", err)
			// Write a simple error message if JSON encoding fails
			w.Write([]byte(`{"error": "Internal server error"}`))
		}
	}
}

// HomeHandler handles requests to the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received request on root endpoint")
	RespondJSON(w, http.StatusOK, map[string]string{"message": "This is API v1"})
}

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received request on health endpoint")
	RespondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}
