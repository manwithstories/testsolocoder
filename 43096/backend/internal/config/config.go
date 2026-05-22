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
	Redis    RedisConfig
	JWT      JWTConfig
	Seat     SeatConfig
	Log      LogConfig
	Email    EmailConfig
	SMS      SMSConfig
	Payment  PaymentConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type SeatConfig struct {
	LockMinutes int
}

type LogConfig struct {
	Level      string
	File       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

type SMSConfig struct {
	AccessKey string
	SecretKey string
	SignName  string
}

type PaymentConfig struct {
	AlipayAppID string
	WechatAppID string
}

var AppConfig *Config

func LoadConfig() {
	_ = godotenv.Load()

	AppConfig = &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnvInt("DB_PORT", 3306),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", "123456"),
			Database:        getEnv("DB_NAME", "ticket_system"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 3600),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "127.0.0.1"),
			Port:         getEnvInt("REDIS_PORT", 6379),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvInt("REDIS_DB", 0),
			PoolSize:     getEnvInt("REDIS_POOL_SIZE", 100),
			MinIdleConns: getEnvInt("REDIS_MIN_IDLE_CONNS", 10),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key-here"),
			ExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 24),
		},
		Seat: SeatConfig{
			LockMinutes: getEnvInt("SEAT_LOCK_MINUTES", 15),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			File:       getEnv("LOG_FILE", "./logs/app.log"),
			MaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 10),
			MaxAge:     getEnvInt("LOG_MAX_AGE", 30),
		},
		Email: EmailConfig{
			SMTPHost: getEnv("EMAIL_SMTP_HOST", "smtp.example.com"),
			SMTPPort: getEnvInt("EMAIL_SMTP_PORT", 587),
			Username: getEnv("EMAIL_USERNAME", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
			From:     getEnv("EMAIL_FROM", "ticket@example.com"),
		},
		SMS: SMSConfig{
			AccessKey: getEnv("SMS_ACCESS_KEY", ""),
			SecretKey: getEnv("SMS_SECRET_KEY", ""),
			SignName:  getEnv("SMS_SIGN_NAME", "票务系统"),
		},
		Payment: PaymentConfig{
			AlipayAppID: getEnv("ALIPAY_APP_ID", ""),
			WechatAppID: getEnv("WECHAT_APP_ID", ""),
		},
	}

	log.Println("Configuration loaded successfully")
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
