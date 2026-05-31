package main

import (
	"fmt"
	"log"
	"os"

	"coffee-platform/config"
	"coffee-platform/database"
	"coffee-platform/middleware"
	"coffee-platform/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default config")
	}

	cfg := config.LoadConfig()

	gin.SetMode(cfg.Server.Mode)

	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := os.MkdirAll(cfg.App.UploadDir, 0755); err != nil {
		log.Printf("Warning: Failed to create upload directory: %v", err)
	}

	r := gin.New()

	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.OperationLogMiddleware())

	routes.SetupRoutes(r, cfg)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("API Base URL: http://localhost%s/api", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
