package handler

import "github.com/hovanhoa/go-url-shortener/internal/service"

type Handler struct {
	URLHandler
}

func New(s *service.Service) *Handler {
	return &Handler{
		URLHandler: &urlHandler{s.URL},
	}
}
