package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret     string
	ExpireHour int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type UploadConfig struct {
	MaxSize int64
	Dir     string
}

var Cfg *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	Cfg = &Config{
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 8080),
			Mode: getEnvStr("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			DSN: getEnvStr("DB_DSN", "root:password@tcp(localhost:3306)/matchmaking?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		JWT: JWTConfig{
			Secret:     getEnvStr("JWT_SECRET", "matchmaking-secret-key-change-in-production"),
			ExpireHour: getEnvInt("JWT_EXPIRE_HOUR", 24),
		},
		Redis: RedisConfig{
			Addr:     getEnvStr("REDIS_ADDR", "localhost:6379"),
			Password: getEnvStr("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Upload: UploadConfig{
			MaxSize: int64(getEnvInt("UPLOAD_MAX_SIZE", 10)) * 1024 * 1024,
			Dir:     getEnvStr("UPLOAD_DIR", "uploads"),
		},
	}
}

func getEnvStr(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return fallback
}
