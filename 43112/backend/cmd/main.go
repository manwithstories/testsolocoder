package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/database"
	"e-learning-platform/internal/middleware"
	"e-learning-platform/internal/routes"
	"e-learning-platform/internal/utils"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	utils.InitLogger(cfg.Log.Level, cfg.Log.LogDir)
	utils.InitJWT(cfg.JWT.Secret)

	if err := os.MkdirAll(cfg.Upload.UploadDir+"/image", 0755); err != nil {
		log.Fatalf("Failed to create upload dir: %v", err)
	}
	os.MkdirAll(cfg.Upload.UploadDir+"/video", 0755)
	os.MkdirAll(cfg.Upload.UploadDir+"/document", 0755)

	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	r.Static("/uploads", cfg.Upload.UploadDir)

	routes.SetupRoutes(r, cfg)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	utils.Logger.Infof("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
