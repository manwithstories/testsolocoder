package database

import (
	"fmt"
	"multishop/internal/config"
	"multishop/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Shop{},
		&models.Category{},
		&models.Product{},
		&models.ProductImage{},
		&models.ProductSpec{},
		&models.SKU{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Refund{},
		&models.Review{},
		&models.Favorite{},
		&models.Notification{},
		&models.Dispute{},
	)
}
