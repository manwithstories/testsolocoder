package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"ticket-system/internal/config"
	"ticket-system/internal/database"
	"ticket-system/internal/handler"
	"ticket-system/internal/logging"
	"ticket-system/internal/middleware"
	"ticket-system/internal/redis"
)

func main() {
	config.LoadConfig()
	logging.InitLogger()
	database.InitDB()
	redis.InitRedis()

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	h := handler.NewHandler()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		user := api.Group("/user")
		user.Use(middleware.JWTAuth())
		{
			user.GET("/info", h.GetUserInfo)
			user.PUT("/info", h.UpdateUser)
			user.POST("/change-password", h.ChangePassword)
			user.GET("/member-levels", h.GetMemberLevels)
			user.POST("/exchange-coupon", h.ExchangeCoupon)
			user.GET("/coupons", h.GetUserCoupons)
		}

		shows := api.Group("/shows")
		{
			shows.GET("", h.ListShows)
			shows.GET("/:id", h.GetShow)
			shows.POST("", middleware.JWTAuth(), middleware.AdminAuth(), h.CreateShow)
			shows.PUT("/:id", middleware.JWTAuth(), middleware.AdminAuth(), h.UpdateShow)
			shows.DELETE("/:id", middleware.JWTAuth(), middleware.AdminAuth(), h.DeleteShow)
			shows.GET("/:showId/sessions", h.GetSessions)
			shows.POST("/sessions", middleware.JWTAuth(), middleware.AdminAuth(), h.CreateSession)
		}

		seatAreas := api.Group("/seat-areas")
		seatAreas.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			seatAreas.GET("/:sessionId", h.GetSeatAreas)
			seatAreas.POST("", h.CreateSeatArea)
		}

		seats := api.Group("/seats")
		{
			seats.GET("/:sessionId", h.GetSeats)
			seats.POST("", middleware.JWTAuth(), middleware.AdminAuth(), h.CreateSeat)
			seats.POST("/batch", middleware.JWTAuth(), middleware.AdminAuth(), h.BatchCreateSeats)
			seats.POST("/lock", middleware.JWTAuth(), h.LockSeats)
			seats.POST("/unlock", middleware.JWTAuth(), h.UnlockSeats)
		}

		seatCharts := api.Group("/seat-charts")
		{
			seatCharts.GET("/:sessionId", h.GetSeatChart)
			seatCharts.POST("", middleware.JWTAuth(), middleware.AdminAuth(), h.UpdateSeatChart)
		}

		orders := api.Group("/orders")
		orders.Use(middleware.JWTAuth())
		{
			orders.GET("", h.ListOrders)
			orders.GET("/:orderNo", h.GetOrder)
			orders.POST("", h.CreateOrder)
			orders.POST("/pay", h.PayOrder)
			orders.POST("/:orderNo/cancel", h.CancelOrder)
			orders.POST("/refund", h.RequestRefund)
			orders.POST("/refund/audit", middleware.AdminAuth(), h.AuditRefund)
			orders.GET("/export/excel", middleware.AdminAuth(), h.ExportOrders)
		}

		checkin := api.Group("/checkin")
		checkin.Use(middleware.JWTAuth())
		{
			checkin.POST("", h.Checkin)
		}

		statistics := api.Group("/statistics")
		statistics.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			statistics.GET("/sales", h.GetSalesStatistics)
			statistics.GET("/area-sales/:sessionId", h.GetAreaSales)
			statistics.GET("/seat-heatmap/:sessionId", h.GetSeatHeatmap)
			statistics.GET("/audience-profile", h.GetAudienceProfile)
			statistics.GET("/daily-sales", h.GetDailySales)
			statistics.GET("/export/pdf", h.ExportStatisticsPDF)
		}
	}

	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	r.Run(addr)
}
