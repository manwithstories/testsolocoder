package utils

import (
	"os"
	"strconv"
	"survey-platform/config"
	"sync"
)

var (
	cfg  *config.Config
	once sync.Once
)

func LoadConfig() *config.Config {
	once.Do(func() {
		cfg = &config.Config{
			ServerPort:      getEnv("SERVER_PORT", ":8080"),
			DBHost:          getEnv("DB_HOST", "localhost"),
			DBPort:          getEnv("DB_PORT", "3306"),
			DBUser:          getEnv("DB_USER", "root"),
			DBPassword:      getEnv("DB_PASSWORD", "123456"),
			DBName:          getEnv("DB_NAME", "survey_platform"),
			RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
			RedisPassword:   getEnv("REDIS_PASSWORD", ""),
			RedisDB:         getEnvInt("REDIS_DB", 0),
			JWTSecret:       getEnv("JWT_SECRET", "survey-platform-secret-key-2024"),
			EmailHost:       getEnv("EMAIL_HOST", "smtp.example.com"),
			EmailPort:       getEnvInt("EMAIL_PORT", 465),
			EmailUser:       getEnv("EMAIL_USER", "noreply@example.com"),
			EmailPassword:   getEnv("EMAIL_PASSWORD", ""),
			FileStoragePath: getEnv("FILE_STORAGE_PATH", "./storage"),
			ExportTemplate:  getEnv("EXPORT_TEMPLATE", "./templates"),
		}
	})
	return cfg
}

func GetConfig() *config.Config {
	if cfg == nil {
		return LoadConfig()
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
