package api

import (
	"venue-booking/internal/api/handler"
	"venue-booking/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	venueHandler := handler.NewVenueHandler()
	deviceHandler := handler.NewDeviceHandler()
	orderHandler := handler.NewOrderHandler()
	paymentHandler := handler.NewPaymentHandler()
	reviewHandler := handler.NewReviewHandler()
	statsHandler := handler.NewStatsHandler()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/send-verify-email", authHandler.SendVerifyEmail)
			auth.POST("/verify-email", authHandler.VerifyEmail)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)

			auth.Use(middleware.JWTAuth())
			auth.POST("/logout", authHandler.Logout)
		}

		users := api.Group("/users")
		users.Use(middleware.JWTAuth())
		{
			users.GET("/me", userHandler.GetMe)
			users.PUT("/me", userHandler.UpdateMe)

			users.Use(middleware.AdminRequired())
			users.GET("", userHandler.ListUsers)

			users.Use(middleware.SuperAdminRequired())
			users.PUT("/:id/role", userHandler.UpdateRole)
		}

		venues := api.Group("/venues")
		{
			venues.GET("", venueHandler.List)
			venues.GET("/:id", venueHandler.GetByID)
			venues.GET("/:id/availability", venueHandler.GetAvailability)

			venues.Use(middleware.JWTAuth(), middleware.AdminRequired())
			venues.POST("", venueHandler.Create)
			venues.PUT("/:id", venueHandler.Update)
			venues.DELETE("/:id", venueHandler.Delete)
			venues.PATCH("/:id/status", venueHandler.UpdateStatus)
			venues.POST("/:id/prices", venueHandler.SetPrice)
		}

		devices := api.Group("/devices")
		{
			devices.GET("/categories", deviceHandler.ListCategories)
			devices.GET("", deviceHandler.List)
			devices.GET("/:id", deviceHandler.GetByID)
			devices.GET("/:id/availability", deviceHandler.GetAvailability)

			devices.Use(middleware.JWTAuth(), middleware.AdminRequired())
			devices.POST("/categories", deviceHandler.CreateCategory)
			devices.POST("", deviceHandler.Create)
			devices.PUT("/:id", deviceHandler.Update)
			devices.PATCH("/:id/status", deviceHandler.UpdateStatus)
			devices.POST("/batch-import", deviceHandler.BatchImport)
		}

		bookings := api.Group("/bookings")
		bookings.Use(middleware.JWTAuth())
		{
			bookings.POST("", orderHandler.Create)
			bookings.GET("/calendar", orderHandler.GetCalendar)
		}

		orders := api.Group("/orders")
		orders.Use(middleware.JWTAuth())
		{
			orders.GET("", orderHandler.List)
			orders.GET("/:id", orderHandler.GetByID)
			orders.PUT("/:id/cancel", orderHandler.Cancel)

			orders.Use(middleware.AdminRequired())
			orders.PUT("/:id/confirm", orderHandler.Confirm)
			orders.PUT("/:id/complete", orderHandler.Complete)
		}

		payments := api.Group("/payments")
		payments.Use(middleware.JWTAuth(), middleware.AdminRequired())
		{
			payments.GET("", paymentHandler.List)
			payments.POST("/:id/confirm", paymentHandler.ConfirmPayment)
			payments.GET("/export", paymentHandler.Export)
		}

		reviews := api.Group("/reviews")
		reviews.Use(middleware.JWTAuth())
		{
			reviews.POST("", reviewHandler.Create)
			reviews.GET("", reviewHandler.List)
			reviews.GET("/:id", reviewHandler.GetByID)

			reviews.Use(middleware.AdminRequired())
			reviews.PUT("/:id/approve", reviewHandler.Approve)
			reviews.PUT("/:id/reject", reviewHandler.Reject)
		}

		stats := api.Group("/stats")
		stats.Use(middleware.JWTAuth(), middleware.AdminRequired())
		{
			stats.GET("/overview", statsHandler.GetOverview)
			stats.GET("/bookings", statsHandler.GetBookingStats)
			stats.GET("/revenue", statsHandler.GetRevenueStats)
			stats.GET("/popular-venues", statsHandler.GetPopularVenues)
		}
	}
}
