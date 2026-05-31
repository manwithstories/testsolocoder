package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret    string
	ExpireHour int
}

type UploadConfig struct {
	Path     string
	MaxSize  int64
	Allowed  []string
}

var GlobalConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: config file not found, using environment variables: %v\n", err)
	}

	GlobalConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "123456"),
			DBName:   getEnv("DB_NAME", "health_platform"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "health-platform-secret-key-2024"),
			ExpireHour: getEnvInt("JWT_EXPIRE_HOUR", 24),
		},
		Upload: UploadConfig{
			Path:    getEnv("UPLOAD_PATH", "./uploads"),
			MaxSize: getEnvInt64("UPLOAD_MAX_SIZE", 50),
			Allowed: []string{".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"},
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := fmt.Sscanf(value, "%d", new(int)); err == nil {
			return intVal
		}
	}
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}
	return defaultVal
}

func getEnvInt64(key string, defaultVal int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := fmt.Sscanf(value, "%d", new(int64)); err == nil {
			return intVal
		}
	}
	if viper.IsSet(key) {
		return viper.GetInt64(key)
	}
	return defaultVal
}

func (db *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		db.User, db.Password, db.Host, db.Port, db.DBName, db.Charset)
}

func (r *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
