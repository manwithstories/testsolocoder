package database

import (
	"os"
	"path/filepath"

	"booklibrary/internal/config"
	applogger "booklibrary/internal/logger"
	"booklibrary/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	dbDir := filepath.Dir(config.AppConfig.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.Database.Path), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	err = DB.AutoMigrate(
		&models.Book{},
		&models.Tag{},
		&models.Category{},
		&models.ReadingNote{},
		&models.BorrowRecord{},
		&models.ReadingGoal{},
		&models.ReadSession{},
	)
	if err != nil {
		return err
	}

	applogger.Info("Database initialized successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
