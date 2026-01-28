package handler

import (
	"github.com/gin-gonic/gin"
)

// PanicHandler triggers a panic to test error handling and observability
func PanicHandler(c *gin.Context) {
	// Cause a nil pointer dereference (invalid memory access)
	var ptr *int
	*ptr = 42 // This will cause a panic: runtime error: invalid memory address or nil pointer dereference
}
