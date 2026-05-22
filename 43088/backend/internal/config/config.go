package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	RSS      RSSConfig
	Notify   NotifyConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LogConfig struct {
	Level string
}

type RSSConfig struct {
	Timeout      int
	MaxRetries   int
	UserAgent    string
	UpdatePeriod int
}

type NotifyConfig struct {
	Enabled      bool
	PollInterval int
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "podcast_manager"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		RSS: RSSConfig{
			Timeout:      getEnvAsInt("RSS_TIMEOUT", 30),
			MaxRetries:   getEnvAsInt("RSS_MAX_RETRIES", 3),
			UserAgent:    getEnv("RSS_USER_AGENT", "PodcastManager/1.0"),
			UpdatePeriod: getEnvAsInt("RSS_UPDATE_PERIOD", 60),
		},
		Notify: NotifyConfig{
			Enabled:      getEnvAsBool("NOTIFY_ENABLED", true),
			PollInterval: getEnvAsInt("NOTIFY_POLL_INTERVAL", 30),
		},
	}
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

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
