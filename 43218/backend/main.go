package main

import (
	"flag"
	"fmt"
	"os"

	"secondhand-platform/cache"
	"secondhand-platform/config"
	"secondhand-platform/database"
	"secondhand-platform/middlewares"
	"secondhand-platform/models"
	"secondhand-platform/routes"
	"secondhand-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	utils.InitLogger(&cfg.Log)

	if err := database.InitDB(&cfg.Database.MySQL); err != nil {
		logrus.Fatalf("Failed to init database: %v", err)
	}

	if err := cache.InitRedis(&cfg.Database.Redis); err != nil {
		logrus.Fatalf("Failed to init redis: %v", err)
	}

	database.AutoMigrate(
		&models.User{},
		&models.TechnicianCert{},
		&models.Product{},
		&models.ProductImage{},
		&models.Favorite{},
		&models.RepairService{},
		&models.RepairOrder{},
		&models.Order{},
		&models.Negotiation{},
		&models.WalletLog{},
		&models.Review{},
		&models.Report{},
		&models.Warranty{},
		&models.Notification{},
		&models.Message{},
		&models.Transaction{},
		&models.AdminLog{},
	)

	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	r.Use(middlewares.RecoveryMiddleware())
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.CORSMiddleware())

	r.Static("/uploads", "./uploads")

	routes.SetupRoutes(r)

	utils.LogInfo("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped")
}
