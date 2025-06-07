package models

import "github.com/google/uuid"

//go:generate easyjson -all url.go

type UrlDTO struct {
	URL string `json:"url"`
}

func NewUrlDTO(url string) *UrlDTO {
	return &UrlDTO{
		URL: url,
	}
}

type URL struct {
	ID       uuid.UUID `json:"id"`
	BaseURL  string    `json:"baseUrl"`
	ShortURL string    `json:"shortUrl"`
}

func NewURL(url, shortURL string) *URL {
	return &URL{
		ID:       uuid.New(),
		BaseURL:  url,
		ShortURL: shortURL,
	}
}
