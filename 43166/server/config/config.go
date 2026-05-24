package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	File     FileConfig     `yaml:"file"`
	Email    EmailConfig    `yaml:"email"`
	SMS      SMSConfig      `yaml:"sms"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Name         string `yaml:"name"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Charset      string `yaml:"charset"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
}

type JWTConfig struct {
	Secret    string `yaml:"secret"`
	ExpiresIn int    `yaml:"expiresIn"`
}

type FileConfig struct {
	UploadPath string `yaml:"uploadPath"`
	MaxSize    int64  `yaml:"maxSize"`
}

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type SMSConfig struct {
	APIKey    string `yaml:"apiKey"`
	APISecret string `yaml:"apiSecret"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

var AppConfig *Config

func LoadConfig(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	AppConfig = &config
	return nil
}
