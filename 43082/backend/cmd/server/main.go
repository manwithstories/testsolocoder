package main

import (
	"gym-management/config"
	"gym-management/internal/api"
	"gym-management/internal/middleware"
	"gym-management/internal/models"
	"gym-management/internal/pkg/database"
	"gym-management/internal/pkg/logger"
	"gym-management/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()

	logger.InitLogger(config.AppConfig.Log.Level, config.AppConfig.Log.File)
	defer logger.Sync()

	database.InitDB()

	models.AutoMigrate()

	gin.SetMode(config.AppConfig.Server.Mode)
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.OperationLogMiddleware())

	apiV1 := r.Group("/api/v1")
	{
		api.NewMemberHandler().RegisterRoutes(apiV1)
		api.NewMembershipHandler().RegisterRoutes(apiV1)
		api.NewCourseHandler().RegisterRoutes(apiV1)
		api.NewBookingHandler().RegisterRoutes(apiV1)
		api.NewCheckInHandler().RegisterRoutes(apiV1)
		api.NewStatsHandler().RegisterRoutes(apiV1)
		api.NewCoachHandler().RegisterRoutes(apiV1)
	}

	go startReminderScheduler()

	log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
	if err := r.Run(":" + config.AppConfig.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startReminderScheduler() {
	c := cron.New()

	reminderService := service.NewReminderService()

	_, err := c.AddFunc("0 0 9 * * *", func() {
		logger.Info("Running daily reminder generation job")
		_ = reminderService.SendReminders()
	})
	if err != nil {
		logger.Error("Failed to add daily reminder cron job", zap.Error(err))
	}

	_, err = c.AddFunc("0 */30 * * * *", func() {
		logger.Info("Running pending reminder processing job")
		_ = reminderService.ProcessPendingReminders()
	})
	if err != nil {
		logger.Error("Failed to add reminder processing cron job", zap.Error(err))
	}

	c.Start()
	logger.Info("Reminder scheduler started")

	select {}
}
