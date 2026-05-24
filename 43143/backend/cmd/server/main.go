package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"skillshare/internal/config"
	"skillshare/internal/database"
	"skillshare/internal/handlers"
	"skillshare/internal/middleware"
	"skillshare/internal/repository"
	"skillshare/internal/service"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.Close()

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	userRepo := repository.NewUserRepository(database.DB)
	skillRepo := repository.NewSkillRepository(database.DB)
	bookingRepo := repository.NewBookingRepository(database.DB)
	reviewRepo := repository.NewReviewRepository(database.DB)
	messageRepo := repository.NewMessageRepository(database.DB)
	paymentRepo := repository.NewPaymentRepository(database.DB)
	scheduleRepo := repository.NewScheduleRepository(database.DB)

	userService := service.NewUserService(userRepo, cfg.JWT.Secret)
	skillService := service.NewSkillService(skillRepo)
	bookingService := service.NewBookingService(bookingRepo, reviewRepo, skillRepo, userRepo)
	messageService := service.NewMessageService(messageRepo, cfg.Encryption.MessageKey)
	paymentService := service.NewPaymentService(paymentRepo, bookingRepo, cfg.Payment.EscrowDays)
	scheduleService := service.NewScheduleService(scheduleRepo)
	statsService := service.NewStatsService(bookingRepo, paymentRepo)

	userHandler := handlers.NewUserHandler(userService)
	skillHandler := handlers.NewSkillHandler(skillService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	messageHandler := handlers.NewMessageHandler(messageService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)
	statsHandler := handlers.NewStatsHandler(statsService)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
		}

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
			users.GET("/:id", userHandler.GetUser)
			users.GET("", userHandler.ListUsers)
			users.POST("/me/tags", userHandler.AddSkillTags)
			users.DELETE("/me/tags/:tag_id", userHandler.RemoveSkillTag)
		}

		categories := api.Group("/categories")
		categories.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			categories.POST("", skillHandler.CreateCategory)
			categories.GET("", skillHandler.GetCategories)
			categories.PUT("/:id", skillHandler.UpdateCategory)
			categories.DELETE("/:id", skillHandler.DeleteCategory)
		}

		tags := api.Group("/tags")
		tags.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			tags.POST("", skillHandler.CreateTag)
			tags.GET("", skillHandler.GetTags)
			tags.PUT("/:id", skillHandler.UpdateTag)
			tags.DELETE("/:id", skillHandler.DeleteTag)
		}

		skills := api.Group("/skills")
		skills.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			skills.POST("", skillHandler.CreateSkill)
			skills.GET("", skillHandler.GetSkills)
			skills.GET("/popular", skillHandler.GetPopularSkills)
			skills.GET("/:id", skillHandler.GetSkill)
			skills.PUT("/:id", skillHandler.UpdateSkill)
			skills.DELETE("/:id", skillHandler.DeleteSkill)
		}

		postings := api.Group("/postings")
		postings.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			postings.POST("", skillHandler.CreatePosting)
			postings.GET("", skillHandler.GetPostings)
			postings.GET("/match", skillHandler.MatchSkills)
			postings.GET("/:id", skillHandler.GetPosting)
			postings.PUT("/:id", skillHandler.UpdatePosting)
			postings.DELETE("/:id", skillHandler.DeletePosting)
		}

		bookings := api.Group("/bookings")
		bookings.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			bookings.POST("", bookingHandler.CreateBooking)
			bookings.GET("", bookingHandler.ListBookings)
			bookings.GET("/:id", bookingHandler.GetBooking)
			bookings.PUT("/:id/confirm", bookingHandler.ConfirmBooking)
			bookings.PUT("/:id/reject", bookingHandler.RejectBooking)
			bookings.PUT("/:id/cancel", bookingHandler.CancelBooking)
			bookings.PUT("/:id/complete", bookingHandler.CompleteBooking)
		}

		reviews := api.Group("/reviews")
		reviews.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			reviews.POST("", bookingHandler.CreateReview)
			reviews.GET("/posting/:posting_id", bookingHandler.GetReviewsByPosting)
			reviews.GET("/user/:user_id", bookingHandler.GetReviewsByUser)
		}

		complaints := api.Group("/complaints")
		complaints.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			complaints.POST("", bookingHandler.CreateComplaint)
			complaints.GET("", bookingHandler.GetComplaints)
			complaints.PUT("/:id/handle", bookingHandler.HandleComplaint)
		}

		messages := api.Group("/messages")
		messages.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			messages.POST("", messageHandler.SendMessage)
			messages.GET("/conversations", messageHandler.GetConversations)
			messages.GET("/unread", messageHandler.GetUnreadCount)
			messages.GET("/:user_id", messageHandler.GetMessages)
			messages.PUT("/:user_id/read", messageHandler.MarkAsRead)
		}

		payments := api.Group("/payments")
		payments.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			payments.POST("", paymentHandler.CreatePayment)
			payments.GET("", paymentHandler.GetUserPayments)
			payments.GET("/:id", paymentHandler.GetPayment)
			payments.GET("/wallet", paymentHandler.GetWallet)
			payments.POST("/withdraw", paymentHandler.Withdraw)
		}

		schedules := api.Group("/schedules")
		schedules.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			schedules.POST("", scheduleHandler.CreateSchedule)
			schedules.GET("", scheduleHandler.GetUserSchedules)
			schedules.GET("/availability/:user_id", scheduleHandler.GetUserAvailability)
			schedules.PUT("/:id", scheduleHandler.UpdateSchedule)
			schedules.DELETE("/:id", scheduleHandler.DeleteSchedule)
		}

		stats := api.Group("/stats")
		stats.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			stats.GET("/teacher", statsHandler.GetTeacherStats)
			stats.GET("/monthly", statsHandler.GetMonthlyReport)
			stats.GET("/export", statsHandler.ExportReport)
		}
	}

	log.Printf("服务器启动在端口 %s", cfg.Server.Port)
	r.Run(":" + cfg.Server.Port)
}
