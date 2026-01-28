package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// GinJSONLogger returns a Gin middleware for JSON structured logging
func GinJSONLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get trace ID from context if available
		span := trace.SpanFromContext(c.Request.Context())
		traceID := ""
		if span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
		}

		// Build log attributes
		attrs := []interface{}{
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"latency", latency.String(),
			"latency_ms", latency.Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		}

		if raw != "" {
			attrs = append(attrs, "query", raw)
		}

		if traceID != "" {
			attrs = append(attrs, "trace_id", traceID)
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
			// Add individual error details
			for i, err := range c.Errors {
				attrs = append(attrs, fmt.Sprintf("error_%d", i), err.Error())
			}
		}

		// Log based on status code
		if c.Writer.Status() >= 500 {
			GetLogger().Error("HTTP request", attrs...)
		} else if c.Writer.Status() >= 400 {
			GetLogger().Warn("HTTP request", attrs...)
		} else {
			GetLogger().Info("HTTP request", attrs...)
		}
	}
}
