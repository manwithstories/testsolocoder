package main

import (
	"venue-booking/internal/api"
	"venue-booking/internal/api/middleware"
	"venue-booking/internal/config"
	"venue-booking/pkg/database"
	"venue-booking/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	logger.Init()
	database.Init()

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api.SetupRoutes(r)

	addr := config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	logger.Info("Server starting on ", addr)
	r.Run(addr)
}
