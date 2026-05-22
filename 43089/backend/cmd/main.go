package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"travel-planner/internal/config"
	"travel-planner/internal/database"
	"travel-planner/internal/logger"
	"travel-planner/internal/middleware"
	"travel-planner/internal/routes"
	"travel-planner/internal/services"
	"travel-planner/internal/utils"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	utils.InitValidator()

	emailService := services.NewEmailService(&config.AppConfig.Email)
	reminderScheduler := services.GetReminderScheduler(emailService)
	reminderScheduler.Start()
	defer reminderScheduler.Stop()

	gin.SetMode(config.AppConfig.Server.Mode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.AuditLog())

	routes.SetupRoutes(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	logger.Infof("Server starting on %s", addr)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := r.Run(addr); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down server...")
	reminderScheduler.Stop()
	logger.Info("Server shutdown complete")
	os.Exit(0)
}
