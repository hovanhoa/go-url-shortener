package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// RecoveryWithOTel returns a middleware that recovers from panics and records them in OpenTelemetry traces
func RecoveryWithOTel() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get the span from context
				span := trace.SpanFromContext(c.Request.Context())
				if span.IsRecording() {
					// Record the panic error in the trace
					panicErr := fmt.Errorf("panic: %v", err)
					span.RecordError(panicErr)
					span.SetStatus(codes.Error, panicErr.Error())

					// Add panic details as attributes
					span.SetAttributes(
						attribute.String("error.type", "panic"),
						attribute.String("error.message", fmt.Sprintf("%v", err)),
						attribute.String("error.stack", string(debug.Stack())),
					)
				}

				// Log the panic
				c.Error(fmt.Errorf("panic: %v", err))

				// Return error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
					"panic": fmt.Sprintf("%v", err),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
