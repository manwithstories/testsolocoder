package models

import (
	"fmt"
	"ticket-system/config"
	"ticket-system/internal/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.App.DBUser,
		config.App.DBPassword,
		config.App.DBHost,
		config.App.DBPort,
		config.App.DBName,
		config.App.DBCharset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Log.Fatalf("Failed to get database instance: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := autoMigrate(); err != nil {
		logger.Log.Fatalf("Failed to migrate database: %v", err)
	}

	logger.Log.Info("Database connected successfully")
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&User{},
		&Activity{},
		&TicketType{},
		&Coupon{},
		&Order{},
		&OrderItem{},
		&CheckIn{},
	)
}
