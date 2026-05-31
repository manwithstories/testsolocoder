package routes

import (
	"beehive-platform/config"
	"beehive-platform/handlers"
	"beehive-platform/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	authHandler := handlers.NewAuthHandler()
	beehiveHandler := handlers.NewBeehiveHandler()
	healthHandler := handlers.NewHealthHandler()
	harvestHandler := handlers.NewHarvestHandler()
	inventoryHandler := handlers.NewInventoryHandler()
	productHandler := handlers.NewProductHandler()
	orderHandler := handlers.NewOrderHandler()
	inspectionHandler := handlers.NewInspectionHandler()
	postHandler := handlers.NewPostHandler()
	notificationHandler := handlers.NewNotificationHandler()
	analyticsHandler := handlers.NewAnalyticsHandler()
	uploadHandler := handlers.NewUploadHandler(cfg)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/auth/profile", authHandler.GetProfile)
			auth.PUT("/auth/profile", authHandler.UpdateProfile)

			auth.GET("/notifications", notificationHandler.List)
			auth.PUT("/notifications/:id/read", notificationHandler.Read)
			auth.PUT("/notifications/read-all", notificationHandler.ReadAll)
			auth.GET("/notifications/unread-count", notificationHandler.UnreadCount)

			auth.GET("/analytics/overview", analyticsHandler.GetOverview)
			auth.GET("/analytics/production", analyticsHandler.GetProductionStats)
			auth.GET("/analytics/disease", analyticsHandler.GetDiseaseStats)
			auth.GET("/analytics/sales", analyticsHandler.GetSalesStats)
			auth.GET("/analytics/export", analyticsHandler.ExportReport)

			auth.POST("/uploads", uploadHandler.Upload)
			auth.GET("/uploads", uploadHandler.List)
			auth.DELETE("/uploads/:id", uploadHandler.Delete)

			beekeeper := auth.Group("")
			beekeeper.Use(middleware.RoleMiddleware("beekeeper"))
			{
				beekeeper.POST("/beehives", beehiveHandler.Create)
				beekeeper.GET("/beehives", beehiveHandler.List)
				beekeeper.GET("/beehives/groups", beehiveHandler.Groups)
				beekeeper.GET("/beehives/:id", beehiveHandler.Get)
				beekeeper.PUT("/beehives/:id", beehiveHandler.Update)
				beekeeper.DELETE("/beehives/:id", beehiveHandler.Delete)

				beekeeper.POST("/health-records", healthHandler.CreateRecord)
				beekeeper.GET("/health-records", healthHandler.List)
				beekeeper.GET("/health-records/:id", healthHandler.Get)
				beekeeper.GET("/disease-warnings", healthHandler.GetDiseaseWarnings)
				beekeeper.GET("/seasonal-tips", healthHandler.GetSeasonalTips)

				beekeeper.POST("/harvests", harvestHandler.Create)
				beekeeper.GET("/harvests", harvestHandler.List)
				beekeeper.GET("/harvests/:id", harvestHandler.Get)
				beekeeper.DELETE("/harvests/:id", harvestHandler.Delete)

				beekeeper.GET("/inventory", inventoryHandler.List)
				beekeeper.GET("/inventory/alerts", inventoryHandler.GetAlerts)
				beekeeper.GET("/inventory/:id", inventoryHandler.Get)
				beekeeper.PUT("/inventory/:id", inventoryHandler.Update)
				beekeeper.DELETE("/inventory/:id", inventoryHandler.Delete)

				beekeeper.POST("/inspections", inspectionHandler.Create)
				beekeeper.GET("/inspections", inspectionHandler.List)
				beekeeper.GET("/inspections/:id", inspectionHandler.Get)
				beekeeper.PUT("/inspections/:id/cancel", inspectionHandler.Cancel)

				beekeeper.POST("/products", productHandler.Create)
				beekeeper.GET("/my-products", productHandler.ListMyProducts)
				beekeeper.PUT("/products/:id", productHandler.Update)
				beekeeper.DELETE("/products/:id", productHandler.Delete)

				beekeeper.GET("/seller-orders", orderHandler.List)
				beekeeper.PUT("/orders/:id/ship", orderHandler.Ship)
				beekeeper.PUT("/orders/:id/rate-buyer", orderHandler.RateBuyer)
			}

			buyer := auth.Group("")
			buyer.Use(middleware.RoleMiddleware("buyer"))
			{
				buyer.POST("/orders", orderHandler.Create)
				buyer.GET("/orders", orderHandler.List)
				buyer.GET("/orders/:id", orderHandler.Get)
				buyer.PUT("/orders/:id/pay", orderHandler.Pay)
				buyer.PUT("/orders/:id/deliver", orderHandler.Deliver)
				buyer.PUT("/orders/:id/complete", orderHandler.Complete)
				buyer.PUT("/orders/:id/cancel", orderHandler.Cancel)
				buyer.PUT("/orders/:id/rate", orderHandler.RateSeller)
			}

			inspector := auth.Group("")
			inspector.Use(middleware.RoleMiddleware("inspector"))
			{
				inspector.GET("/inspections", inspectionHandler.List)
				inspector.GET("/inspections/:id", inspectionHandler.Get)
				inspector.PUT("/inspections/:id/assign", inspectionHandler.AssignInspector)
				inspector.PUT("/inspections/:id/result", inspectionHandler.SubmitResult)
			}
		}

		api.GET("/products", productHandler.List)
		api.GET("/products/:id", productHandler.Get)

		api.GET("/posts/categories", postHandler.GetCategories)
		api.GET("/posts", postHandler.List)
		api.GET("/posts/:id", postHandler.Get)

		auth2 := api.Group("")
		auth2.Use(middleware.AuthMiddleware())
		{
			auth2.POST("/posts", postHandler.Create)
			auth2.PUT("/posts/:id", postHandler.Update)
			auth2.DELETE("/posts/:id", postHandler.Delete)
			auth2.POST("/posts/:id/like", postHandler.Like)
			auth2.POST("/comments", postHandler.CreateComment)
			auth2.GET("/posts/:id/comments", postHandler.ListComments)
		}
	}

	r.Static("/static", "./uploads")
}
