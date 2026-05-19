package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Log       LogConfig       `mapstructure:"log"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Queue     QueueConfig     `mapstructure:"queue"`
	Retry     RetryConfig     `mapstructure:"retry"`
	Webhook   WebhookConfig   `mapstructure:"webhook"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	User             string `mapstructure:"user"`
	Password         string `mapstructure:"password"`
	DBName           string `mapstructure:"dbname"`
	Charset          string `mapstructure:"charset"`
	ParseTime        bool   `mapstructure:"parse_time"`
	MaxOpenConns     int    `mapstructure:"max_open_conns"`
	MaxIdleConns     int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime  int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type RateLimitConfig struct {
	Global         RateLimitRule `mapstructure:"global"`
	DefaultChannel RateLimitRule `mapstructure:"default_channel"`
}

type RateLimitRule struct {
	Enabled           bool    `mapstructure:"enabled"`
	RequestsPerSecond float64 `mapstructure:"requests_per_second"`
	Burst             int     `mapstructure:"burst"`
}

type QueueConfig struct {
	WorkerCount  int `mapstructure:"worker_count"`
	MaxQueueSize int `mapstructure:"max_queue_size"`
	PollInterval int `mapstructure:"poll_interval"`
}

type RetryConfig struct {
	MaxAttempts      int      `mapstructure:"max_attempts"`
	InitialBackoff   int      `mapstructure:"initial_backoff"`
	MaxBackoff       int      `mapstructure:"max_backoff"`
	BackoffMultiplier float64 `mapstructure:"backoff_multiplier"`
	RetryableErrors  []string `mapstructure:"retryable_errors"`
}

type WebhookConfig struct {
	Timeout    int `mapstructure:"timeout"`
	MaxRetries int `mapstructure:"max_retries"`
	Backoff    int `mapstructure:"backoff"`
}

var AppConfig *Config

func Load(env string) error {
	v := viper.New()

	configFile := fmt.Sprintf("configs/config.%s.yaml", env)
	v.SetConfigFile(configFile)

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file failed: %w", err)
	}

	AppConfig = &Config{}
	if err := v.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("unmarshal config failed: %w", err)
	}

	return nil
}

func GetEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	return env
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset, c.ParseTime)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
