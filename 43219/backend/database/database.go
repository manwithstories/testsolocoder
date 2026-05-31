package database

import (
	"housekeeping/config"
	"housekeeping/models"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.Type {
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	default:
		dialector = sqlite.Open(cfg.DSN)
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	DB = db
	if err := db.AutoMigrate(
		&models.User{},
		&models.StaffProfile{},
		&models.Service{},
		&models.Booking{},
		&models.Order{},
		&models.Review{},
		&models.Ticket{},
		&models.Wallet{},
		&models.Settlement{},
		&models.Withdrawal{},
	); err != nil {
		zap.L().Warn("auto migrate warn", zap.Error(err))
	}
	return db, nil
}
