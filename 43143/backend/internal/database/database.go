package database

import (
	"fmt"
	"log"
	"skillshare/internal/config"
	"skillshare/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	log.Println("数据库连接成功")

	err = migrate()
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Println("数据库迁移完成")
	return nil
}

func migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Skill{},
		&models.SkillCategory{},
		&models.SkillTag{},
		&models.SkillPosting{},
		&models.Booking{},
		&models.Review{},
		&models.Message{},
		&models.Payment{},
		&models.Schedule{},
		&models.OperationLog{},
		&models.Certification{},
		&models.Complaint{},
	)
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
