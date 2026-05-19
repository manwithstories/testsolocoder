package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/notification-center/internal/cache"
	"github.com/notification-center/internal/channels"
	"github.com/notification-center/internal/config"
	"github.com/notification-center/internal/database"
	"github.com/notification-center/internal/handlers"
	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/ratelimit"
	"github.com/notification-center/internal/router"
	"github.com/notification-center/internal/services"
	"go.uber.org/zap"
)

func main() {
	env := config.GetEnv()
	fmt.Printf("Starting notification center server in %s mode...\n", env)

	if err := config.Load(env); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(&config.AppConfig.Log); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	gin.SetMode(config.AppConfig.Server.Mode)

	if err := database.Init(&config.AppConfig.Database); err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
	}
	defer database.Close()

	if err := database.AutoMigrate(
		&models.Channel{},
		&models.Template{},
		&models.Recipient{},
		&models.Tag{},
		&models.RecipientGroup{},
		&models.Message{},
		&models.MessageQueue{},
		&models.Webhook{},
		&models.WebhookLog{},
	); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	if err := cache.Init(&config.AppConfig.Redis); err != nil {
		logger.Warn("Failed to init redis", zap.Error(err))
	}
	defer cache.Close()

	rateLimiter := ratelimit.GetInstance(&config.AppConfig.RateLimit)

	senderSvc := services.NewSenderService()
	senderSvc.RegisterAdapter(models.ChannelTypeEmail, channels.NewEmailAdapter())
	senderSvc.RegisterAdapter(models.ChannelTypeSMS, channels.NewSMSAdapter())
	senderSvc.RegisterAdapter(models.ChannelTypeWeChat, channels.NewWeChatAdapter())
	senderSvc.RegisterAdapter(models.ChannelTypeDingTalk, channels.NewDingTalkAdapter())
	senderSvc.RegisterAdapter(models.ChannelTypeWebhook, channels.NewWebhookChannelAdapter())

	channelSvc := services.NewChannelService(senderSvc)
	templateSvc := services.NewTemplateService(channelSvc)
	recipientSvc := services.NewRecipientService()
	webhookSvc := services.NewWebhookService(&config.AppConfig.Webhook)
	statisticsSvc := services.NewStatisticsService(channelSvc)
	queueSvc := services.NewQueueService(
		&config.AppConfig.Queue,
		&config.AppConfig.Retry,
		rateLimiter,
		channelSvc,
		senderSvc,
		webhookSvc,
	)

	queueSvc.Start()
	defer queueSvc.Stop()

	channelHandler := handlers.NewChannelHandler(channelSvc)
	templateHandler := handlers.NewTemplateHandler(templateSvc)
	recipientHandler := handlers.NewRecipientHandler(recipientSvc)
	messageHandler := handlers.NewMessageHandler(queueSvc, templateSvc, channelSvc, statisticsSvc)
	webhookHandler := handlers.NewWebhookHandler(webhookSvc)

	r := router.SetupRouter(
		channelHandler,
		templateHandler,
		recipientHandler,
		messageHandler,
		webhookHandler,
	)

	addr := config.AppConfig.Server.Addr()
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		logger.Info("server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited gracefully")
}
