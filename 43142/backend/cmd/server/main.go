package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"recruitment-platform/internal/config"
	"recruitment-platform/internal/handlers"
	"recruitment-platform/internal/middleware"
	"recruitment-platform/internal/models"
	"recruitment-platform/internal/repository"
	"recruitment-platform/internal/routes"
	"recruitment-platform/internal/services"
	"recruitment-platform/internal/utils"
	appLogger "recruitment-platform/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	if err := appLogger.InitLogger("./logs"); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer appLogger.CloseLogger()

	gin.SetMode(cfg.Server.Mode)

	db, err := gorm.Open(sqlite.Open(cfg.Database.SQLitePath), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.ApplicantProfile{},
		&models.Company{},
		&models.Job{},
		&models.Resume{},
		&models.Application{},
		&models.ApplicationHistory{},
		&models.Interview{},
		&models.JobViewLog{},
		&models.DailyStatistics{},
	)

	appLogger.Info("数据库迁移完成")

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	resumeRepo := repository.NewResumeRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	interviewRepo := repository.NewInterviewRepository(db)
	statsRepo := repository.NewStatisticsRepository(db)

	emailService := utils.NewEmailService(utils.EmailConfig{
		SMTPHost: cfg.Email.SMTPHost,
		SMTPPort: cfg.Email.SMTPPort,
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
		From:     cfg.Email.From,
	})

	userService := services.NewUserService(userRepo)
	jobService := services.NewJobService(jobRepo)
	resumeService := services.NewResumeService(resumeRepo)
	applicationService := services.NewApplicationService(applicationRepo, jobRepo, resumeRepo)
	interviewService := services.NewInterviewService(interviewRepo, applicationRepo, jobRepo, emailService)
	statsService := services.NewStatisticsService(statsRepo, applicationRepo, jobRepo)

	userHandler := handlers.NewUserHandler(userService)
	jobHandler := handlers.NewJobHandler(jobService)
	resumeHandler := handlers.NewResumeHandler(resumeService)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	interviewHandler := handlers.NewInterviewHandler(interviewService)
	statsHandler := handlers.NewStatsHandler(statsService)

	r := gin.New()

	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())

	routes.SetupRoutes(r, userHandler, jobHandler, resumeHandler, applicationHandler, interviewHandler, statsHandler)

	r.Static("/uploads", "./uploads")

	go startCleanupRoutine()

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	appLogger.Info("服务器启动于 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

func startCleanupRoutine() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		registerLimiter := middleware.GetRateLimiter("register")
		registerLimiter.Cleanup(2 * time.Hour)

		loginLimiter := middleware.GetRateLimiter("login")
		loginLimiter.Cleanup(1 * time.Hour)
	}
}
