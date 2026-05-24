package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"repair-platform/internal/handlers"
	"repair-platform/internal/middleware"
	"repair-platform/internal/models"
	"repair-platform/internal/routes"
	"repair-platform/pkg/config"
	"repair-platform/pkg/logger"
	redisClient "repair-platform/pkg/redis"
)

func main() {
	logger.InitLogger()
	logger.Info("Starting Repair Platform Server...")

	cfg := config.LoadConfig()

	if err := models.InitDB(cfg.DatabasePath); err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	logger.Info("Database initialized successfully")

	models.InitSeedData()

	if err := redisClient.InitRedis(cfg); err != nil {
		logger.Warnf("Redis connection failed (optional): %v", err)
	} else {
		logger.Info("Redis connected successfully")
	}

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.RateLimitMiddleware(100, 60))

	routes.SetupRoutes(r)

	orderHandler := handlers.NewOrderHandler()
	partHandler := handlers.NewPartHandler()

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			orderHandler.CheckOrderTimeout()
			partHandler.CheckLowStock()
		}
	}()

	logger.Infof("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
