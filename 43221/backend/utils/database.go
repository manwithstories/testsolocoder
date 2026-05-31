package utils

import (
	"fmt"

	"consultation-platform/config"
	"consultation-platform/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"go.uber.org/zap"
)

var DB *gorm.DB

func InitDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Service{},
		&models.Schedule{},
		&models.Appointment{},
		&models.Payment{},
		&models.ConsultRecord{},
		&models.Review{},
		&models.Notification{},
		&models.NotificationTemplate{},
	)
	if err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	Logger.Info("Database connected and migrated successfully")
	return DB, nil
}

func GetDB() *gorm.DB {
	return DB
}

func LogInfo(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func LogError(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}
