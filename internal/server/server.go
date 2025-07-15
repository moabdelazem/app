package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/moabdelazem/app/internal/config"
	"github.com/moabdelazem/app/internal/database"
	"github.com/moabdelazem/app/internal/handlers"
	"github.com/moabdelazem/app/internal/repository"
	"github.com/moabdelazem/app/internal/service"
)

type Server struct {
	router           *mux.Router
	config           config.Config
	server           *http.Server
	db               *database.DB
	guestBookHandler *handlers.GuestBookHandler
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

	// Root endpoint - API information
	s.router.HandleFunc("/", handlers.APIInfoHandler).Methods("GET")

	// Health endpoint (basic)
	s.router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	// Health endpoint with database check
	api.HandleFunc("/health", handlers.HealthHandlerWithDB(s.db)).Methods("GET")

	// Guest book endpoints
	// GET /api/v1/guestbook - Get all messages with pagination
	api.HandleFunc("/guestbook", s.guestBookHandler.GetGuestBookMessages).Methods("GET")

	// POST /api/v1/guestbook - Create a new message
	api.HandleFunc("/guestbook", s.guestBookHandler.CreateGuestBookMessage).Methods("POST")

	// GET /api/v1/guestbook/{id} - Get specific message (only numeric IDs)
	api.HandleFunc("/guestbook/{id:[0-9]+}", s.guestBookHandler.GetGuestBookMessage).Methods("GET")

	// Set custom 404 and 405 handlers
	s.router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	s.router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)

	// Add middleware for logging
	s.router.Use(s.loggingMiddleware)

	// Add CORS middleware
	s.router.Use(s.corsMiddleware)
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

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start() error {
	slog.Info("Starting server", "port", s.config.Port)

	// Connect to database
	if err := s.initializeDatabase(); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		return err
	}

	// Register routes after database is initialized
	s.RegisterRoutes()

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	return nil
}

func (s *Server) initializeDatabase() error {
	ctx := context.Background()

	// Create database connection
	db, err := database.NewConnection(ctx, &s.config)
	if err != nil {
		return err
	}
	s.db = db

	// Create guest book handler
	s.guestBookHandler = handlers.NewGuestBookHandler(db)

	// Initialize database tables
	guestBookService := service.NewGuestBookService(repository.NewGuestBookRepository(db))
	if err := guestBookService.InitializeDatabase(ctx); err != nil {
		return err
	}

	slog.Info("Database initialized successfully")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server...")

	// Close database connection
	if s.db != nil {
		s.db.Close()
	}

	return s.server.Shutdown(ctx)
}
