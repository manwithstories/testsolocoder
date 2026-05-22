package router

import (
	"car-rental/internal/config"
	"car-rental/internal/handler"
	"car-rental/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())
	r.Use(middleware.Recovery())

	r.Static("/uploads", "./uploads")

	userHandler := handler.NewUserHandler(cfg)
	carHandler := handler.NewCarHandler(cfg)
	storeHandler := handler.NewStoreHandler()
	bookingHandler := handler.NewBookingHandler(&cfg.Email)
	orderHandler := handler.NewOrderHandler()
	reviewHandler := handler.NewReviewHandler()
	maintenanceHandler := handler.NewMaintenanceHandler()
	messageHandler := handler.NewMessageHandler(&cfg.Email)
	statsHandler := handler.NewStatsHandler()
	promoHandler := handler.NewPromoHandler()
	pricingHandler := handler.NewPricingHandler()

	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/refresh-token", userHandler.RefreshToken)

		api.GET("/cars", carHandler.GetAllCars)
		api.GET("/cars/available", carHandler.GetAvailableCars)
		api.GET("/cars/:id", carHandler.GetCarByID)
		api.GET("/cars/:id/images", carHandler.GetCarImages)
		api.GET("/cars/:id/reviews", reviewHandler.GetCarReviews)

		api.GET("/cities", storeHandler.GetAllCities)
		api.GET("/cities/:id/stores", storeHandler.GetStoresByCity)
		api.GET("/stores", storeHandler.GetAllStores)
		api.GET("/stores/:id", storeHandler.GetStoreByID)

		api.GET("/bookings/check-availability", bookingHandler.CheckCarAvailability)
		api.GET("/bookings/calculate-price", bookingHandler.CalculatePrice)

		api.GET("/promos/:code", promoHandler.GetPromoByCode)

		auth := api.Group("")
		auth.Use(middleware.Auth())
		{
			auth.GET("/profile", userHandler.GetProfile)
			auth.PUT("/profile", userHandler.UpdateProfile)
			auth.PUT("/password", userHandler.ChangePassword)
			auth.POST("/upload-license", userHandler.UploadLicense)
			auth.POST("/upload-avatar", userHandler.UploadAvatar)

			auth.POST("/bookings", bookingHandler.CreateBooking)
			auth.GET("/bookings", bookingHandler.GetUserBookings)
			auth.GET("/bookings/:id", bookingHandler.GetBookingByID)
			auth.GET("/bookings/no/:no", bookingHandler.GetBookingByNo)
			auth.PUT("/bookings/:id/cancel", bookingHandler.CancelBooking)

			auth.GET("/orders", orderHandler.GetUserOrders)
			auth.GET("/orders/:id", orderHandler.GetOrderByID)
			auth.GET("/orders/no/:no", orderHandler.GetOrderByNo)

			auth.POST("/reviews", reviewHandler.CreateReview)
			auth.GET("/my-reviews", reviewHandler.GetUserReviews)
			auth.PUT("/reviews/:id", reviewHandler.UpdateReview)
			auth.DELETE("/reviews/:id", reviewHandler.DeleteReview)
			auth.POST("/reviews/:id/like", reviewHandler.LikeReview)

			auth.GET("/messages", messageHandler.GetUserMessages)
			auth.GET("/messages/unread-count", messageHandler.GetUnreadCount)
			auth.PUT("/messages/:id/read", messageHandler.MarkAsRead)
			auth.PUT("/messages/read-all", messageHandler.MarkAllAsRead)
			auth.DELETE("/messages/:id", messageHandler.DeleteMessage)

			admin := auth.Group("")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/users", userHandler.GetAllUsers)
				admin.GET("/users/:id", userHandler.GetUserByID)
				admin.PUT("/users/:id/auth-status", userHandler.UpdateUserAuthStatus)
				admin.PUT("/users/:id/status", userHandler.UpdateUserStatus)
				admin.DELETE("/users/:id", userHandler.DeleteUser)

				admin.POST("/cars", carHandler.CreateCar)
				admin.PUT("/cars/:id", carHandler.UpdateCar)
				admin.PUT("/cars/:id/status", carHandler.UpdateCarStatus)
				admin.DELETE("/cars/:id", carHandler.DeleteCar)
				admin.POST("/cars/:id/upload", carHandler.UploadCarImage)
				admin.POST("/cars/:id/batch-upload", carHandler.BatchUploadCarImages)
				admin.DELETE("/cars/:id/images/:imageId", carHandler.DeleteCarImage)

				admin.POST("/cities", storeHandler.CreateCity)
				admin.PUT("/cities/:id", storeHandler.UpdateCity)
				admin.DELETE("/cities/:id", storeHandler.DeleteCity)

				admin.POST("/stores", storeHandler.CreateStore)
				admin.PUT("/stores/:id", storeHandler.UpdateStore)
				admin.DELETE("/stores/:id", storeHandler.DeleteStore)

				admin.GET("/admin/bookings", bookingHandler.GetAllBookings)
				admin.PUT("/admin/bookings/:id/confirm", bookingHandler.ConfirmBooking)
				admin.PUT("/admin/bookings/:id/complete", bookingHandler.CompleteBooking)

				admin.GET("/admin/orders", orderHandler.GetAllOrders)
				admin.PUT("/admin/orders/:id/status", orderHandler.UpdateOrderStatus)
				admin.PUT("/admin/orders/:id/refund", orderHandler.RefundOrder)
				admin.GET("/admin/orders/export", orderHandler.ExportOrders)

				admin.GET("/admin/reviews", reviewHandler.GetAllReviews)
				admin.PUT("/admin/reviews/:id/hidden", reviewHandler.ToggleReviewHidden)

				admin.POST("/maintenance", maintenanceHandler.CreateMaintenance)
				admin.GET("/maintenance", maintenanceHandler.GetAllMaintenance)
				admin.GET("/maintenance/:id", maintenanceHandler.GetMaintenanceByID)
				admin.GET("/cars/:carId/maintenance", maintenanceHandler.GetCarMaintenance)
				admin.PUT("/maintenance/:id", maintenanceHandler.UpdateMaintenance)
				admin.PUT("/maintenance/:id/start", maintenanceHandler.StartMaintenance)
				admin.PUT("/maintenance/:id/complete", maintenanceHandler.CompleteMaintenance)
				admin.PUT("/maintenance/:id/cancel", maintenanceHandler.CancelMaintenance)
				admin.DELETE("/maintenance/:id", maintenanceHandler.DeleteMaintenance)
				admin.GET("/maintenance/upcoming", maintenanceHandler.GetUpcomingMaintenance)

				admin.GET("/stats/dashboard", statsHandler.GetDashboardStats)
				admin.GET("/stats/revenue", statsHandler.GetRevenueStats)

				admin.POST("/promos", promoHandler.CreatePromo)
				admin.GET("/promos", promoHandler.GetAllPromos)
				admin.GET("/promos/:id", promoHandler.GetPromoByID)
				admin.PUT("/promos/:id", promoHandler.UpdatePromo)
				admin.DELETE("/promos/:id", promoHandler.DeletePromo)

				admin.POST("/pricing-rules", pricingHandler.CreateRule)
				admin.GET("/pricing-rules", pricingHandler.GetAllRules)
				admin.GET("/pricing-rules/:id", pricingHandler.GetRuleByID)
				admin.PUT("/pricing-rules/:id", pricingHandler.UpdateRule)
				admin.DELETE("/pricing-rules/:id", pricingHandler.DeleteRule)
				admin.PUT("/pricing-rules/:id/toggle", pricingHandler.ToggleRuleActive)
			}
		}
	}

	return r
}
