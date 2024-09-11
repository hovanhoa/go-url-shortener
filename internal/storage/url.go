package storage

import "database/sql"

type (
	URLRepository interface {
	}

	urlRepository struct {
		*sql.DB
	}
)
