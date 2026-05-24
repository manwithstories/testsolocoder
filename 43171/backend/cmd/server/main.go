package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"drone-rental/internal/config"
	"drone-rental/internal/middleware"
	"drone-rental/internal/pkg/logger"
	"drone-rental/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if _, err := logger.New(cfg.Log); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.L.Close()

	if err := config.InitDB(cfg.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	r.Use(middleware.CORS())
	r.Use(logger.GinLogger())
	r.Use(gin.Recovery())

	if err := os.MkdirAll(cfg.Upload.SavePath, 0755); err != nil {
		log.Fatalf("Failed to create upload dir: %v", err)
	}

	router.RegisterRoutes(r)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
