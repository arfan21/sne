package urlshortenersvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arfan21/backend-test/config"
	"github.com/arfan21/backend-test/internal/entity"
	"github.com/arfan21/backend-test/internal/model"
	"github.com/arfan21/backend-test/internal/urlshortener"
	"github.com/arfan21/backend-test/pkg/constant"
	"github.com/arfan21/backend-test/pkg/validation"
	"github.com/jaevor/go-nanoid"
	"github.com/redis/go-redis/v9"
)

type GeneratorShorted func(len int) string

type Option struct {
	generatorShorted GeneratorShorted
}

type OptionFunc func(s *Service)

var defaultOption = func() Option {
	return Option{
		generatorShorted: func(len int) string {
			nanoId, _ := nanoid.Standard(len)
			return nanoId()
		},
	}
}

func WithGeneratorShortedURL(generator GeneratorShorted) OptionFunc {
	return func(s *Service) {
		s.opt.generatorShorted = generator
	}
}

type Service struct {
	repoRedis urlshortener.Repository
	repoPg    urlshortener.Repository

	opt Option
}

func New(repoRedis urlshortener.Repository, repoPg urlshortener.Repository, opts ...OptionFunc) *Service {
	opt := defaultOption()
	svc := &Service{
		repoRedis: repoRedis,
		repoPg:    repoPg,
		opt:       opt,
	}

	for _, oFn := range opts {
		oFn(svc)
	}

	return svc
}

func (s Service) Create(ctx context.Context, value model.CreateURLShortedRequest) (res model.CreateURLShortedResponse, err error) {
	err = validation.Validate(value)
	if err != nil {
		err = fmt.Errorf("urlshorted.service.Create: failed to validate %w", err)
		return
	}

	charLength := config.GetConfig().Service.ShortedURLCharLength
	generatedShorted := s.opt.generatorShorted(charLength)

	// check if shorted url already exist
	_, err = s.repoPg.Get(ctx, generatedShorted)
	if err != nil && !errors.Is(err, constant.ErrURLNotFound) {
		err = fmt.Errorf("urlshorted.service.Create: failed to get url shorted %w", err)
		return
	}

	if err == nil {
		// if shorted url already exist, generate new shorted url
		maxAttempt := config.GetConfig().Service.MaxAttemptGenerateShortedURL
		totalAttempt := 0
		// handle collision
		// loop until generated shorted url not exist in database
		for ; totalAttempt < maxAttempt; totalAttempt++ {
			generatedShorted = s.opt.generatorShorted(charLength + totalAttempt)

			_, err = s.repoPg.Get(ctx, generatedShorted)
			if err != nil && errors.Is(err, constant.ErrURLNotFound) {
				break
			}
		}

	}

	ttl := time.Duration(config.GetConfig().Service.ShortedURLTTL) * time.Hour
	expiredAtUnix := time.Now().Add(ttl).Unix()

	urlShorted := entity.URLShortener{
		LongURL:   value.LongUrl,
		ShortURL:  generatedShorted,
		ExpiredAt: expiredAtUnix,
	}

	err = s.repoPg.Set(ctx, urlShorted)
	if err != nil {
		err = fmt.Errorf("urlshorted.service.Create: failed to set url shorted %w", err)
		return
	}

	res = model.CreateURLShortedResponse{
		ShortUrl: config.GetConfig().Service.BaseURL + "/" + urlShorted.ShortURL,
	}

	return
}

func (s Service) Get(ctx context.Context, shortedUrl string) (res string, err error) {
	// get from redis
	data, err := s.repoRedis.Get(ctx, shortedUrl)
	if err != nil && !errors.Is(err, redis.Nil) {
		err = fmt.Errorf("urlshorted.service.Get: failed to get url shorted from redis %w", err)
		return
	}

	if errors.Is(err, redis.Nil) {
		// get from postgres
		data, err = s.repoPg.Get(ctx, shortedUrl)
		if err != nil {
			err = fmt.Errorf("urlshorted.service.Get: failed to get url shorted from pg %w", err)
			return
		}

		// set to redis
		err = s.repoRedis.Set(ctx, data)
		if err != nil {
			err = fmt.Errorf("urlshorted.service.Get: failed to set url shorted %w", err)
			return
		}
	}

	now := time.Now().Unix()
	if data.ExpiredAt < now {
		err = constant.ErrURLNotFound
		return
	}

	res = data.LongURL

	return
}
