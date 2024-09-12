package handler

import (
	"github.com/hovanhoa/go-url-shortener/internal/service"
	"github.com/hovanhoa/go-url-shortener/pkg/snowflake"
)

type Handler struct {
	URLHandler
}

func New(s *service.Service, n *snowflake.Node) *Handler {
	return &Handler{
		URLHandler: &urlHandler{s.URL, n},
	}
}
