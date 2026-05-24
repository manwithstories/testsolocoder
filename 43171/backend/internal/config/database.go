package config

import (
	"drone-rental/internal/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg DatabaseConfig) error {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	err = db.AutoMigrate(
		&model.User{},
		&model.Drone{},
		&model.RentalOrder{},
		&model.Payment{},
		&model.AerialService{},
		&model.ServiceBid{},
		&model.FlightRecord{},
		&model.InsuranceClaim{},
		&model.Review{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	DB = db
	return nil
}
