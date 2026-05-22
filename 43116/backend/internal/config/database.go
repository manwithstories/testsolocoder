package config

import (
	"car-rental/internal/model"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *DatabaseConfig) error {
	var err error

	DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("数据库连接成功")

	err = autoMigrate()
	if err != nil {
		return err
	}

	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.Car{},
		&model.CarImage{},
		&model.City{},
		&model.Store{},
		&model.Booking{},
		&model.PromoCode{},
		&model.PricingRule{},
		&model.Order{},
		&model.Review{},
		&model.MaintenancePlan{},
		&model.Message{},
		&model.OperationLog{},
		&model.Notification{},
	)
}

func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}