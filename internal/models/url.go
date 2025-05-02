package models

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
	BaseURL  string `json:"baseUrl"`
	ShortURL string `json:"shortUrl"`
}

func NewURL(url, shortURL string) *URL {
	return &URL{
		BaseURL:  url,
		ShortURL: shortURL,
	}
}
