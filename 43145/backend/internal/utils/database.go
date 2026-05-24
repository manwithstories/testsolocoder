package utils

import (
	"fmt"
	"log"
	"survey-platform/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	cfg := GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Survey{},
		&model.Question{},
		&model.Option{},
		&model.LogicJump{},
		&model.Response{},
		&model.Answer{},
		&model.DistributionLink{},
		&model.Invitation{},
	)

	log.Println("Database connection established and tables migrated")
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
