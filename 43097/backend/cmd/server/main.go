package main

import (
	"hotel-system/internal/config"
	"hotel-system/internal/database"
	"hotel-system/internal/middleware"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/jwt"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/router"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	loggerConfig := logger.Config{
		LogLevel: config.AppConfig.Log.Level,
		LogFile:  config.AppConfig.Log.File,
	}
	if err := logger.Init(loggerConfig); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	logger.Info("Logger initialized successfully")

	jwt.Init(config.AppConfig.JWT.Secret, config.AppConfig.JWT.ExpireHours*60)
	logger.Info("JWT initialized successfully")

	database.InitDB()
	db := database.GetDB()

	err := db.AutoMigrate(
		&model.User{},
		&model.RoomType{},
		&model.Room{},
		&model.MemberLevel{},
		&model.Member{},
		&model.MemberPointsLog{},
		&model.Booking{},
		&model.CheckIn{},
		&model.Payment{},
	)
	if err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}
	logger.Info("Database migration completed")

	database.InitSeedData(db)
	logger.Info("Seed data initialization completed")

	gin.SetMode(config.AppConfig.Server.Mode)
	r := gin.New()

	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	router.SetupRouter(r)

	server := &http.Server{
		Addr:         ":" + config.AppConfig.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Infof("Server starting on port %s", config.AppConfig.Server.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
