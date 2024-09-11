package storage

import "database/sql"

type Repository struct {
	DB  *sql.DB
	URL URLRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{
		DB:  db,
		URL: &urlRepository{db},
	}
}
