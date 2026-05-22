package routes

import (
	"multishop/internal/config"
	"multishop/internal/handlers"
	"multishop/internal/middleware"
	"multishop/internal/models"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, cfg *config.Config) {
	r.Use(middleware.CORSMiddleware())
	r.Static("/uploads", "./uploads")

	authHandler := handlers.NewAuthHandler(cfg)
	shopHandler := handlers.NewShopHandler()
	productHandler := handlers.NewProductHandler()
	orderHandler := handlers.NewOrderHandler()
	reviewHandler := handlers.NewReviewHandler()
	favoriteHandler := handlers.NewFavoriteHandler()
	notificationHandler := handlers.NewNotificationHandler()
	categoryHandler := handlers.NewCategoryHandler()
	adminHandler := handlers.NewAdminHandler()
	statisticsHandler := handlers.NewStatisticsHandler()
	uploadHandler := handlers.NewUploadHandler(cfg)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(cfg), authHandler.GetProfile)
			auth.PUT("/profile", middleware.AuthMiddleware(cfg), authHandler.UpdateProfile)
			auth.PUT("/password", middleware.AuthMiddleware(cfg), authHandler.ChangePassword)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.List)
			categories.POST("", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleAdmin), categoryHandler.Create)
			categories.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleAdmin), categoryHandler.Update)
			categories.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleAdmin), categoryHandler.Delete)
			categories.GET("/all", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleAdmin), categoryHandler.GetAll)
		}

		shops := api.Group("/shops")
		{
			shops.GET("", shopHandler.List)
			shops.GET("/:id", shopHandler.GetByID)
			shops.POST("/apply", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller), shopHandler.Apply)
			shops.GET("/my", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller), shopHandler.GetMyShop)
			shops.PUT("/my", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller), shopHandler.Update)
			shops.PUT("/:id/review", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleAdmin), shopHandler.Review)
		}

		products := api.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.GetByID)
			products.POST("", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller), productHandler.Create)
			products.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller), productHandler.Update)
			products.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller), productHandler.Delete)
			products.GET("/my/list", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller), productHandler.MyProducts)
		}

		reviews := api.Group("/reviews")
		{
			reviews.GET("/product/:product_id", reviewHandler.GetProductReviews)
			reviews.POST("", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleBuyer), reviewHandler.Create)
			reviews.PUT("/:id/reply", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller), reviewHandler.Reply)
			reviews.GET("/my/shop", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller), reviewHandler.GetShopReviews)
		}

		cart := api.Group("/cart")
		{
			cart.Use(middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleBuyer))
			cart.GET("", orderHandler.GetCart)
			cart.POST("", orderHandler.AddToCart)
			cart.PUT("/:id", orderHandler.UpdateCart)
			cart.DELETE("/:id", orderHandler.RemoveCartItem)
		}

		orders := api.Group("/orders")
		{
			orders.Use(middleware.AuthMiddleware(cfg))
			orders.GET("", orderHandler.GetOrders)
			orders.GET("/:id", orderHandler.GetOrderDetail)
			orders.POST("", middleware.RoleMiddleware(models.RoleBuyer), orderHandler.CreateOrder)
			orders.POST("/:id/pay", middleware.RoleMiddleware(models.RoleBuyer), orderHandler.PayOrder)
			orders.POST("/:id/ship", middleware.RoleMiddleware(models.RoleSeller), orderHandler.ShipOrder)
			orders.POST("/:id/confirm", middleware.RoleMiddleware(models.RoleBuyer), orderHandler.ConfirmReceive)
			orders.POST("/:id/cancel", middleware.RoleMiddleware(models.RoleBuyer), orderHandler.CancelOrder)
			orders.POST("/refund", middleware.RoleMiddleware(models.RoleBuyer), orderHandler.ApplyRefund)
			orders.PUT("/refund/:id/review", middleware.RoleMiddleware(models.RoleSeller), orderHandler.ReviewRefund)
		}

		favorites := api.Group("/favorites")
		{
			favorites.Use(middleware.AuthMiddleware(cfg))
			favorites.POST("/shop/:shop_id", favoriteHandler.ToggleShop)
			favorites.POST("/product/:product_id", favoriteHandler.ToggleProduct)
			favorites.GET("/shops", favoriteHandler.GetMyShops)
			favorites.GET("/products", favoriteHandler.GetMyProducts)
		}

		notifications := api.Group("/notifications")
		{
			notifications.Use(middleware.AuthMiddleware(cfg))
			notifications.GET("", notificationHandler.List)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
		}

		admin := api.Group("/admin")
		{
			admin.Use(middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleAdmin))
			admin.GET("/users", adminHandler.GetUsers)
			admin.GET("/disputes", adminHandler.GetDisputes)
			admin.PUT("/disputes/:id/resolve", adminHandler.ResolveDispute)
			admin.GET("/statistics", adminHandler.GetStatistics)
			admin.GET("/export/orders", adminHandler.ExportOrders)
		}

		disputes := api.Group("/disputes")
		{
			disputes.Use(middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleBuyer))
			disputes.POST("", adminHandler.CreateDispute)
		}

		statistics := api.Group("/statistics")
		{
			statistics.Use(middleware.AuthMiddleware(cfg), middleware.RoleMiddleware(models.RoleSeller))
			statistics.GET("/shop", statisticsHandler.GetShopStatistics)
			statistics.GET("/shop/export", statisticsHandler.ExportShopOrders)
		}

		upload := api.Group("/upload")
		{
			upload.Use(middleware.AuthMiddleware(cfg))
			upload.POST("/image", uploadHandler.UploadImage)
			upload.POST("/multiple", uploadHandler.UploadMultiple)
		}
	}
}
