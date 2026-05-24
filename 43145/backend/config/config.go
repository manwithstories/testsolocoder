package config

type Config struct {
	ServerPort      string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	JWTSecret       string
	EmailHost       string
	EmailPort       int
	EmailUser       string
	EmailPassword   string
	FileStoragePath string
	ExportTemplate  string
}
