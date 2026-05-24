package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"errand-service/internal/config"
	"errand-service/internal/handlers"
	"errand-service/internal/middleware"
	"errand-service/pkg/database"
	"errand-service/pkg/logger"
	"errand-service/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(cfg.LogLevel)
	logger.Info("Starting errand-service...")

	db, err := database.Init(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.AutoMigrate(db); err != nil {
		logger.Fatalf("Failed to auto migrate: %v", err)
	}

	redisClient, err := redis.Init(cfg.Redis)
	if err != nil {
		logger.Fatalf("Failed to connect to redis: %v", err)
	}
	defer redis.Close()

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())

	v1 := r.Group("/api/v1")
	{
		authHandler := handlers.NewAuthHandler(db)
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/refresh", authHandler.RefreshToken)

		authGroup := v1.Group("")
		authGroup.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			userHandler := handlers.NewUserHandler(db)
			authGroup.GET("/user/profile", userHandler.GetProfile)
			authGroup.PUT("/user/profile", userHandler.UpdateProfile)
			authGroup.POST("/user/verify", userHandler.SubmitVerification)

			courierHandler := handlers.NewCourierHandler(db)
			authGroup.POST("/courier/apply", courierHandler.Apply)
			authGroup.GET("/courier/tasks", courierHandler.GetMyTasks)

			taskHandler := handlers.NewTaskHandler(db, redisClient)
			authGroup.POST("/tasks", taskHandler.Create)
			authGroup.GET("/tasks", taskHandler.List)
			authGroup.GET("/tasks/:id", taskHandler.Get)
			authGroup.PUT("/tasks/:id", taskHandler.Update)
			authGroup.DELETE("/tasks/:id", taskHandler.Delete)
			authGroup.POST("/tasks/:id/accept", taskHandler.Accept)
			authGroup.POST("/tasks/:id/complete", taskHandler.Complete)
			authGroup.POST("/tasks/:id/cancel", taskHandler.Cancel)
			authGroup.GET("/tasks/nearby", taskHandler.GetNearby)

			orderHandler := handlers.NewOrderHandler(db, redisClient)
			authGroup.GET("/orders", orderHandler.List)
			authGroup.GET("/orders/:id", orderHandler.Get)
			authGroup.POST("/orders/:id/track", orderHandler.Track)

			paymentHandler := handlers.NewPaymentHandler(db)
			authGroup.POST("/payments/deposit", paymentHandler.Deposit)
			authGroup.POST("/payments/withdraw", paymentHandler.Withdraw)
			authGroup.POST("/payments/refund", paymentHandler.Refund)
			authGroup.GET("/payments/history", paymentHandler.History)

			reviewHandler := handlers.NewReviewHandler(db)
			authGroup.POST("/reviews", reviewHandler.Create)
			authGroup.GET("/reviews", reviewHandler.List)

			wsHandler := handlers.NewWebSocketHandler(db)
			authGroup.GET("/ws/:orderId", wsHandler.HandleConnection)
		}

		adminGroup := v1.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		adminGroup.Use(middleware.RoleMiddleware("admin"))
		{
			adminHandler := handlers.NewAdminHandler(db)
			adminGroup.GET("/users", adminHandler.ListUsers)
			adminGroup.PUT("/users/:id/freeze", adminHandler.FreezeUser)
			adminGroup.PUT("/users/:id/unfreeze", adminHandler.UnfreezeUser)
			adminGroup.PUT("/couriers/:id/approve", adminHandler.ApproveCourier)
			adminGroup.PUT("/couriers/:id/reject", adminHandler.RejectCourier)
		}
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Infof("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-quit.Done()

	logger.Info("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server stopped")
}
