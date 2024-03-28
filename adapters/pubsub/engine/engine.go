package engine

import (
	"context"
	"fmt"

	"github.com/SeaCloudHub/notification-hub/adapters/pubsub"
	redisstore "github.com/SeaCloudHub/notification-hub/adapters/redis_store"
	"github.com/SeaCloudHub/notification-hub/domain/notification"
)

type redisPubsub struct {
	engineRedis redisstore.RedisStorage
	notiStorage notification.Storage
}

func (ps *redisPubsub) Publish(ctx context.Context, topic pubsub.Topic, data *pubsub.Message) error {
	pub := ps.engineRedis.Store.Publish(ctx, string(topic), data)
	return pub.Err()
}

func (ps *redisPubsub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	c := make(chan *pubsub.Message)
	sub := ps.engineRedis.Store.Subscribe(ctx, string(topic))

	go func() {

		for {
			message, err := sub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Println(err)
				break
			}

			msg := &pubsub.Message{}
			msg.SetChannel(pubsub.Topic(message.Channel))
			msg.SetData([]byte(message.Payload))

			c <- msg
		}
	}()

	return c, func() {
		sub.Close()
	}
}
