package service

import "github.com/hovanhoa/go-url-shortener/internal/storage"

type (
	URLService interface {
	}

	urlService struct {
		storage *storage.Repository
	}
)
