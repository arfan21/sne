package dbredis

import (
	"context"

	"github.com/arfan21/backend-test/config"
	"github.com/arfan21/backend-test/pkg/logger"
	"github.com/gofiber/storage/redis/v3"
)

func New() (*redis.Storage, error) {
	client := redis.New(redis.Config{
		Host:     config.GetConfig().Redis.Host,
		Port:     config.GetConfig().Redis.Port,
		Username: config.GetConfig().Redis.Username,
		Password: config.GetConfig().Redis.Password,
	})

	err := client.Conn().Ping(context.Background()).Err()
	if err != nil {
		logger.Log(context.Background()).Error().Err(err).Msg("failed to ping redis")
		return nil, err
	}

	logger.Log(context.Background()).Info().Msg("dbredis: connection established")

	return client, nil
}
