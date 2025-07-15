package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/moabdelazem/app/internal/models"
)

func TestGuestBookHandler_GetGuestBookMessages(t *testing.T) {
	mockService := NewMockGuestBookService()
	handler := NewGuestBookHandlerWithService(mockService)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Get all messages - default pagination",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Get messages with pagination",
			queryParams:    "?page=1&page_size=1",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "Get messages with invalid page",
			queryParams:    "?page=0&page_size=10",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Get messages with large page size",
			queryParams:    "?page=1&page_size=1000",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guestbook"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler.GetGuestBookMessages(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			messages, ok := response["messages"].([]interface{})
			if !ok {
				t.Fatal("Expected messages to be an array")
			}

			if len(messages) != tt.expectedCount {
				t.Errorf("Expected %d messages, got %d", tt.expectedCount, len(messages))
			}

			// Check pagination structure
			pagination, ok := response["pagination"].(map[string]interface{})
			if !ok {
				t.Fatal("Expected pagination to be an object")
			}

			expectedPaginationFields := []string{"page", "page_size", "total", "total_pages"}
			for _, field := range expectedPaginationFields {
				if _, exists := pagination[field]; !exists {
					t.Errorf("Expected pagination field %q to exist", field)
				}
			}
		})
	}
}

func TestGuestBookHandler_GetGuestBookMessage(t *testing.T) {
	mockService := NewMockGuestBookService()
	handler := NewGuestBookHandlerWithService(mockService)

	tests := []struct {
		name           string
		messageID      string
		expectedStatus int
		expectedID     int
	}{
		{
			name:           "Get existing message",
			messageID:      "1",
			expectedStatus: http.StatusOK,
			expectedID:     1,
		},
		{
			name:           "Get non-existent message",
			messageID:      "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Get message with invalid ID",
			messageID:      "invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guestbook/"+tt.messageID, nil)
			w := httptest.NewRecorder()

			// Set up mux vars to simulate route parameters
			req = mux.SetURLVars(req, map[string]string{"id": tt.messageID})

			handler.GetGuestBookMessage(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.GuestBookMessage
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.ID != tt.expectedID {
					t.Errorf("Expected message ID %d, got %d", tt.expectedID, response.ID)
				}

				// Check required fields
				if response.Name == "" {
					t.Error("Expected name to not be empty")
				}
				if response.Email == "" {
					t.Error("Expected email to not be empty")
				}
				if response.Message == "" {
					t.Error("Expected message to not be empty")
				}
			} else {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}

				if _, exists := response["error"]; !exists {
					t.Error("Expected error field in response")
				}
			}
		})
	}
}

func TestGuestBookHandler_CreateGuestBookMessage(t *testing.T) {
	mockService := NewMockGuestBookService()
	handler := NewGuestBookHandlerWithService(mockService)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, response []byte)
	}{
		{
			name: "Create valid message",
			requestBody: models.CreateGuestBookMessage{
				Name:    "Test User",
				Email:   "test@example.com",
				Message: "This is a test message for the guest book.",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, response []byte) {
				var msg models.GuestBookMessage
				if err := json.Unmarshal(response, &msg); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if msg.ID == 0 {
					t.Error("Expected non-zero ID")
				}
				if msg.Name != "Test User" {
					t.Errorf("Expected name 'Test User', got %q", msg.Name)
				}
				if msg.Email != "test@example.com" {
					t.Errorf("Expected email 'test@example.com', got %q", msg.Email)
				}
			},
		},
		{
			name: "Create message with name too short",
			requestBody: models.CreateGuestBookMessage{
				Name:    "A",
				Email:   "test@example.com",
				Message: "This is a test message for the guest book.",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response []byte) {
				var errorResp map[string]string
				if err := json.Unmarshal(response, &errorResp); err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}

				if !strings.Contains(errorResp["error"], "name must be between") {
					t.Errorf("Expected name validation error, got %q", errorResp["error"])
				}
			},
		},
		{
			name: "Create message with empty email",
			requestBody: models.CreateGuestBookMessage{
				Name:    "Test User",
				Email:   "",
				Message: "This is a test message for the guest book.",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response []byte) {
				var errorResp map[string]string
				if err := json.Unmarshal(response, &errorResp); err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}

				if !strings.Contains(errorResp["error"], "email must be between") {
					t.Errorf("Expected email validation error, got %q", errorResp["error"])
				}
			},
		},
		{
			name: "Create message with message too short",
			requestBody: models.CreateGuestBookMessage{
				Name:    "Test User",
				Email:   "test@example.com",
				Message: "Short",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response []byte) {
				var errorResp map[string]string
				if err := json.Unmarshal(response, &errorResp); err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}

				if !strings.Contains(errorResp["error"], "message must be between") {
					t.Errorf("Expected message validation error, got %q", errorResp["error"])
				}
			},
		},
		{
			name:           "Create message with invalid JSON",
			requestBody:    `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response []byte) {
				var errorResp map[string]string
				if err := json.Unmarshal(response, &errorResp); err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}

				if errorResp["error"] != "Invalid request body" {
					t.Errorf("Expected 'Invalid request body' error, got %q", errorResp["error"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/guestbook", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateGuestBookMessage(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestGuestBookHandler_CreateGuestBookMessage_EdgeCases(t *testing.T) {
	mockService := NewMockGuestBookService()
	handler := NewGuestBookHandlerWithService(mockService)

	tests := []struct {
		name           string
		requestBody    models.CreateGuestBookMessage
		expectedStatus int
	}{
		{
			name: "Name at minimum length",
			requestBody: models.CreateGuestBookMessage{
				Name:    "AB",
				Email:   "test@example.com",
				Message: "This is a test message for the guest book.",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Name at maximum length",
			requestBody: models.CreateGuestBookMessage{
				Name:    strings.Repeat("A", 100),
				Email:   "test@example.com",
				Message: "This is a test message for the guest book.",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Message at minimum length",
			requestBody: models.CreateGuestBookMessage{
				Name:    "Test User",
				Email:   "test@example.com",
				Message: "1234567890", // exactly 10 characters
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Message at maximum length",
			requestBody: models.CreateGuestBookMessage{
				Name:    "Test User",
				Email:   "test@example.com",
				Message: strings.Repeat("A", 1000),
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/guestbook", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateGuestBookMessage(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
