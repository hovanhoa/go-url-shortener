package logger

import (
	"log/slog"
	"os"
)

var (
	// Logger is the global structured logger instance
	Logger *slog.Logger
)

// Init initializes the global logger with JSON output
func Init(serviceName string) {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	Logger = slog.New(handler).With(
		slog.String("service", serviceName),
	)

	// Set as default logger
	slog.SetDefault(Logger)
}

// InitWithLevel initializes the logger with a specific log level
func InitWithLevel(serviceName string, level slog.Level) {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	Logger = slog.New(handler).With(
		slog.String("service", serviceName),
	)

	slog.SetDefault(Logger)
}

// GetLogger returns the global logger instance
func GetLogger() *slog.Logger {
	if Logger == nil {
		// Fallback to default if not initialized
		Init("go-url-shortener")
	}
	return Logger
}
