package database

import (
	"time"

	"event-platform/internal/config"
	applog "event-platform/internal/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg config.MySQLConfig) (*gorm.DB, error) {
	var err error
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Info),
		})
		if err == nil {
			break
		}
		applog.Warnf("mysql connect attempt %d failed: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		return nil, err
	}
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)
	applog.Infof("mysql connected")
	return DB, nil
}
