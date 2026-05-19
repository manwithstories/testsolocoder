package config

import (
	"os"

	"github.com/joho/godotenv"
)

var JWTSecret string

func InitConfig() {
	_ = godotenv.Load()
	JWTSecret = getEnv("JWT_SECRET", "default_secret_key_change_in_production")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
