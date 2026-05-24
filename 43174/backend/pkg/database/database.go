package database

import (
	"fmt"
	"log"
	"time"

	"campus-trade-platform/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(databaseURL string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Failed to connect to database, retrying...", i+1)
		time.Sleep(time.Second * 5)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after 10 attempts: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = migrateDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	log.Println("Database connection established successfully")
	return db, nil
}

func migrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Textbook{},
		&models.Note{},
		&models.Transaction{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderStatusHistory{},
		&models.Message{},
		&models.Review{},
		&models.Notification{},
	)
	if err != nil {
		return err
	}
	log.Println("Database migration completed successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("Database connection closed")
		}
	}
}
