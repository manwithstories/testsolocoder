package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Upload UploadConfig `mapstructure:"upload"`
}

type ServerConfig struct {
	Port     string `mapstructure:"port"`
	Mode     string `mapstructure:"mode"`
	LogLevel string `mapstructure:"log_level"`
}

type MySQLConfig struct {
	DSN     string `mapstructure:"dsn"`
	MaxOpen int    `mapstructure:"max_open"`
	MaxIdle int    `mapstructure:"max_idle"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Secret         string `mapstructure:"secret"`
	ExpireMinutes  int    `mapstructure:"expire_minutes"`
}

type UploadConfig struct {
	Dir            string `mapstructure:"dir"`
	CertificateDir string `mapstructure:"certificate_dir"`
	MaxFileSize    int64  `mapstructure:"max_file_size"`
}

var Cfg *Config

func Load(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("read config error: %v", err)
	}
	Cfg = &Config{}
	if err := v.Unmarshal(Cfg); err != nil {
		log.Fatalf("parse config error: %v", err)
	}
	return Cfg
}
