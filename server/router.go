package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hovanhoa/go-url-shortener/internal/handler"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", handler.Health)

	return router
}
