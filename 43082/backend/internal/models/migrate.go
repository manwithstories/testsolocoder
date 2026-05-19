package models

import (
	"gym-management/internal/pkg/database"
	"gym-management/internal/pkg/logger"

	"go.uber.org/zap"
)

func AutoMigrate() {
	db := database.GetDB()

	err := db.AutoMigrate(
		&Member{},
		&Membership{},
		&Coach{},
		&Course{},
		&CourseSchedule{},
		&Booking{},
		&Waitlist{},
		&CheckIn{},
		&Reminder{},
		&OperationLog{},
	)

	if err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	logger.Info("Database migration completed successfully")
}
