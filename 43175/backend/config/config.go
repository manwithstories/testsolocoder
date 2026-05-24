package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	AppPort         string
	AppSecret       string
	DBDriver        string
	DBPath          string
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	RedisDB         int
	JWTSecret       string
	JWTExpireHour   int
	EmailSMTPHost   string
	EmailSMTPPort   int
	EmailFrom       string
	EmailPassword   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}

	return &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		AppSecret:     getEnv("APP_SECRET", "default-secret-key"),
		DBDriver:      getEnv("DB_DRIVER", "sqlite"),
		DBPath:        getEnv("DB_PATH", "./smart_energy.db"),
		RedisHost:     getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		JWTSecret:     getEnv("JWT_SECRET", "jwt-secret-key"),
		JWTExpireHour: getEnvInt("JWT_EXPIRE_HOUR", 24),
		EmailSMTPHost: getEnv("EMAIL_SMTP_HOST", ""),
		EmailSMTPPort: getEnvInt("EMAIL_SMTP_PORT", 587),
		EmailFrom:     getEnv("EMAIL_FROM", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
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
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
