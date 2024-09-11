package handler

import "github.com/hovanhoa/go-url-shortener/internal/service"

type (
	URLHandler interface {
	}

	urlHandler struct {
		service.URLService
	}
)
