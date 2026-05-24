package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"luxury-trading-platform/internal/cache"
	"luxury-trading-platform/internal/config"
	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/router"
	"luxury-trading-platform/internal/utils/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	log := logger.InitLogger(cfg.Log.Level)
	log.Info("Logger initialized")

	db, err := model.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	log.Info("Database connected successfully")

	rdb, err := cache.InitRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}
	log.Info("Redis connected successfully")

	gin.SetMode(cfg.Server.Mode)
	r := router.SetupRouter(db, rdb, cfg, log)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Infof("Server listening on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}
