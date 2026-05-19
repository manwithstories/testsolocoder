package database

import (
	"fmt"
	"log"
	"time"

	"ticket-system/internal/config"
	"ticket-system/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	cfg := config.AppConfig.Database
	dsn := cfg.DSN()

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err = autoMigrate(); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	if err = initDefaultData(); err != nil {
		log.Printf("Warning: failed to init default data: %v", err)
	}

	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.Department{},
		&models.SkillGroup{},
		&models.User{},
		&models.Customer{},
		&models.TicketCounter{},
		&models.Ticket{},
		&models.TicketLog{},
		&models.Comment{},
		&models.Attachment{},
		&models.SLARecord{},
		&models.AssignmentRule{},
	)
}

func initDefaultData() error {
	var count int64
	DB.Model(&models.Department{}).Count(&count)
	if count == 0 {
		defaultDept := &models.Department{
			Name:        "技术支持部",
			Description: "默认技术支持部门",
		}
		if err := DB.Create(defaultDept).Error; err != nil {
			return err
		}

		defaultSkillGroup := &models.SkillGroup{
			Name:        "通用技术组",
			Description: "默认通用技术技能组",
		}
		if err := DB.Create(defaultSkillGroup).Error; err != nil {
			return err
		}

		defaultRule := &models.AssignmentRule{
			Name:         "默认分配规则",
			Description:  "默认负载均衡分配规则",
			Mode:         models.AssignmentModeLoadBalance,
			SkillGroupID: &defaultSkillGroup.ID,
			IsDefault:    true,
			Enabled:      true,
		}
		if err := DB.Create(defaultRule).Error; err != nil {
			return err
		}
	}
	return nil
}

func Transaction(f func(tx *gorm.DB) error) error {
	return DB.Transaction(f)
}

type ContextKey string

const DBContextKey ContextKey = "db"
