package httpserver

import (
	"net/http"
	"strings"

	realtimePubsub "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub"
	localEngine "github.com/SeaCloudHub/notification-hub/adapters/realtime_pubsub/local_engine"
	"github.com/SeaCloudHub/notification-hub/adapters/subcriber"

	"github.com/SeaCloudHub/notification-hub/adapters/skio"
	"github.com/SeaCloudHub/notification-hub/domain/identity"
	"github.com/SeaCloudHub/notification-hub/domain/notification"
	"github.com/SeaCloudHub/notification-hub/domain/permission"
	"github.com/SeaCloudHub/notification-hub/pkg/config"
	"github.com/SeaCloudHub/notification-hub/pkg/sentry"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Options func(s *Server) error

type Server struct {
	router *echo.Echo
	Config *config.Config
	Logger *zap.SugaredLogger

	//storage adapters
	NotificationStore notification.Store

	// pubsub
	pubsub realtimePubsub.Pubsub
	// services
	IdentityService   identity.Service
	PermissionService permission.Service
}

func New(cfg *config.Config, logger *zap.SugaredLogger, options ...Options) (*Server, error) {
	s := Server{
		router: echo.New(),
		Config: cfg,
		Logger: logger,
	}

	for _, fn := range options {
		if err := fn(&s); err != nil {
			return nil, err
		}
	}

	s.RegisterGlobalMiddlewares()
	s.RegisterHealthCheck(s.router.Group(""))
	s.RegisterWebSocket(s.router.Group(""))

	// authMiddleware := s.NewAuthentication("header:Authorization", "Bearer",
	// 	[]string{
	// 		"/healthz",
	// 		"/api/users/login",
	// 		"/demo",
	// 	},
	// ).Middleware()

	// s.router.Use(authMiddleware)

	s.RegisterUserRoutes(s.router.Group("/api/users"))
	s.RegisterNotificationRoutes(s.router.Group("api/internal"))

	return &s, nil
}

func (s *Server) SetupEngineForPubsubAndSocket() error {
	skioEngine := skio.NewEngine()
	s.pubsub = localEngine.NewPubSub()

	if err := skioEngine.Run(s.router, s.IdentityService); err != nil {
		return err
	}

	if err := subcriber.NewEngine(s.pubsub, skioEngine, s.Logger, s.NotificationStore).Start(); err != nil {
		return err
	}

	return nil

}
func (s *Server) RegisterGlobalMiddlewares() {
	s.router.Use(middleware.Recover())
	s.router.Use(middleware.Secure())
	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Gzip())
	s.router.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	// CORS
	if s.Config.AllowOrigins != "" {
		aos := strings.Split(s.Config.AllowOrigins, ",")
		s.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: aos,
		}))
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) RegisterHealthCheck(router *echo.Group) {
	router.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK!!!")
	})
}

func (s *Server) RegisterWebSocket(router *echo.Group) {
	// router.GET("/demo", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "OK!!!")
	// })
	router.File("/demo", "./demo.html")
}

func (s *Server) handleError(c echo.Context, err error, status int) error {
	s.Logger.Errorw(
		err.Error(),
		zap.String("request_id", s.requestID(c)),
	)

	if status >= http.StatusInternalServerError {
		sentry.WithContext(c).Error(err)
	}

	return c.JSON(status, map[string]string{
		"message": http.StatusText(status),
		"info":    err.Error(),
	})
}

func (s *Server) success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    data,
	})
}

func (s *Server) requestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}
