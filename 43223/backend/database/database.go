package database

import (
	"fmt"
	"log"

	"coffee-platform/config"
	"coffee-platform/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connected successfully")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductImage{},
		&models.RoastingRecord{},
		&models.RoastingDataPoint{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.CuppingScore{},
		&models.CuppingCriteria{},
		&models.RoasterCertification{},
		&models.OperationLog{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
