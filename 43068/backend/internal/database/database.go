package database

import (
	"log"

	"freelancer-management/internal/config"
	"freelancer-management/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	cfg := config.GetConfig()
	var err error

	switch cfg.Database.Driver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	default:
		log.Fatalf("Unsupported database driver: %s", cfg.Database.Driver)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.Project{},
		&models.Milestone{},
		&models.TimeEntry{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.InvoiceCounter{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	initInvoiceCounter()

	log.Println("Database initialized successfully")
	return DB
}

func initInvoiceCounter() {
	var count int64
	DB.Model(&models.InvoiceCounter{}).Count(&count)
	if count == 0 {
		year := 2026
		DB.Create(&models.InvoiceCounter{
			Year:  year,
			Count: 0,
		})
	}
}

func GetDB() *gorm.DB {
	return DB
}
