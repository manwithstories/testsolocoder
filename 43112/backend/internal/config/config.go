package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
	Log      LogConfig      `yaml:"log"`
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
	Timezone string `yaml:"timezone"`
}

type JWTConfig struct {
	Secret       string `yaml:"secret"`
	ExpireHours  int    `yaml:"expire_hours"`
}

type UploadConfig struct {
	MaxSizeMB      int    `yaml:"max_size_mb"`
	AllowedImage   string `yaml:"allowed_image"`
	AllowedVideo   string `yaml:"allowed_video"`
	AllowedDocument string `yaml:"allowed_document"`
	UploadDir     string `yaml:"upload_dir"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	LogDir string `yaml:"log_dir"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
