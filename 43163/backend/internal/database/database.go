package database

import (
	"log"

	"printshop/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	log.Println("database migrated successfully")
	return db, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Customer{},
		&models.Template{},
		&models.TemplateOption{},
		&models.TemplateMaterial{},
		&models.TemplateProcess{},
		&models.Order{},
		&models.OrderItem{},
		&models.PriceRule{},
		&models.ProductionLine{},
		&models.ProductionSchedule{},
		&models.FileAsset{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.AuditLog{},
	)
}
