package config

import (
	"os"
)

type Config struct {
	ServerPort   string
	DatabasePath string
	RedisAddr    string
	RedisPass    string
	RedisDB      int
	JWTSecret    string
	PlatformFee  float64
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DatabasePath: getEnv("DATABASE_PATH", "repair.db"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:      0,
		JWTSecret:    getEnv("JWT_SECRET", "repair-platform-secret-key-2024"),
		PlatformFee:  0.15,
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
