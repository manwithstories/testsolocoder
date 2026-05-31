package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ship-rental-platform/internal/config"
	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/middleware"
	"ship-rental-platform/internal/model"
	"ship-rental-platform/internal/routes"
	"ship-rental-platform/internal/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	middleware.JWTSecret = []byte(cfg.JWT.Secret)

	db, err := database.InitDatabase(&cfg.Database)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	err = database.AutoMigrate(db,
		&model.User{},
		&model.Ship{},
		&model.ShipImage{},
		&model.Dock{},
		&model.Berth{},
		&model.BerthReservation{},
		&model.WaterLevel{},
		&model.Rental{},
		&model.PaymentRecord{},
		&model.VoyageLog{},
		&model.MaintenanceRecord{},
		&model.MaintenanceSchedule{},
		&model.Transaction{},
		&model.Settlement{},
		&model.Review{},
		&model.AuditLog{},
	)
	if err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	logrus.Info("Database migrations completed")

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLogger())

	routes.SetupRoutes(r)

	scheduler.InitScheduler()
	defer scheduler.StopScheduler()

	go func() {
		addr := cfg.Server.Port
		logrus.Infof("Server starting on port %s", addr)
		if err := r.Run(addr); err != nil {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")
	logrus.Info("Server shutdown complete")
}
