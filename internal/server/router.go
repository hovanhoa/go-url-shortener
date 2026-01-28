package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hovanhoa/go-url-shortener/config"
	"github.com/hovanhoa/go-url-shortener/internal/handler"
	"github.com/hovanhoa/go-url-shortener/internal/middleware/ratelimit"
	"github.com/hovanhoa/go-url-shortener/internal/middleware/timeout"
	"github.com/hovanhoa/go-url-shortener/pkg/apm"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"err": "Too many requests",
	})
}

func NewRouter(h *handler.Handler) *gin.Engine {
	cfg := config.GetConfig()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Prometheus APM
	prom := apm.NewPrometheus("go-url-shortener")
	router.Use(prom.GinMiddleware())
	router.GET("/metrics", gin.WrapH(prom.Handler()))

	// middleware for rate limiter
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  cfg.RateLimit.Rate,
		Limit: cfg.RateLimit.Limit,
	})
	rateLimiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	// middleware for timeout request
	timeoutMiddleware := timeout.TimeOutMiddleware(cfg.TimeOut.Time)
	router.Use(rateLimiter, timeoutMiddleware)

	router.GET("/health", handler.Health)
	router.POST("/sl", h.URLHandler.AddNewURL)
	router.GET("/sl/:url", h.URLHandler.GetURL)

	return router
}
