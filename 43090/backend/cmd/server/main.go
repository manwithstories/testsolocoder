package main

import (
	"fmt"
	"log"

	"auction-system/config"
	"auction-system/internal/models"
	"auction-system/internal/routers"
	"auction-system/pkg/logger"
	"auction-system/pkg/redis"
)

func main() {
	config.LoadConfig()

	logFileName := fmt.Sprintf("./logs/%s", logger.GetLogFileName())
	if err := logger.InitLogger(logFileName); err != nil {
		log.Printf("Warning: Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	if err := models.InitDB(); err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}

	if err := redis.InitRedis(); err != nil {
		logger.Fatal("Failed to initialize Redis: %v", err)
	}

	r := routers.SetupRouter()

	addr := fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	logger.Info("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}
