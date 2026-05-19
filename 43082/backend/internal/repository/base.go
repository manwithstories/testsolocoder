package repository

import (
	"gym-management/internal/pkg/database"

	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	return database.GetDB()
}
