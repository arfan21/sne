package urlshortenerrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arfan21/backend-test/config"
	"github.com/arfan21/backend-test/internal/entity"
	"github.com/arfan21/backend-test/pkg/constant"
	"github.com/gofiber/storage/redis/v3"
)

type RepositoryRedis struct {
	dbRedis *redis.Storage
}

func NewRedis(dbRedis *redis.Storage) *RepositoryRedis {
	return &RepositoryRedis{
		dbRedis: dbRedis,
	}
}

func (r *RepositoryRedis) Get(ctx context.Context, key string) (res entity.URLShortener, err error) {
	key = constant.URLShortedKey + key
	resByte, err := r.dbRedis.Conn().Get(ctx, key).Bytes()
	if err != nil {
		err = fmt.Errorf("urlshortener.repositoryRedis.Get: failed to get redis %w", err)
		return
	}

	err = json.Unmarshal(resByte, &res)
	if err != nil {
		err = fmt.Errorf("urlshortener.repositoryRedis.Get: failed to unmarshal %w", err)
		return
	}

	return
}

func (r *RepositoryRedis) Set(ctx context.Context, value entity.URLShortener) (err error) {
	key := constant.URLShortedKey + value.ShortURL
	valueByte, err := json.Marshal(value)
	if err != nil {
		err = fmt.Errorf("urlshortener.repositoryRedis.Set: failed to marshal %w", err)
		return
	}

	expired := time.Duration(config.GetConfig().Service.ShortedURLTTLCache) * time.Second

	err = r.dbRedis.Conn().Set(ctx, key, string(valueByte), expired).Err()
	if err != nil {
		err = fmt.Errorf("urlshortener.repositoryRedis.Set: failed to set redis %w", err)
		return
	}

	return
}
