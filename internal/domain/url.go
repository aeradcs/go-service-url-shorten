package domain

import (
	"crypto/sha256"
	"encoding/base64"
)

type Url struct {
	Original string
	ShortKey string
}

func NewUrl(original string) (*Url, error) {
	return &Url{
		Original: original,
		ShortKey: GenerateShortKey(original),
	}, nil
}

func GenerateShortKey(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	hashString := base64.URLEncoding.EncodeToString(hashBytes)
	shortKey := hashString[:7]
	return shortKey
}
