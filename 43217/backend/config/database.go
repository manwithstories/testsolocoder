package config

import (
	"fmt"
	"log"

	"health-platform/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := GlobalConfig.Database.GetDSN()
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	DB = db

	err = DB.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Department{},
		&models.Employee{},
		&models.Agency{},
		&models.Package{},
		&models.PackageItem{},
		&models.TimeSlot{},
		&models.Appointment{},
		&models.Report{},
		&models.ReportItem{},
		&models.HealthRecord{},
		&models.AbnormalItem{},
		&models.RecheckReminder{},
		&models.Billing{},
		&models.BillingItem{},
		&models.Transaction{},
		&models.CompanyBudget{},
		&models.DepartmentAppointment{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database connection and migration successful")
}
