package usecase

import (
	"fmt"

	"github.com/MatiXxD/url-shortener/internal/repository"
)

type UrlUsecase struct {
	repo *repository.MapRepository
}

func NewUrlUsecase(r *repository.MapRepository) *UrlUsecase {
	return &UrlUsecase{
		repo: r,
	}
}

func (uu *UrlUsecase) ReduceURL(url string) (string, error) {
	s, err := uu.repo.ReduceURL(url)
	if err != nil {
		return "", fmt.Errorf("Can't reduce url")
	}
	return s, nil
}

func (uu *UrlUsecase) GetURL(shortURL string) (string, bool) {
	url, ok := uu.repo.GetURL(shortURL)
	if !ok {
		return "", false
	}
	return url, true
}
