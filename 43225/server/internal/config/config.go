package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server       ServerConfig       `yaml:"server"`
	Database     DatabaseConfig     `yaml:"database"`
	JWT          JWTConfig          `yaml:"jwt"`
	Upload       UploadConfig       `yaml:"upload"`
	Finance      FinanceConfig      `yaml:"finance"`
	Notification NotificationConfig `yaml:"notification"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpiresHour int    `yaml:"expires_hour"`
}

type UploadConfig struct {
	MaxSize     int      `yaml:"max_size"`
	AllowedTypes []string `yaml:"allowed_types"`
	Path        string   `yaml:"path"`
}

type FinanceConfig struct {
	DefaultCurrency   string   `yaml:"default_currency"`
	SupportedCurrencies []string `yaml:"supported_currencies"`
	PlatformFeeRate  float64  `yaml:"platform_fee_rate"`
	DockFeeRate      float64  `yaml:"dock_fee_rate"`
}

type NotificationConfig struct {
	ReminderDays int `yaml:"reminder_days"`
}

var AppConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	AppConfig = &config
	return &config, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
