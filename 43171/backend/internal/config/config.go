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
	Upload   UploadConfig   `yaml:"upload"`
	Log      LogConfig      `yaml:"log"`
	Insurance InsuranceConfig `yaml:"insurance"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	DSN          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type JWTConfig struct {
	Secret       string `yaml:"secret"`
	ExpireHours  int    `yaml:"expire_hours"`
}

type UploadConfig struct {
	MaxSize    int      `yaml:"max_size"`
	SavePath   string   `yaml:"save_path"`
	AllowedExt []string `yaml:"allowed_ext"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

type InsuranceConfig struct {
	DepositRate   float64 `yaml:"deposit_rate"`
	LateFeeRate   float64 `yaml:"late_fee_rate"`
	InsuranceRate float64 `yaml:"insurance_rate"`
}

var Cfg *Config

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}
	Cfg = &cfg
	return &cfg, nil
}
