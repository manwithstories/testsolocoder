package routes

import (
	"github.com/gin-gonic/gin"
	"temp-staff-platform/handlers"
	"temp-staff-platform/middleware"
	"temp-staff-platform/models"
)

func SetupRoutes(r *gin.Engine) {
	authHandler := handlers.NewAuthHandler()
	jobHandler := handlers.NewJobHandler()
	scheduleHandler := handlers.NewScheduleHandler()
	checkInHandler := handlers.NewCheckInHandler()
	salaryHandler := handlers.NewSalaryHandler()
	evaluationHandler := handlers.NewEvaluationHandler()
	statsHandler := handlers.NewStatsHandler()
	templateHandler := handlers.NewJobTemplateHandler()
	matchHandler := handlers.NewMatchHandler()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   gin.H{},
		})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/auth/me", authHandler.GetCurrentUser)
			protected.PUT("/auth/profile", authHandler.UpdateProfile)
			protected.PUT("/auth/password", authHandler.ChangePassword)

			jobs := protected.Group("/jobs")
			{
				jobs.GET("", jobHandler.GetJobs)
				jobs.GET("/mine", jobHandler.GetMyJobs)
				jobs.GET("/:id", middleware.ParseUUIDParam("id"), jobHandler.GetJob)
				jobs.POST("", middleware.RoleAuth(models.RoleEmployer), jobHandler.CreateJob)
				jobs.PUT("/:id", middleware.ParseUUIDParam("id"), jobHandler.UpdateJob)
				jobs.DELETE("/:id", middleware.ParseUUIDParam("id"), jobHandler.DeleteJob)
				jobs.POST("/:id/apply", middleware.ParseUUIDParam("id"), middleware.RoleAuth(models.RoleTemporary), jobHandler.ApplyJob)
				jobs.GET("/:id/applications", middleware.ParseUUIDParam("id"), middleware.RoleAuth(models.RoleEmployer, models.RoleAgent), jobHandler.GetJobApplications)
			}

			applications := protected.Group("/applications")
			{
				applications.GET("/mine", middleware.RoleAuth(models.RoleTemporary), jobHandler.GetMyApplications)
				applications.PUT("/:id/review", middleware.ParseUUIDParam("id"), middleware.RoleAuth(models.RoleEmployer, models.RoleAgent), jobHandler.ReviewApplication)
			}

			schedules := protected.Group("/schedules")
			{
				schedules.GET("", scheduleHandler.GetSchedules)
				schedules.GET("/mine", middleware.RoleAuth(models.RoleTemporary), scheduleHandler.GetMySchedules)
				schedules.GET("/:id", middleware.ParseUUIDParam("id"), scheduleHandler.GetSchedule)
				schedules.POST("", middleware.RoleAuth(models.RoleEmployer, models.RoleAgent), scheduleHandler.CreateSchedule)
				schedules.POST("/batch", middleware.RoleAuth(models.RoleEmployer, models.RoleAgent), scheduleHandler.BatchCreateSchedules)
				schedules.PUT("/:id", middleware.ParseUUIDParam("id"), scheduleHandler.UpdateSchedule)
				schedules.DELETE("/:id", middleware.ParseUUIDParam("id"), scheduleHandler.DeleteSchedule)
				schedules.POST("/check-conflict", middleware.RoleAuth(models.RoleEmployer, models.RoleAgent), scheduleHandler.CheckConflict)
				schedules.GET("/export", scheduleHandler.ExportSchedules)
				schedules.GET("/:id/qrcode", middleware.ParseUUIDParam("id"), middleware.RoleAuth(models.RoleTemporary), checkInHandler.GenerateQRCode)
			}

			checkins := protected.Group("/checkins")
			{
				checkins.GET("", checkInHandler.GetCheckInRecords)
				checkins.GET("/stats", checkInHandler.GetCheckInStats)
				checkins.POST("", middleware.RoleAuth(models.RoleTemporary), checkInHandler.CheckIn)
				checkins.POST("/checkout", middleware.RoleAuth(models.RoleTemporary), checkInHandler.CheckOut)
				checkins.POST("/verify-face", checkInHandler.VerifyFace)
				checkins.POST("/register-face", middleware.RoleAuth(models.RoleTemporary), checkInHandler.RegisterFaceData)
			}

			salaries := protected.Group("/salaries")
			{
				salaries.GET("", salaryHandler.GetSalaryRecords)
				salaries.GET("/:id", middleware.ParseUUIDParam("id"), salaryHandler.GetSalaryRecord)
				salaries.POST("", middleware.RoleAuth(models.RoleEmployer, models.RoleAgent), salaryHandler.CalculateSalary)
				salaries.POST("/:id/pay", middleware.ParseUUIDParam("id"), middleware.RoleAuth(models.RoleEmployer), salaryHandler.PaySalary)
				salaries.POST("/batch-pay", middleware.RoleAuth(models.RoleEmployer), salaryHandler.BatchPaySalary)
				salaries.GET("/:id/export", middleware.ParseUUIDParam("id"), salaryHandler.ExportSalary)
			}

			evaluations := protected.Group("/evaluations")
			{
				evaluations.GET("", evaluationHandler.GetEvaluations)
				evaluations.GET("/mine", evaluationHandler.GetMyEvaluations)
				evaluations.GET("/:id/stats", middleware.ParseUUIDParam("id"), evaluationHandler.GetEvaluationStats)
				evaluations.POST("", evaluationHandler.CreateEvaluation)
				evaluations.PUT("/:id", middleware.ParseUUIDParam("id"), evaluationHandler.UpdateEvaluation)
			}

			templates := protected.Group("/job-templates")
			{
				templates.GET("", middleware.RoleAuth(models.RoleEmployer), templateHandler.GetTemplates)
				templates.GET("/:id", middleware.ParseUUIDParam("id"), templateHandler.GetTemplate)
				templates.POST("", middleware.RoleAuth(models.RoleEmployer), templateHandler.CreateTemplate)
				templates.PUT("/:id", middleware.ParseUUIDParam("id"), templateHandler.UpdateTemplate)
				templates.DELETE("/:id", middleware.ParseUUIDParam("id"), templateHandler.DeleteTemplate)
				templates.POST("/batch-import", middleware.RoleAuth(models.RoleEmployer), templateHandler.BatchImportTemplates)
				templates.POST("/apply", middleware.RoleAuth(models.RoleEmployer), templateHandler.ApplyTemplate)
			}

			match := protected.Group("/match")
			{
				match.POST("/temporaries", middleware.RoleAuth(models.RoleEmployer, models.RoleAgent), matchHandler.MatchTemporaries)
				match.POST("/quick-assign", middleware.RoleAuth(models.RoleEmployer, models.RoleAgent), matchHandler.QuickAssign)
				match.GET("/history", middleware.RoleAuth(models.RoleAgent), matchHandler.GetMatchHistory)
			}

			stats := protected.Group("/stats")
			{
				stats.GET("/overview", statsHandler.GetOverviewStats)
				stats.GET("/activities", statsHandler.GetActivityStats)
				stats.GET("/personnel", statsHandler.GetPersonnelStats)
				stats.GET("/salary", statsHandler.GetSalaryStats)
			}
		}
	}
}
