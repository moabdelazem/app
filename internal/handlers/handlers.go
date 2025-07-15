package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/moabdelazem/app/internal/database"
	"github.com/moabdelazem/app/internal/models"
	"github.com/moabdelazem/app/internal/repository"
	"github.com/moabdelazem/app/internal/service"
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

type GuestBookHandler struct {
	service *service.GuestBookService
}

func NewGuestBookHandler(db *database.DB) *GuestBookHandler {
	return &GuestBookHandler{
		service: service.NewGuestBookService(repository.NewGuestBookRepository(db)),
	}
}

// GetGuestBookMessages handles GET /api/v1/guestbook
func (h *GuestBookHandler) GetGuestBookMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	messages, total, err := h.service.GetMessages(ctx, page, pageSize)
	if err != nil {
		slog.Error("Failed to get guest book messages", "error", err)
		RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve messages",
		})
		return
	}

	// Calculate pagination info
	totalPages := (total + pageSize - 1) / pageSize

	response := map[string]interface{}{
		"messages": messages,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	}

	RespondJSON(w, http.StatusOK, response)
}

// GetGuestBookMessage handles GET /api/v1/guestbook/{id}
func (h *GuestBookHandler) GetGuestBookMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	message, err := h.service.GetMessageByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get guest book message", "id", id, "error", err)
		RespondJSON(w, http.StatusNotFound, map[string]string{
			"error": "Message not found",
		})
		return
	}

	RespondJSON(w, http.StatusOK, message)
}

// CreateGuestBookMessage handles POST /api/v1/guestbook
func (h *GuestBookHandler) CreateGuestBookMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createMsg models.CreateGuestBookMessage
	if err := json.NewDecoder(r.Body).Decode(&createMsg); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	message, err := h.service.CreateMessage(ctx, &createMsg)
	if err != nil {
		slog.Error("Failed to create guest book message", "error", err)
		RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	slog.Info("Created new guest book message", "id", message.ID, "name", message.Name)
	RespondJSON(w, http.StatusCreated, message)
}

// HealthHandler handles health check requests with database connectivity check
func HealthHandlerWithDB(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// Check database health
		if err := db.Health(ctx); err != nil {
			slog.Error("Database health check failed", "error", err)
			RespondJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unhealthy",
				"error":  "Database connection failed",
			})
			return
		}

		RespondJSON(w, http.StatusOK, map[string]string{
			"status":   "healthy",
			"database": "connected",
		})
	}
}

// NotFoundHandler handles 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	slog.Warn("Route not found", "method", r.Method, "path", r.URL.Path)
	RespondJSON(w, http.StatusNotFound, map[string]interface{}{
		"error":   "Not Found",
		"message": "The requested resource was not found",
		"path":    r.URL.Path,
		"method":  r.Method,
	})
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	slog.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
	RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
		"error":   "Method Not Allowed",
		"message": "The request method is not supported for this resource",
		"path":    r.URL.Path,
		"method":  r.Method,
	})
}

// APIInfoHandler provides information about available endpoints
func APIInfoHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received request on API info endpoint")

	apiInfo := map[string]interface{}{
		"name":        "Guest Book API",
		"version":     "v1",
		"description": "A simple guest book API for managing messages",
		"endpoints": map[string]interface{}{
			"GET /":                      "API information",
			"GET /health":                "Basic health check",
			"GET /api/v1/health":         "Health check with database connectivity",
			"GET /api/v1/guestbook":      "Get all guest book messages (supports pagination: ?page=1&page_size=10)",
			"POST /api/v1/guestbook":     "Create a new guest book message",
			"GET /api/v1/guestbook/{id}": "Get a specific guest book message by ID",
		},
		"example_request": map[string]interface{}{
			"POST /api/v1/guestbook": map[string]interface{}{
				"name":    "John Doe",
				"email":   "john.doe@example.com",
				"message": "Hello! This is my message in the guest book.",
			},
		},
	}

	RespondJSON(w, http.StatusOK, apiInfo)
}
