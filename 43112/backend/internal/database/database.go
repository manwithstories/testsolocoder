package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/models"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode, cfg.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Chapter{},
		&models.Lesson{},
		&models.Quiz{},
		&models.QuizQuestion{},
		&models.QuizOption{},
		&models.UserQuiz{},
		&models.Order{},
		&models.Coupon{},
		&models.CouponUsed{},
		&models.Question{},
		&models.Answer{},
		&models.AnswerLike{},
		&models.QuestionLike{},
		&models.Review{},
		&models.Progress{},
		&models.Note{},
		&models.InstructorApplication{},
		&models.File{},
	)
}
