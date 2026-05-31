package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Museum   MuseumConfig   `mapstructure:"museum"`
}

type ServerConfig struct {
	Port          int    `mapstructure:"port"`
	Mode          string `mapstructure:"mode"`
	UploadDir     string `mapstructure:"upload_dir"`
	MaxUploadSize int    `mapstructure:"max_upload_size"`
}

type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"name"`
	MaxOpenConns   int    `mapstructure:"max_open_conns"`
	MaxIdleConns   int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Secret       string `mapstructure:"secret"`
	ExpireHours  int    `mapstructure:"expire_hours"`
}

type MuseumConfig struct {
	OpenTime              string `mapstructure:"open_time"`
	CloseTime             string `mapstructure:"close_time"`
	MaxVisitorsPerSlot    int    `mapstructure:"max_visitors_per_slot"`
	SlotDurationMinutes   int    `mapstructure:"slot_duration_minutes"`
	AdvanceBookingDays    int    `mapstructure:"advance_booking_days"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
