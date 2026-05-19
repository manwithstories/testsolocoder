package database

import (
	"gym-management/config"
	mylogger "gym-management/internal/pkg/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := config.AppConfig.DB.DSN()

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		mylogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		mylogger.Fatal("Failed to get database instance", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	mylogger.Info("Database connection established successfully")
}

func GetDB() *gorm.DB {
	return DB
}
