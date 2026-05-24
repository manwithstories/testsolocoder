package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName      string
	Port         string
	DBPath       string
	JWTSecret    string
	JWTExpire    time.Duration
	UploadDir    string
	LogDir       string
	CORSOrigin   string
	MaxFileSize  int64
}

var Cfg *Config

func Load() {
	_ = godotenv.Load()
	exp, _ := strconv.Atoi(get("JWT_EXPIRE_HOURS", "72"))
	maxSize, _ := strconv.Atoi(get("MAX_FILE_SIZE_MB", "20"))
	Cfg = &Config{
		AppName:     get("APP_NAME", "WatchPlatform"),
		Port:        get("PORT", "8080"),
		DBPath:      get("DB_PATH", "watchplatform.db"),
		JWTSecret:   get("JWT_SECRET", "change-me-in-prod"),
		JWTExpire:   time.Duration(exp) * time.Hour,
		UploadDir:   get("UPLOAD_DIR", "./uploads"),
		LogDir:      get("LOG_DIR", "./logs"),
		CORSOrigin:  get("CORS_ORIGIN", "*"),
		MaxFileSize: int64(maxSize) * 1024 * 1024,
	}
}

func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
