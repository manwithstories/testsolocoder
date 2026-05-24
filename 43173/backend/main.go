package main

import (
	"fmt"
	"log"

	"music-platform/internal/model"
	"music-platform/internal/middleware"
	"music-platform/internal/router"
	"music-platform/pkg/config"
	"music-platform/pkg/database"
	applogger "music-platform/pkg/logger"
	"music-platform/pkg/redis"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	if err := applogger.InitLogger(&cfg.Log); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer applogger.Sync()

	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer database.CloseDB()

	applogger.Info("Database connected successfully")

	if err := redis.InitRedis(&cfg.Redis); err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}
	defer redis.CloseRedis()

	applogger.Info("Redis connected successfully")

	err := database.AutoMigrate(
		&model.User{},
		&model.ArtistInfo{},
		&model.Work{},
		&model.Album{},
		&model.Tag{},
		&model.WorkTag{},
		&model.Copyright{},
		&model.Follow{},
		&model.Comment{},
		&model.Playlist{},
		&model.PlaylistWork{},
		&model.Like{},
		&model.PlayRecord{},
		&model.Share{},
		&model.Notification{},
		&model.Event{},
		&model.Ticket{},
		&model.Order{},
		&model.RevenueRecord{},
		&model.WithdrawRequest{},
		&model.Subscription{},
		&model.DailyStats{},
		&model.OperationLog{},
	)
	if err != nil {
		applogger.Error("Failed to auto migrate", zap.Error(err))
	} else {
		applogger.Info("Database migration completed")
	}

	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.OperationLog())

	router.SetupRouter(r)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	applogger.Info(fmt.Sprintf("Server starting on %s", addr))

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
