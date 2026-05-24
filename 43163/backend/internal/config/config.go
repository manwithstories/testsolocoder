package config

import (
	"os"
	"time"
)

type Config struct {
	Server   Server
	Database Database
	JWT      JWT
}

type Server struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	Driver string
	DSN    string
}

type JWT struct {
	Secret     string
	ExpireHour int
}

func Load() *Config {
	secret := env("JWT_SECRET", "printshop-secret-change-me")
	return &Config{
		Server: Server{
			Port:         env("SERVER_PORT", "8080"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: Database{
			Driver: env("DB_DRIVER", "sqlite"),
			DSN:    env("DB_DSN", "printshop.db"),
		},
		JWT: JWT{
			Secret:     secret,
			ExpireHour: 24,
		},
	}
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
