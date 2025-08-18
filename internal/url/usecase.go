package url

import (
	"context"

	"github.com/MatiXxD/url-shortener/internal/models"
)

type Usecase interface {
	ReduceURL(context.Context, *models.UrlDTO) (string, error)
	BatchReduceURL(context.Context, []*models.UrlDTO) ([]*models.UrlDTO, error)
	GetURL(context.Context, string) (string, bool)
}
