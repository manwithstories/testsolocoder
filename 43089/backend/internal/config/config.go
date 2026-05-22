package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
	Email    EmailConfig    `yaml:"email"`
}

type EmailConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	From       string `yaml:"from"`
	UseTLS     bool   `yaml:"use_tls"`
	SkipVerify bool   `yaml:"skip_verify"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type JWTConfig struct {
	Secret    string `yaml:"secret"`
	ExpiresIn int    `yaml:"expires_in"`
}

type UploadConfig struct {
	Path        string   `yaml:"path"`
	MaxSize     int64    `yaml:"max_size"`
	AllowedTypes []string `yaml:"allowed_types"`
}

var AppConfig *Config

func LoadConfig() error {
	_ = godotenv.Load()

	configFile := "configs/config.yaml"
	if os.Getenv("CONFIG_FILE") != "" {
		configFile = os.Getenv("CONFIG_FILE")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	overrideFromEnv(&config)
	AppConfig = &config

	log.Printf("Config loaded successfully. Server mode: %s, Port: %d", config.Server.Mode, config.Server.Port)
	return nil
}

func overrideFromEnv(config *Config) {
	if os.Getenv("SERVER_PORT") != "" {
		fmt.Sscanf(os.Getenv("SERVER_PORT"), "%d", &config.Server.Port)
	}
	if os.Getenv("SERVER_MODE") != "" {
		config.Server.Mode = os.Getenv("SERVER_MODE")
	}
	if os.Getenv("DB_HOST") != "" {
		config.Database.Host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		fmt.Sscanf(os.Getenv("DB_PORT"), "%d", &config.Database.Port)
	}
	if os.Getenv("DB_USER") != "" {
		config.Database.User = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		config.Database.Password = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_NAME") != "" {
		config.Database.DBName = os.Getenv("DB_NAME")
	}
	if os.Getenv("JWT_SECRET") != "" {
		config.JWT.Secret = os.Getenv("JWT_SECRET")
	}
	if os.Getenv("JWT_EXPIRES_IN") != "" {
		fmt.Sscanf(os.Getenv("JWT_EXPIRES_IN"), "%d", &config.JWT.ExpiresIn)
	}
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
