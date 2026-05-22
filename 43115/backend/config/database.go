package config

import (
	"fmt"
	"log"
	"housekeeping-platform/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(cfg *DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connection established")

	err = autoMigrate()
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.ServiceProviderCert{},
		&models.ServiceCategory{},
		&models.ServiceItem{},
		&models.ServiceArea{},
		&models.Order{},
		&models.OrderInvitation{},
		&models.Review{},
		&models.Complaint{},
		&models.Bill{},
		&models.WithdrawRequest{},
		&models.Message{},
		&models.OperationLog{},
	)
}
