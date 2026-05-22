package database

import (
	"fmt"
	"podcast-manager/internal/config"
	"podcast-manager/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	logrus.Info("Database connection established")

	err = DB.AutoMigrate(
		&models.Podcast{},
		&models.Episode{},
		&models.PlaybackProgress{},
		&models.Note{},
		&models.Playlist{},
		&models.PlaylistItem{},
		&models.ListeningHistory{},
	)
	if err != nil {
		logrus.Fatalf("Failed to migrate database: %v", err)
	}

	logrus.Info("Database migration completed")
}

func BeginTransaction() *gorm.DB {
	return DB.Begin()
}
