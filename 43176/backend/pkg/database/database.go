package database

import (
	"errand-service/internal/config"
	"errand-service/internal/models"
	"errand-service/pkg/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connected successfully")
	return db, nil
}

func GetDB() *gorm.DB {
	return db
}

func Close() {
	if db != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		logger.Info("Database connection closed")
	}
}

func AutoMigrate(d *gorm.DB) error {
	logger.Info("Running database migrations...")

	err := d.AutoMigrate(
		&models.User{},
		&models.Verification{},
		&models.CourierProfile{},
		&models.Task{},
		&models.TaskImage{},
		&models.TaskAcceptLog{},
		&models.Order{},
		&models.OrderTrack{},
		&models.OrderProof{},
		&models.OrderStatusLog{},
		&models.Transaction{},
		&models.WithdrawRequest{},
		&models.RefundRequest{},
		&models.Review{},
		&models.ChatMessage{},
		&models.Notification{},
	)

	if err != nil {
		return err
	}

	logger.Info("Database migrations completed")
	return nil
}
