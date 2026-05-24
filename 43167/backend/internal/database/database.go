package database

import (
	"log"

	"watchplatform/internal/config"
	"watchplatform/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.Cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	err = DB.AutoMigrate(
		&model.User{},
		&model.Watch{},
		&model.WatchPhoto{},
		&model.AuthOrder{},
		&model.AuthReport{},
		&model.Trade{},
		&model.TradeBid{},
		&model.Favorite{},
		&model.FavoriteGroup{},
		&model.Review{},
		&model.Message{},
	)
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
}
