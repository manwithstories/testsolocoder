package main

import (
	"ticket-system/config"
	"ticket-system/internal/logger"
	"ticket-system/internal/middleware"
	"ticket-system/internal/models"
	"ticket-system/internal/redis"
	"ticket-system/internal/router"
	"ticket-system/internal/seeder"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	logger.Init()
	models.InitDB()
	redis.Init()
	seeder.Run()

	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.RequestID())

	router.Register(r)

	logger.Log.Infof("Server starting on %s:%s", config.App.ServerHost, config.App.ServerPort)
	if err := r.Run(config.App.ServerHost + ":" + config.App.ServerPort); err != nil {
		logger.Log.Fatalf("Failed to start server: %v", err)
	}
}
