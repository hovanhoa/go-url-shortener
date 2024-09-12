package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hovanhoa/go-url-shortener/internal/entities"
	"github.com/hovanhoa/go-url-shortener/internal/service"
	"github.com/hovanhoa/go-url-shortener/pkg/base62"
	"github.com/hovanhoa/go-url-shortener/pkg/snowflake"
	"github.com/hovanhoa/go-url-shortener/pkg/url"
	"gorm.io/gorm"
	"net/http"
)

type (
	URLHandler interface {
		AddNewURL(c *gin.Context)
		GetURL(c *gin.Context)
	}

	urlHandler struct {
		service.URLService
		*snowflake.Node
	}
)

func (urlHandler *urlHandler) AddNewURL(c *gin.Context) {
	var u *entities.URL
	if err := c.Bind(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the format is invalid"})
		return
	}

	if u.LongURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the long url is missing"})
		return
	}

	if !url.IsValidURL(u.LongURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the url is invalid format"})
		return
	}

	existURL, err := urlHandler.URLService.FindOneByLongURL(u.LongURL)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server"})
		return
	} else if existURL != nil {
		c.JSON(http.StatusOK, existURL)
		return
	}

	u.SortURL = base62.Encode(uint64(urlHandler.Node.Generate().Int64()))
	newURL, err := urlHandler.URLService.AddNewURL(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server"})
		return
	}

	c.JSON(http.StatusOK, newURL)
	return
}

func (urlHandler *urlHandler) GetURL(c *gin.Context) {
	shortURL := c.Param("url")
	u, err := urlHandler.URLService.FindOneByShortURL(shortURL)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "the url is not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server"})
		return
	}

	c.JSON(http.StatusOK, u)
	return
}
