package service

import "github.com/hovanhoa/go-url-shortener/internal/storage"

type Service struct {
	URL URLService
}

func New(s *storage.Repository) *Service {
	return &Service{
		URL: &urlService{storage: s},
	}
}
