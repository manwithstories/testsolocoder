package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"temp-staff-platform/config"
	"temp-staff-platform/database"
	"temp-staff-platform/middleware"
	"temp-staff-platform/models"
	"temp-staff-platform/routes"
)

func main() {
	godotenv.Load()

	config.LoadConfig()

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(
		&models.User{},
		&models.JobPosting{},
		&models.JobApplication{},
		&models.Schedule{},
		&models.CheckIn{},
		&models.SalaryRecord{},
		&models.SalaryDetail{},
		&models.Evaluation{},
		&models.JobTemplate{},
		&models.Notification{},
	); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	os.MkdirAll(config.AppConfig.UploadDir+"/exports", 0755)

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	routes.SetupRoutes(r)

	addr := fmt.Sprintf(":%s", config.AppConfig.ServerPort)
	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
