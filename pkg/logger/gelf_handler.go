package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gopkg.in/Graylog2/go-gelf.v1/gelf"
)

// gelfHandler is a slog handler that writes to both stdout (JSON) and Graylog (GELF)
type gelfHandler struct {
	jsonHandler slog.Handler
	gelfWriter  *gelf.Writer
	level       slog.Level
}

// NewGELFHandler creates a new handler that writes to both JSON and GELF
func NewGELFHandler(gelfWriter *gelf.Writer, level slog.Level) slog.Handler {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)

	return &gelfHandler{
		jsonHandler: jsonHandler,
		gelfWriter:  gelfWriter,
		level:       level,
	}
}

func (h *gelfHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *gelfHandler) Handle(ctx context.Context, record slog.Record) error {
	// Write to stdout as JSON
	if err := h.jsonHandler.Handle(ctx, record); err != nil {
		return err
	}

	// Write to Graylog as GELF
	if h.gelfWriter != nil {
		// Convert slog level to syslog level
		var syslogLevel int32
		switch {
		case record.Level >= slog.LevelError:
			syslogLevel = 3 // ERROR
		case record.Level >= slog.LevelWarn:
			syslogLevel = 4 // WARNING
		case record.Level >= slog.LevelInfo:
			syslogLevel = 6 // INFO
		default:
			syslogLevel = 7 // DEBUG
		}

		// Build the message
		msg := record.Message

		// Create GELF message
		gelfMsg := &gelf.Message{
			Version:  "1.1",
			Host:     getHostname(),
			Short:    msg,
			Full:     msg,
			TimeUnix: float64(record.Time.UnixNano()) / 1e9,
			Level:    syslogLevel,
			Extra:    make(map[string]interface{}),
		}

		// Add attributes as additional fields
		record.Attrs(func(a slog.Attr) bool {
			key := "_" + a.Key // GELF custom fields start with underscore
			switch a.Value.Kind() {
			case slog.KindString:
				gelfMsg.Extra[key] = a.Value.String()
			case slog.KindInt64:
				gelfMsg.Extra[key] = a.Value.Int64()
			case slog.KindUint64:
				gelfMsg.Extra[key] = a.Value.Uint64()
			case slog.KindFloat64:
				gelfMsg.Extra[key] = a.Value.Float64()
			case slog.KindBool:
				gelfMsg.Extra[key] = a.Value.Bool()
			case slog.KindTime:
				gelfMsg.Extra[key] = a.Value.Time().Format(time.RFC3339)
			case slog.KindDuration:
				gelfMsg.Extra[key] = a.Value.Duration().String()
			case slog.KindGroup:
				// Handle nested groups
				for _, attr := range a.Value.Group() {
					nestedKey := key + "_" + attr.Key
					switch attr.Value.Kind() {
					case slog.KindString:
						gelfMsg.Extra[nestedKey] = attr.Value.String()
					case slog.KindInt64:
						gelfMsg.Extra[nestedKey] = attr.Value.Int64()
					case slog.KindUint64:
						gelfMsg.Extra[nestedKey] = attr.Value.Uint64()
					case slog.KindFloat64:
						gelfMsg.Extra[nestedKey] = attr.Value.Float64()
					case slog.KindBool:
						gelfMsg.Extra[nestedKey] = attr.Value.Bool()
					}
				}
			}
			return true
		})

		// Add source information if available
		if record.PC != 0 {
			fs := record.PC
			_ = fs // Source info
		}

		// Send to Graylog
		if err := h.gelfWriter.WriteMessage(gelfMsg); err != nil {
			// Log error to stderr (don't use slog to avoid recursion)
			os.Stderr.WriteString(fmt.Sprintf("[GELF] Failed to write message: %v\n", err))
		}
	}

	return nil
}

func (h *gelfHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &gelfHandler{
		jsonHandler: h.jsonHandler.WithAttrs(attrs),
		gelfWriter:  h.gelfWriter,
		level:       h.level,
	}
}

func (h *gelfHandler) WithGroup(name string) slog.Handler {
	return &gelfHandler{
		jsonHandler: h.jsonHandler.WithGroup(name),
		gelfWriter:  h.gelfWriter,
		level:       h.level,
	}
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
