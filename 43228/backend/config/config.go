package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	Charset      string `mapstructure:"charset"`
	ParseTime    bool   `mapstructure:"parse_time"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type UploadConfig struct {
	MaxSize             int      `mapstructure:"max_size"`
	AllowedTypes        []string `mapstructure:"allowed_types"`
	TeaImagesPath       string   `mapstructure:"tea_images_path"`
	PackagingImagesPath string   `mapstructure:"packaging_images_path"`
	EvaluationImagesPath string  `mapstructure:"evaluation_images_path"`
}

type ScheduleConfig struct {
	InventoryCheckInterval string `mapstructure:"inventory_check_interval"`
	ShelfLifeCheckInterval string `mapstructure:"shelf_life_check_interval"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Schedule ScheduleConfig `mapstructure:"schedule"`
}

var (
	Cfg  *Config
	once sync.Once
)

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name, d.Charset, d.ParseTime,
	)
}

func Load(path string) (*Config, error) {
	var loadErr error
	once.Do(func() {
		v := viper.New()
		v.SetConfigFile(path)
		v.SetConfigType("yaml")

		if err := v.ReadInConfig(); err != nil {
			loadErr = fmt.Errorf("读取配置文件失败: %w", err)
			return
		}

		var cfg Config
		if err := v.Unmarshal(&cfg); err != nil {
			loadErr = fmt.Errorf("解析配置文件失败: %w", err)
			return
		}

		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("[config] 检测到配置文件变化: %s, 重新加载...", e.Name)
			var newCfg Config
			if err := v.Unmarshal(&newCfg); err != nil {
				log.Printf("[config] 重新加载配置失败: %v", err)
				return
			}
			Cfg = &newCfg
			log.Printf("[config] 配置已热加载完成")
		})

		Cfg = &cfg
	})
	return Cfg, loadErr
}

func Get() *Config {
	return Cfg
}
