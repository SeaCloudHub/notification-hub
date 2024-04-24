package subcriber

import (
	"context"
	"fmt"

	pubsub "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub"
	"github.com/SeaCloudHub/notification-hub/adapters/skio"
	noti "github.com/SeaCloudHub/notification-hub/domain/notification"
	"go.uber.org/zap"
)

func NotifyWhenUserCall(rtEngine skio.RealtimeEngine, logger *zap.SugaredLogger, store noti.Store) consumerJob {
	return consumerJob{
		Title: "Notification for user",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			fmt.Printf("day ne: %v\n", message.Data())
			setNotifications := message.Data().(noti.SetOfNotifications)
			fmt.Printf("Set of notifications: %v\n", setNotifications)
			for _, notification := range setNotifications.Noitications {
				fmt.Print("notication tung cai", notification)
				fmt.Printf("content day ne: %v\n", notification.Content)
				err := rtEngine.EmitToUser(notification.To, "notification", notification.Content)
				if err != nil {
					// fmt.Printf("loi content day ne: %v\n", notification.Content)
					if err := store.UpdateStatusByUid(ctx, notification.Uid, noti.StatusFailure); err != nil {
						logger.Errorf("Error updating notification status failure: %v", err)
					}

				}
				if err := store.UpdateStatusByUid(ctx, notification.Uid, noti.StatusSuccess); err != nil {
					logger.Errorf("Error updating notification status success: %v", err)
				}

			}

			return nil
		},
	}
}
