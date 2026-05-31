package main

import (
	"garden-planner/config"
	"garden-planner/database"
	"garden-planner/handlers"
	"garden-planner/middleware"
	"garden-planner/models"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	gin.SetMode(config.AppConfig.GinMode)

	database.Connect()

	database.Migrate(
		&models.User{},
		&models.Plot{},
		&models.Plant{},
		&models.PlantingRecord{},
		&models.GrowthLog{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
		&models.Follow{},
		&models.DiseaseDiagnosis{},
		&models.SeedExchange{},
		&models.ExchangeOffer{},
		&models.Product{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
		&models.CalendarEvent{},
		&models.OperationLog{},
	)

	if err := os.MkdirAll(config.AppConfig.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Static("/uploads", config.AppConfig.UploadDir)

	router.Use(middleware.LoggingMiddleware())

	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "garden-planner-api",
				"version": "1.0.0",
			})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		users := api.Group("/users")
		users.Use(middleware.JWTAuth())
		{
			users.GET("/profile", handlers.GetProfile)
			users.PUT("/profile", handlers.UpdateProfile)
			users.PUT("/password", handlers.ChangePassword)
			users.GET("/:id", handlers.GetUserByID)
		}

		plots := api.Group("/plots")
		plots.Use(middleware.JWTAuth())
		{
			plots.POST("", handlers.CreatePlot)
			plots.GET("", handlers.GetPlots)
			plots.GET("/:id", handlers.GetPlot)
			plots.PUT("/:id", handlers.UpdatePlot)
			plots.DELETE("/:id", handlers.DeletePlot)
		}

		plants := api.Group("/plants")
		{
			plants.GET("", handlers.GetPlants)
			plants.GET("/:id", handlers.GetPlant)
			plants.POST("", middleware.JWTAuth(), handlers.CreatePlant)
		}

		planting := api.Group("/planting-records")
		planting.Use(middleware.JWTAuth())
		{
			planting.POST("", handlers.CreatePlantingRecord)
			planting.GET("", handlers.GetPlantingRecords)
			planting.GET("/:id", handlers.GetPlantingRecord)
			planting.PUT("/:id", handlers.UpdatePlantingRecord)
			planting.DELETE("/:id", handlers.DeletePlantingRecord)
			planting.GET("/:id/report", handlers.ExportPlantingReport)
		}

		growth := api.Group("/growth-logs")
		growth.Use(middleware.JWTAuth())
		{
			growth.POST("", handlers.CreateGrowthLog)
			growth.GET("", handlers.GetGrowthLogs)
			growth.GET("/:id", handlers.GetGrowthLog)
			growth.PUT("/:id", handlers.UpdateGrowthLog)
			growth.DELETE("/:id", handlers.DeleteGrowthLog)
		}

		uploads := api.Group("/upload")
		uploads.Use(middleware.JWTAuth())
		{
			uploads.POST("", handlers.UploadFile)
		}

		calendar := api.Group("/calendar")
		calendar.Use(middleware.JWTAuth())
		{
			calendar.POST("/events", handlers.CreateCalendarEvent)
			calendar.GET("/events", handlers.GetCalendarEvents)
			calendar.GET("/events/:id", handlers.GetCalendarEvent)
			calendar.PUT("/events/:id", handlers.UpdateCalendarEvent)
			calendar.DELETE("/events/:id", handlers.DeleteCalendarEvent)
			calendar.GET("/recommendations", handlers.GetPlantingRecommendations)
		}

		disease := api.Group("/disease-diagnosis")
		disease.Use(middleware.JWTAuth())
		{
			disease.POST("", handlers.CreateDiseaseDiagnosis)
			disease.GET("", handlers.GetDiseaseDiagnoses)
			disease.GET("/:id", handlers.GetDiseaseDiagnosis)
			disease.PUT("/:id", handlers.UpdateDiseaseDiagnosis)
			disease.DELETE("/:id", handlers.DeleteDiseaseDiagnosis)
		}

		posts := api.Group("/posts")
		{
			posts.GET("", handlers.GetPosts)
			posts.GET("/:id", handlers.GetPost)
			posts.GET("/:id/comments", handlers.GetComments)
		}

		postsAuth := api.Group("/posts")
		postsAuth.Use(middleware.JWTAuth())
		{
			postsAuth.POST("", handlers.CreatePost)
			postsAuth.PUT("/:id", handlers.UpdatePost)
			postsAuth.DELETE("/:id", handlers.DeletePost)
			postsAuth.POST("/:id/like", handlers.LikePost)
			postsAuth.POST("/:id/comments", handlers.CreateComment)
			postsAuth.DELETE("/comments/:id", handlers.DeleteComment)
		}

		follows := api.Group("/follows")
		follows.Use(middleware.JWTAuth())
		{
			follows.POST("/:id", handlers.FollowUser)
			follows.GET("/:id/followers", handlers.GetFollowers)
			follows.GET("/:id/following", handlers.GetFollowing)
		}

		exchanges := api.Group("/seed-exchanges")
		{
			exchanges.GET("", handlers.GetSeedExchanges)
			exchanges.GET("/:id", handlers.GetSeedExchange)
		}

		exchangesAuth := api.Group("/seed-exchanges")
		exchangesAuth.Use(middleware.JWTAuth())
		{
			exchangesAuth.POST("", handlers.CreateSeedExchange)
			exchangesAuth.PUT("/:id", handlers.UpdateSeedExchange)
			exchangesAuth.DELETE("/:id", handlers.DeleteSeedExchange)
			exchangesAuth.POST("/:id/offers", handlers.CreateExchangeOffer)
			exchangesAuth.GET("/:id/offers", handlers.GetExchangeOffers)
			exchangesAuth.PUT("/offers/:id", handlers.UpdateExchangeOffer)
			exchangesAuth.DELETE("/offers/:id", handlers.DeleteExchangeOffer)
		}

		products := api.Group("/products")
		{
			products.GET("", handlers.GetProducts)
			products.GET("/:id", handlers.GetProduct)
		}

		productsAuth := api.Group("/products")
		productsAuth.Use(middleware.JWTAuth())
		{
			productsAuth.POST("", handlers.CreateProduct)
			productsAuth.PUT("/:id", handlers.UpdateProduct)
			productsAuth.DELETE("/:id", handlers.DeleteProduct)
		}

		cart := api.Group("/cart")
		cart.Use(middleware.JWTAuth())
		{
			cart.POST("", handlers.AddToCart)
			cart.GET("", handlers.GetCart)
			cart.PUT("/:id", handlers.UpdateCart)
			cart.DELETE("/:id", handlers.RemoveFromCart)
			cart.DELETE("", handlers.ClearCart)
		}

		orders := api.Group("/orders")
		orders.Use(middleware.JWTAuth())
		{
			orders.POST("", handlers.CreateOrder)
			orders.GET("", handlers.GetOrders)
			orders.GET("/:id", handlers.GetOrder)
			orders.PUT("/:id", handlers.UpdateOrder)
			orders.POST("/:id/cancel", handlers.CancelOrder)
		}
	}

	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
	if err := router.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
