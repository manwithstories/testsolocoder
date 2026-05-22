package models

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"auction-system/config"
)

var DB *gorm.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Name,
		config.AppConfig.Database.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return err
	}

	err = DB.AutoMigrate(
		&User{},
		&Category{},
		&AuctionItem{},
		&AuctionImage{},
		&AuctionSession{},
		&AuctionItemSession{},
		&Bid{},
		&AutoBid{},
		&Order{},
		&Payment{},
		&Review{},
		&Notification{},
		&SystemLog{},
	)

	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return err
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

func BeginTransaction() *gorm.DB {
	return DB.Begin()
}
