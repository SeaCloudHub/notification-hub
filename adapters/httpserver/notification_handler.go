package httpserver

import (
	"net/http"

	"github.com/SeaCloudHub/notification-hub/adapters/httpserver/model"
	realtimePubsub "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub"
	"github.com/SeaCloudHub/notification-hub/adapters/subcriber"
	"github.com/SeaCloudHub/notification-hub/pkg/mycontext"
	"github.com/labstack/echo/v4"
)

func (s *Server) Notification(c echo.Context) error {
	var (
		ctx = mycontext.NewEchoContextAdapter(c)
		req model.NotificationRequest
	)

	if err := c.Bind(&req); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	s.pubsub.Publish(ctx, subcriber.UserNotificationChannel, realtimePubsub.NewMessage(req))

	return s.success(c, model.NotificationResponse{Status: "processing"})
}
