package url

type Repository interface {
	ReduceURL(url string) (string, error)
	GetURL(shortURL string) (string, bool)
}
