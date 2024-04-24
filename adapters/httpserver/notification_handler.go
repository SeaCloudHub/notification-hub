package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SeaCloudHub/notification-hub/adapters/httpserver/model"
	realtimePubsub "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub"
	"github.com/SeaCloudHub/notification-hub/adapters/subcriber"
	"github.com/SeaCloudHub/notification-hub/domain/notification"
	"github.com/SeaCloudHub/notification-hub/pkg/mycontext"
	"github.com/labstack/echo/v4"
)

func (s *Server) PushNotification(c echo.Context) error {
	var (
		ctx = mycontext.NewEchoContextAdapter(c)
		req model.NotificationRequest
	)

	if err := c.Bind(&req); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	var setOfNotification notification.SetOfNotifications

	// init notification
	for _, noti := range req.Notifications {
		uid := fmt.Sprintf("%v.%v", time.Now(), noti.UserId)
		notiEntity := notification.NewNotification(uid, req.From, noti.UserId, noti.Content)
		setOfNotification.Noitications = append(setOfNotification.Noitications, &notiEntity)
	}

	s.pubsub.Publish(ctx, subcriber.UserNotificationChannel, realtimePubsub.NewMessage(setOfNotification))

	return s.success(c, model.NotificationResponse{Status: "processing"})
}

func (s *Server) RegisterNotificationRoutes(router *echo.Group) {
	router.POST("/user", s.PushNotification)
}
