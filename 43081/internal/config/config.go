package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultAccount string      `yaml:"default_account"`
	DefaultCurrency string     `yaml:"default_currency"`
	BudgetAlert    bool        `yaml:"budget_alert"`
	Backup         BackupConfig `yaml:"backup"`
	DatabasePath   string      `yaml:"database_path"`
	LogPath        string      `yaml:"log_path"`
}

type BackupConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Schedule      string `yaml:"schedule"`
	LocalPath     string `yaml:"local_path"`
	RetentionDays int    `yaml:"retention_days"`
	MaxBackups    int    `yaml:"max_backups"`
}

var defaultConfig = Config{
	DefaultAccount:  "cash",
	DefaultCurrency: "CNY",
	BudgetAlert:     true,
	Backup: BackupConfig{
		Enabled:       false,
		Schedule:      "daily",
		LocalPath:     "~/.finance-tracker/backups",
		RetentionDays: 30,
		MaxBackups:    10,
	},
	DatabasePath: "~/.finance-tracker/finance.db",
	LogPath:      "~/.finance-tracker/logs",
}

func Load() (*Config, error) {
	configPath := getConfigPath()
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := Save(&defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return &defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	applyDefaults(&cfg)
	return &cfg, nil
}

func Save(cfg *Config) error {
	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".finance-tracker/config.yaml"
	}
	return filepath.Join(home, ".finance-tracker", "config.yaml")
}

func applyDefaults(cfg *Config) {
	if cfg.DefaultCurrency == "" {
		cfg.DefaultCurrency = defaultConfig.DefaultCurrency
	}
	if cfg.DatabasePath == "" {
		cfg.DatabasePath = defaultConfig.DatabasePath
	}
	if cfg.LogPath == "" {
		cfg.LogPath = defaultConfig.LogPath
	}
	if cfg.Backup.LocalPath == "" {
		cfg.Backup.LocalPath = defaultConfig.Backup.LocalPath
	}
	if cfg.Backup.Schedule == "" {
		cfg.Backup.Schedule = defaultConfig.Backup.Schedule
	}
	if cfg.Backup.RetentionDays == 0 {
		cfg.Backup.RetentionDays = defaultConfig.Backup.RetentionDays
	}
	if cfg.Backup.MaxBackups == 0 {
		cfg.Backup.MaxBackups = defaultConfig.Backup.MaxBackups
	}
}

func ExpandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}
	return path
}
