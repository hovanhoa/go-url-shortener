package service

import (
	"github.com/hovanhoa/go-url-shortener/internal/storage"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	URL URLService
}

func New(s *storage.Repository, r *redis.Client) *Service {
	return &Service{
		URL: &urlService{storage: s, redis: r},
	}
}
