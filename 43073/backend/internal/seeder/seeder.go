package seeder

import (
	"ticket-system/internal/logger"
	"ticket-system/internal/models"
	"ticket-system/internal/util"
)

func Run() {
	seedAdmin()
}

func seedAdmin() {
	var count int64
	models.DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, err := util.HashPassword("admin123")
	if err != nil {
		logger.Log.Errorf("Hash admin password failed: %v", err)
		return
	}

	admin := &models.User{
		Username: "admin",
		Password: hashedPassword,
		Email:    "admin@example.com",
		Role:     models.RoleAdmin,
	}

	if err := models.DB.Create(admin).Error; err != nil {
		logger.Log.Errorf("Create admin failed: %v", err)
		return
	}

	logger.Log.Info("Admin user created: admin/admin123")
}
