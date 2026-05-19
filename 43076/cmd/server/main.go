package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ticket-system/internal/config"
	"ticket-system/internal/database"
	"ticket-system/internal/middleware"
	"ticket-system/internal/router"
	"ticket-system/internal/scheduler"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.Load("./configs/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	utils.InitValidator()

	if err := database.Init(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	log.Println("Database initialized successfully")

	scheduler.InitScheduler()

	gin.SetMode(config.AppConfig.Server.Mode)
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(gin.Recovery())

	router.SetupRoutes(r)

	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	scheduler.StopScheduler()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
