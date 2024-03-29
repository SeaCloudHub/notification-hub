package engine

import (
	"context"
	"testing"

	"github.com/SeaCloudHub/notification-hub/adapters/pubsub"
	redisstore "github.com/SeaCloudHub/notification-hub/adapters/redis_store"
	"github.com/SeaCloudHub/notification-hub/pkg/config"
)

const DefaultAddress = "127.0.0.1:6379"
const DefaultDb = 1

func TestPublish(t *testing.T) {
	ctx := context.TODO()
	cfg := config.Config{}
	cfg.Redis.Addr = DefaultAddress
	cfg.Redis.Db = DefaultDb

	flag := make(chan bool)

	redisStorage, err := redisstore.NewRedisStorage(&cfg)
	if err != nil {
		t.Error(err)
	}

	redisStoragePing := redisStorage.Store.Ping(ctx)
	if redisStoragePing.Err() != nil {
		t.Error(err)
	}

	redisPubsub := NewRedisPubsub(redisStorage, nil)

	message := &pubsub.Message{}
	message.SetData("test")

	msg, _ := redisPubsub.Subscribe(ctx, "test")

	if err := redisPubsub.Publish(ctx, "test", message); err != nil {
		t.Error(err)
	}

	go func() {
		for {
			c := <-msg
			if c.Data() != "test" {
				t.Error()
			}

			flag <- true
		}
	}()

	if result := <-flag; result {
		t.Log("ok")
	}

}
