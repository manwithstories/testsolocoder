package main

import (
	"log"

	"beauty-salon-system/config"
	"beauty-salon-system/internal/handler"
	"beauty-salon-system/internal/middleware"
	"beauty-salon-system/internal/repository"
	redisclient "beauty-salon-system/internal/repository/redis"
	"beauty-salon-system/internal/repository/mysql"
	"beauty-salon-system/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := mysql.Init(cfg.MySQL.DSN(), cfg.MySQL.MaxOpenConns, cfg.MySQL.MaxIdleConns); err != nil {
		log.Fatalf("Failed to init mysql: %v", err)
	}
	defer mysql.Close()

	if err := redisclient.Init(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.Database, cfg.Redis.PoolSize); err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}
	defer redisclient.Close()

	db := mysql.DB

	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	technicianRepo := repository.NewTechnicianRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	customerPkgRepo := repository.NewCustomerPackageRepository(db)
	memberCardRepo := repository.NewMemberCardRepository(db)
	productRepo := repository.NewProductRepository(db)
	productRecordRepo := repository.NewProductRecordRepository(db)
	productSaleRepo := repository.NewProductSaleRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.Expire)
	customerService := service.NewCustomerService(customerRepo, userRepo)
	technicianService := service.NewTechnicianService(technicianRepo, userRepo, appointmentRepo, nil)
	serviceItemService := service.NewServiceItemService(serviceRepo)
	appointmentService := service.NewAppointmentService(
		appointmentRepo, technicianRepo, customerRepo, serviceRepo,
		customerPkgRepo, productRepo, productRecordRepo,
		cfg.Appointment.CancelFreeHours, cfg.Appointment.DeductPoints,
	)
	paymentService := service.NewPaymentService(
		paymentRepo, appointmentRepo, customerRepo,
		customerPkgRepo, memberCardRepo, serviceItemService,
	)
	productService := service.NewProductService(productRepo, productRecordRepo, productSaleRepo)
	notificationService := service.NewNotificationService(notificationRepo, userRepo)
	reportService := service.NewReportService(paymentRepo, appointmentRepo, serviceRepo, technicianRepo)
	auditService := service.NewAuditService(auditRepo)
	reviewService := service.NewReviewService(reviewRepo, technicianRepo)

	technicianService.SetNotificationService(notificationService)

	authHandler := handler.NewAuthHandler(authService)
	customerHandler := handler.NewCustomerHandler(customerService)
	technicianHandler := handler.NewTechnicianHandler(technicianService)
	serviceItemHandler := handler.NewServiceItemHandler(serviceItemService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService, notificationService, auditService)
	paymentHandler := handler.NewPaymentHandler(paymentService, notificationService, auditService, customerService)
	productHandler := handler.NewProductHandler(productService, notificationService, auditService)
	reportHandler := handler.NewReportHandler(reportService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	auditHandler := handler.NewAuditHandler(auditService)
	reviewHandler := handler.NewReviewHandler(reviewService)

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/profile", middleware.JWTAuth(cfg.JWT.Secret), authHandler.Profile)
		}

		customers := api.Group("/customers")
		customers.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			customers.POST("", customerHandler.Create)
			customers.GET("", customerHandler.List)
			customers.GET("/my", customerHandler.GetMyProfile)
			customers.PUT("/my", customerHandler.UpdateMyProfile)
			customers.GET("/:id", customerHandler.GetByID)
			customers.PUT("/:id", customerHandler.Update)
		}

		technicians := api.Group("/technicians")
		technicians.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			technicians.POST("", technicianHandler.Create)
			technicians.GET("", technicianHandler.List)
			technicians.GET("/all", technicianHandler.ListAll)
			technicians.GET("/my", technicianHandler.GetMyProfile)
			technicians.GET("/:id", technicianHandler.GetByID)
			technicians.PUT("/:id", technicianHandler.Update)
			technicians.POST("/:id/leave", technicianHandler.AddLeave)
			technicians.GET("/:id/leaves", technicianHandler.GetLeaves)
			technicians.GET("/:id/schedule", appointmentHandler.GetByTechnicianAndDate)
		}

		services := api.Group("/services")
		services.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			services.POST("", serviceItemHandler.Create)
			services.GET("", serviceItemHandler.List)
			services.GET("/all", serviceItemHandler.ListAll)
			services.GET("/:id", serviceItemHandler.GetByID)
			services.PUT("/:id", serviceItemHandler.Update)
			services.DELETE("/:id", serviceItemHandler.Delete)
			services.POST("/package", serviceItemHandler.AddPackageService)
			services.GET("/:id/package-services", serviceItemHandler.GetPackageServices)
			services.DELETE("/:id/package-services", serviceItemHandler.DeletePackageServices)
		}

		appointments := api.Group("/appointments")
		appointments.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			appointments.POST("", appointmentHandler.Create)
			appointments.GET("", appointmentHandler.List)
			appointments.GET("/my", appointmentHandler.GetMyAppointments)
			appointments.GET("/:id", appointmentHandler.GetByID)
			appointments.POST("/cancel", appointmentHandler.Cancel)
			appointments.POST("/reschedule", appointmentHandler.Reschedule)
			appointments.POST("/:id/complete", appointmentHandler.Complete)
		}

		available := api.Group("/available")
		available.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			available.GET("/:id/slots", appointmentHandler.GetAvailableSlots)
		}

		payments := api.Group("/payments")
		payments.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			payments.POST("", paymentHandler.CreatePayment)
			payments.GET("", paymentHandler.List)
			payments.GET("/:id", paymentHandler.GetByID)
			payments.POST("/member-card", paymentHandler.CreateMemberCard)
			payments.GET("/member-card/:id", paymentHandler.GetMemberCards)
			payments.POST("/member-card/:id/recharge", paymentHandler.RechargeCard)
			payments.POST("/package", paymentHandler.PurchasePackage)
			payments.GET("/package/:id", paymentHandler.GetCustomerPackages)
		}

		products := api.Group("/products")
		products.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			products.POST("", productHandler.Create)
			products.GET("", productHandler.List)
			products.GET("/all", productHandler.ListAll)
			products.GET("/low-stock", productHandler.GetLowStock)
			products.GET("/:id", productHandler.GetByID)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
			products.POST("/add-stock", productHandler.AddStock)
			products.POST("/deduct-stock", productHandler.DeductStock)
			products.GET("/records/list", productHandler.GetRecords)
			products.POST("/sale", productHandler.Sale)
			products.GET("/sales/list", productHandler.GetSales)
			products.POST("/stock-take", productHandler.StockTake)
		}

		reports := api.Group("/reports")
		reports.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			reports.GET("/revenue", reportHandler.GetRevenueReport)
			reports.GET("/technician-performance", reportHandler.GetTechnicianPerformance)
			reports.GET("/service-ranking", reportHandler.GetServiceRanking)
			reports.GET("/full", reportHandler.GetFullReport)
			reports.GET("/export/excel", reportHandler.ExportExcel)
			reports.GET("/export/pdf", reportHandler.ExportPDF)
		}

		notifications := api.Group("/notifications")
		notifications.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			notifications.GET("", notificationHandler.GetMyNotifications)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
		}

		audits := api.Group("/audits")
		audits.Use(middleware.JWTAuth(cfg.JWT.Secret), middleware.RoleAuth("admin"))
		{
			audits.GET("", auditHandler.List)
		}

		reviews := api.Group("/reviews")
		reviews.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			reviews.POST("", reviewHandler.Create)
			reviews.GET("/technician/:id", reviewHandler.GetByTechnicianID)
		}
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
