package redisstore

import (
	"context"
	"testing"

	"github.com/SeaCloudHub/notification-hub/pkg/config"
)

const DefaultAddress = "127.0.0.1:6379"
const DefaultDb = 1

func TestConnectRedis(t *testing.T) {
	ctx := context.TODO()
	cfg := config.Config{}
	cfg.Redis.Addr = DefaultAddress
	cfg.Redis.Db = DefaultDb

	redisStorage, err := NewRedisStorage(&cfg)
	if err != nil {
		t.Error(err)
	}

	result := redisStorage.Store.Ping(ctx)
	if result.Err() != nil {
		t.Error(err)
	}

}
