package seeds

import (
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"sports-league/models"
)

func SeedAdmin(db *gorm.DB) {
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		return
	}
	var count int64
	db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	admin := models.User{
		Email:    email,
		Password: string(hashed),
		FullName: "System Admin",
		Role:     models.RoleAdmin,
		IsActive: true,
	}
	db.Create(&admin)
}
