package database

import (
	"fmt"
	"log"
	"tutoring-platform/config"
	"tutoring-platform/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Database connection established successfully")
	return db, nil
}

func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database connection is not established")
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.TeacherProfile{},
		&models.Subject{},
		&models.TeacherSubject{},
		&models.AvailabilitySlot{},
		&models.StudentProfile{},
		&models.LearningGoal{},
		&models.AssessmentQuestion{},
		&models.AssessmentAnswer{},
		&models.Booking{},
		&models.BookingRescheduleHistory{},
		&models.VideoSession{},
		&models.VideoSessionLog{},
		&models.Wallet{},
		&models.Transaction{},
		&models.WithdrawRequest{},
		&models.PaymentConfig{},
		&models.LessonNote{},
		&models.Homework{},
		&models.HomeworkSubmission{},
		&models.Feedback{},
		&models.Milestone{},
		&models.Review{},
		&models.Message{},
		&models.MessageFile{},
		&models.Notification{},
		&models.AdminAction{},
		&models.SystemLog{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
