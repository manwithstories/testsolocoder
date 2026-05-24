package main

import (
	"fmt"
	"log"
	"matchmaking-platform/config"
	"matchmaking-platform/internal/handler"
	"matchmaking-platform/internal/middleware"
	"matchmaking-platform/internal/service"
	"matchmaking-platform/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	utils.InitDB()

	go handler.RunChatServer()

	go startCronJobs()

	gin.SetMode(config.Cfg.Server.Mode)
	r := gin.New()

	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Logger())

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", handler.NewUserHandler().Register)
		api.POST("/auth/login", handler.NewUserHandler().Login)
		api.GET("/auth/sms-code", handler.NewUserHandler().SendSmsCode)

		api.GET("/matchmakers", handler.NewMatchmakerHandler().ListAll)

		api.GET("/member/benefits", handler.NewMemberHandler().GetBenefits)
	}

	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user/info", handler.NewUserHandler().GetUserInfo)
		auth.PUT("/user/profile", handler.NewUserHandler().UpdateProfile)
		auth.POST("/user/verify", handler.NewUserHandler().Verify)
		auth.POST("/user/avatar", handler.NewUserHandler().UploadAvatar)
		auth.POST("/user/photos", handler.NewUserHandler().UploadPhotos)
		auth.GET("/user/:id", handler.NewUserHandler().GetUserProfile)

		auth.GET("/match/smart", handler.NewMatchHandler().SmartMatch)
		auth.GET("/match/filter", handler.NewMatchHandler().FilterMatch)
		auth.POST("/match/:id/favorite", handler.NewMatchHandler().Favorite)
		auth.POST("/match/:id/block", handler.NewMatchHandler().Block)
		auth.GET("/match/favorites", handler.NewMatchHandler().GetFavorites)
		auth.GET("/match/blocked", handler.NewMatchHandler().GetBlocked)

		auth.POST("/dates", handler.NewDateHandler().CreateInvite)
		auth.POST("/dates/:id/accept", handler.NewDateHandler().Accept)
		auth.POST("/dates/:id/reject", handler.NewDateHandler().Reject)
		auth.POST("/dates/:id/cancel", handler.NewDateHandler().Cancel)
		auth.POST("/dates/:id/complete", handler.NewDateHandler().Complete)
		auth.GET("/dates", handler.NewDateHandler().GetUserDates)
		auth.POST("/dates/reviews", handler.NewDateHandler().CreateReview)
		auth.GET("/dates/reviews", handler.NewDateHandler().GetReviews)

		auth.POST("/chat/send", handler.NewChatHandler().SendMessage)
		auth.GET("/chat/history/:id", handler.NewChatHandler().GetHistory)
		auth.GET("/chat/unread", handler.NewChatHandler().GetUnreadCount)
		auth.GET("/chat/sessions", handler.NewChatHandler().GetSessions)
		auth.POST("/chat/:id/read", handler.NewChatHandler().MarkAsRead)
		auth.POST("/chat/upload", handler.NewChatHandler().UploadChatFile)

		auth.POST("/member/orders", handler.NewMemberHandler().CreateOrder)
		auth.POST("/member/orders/:id/pay", handler.NewMemberHandler().PayOrder)
		auth.GET("/member/orders", handler.NewMemberHandler().GetUserOrders)
		auth.GET("/member/interact-limit", handler.NewMemberHandler().CheckInteractLimit)
	}

	verified := auth.Group("")
	verified.Use(middleware.VerifiedMiddleware())
	{
	}

	matchmaker := auth.Group("/matchmaker")
	matchmaker.Use(middleware.MatchmakerMiddleware())
	{
		matchmaker.POST("/members", handler.NewMatchmakerHandler().AddMember)
		matchmaker.DELETE("/members/:id", handler.NewMatchmakerHandler().RemoveMember)
		matchmaker.GET("/members", handler.NewMatchmakerHandler().ListMembers)
		matchmaker.POST("/services", handler.NewMatchmakerHandler().CreateService)
		matchmaker.PUT("/services/:id/progress", handler.NewMatchmakerHandler().UpdateProgress)
		matchmaker.GET("/services", handler.NewMatchmakerHandler().ListServices)
		matchmaker.GET("/stats", handler.NewMatchmakerHandler().GetStats)
	}

	admin := auth.Group("/admin")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("/users", handler.NewUserHandler().ListUsers)
		admin.PUT("/users/:id/disable", handler.NewUserHandler().DisableUser)
		admin.PUT("/users/:id/enable", handler.NewUserHandler().EnableUser)
		admin.PUT("/users/:id/verify/approve", handler.NewUserHandler().ApproveVerify)
		admin.PUT("/users/:id/verify/reject", handler.NewUserHandler().RejectVerify)

		admin.GET("/stats/platform", handler.NewStatsHandler().GetPlatformStats)
		admin.GET("/stats/daily", handler.NewStatsHandler().GetDailyStats)
		admin.GET("/stats/matchmaker", handler.NewStatsHandler().GetMatchmakerStats)
		admin.GET("/stats/export/excel", handler.NewStatsHandler().ExportExcel)
		admin.GET("/stats/export/pdf", handler.NewStatsHandler().ExportPDF)
		admin.GET("/logs", handler.NewStatsHandler().GetSystemLogs)
	}

	api.GET("/ws", handler.WebSocketHandler)
	api.GET("/ws/online", handler.GetOnlineUsers)

	addr := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	r.Run(addr)
}

func startCronJobs() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		memberService := service.NewMemberService()
		memberService.CheckAndDowngradeExpired()

		sendDateReminders()
	}
}

func sendDateReminders() {
	now := time.Now()
	oneHourLater := now.Add(1 * time.Hour)

	var dates []struct {
		ID          uint
		InitiatorID uint
		ReceiverID  uint
		Title       string
		DateAt      time.Time
	}

	utils.DB.Table("date_records").
		Where("status = ? AND date_at BETWEEN ? AND ? AND reminded = ?", "accepted", now, oneHourLater, false).
		Scan(&dates)

	for _, d := range dates {
		handler.NotifyUser(d.InitiatorID, "date_reminder", map[string]interface{}{
			"date_id": d.ID,
			"title":   d.Title,
			"date_at": d.DateAt,
		})
		handler.NotifyUser(d.ReceiverID, "date_reminder", map[string]interface{}{
			"date_id": d.ID,
			"title":   d.Title,
			"date_at": d.DateAt,
		})

		utils.DB.Table("date_records").Where("id = ?", d.ID).Update("reminded", true)
	}
}
