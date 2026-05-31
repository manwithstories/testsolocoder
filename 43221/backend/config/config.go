package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server        ServerConfig   `yaml:"server"`
	Database      DatabaseConfig `yaml:"database"`
	Redis         RedisConfig    `yaml:"redis"`
	JWT           JWTConfig      `yaml:"jwt"`
	Email         EmailConfig    `yaml:"email"`
	Payment       PaymentConfig  `yaml:"payment"`
	SensitiveWords []string      `yaml:"sensitive_words"`
	RateLimit     RateLimitConfig `yaml:"rate_limit"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslmode"`
	Timezone string `yaml:"timezone"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.Timezone)
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type JWTConfig struct {
	Secret             string `yaml:"secret"`
	AccessTokenExpire  int    `yaml:"access_token_expire"`
	RefreshTokenExpire int    `yaml:"refresh_token_expire"`
}

type EmailConfig struct {
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	FromName string `yaml:"from_name"`
}

type PaymentConfig struct {
	TimeoutMinutes    int    `yaml:"timeout_minutes"`
	AlipayAppID       string `yaml:"alipay_app_id"`
	AlipayPrivateKey  string `yaml:"alipay_private_key"`
	WechatAppID       string `yaml:"wechat_app_id"`
	WechatMchID       string `yaml:"wechat_mch_id"`
}

type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
}

var AppConfig *Config

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	AppConfig = &cfg
	return &cfg, nil
}
