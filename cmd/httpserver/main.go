package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SeaCloudHub/notification-hub/adapters/httpserver"
	"github.com/SeaCloudHub/notification-hub/adapters/postgrestore"
	"github.com/SeaCloudHub/notification-hub/adapters/services"
	"github.com/SeaCloudHub/notification-hub/pkg/config"
	"github.com/SeaCloudHub/notification-hub/pkg/logger"
	"github.com/SeaCloudHub/notification-hub/pkg/sentry"
	sentrygo "github.com/getsentry/sentry-go"
)

func main() {
	applog, err := logger.NewAppLogger()
	if err != nil {
		log.Fatalf("cannot load config: %v\n", err)
	}
	defer logger.Sync(applog)

	cfg, err := config.LoadConfig()
	if err != nil {
		applog.Fatal(err)
	}

	applog.Infof("env: %v", cfg)

	err = sentrygo.Init(sentrygo.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.AppEnv,
		AttachStacktrace: true,
	})
	if err != nil {
		applog.Fatalf("cannot init sentry: %v", err)
	}
	defer sentrygo.Flush(sentry.FlushTime)

	db, err := postgrestore.NewConnection(postgrestore.ParseFromConfig(cfg))
	if err != nil {
		applog.Fatal(err)
	}

	server, err := httpserver.New(cfg, applog)
	server.NotificationStore = postgrestore.NewNotificationStore(db)
	server.IdentityService = services.NewIdentityService(cfg)

	if err != nil {
		applog.Fatal(err)
	}

	if err := server.SetupEngineForPubsubAndSocket(); err != nil {
		applog.Fatal(err)
	}

	fmt.Print(server)

	addr := fmt.Sprintf(":%s", cfg.Port)
	applog.Info("server started!")
	applog.Fatal(http.ListenAndServe(addr, server))
}
