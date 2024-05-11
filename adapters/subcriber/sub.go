package subcriber

import (
	"context"
	pubsub "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub"
	"github.com/SeaCloudHub/notification-hub/adapters/skio"
	"github.com/SeaCloudHub/notification-hub/domain/notification"
	"go.uber.org/zap"
)

const (
	UserNotificationChannel  = pubsub.Topic("user-notification-channel")
	AdminNotificationChannel = pubsub.Topic("admin-notification-channel")
)

type consumerJob struct {
	Title string
	Hld   func(ctx context.Context, message *pubsub.Message) error
}

type consumerEngine struct {
	pubsub    pubsub.Pubsub
	rtEngine  skio.RealtimeEngine
	logger    *zap.SugaredLogger
	notiStore notification.Store
}

func NewEngine(pubsub pubsub.Pubsub, rtEngine skio.RealtimeEngine, logger *zap.SugaredLogger, store notification.Store) *consumerEngine {
	return &consumerEngine{
		pubsub:    pubsub,
		rtEngine:  rtEngine,
		logger:    logger,
		notiStore: store,
	}
}

func (engine *consumerEngine) startSubTopic(topic pubsub.Topic, consumerJobs ...consumerJob) error {
	ctx := context.Background()
	c, _ := engine.pubsub.Subscribe(ctx, topic)

	go func() {
		for _, item := range consumerJobs {
			for {
				message := <-c
				go item.Hld(ctx, message)
			}
		}

	}()

	return nil
}

func (engine *consumerEngine) Start() error {
	engine.startSubTopic(UserNotificationChannel, NotifyWhenUserCall(engine.rtEngine, engine.logger, engine.notiStore))
	return nil
}
