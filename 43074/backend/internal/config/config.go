package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Upload   UploadConfig
	ISBN     ISBNConfig
	Log      LogConfig
}

type ServerConfig struct {
	Host string
	Port int
	Mode string
}

type DatabaseConfig struct {
	Path string
}

type UploadConfig struct {
	MaxSize    int64
	AllowedExt []string
	SavePath   string
	AccessURL  string
}

type ISBNConfig struct {
	Timeout int
}

type LogConfig struct {
	Level  string
	Format string
	Output string
}

var AppConfig *Config

func Load() error {
	_ = godotenv.Load()

	AppConfig = &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "data/booklibrary.db"),
		},
		Upload: UploadConfig{
			MaxSize:    int64(getEnvAsInt("UPLOAD_MAX_SIZE", 5)) * 1024 * 1024,
			AllowedExt: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
			SavePath:   getEnv("UPLOAD_PATH", "uploads"),
			AccessURL:  getEnv("UPLOAD_ACCESS_URL", "/uploads"),
		},
		ISBN: ISBNConfig{
			Timeout: getEnvAsInt("ISBN_TIMEOUT", 10),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
			Output: getEnv("LOG_OUTPUT", "stdout"),
		},
	}

	if err := os.MkdirAll(AppConfig.Upload.SavePath, 0755); err != nil {
		return fmt.Errorf("create upload directory: %w", err)
	}
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
