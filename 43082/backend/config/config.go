package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Log      LogConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Reminder ReminderConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Charset  string
	ParseTime string
	Loc      string
}

type LogConfig struct {
	Level string
	File  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type ReminderConfig struct {
	EmailHost          string
	EmailPort          int
	EmailUser          string
	EmailPassword      string
	SMSAPIKey          string
	TemplateExpire     string
	TemplateCourse1D   string
	TemplateCourse2H   string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	emailPort, _ := strconv.Atoi(getEnv("REMINDER_EMAIL_PORT", "587"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		DB: DBConfig{
			Host:      getEnv("DB_HOST", "127.0.0.1"),
			Port:      getEnv("DB_PORT", "3306"),
			User:      getEnv("DB_USER", "root"),
			Password:  getEnv("DB_PASSWORD", "123456"),
			Name:      getEnv("DB_NAME", "gym_management"),
			Charset:   getEnv("DB_CHARSET", "utf8mb4"),
			ParseTime: getEnv("DB_PARSE_TIME", "True"),
			Loc:       getEnv("DB_LOC", "Local"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "logs/app.log"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "gym_management_secret_key"),
			ExpireHours: jwtExpireHours,
		},
		Reminder: ReminderConfig{
			EmailHost:        getEnv("REMINDER_EMAIL_HOST", "smtp.example.com"),
			EmailPort:        emailPort,
			EmailUser:        getEnv("REMINDER_EMAIL_USER", ""),
			EmailPassword:    getEnv("REMINDER_EMAIL_PASSWORD", ""),
			SMSAPIKey:        getEnv("REMINDER_SMS_API_KEY", ""),
			TemplateExpire:   getEnv("REMINDER_TEMPLATE_EXPIRE", "您的会员卡将于%s到期，请及时续费。"),
			TemplateCourse1D: getEnv("REMINDER_TEMPLATE_COURSE_1D", "您预约的课程%s将于明天%s开始，请准时参加。"),
			TemplateCourse2H: getEnv("REMINDER_TEMPLATE_COURSE_2H", "您预约的课程%s将于2小时后开始，请准时参加。"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *DBConfig) DSN() string {
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Name + "?charset=" + c.Charset + "&parseTime=" + c.ParseTime + "&loc=" + c.Loc
}
