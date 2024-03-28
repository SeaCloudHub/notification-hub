package redisstore

import (
	"context"
	"errors"

	"github.com/SeaCloudHub/notification-hub/pkg/config"
	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	Store *redis.Client
}

func NewRedisStorage(cfg *config.Config) (*RedisStorage, error) {
	if cfg.Redis.Addr == "." {
		return nil, errors.New("empty redis address")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		DB:       cfg.Redis.Db,
		Password: cfg.Redis.Pass,
	})

	_, err := redisClient.Ping(context.TODO()).Result()
	if err != nil {
		return &RedisStorage{
			Store: redisClient,
		}, errors.New("ping error: " + err.Error())
	}

	return &RedisStorage{
		Store: redisClient,
	}, nil
}
