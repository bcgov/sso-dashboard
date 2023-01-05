package model

import (
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"sso-dashboard.bcgov.com/aggregator/config"
)

var (
	pgdb *pg.DB
)

func init() {
	pgdb = connect()
}

func connect() *pg.DB {
	return pg.Connect(pgOptions())
}

func pgOptions() *pg.Options {
	cfg := config.LoadDatabaseConfig()

	return &pg.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.DBName,

		MaxRetries:      3,
		MinRetryBackoff: -1,

		DialTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		PoolSize:           100,
		MaxConnAge:         10 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        10 * time.Second,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

func GetDB() *pg.DB {
	return pgdb
}
