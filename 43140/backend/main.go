package main

import (
	"log"
	"os"

	"recruitment-platform/config"
	"recruitment-platform/database"
	"recruitment-platform/models"
	"recruitment-platform/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.JobSeeker{},
		&models.Department{},
		&models.Job{},
		&models.Resume{},
		&models.Education{},
		&models.WorkExperience{},
		&models.Skill{},
		&models.Application{},
		&models.Interview{},
		&models.Review{},
		&models.JobView{},
		&models.Notification{},
	)

	r := gin.Default()

	routes.SetupRoutes(r, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
