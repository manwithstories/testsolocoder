package database

import (
	"fmt"
	"log"
	"business-registration-platform/config"
	"business-registration-platform/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	cfg := config.AppConfig.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	err = migrateTables()
	if err != nil {
		return fmt.Errorf("failed to migrate tables: %v", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

func migrateTables() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.AgentProfile{},
		&models.Application{},
		&models.ProcessStep{},
		&models.ApplicationFee{},
		&models.FeeItem{},
		&models.FeeStandard{},
		&models.DiscountPolicy{},
		&models.Notification{},
		&models.NotificationTemplate{},
		&models.ExportTask{},
		&models.DownloadLog{},
		&models.OperationLog{},
	)
}

func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Println("Error getting database instance:", err)
			return
		}
		sqlDB.Close()
		log.Println("Database connection closed")
	}
}
