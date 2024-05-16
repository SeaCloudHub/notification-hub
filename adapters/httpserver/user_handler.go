package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) Me(c echo.Context) error {
	return c.JSON(http.StatusOK, c.Get(ContextKeyIdentity))
}

func (s *Server) RegisterUserRoutes(router *echo.Group) {
	router.GET("/notifications", s.ListPageEntries)
	router.PATCH("/notifications", s.UpdateViewedTime)
	router.PATCH("/marking-viewed-notifications", s.MarkAllAsView)
}
