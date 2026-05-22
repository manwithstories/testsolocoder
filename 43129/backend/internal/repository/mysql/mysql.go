package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"beauty-salon-system/internal/model"
)

var DB *gorm.DB

func Init(dsn string, maxOpenConns, maxIdleConns int) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("connect mysql: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := autoMigrate(); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Customer{},
		&model.Technician{},
		&model.TechnicianLeave{},
		&model.Service{},
		&model.PackageService{},
		&model.CustomerPackage{},
		&model.Appointment{},
		&model.Payment{},
		&model.MemberCard{},
		&model.Product{},
		&model.ProductRecord{},
		&model.ProductSale{},
		&model.Review{},
		&model.Notification{},
		&model.AuditLog{},
	)
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
