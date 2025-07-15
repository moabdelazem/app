package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondJSON(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		payload        interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success with payload",
			status:         http.StatusOK,
			payload:        map[string]string{"message": "success"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name:           "Success with nil payload",
			status:         http.StatusNoContent,
			payload:        nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:           "Error status",
			status:         http.StatusBadRequest,
			payload:        map[string]string{"error": "bad request"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"bad request"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondJSON(w, tt.status, tt.payload)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" {
				var expected, actual map[string]interface{}
				if err := json.Unmarshal([]byte(tt.expectedBody), &expected); err != nil {
					t.Fatalf("Failed to unmarshal expected body: %v", err)
				}
				if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
					t.Fatalf("Failed to unmarshal actual body: %v", err)
				}

				if len(expected) != len(actual) {
					t.Errorf("Expected %v, got %v", expected, actual)
				}

				for key, expectedValue := range expected {
					if actual[key] != expectedValue {
						t.Errorf("Expected %s to be %v, got %v", key, expectedValue, actual[key])
					}
				}
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type to be application/json, got %s", contentType)
			}
		})
	}
}

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	HomeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedMessage := "This is API v1"
	if response["message"] != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, response["message"])
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedStatus := "healthy"
	if response["status"] != expectedStatus {
		t.Errorf("Expected status %q, got %q", expectedStatus, response["status"])
	}
}

func TestAPIInfoHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	APIInfoHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check required fields
	requiredFields := []string{"name", "version", "description", "endpoints"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Expected field %q to exist in response", field)
		}
	}

	// Check API name
	if response["name"] != "Guest Book API" {
		t.Errorf("Expected name to be 'Guest Book API', got %v", response["name"])
	}

	// Check version
	if response["version"] != "v1" {
		t.Errorf("Expected version to be 'v1', got %v", response["version"])
	}
}

func TestNotFoundHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	NotFoundHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check error fields
	expectedFields := map[string]interface{}{
		"error":  "Not Found",
		"path":   "/nonexistent",
		"method": "GET",
	}

	for field, expectedValue := range expectedFields {
		if response[field] != expectedValue {
			t.Errorf("Expected %s to be %v, got %v", field, expectedValue, response[field])
		}
	}
}

func TestMethodNotAllowedHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	MethodNotAllowedHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check error fields
	expectedFields := map[string]interface{}{
		"error":  "Method Not Allowed",
		"path":   "/health",
		"method": "POST",
	}

	for field, expectedValue := range expectedFields {
		if response[field] != expectedValue {
			t.Errorf("Expected %s to be %v, got %v", field, expectedValue, response[field])
		}
	}
}
