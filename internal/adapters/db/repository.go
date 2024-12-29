package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type UrlRepository interface {
	GetShortUrl(original string) (string, error)
	InsertUrl(original string, shortKey string) error
}

type urlRepository struct {
	db *sql.DB
}

func NewUrlRepository(db *sql.DB) UrlRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) GetShortUrl(original string) (string, error) {
	var shortUrl string
	query := "SELECT url_shortened FROM urls WHERE url_original = $1"
	err := r.db.QueryRow(query, original).Scan(&shortUrl)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

func (r *urlRepository) InsertUrl(original string, shortKey string) error {
	query := "INSERT INTO urls (url_original, url_shortened) VALUES ($1, $2)"
	_, err := r.db.Exec(query, original, "http://localhost:8080/"+shortKey)
	return err
}
