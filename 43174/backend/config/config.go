package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	ServerMode     string
	DatabaseURL    string
	JWTSecret      string
	JWTExpireHours int
	UploadPath     string
	MaxUploadSize  int64
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		ServerMode:     getEnv("SERVER_MODE", "release"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/campus_trade?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 24),
		UploadPath:     getEnv("UPLOAD_PATH", "./uploads"),
		MaxUploadSize:  getEnvInt64("MAX_UPLOAD_SIZE", 10<<20),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		return parseInt(value, fallback)
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		return parseInt64(value, fallback)
	}
	return fallback
}
