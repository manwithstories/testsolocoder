package main

import (
	"log"

	"pet-adoption-platform/config"
	"pet-adoption-platform/database"
	"pet-adoption-platform/models"
	"pet-adoption-platform/routes"
	"pet-adoption-platform/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	utils.InitLogger()

	cfg := config.Load()

	utils.SetJWTSecret(cfg.JWTSecret)

	if err := database.Init(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	utils.Logger.Info("Database connected successfully")

	err := database.AutoMigrate(
		&models.User{},
		&models.RescueStation{},
		&models.Pet{},
		&models.AdoptionApplication{},
		&models.AdoptionAgreement{},
		&models.FollowUpRecord{},
		&models.HealthRecord{},
		&models.HealthReminder{},
		&models.Appointment{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	utils.Logger.Info("Database migration completed")

	createDefaultAdmin()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.SetupRouter(r)

	addr := ":" + cfg.ServerPort
	utils.Logger.Infof("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createDefaultAdmin() {
	var count int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)
	if count > 0 {
		return
	}

	admin := &models.User{
		Email:    "admin@petadoption.com",
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		Name:     "系统管理员",
		Role:     models.RoleAdmin,
		IsVerified: true,
	}

	if err := database.DB.Create(admin).Error; err != nil {
		utils.Logger.Warn("Failed to create default admin: %v", err)
	} else {
		utils.Logger.Info("Default admin created: admin@petadoption.com / admin123")
	}
}
