package main

import (
	"log"
	"strconv"

	"github.com/SeaCloudHub/notification-hub/adapters/postgrestore"
	"github.com/SeaCloudHub/notification-hub/pkg/config"
	"github.com/SeaCloudHub/notification-hub/pkg/logger"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	applogger, err := logger.NewAppLogger()
	// defer logger.Sync(applogger)
	if err != nil {
		log.Fatalf("cannot load config: %v\n", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		applogger.Fatalf("cannot load config: %v\n", err)
	}

	db, err := postgrestore.NewConnection(postgrestore.Options{
		DBName:   cfg.DB.Name,
		DBUser:   cfg.DB.User,
		Password: cfg.DB.Pass,
		Host:     cfg.DB.Host,
		Port:     strconv.Itoa(cfg.DB.Port),
		SSLMode:  false,
	})
	if err != nil {
		applogger.Fatalf("cannot connecting to db: %v\n", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	total, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		applogger.Fatalf("cannot execute migration: %v\n", err)
	}

	applogger.Infof("applied %d migrations\n", total)
}
