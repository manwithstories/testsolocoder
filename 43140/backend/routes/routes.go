package routes

import (
	"recruitment-platform/handlers"
	"recruitment-platform/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	authHandler := handlers.NewAuthHandler(db)
	jobHandler := handlers.NewJobHandler(db)
	departmentHandler := handlers.NewDepartmentHandler(db)
	resumeHandler := handlers.NewResumeHandler(db)
	applicationHandler := handlers.NewApplicationHandler(db)
	interviewHandler := handlers.NewInterviewHandler(db)
	reviewHandler := handlers.NewReviewHandler(db)
	searchHandler := handlers.NewSearchHandler(db)
	recommendationHandler := handlers.NewRecommendationHandler(db)
	statisticsHandler := handlers.NewStatisticsHandler(db)
	exportHandler := handlers.NewExportHandler(db)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		api.GET("/jobs/search", searchHandler.SearchJobs)
		api.GET("/jobs/:id", jobHandler.GetJob)
		api.GET("/jobs/:id/similar", recommendationHandler.GetSimilarJobs)
		api.GET("/job-types", searchHandler.GetJobTypes)
		api.GET("/locations", searchHandler.GetLocations)
		api.GET("/salary-ranges", searchHandler.GetSalaryRanges)
		api.GET("/skills", searchHandler.GetAllSkills)

		api.GET("/companies/:company_id/departments", departmentHandler.GetPublicDepartments)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/auth/profile", authHandler.GetProfile)
			auth.PUT("/auth/profile", authHandler.UpdateProfile)
			auth.PUT("/auth/password", authHandler.ChangePassword)

			auth.GET("/notifications", statisticsHandler.GetNotifications)
			auth.PUT("/notifications/:id/read", statisticsHandler.MarkNotificationRead)
			auth.PUT("/notifications/read-all", statisticsHandler.MarkAllNotificationsRead)

			company := auth.Group("")
			company.Use(middleware.CompanyMiddleware())
			{
				company.POST("/jobs", jobHandler.CreateJob)
				company.PUT("/jobs/:id", jobHandler.UpdateJob)
				company.DELETE("/jobs/:id", jobHandler.DeleteJob)
				company.GET("/company/jobs", jobHandler.ListCompanyJobs)
				company.PUT("/jobs/:id/status", jobHandler.UpdateJobStatus)
				company.GET("/company/jobs/stats", jobHandler.GetJobStats)
				company.GET("/company/jobs/:id/statistics", statisticsHandler.GetJobStatistics)

				company.POST("/departments", departmentHandler.CreateDepartment)
				company.PUT("/departments/:id", departmentHandler.UpdateDepartment)
				company.DELETE("/departments/:id", departmentHandler.DeleteDepartment)
				company.GET("/departments", departmentHandler.ListDepartments)

				company.GET("/company/applications", applicationHandler.ListCompanyApplications)
				company.GET("/company/applications/:id", applicationHandler.GetApplication)
				company.PUT("/company/applications/:id/status", applicationHandler.UpdateApplicationStatus)

				company.POST("/interviews", interviewHandler.CreateInterview)
				company.GET("/company/interviews", interviewHandler.ListCompanyInterviews)
				company.PUT("/interviews/:id", interviewHandler.UpdateInterview)
				company.PUT("/interviews/:id/cancel", interviewHandler.CancelInterview)

				company.POST("/reviews", reviewHandler.CreateReview)
				company.PUT("/reviews/:id", reviewHandler.UpdateReview)
				company.GET("/company/reviews", reviewHandler.ListCompanyReviews)

				company.GET("/company/statistics", statisticsHandler.GetCompanyStatistics)
				company.GET("/company/trends", statisticsHandler.GetApplicationTrend)

				company.GET("/export/applications", exportHandler.ExportApplications)
				company.GET("/export/interviews", exportHandler.ExportInterviews)
				company.GET("/export/jobs", exportHandler.ExportJobs)
			}

			jobSeeker := auth.Group("")
			jobSeeker.Use(middleware.JobSeekerMiddleware())
			{
				jobSeeker.POST("/resumes", resumeHandler.CreateResume)
				jobSeeker.PUT("/resumes/:id", resumeHandler.UpdateResume)
				jobSeeker.DELETE("/resumes/:id", resumeHandler.DeleteResume)
				jobSeeker.GET("/resumes/:id", resumeHandler.GetResume)
				jobSeeker.GET("/resumes", resumeHandler.ListResumes)
				jobSeeker.POST("/resumes/:id/upload", resumeHandler.UploadResumeFile)
				jobSeeker.PUT("/resumes/:id/default", resumeHandler.SetDefaultResume)

				jobSeeker.POST("/applications", applicationHandler.CreateApplication)
				jobSeeker.GET("/jobseeker/applications", applicationHandler.ListJobSeekerApplications)
				jobSeeker.DELETE("/applications/:id", applicationHandler.DeleteApplication)

				jobSeeker.GET("/jobseeker/interviews", interviewHandler.ListJobSeekerInterviews)
				jobSeeker.PUT("/interviews/:id/confirm", interviewHandler.ConfirmInterview)
				jobSeeker.PUT("/interviews/:id/decline", interviewHandler.DeclineInterview)

				jobSeeker.GET("/jobseeker/reviews", reviewHandler.ListJobSeekerReviews)
				jobSeeker.GET("/interviews/:interview_id/review", reviewHandler.GetReviewByInterview)

				jobSeeker.GET("/jobseeker/statistics", statisticsHandler.GetJobSeekerStatistics)
				jobSeeker.GET("/recommendations", recommendationHandler.GetRecommendedJobs)
			}

			admin := auth.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/admin/users", authHandler.ListUsers)
				admin.GET("/admin/users/:id", authHandler.GetUserByID)
				admin.PUT("/admin/users/:id/status", authHandler.UpdateUserStatus)
			}
		}
	}
}
