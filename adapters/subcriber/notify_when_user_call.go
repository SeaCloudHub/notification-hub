package subcriber

import (
	"context"
	"fmt"

	pubsub "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub"
	"github.com/SeaCloudHub/notification-hub/adapters/skio"
	"github.com/SeaCloudHub/notification-hub/domain/notification"
	"go.uber.org/zap"
)

func NotifyWhenUserCall(rtEngine skio.RealtimeEngine, logger *zap.SugaredLogger) consumerJob {
	return consumerJob{
		Title: "Notification for user",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			fmt.Printf("day ne: %v\n", message.Data())
			setNotifications := message.Data().(notification.SetOfNotifications)
			fmt.Printf("Set of notifications: %v\n", setNotifications)
			for _, notification := range setNotifications.Noitications {
				fmt.Print("notication tung cai", notification)
				fmt.Printf("content day ne: %v\n", notification.Content)
				err := rtEngine.EmitToUser(notification.To, "notification", notification.Content)
				if err != nil {
					// fmt.Printf("loi content day ne: %v\n", notification.Content)

					// do some thing
					logger.Errorw(
						err.Error(),
						zap.String("message", message.Data().(string)),
					)
				} else {

				}
				// TODO: update status of notification

			}

			return nil
		},
	}
}
