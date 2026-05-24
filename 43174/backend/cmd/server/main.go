package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"campus-trade-platform/config"
	"campus-trade-platform/internal/handlers"
	"campus-trade-platform/internal/middleware"
	"campus-trade-platform/internal/models"
	"campus-trade-platform/internal/repository"
	"campus-trade-platform/internal/services"
	"campus-trade-platform/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	gin.SetMode(cfg.ServerMode)

	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	os.MkdirAll(cfg.UploadPath, 0755)
	os.MkdirAll("./exports", 0755)

	userRepo := repository.NewUserRepository(db)
	textbookRepo := repository.NewTextbookRepository(db)
	noteRepo := repository.NewNoteRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	userService := services.NewUserService(userRepo, db)
	textbookService := services.NewTextbookService(textbookRepo, categoryRepo, userRepo, db)
	noteService := services.NewNoteService(noteRepo, userRepo, db)
	transactionService := services.NewTransactionService(transactionRepo, textbookRepo, userRepo, db)
	orderService := services.NewOrderService(orderRepo, textbookRepo, userRepo, db)
	messageService := services.NewMessageService(messageRepo, userRepo, db)
	reviewService := services.NewReviewService(reviewRepo, userRepo, textbookRepo, noteRepo, db)
	categoryService := services.NewCategoryService(categoryRepo, db)
	statisticsService := services.NewStatisticsService(db)
	notificationService := services.NewNotificationService(db)

	userHandler := handlers.NewUserHandler(userService, cfg.JWTSecret, cfg.JWTExpireHours, cfg.UploadPath, cfg.MaxUploadSize)
	textbookHandler := handlers.NewTextbookHandler(textbookService, cfg.UploadPath, cfg.MaxUploadSize)
	noteHandler := handlers.NewNoteHandler(noteService, cfg.UploadPath, cfg.MaxUploadSize)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	orderHandler := handlers.NewOrderHandler(orderService)
	messageHandler := handlers.NewMessageHandler(messageService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	statisticsHandler := handlers.NewStatisticsHandler(statisticsService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	healthHandler := handlers.NewHealthHandler()

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler())

	r.Static("/uploads", cfg.UploadPath)

	api := r.Group("/api/v1")
	{
		api.GET("/health", healthHandler.HealthCheck)

		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.PUT("/password", userHandler.ChangePassword)
			users.POST("/avatar", userHandler.UploadAvatar)
			users.GET("/top-rated", userHandler.GetTopRatedUsers)

			admin := users.Group("")
			admin.Use(middleware.RoleMiddleware(string(models.RoleAdmin)))
			{
				admin.GET("", userHandler.GetAllUsers)
				admin.GET("/:id", userHandler.GetUserByID)
				admin.PUT("/:id/status", userHandler.UpdateUserStatus)
				admin.DELETE("/:id", userHandler.DeleteUser)
			}
		}

		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)

			admin := categories.Group("")
			admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			admin.Use(middleware.RoleMiddleware(string(models.RoleAdmin)))
			{
				admin.POST("", categoryHandler.CreateCategory)
				admin.PUT("/:id", categoryHandler.UpdateCategory)
				admin.DELETE("/:id", categoryHandler.DeleteCategory)
			}
		}

		textbooks := api.Group("/textbooks")
		{
			textbooks.GET("", textbookHandler.GetAllTextbooks)
			textbooks.GET("/search/isbn", textbookHandler.SearchByISBN)
			textbooks.GET("/popular", textbookHandler.GetPopularTextbooks)
			textbooks.GET("/seller/:seller_id", textbookHandler.GetSellerTextbooks)
			textbooks.GET("/:id", textbookHandler.GetTextbookByID)

			auth := textbooks.Group("")
			auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				auth.POST("", textbookHandler.CreateTextbook)
				auth.POST("/cover-image", textbookHandler.UploadCoverImage)
				auth.PUT("/:id", textbookHandler.UpdateTextbook)
				auth.PUT("/:id/status", textbookHandler.UpdateTextbookStatus)
				auth.DELETE("/:id", textbookHandler.DeleteTextbook)
			}
		}

		notes := api.Group("/notes")
		{
			notes.GET("", noteHandler.GetAllNotes)
			notes.GET("/featured", noteHandler.GetFeaturedNotes)
			notes.GET("/uploader/:uploader_id", noteHandler.GetUploaderNotes)
			notes.GET("/:id", noteHandler.GetNoteByID)

			auth := notes.Group("")
			auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				auth.POST("", noteHandler.CreateNote)
				auth.POST("/upload", noteHandler.UploadNoteFile)
				auth.PUT("/:id", noteHandler.UpdateNote)
				auth.DELETE("/:id", noteHandler.DeleteNote)
				auth.POST("/:id/download", noteHandler.IncrementDownload)

				admin := auth.Group("")
				admin.Use(middleware.RoleMiddleware(string(models.RoleAdmin)))
				{
					admin.PUT("/:id/featured", noteHandler.SetFeatured)
				}
			}
		}

		transactions := api.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetAllTransactions)
			transactions.GET("/buyer/:buyer_id", transactionHandler.GetBuyerTransactions)
			transactions.GET("/seller/:seller_id", transactionHandler.GetSellerTransactions)
			transactions.GET("/:id", transactionHandler.GetTransactionByID)
			transactions.PUT("/:id/confirm", transactionHandler.ConfirmTransaction)
			transactions.PUT("/:id/complete", transactionHandler.CompleteTransaction)
			transactions.PUT("/:id/cancel", transactionHandler.CancelTransaction)
			transactions.PUT("/:id/negotiate", transactionHandler.StartNegotiation)
		}

		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetAllOrders)
			orders.GET("/my", orderHandler.GetUserOrders)
			orders.GET("/order-no/:order_no", orderHandler.GetOrderByOrderNo)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.PUT("/:id/pay", orderHandler.PayOrder)
			orders.PUT("/:id/ship", orderHandler.ShipOrder)
			orders.PUT("/:id/deliver", orderHandler.DeliverOrder)
			orders.PUT("/:id/complete", orderHandler.CompleteOrder)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
			orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		}

		messages := api.Group("/messages")
		messages.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			messages.POST("", messageHandler.CreateMessage)
			messages.GET("/conversation", messageHandler.GetConversation)
			messages.GET("/unread-count", messageHandler.GetUnreadCount)
			messages.PUT("/mark-read", messageHandler.MarkAsRead)
		}

		reviews := api.Group("/reviews")
		{
			reviews.GET("/textbook/:textbook_id", reviewHandler.GetTextbookReviews)
			reviews.GET("/note/:note_id", reviewHandler.GetNoteReviews)

			auth := reviews.Group("")
			auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				auth.POST("", reviewHandler.CreateReview)

				admin := auth.Group("")
				admin.Use(middleware.RoleMiddleware(string(models.RoleAdmin)))
				{
					admin.GET("", reviewHandler.GetAllReviews)
					admin.PUT("/:id/hide", reviewHandler.HideReview)
					admin.PUT("/:id/malicious", reviewHandler.MarkMalicious)
				}
			}
		}

		notifications := api.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			notifications.GET("", notificationHandler.GetUserNotifications)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
		}

		statistics := api.Group("/statistics")
		statistics.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		statistics.Use(middleware.RoleMiddleware(string(models.RoleAdmin)))
		{
			statistics.GET("/textbooks", statisticsHandler.GetTextbookStats)
			statistics.GET("/users", statisticsHandler.GetUserStats)
			statistics.GET("/orders", statisticsHandler.GetOrderStats)
			statistics.GET("/popular-textbooks", statisticsHandler.GetPopularTextbooks)
			statistics.GET("/top-users", statisticsHandler.GetTopUsers)
			statistics.GET("/monthly", statisticsHandler.GetMonthlyStats)
			statistics.GET("/export", statisticsHandler.ExportMonthlyReport)
		}
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
