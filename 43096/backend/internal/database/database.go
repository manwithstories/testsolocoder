package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ticket-system/internal/config"
	"ticket-system/internal/models"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database instance: %v", err))
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	err = DB.AutoMigrate(
		&models.User{},
		&models.MemberLevel{},
		&models.Coupon{},
		&models.Show{},
		&models.Session{},
		&models.SeatArea{},
		&models.Seat{},
		&models.SeatChart{},
		&models.Order{},
		&models.Ticket{},
		&models.Refund{},
		&models.CheckinLog{},
		&models.PaymentLog{},
		&models.OperationLog{},
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	initMemberLevels()

	fmt.Println("Database connected and migrated successfully")
}

func initMemberLevels() {
	var count int64
	DB.Model(&models.MemberLevel{}).Count(&count)
	if count == 0 {
		levels := []models.MemberLevel{
			{Name: "普通会员", Level: 1, Discount: 1.0, MinPoints: 0, Priority: 0, Description: "注册即成为普通会员"},
			{Name: "银卡会员", Level: 2, Discount: 0.95, MinPoints: 1000, Priority: 1, Description: "积分达到1000升级"},
			{Name: "金卡会员", Level: 3, Discount: 0.9, MinPoints: 5000, Priority: 2, Description: "积分达到5000升级"},
			{Name: "铂金会员", Level: 4, Discount: 0.85, MinPoints: 20000, Priority: 3, Description: "积分达到20000升级"},
		}
		DB.Create(&levels)
		fmt.Println("Default member levels initialized")
	}
}
