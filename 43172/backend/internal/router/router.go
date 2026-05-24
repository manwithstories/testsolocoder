package router

import (
	"net/http"
	"os"

	"luxury-trading-platform/internal/cache"
	"luxury-trading-platform/internal/config"
	"luxury-trading-platform/internal/handler"
	"luxury-trading-platform/internal/middleware"
	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/repository"
	"luxury-trading-platform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, rdb *cache.RedisClient, cfg *config.Config, log *logrus.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware(log))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.AuditLogMiddleware(log))

	r.Static("/uploads", "./uploads")

	ensureUploadDirs(cfg.Upload.Path)

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	authRepo := repository.NewAuthenticationRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	statRepo := repository.NewStatisticRepository(db)

	userService := service.NewUserService(userRepo, rdb, db)
	productService := service.NewProductService(productRepo, rdb, db, log, cfg.Upload.Path)
	orderService := service.NewOrderService(orderRepo, productRepo, userRepo, rdb, db)
	authService := service.NewAuthenticationService(authRepo, orderRepo, productRepo, userRepo, rdb, db)
	reviewService := service.NewReviewService(reviewRepo, orderRepo, userRepo, rdb, db)
	statService := service.NewStatisticService(statRepo, orderRepo, productRepo, userRepo, rdb, db)
	pdfService := service.NewPDFService(cfg.Upload.Path)

	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	authHandler := handler.NewAuthenticationHandler(authService, pdfService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	statHandler := handler.NewStatisticHandler(statService)

	r.GET("/health", userHandler.HealthCheck)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/register/authenticator", userHandler.RegisterAuthenticator)
			auth.POST("/login", userHandler.Login)
		}

		products := api.Group("/products")
		{
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.GET("/:id/images", productHandler.GetProductImages)
			products.GET("/brands/list", productHandler.ListBrands)
		}

		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			users := authenticated.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.GET("/:id", userHandler.GetUser)
			}

			products := authenticated.Group("/products")
			{
				products.POST("", middleware.RoleMiddleware(model.RoleSeller), productHandler.CreateProduct)
				products.PUT("/:id", middleware.RoleMiddleware(model.RoleSeller), productHandler.UpdateProduct)
				products.DELETE("/:id", middleware.RoleMiddleware(model.RoleSeller), productHandler.DeleteProduct)
				products.PATCH("/:id/status", middleware.RoleMiddleware(model.RoleSeller), productHandler.UpdateProductStatus)
				products.POST("/:id/images", middleware.RoleMiddleware(model.RoleSeller), productHandler.UploadImages)
				products.DELETE("/:id/images", middleware.RoleMiddleware(model.RoleSeller), productHandler.DeleteProductImages)
				products.GET("/seller/my", middleware.RoleMiddleware(model.RoleSeller), productHandler.ListSellerProducts)
				products.POST("/brands", middleware.RoleMiddleware(model.RoleAdmin), productHandler.CreateBrand)
				products.GET("/brands/:id", productHandler.GetBrand)
			}

			orders := authenticated.Group("/orders")
			{
				orders.POST("", middleware.RoleMiddleware(model.RoleBuyer), orderHandler.CreateOrder)
				orders.GET("", orderHandler.ListOrders)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.GET("/number/:order_number", orderHandler.GetOrderByNumber)
				orders.POST("/:id/pay", middleware.RoleMiddleware(model.RoleBuyer), orderHandler.PayOrder)
				orders.POST("/:id/ship", middleware.RoleMiddleware(model.RoleSeller), orderHandler.ShipOrder)
				orders.POST("/:id/confirm", middleware.RoleMiddleware(model.RoleBuyer), orderHandler.ConfirmDelivery)
				orders.POST("/:id/cancel", orderHandler.CancelOrder)
			}

			authentications := authenticated.Group("/authentications")
			{
				authentications.POST("", middleware.RoleMiddleware(model.RoleBuyer), authHandler.CreateAuthentication)
				authentications.GET("", authHandler.ListAuthentications)
				authentications.GET("/:id", authHandler.GetAuthentication)
				authentications.GET("/order/:order_id", authHandler.GetAuthenticationByOrder)
				authentications.POST("/:id/accept", middleware.RoleMiddleware(model.RoleAuthenticator), authHandler.AcceptAuthentication)
				authentications.POST("/:id/complete", middleware.RoleMiddleware(model.RoleAuthenticator), authHandler.CompleteAuthentication)
				authentications.POST("/:id/reject", middleware.RoleMiddleware(model.RoleAuthenticator), authHandler.RejectAuthentication)
				authentications.POST("/:id/cancel", middleware.RoleMiddleware(model.RoleBuyer), authHandler.CancelAuthentication)
				authentications.GET("/:id/report/download", authHandler.DownloadReport)
			}

			reviews := authenticated.Group("/reviews")
			{
				reviews.POST("", reviewHandler.CreateReview)
				reviews.GET("", reviewHandler.ListReviews)
				reviews.GET("/:id", reviewHandler.GetReview)
				reviews.GET("/user/:id/rating", reviewHandler.GetUserAverageRating)
			}

			admin := authenticated.Group("/admin")
			admin.Use(middleware.RoleMiddleware(model.RoleAdmin))
			{
				admin.GET("/users", userHandler.ListUsers)
				admin.GET("/authenticators", userHandler.ListAuthenticators)
				admin.POST("/authenticators/:id/approve", userHandler.ApproveAuthenticator)
				admin.POST("/authenticators/:id/reject", userHandler.RejectAuthenticator)
			}

			stats := authenticated.Group("/statistics")
			{
				stats.GET("/dashboard", statHandler.GetDashboardStats)
				stats.POST("/invalidate-cache", statHandler.InvalidateCache)
			}
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "route not found",
		})
	})

	return r
}

func ensureUploadDirs(uploadPath string) {
	dirs := []string{
		uploadPath,
		uploadPath + "/products",
		uploadPath + "/reports",
		uploadPath + "/avatars",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logrus.Errorf("Failed to create directory %s: %v", dir, err)
		}
	}
}
