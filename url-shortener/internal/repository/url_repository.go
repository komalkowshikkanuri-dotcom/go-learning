package repository

import (
	"database/sql"
	"fmt"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) Create(shortCode, originalURL string) error {
	_, err := r.db.Exec(
		`INSERT INTO urls (short_code, original_url)
		 VALUES ($1, $2)`,
		shortCode,
		originalURL,
	)

	if err != nil {
		return fmt.Errorf("create URL mapping: %w", err)
	}

	return nil
}

func (r *URLRepository) GetOriginalURL(shortCode string) (string, error) {
	var originalURL string

	err := r.db.QueryRow(
		`SELECT original_url
		 FROM urls
		 WHERE short_code = $1`,
		shortCode,
	).Scan(&originalURL)

	if err != nil {
		return "", err
	}

	return originalURL, nil
}