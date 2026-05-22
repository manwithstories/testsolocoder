package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

type JWTConfig struct {
	Secret    string
	ExpireHour int
}

type LogConfig struct {
	Level  string
	Output string
}

func LoadConfig(env string) *Config {
	configs := map[string]*Config{
		"development": {
			Server: ServerConfig{
				Port: ":8080",
				Mode: "debug",
			},
			Database: DatabaseConfig{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     "3306",
				User:     "root",
				Password: "password",
				DBName:   "qa_platform",
				Charset:  "utf8mb4",
			},
			JWT: JWTConfig{
				Secret:    "qa-platform-dev-secret-key-2024",
				ExpireHour: 24,
			},
			Log: LogConfig{
				Level:  "debug",
				Output: "stdout",
			},
		},
		"production": {
			Server: ServerConfig{
				Port: ":8080",
				Mode: "release",
			},
			Database: DatabaseConfig{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     "3306",
				User:     "root",
				Password: "password",
				DBName:   "qa_platform_prod",
				Charset:  "utf8mb4",
			},
			JWT: JWTConfig{
				Secret:    "qa-platform-prod-secret-key-2024-prod",
				ExpireHour: 24,
			},
			Log: LogConfig{
				Level:  "info",
				Output: "file",
			},
		},
	}

	if config, ok := configs[env]; ok {
		return config
	}
	return configs["development"]
}
