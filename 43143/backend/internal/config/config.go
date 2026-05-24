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
	Payment      PaymentConfig      `yaml:"payment"`
	Notification NotificationConfig `yaml:"notification"`
	Encryption   EncryptionConfig   `yaml:"encryption"`
	Storage      StorageConfig      `yaml:"storage"`
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

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

type JWTConfig struct {
	Secret              string `yaml:"secret"`
	AccessTokenExpire   int    `yaml:"access_token_expire"`
	RefreshTokenExpire  int    `yaml:"refresh_token_expire"`
}

type PaymentConfig struct {
	ServiceFeeRate float64     `yaml:"service_fee_rate"`
	EscrowDays     int         `yaml:"escrow_days"`
	Alipay         AlipayConfig `yaml:"alipay"`
	Wechat         WechatConfig `yaml:"wechat"`
}

type AlipayConfig struct {
	AppID      string `yaml:"app_id"`
	PrivateKey string `yaml:"private_key"`
}

type WechatConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type NotificationConfig struct {
	Email EmailConfig `yaml:"email"`
	SMS   SMSConfig   `yaml:"sms"`
}

type EmailConfig struct {
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SMSConfig struct {
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
}

type EncryptionConfig struct {
	MessageKey string `yaml:"message_key"`
}

type StorageConfig struct {
	UploadDir   string `yaml:"upload_dir"`
	MaxFileSize int64  `yaml:"max_file_size"`
}

var globalConfig *Config

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

func Get() *Config {
	return globalConfig
}
