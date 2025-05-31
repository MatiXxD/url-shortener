package url

type Usecase interface {
	ReduceURL(url string) (string, error)
	GetURL(shortURL string) (string, bool)
}
