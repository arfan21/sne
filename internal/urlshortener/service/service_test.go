package urlshortenersvc

import (
	"context"
	"errors"
	"testing"

	"github.com/arfan21/backend-test/config"
	"github.com/arfan21/backend-test/internal/entity"
	"github.com/arfan21/backend-test/internal/model"
	repoMock "github.com/arfan21/backend-test/mocks/internal_/urlshortener"
	"github.com/arfan21/backend-test/pkg/constant"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Create(t *testing.T) {
	repoRedis := repoMock.NewRepository(t)
	repoPg := repoMock.NewRepository(t)

	svc := New(repoRedis, repoPg, func(s *Service) {
		s.opt.generatorShorted = func(length int) string {
			return "shortedUrl"
		}
	})

	testCases := []struct {
		name          string
		input         model.CreateURLShortedRequest
		setupMocks    func()
		expectedError error
		expectedRes   model.CreateURLShortedResponse
	}{
		{
			name: "Success case",
			input: model.CreateURLShortedRequest{
				LongUrl: "http://example.com",
			},
			setupMocks: func() {
				repoPg.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, constant.ErrURLNotFound).Once()
				repoPg.On("Set", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedError: nil,
			expectedRes: model.CreateURLShortedResponse{
				ShortUrl: config.GetConfig().Service.BaseURL + "/shortedUrl",
			},
		},
		{
			name: "Validation failure",
			input: model.CreateURLShortedRequest{
				LongUrl: "",
			},
			setupMocks:    func() {},
			expectedError: errors.New("urlshorted.service.Create: failed to validate"),
		},
		{
			name: "Repository Get failure",
			input: model.CreateURLShortedRequest{
				LongUrl: "http://example.com",
			},
			setupMocks: func() {
				repoPg.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, errors.New("db error")).Once()
			},
			expectedError: errors.New("urlshorted.service.Create: failed to get url shorted db error"),
		},
		{
			name: "Repository Set failure",
			input: model.CreateURLShortedRequest{
				LongUrl: "http://example.com",
			},
			setupMocks: func() {
				repoPg.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, constant.ErrURLNotFound).Once()
				repoPg.On("Set", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
			},
			expectedError: errors.New("urlshorted.service.Create: failed to set url shorted db error"),
		},
		{
			name: "Duplicate shorted url",
			input: model.CreateURLShortedRequest{
				LongUrl: "http://example.com",
			},
			setupMocks: func() {
				repoPg.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, nil).Once()
				repoPg.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, constant.ErrURLNotFound).Once()
				repoPg.On("Set", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedError: nil,
			expectedRes: model.CreateURLShortedResponse{
				ShortUrl: config.GetConfig().Service.BaseURL + "/shortedUrl",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			res, err := svc.Create(context.Background(), tc.input)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}

			repoPg.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	repoRedis := repoMock.NewRepository(t)
	repoPg := repoMock.NewRepository(t)

	svc := New(repoRedis, repoPg)

	testCases := []struct {
		name          string
		shortedUrl    string
		setupMocks    func()
		expectedError error
		expectedRes   string
	}{
		{
			name:       "Success - Found in Redis",
			shortedUrl: "shortedUrl",
			setupMocks: func() {
				repoRedis.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{LongURL: "http://example.com"}, nil).Once()
			},
			expectedError: nil,
			expectedRes:   "http://example.com",
		},
		{
			name:       "Success - Found in Postgres and set to Redis",
			shortedUrl: "shortedUrl",
			setupMocks: func() {
				repoRedis.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, redis.Nil).Once()
				repoPg.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{LongURL: "http://example.com"}, nil).Once()
				repoRedis.On("Set", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedError: nil,
			expectedRes:   "http://example.com",
		},
		{
			name:       "Failure - Redis get error",
			shortedUrl: "shortedUrl",
			setupMocks: func() {
				repoRedis.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, errors.New("redis error")).Once()
			},
			expectedError: errors.New("urlshorted.service.Get: failed to get url shorted from redis"),
		},
		{
			name:       "Failure - Postgres get error",
			shortedUrl: "shortedUrl",
			setupMocks: func() {
				repoRedis.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, redis.Nil).Once()
				repoPg.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, errors.New("postgres error")).Once()
			},
			expectedError: errors.New("urlshorted.service.Get: failed to get url shorted from pg"),
		},
		{
			name:       "Failure - Redis set error",
			shortedUrl: "shortedUrl",
			setupMocks: func() {
				repoRedis.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{}, redis.Nil).Once()
				repoPg.On("Get", mock.Anything, "shortedUrl").Return(entity.URLShortener{LongURL: "http://example.com"}, nil).Once()
				repoRedis.On("Set", mock.Anything, mock.Anything).Return(errors.New("redis set error")).Once()
			},
			expectedError: errors.New("urlshorted.service.Get: failed to set url shorted"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			res, err := svc.Get(context.Background(), tc.shortedUrl)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}

			repoRedis.AssertExpectations(t)
			repoPg.AssertExpectations(t)
		})
	}
}
