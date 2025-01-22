package models

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
