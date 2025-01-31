package repository

import (
	"database/sql"
	"errors"
	"github.com/example-module/url-shortener/config"
	"os"
	"strings"
)

type PostgresUrlRepository struct {
	Db *sql.DB
}

func NewPostgresRepository() (*PostgresUrlRepository, error) {
	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		return nil, errors.New("DB_CONN_STR environment variable is not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresUrlRepository{Db: db}, nil
}

func (p *PostgresUrlRepository) GetShortUrl(original string) (string, error) {
	query := "SELECT url_shortened FROM urls WHERE url_original = $1"
	var shortUrl string
	err := p.Db.QueryRow(query, original).Scan(&shortUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("url not found")
		}
		return "", err
	}
	return shortUrl, nil
}

func (p *PostgresUrlRepository) SaveUrl(original, shortKey string) (string, error) {
	query := "INSERT INTO urls (url_original, url_shortened) VALUES ($1, $2)"
	shortUrl := *config.BaseUrl + shortKey
	_, err := p.Db.Exec(query, original, shortUrl)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "", errors.New("duplicate key error")
		}
		return "", err
	}
	return shortUrl, nil
}

func (p *PostgresUrlRepository) GetOriginalUrl(short string) (string, error) {
	short = *config.BaseUrl + short
	query := "SELECT url_original FROM urls WHERE url_shortened = $1"
	var originalUrl string
	err := p.Db.QueryRow(query, short).Scan(&originalUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("url not found")
		}
		return "", err
	}
	return originalUrl, nil
}
