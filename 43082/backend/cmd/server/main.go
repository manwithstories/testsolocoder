package main

import (
	"log"
	"multishop/internal/config"
	"multishop/internal/database"
	"multishop/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := database.Init(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	r := gin.Default()

	routes.Setup(r, cfg)

	log.Printf("Server starting on %s:%s", cfg.ServerHost, cfg.ServerPort)
	if err := r.Run(cfg.ServerHost + ":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
