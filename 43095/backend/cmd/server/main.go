package main

import (
	"medical-platform/internal/config"
	"medical-platform/internal/handlers"
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/migrations"
	"medical-platform/pkg/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	if _, err := config.Load(); err != nil {
		panic("Failed to load config: " + err.Error())
	}
	defer config.Logger.Sync()

	db, err := database.Init()
	if err != nil {
		config.Logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	if err := migrations.RunMigrations(db); err != nil {
		config.Logger.Fatal("Failed to run migrations", zap.Error(err))
	}
	config.Logger.Info("Migrations completed successfully")

	scheduler := services.NewScheduler()
	scheduler.Start()
	defer scheduler.Stop()

	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")
	{
		handlers.RegisterAuthRoutes(api)
		handlers.RegisterDepartmentRoutes(api)
		handlers.RegisterDoctorRoutes(api)
		handlers.RegisterPatientRoutes(api)
		handlers.RegisterAppointmentRoutes(api)
		handlers.RegisterConsultationRoutes(api)
		handlers.RegisterHealthRecordRoutes(api)
		handlers.RegisterNotificationRoutes(api)
		handlers.RegisterPaymentRoutes(api)
		handlers.RegisterReviewRoutes(api)
		handlers.RegisterUploadRoutes(api)
	}

	addr := config.AppConfig.ServerHost + ":" + config.AppConfig.ServerPort
	config.Logger.Info("Server starting", zap.String("address", addr))
	if err := r.Run(addr); err != nil {
		config.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
