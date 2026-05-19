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
	Email    EmailConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Charset  string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type EmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type UploadConfig struct {
	Path         string
	MaxSizeBytes int64
}

var AppConfig *Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	jwtExpireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)

	AppConfig = &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "venue_booking"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "default-secret-key"),
			ExpireHours: jwtExpireHours,
		},
		Email: EmailConfig{
			Host:     getEnv("SMTP_HOST", "smtp.example.com"),
			Port:     smtpPort,
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
		},
		Upload: UploadConfig{
			Path:         getEnv("UPLOAD_PATH", "./uploads"),
			MaxSizeBytes: maxUploadSize,
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
