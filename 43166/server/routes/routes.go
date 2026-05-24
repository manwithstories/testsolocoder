package routes

import (
	"github.com/gin-gonic/gin"
	"business-registration-platform/controllers"
	"business-registration-platform/middleware"
	"business-registration-platform/models"
	"business-registration-platform/services"
)

func SetupRoutes(r *gin.Engine) {
	authCtrl := controllers.NewAuthController()
	appCtrl := controllers.NewApplicationController()
	processCtrl := controllers.NewProcessController()
	feeCtrl := controllers.NewFeeController()
	notificationCtrl := controllers.NewNotificationController()
	statisticsCtrl := controllers.NewStatisticsController()
	exportCtrl := controllers.NewExportController()
	agentCtrl := controllers.NewAgentController()
	agentService := services.NewAgentService()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
			auth.POST("/logout", authCtrl.Logout)
		}

		authRequired := api.Group("")
		authRequired.Use(middleware.JWTAuthMiddleware())
		{
			user := authRequired.Group("/user")
			{
				user.GET("/profile", authCtrl.GetProfile)
				user.PUT("/profile", authCtrl.UpdateProfile)
				user.PUT("/password", authCtrl.ChangePassword)
			}

			applications := authRequired.Group("/applications")
			{
				applications.POST("", appCtrl.CreateApplication)
				applications.GET("", appCtrl.GetApplicationList)
				applications.GET("/:id", appCtrl.GetApplication)
				applications.PUT("/:id", appCtrl.UpdateApplication)
				applications.POST("/:id/submit", appCtrl.SubmitApplication)
				applications.POST("/:id/cancel", appCtrl.CancelApplication)
				applications.POST("/:id/upload", appCtrl.UploadMaterials)
				applications.POST("/:id/review", appCtrl.ReviewApplication)
				applications.POST("/:id/assign", appCtrl.AssignAgent)
			}

			process := authRequired.Group("/applications/:applicationId/steps")
			{
				process.GET("", processCtrl.GetProcessSteps)
				process.PUT("/:stepId", processCtrl.UpdateProcessStep)
				process.POST("/:stepId/start", processCtrl.StartStep)
				process.POST("/:stepId/complete", processCtrl.CompleteStep)
				process.POST("/:stepId/skip", processCtrl.SkipStep)
				process.POST("/:stepId/upload", processCtrl.UploadCertificate)
			}

			fees := authRequired.Group("/fees")
			{
				fees.POST("/calculate", feeCtrl.CalculateFee)
				fees.POST("", feeCtrl.CreateApplicationFee)
				fees.POST("/pay", feeCtrl.PayFee)
				fees.GET("", feeCtrl.GetFeeList)
				fees.GET("/:applicationId", feeCtrl.GetApplicationFee)
			}

			notifications := authRequired.Group("/notifications")
			{
				notifications.GET("", notificationCtrl.GetUserNotifications)
				notifications.PUT("/:id/read", notificationCtrl.MarkAsRead)
				notifications.PUT("/read-all", notificationCtrl.MarkAllAsRead)
				notifications.GET("/unread-count", notificationCtrl.GetUnreadCount)
				notifications.POST("/send", notificationCtrl.SendNotification)
			}

			exports := authRequired.Group("/exports")
			{
				exports.POST("", exportCtrl.CreateExportTask)
				exports.GET("", exportCtrl.GetExportTasks)
				exports.GET("/:id", exportCtrl.GetExportTask)
				exports.GET("/:id/download", exportCtrl.DownloadExport)
			}
		}

		adminRequired := api.Group("")
		adminRequired.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware(models.RoleAdmin))
		{
			admin := adminRequired.Group("/admin")
			{
				agents := admin.Group("/agents")
				{
					agents.GET("", agentCtrl.GetAgentList)
					agents.GET("/available", agentCtrl.GetAvailableAgents)
					agents.GET("/:id", agentCtrl.GetAgent)
					agents.POST("", agentCtrl.CreateAgent)
					agents.PUT("/:id", agentCtrl.UpdateAgentProfile)
					agents.DELETE("/:id", agentCtrl.DeleteAgent)
					agents.POST("/:id/schedule", agentCtrl.UpdateAgentWorkSchedule)
					agents.PUT("/:id/max-apps", agentCtrl.UpdateAgentMaxApps)
					agents.GET("/:id/stats", agentCtrl.GetAgentStats)
					agents.GET("/:id/performance", agentCtrl.GetAgentPerformanceReport)
					agents.POST("/auto-assign/:applicationId", agentCtrl.AutoAssignAgent)
				}

				feeStandards := admin.Group("/fee-standards")
				{
					feeStandards.GET("", feeCtrl.GetFeeStandards)
					feeStandards.POST("", feeCtrl.CreateFeeStandard)
					feeStandards.PUT("/:id", feeCtrl.UpdateFeeStandard)
				}

				discounts := admin.Group("/discounts")
				{
					discounts.GET("", feeCtrl.GetDiscountPolicies)
					discounts.POST("", feeCtrl.CreateDiscountPolicy)
					discounts.PUT("/:id", feeCtrl.UpdateDiscountPolicy)
					discounts.DELETE("/:id", feeCtrl.DeleteDiscountPolicy)
				}

				templates := admin.Group("/notification-templates")
				{
					templates.GET("", notificationCtrl.GetNotificationTemplates)
					templates.POST("", notificationCtrl.CreateNotificationTemplate)
					templates.PUT("/:id", notificationCtrl.UpdateNotificationTemplate)
					templates.DELETE("/:id", notificationCtrl.DeleteNotificationTemplate)
				}

				statistics := admin.Group("/statistics")
				{
					statistics.GET("/overview", statisticsCtrl.GetOverviewStats)
					statistics.GET("/status-distribution", statisticsCtrl.GetStatusDistribution)
					statistics.GET("/company-type-distribution", statisticsCtrl.GetCompanyTypeDistribution)
					statistics.GET("/agent-performance", statisticsCtrl.GetAgentPerformance)
					statistics.GET("/application-time-series", statisticsCtrl.GetApplicationTimeSeries)
					statistics.GET("/revenue", statisticsCtrl.GetRevenueStats)
				}
			}
		}

		agentRequired := api.Group("")
		agentRequired.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware(models.RoleAgent))
		{
			agent := agentRequired.Group("/agent")
			{
				agent.GET("/applications", agentCtrl.GetAgentApplications)
				agent.GET("/stats", func(c *gin.Context) {
					userID, _ := c.Get("userID")
					stats, err := agentService.GetAgentStats(userID.(uint))
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					c.JSON(200, gin.H{"code": 200, "data": stats})
				})
			}
		}
	}
}
