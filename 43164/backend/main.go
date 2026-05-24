package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"tutoring-platform/config"
	"tutoring-platform/database"
	"tutoring-platform/middleware"
	"tutoring-platform/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	routes.SetupRoutes(r)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("Database connected: %s", cfg.DB.Name)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	_ = db
}
