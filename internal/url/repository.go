package url

type Repository interface {
	AddURL(url, shortURL string) (string, error)
	GetURL(shortURL string) (string, bool)
}
