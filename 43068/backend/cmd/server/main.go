package main

import (
	"log"

	"freelancer-management/internal/config"
	"freelancer-management/internal/database"
	"freelancer-management/internal/logger"
	"freelancer-management/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger()
	database.InitDB()

	gin.SetMode(cfg.Server.Mode)

	r := router.SetupRouter()

	logger.LogInfo("Server starting on port %s in %s mode", cfg.Server.Port, cfg.Server.Mode)
	log.Printf("Server starting on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
