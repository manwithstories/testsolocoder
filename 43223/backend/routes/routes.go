package routes

import (
	"coffee-platform/config"
	"coffee-platform/handlers"
	"coffee-platform/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	authHandler := handlers.NewAuthHandler(cfg)
	userHandler := handlers.NewUserHandler(cfg)
	productHandler := handlers.NewProductHandler(cfg)
	roastingHandler := handlers.NewRoastingHandler(cfg)
	orderHandler := handlers.NewOrderHandler(cfg)
	cuppingHandler := handlers.NewCuppingHandler(cfg)
	certHandler := handlers.NewCertificationHandler(cfg)
	statsHandler := handlers.NewStatsHandler(cfg)
	searchHandler := handlers.NewSearchHandler(cfg)

	r.Static("/uploads", cfg.App.UploadDir)

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		api.GET("/products", productHandler.List)
		api.GET("/products/:id", productHandler.Get)
		api.GET("/origins", productHandler.GetOrigins)
		api.GET("/roast-levels", productHandler.GetRoastLevels)
		api.GET("/process-methods", productHandler.GetProcessMethods)

		api.GET("/search", searchHandler.SearchProducts)
		api.GET("/search/suggest", searchHandler.Suggest)
		api.POST("/search/advanced", searchHandler.AdvancedSearch)

		api.GET("/roasters", certHandler.ListCertifiedRoasters)
		api.GET("/roasters/:id", certHandler.GetRoasterProfile)

		api.GET("/cupping/scores", cuppingHandler.List)
		api.GET("/cupping/scores/:id", cuppingHandler.Get)
		api.GET("/cupping/stats", cuppingHandler.GetProductStats)
		api.GET("/cupping/trend", cuppingHandler.GetScoreTrend)

		api.GET("/roasting/records", roastingHandler.List)
		api.GET("/roasting/records/:id", roastingHandler.Get)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(&cfg.JWT))
		{
			protected.GET("/auth/profile", authHandler.Profile)
			protected.PUT("/auth/profile", authHandler.UpdateProfile)

			protected.POST("/products", middleware.RoasterMiddleware(), productHandler.Create)
			protected.PUT("/products/:id", middleware.RoasterMiddleware(), productHandler.Update)
			protected.PATCH("/products/:id/status", middleware.RoasterMiddleware(), productHandler.UpdateStatus)
			protected.DELETE("/products/:id", middleware.RoasterMiddleware(), productHandler.Delete)
			protected.POST("/products/:id/images", middleware.RoasterMiddleware(), productHandler.UploadImage)
			protected.DELETE("/products/:id/images/:image_id", middleware.RoasterMiddleware(), productHandler.DeleteImage)
			protected.POST("/products/import", middleware.AdminMiddleware(), productHandler.ImportCSV)

			protected.POST("/roasting/records", middleware.RoasterMiddleware(), roastingHandler.Create)
			protected.PUT("/roasting/records/:id", middleware.RoasterMiddleware(), roastingHandler.Update)
			protected.DELETE("/roasting/records/:id", middleware.RoasterMiddleware(), roastingHandler.Delete)
			protected.POST("/roasting/compare", roastingHandler.Compare)
			protected.GET("/roasting/stats", roastingHandler.GetStats)

			protected.GET("/cart", orderHandler.GetCart)
			protected.POST("/cart", orderHandler.AddToCart)
			protected.PUT("/cart/:id", orderHandler.UpdateCartItem)
			protected.DELETE("/cart/:id", orderHandler.RemoveFromCart)
			protected.DELETE("/cart", orderHandler.ClearCart)

			protected.GET("/orders", orderHandler.List)
			protected.GET("/orders/:id", orderHandler.Get)
			protected.POST("/orders", orderHandler.Create)
			protected.PATCH("/orders/:id/status", middleware.AdminMiddleware(), orderHandler.UpdateStatus)
			protected.POST("/orders/:id/cancel", orderHandler.CancelOrder)
			protected.POST("/orders/pay", orderHandler.Pay)

			protected.POST("/cupping/scores", cuppingHandler.Create)
			protected.PUT("/cupping/scores/:id", cuppingHandler.Update)
			protected.DELETE("/cupping/scores/:id", cuppingHandler.Delete)
			protected.GET("/cupping/my-scores", cuppingHandler.GetUserScores)

			protected.POST("/certification/apply", certHandler.Apply)
			protected.PUT("/certification/apply", certHandler.UpdateApplication)
			protected.GET("/certification/my", certHandler.GetMyCertification)

			protected.GET("/search/history", searchHandler.GetSearchHistory)

			admin := protected.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", userHandler.List)
				admin.GET("/users/:id", userHandler.Get)
				admin.PATCH("/users/:id/status", userHandler.UpdateStatus)
				admin.PATCH("/users/:id/role", userHandler.UpdateRole)
				admin.DELETE("/users/:id", userHandler.Delete)
				admin.GET("/users/activity", userHandler.GetUserActivity)

				admin.GET("/certification", certHandler.List)
				admin.POST("/certification/:id/review", certHandler.Review)

				admin.GET("/stats/sales", statsHandler.GetSalesStats)
				admin.GET("/stats/sales-trend", statsHandler.GetSalesTrend)
				admin.GET("/stats/origins", statsHandler.GetOriginDistribution)
				admin.GET("/stats/user-activity", statsHandler.GetUserActivityStats)
				admin.GET("/stats/top-products", statsHandler.GetTopProducts)
				admin.GET("/stats/export/excel", statsHandler.ExportExcel)
				admin.GET("/stats/export/pdf", statsHandler.ExportPDF)
			}
		}
	}
}
