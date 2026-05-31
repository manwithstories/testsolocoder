package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Payment  PaymentConfig  `mapstructure:"payment"`
	Order    OrderConfig    `mapstructure:"order"`
	Ticket   TicketConfig   `mapstructure:"ticket"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Type string `mapstructure:"type"`
	DSN  string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireHour int    `mapstructure:"expire_hour"`
}

type PaymentConfig struct {
	CompanyRatio float64 `mapstructure:"company_ratio"`
	StaffRatio   float64 `mapstructure:"staff_ratio"`
}

type OrderConfig struct {
	AutoConfirmHours  int `mapstructure:"auto_confirm_hours"`
	MaxReschedule     int `mapstructure:"max_reschedule"`
}

type TicketConfig struct {
	EscalateHours int `mapstructure:"escalate_hours"`
}

var C *Config

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.SetEnvPrefix("HK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	C = cfg
	return cfg, nil
}
