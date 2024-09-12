package storage

import (
	"github.com/hovanhoa/go-url-shortener/internal/entities"
	"gorm.io/gorm"
)

type (
	URLRepository interface {
		AddNewURL(url *entities.URL) (*entities.URL, error)
		FindOneURL(url *entities.URL) (*entities.URL, error)
	}

	urlRepository struct {
		*gorm.DB
	}
)

func (u *urlRepository) AddNewURL(url *entities.URL) (*entities.URL, error) {
	if err := u.DB.Create(&url).Error; err != nil {
		return nil, err
	}

	return url, nil
}

func (u *urlRepository) FindOneURL(url *entities.URL) (*entities.URL, error) {
	if err := u.DB.Where(url).First(&url).Error; err != nil {
		return nil, err
	}

	return url, nil
}
