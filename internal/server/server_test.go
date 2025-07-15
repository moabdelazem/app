package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/moabdelazem/app/internal/config"
)

func TestServer_Routes(t *testing.T) {
	// Create a test server without database
	cfg := config.Config{
		Port:  "8080",
		Debug: false,
	}

	server := NewServer(cfg)

	// Manually register routes without database initialization
	server.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "API info"})
	}).Methods("GET")

	server.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	// Test cases
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "GET root endpoint",
			method:         http.MethodGet,
			url:            "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET health endpoint",
			method:         http.MethodGet,
			url:            "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET non-existent endpoint",
			method:         http.MethodGet,
			url:            "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestServer_Middleware(t *testing.T) {
	cfg := config.Config{
		Port:  "8080",
		Debug: false,
	}

	server := NewServer(cfg)

	// Add a test route
	server.router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}).Methods("GET")

	// Add middleware
	server.router.Use(server.loggingMiddleware)
	server.router.Use(server.corsMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	// Check CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header to be set")
	}

	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected CORS methods header to be set")
	}
}

func TestServer_CORSMiddleware(t *testing.T) {
	cfg := config.Config{
		Port:  "8080",
		Debug: false,
	}

	server := NewServer(cfg)

	// Add a test route
	server.router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET", "POST", "OPTIONS")

	server.router.Use(server.corsMiddleware)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   bool
	}{
		{
			name:           "GET request with CORS",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "OPTIONS preflight request",
			method:         http.MethodOptions,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "POST request with CORS",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkHeaders {
				expectedHeaders := map[string]string{
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
					"Access-Control-Allow-Headers": "Content-Type, Authorization",
				}

				for header, expectedValue := range expectedHeaders {
					if w.Header().Get(header) != expectedValue {
						t.Errorf("Expected %s header to be %q, got %q", header, expectedValue, w.Header().Get(header))
					}
				}
			}
		})
	}
}

func TestServer_LoggingMiddleware(t *testing.T) {
	cfg := config.Config{
		Port:  "8080",
		Debug: false,
	}

	server := NewServer(cfg)

	// Add a test route that takes some time
	server.router.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	server.router.Use(server.loggingMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// The logging middleware should have logged the request
	// In a real test, you might want to capture the log output
}

func TestServer_Shutdown(t *testing.T) {
	cfg := config.Config{
		Port:  "0", // Use random port
		Debug: false,
	}

	server := NewServer(cfg)

	// Test shutdown without starting
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown should not return error: %v", err)
	}
}
