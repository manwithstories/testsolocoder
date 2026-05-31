package database

import (
	"fmt"
	"log"
	"translation-platform/internal/config"
	"translation-platform/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	cfg := config.Cfg.Database

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	DB = db

	log.Println("数据库连接成功")
	return nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.LanguagePair{},
		&models.ExpertiseTag{},
		&models.Project{},
		&models.ProjectComment{},
		&models.Document{},
		&models.TranslationSegment{},
		&models.TranslationMemory{},
		&models.GlossaryTerm{},
		&models.ReviewTask{},
		&models.Payment{},
		&models.OperationLog{},
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Println("数据库迁移成功")
	return nil
}
