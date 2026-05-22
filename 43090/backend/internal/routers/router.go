package routers

import (
	"github.com/gin-gonic/gin"
	"auction-system/internal/controllers"
	"auction-system/internal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.Static("/uploads", "./uploads")
	r.Static("/exports", "./exports")

	userCtrl := controllers.NewUserController()
	itemCtrl := controllers.NewAuctionItemController()
	sessionCtrl := controllers.NewAuctionSessionController()
	bidCtrl := controllers.NewBidController()
	orderCtrl := controllers.NewOrderController()
	notifyCtrl := controllers.NewNotificationController()
	reviewCtrl := controllers.NewReviewController()
	statsCtrl := controllers.NewStatisticsController()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userCtrl.Register)
			auth.POST("/login", userCtrl.Login)
		}

		items := api.Group("/items")
		{
			items.GET("", itemCtrl.GetList)
			items.GET("/:id", itemCtrl.GetByID)
			items.GET("/:id/bids", bidCtrl.GetBidHistory)
			items.GET("/:id/current-bid", bidCtrl.GetCurrentBid)

			items.Use(middleware.JWTAuth())
			{
				items.POST("", itemCtrl.Create)
				items.PUT("/:id", itemCtrl.Update)
				items.DELETE("/:id", itemCtrl.Delete)
				items.POST("/:id/online", itemCtrl.Online)
				items.POST("/:id/offline", itemCtrl.Offline)
				items.POST("/:id/images", itemCtrl.UploadImages)
				items.POST("/:id/bid", bidCtrl.PlaceBid)
				items.GET("/:id/auto-bid", bidCtrl.GetAutoBid)
			}
		}

		my := api.Group("/my")
		my.Use(middleware.JWTAuth())
		{
			my.GET("/items", itemCtrl.GetMyItems)
			my.GET("/bids", bidCtrl.GetMyBids)
			my.GET("/auto-bids", bidCtrl.GetMyAutoBids)
			my.POST("/auto-bids", bidCtrl.SetAutoBid)
			my.DELETE("/auto-bids/:id", bidCtrl.CancelAutoBid)
			my.GET("/orders/buyer", orderCtrl.GetBuyerOrders)
			my.GET("/orders/seller", orderCtrl.GetSellerOrders)
			my.GET("/notifications", notifyCtrl.GetMyNotifications)
			my.GET("/notifications/unread-count", notifyCtrl.GetUnreadCount)
			my.POST("/notifications/mark-read", notifyCtrl.MarkAsRead)
			my.POST("/notifications/mark-all-read", notifyCtrl.MarkAllAsRead)
			my.GET("/statistics", statsCtrl.GetMyStatistics)
		}

		orders := api.Group("/orders")
		orders.Use(middleware.JWTAuth())
		{
			orders.POST("", orderCtrl.Create)
			orders.GET("/:id", orderCtrl.GetByID)
			orders.GET("/no/:order_no", orderCtrl.GetByNo)
			orders.POST("/:id/pay", orderCtrl.Pay)
			orders.POST("/:id/ship", orderCtrl.Ship)
			orders.POST("/:id/confirm-delivery", orderCtrl.ConfirmDelivery)
			orders.POST("/:id/complete", orderCtrl.Complete)
		}

		reviews := api.Group("/reviews")
		{
			reviews.GET("/order/:order_id", reviewCtrl.GetOrderReviews)
			reviews.GET("/user/:user_id", reviewCtrl.GetUserReviews)
			reviews.GET("/user/:user_id/rating", reviewCtrl.GetUserRating)
			reviews.GET("/:id", reviewCtrl.GetByID)

			reviews.Use(middleware.JWTAuth())
			{
				reviews.POST("", reviewCtrl.Create)
			}
		}

		sessions := api.Group("/sessions")
		{
			sessions.GET("", sessionCtrl.GetList)
			sessions.GET("/active", sessionCtrl.GetActive)
			sessions.GET("/:id", sessionCtrl.GetByID)

			sessions.Use(middleware.JWTAuth(), middleware.AdminRequired())
			{
				sessions.POST("", sessionCtrl.Create)
				sessions.PUT("/:id", sessionCtrl.Update)
				sessions.POST("/:id/items", sessionCtrl.AddItems)
				sessions.DELETE("/:id/items/:item_id", sessionCtrl.RemoveItem)
				sessions.POST("/:id/start", sessionCtrl.Start)
				sessions.POST("/:id/end", sessionCtrl.End)
				sessions.POST("/:id/cancel", sessionCtrl.Cancel)
			}
		}

		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.AdminRequired())
		{
			admin.GET("/users", userCtrl.GetUserList)
			admin.PUT("/users/:id/status", userCtrl.UpdateUserStatus)
			admin.GET("/orders", orderCtrl.GetAllOrders)
			admin.GET("/statistics/overall", statsCtrl.GetOverall)
			admin.GET("/export/orders", statsCtrl.ExportOrders)
			admin.GET("/export/bids", statsCtrl.ExportBids)
		}

		user := api.Group("/user")
		user.Use(middleware.JWTAuth())
		{
			user.GET("/info", userCtrl.GetUserInfo)
			user.PUT("/info", userCtrl.UpdateUser)
			user.POST("/change-password", userCtrl.ChangePassword)
		}
	}

	return r
}
