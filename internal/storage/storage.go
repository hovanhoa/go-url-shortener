package storage

import (
	"github.com/hovanhoa/go-url-shortener/internal/entities"
	"gorm.io/gorm"
)

type Repository struct {
	DB  *gorm.DB
	URL URLRepository
}

func New(db *gorm.DB) *Repository {
	db.AutoMigrate(&entities.URL{})

	return &Repository{
		DB:  db,
		URL: &urlRepository{db},
	}
}
