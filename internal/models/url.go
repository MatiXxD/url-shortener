package models

import (
	"time"
)

//go:generate easyjson -all url.go

type UrlDTO struct {
	CorrelationID string `json:"correlation_id"`
	OriginURL     string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}

type URL struct {
	ID            int       `json:"id,omitempty"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	BaseURL       string    `json:"original_url"`
	ShortURL      string    `json:"short_url"`
	CreateAt      time.Time `json:"created_ad,omitempty"`
	IsDeleted     bool      `json:"deleted,omitempty"`
}

type ShortenURLReqBody struct {
	URL string `json:"url"`
}

type ShortenURLRespBody struct {
	ShortURL string `json:"short_url"`
}
