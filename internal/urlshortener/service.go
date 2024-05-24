package urlshortener

import (
	"context"

	"github.com/arfan21/backend-test/internal/model"
)

type Service interface {
	Create(ctx context.Context, value model.CreateURLShortedRequest) (res model.CreateURLShortedResponse, err error)
	Get(ctx context.Context, key string) (res string, err error)
}
