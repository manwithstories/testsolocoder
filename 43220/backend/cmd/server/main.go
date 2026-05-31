package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"pet-board/internal/config"
	"pet-board/internal/database"
	"pet-board/internal/handlers"
	"pet-board/internal/middleware"
	"pet-board/internal/models"
	"pet-board/internal/repository"
	"pet-board/internal/service"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	gin.SetMode(cfg.Server.Mode)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := runMigrations(database.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	uploadDir := cfg.Upload.UploadDir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	userRepo := repository.NewUserRepository()
	petRepo := repository.NewPetRepository()
	pkgRepo := repository.NewPackageRepository()
	resRepo := repository.NewReservationRepository()
	recordRepo := repository.NewDailyRecordRepository()
	reviewRepo := repository.NewReviewRepository()
	orderRepo := repository.NewOrderRepository()
	alertRepo := repository.NewHealthAlertRepository()
	statRepo := repository.NewStatisticsRepository()
	logRepo := repository.NewOperationLogRepository()

	userService := service.NewUserService(userRepo, cfg.JWT, logger)
	petService := service.NewPetService(petRepo, logger)
	pkgService := service.NewPackageService(pkgRepo, logger)
	resService := service.NewReservationService(resRepo, pkgRepo, petRepo, orderRepo, logger)
	recordService := service.NewDailyRecordService(recordRepo, resRepo, logger)
	reviewService := service.NewReviewService(reviewRepo, resRepo, userRepo, logger)
	orderService := service.NewOrderService(orderRepo, resRepo, logger)
	alertService := service.NewHealthAlertService(alertRepo, petRepo, logger)
	statService := service.NewStatisticsService(statRepo, logger)
	logService := service.NewOperationLogService(logRepo, logger)

	userHandler := handlers.NewUserHandler(userService)
	petHandler := handlers.NewPetHandler(petService)
	pkgHandler := handlers.NewPackageHandler(pkgService)
	resHandler := handlers.NewReservationHandler(resService)
	recordHandler := handlers.NewDailyRecordHandler(recordService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	orderHandler := handlers.NewOrderHandler(orderService)
	alertHandler := handlers.NewHealthAlertHandler(alertService)
	statHandler := handlers.NewStatisticsHandler(statService)
	exportHandler := handlers.NewExportHandler(statService)
	uploadHandler := handlers.NewUploadHandler(cfg.Upload)

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.OperationLogMiddleware(logService))

	r.Static("/uploads", uploadDir)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)

	api.GET("/packages", pkgHandler.ListByStore)
	api.GET("/packages/:id", pkgHandler.GetByID)

	api.GET("/stores/reviews", reviewHandler.ListByStore)
	api.GET("/keepers/reviews", reviewHandler.ListByKeeper)

	auth := api.Group("")
	auth.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		auth.GET("/auth/profile", userHandler.GetProfile)
		auth.PUT("/auth/profile", userHandler.UpdateProfile)
		auth.PUT("/auth/password", userHandler.ChangePassword)

		auth.POST("/upload", uploadHandler.Upload)
		auth.POST("/upload/multiple", uploadHandler.UploadMultiple)

		auth.POST("/pets", petHandler.Create)
		auth.GET("/pets", petHandler.ListByOwner)
		auth.GET("/pets/:id", petHandler.GetByID)
		auth.PUT("/pets/:id", petHandler.Update)
		auth.DELETE("/pets/:id", petHandler.Delete)
		auth.POST("/pets/vaccines", petHandler.AddVaccineRecord)
		auth.GET("/pets/:id/vaccines", petHandler.GetVaccineRecords)
		auth.POST("/pets/deworms", petHandler.AddDewormRecord)
		auth.GET("/pets/:id/deworms", petHandler.GetDewormRecords)

		auth.GET("/reservations", resHandler.List)
		auth.GET("/reservations/:id", resHandler.GetByID)

		auth.POST("/reviews", reviewHandler.Create)
		auth.GET("/reviews/:id", reviewHandler.GetByID)

		auth.GET("/orders", orderHandler.List)
		auth.GET("/orders/:id", orderHandler.GetByID)
		auth.GET("/reservations/:reservationId/orders", orderHandler.GetByReservationID)

		auth.GET("/alerts", alertHandler.List)
		auth.PUT("/alerts/:id/read", alertHandler.MarkAsRead)
		auth.PUT("/alerts/read-all", alertHandler.MarkAllAsRead)

		auth.GET("/daily-records/pet/:petId", recordHandler.ListByPet)
		auth.GET("/daily-records", recordHandler.ListByReservation)
		auth.GET("/daily-records/:id", recordHandler.GetByID)

		auth.GET("/statistics/revenue", statHandler.GetRevenueTrend)
		auth.GET("/statistics/occupancy", statHandler.GetOccupancyRate)
		auth.GET("/statistics/pet-types", statHandler.GetPetTypeDistribution)
		auth.GET("/statistics/orders", statHandler.GetOrderStatistics)

		auth.GET("/export/excel", exportHandler.ExportExcel)
		auth.GET("/export/pdf", exportHandler.ExportPDF)
	}

	ownerAuth := api.Group("")
	ownerAuth.Use(middleware.JWTAuth(cfg.JWT.Secret))
	ownerAuth.Use(middleware.RBAC(string(models.RoleOwner)))
	{
		ownerAuth.POST("/reservations", resHandler.Create)
		ownerAuth.PUT("/reservations/:id/cancel", resHandler.Cancel)
		ownerAuth.POST("/orders/pay", orderHandler.Pay)
	}

	storeAuth := api.Group("")
	storeAuth.Use(middleware.JWTAuth(cfg.JWT.Secret))
	storeAuth.Use(middleware.RBAC(string(models.RoleStore)))
	{
		storeAuth.PUT("/store/info", userHandler.UpdateStoreInfo)
		storeAuth.POST("/packages", pkgHandler.Create)
		storeAuth.PUT("/packages/:id", pkgHandler.Update)
		storeAuth.DELETE("/packages/:id", pkgHandler.Delete)

		storeAuth.PUT("/reservations/:id/confirm", resHandler.Confirm)
		storeAuth.PUT("/reservations/:id/check-in", resHandler.CheckIn)
		storeAuth.PUT("/reservations/:id/check-out", resHandler.CheckOut)

		storeAuth.POST("/orders/settle", orderHandler.Settle)
		storeAuth.POST("/orders/:id/refund", orderHandler.Refund)

		storeAuth.PUT("/reviews/:id/reply", reviewHandler.Reply)

		storeAuth.GET("/admin/users", userHandler.ListUsers)
	}

	keeperAuth := api.Group("")
	keeperAuth.Use(middleware.JWTAuth(cfg.JWT.Secret))
	keeperAuth.Use(middleware.RBAC(string(models.RoleKeeper)))
	{
		keeperAuth.PUT("/keeper/info", userHandler.UpdateKeeperInfo)
		keeperAuth.POST("/daily-records", recordHandler.Create)
		keeperAuth.PUT("/daily-records/:id", recordHandler.Update)
	}

	go alertService.CheckAndCreateAlerts()

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("Upload directory: %s", filepath.Abs(uploadDir))

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations(db *gorm.DB) error {
	tables := []interface{}{
		&models.User{},
		&models.StoreInfo{},
		&models.KeeperInfo{},
		&models.Pet{},
		&models.VaccineRecord{},
		&models.DewormRecord{},
		&models.BoardingPackage{},
		&models.Reservation{},
		&models.DailyRecord{},
		&models.Review{},
		&models.Order{},
		&models.HealthAlert{},
		&models.OperationLog{},
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", table, err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
