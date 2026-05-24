package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"business-registration-platform/config"
	"business-registration-platform/database"
	"business-registration-platform/middleware"
	"business-registration-platform/routes"
	"business-registration-platform/utils"
)

func main() {
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logFile, err := utils.InitLogger()
	if err != nil {
		log.Printf("Warning: Failed to initialize logger: %v", err)
	}
	if logFile != nil {
		defer logFile.Close()
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	gin.SetMode(config.AppConfig.Server.Mode)

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.LoggingMiddleware())

	routes.SetupRoutes(r)

	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
