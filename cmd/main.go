package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moabdelazem/app/internal/config"
	"github.com/moabdelazem/app/internal/logger"
	"github.com/moabdelazem/app/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger with config
	logger.Initialize(cfg)

	// Create and configure server
	srv := server.NewServer(cfg)

	// Start server (this will now initialize database and register routes)
	if err := srv.Start(); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server gracefully stopped")
}
