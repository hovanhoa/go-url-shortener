package storage

import (
	"context"

	"github.com/hovanhoa/go-url-shortener/internal/entities"
	"github.com/hovanhoa/go-url-shortener/pkg/otel"
	"gorm.io/gorm"
)

type (
	URLRepository interface {
		AddNewURL(ctx context.Context, url *entities.URL) (*entities.URL, error)
		FindOneURL(ctx context.Context, url *entities.URL) (*entities.URL, error)
	}

	urlRepository struct {
		*gorm.DB
	}
)

func (u *urlRepository) AddNewURL(ctx context.Context, url *entities.URL) (*entities.URL, error) {
	var result *entities.URL
	err := otel.GormWithTracing(ctx, u.DB, "create", func(db *gorm.DB) error {
		if err := db.Create(&url).Error; err != nil {
			return err
		}
		result = url
		return nil
	})
	return result, err
}

func (u *urlRepository) FindOneURL(ctx context.Context, url *entities.URL) (*entities.URL, error) {
	var result *entities.URL
	err := otel.GormWithTracing(ctx, u.DB, "find", func(db *gorm.DB) error {
		if err := db.Where(url).First(&url).Error; err != nil {
			return err
		}
		result = url
		return nil
	})
	return result, err
}
