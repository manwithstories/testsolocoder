package routes

import (
	"secondhand-trading/controllers"
	"secondhand-trading/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		user := api.Group("/users")
		user.Use(middleware.AuthRequired())
		{
			user.GET("/profile", controllers.GetProfile)
			user.PUT("/profile", controllers.UpdateProfile)
			user.GET("/credit-score", controllers.GetUserCreditScore)
			user.GET("/:id", controllers.GetUserByID)
			user.GET("/:id/reviews", controllers.GetUserReviews)
		}

		products := api.Group("/products")
		{
			products.GET("", controllers.SearchProducts)
			products.GET("/categories", controllers.GetCategories)
			products.GET("/:id", controllers.GetProduct)

			products.POST("", middleware.AuthRequired(), controllers.CreateProduct)
			products.POST("/:id/images", middleware.AuthRequired(), controllers.UploadProductImages)
			products.PUT("/:id", middleware.AuthRequired(), controllers.UpdateProduct)
			products.DELETE("/:id", middleware.AuthRequired(), controllers.DeleteProduct)
			products.POST("/:id/relist", middleware.AuthRequired(), controllers.RelistProduct)
			products.GET("/mine/list", middleware.AuthRequired(), controllers.GetMyProducts)
		}

		transactions := api.Group("/transactions")
		transactions.Use(middleware.AuthRequired())
		{
			transactions.POST("", controllers.CreateTransaction)
			transactions.GET("", controllers.GetMyTransactions)
			transactions.GET("/:id", controllers.GetTransaction)
			transactions.POST("/:id/confirm", controllers.ConfirmTransaction)
			transactions.POST("/:id/cancel", controllers.CancelTransaction)
		}

		favorites := api.Group("/favorites")
		favorites.Use(middleware.AuthRequired())
		{
			favorites.POST("", controllers.AddFavorite)
			favorites.DELETE("/:product_id", controllers.RemoveFavorite)
			favorites.GET("", controllers.GetFavorites)
			favorites.GET("/check/:product_id", controllers.CheckFavorite)
		}

		reviews := api.Group("/reviews")
		reviews.Use(middleware.AuthRequired())
		{
			reviews.POST("", controllers.CreateReview)
			reviews.GET("/transaction/:transaction_id", controllers.GetTransactionReviews)
			reviews.GET("/mine/received", controllers.GetMyReceivedReviews)
		}

		notifications := api.Group("/notifications")
		notifications.Use(middleware.AuthRequired())
		{
			notifications.GET("", controllers.GetNotifications)
			notifications.PUT("/:id/read", controllers.MarkNotificationRead)
			notifications.PUT("/read-all", controllers.MarkAllNotificationsRead)
			notifications.DELETE("/:id", controllers.DeleteNotification)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
		{
			admin.GET("/products/pending-review", controllers.GetPendingReviewProducts)
			admin.POST("/products/:id/review", controllers.ReviewProduct)
			admin.POST("/products/batch-review", controllers.BatchReviewProducts)
		}
	}
}
