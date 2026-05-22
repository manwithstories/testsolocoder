package database

import (
	"log"
	"wedding-planner/config"
	"wedding-planner/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	cfg := config.AppConfig

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	err = db.AutoMigrate(
		&models.User{},
		&models.Wedding{},
		&models.Vendor{},
		&models.VendorReview{},
		&models.Guest{},
		&models.GuestTable{},
		&models.BudgetItem{},
		&models.Payment{},
		&models.Task{},
		&models.TaskTemplate{},
		&models.Document{},
		&models.OperationLog{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = db
	log.Println("Database connected and migrated successfully")
	return db
}

func GetDB() *gorm.DB {
	return DB
}
