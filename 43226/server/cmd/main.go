package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"museum-server/internal/config"
	"museum-server/internal/handlers"
	"museum-server/internal/middleware"
	"museum-server/internal/models"
	"museum-server/internal/repository"
	"museum-server/internal/services"
	databasePkg "museum-server/pkg/database"
	loggerPkg "museum-server/pkg/logger"
	redisPkg "museum-server/pkg/redis"
)

func main() {
	cfgPath := "config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := loggerPkg.Init("logs"); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer loggerPkg.Close()

	db, err := databasePkg.Init(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	loggerPkg.Info("Database connected successfully")

	if err := redisPkg.Init(&cfg.Redis); err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}
	loggerPkg.Info("Redis connected successfully")

	db.AutoMigrate(
		&models.Museum{},
		&models.User{},
		&models.CollectionCategory{},
		&models.Collection{},
		&models.CollectionTag{},
		&models.Exhibition{},
		&models.ExhibitionCollection{},
		&models.TimeSlot{},
		&models.Reservation{},
		&models.VisitRecord{},
		&models.GuideSchedule{},
		&models.GuideContent{},
		&models.ResearchApplication{},
		&models.Statistic{},
	)
	loggerPkg.Info("Database migrated successfully")

	os.MkdirAll(cfg.Server.UploadDir, 0755)

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", cfg.Server.UploadDir)

	userRepo := repository.NewUserRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)
	exhibitionRepo := repository.NewExhibitionRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	guideRepo := repository.NewGuideRepository(db)
	researchRepo := repository.NewResearchRepository(db)
	museumRepo := repository.NewMuseumRepository(db)

	userService := services.NewUserService(userRepo, &cfg.JWT)
	collectionService := services.NewCollectionService(collectionRepo)
	exhibitionService := services.NewExhibitionService(exhibitionRepo)
	reservationService := services.NewReservationService(reservationRepo, exhibitionRepo, db)
	guideService := services.NewGuideService(guideRepo)
	researchService := services.NewResearchService(researchRepo)
	statisticsService := services.NewStatisticsService(reservationRepo, exhibitionRepo)
	museumService := services.NewMuseumService(museumRepo)

	userHandler := handlers.NewUserHandler(userService)
	collectionHandler := handlers.NewCollectionHandler(collectionService)
	exhibitionHandler := handlers.NewExhibitionHandler(exhibitionService)
	reservationHandler := handlers.NewReservationHandler(reservationService)
	guideHandler := handlers.NewGuideHandler(guideService)
	researchHandler := handlers.NewResearchHandler(researchService)
	statisticsHandler := handlers.NewStatisticsHandler(statisticsService)
	museumHandler := handlers.NewMuseumHandler(museumService)
	uploadHandler := handlers.NewUploadHandler(cfg.Server.UploadDir)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		museums := api.Group("/museums")
		{
			museums.GET("", museumHandler.List)
			museums.GET("/:id", museumHandler.GetByID)
			museums.POST("", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), museumHandler.Create)
			museums.PUT("/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), museumHandler.Update)
			museums.DELETE("/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), museumHandler.Delete)
		}

		collections := api.Group("/collections")
		{
			collections.GET("", collectionHandler.List)
			collections.GET("/:id", collectionHandler.GetByID)
			collections.POST("", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), collectionHandler.Create)
			collections.PUT("/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), collectionHandler.Update)
			collections.DELETE("/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), collectionHandler.Delete)
			collections.POST("/batch-import", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), collectionHandler.BatchImport)

			categories := collections.Group("/categories")
			{
				categories.GET("", collectionHandler.ListCategories)
				categories.POST("", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), collectionHandler.CreateCategory)
				categories.PUT("/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), collectionHandler.UpdateCategory)
				categories.DELETE("/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), collectionHandler.DeleteCategory)
			}

			tags := collections.Group("/tags")
			{
				tags.GET("", collectionHandler.ListTags)
				tags.POST("", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), collectionHandler.CreateTag)
			}
		}

		exhibitions := api.Group("/exhibitions")
		{
			exhibitions.GET("/hot", exhibitionHandler.GetHotExhibitions)
			exhibitions.GET("", exhibitionHandler.List)
			exhibitions.GET("/:id", exhibitionHandler.GetByID)
			exhibitions.POST("", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), exhibitionHandler.Create)
			exhibitions.PUT("/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), exhibitionHandler.Update)
			exhibitions.DELETE("/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), exhibitionHandler.Delete)
			exhibitions.POST("/:id/collections", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), exhibitionHandler.AddCollections)
			exhibitions.DELETE("/:id/collections", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), exhibitionHandler.RemoveCollections)
			exhibitions.GET("/:id/collections", exhibitionHandler.GetCollections)
			exhibitions.POST("/:id/time-slots", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), exhibitionHandler.CreateTimeSlot)
			exhibitions.POST("/time-slots/batch", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), exhibitionHandler.BatchCreateTimeSlots)
			exhibitions.GET("/:id/time-slots", exhibitionHandler.ListTimeSlots)
		}

		reservations := api.Group("/reservations")
		{
			reservations.GET("/my", middleware.JWTAuth(&cfg.JWT), reservationHandler.ListByUser)
			reservations.GET("/qr/:qr_code", reservationHandler.GetByQRCode)
			reservations.GET("/:id", middleware.JWTAuth(&cfg.JWT), reservationHandler.GetByID)
			reservations.POST("", middleware.JWTAuth(&cfg.JWT), reservationHandler.Create)
			reservations.PUT("/:id/confirm", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), reservationHandler.Confirm)
			reservations.PUT("/:id/cancel", middleware.JWTAuth(&cfg.JWT), reservationHandler.Cancel)
			reservations.PUT("/:id/reschedule", middleware.JWTAuth(&cfg.JWT), reservationHandler.Reschedule)
			reservations.GET("/exhibition/:id", middleware.JWTAuth(&cfg.JWT), reservationHandler.ListByExhibition)
			reservations.POST("/check-in", middleware.JWTAuth(&cfg.JWT), reservationHandler.CheckIn)
			reservations.PUT("/:id/check-out", middleware.JWTAuth(&cfg.JWT), reservationHandler.CheckOut)
			reservations.PUT("/:id/rate", middleware.JWTAuth(&cfg.JWT), reservationHandler.RateVisit)
			reservations.GET("/status/:id", reservationHandler.GetReservationStatus)
		}

		visits := api.Group("/visits")
		{
			visits.GET("/records", middleware.JWTAuth(&cfg.JWT), reservationHandler.ListVisitRecords)
			visits.GET("/stats", middleware.JWTAuth(&cfg.JWT), reservationHandler.GetUserVisitStats)
		}

		guides := api.Group("/guides")
		{
			guides.GET("/schedules", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleGuide), guideHandler.ListSchedules)
			guides.POST("/schedules", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleGuide), guideHandler.CreateSchedule)
			guides.PUT("/schedules/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleGuide), guideHandler.UpdateSchedule)
			guides.DELETE("/schedules/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleGuide), guideHandler.DeleteSchedule)

			guides.GET("/contents", guideHandler.ListContents)
			guides.POST("/contents", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin, models.UserRoleGuide), guideHandler.CreateContent)
			guides.PUT("/contents/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin, models.UserRoleGuide), guideHandler.UpdateContent)
			guides.DELETE("/contents/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin, models.UserRoleGuide), guideHandler.DeleteContent)
		}

		research := api.Group("/research")
		{
			research.GET("/applications/my", middleware.JWTAuth(&cfg.JWT), researchHandler.ListMyApplications)
			research.GET("/applications/:id", middleware.JWTAuth(&cfg.JWT), researchHandler.GetApplication)
			research.POST("/applications", middleware.JWTAuth(&cfg.JWT), researchHandler.CreateApplication)
			research.PUT("/applications/:id/review", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), researchHandler.ReviewApplication)
			research.GET("/applications", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), researchHandler.ListApplications)
		}

		statistics := api.Group("/statistics")
		{
			statistics.GET("", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), statisticsHandler.GetStatistics)
			statistics.GET("/export/excel", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), statisticsHandler.ExportExcel)
			statistics.GET("/export/pdf", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), statisticsHandler.ExportPDF)
		}

		users := api.Group("/users")
		{
			users.GET("/me", middleware.JWTAuth(&cfg.JWT), userHandler.GetProfile)
			users.PUT("/me", middleware.JWTAuth(&cfg.JWT), userHandler.UpdateProfile)
			users.GET("", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), userHandler.ListUsers)
			users.GET("/:id", middleware.JWTAuth(&cfg.JWT), middleware.RequireRole(models.UserRoleAdmin), userHandler.GetUserByID)
			users.GET("/guides/list", middleware.JWTAuth(&cfg.JWT), userHandler.GetGuides)
		}

		uploads := api.Group("/uploads")
		{
			uploads.POST("/image", middleware.JWTAuth(&cfg.JWT), uploadHandler.UploadImage)
		}
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	loggerPkg.Info("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
