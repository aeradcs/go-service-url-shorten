package usecases

import (
	"service-url-shortener/internal/adapters/db"
	"service-url-shortener/pkg/utils"
)

type UrlShortener interface {
	GenerateShortUrl(original string) (string, error)
	GetShortUrl(original string) (string, error)
}

type urlShortener struct {
	repo db.UrlRepository
}

func NewUrlShortener(repo db.UrlRepository) UrlShortener {
	return &urlShortener{repo: repo}
}

func (u *urlShortener) GenerateShortUrl(original string) (string, error) {
	shortKey := utils.GenerateShortKey()

	err := u.repo.InsertUrl(original, shortKey)
	if err != nil {
		return "", err
	}

	return "http://localhost:8080/" + shortKey, nil
}

func (u *urlShortener) GetShortUrl(original string) (string, error) {
	return u.repo.GetShortUrl(original)
}
