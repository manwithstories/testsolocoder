package database

import (
	"fmt"
	"medical-platform/internal/config"
	"medical-platform/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() (*gorm.DB, error) {
	cfg := config.AppConfig

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	logLevel := logger.Info
	if cfg.Environment == "production" {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db

	config.Logger.Info("Database connected successfully")

	return db, nil
}

func GetDB() *gorm.DB {
	return DB
}

func WithTransaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func PreloadAll() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("User").
			Preload("Department").
			Preload("Schedules").
			Preload("Appointments", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Patient").
					Preload("Doctor").
					Preload("Consultation", func(db *gorm.DB) *gorm.DB {
						return db.Preload("Prescription.Items").
							Preload("Reports")
					}).
					Preload("Payment").
					Preload("Review")
			}).
			Preload("Reviews")
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Department{},
		&models.Doctor{},
		&models.Schedule{},
		&models.Patient{},
		&models.HealthRecord{},
		&models.Appointment{},
		&models.Consultation{},
		&models.Prescription{},
		&models.PrescriptionItem{},
		&models.ExaminationReport{},
		&models.Notification{},
		&models.Payment{},
		&models.Review{},
	)
}
