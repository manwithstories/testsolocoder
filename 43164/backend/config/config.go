package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
	Video  VideoConfig
	Redis  RedisConfig
	SMTP   SMTPConfig
	Payment PaymentConfig
	Platform PlatformConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret       string
	ExpireHours  int
}

type VideoConfig struct {
	APIKey    string
	APISecret string
	Endpoint  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type PaymentConfig struct {
	GatewayKey     string
	GatewaySecret  string
	WebhookSecret  string
}

type PlatformConfig struct {
	CommissionRate   float64
	MinWithdrawAmount float64
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	cfg := &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "tutoring_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "default-secret-key"),
			ExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 24),
		},
		Video: VideoConfig{
			APIKey:    getEnv("VIDEO_API_KEY", ""),
			APISecret: getEnv("VIDEO_API_SECRET", ""),
			Endpoint:  getEnv("VIDEO_API_ENDPOINT", "https://api.video.service"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnvInt("SMTP_PORT", 587),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
		},
		Payment: PaymentConfig{
			GatewayKey:    getEnv("PAYMENT_GATEWAY_KEY", ""),
			GatewaySecret: getEnv("PAYMENT_GATEWAY_SECRET", ""),
			WebhookSecret: getEnv("PAYMENT_GATEWAY_WEBHOOK_SECRET", ""),
		},
		Platform: PlatformConfig{
			CommissionRate:    getEnvFloat("PLATFORM_COMMISSION_RATE", 0.1),
			MinWithdrawAmount: getEnvFloat("MIN_WITHDRAW_AMOUNT", 100),
		},
	}

	AppConfig = cfg
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}
