package main

import (
	"fmt"
	"log"
	"smart-energy-platform/config"
	"smart-energy-platform/handlers"
	"smart-energy-platform/middleware"
	"smart-energy-platform/models"
	"smart-energy-platform/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	models.InitDB(cfg)
	services.InitRedis(cfg)

	go services.InitEnergyCollector()
	go services.InitScheduleExecutor()
	go services.InitSceneTriggerWatcher()

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		authenticated := api.Group("")
		authenticated.Use(middleware.Auth())
		{
			user := authenticated.Group("/user")
			{
				user.GET("/profile", handlers.GetProfile)
				user.PUT("/profile", handlers.UpdateProfile)
				user.PUT("/password", handlers.ChangePassword)
			}

			family := authenticated.Group("/families")
			{
				family.GET("", handlers.ListFamilies)
				family.POST("", handlers.CreateFamily)
				family.GET("/:id", handlers.GetFamily)
				family.PUT("/:id", handlers.UpdateFamily)
				family.DELETE("/:id", handlers.DeleteFamily)
				family.POST("/:id/invite", handlers.InviteMember)
				family.DELETE("/:id/members/:memberId", handlers.RemoveMember)
				family.PUT("/:id/members/:memberId/role", handlers.UpdateMemberRole)
			}

			invitation := authenticated.Group("/invitations")
			{
				invitation.GET("", handlers.ListInvitations)
				invitation.POST("/:id/accept", handlers.AcceptInvitation)
				invitation.POST("/:id/reject", handlers.RejectInvitation)
			}

			devices := authenticated.Group("/devices")
			{
				devices.GET("", handlers.ListDevices)
				devices.POST("", handlers.CreateDevice)
				devices.GET("/:id", handlers.GetDevice)
				devices.PUT("/:id", handlers.UpdateDevice)
				devices.DELETE("/:id", handlers.DeleteDevice)
				devices.PUT("/:id/status", handlers.UpdateDeviceStatus)
				devices.GET("/:id/energy", handlers.GetDeviceEnergy)
			}

			groups := authenticated.Group("/groups")
			{
				groups.GET("", handlers.ListGroups)
				groups.POST("", handlers.CreateGroup)
				groups.GET("/:id", handlers.GetGroup)
				groups.PUT("/:id", handlers.UpdateGroup)
				groups.DELETE("/:id", handlers.DeleteGroup)
				groups.POST("/:id/devices", handlers.AddDeviceToGroup)
				groups.DELETE("/:id/devices/:deviceId", handlers.RemoveDeviceFromGroup)
				groups.PUT("/:id/control", handlers.BatchControlGroup)
				groups.GET("/:id/energy", handlers.GetGroupEnergy)
			}

			energy := authenticated.Group("/energy")
			{
				energy.GET("/realtime", handlers.GetRealtimeEnergy)
				energy.GET("/statistics", handlers.GetEnergyStatistics)
				energy.GET("/trend", handlers.GetEnergyTrend)
				energy.GET("/alerts", handlers.ListEnergyAlerts)
				energy.GET("/export", handlers.ExportEnergyReport)
			}

			scenes := authenticated.Group("/scenes")
			{
				scenes.GET("", handlers.ListScenes)
				scenes.POST("", handlers.CreateScene)
				scenes.GET("/:id", handlers.GetScene)
				scenes.PUT("/:id", handlers.UpdateScene)
				scenes.DELETE("/:id", handlers.DeleteScene)
				scenes.POST("/:id/execute", handlers.ExecuteScene)
			}

			schedules := authenticated.Group("/schedules")
			{
				schedules.GET("", handlers.ListSchedules)
				schedules.POST("", handlers.CreateSchedule)
				schedules.GET("/:id", handlers.GetSchedule)
				schedules.PUT("/:id", handlers.UpdateSchedule)
				schedules.DELETE("/:id", handlers.DeleteSchedule)
				schedules.GET("/:id/logs", handlers.ListScheduleLogs)
			}

			notifications := authenticated.Group("/notifications")
			{
				notifications.GET("", handlers.ListNotifications)
				notifications.PUT("/:id/read", handlers.MarkNotificationRead)
				notifications.PUT("/read-all", handlers.MarkAllRead)
				notifications.DELETE("/:id", handlers.DeleteNotification)
			}
		}
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
