package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}
	ensureEnv("JWT_SECRET", "dev-secret-change-me")
	ensureEnv("JWT_EXPIRE_HOURS", "72")
	ensureEnv("SERVER_PORT", "8080")
}

func ensureEnv(key, def string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, def)
	}
}
