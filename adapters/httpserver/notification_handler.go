package httpserver

import (
	"errors"
	"net/http"

	"github.com/SeaCloudHub/notification-hub/adapters/httpserver/model"
	"github.com/SeaCloudHub/notification-hub/domain/identity"
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

	session, err := s.IdentityService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, identity.ErrInvalidCredentials) {
			return s.handleError(c, err, http.StatusBadRequest)
		}

		return s.handleError(c, err, http.StatusInternalServerError)
	}

	return s.success(c, model.LoginResponse{SessionToken: *session.Token})
}
