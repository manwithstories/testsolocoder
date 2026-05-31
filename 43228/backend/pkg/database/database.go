package database

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"tea-platform/config"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func Init() {
	once.Do(func() {
		cfg := config.Get()

		gormLevel := parseGormLevel(cfg.Log.Level)

		newLogger := gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  gormLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)

		var err error
		DB, err = gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			log.Fatalf("连接数据库失败: %v", err)
		}

		sqlDB, err := DB.DB()
		if err != nil {
			log.Fatalf("获取底层 sql.DB 失败: %v", err)
		}
		sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Hour)

		log.Printf("[database] 数据库连接成功: %s@%s:%s/%s",
			cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	})
}

func parseGormLevel(level string) gormlogger.LogLevel {
	switch level {
	case "debug":
		return gormlogger.Info
	case "info":
		return gormlogger.Warn
	case "error":
		return gormlogger.Error
	default:
		return gormlogger.Warn
	}
}

func GetDB() *gorm.DB {
	return DB
}
