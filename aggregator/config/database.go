package config

import (
	"sso-dashboard.bcgov.com/aggregator/utils"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

var (
	databaseConfig DatabaseConfig
)

func init() {
	hostname := utils.GetEnv("DB_HOSTNAME", "localhost")
	port := utils.GetEnv("DB_PORT", "5432")
	database := utils.GetEnv("DB_DATABASE", "aggregation")
	username := utils.GetEnv("DB_USERNAME", "")
	password := utils.GetEnv("DB_PASSWORD", "")

	databaseConfig.Host = hostname
	databaseConfig.Port = port
	databaseConfig.DBName = database
	databaseConfig.User = username
	databaseConfig.Password = password
}

func LoadDatabaseConfig() *DatabaseConfig {
	return &databaseConfig
}
