package subcriber

import (
	"context"

	"github.com/SeaCloudHub/notification-hub/adapters/pubsub"
	"github.com/SeaCloudHub/notification-hub/adapters/skio"
)

type consumerJob struct {
	Title string
	Hld   func(ctx context.Context, message *pubsub.Message) error
}

type consumerEngine struct {
	pubsub pubsub.Pubsub
	skio   skio.AppSocket
}

func (engine *consumerEngine) startSubTopic(topic pubsub.Topic, consuconsumerJobs ...consumerJob) error {
	ctx := context.Background()
	c, _ := engine.pubsub.Subscribe(ctx, topic)

	for _, item := range consuconsumerJobs {
		message := <-c
		go item.Hld(ctx, message)
	}

	return nil
}
