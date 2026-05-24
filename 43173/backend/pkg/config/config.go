package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireHour int    `mapstructure:"expire_hour"`
	Issuer     string `mapstructure:"issuer"`
}

type UploadConfig struct {
	MaxSize    int64    `mapstructure:"max_size"`
	AudioPath  string   `mapstructure:"audio_path"`
	ImagePath  string   `mapstructure:"image_path"`
	AllowedExt []string `mapstructure:"allowed_ext"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Filename string `mapstructure:"filename"`
	MaxSize  int    `mapstructure:"max_size"`
	MaxAge   int    `mapstructure:"max_age"`
}

var AppConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	AppConfig = &config
	return &config, nil
}

func LoadConfigFromEnv() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			Mode:         getEnv("SERVER_MODE", "release"),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 60),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "music_platform"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Shanghai"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			PoolSize: getEnvInt("REDIS_POOL_SIZE", 100),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "music-platform-secret-key"),
			ExpireHour: getEnvInt("JWT_EXPIRE_HOUR", 72),
			Issuer:     getEnv("JWT_ISSUER", "music-platform"),
		},
		Upload: UploadConfig{
			MaxSize:    getEnvInt64("UPLOAD_MAX_SIZE", 50<<20),
			AudioPath:  getEnv("UPLOAD_AUDIO_PATH", "./uploads/audio"),
			ImagePath:  getEnv("UPLOAD_IMAGE_PATH", "./uploads/images"),
			AllowedExt: []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a"},
		},
		Log: LogConfig{
			Level:    getEnv("LOG_LEVEL", "info"),
			Filename: getEnv("LOG_FILENAME", "./logs/app.log"),
			MaxSize:  getEnvInt("LOG_MAX_SIZE", 100),
			MaxAge:   getEnvInt("LOG_MAX_AGE", 30),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		return atoi(value, defaultValue)
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		return atoi64(value, defaultValue)
	}
	return defaultValue
}

func atoi(s string, defaultVal int) int {
	result := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	if result == 0 {
		return defaultVal
	}
	return result
}

func atoi64(s string, defaultVal int64) int64 {
	var result int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int64(c-'0')
		}
	}
	if result == 0 {
		return defaultVal
	}
	return result
}
