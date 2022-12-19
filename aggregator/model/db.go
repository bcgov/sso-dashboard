package model

import (
	"fmt"

	"github.com/go-pg/pg"
	"sso-dashboard.bcgov.com/aggregator/config"
)

// ConnectDB is
func ConnectDB() *pg.DB {
	cfg := config.LoadDatabaseConfig()

	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.DBName,
	})

	return db
}
