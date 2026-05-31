package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Storage  StorageConfig  `yaml:"storage"`
	Pricing  PricingConfig  `yaml:"pricing"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
	Port         int    `yaml:"port"`
	Mode         string `yaml:"mode"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	MaxUploadSize int64 `yaml:"max_upload_size"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	Timezone string `yaml:"timezone"`
}

type JWTConfig struct {
	SecretKey              string `yaml:"secret_key"`
	AccessTokenExpireHours int    `yaml:"access_token_expire_hours"`
	RefreshTokenExpireHours int   `yaml:"refresh_token_expire_hours"`
}

type StorageConfig struct {
	UploadPath        string   `yaml:"upload_path"`
	MaxFileSize       int64    `yaml:"max_file_size"`
	AllowedExtensions []string `yaml:"allowed_extensions"`
}

type PricingConfig struct {
	PlatformFeeRate        float64 `yaml:"platform_fee_rate"`
	DesignerFeeRate        float64 `yaml:"designer_fee_rate"`
	PrinterFeeRate         float64 `yaml:"printer_fee_rate"`
	VolumePriceMultiplier  float64 `yaml:"volume_price_multiplier"`
	TimePriceMultiplier    float64 `yaml:"time_price_multiplier"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode, d.Timezone)
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

func GetConfig() *Config {
	return AppConfig
}
