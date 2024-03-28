package main

import (
	"log"

	"github.com/SeaCloudHub/notification-hub/pkg/config"
	"github.com/SeaCloudHub/notification-hub/pkg/logger"
	"github.com/SeaCloudHub/notification-hub/pkg/sentry"
	sentrygo "github.com/getsentry/sentry-go"
	_ "github.com/lib/pq"
)

func main() {
	applog, err := logger.NewAppLogger()
	if err != nil {
		log.Fatalf("cannot load config: %v\n", err)
	}
	// defer logger.Sync(applog)

	cfg, err := config.LoadConfig()
	if err != nil {
		applog.Fatal(err)
	}

	err = sentrygo.Init(sentrygo.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.AppEnv,
		AttachStacktrace: true,
	})
	if err != nil {
		applog.Fatalf("cannot init sentry: %v", err)
	}
	defer sentrygo.Flush(sentry.FlushTime)

	email := "admin@seacloudhub.com"
	password := "plzdonthackme"

	applog.Info("admin user created successfully")
	applog.Infof("email: %s - password: %s", email, password)
}
