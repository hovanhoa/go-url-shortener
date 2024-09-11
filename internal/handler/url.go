package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hovanhoa/go-url-shortener/internal/entities"
	"github.com/hovanhoa/go-url-shortener/internal/service"
	"net/http"
)

type (
	URLHandler interface {
		AddNewURL(c *gin.Context)
	}

	urlHandler struct {
		service.URLService
	}
)

func (u *urlHandler) AddNewURL(c *gin.Context) {
	var url *entities.URL
	if err := c.Bind(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The format is invalid"})
		return
	}

	if url.LongURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The URL is missing"})
		return
	}

	err := u.URLService.AddNewURL(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
