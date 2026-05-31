package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"print3d-platform/internal/config"
	"print3d-platform/internal/database"
	"print3d-platform/internal/handlers"
	"print3d-platform/internal/middleware"
	"print3d-platform/internal/models"
	"print3d-platform/internal/repository"
	"print3d-platform/internal/service"
	"print3d-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = utils.InitLogger(cfg.Logging.Level, "logs",
		cfg.Logging.MaxSize, cfg.Logging.MaxBackups, cfg.Logging.MaxAge)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	err = database.Migrate()
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = database.SeedData()
	if err != nil {
		utils.LogWarn("Failed to seed data: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	modelRepo := repository.NewModelRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	printerRepo := repository.NewPrinterRepository(db)
	fileRepo := repository.NewFileRepository(db)

	userService := service.NewUserService(userRepo, cfg)
	modelService := service.NewModelService(modelRepo, userRepo, orderRepo, cfg)
	orderService := service.NewOrderService(orderRepo, modelRepo, userRepo, printerRepo, cfg)
	printerService := service.NewPrinterService(printerRepo, userRepo, orderRepo, modelRepo)
	fileService := service.NewFileService(fileRepo, cfg)

	userHandler := handlers.NewUserHandler(userService)
	modelHandler := handlers.NewModelHandler(modelService)
	orderHandler := handlers.NewOrderHandler(orderService)
	printerHandler := handlers.NewPrinterHandler(printerService)
	fileHandler := handlers.NewFileHandler(fileService)
	statsHandler := handlers.NewStatsHandler(orderService, modelService, printerService)

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	r.MaxMultipartMemory = cfg.Server.MaxUploadSize

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
		}

		api.GET("/models", modelHandler.ListModels)
		api.GET("/models/hot", modelHandler.GetHotModels)
		api.GET("/models/:id", modelHandler.GetModel)
		api.GET("/models/:id/versions", modelHandler.GetVersions)
		api.GET("/models/:id/reviews", printerHandler.GetModelReviews)
		api.GET("/designers/:designer_id/models", modelHandler.GetDesignerModels)
		api.GET("/designers", userHandler.ListDesigners)
		api.GET("/printers", userHandler.ListPrinters)
		api.GET("/materials", printerHandler.GetMaterials)
		api.GET("/materials/:id", printerHandler.GetMaterial)
		api.POST("/orders/estimate", orderHandler.EstimatePrice)
		api.GET("/models/validate", modelHandler.ValidateModelFile)

		authenticated := api.Group("")
		authenticated.Use(middleware.JWTAuthMiddleware(&cfg.JWT))
		{
			user := authenticated.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.PUT("/profile/designer", userHandler.UpdateDesignerProfile)
				user.PUT("/profile/printer", userHandler.UpdatePrinterProfile)
				user.GET("/stats", userHandler.GetUserStats)
				user.GET("/notifications", userHandler.GetNotifications)
				user.PUT("/notifications/:id/read", userHandler.MarkNotificationRead)
				user.GET("/transactions", userHandler.GetTransactions)
			}

			modelsGroup := authenticated.Group("/models")
			{
				modelsGroup.POST("", middleware.RoleMiddleware(models.RoleDesigner), modelHandler.CreateModel)
				modelsGroup.PUT("/:id", middleware.RoleMiddleware(models.RoleDesigner), modelHandler.UpdateModel)
				modelsGroup.DELETE("/:id", middleware.RoleMiddleware(models.RoleDesigner), modelHandler.DeleteModel)
				modelsGroup.POST("/:id/file", middleware.RoleMiddleware(models.RoleDesigner), modelHandler.UploadModelFile)
				modelsGroup.POST("/:id/thumbnail", middleware.RoleMiddleware(models.RoleDesigner), modelHandler.UploadThumbnail)
				modelsGroup.POST("/:id/purchase", middleware.RoleMiddleware(models.RoleCustomer, models.RoleAdmin), modelHandler.PurchaseModel)
				modelsGroup.GET("/:id/download", modelHandler.DownloadModel)
				modelsGroup.POST("/:id/favorite", modelHandler.AddFavorite)
				modelsGroup.DELETE("/:id/favorite", modelHandler.RemoveFavorite)
				modelsGroup.POST("/:id/version", middleware.RoleMiddleware(models.RoleDesigner), modelHandler.CreateNewVersion)
				modelsGroup.GET("/designer/my", middleware.RoleMiddleware(models.RoleDesigner), modelHandler.GetDesignerModels)
			}

			authenticated.GET("/favorites", modelHandler.GetUserFavorites)
			authenticated.GET("/purchases", modelHandler.GetUserPurchases)

			orders := authenticated.Group("/orders")
			{
				orders.POST("", middleware.RoleMiddleware(models.RoleCustomer, models.RoleAdmin), orderHandler.CreateOrder)
				orders.GET("/my/customer", middleware.RoleMiddleware(models.RoleCustomer, models.RoleAdmin), orderHandler.ListCustomerOrders)
				orders.GET("/my/printer", middleware.RoleMiddleware(models.RolePrinter, models.RoleAdmin), orderHandler.ListPrinterOrders)
				orders.GET("/pending", middleware.RoleMiddleware(models.RolePrinter, models.RoleAdmin), orderHandler.GetPendingOrders)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.GET("/no/:order_no", orderHandler.GetOrderByNo)
				orders.GET("/:id/history", orderHandler.GetOrderHistory)
				orders.PUT("/:id/assign", middleware.RoleMiddleware(models.RolePrinter, models.RoleAdmin), orderHandler.AssignPrinter)
				orders.PUT("/:id/print/start", middleware.RoleMiddleware(models.RolePrinter, models.RoleAdmin), orderHandler.StartPrinting)
				orders.PUT("/:id/print/complete", middleware.RoleMiddleware(models.RolePrinter, models.RoleAdmin), orderHandler.CompletePrinting)
				orders.PUT("/:id/quality/approve", middleware.RoleMiddleware(models.RolePrinter, models.RoleAdmin), orderHandler.ApproveQuality)
				orders.PUT("/:id/ship", middleware.RoleMiddleware(models.RolePrinter, models.RoleAdmin), orderHandler.ShipOrder)
				orders.PUT("/:id/deliver", orderHandler.DeliverOrder)
				orders.PUT("/:id/complete", orderHandler.CompleteOrder)
				orders.PUT("/:id/cancel", orderHandler.CancelOrder)
			}

			printer := authenticated.Group("/printer")
			printer.Use(middleware.RoleMiddleware(models.RolePrinter, models.RoleAdmin))
			{
				printer.GET("/devices", printerHandler.GetPrinterDevices)
				printer.POST("/devices", printerHandler.CreateDevice)
				printer.GET("/devices/:id", printerHandler.GetDevice)
				printer.PUT("/devices/:id", printerHandler.UpdateDevice)
				printer.DELETE("/devices/:id", printerHandler.DeleteDevice)
				printer.GET("/devices/idle", printerHandler.GetIdleDevices)

				printer.GET("/inventory", printerHandler.GetPrinterInventory)
				printer.POST("/inventory", printerHandler.CreateInventory)
				printer.PUT("/inventory/:id", printerHandler.UpdateInventoryQuantity)
				printer.DELETE("/inventory/:id", printerHandler.DeleteInventory)

				printer.GET("/schedules", printerHandler.GetPrinterSchedules)
				printer.POST("/schedules", printerHandler.CreateSchedule)
			}

			reviews := authenticated.Group("/reviews")
			{
				reviews.POST("", middleware.RoleMiddleware(models.RoleCustomer, models.RoleAdmin), printerHandler.CreateReview)
				reviews.GET("/:id", printerHandler.GetReview)
				reviews.GET("/printer/:printer_id", printerHandler.GetPrinterReviews)
			}

			files := authenticated.Group("/files")
			{
				files.POST("/initiate", fileHandler.InitiateUpload)
				files.POST("/:id/chunk", fileHandler.UploadChunk)
				files.POST("/:id/complete", fileHandler.CompleteUpload)
				files.GET("/:id", fileHandler.GetUpload)
				files.GET("/my", fileHandler.GetUserUploads)
				files.DELETE("/:id", fileHandler.DeleteUpload)
				files.GET("/:id/logs", fileHandler.GetAccessLogs)
			}

			stats := authenticated.Group("/stats")
			stats.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleDesigner, models.RolePrinter))
			{
				stats.GET("/platform", statsHandler.GetPlatformStats)
				stats.GET("/revenue", statsHandler.GetRevenueStats)
				stats.GET("/materials", statsHandler.GetMaterialStats)
				stats.GET("/export", statsHandler.ExportStats)
			}
		}
	}

	health := r.Group("/health")
	{
		health.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"time":   time.Now().Format(time.RFC3339),
			})
		})
	}

	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	utils.LogInfo("Server starting on %s", serverAddr)
	utils.LogInfo("Server mode: %s", cfg.Server.Mode)

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
