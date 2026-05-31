package main

import (
	"fmt"
	"log"

	"consultation-platform/config"
	"consultation-platform/middlewares"
	"consultation-platform/routes"
	"consultation-platform/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := utils.InitLogger(cfg.Server.Mode)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer utils.SyncLogger()

	_, err = utils.InitDatabase(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
	}

	_, err = utils.InitRedis(cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to init redis", zap.Error(err))
	}

	gin.SetMode(cfg.Server.Mode)

	router := gin.New()

	router.Use(middlewares.RecoveryMiddleware(logger))
	router.Use(gin.Logger())
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.RateLimitMiddleware(cfg.RateLimit))

	routes.SetupRoutes(router, cfg)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Info("Server starting", zap.String("addr", addr))

	if err := router.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
