package main

import (
	"log"
	"photo-rental/internal/config"
	"photo-rental/internal/handlers"
	"photo-rental/internal/middleware"
	"photo-rental/internal/utils"
	"photo-rental/pkg/database"
	redispkg "photo-rental/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	utils.InitLogger()
	defer utils.Logger.Close()

	gin.SetMode(cfg.Server.Mode)

	database.InitDB(cfg.Database)
	redispkg.InitRedis(cfg.Redis)

	middleware.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpireHour)

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.Static("/uploads", "./uploads")

	userHandler := handlers.NewUserHandler()
	equipmentHandler := handlers.NewEquipmentHandler()
	orderHandler := handlers.NewOrderHandler()
	searchHandler := handlers.NewSearchHandler()
	reviewHandler := handlers.NewReviewHandler()
	settlementHandler := handlers.NewSettlementHandler()
	exportHandler := handlers.NewExportHandler()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		users := api.Group("/users")
		users.Use(middleware.JWTAuth())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.GET("/", middleware.RequireRole("admin"), userHandler.GetAllUsers)
			users.PUT("/verify", middleware.RequireRole("admin"), userHandler.VerifyUser)
		}

		equipments := api.Group("/equipments")
		{
			equipments.GET("/:id", equipmentHandler.GetEquipment)
			equipments.GET("/", searchHandler.SearchEquipments)
			equipments.GET("/categories", searchHandler.GetCategories)
			equipments.GET("/brands", searchHandler.GetBrands)
			equipments.GET("/:id/reserved-dates", searchHandler.GetEquipmentReservedDates)

			equipments.Use(middleware.JWTAuth())
			{
				equipments.POST("/", middleware.RequireRole("owner", "admin"), equipmentHandler.CreateEquipment)
				equipments.GET("/my", equipmentHandler.GetMyEquipments)
				equipments.PUT("/:id", equipmentHandler.UpdateEquipment)
				equipments.DELETE("/:id", equipmentHandler.DeleteEquipment)
				equipments.POST("/:id/images", equipmentHandler.UploadImage)
				equipments.DELETE("/:id/images/:imageId", equipmentHandler.DeleteImage)
			}
		}

		orders := api.Group("/orders")
		orders.Use(middleware.JWTAuth())
		{
			orders.POST("/", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.GET("/", orderHandler.GetMyOrders)
			orders.PUT("/:id/confirm", orderHandler.ConfirmOrder)
			orders.PUT("/:id/reject", orderHandler.RejectOrder)
			orders.PUT("/:id/start", orderHandler.StartRental)
			orders.PUT("/:id/complete", orderHandler.CompleteOrder)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
			orders.GET("/:id/status-history", orderHandler.GetOrderStatusHistory)
		}

		search := api.Group("/search")
		{
			search.POST("/equipments", searchHandler.SearchEquipments)
			search.POST("/orders", middleware.JWTAuth(), searchHandler.SearchOrders)
		}

		reviews := api.Group("/reviews")
		{
			reviews.GET("/equipment/:id", reviewHandler.GetEquipmentReviews)
			reviews.GET("/user/:id", reviewHandler.GetUserReviews)

			reviews.Use(middleware.JWTAuth())
			{
				reviews.POST("/", reviewHandler.CreateReview)
				reviews.GET("/my", reviewHandler.GetMyReviews)
			}
		}

		settlements := api.Group("/settlements")
		settlements.Use(middleware.JWTAuth())
		{
			settlements.POST("/", settlementHandler.CreateSettlement)
			settlements.GET("/:id", settlementHandler.GetSettlement)
			settlements.GET("/order/:orderId", settlementHandler.GetOrderSettlement)
			settlements.GET("/", settlementHandler.GetMySettlements)
		}

		export := api.Group("/export")
		export.Use(middleware.JWTAuth())
		{
			export.GET("/rentals", exportHandler.ExportRentalRecords)
			export.GET("/revenue", exportHandler.ExportRevenueReport)
		}
	}

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
