package routes

import (
	"github.com/gin-gonic/gin"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/handlers"
	"e-learning-platform/internal/middleware"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	userHandler := handlers.NewUserHandler(cfg)
	courseHandler := handlers.NewCourseHandler(cfg)
	progressHandler := handlers.NewProgressHandler(cfg)
	orderHandler := handlers.NewOrderHandler(cfg)
	qaHandler := handlers.NewQAHandler(cfg)
	analyticsHandler := handlers.NewAnalyticsHandler(cfg)
	reviewHandler := handlers.NewReviewHandler(cfg)
	fileHandler := handlers.NewFileHandler(cfg)

	api := r.Group("/api/v1")

	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/refresh", middleware.RefreshTokenAuth(), userHandler.RefreshToken)

	api.GET("/courses", courseHandler.ListCourses)
	api.GET("/courses/categories", courseHandler.ListCategories)
	api.GET("/courses/:id", courseHandler.GetCourse)

	api.GET("/reviews", reviewHandler.ListReviews)
	api.GET("/reviews/summary/:course_id", reviewHandler.GetCourseSummary)

	auth := api.Group("")
	auth.Use(middleware.JWTAuth(), middleware.TokenExpireCheck())
	{
		auth.GET("/auth/profile", userHandler.GetProfile)
		auth.PUT("/auth/profile", userHandler.UpdateProfile)
		auth.PUT("/auth/password", userHandler.ChangePassword)
		auth.POST("/auth/instructor-apply", userHandler.ApplyInstructor)

		auth.GET("/courses/my", middleware.RequireInstructor(), courseHandler.ListMyCourses)
		auth.POST("/courses", middleware.RequireInstructor(), courseHandler.CreateCourse)
		auth.PUT("/courses/:id", courseHandler.UpdateCourse)
		auth.PUT("/courses/:id/status", courseHandler.UpdateCourseStatus)
		auth.DELETE("/courses/:id", courseHandler.DeleteCourse)
		auth.PUT("/courses/:id/update-hours", courseHandler.UpdateCourseHours)

		auth.POST("/courses/:course_id/chapters", middleware.RequireInstructor(), courseHandler.CreateChapter)
		auth.PUT("/chapters/:id", courseHandler.UpdateChapter)
		auth.DELETE("/chapters/:id", courseHandler.DeleteChapter)

		auth.POST("/chapters/:chapter_id/lessons", middleware.RequireInstructor(), courseHandler.CreateLesson)
		auth.PUT("/lessons/:id", courseHandler.UpdateLesson)
		auth.DELETE("/lessons/:id", courseHandler.DeleteLesson)

		auth.POST("/lessons/:lesson_id/quiz", middleware.RequireInstructor(), courseHandler.CreateQuiz)
		auth.POST("/quizzes/:quiz_id/submit", courseHandler.SubmitQuiz)

		auth.PUT("/progress/lessons/:lesson_id", progressHandler.UpdateProgress)
		auth.GET("/progress/courses/:course_id", progressHandler.GetCourseProgress)
		auth.GET("/progress/lessons/:lesson_id", progressHandler.GetLessonProgress)

		auth.POST("/notes", progressHandler.CreateNote)
		auth.GET("/notes", progressHandler.ListNotes)
		auth.PUT("/notes/:id", progressHandler.UpdateNote)
		auth.DELETE("/notes/:id", progressHandler.DeleteNote)

		auth.POST("/orders", orderHandler.CreateOrder)
		auth.POST("/orders/:id/pay", orderHandler.PayOrder)
		auth.GET("/orders/:id", orderHandler.GetOrder)
		auth.GET("/orders/my", orderHandler.ListMyOrders)
		auth.POST("/orders/:id/refund", orderHandler.ApplyRefund)
		auth.GET("/coupons/validate", orderHandler.ValidateCoupon)

		auth.POST("/questions", qaHandler.CreateQuestion)
		auth.GET("/questions", qaHandler.ListQuestions)
		auth.GET("/questions/:id", qaHandler.GetQuestion)
		auth.DELETE("/questions/:id", qaHandler.DeleteQuestion)
		auth.POST("/questions/:id/like", qaHandler.LikeQuestion)

		auth.POST("/answers", qaHandler.CreateAnswer)
		auth.POST("/answers/:answer_id/best", qaHandler.MarkBestAnswer)
		auth.POST("/answers/:answer_id/like", qaHandler.LikeAnswer)
		auth.DELETE("/answers/:answer_id", qaHandler.DeleteAnswer)

		auth.POST("/reviews", reviewHandler.CreateReview)
		auth.GET("/reviews/my", reviewHandler.GetMyReview)
		auth.DELETE("/reviews/:id", reviewHandler.DeleteReview)

		auth.GET("/analytics/instructor/dashboard", middleware.RequireInstructor(), analyticsHandler.GetInstructorDashboard)
		auth.GET("/analytics/instructor/revenue", middleware.RequireInstructor(), analyticsHandler.GetInstructorRevenueReport)
		auth.GET("/analytics/instructor/export", middleware.RequireInstructor(), analyticsHandler.ExportInstructorReport)

		auth.GET("/analytics/instructor/stats", middleware.RequireInstructor(), userHandler.GetInstructorStats)

		auth.POST("/upload", fileHandler.Upload)
		auth.GET("/files", fileHandler.ListFiles)
		auth.DELETE("/files/:id", fileHandler.DeleteFile)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.RequireAdmin())
	{
		admin.GET("/users", userHandler.ListUsers)
		admin.PUT("/users/:id/status", userHandler.UpdateUserStatus)

		admin.GET("/instructor-applications", userHandler.ListInstructorApplications)
		admin.PUT("/instructor-applications/:id/review", userHandler.ReviewInstructor)

		admin.GET("/orders", orderHandler.ListAllOrders)
		admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
		admin.PUT("/orders/:id/refund/process", orderHandler.ProcessRefund)

		admin.POST("/coupons", orderHandler.CreateCoupon)
		admin.GET("/coupons", orderHandler.ListCoupons)
		admin.PUT("/coupons/:id", orderHandler.UpdateCoupon)
		admin.DELETE("/coupons/:id", orderHandler.DeleteCoupon)

		admin.GET("/analytics/dashboard", analyticsHandler.GetAdminDashboard)
		admin.GET("/analytics/revenue", analyticsHandler.GetAdminRevenueReport)
		admin.GET("/analytics/export", analyticsHandler.ExportAdminReport)
		admin.GET("/analytics/users/export", analyticsHandler.ExportUserReport)
		admin.GET("/analytics/stats", userHandler.GetAdminStats)
	}
}
