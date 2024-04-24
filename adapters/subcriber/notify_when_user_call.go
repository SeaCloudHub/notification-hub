package subcriber

import (
	"context"
	"fmt"

	"github.com/SeaCloudHub/notification-hub/adapters/httpserver/model"
	pubsub "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub"
	"github.com/SeaCloudHub/notification-hub/adapters/skio"
	"go.uber.org/zap"
)

func NotifyWhenUserCall(rtEngine skio.RealtimeEngine, logger *zap.SugaredLogger) consumerJob {
	return consumerJob{
		Title: "Notification for user",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			fmt.Printf("day ne: %v\n", message.Data())
			notifications := message.Data().(model.NotificationRequest)
			for _, notification := range notifications.Notifications {
				fmt.Printf("content day ne: %v\n", notification.Content)
				err := rtEngine.EmitToUser(notification.UserId, "notification", notification.Content)
				if err != nil {
					// fmt.Printf("loi content day ne: %v\n", notification.Content)

					// do some thing
					logger.Errorw(
						err.Error(),
						zap.String("message", message.Data().(string)),
					)
				}
				// TODO: update status of notification

			}

			return nil
		},
	}
}
