package service

import (
	"github.com/hovanhoa/go-url-shortener/internal/entities"
	"github.com/hovanhoa/go-url-shortener/internal/storage"
)

type (
	URLService interface {
		AddNewURL(url *entities.URL) error
	}

	urlService struct {
		storage *storage.Repository
	}
)

func (u *urlService) AddNewURL(url *entities.URL) error {
	if err := u.storage.URL.AddNewURL(url); err != nil {
		return err
	}

	return nil
}
