package routes

import (
	"github.com/gin-gonic/gin"

	"recruitment-platform/internal/handlers"
	"recruitment-platform/internal/middleware"
	"recruitment-platform/internal/models"
)

func SetupRoutes(
	r *gin.Engine,
	userHandler *handlers.UserHandler,
	jobHandler *handlers.JobHandler,
	resumeHandler *handlers.ResumeHandler,
	applicationHandler *handlers.ApplicationHandler,
	interviewHandler *handlers.InterviewHandler,
	statsHandler *handlers.StatsHandler,
) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"code":    0,
				"message": "ok",
				"data":    gin.H{"status": "running"},
			})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			users := authenticated.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.GET("/company", userHandler.GetCompany)
				users.PUT("/company", userHandler.UpdateCompany)
			}

			adminUsers := authenticated.Group("/admin/users")
			adminUsers.Use(middleware.RoleMiddleware(models.RoleAdmin))
			{
				adminUsers.GET("", userHandler.ListUsers)
				adminUsers.PUT("/:id/status", userHandler.UpdateUserStatus)
			}
		}

		jobs := api.Group("/jobs")
		{
			jobs.GET("", jobHandler.ListJobs)
			jobs.GET("/:id", jobHandler.GetJob)
			jobs.GET("/:id/stats/views", jobHandler.GetViewStats)
		}

		companyJobs := api.Group("/company/jobs")
		companyJobs.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleCompany))
		{
			companyJobs.POST("", jobHandler.CreateJob)
			companyJobs.GET("", jobHandler.ListMyJobs)
			companyJobs.PUT("/:id", jobHandler.UpdateJob)
			companyJobs.DELETE("/:id", jobHandler.DeleteJob)
			companyJobs.PUT("/:id/publish", jobHandler.PublishJob)
			companyJobs.PUT("/:id/close", jobHandler.CloseJob)
			companyJobs.POST("/bulk/import", jobHandler.BulkImport)
			companyJobs.DELETE("/bulk", jobHandler.BulkDelete)
			companyJobs.GET("/export", jobHandler.ExportJobs)
		}

		resumes := api.Group("/resumes")
		resumes.Use(middleware.AuthMiddleware())
		{
			resumes.POST("", resumeHandler.CreateResume)
			resumes.GET("", resumeHandler.ListResumes)
			resumes.GET("/default", resumeHandler.GetDefaultResume)
			resumes.GET("/:id", resumeHandler.GetResume)
			resumes.PUT("/:id", resumeHandler.UpdateResume)
			resumes.DELETE("/:id", resumeHandler.DeleteResume)
			resumes.PUT("/:id/default", resumeHandler.SetDefaultResume)
			resumes.POST("/:id/upload", resumeHandler.UploadResumeFile)
		}

		resumeSearch := api.Group("/search/resumes")
		resumeSearch.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleCompany))
		{
			resumeSearch.GET("", resumeHandler.SearchResumes)
		}

		applications := api.Group("/applications")
		applications.Use(middleware.AuthMiddleware())
		{
			applications.POST("", applicationHandler.Apply)
			applications.GET("/my", applicationHandler.ListMyApplications)
			applications.GET("/:id", applicationHandler.GetApplication)
			applications.PUT("/:id/status", applicationHandler.UpdateStatus)
			applications.PUT("/:id/withdraw", applicationHandler.Withdraw)
			applications.GET("/:id/history", applicationHandler.GetHistory)
		}

		companyApplications := api.Group("/company/applications")
		companyApplications.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleCompany))
		{
			companyApplications.GET("", applicationHandler.ListCompanyApplications)
			companyApplications.GET("/job/:jobId", applicationHandler.ListJobApplications)
			companyApplications.PUT("/bulk/status", applicationHandler.BulkUpdateStatus)
			companyApplications.GET("/job/:jobId/status-count", applicationHandler.GetStatusCount)
		}

		interviews := api.Group("/interviews")
		interviews.Use(middleware.AuthMiddleware())
		{
			interviews.GET("/my", interviewHandler.ListMyInterviews)
			interviews.GET("/:id", interviewHandler.GetInterview)
			interviews.PUT("/:id/accept", interviewHandler.AcceptInterview)
			interviews.PUT("/:id/reject", interviewHandler.RejectInterview)
		}

		companyInterviews := api.Group("/company/interviews")
		companyInterviews.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleCompany))
		{
			companyInterviews.POST("", interviewHandler.ScheduleInterview)
			companyInterviews.GET("", interviewHandler.ListCompanyInterviews)
			companyInterviews.PUT("/:id", interviewHandler.UpdateInterview)
			companyInterviews.PUT("/:id/complete", interviewHandler.CompleteInterview)
			companyInterviews.PUT("/:id/cancel", interviewHandler.CancelInterview)
		}

		stats := api.Group("/stats")
		stats.Use(middleware.AuthMiddleware())
		{
			stats.GET("/daily", statsHandler.GetDailyStats)
			stats.GET("/range", statsHandler.GetDateRangeStats)
			stats.GET("/applications", statsHandler.GetApplicationStats)
		}

		companyStats := api.Group("/company/stats")
		companyStats.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware(models.RoleCompany))
		{
			companyStats.GET("/jobs", statsHandler.GetJobStats)
			companyStats.GET("/recruitment-cycle", statsHandler.GetRecruitmentCycleStats)
			companyStats.GET("/export/jobs", statsHandler.ExportJobStats)
			companyStats.GET("/export/applications", statsHandler.ExportApplicationStats)
		}
	}
}
