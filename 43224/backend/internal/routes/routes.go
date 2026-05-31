package routes

import (
	"translation-platform/internal/handlers"
	"translation-platform/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.OperationLog())

	api := r.Group("/api")
	{
		api.POST("/auth/register", handlers.Register)
		api.POST("/auth/login", handlers.Login)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/auth/me", handlers.GetCurrentUser)
			auth.PUT("/auth/profile", handlers.UpdateProfile)
			auth.PUT("/auth/password", handlers.ChangePassword)

			auth.GET("/users", middleware.RoleAuth("admin"), handlers.ListUsers)
			auth.PUT("/users/:id/status", middleware.RoleAuth("admin"), handlers.UpdateUserStatus)

			auth.GET("/language-pairs", handlers.ListLanguagePairs)
			auth.POST("/language-pairs", middleware.RoleAuth("admin"), handlers.CreateLanguagePair)
			auth.GET("/expertise-tags", handlers.ListExpertiseTags)
			auth.POST("/expertise-tags", middleware.RoleAuth("admin"), handlers.CreateExpertiseTag)

			auth.GET("/projects", handlers.ListProjects)
			auth.GET("/projects/:id", handlers.GetProject)
			auth.POST("/projects", middleware.RoleAuth("client", "admin"), handlers.CreateProject)
			auth.PUT("/projects/:id/approve", middleware.RoleAuth("pm", "admin"), handlers.ApproveProject)
			auth.PUT("/projects/:id/assign", middleware.RoleAuth("pm", "admin"), handlers.AssignTranslator)
			auth.PUT("/projects/:id/start", middleware.RoleAuth("translator"), handlers.StartProject)
			auth.PUT("/projects/:id/submit", middleware.RoleAuth("translator"), handlers.SubmitForReview)
			auth.PUT("/projects/:id/complete", middleware.RoleAuth("pm", "admin"), handlers.CompleteProject)
			auth.PUT("/projects/:id/cancel", handlers.CancelProject)
			auth.POST("/projects/:id/comments", handlers.AddProjectComment)

			auth.GET("/projects/:id/recommend-translators", middleware.RoleAuth("pm", "admin"), handlers.RecommendTranslators)
			auth.POST("/projects/:id/auto-assign", middleware.RoleAuth("pm", "admin"), handlers.AutoAssignTranslator)

			auth.GET("/projects/:project_id/documents", handlers.ListDocuments)
			auth.POST("/projects/:project_id/documents", handlers.UploadDocument)
			auth.GET("/documents/:id", handlers.DownloadDocument)
			auth.DELETE("/documents/:id", handlers.DeleteDocument)
			auth.GET("/projects/:project_id/documents/versions", handlers.GetDocumentVersions)
			auth.POST("/documents/:id/extract-segments", handlers.ExtractDocumentSegments)

			auth.GET("/projects/:project_id/segments", handlers.GetProjectSegments)
			auth.PUT("/segments/:id", middleware.RoleAuth("translator"), handlers.UpdateSegmentTranslation)
			auth.GET("/memory/suggestions", handlers.GetMemorySuggestions)

			auth.GET("/translation-memories", handlers.ListTranslationMemories)
			auth.POST("/translation-memories", handlers.CreateTranslationMemory)
			auth.PUT("/translation-memories/:id", handlers.UpdateTranslationMemory)
			auth.DELETE("/translation-memories/:id", handlers.DeleteTranslationMemory)

			auth.GET("/glossary-terms", handlers.ListGlossaryTerms)
			auth.POST("/glossary-terms", handlers.CreateGlossaryTerm)
			auth.PUT("/glossary-terms/:id", handlers.UpdateGlossaryTerm)
			auth.DELETE("/glossary-terms/:id", handlers.DeleteGlossaryTerm)

			auth.GET("/review-tasks", handlers.ListReviewTasks)
			auth.POST("/review-tasks", middleware.RoleAuth("pm", "admin"), handlers.CreateReviewTask)
			auth.PUT("/review-tasks/:id", handlers.ProcessReview)
			auth.POST("/review-tasks/batch", handlers.BatchReview)
			auth.GET("/projects/:id/review-summary", handlers.GetProjectReviewSummary)

			auth.GET("/payments", handlers.ListPayments)
			auth.GET("/payments/:id", handlers.GetPayment)
			auth.PUT("/payments/:id/confirm", middleware.RoleAuth("client", "admin"), handlers.ConfirmPayment)
			auth.POST("/fee/calculate", handlers.CalculateProjectFee)
			auth.GET("/payments/statistics", handlers.GetPaymentStatistics)

			auth.GET("/statistics/projects", handlers.GetProjectStatistics)
			auth.GET("/statistics/translators", handlers.GetTranslatorStatistics)
			auth.GET("/statistics/revenue-trend", handlers.GetRevenueTrend)
			auth.GET("/statistics/language-pairs", handlers.GetLanguagePairStatistics)
			auth.GET("/statistics/export/excel", handlers.ExportStatisticsExcel)
			auth.GET("/statistics/export/pdf", handlers.ExportStatisticsPDF)

			auth.GET("/operation-logs", middleware.RoleAuth("admin"), handlers.ListOperationLogs)
		}
	}
}
