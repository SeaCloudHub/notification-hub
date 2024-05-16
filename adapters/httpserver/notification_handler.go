package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SeaCloudHub/notification-hub/adapters/httpserver/model"
	realtimePubsub "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub"
	"github.com/SeaCloudHub/notification-hub/adapters/subcriber"
	"github.com/SeaCloudHub/notification-hub/domain/identity"
	"github.com/SeaCloudHub/notification-hub/domain/notification"

	"github.com/SeaCloudHub/notification-hub/pkg/mycontext"
	"github.com/SeaCloudHub/notification-hub/pkg/pagination"
	"github.com/labstack/echo/v4"
)

func (s *Server) ListEntries(c echo.Context) error {
	var (
		ctx = mycontext.NewEchoContextAdapter(c)
		req model.ListEntriesRequest
	)

	if err := c.Bind(&req); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	if err := req.Validate(ctx); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	identity, _ := c.Get(ContextKeyIdentity).(*identity.Identity)

	cursor := pagination.NewCursor(req.Cursor, req.Limit)

	notifications, err := s.NotificationStore.ListByUserIdUsingCursor(ctx, identity.ID, cursor)

	if err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	return s.success(c, model.ListEntriesResponse{
		Entries: notifications,
		Cursor:  cursor.NextToken(),
	})

}

func (s *Server) ListPageEntries(c echo.Context) error {
	var (
		ctx = mycontext.NewEchoContextAdapter(c)
		req model.ListPageEntriesRequest
	)

	if err := c.Bind(&req); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	if err := req.Validate(ctx); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	identity, _ := c.Get(ContextKeyIdentity).(string)
	if identity == "" {
		return s.handleError(c, errors.New("identity is nil"), http.StatusNonAuthoritativeInfo)
	}
	fmt.Print("user: ", identity)
	pager := pagination.NewPager(req.Page, req.Limit)
	notifications, err := s.NotificationStore.ListByUserIdUsingPaper(ctx, identity, pager)

	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	return s.success(c, model.ListPageEntriesResponse{
		Entries:    notifications,
		Pagination: pager.PageInfo(),
	})

}

func (s *Server) UpdateViewedTime(c echo.Context) error {
	var (
		ctx = mycontext.NewEchoContextAdapter(c)
		req model.UpdateViewedTimeRequest
	)

	if err := c.Bind(&req); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	identity, _ := c.Get(ContextKeyIdentity).(string)
	if identity == "" {
		return s.handleError(c, errors.New("identity is nil"), http.StatusNonAuthoritativeInfo)
	}

	// check exist
	count, err := s.NotificationStore.CheckExistToUpdateViewedTimeAndStatus(ctx, req.IdNotification, identity)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	if count == 0 {
		return s.handleError(c, errors.New("can not find suitable record to update"), http.StatusBadRequest)
	}

	if err := s.NotificationStore.UpdateViewedTimeAndStatus(ctx, req.IdNotification, identity, time.Now()); err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	return s.success(c, model.NotificationResponse{Status: "updated"})

}
func (s *Server) PushNotification(c echo.Context) error {
	var (
		ctx = mycontext.NewEchoContextAdapter(c)
		req model.NotificationRequest
	)

	if err := c.Bind(&req); err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	var setOfNotification notification.SetOfNotifications

	// init notifications
	for _, noti := range req.Notifications {
		uid := fmt.Sprintf("%v.%v", time.Now(), noti.UserId)
		notiEntity := notification.NewNotification(uid, req.From, noti.UserId, noti.Content)
		err := s.NotificationStore.Create(ctx, &notiEntity)
		if err != nil {
			s.Logger.Errorf("Error creating notification with uid %v: %v", uid, err)
		}
		setOfNotification.Noitications = append(setOfNotification.Noitications, &notiEntity)
	}

	message := realtimePubsub.NewMessage(setOfNotification)

	s.pubsub.Publish(ctx, subcriber.UserNotificationChannel, message)

	return s.success(c, model.NotificationResponse{Status: "processing"})
}

func (s *Server) RegisterNotificationRoutes(router *echo.Group) {
	router.POST("/notifications", s.PushNotification)
}
