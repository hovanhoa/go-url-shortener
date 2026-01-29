package logger

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/Graylog2/go-gelf.v1/gelf"
)

var (
	// Logger is the global structured logger instance
	Logger *slog.Logger
)

// Init initializes the global logger with JSON output
func Init(serviceName string) {
	InitWithGraylog(serviceName, "", 0, false)
}

// InitWithGraylog initializes the logger with optional Graylog support
func InitWithGraylog(serviceName string, graylogAddr string, graylogPort int, graylogEnabled bool) {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	var handler slog.Handler

	// Initialize Graylog if enabled
	if graylogEnabled && graylogAddr != "" && graylogPort > 0 {
		gelfAddr := fmt.Sprintf("%s:%d", graylogAddr, graylogPort)
		gelfWriter, err := gelf.NewWriter(gelfAddr)
		if err != nil {
			// Fallback to JSON only if Graylog fails
			handler = slog.NewJSONHandler(os.Stdout, opts)
			Logger = slog.New(handler).With(slog.String("service", serviceName))
			slog.SetDefault(Logger)
			slog.Error("Failed to initialize Graylog writer", "error", err, "address", gelfAddr)
		} else {
			// Use custom handler that writes to both JSON and GELF
			handler = NewGELFHandler(gelfWriter, slog.LevelInfo)
			Logger = slog.New(handler).With(slog.String("service", serviceName))
			slog.SetDefault(Logger)
			// Use fmt.Printf to ensure this message appears even if logger isn't fully initialized
			fmt.Printf("[LOGGER] Graylog logging enabled at %s\n", gelfAddr)
			slog.Info("Graylog logging enabled", "address", gelfAddr, "service", serviceName)
		}
	} else {
		// Standard JSON handler when Graylog is disabled
		handler = slog.NewJSONHandler(os.Stdout, opts)
		Logger = slog.New(handler).With(
			slog.String("service", serviceName),
		)
		// Set as default logger
		slog.SetDefault(Logger)
		if !graylogEnabled {
			fmt.Printf("[LOGGER] Graylog is disabled in config\n")
		}
	}
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
