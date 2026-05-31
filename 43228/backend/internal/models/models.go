package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&User{},
		&Tea{},
		&TeaImage{},
		&TastingRecord{},
		&TastingImage{},
		&Traceability{},
		&Order{},
		&OrderLog{},
		&Collection{},
		&AppraisalReport{},
		&Post{},
		&Comment{},
		&Like{},
		&TastingActivity{},
		&ActivityParticipant{},
		&OperationLog{},
	)
	if err != nil {
		return err
	}

	if err := initSeedData(db); err != nil {
		log.Printf("初始化基础数据失败: %v", err)
	}

	return nil
}

func initSeedData(db *gorm.DB) error {
	if err := initAdminUser(db); err != nil {
		return err
	}
	return nil
}

func initAdminUser(db *gorm.DB) error {
	var count int64
	db.Model(&User{}).Where("role = ?", UserRoleAdmin).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &User{
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "admin@tea-platform.com",
		Phone:    "13800138000",
		Role:     UserRoleAdmin,
		Status:   UserStatusActive,
	}

	return db.Create(admin).Error
}
