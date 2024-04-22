package subcriber

import (
	"context"

	pubsub "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub"
	"github.com/SeaCloudHub/notification-hub/adapters/skio"
	"go.uber.org/zap"
)

const (
	UserNotificationChannel = pubsub.Topic("notification-channel")
)

type consumerJob struct {
	Title string
	Hld   func(ctx context.Context, message *pubsub.Message) error
}

type consumerEngine struct {
	pubsub   pubsub.Pubsub
	rtEngine skio.RealtimeEngine
	logger   *zap.SugaredLogger
}

func NewEngine(pubsub pubsub.Pubsub, rtEngine skio.RealtimeEngine) *consumerEngine {
	return &consumerEngine{
		pubsub:   pubsub,
		rtEngine: rtEngine,
	}
}

func (engine *consumerEngine) startSubTopic(topic pubsub.Topic, consuconsumerJobs ...consumerJob) error {
	ctx := context.Background()
	c, _ := engine.pubsub.Subscribe(ctx, topic)

	go func() {
		for _, item := range consuconsumerJobs {
			message := <-c
			go item.Hld(ctx, message)
		}
	}()

	return nil
}

func (engine *consumerEngine) Start() error {
	engine.startSubTopic(UserNotificationChannel, NotifyWhenUserCall(engine.rtEngine, engine.logger))
	return nil
}
