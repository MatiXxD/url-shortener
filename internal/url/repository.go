package url

import (
	"context"

	"github.com/MatiXxD/url-shortener/internal/models"
)

type Repository interface {
	AddURL(context.Context, *models.URL) (string, error)
	BatchAddURL(context.Context, []*models.URL) ([]*models.URL, error)
	GetURL(context.Context, string) (*models.URL, error)
}
