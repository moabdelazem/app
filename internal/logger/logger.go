package logger

import (
	"log/slog"
	"os"

	"github.com/moabdelazem/app/internal/config"
)

// Initialize sets up the structured logger with config
func Initialize(cfg config.Config) {
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)
}
