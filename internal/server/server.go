package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/moabdelazem/app/internal/config"
	"github.com/moabdelazem/app/internal/handlers"
)

type Server struct {
	router *mux.Router
	config config.Config
	server *http.Server
}

func NewServer(cfg config.Config) *Server {
	r := mux.NewRouter()
	return &Server{
		router: r,
		config: cfg,
		server: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      r,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) RegisterRoutes() {
	// API v1 routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Root endpoint
	s.router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Health endpoint
	s.router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	api.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	// Add middleware for logging
	s.router.Use(s.loggingMiddleware)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

func (s *Server) Start() error {
	slog.Info("Starting server", "port", s.config.Port)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server...")
	return s.server.Shutdown(ctx)
}
