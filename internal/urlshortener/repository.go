package urlshortener

import (
	"context"

	"github.com/arfan21/backend-test/internal/entity"
)

type Repository interface {
	Get(ctx context.Context, shortURL string) (res entity.URLShortener, err error)
	Set(ctx context.Context, value entity.URLShortener) (err error)
}
