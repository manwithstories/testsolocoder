package router

import (
	"event-platform/internal/config"
	"event-platform/internal/handler"
	"event-platform/internal/middleware"
	"event-platform/internal/queue"
	"event-platform/internal/repository"
	"event-platform/internal/service"
	"event-platform/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, cfg *config.Config, jm *jwt.Manager, q *queue.Queue) {
	userRepo := repository.NewUserRepo()
	eventRepo := repository.NewEventRepo()
	regRepo := repository.NewRegistrationRepo()
	scoreRepo := repository.NewScoreRepo()
	certRepo := repository.NewCertificateRepo()
	msgRepo := repository.NewMessageRepo()
	logRepo := repository.NewLogRepo()

	userSvc := service.NewUserService(userRepo, logRepo, jm)
	eventSvc := service.NewEventService(eventRepo, logRepo, regRepo)
	regSvc := service.NewRegistrationService(regRepo, eventRepo, userRepo, q, logRepo)
	scoreSvc := service.NewScoreService(scoreRepo, regRepo, eventRepo, certRepo, q, logRepo)
	certSvc := service.NewCertificateService(certRepo, scoreRepo, eventRepo, userRepo, q, cfg.Upload.CertificateDir)
	msgSvc := service.NewMessageService(msgRepo, q)

	userH := handler.NewUserHandler(userSvc)
	eventH := handler.NewEventHandler(eventSvc)
	regH := handler.NewRegistrationHandler(regSvc)
	scoreH := handler.NewScoreHandler(scoreSvc)
	scoreImportH := handler.NewScoreImportHandler(scoreSvc)
	certH := handler.NewCertificateHandler(certSvc, scoreSvc, userSvc, eventSvc, q, cfg.Upload.CertificateDir)
	msgH := handler.NewMessageHandler(msgSvc)
	statsH := handler.NewStatsHandler(regSvc, scoreSvc, eventSvc)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", userH.Register)
		api.POST("/auth/login", userH.Login)
		api.GET("/events", eventH.List)
		api.GET("/events/:id", eventH.Get)
	}

	auth := api.Group("", middleware.Auth(jm))
	{
		auth.GET("/users/me", userH.Profile)
		auth.PUT("/users/me/verify", userH.Verify)
		auth.GET("/registrations/me", regH.MyList)
		auth.POST("/registrations", regH.Register)
		auth.GET("/scores/me", scoreH.MyList)
		auth.GET("/certificates/me", certH.MyList)
		auth.GET("/certificates/:id/download", certH.Download)
		auth.POST("/certificates/:score_id/generate", certH.Generate)
		auth.GET("/messages", msgH.List)
		auth.PUT("/messages/:id/read", msgH.MarkRead)
		auth.PUT("/messages/read-all", msgH.MarkAllRead)
		auth.GET("/messages/unread-count", msgH.UnreadCount)
		auth.GET("/stats/overview", statsH.Overview)
	}

	admin := api.Group("/admin", middleware.Auth(jm), middleware.AdminOnly())
	{
		admin.GET("/users", userH.List)
		admin.POST("/events", eventH.Create)
		admin.PUT("/events/:id", eventH.Update)
		admin.PUT("/events/:id/publish", eventH.Publish)
		admin.PUT("/events/:id/unpublish", eventH.Unpublish)
		admin.GET("/events", eventH.List)
		admin.GET("/registrations/event/:event_id", regH.ListByEvent)
		admin.PUT("/registrations/:id/confirm", regH.ConfirmWaitlist)
		admin.POST("/scores", scoreH.Entry)
		admin.POST("/scores/import", scoreImportH.Import)
		admin.GET("/scores/import/template", scoreImportH.Template)
		admin.GET("/scores/item/:item_id", scoreH.ListByItem)
		admin.GET("/stats/export", statsH.Export)
	}
}
