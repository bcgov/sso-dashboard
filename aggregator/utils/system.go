package utils

import (
	"github.com/joho/godotenv"
	"os"
)

func init() {
	godotenv.Load(".env")
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
