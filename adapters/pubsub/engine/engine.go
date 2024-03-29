package engine

import (
	"context"
	"encoding/json"

	"github.com/SeaCloudHub/notification-hub/adapters/pubsub"
	redisstore "github.com/SeaCloudHub/notification-hub/adapters/redis_store"
	"go.uber.org/zap"
)

type redisPubsub struct {
	engineRedis *redisstore.RedisStorage
	logger      *zap.SugaredLogger
}

func NewRedisPubsub(engine *redisstore.RedisStorage, l *zap.SugaredLogger) *redisPubsub {
	return &redisPubsub{
		engineRedis: engine,
		logger:      l,
	}
}

func (ps *redisPubsub) Publish(ctx context.Context, topic pubsub.Topic, data *pubsub.Message) error {
	dataBi, err := json.Marshal(data)
	if err != nil {
		return err
	}
	pub := ps.engineRedis.Store.Publish(ctx, string(topic), dataBi)
	return pub.Err()
}

func (ps *redisPubsub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	c := make(chan *pubsub.Message)
	sub := ps.engineRedis.Store.Subscribe(ctx, string(topic))

	go func() {
		for {
			message, err := sub.ReceiveMessage(ctx)
			if err != nil {
				ps.logger.Error("received error", message)
				break
			}

			msg := &pubsub.Message{}
			msg.SetChannel(pubsub.Topic(message.Channel))
			msg.SetData(message.Payload) // SetData directly with the payload

			c <- msg
		}
	}()

	return c, func() {
		sub.Close()
	}
}
