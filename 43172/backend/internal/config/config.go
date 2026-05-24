package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.Charset)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type UploadConfig struct {
	Path    string   `mapstructure:"path"`
	MaxSize int64    `mapstructure:"max_size"`
	Types   []string `mapstructure:"types"`
}

var (
	config     *Config
	configLock sync.RWMutex
)

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	configLock.Lock()
	config = cfg
	configLock.Unlock()

	watchConfig(v)

	return cfg, nil
}

func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func watchConfig(v *viper.Viper) {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if e.Op&fsnotify.Write == fsnotify.Write {
			newCfg := &Config{}
			if err := v.Unmarshal(newCfg); err != nil {
				return
			}
			configLock.Lock()
			config = newCfg
			configLock.Unlock()
		}
	})
}

func LoadConfigFromEnv() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "release",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "root",
			DBName:   "luxury_trading",
			Charset:  "utf8mb4",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		},
		JWT: JWTConfig{
			Secret:     "luxury-trading-secret-key-change-in-production",
			ExpireTime: 24 * time.Hour,
		},
		Log: LogConfig{
			Level: "info",
			File:  "logs/app.log",
		},
		Upload: UploadConfig{
			Path:    "./uploads",
			MaxSize: 10 * 1024 * 1024,
			Types:   []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"},
		},
	}
}
