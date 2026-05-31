package routes

import (
	"secondhand-platform/controllers"
	"secondhand-platform/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	userCtrl := controllers.NewUserController()
	productCtrl := controllers.NewProductController()
	repairCtrl := controllers.NewRepairController()
	orderCtrl := controllers.NewOrderController()
	reviewCtrl := controllers.NewReviewController()
	reportCtrl := controllers.NewReportController()
	warrantyCtrl := controllers.NewWarrantyController()
	notificationCtrl := controllers.NewNotificationController()
	messageCtrl := controllers.NewMessageController()
	adminCtrl := controllers.NewAdminController()

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", userCtrl.Register)
		api.POST("/auth/login", userCtrl.Login)
		api.POST("/auth/refresh", userCtrl.RefreshToken)
		api.POST("/auth/logout", middlewares.JWTAuth(), userCtrl.Logout)

		api.GET("/products", productCtrl.ListProducts)
		api.GET("/products/hot", productCtrl.GetHotProducts)
		api.GET("/products/categories", productCtrl.GetCategories)
		api.GET("/products/:id", productCtrl.GetProduct)

		api.GET("/services", repairCtrl.ListServices)
		api.GET("/services/types", repairCtrl.GetServiceTypes)
		api.GET("/services/:id", repairCtrl.GetService)

		api.GET("/reviews", reviewCtrl.ListReviews)
		api.GET("/reviews/:id", reviewCtrl.GetReview)
		api.GET("/users/:user_id/rating", reviewCtrl.GetAverageRating)

		auth := api.Group("")
		auth.Use(middlewares.JWTAuth())
		{
			auth.GET("/user/profile", userCtrl.GetProfile)
			auth.PUT("/user/profile", userCtrl.UpdateProfile)
			auth.PUT("/user/password", userCtrl.ChangePassword)

			auth.POST("/user/realname-auth", userCtrl.SubmitRealNameAuth)
			auth.POST("/user/technician-cert", userCtrl.SubmitTechnicianCert)

			auth.GET("/user/wallet/balance", userCtrl.GetWalletBalance)
			auth.POST("/user/wallet/recharge", userCtrl.Recharge)
			auth.POST("/user/wallet/withdraw", userCtrl.Withdraw)
			auth.GET("/user/wallet/logs", userCtrl.ListWalletLogs)

			auth.GET("/user/stats", userCtrl.GetUserStats)

			auth.POST("/upload", controllers.UploadImage)
			auth.POST("/upload/multiple", controllers.UploadMultipleImages)

			seller := auth.Group("")
			seller.Use(middlewares.JWTAuth("seller"))
			{
				seller.POST("/products", productCtrl.CreateProduct)
				seller.PUT("/products/:id", productCtrl.UpdateProduct)
				seller.DELETE("/products/:id", productCtrl.DeleteProduct)
				seller.POST("/products/:id/off-shelf", productCtrl.OffShelfProduct)
				seller.GET("/my/products", productCtrl.ListMyProducts)

				seller.POST("/orders/:order_no/ship", orderCtrl.ShipOrder)
				seller.POST("/orders/:order_no/cancel", orderCtrl.CancelOrder)
				seller.POST("/orders/negotiation/handle", orderCtrl.HandleNegotiation)

				seller.GET("/seller/orders", orderCtrl.ListOrders)
			}

			technician := auth.Group("")
			technician.Use(middlewares.JWTAuth("technician"))
			{
				technician.POST("/services", repairCtrl.CreateService)
				technician.PUT("/services/:id", repairCtrl.UpdateService)
				technician.DELETE("/services/:id", repairCtrl.DeleteService)

				technician.POST("/repair-orders/:id/accept", repairCtrl.AcceptRepairOrder)
				technician.POST("/repair-orders/:id/start", repairCtrl.StartRepair)
				technician.POST("/repair-orders/:id/complete", repairCtrl.CompleteRepair)

				technician.GET("/technician/repair-orders", repairCtrl.ListRepairOrders)
			}

			buyer := auth.Group("")
			buyer.Use(middlewares.JWTAuth("buyer"))
			{
				buyer.POST("/orders", orderCtrl.CreateOrder)
				buyer.POST("/orders/:order_no/pay", orderCtrl.PayOrder)
				buyer.POST("/orders/:order_no/confirm", orderCtrl.ConfirmDelivery)
				buyer.POST("/orders/:order_no/cancel", orderCtrl.CancelOrder)
				buyer.POST("/orders/:order_no/refund", orderCtrl.RefundOrder)
				buyer.POST("/orders/negotiate", orderCtrl.NegotiatePrice)

				buyer.GET("/buyer/orders", orderCtrl.ListOrders)
				buyer.GET("/orders/:id", orderCtrl.GetOrder)
				buyer.GET("/orders/by-no/:order_no", orderCtrl.GetOrderByNo)
				buyer.GET("/order/stats", orderCtrl.GetOrderStats)

				buyer.POST("/repair-orders", repairCtrl.CreateRepairOrder)
				buyer.POST("/repair-orders/:id/pickup", repairCtrl.PickUpDevice)
				buyer.POST("/repair-orders/:id/cancel", repairCtrl.CancelRepairOrder)

				buyer.GET("/buyer/repair-orders", repairCtrl.ListRepairOrders)
			}

			auth.POST("/favorites/:id", productCtrl.ToggleFavorite)
			auth.GET("/favorites", productCtrl.ListFavorites)

			auth.POST("/reviews", reviewCtrl.CreateReview)
			auth.DELETE("/reviews/:id", reviewCtrl.DeleteReview)

			auth.POST("/reports", reportCtrl.CreateReport)

			auth.POST("/warranties", warrantyCtrl.CreateWarranty)
			auth.GET("/warranties", warrantyCtrl.ListWarranties)
			auth.GET("/warranties/:id", warrantyCtrl.GetWarranty)

			auth.GET("/notifications", notificationCtrl.ListNotifications)
			auth.GET("/notifications/:id", notificationCtrl.GetNotification)
			auth.POST("/notifications/:id/read", notificationCtrl.MarkAsRead)
			auth.POST("/notifications/read-all", notificationCtrl.MarkAllAsRead)
			auth.GET("/notifications/unread-count", notificationCtrl.GetUnreadCount)
			auth.DELETE("/notifications/:id", notificationCtrl.DeleteNotification)

			auth.POST("/messages", messageCtrl.SendMessage)
			auth.GET("/messages/:user_id", messageCtrl.ListMessages)
			auth.POST("/messages/:user_id/read", messageCtrl.MarkMessagesAsRead)
			auth.GET("/messages/unread-count", messageCtrl.GetUnreadMessageCount)
			auth.GET("/message-contacts", messageCtrl.GetMessageContacts)
		}

		admin := api.Group("/admin")
		admin.Use(middlewares.JWTAuth("admin"))
		{
			admin.GET("/dashboard/stats", adminCtrl.GetDashboardStats)
			admin.GET("/transaction-stats", adminCtrl.GetTransactionStats)
			admin.GET("/user-activity-stats", adminCtrl.GetUserActivityStats)

			admin.GET("/users", userCtrl.ListUsers)
			admin.PUT("/users/:id/status", userCtrl.UpdateUserStatus)
			admin.POST("/users/:id/review-realname", userCtrl.ReviewRealNameAuth)
			admin.POST("/technician-certs/:id/review", userCtrl.ReviewTechnicianCert)

			admin.GET("/products/pending", productCtrl.ListPendingProducts)
			admin.POST("/products/:id/review", productCtrl.ReviewProduct)

			admin.GET("/reports", reportCtrl.ListReports)
			admin.GET("/reports/:id", reportCtrl.GetReport)
			admin.POST("/reports/:id/handle", reportCtrl.HandleReport)

			admin.GET("/warranties", warrantyCtrl.ListWarranties)
			admin.POST("/warranties/:id/handle", warrantyCtrl.HandleWarranty)

			admin.GET("/transactions", adminCtrl.ListAllTransactions)

			admin.POST("/orders/refund/handle", orderCtrl.HandleRefund)

			admin.GET("/admin-logs", adminCtrl.ListAdminLogs)
		}
	}
}
