package usecase

import (
	"github.com/example-module/url-shortener/internal/domain"
)

type UrlRepository interface {
	GetShortUrl(original string) (string, error)
	SaveUrl(original, shortKey string) (string, error)
	GetOriginalUrl(short string) (string, error)
}

type UrlShortenerUseCase struct {
	Repo UrlRepository
}

func (u *UrlShortenerUseCase) ShortenUrl(original string) (string, error) {
	shortUrl, err := u.Repo.GetShortUrl(original)
	if err == nil {
		return shortUrl, nil
	}

	url, err := domain.NewUrl(original)
	if err != nil {
		return "", err
	}
	shortUrl, err = u.Repo.SaveUrl(url.Original, url.ShortKey)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

func (u *UrlShortenerUseCase) GerOriginalUrl(short string) (string, error) {
	return u.Repo.GetOriginalUrl(short)
}
