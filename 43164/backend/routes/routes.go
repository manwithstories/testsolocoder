package routes

import (
	"github.com/gin-gonic/gin"
	"tutoring-platform/handlers"
	"tutoring-platform/middleware"
	"tutoring-platform/models"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Tutoring Platform API is running",
		})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		subjects := api.Group("/subjects")
		{
			subjects.GET("", handlers.GetSubjects)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			users := protected.Group("/users")
			{
				users.GET("/me", handlers.GetProfile)
				users.PUT("/me", handlers.UpdateProfile)
				users.GET("/:id", handlers.GetUserByID)
			}

			teachers := protected.Group("/teachers")
			{
				teachers.GET("", handlers.ListTeachers)
				teachers.GET("/:id", handlers.GetTeacherByID)
				teachers.GET("/:id/reviews", handlers.GetTeacherReviews)
			}

			teacherProtected := protected.Group("/teacher")
			teacherProtected.Use(middleware.RoleMiddleware(models.RoleTeacher))
			{
				teacherProtected.GET("/profile", handlers.GetTeacherProfile)
				teacherProtected.PUT("/profile", handlers.CreateTeacherProfile)
				teacherProtected.POST("/subjects", handlers.AddTeacherSubject)
				teacherProtected.DELETE("/subjects/:subjectId", handlers.RemoveTeacherSubject)
				teacherProtected.POST("/availability", handlers.AddAvailabilitySlot)
				teacherProtected.DELETE("/availability/:slotId", handlers.RemoveAvailabilitySlot)
				teacherProtected.GET("/earnings", handlers.GetEarningsSummary)
			}

			students := protected.Group("/students")
			{
				students.GET("/:id", handlers.GetStudentByID)
			}

			studentProtected := protected.Group("/student")
			studentProtected.Use(middleware.RoleMiddleware(models.RoleStudent))
			{
				studentProtected.GET("/profile", handlers.GetStudentProfile)
				studentProtected.PUT("/profile", handlers.UpdateStudentProfile)
				studentProtected.POST("/goals", handlers.AddLearningGoal)
				studentProtected.PUT("/goals/:goalId", handlers.UpdateLearningGoal)
				studentProtected.DELETE("/goals/:goalId", handlers.DeleteLearningGoal)
				studentProtected.GET("/assessment/questions", handlers.GetAssessmentQuestions)
				studentProtected.POST("/assessment/submit", handlers.SubmitAssessment)
				studentProtected.GET("/assessment", handlers.GetMyAssessment)
				studentProtected.GET("/match-teachers", handlers.MatchTeachers)
				studentProtected.GET("/milestones", handlers.GetMilestones)
				studentProtected.POST("/milestones", handlers.CreateMilestone)
				studentProtected.PUT("/milestones/:id", handlers.UpdateMilestone)
			}

			bookings := protected.Group("/bookings")
			{
				bookings.GET("", handlers.GetBookings)
				bookings.GET("/:id", handlers.GetBookingByID)
				bookings.POST("", handlers.CreateBooking)
				bookings.POST("/:id/confirm", handlers.ConfirmBooking)
				bookings.POST("/reschedule", handlers.RescheduleBooking)
				bookings.POST("/cancel", handlers.CancelBooking)
				bookings.POST("/:id/complete", handlers.CompleteBooking)
			}

			video := protected.Group("/video")
			{
				video.POST("/sessions", handlers.CreateVideoSession)
				video.GET("/sessions/:id", handlers.GetVideoSession)
				video.GET("/sessions/booking/:bookingId", handlers.GetSessionByBooking)
				video.POST("/sessions/start", handlers.StartVideoSession)
				video.POST("/sessions/end", handlers.EndVideoSession)
				video.POST("/sessions/:sessionId/events", handlers.HandleSessionEvent)
				video.GET("/sessions/:id/quality", handlers.GetSessionQuality)
			}

			wallet := protected.Group("/wallet")
			{
				wallet.GET("", handlers.GetWallet)
				wallet.GET("/transactions", handlers.GetTransactions)
				wallet.GET("/transactions/:id", handlers.GetTransactionByID)
				wallet.POST("/deposit", handlers.Deposit)
				wallet.POST("/withdraw", handlers.Withdraw)
				wallet.GET("/withdraw-requests", handlers.GetWithdrawRequests)
			}

			notes := protected.Group("/notes")
			{
				notes.GET("", handlers.GetLessonNotes)
				notes.GET("/:id", handlers.GetLessonNoteByID)
				notes.POST("", handlers.CreateLessonNote)
				notes.PUT("/:id", handlers.UpdateLessonNote)
				notes.DELETE("/:id", handlers.DeleteLessonNote)
			}

			homework := protected.Group("/homework")
			{
				homework.GET("", handlers.GetHomeworks)
				homework.POST("", handlers.CreateHomework)
				homework.POST("/submit", handlers.SubmitHomework)
				homework.POST("/:id/grade", handlers.GradeHomework)
			}

			feedback := protected.Group("/feedback")
			{
				feedback.GET("", handlers.GetFeedbacks)
				feedback.POST("", handlers.CreateFeedback)
			}

			reviews := protected.Group("/reviews")
			{
				reviews.GET("", handlers.GetReviews)
				reviews.GET("/:id", handlers.GetReviewByID)
				reviews.POST("", handlers.CreateReview)
				reviews.POST("/:id/reply", handlers.ReplyToReview)
			}

			messages := protected.Group("/messages")
			{
				messages.GET("", handlers.GetMessages)
				messages.GET("/conversations", handlers.GetConversationList)
				messages.GET("/unread-count", handlers.GetUnreadCount)
				messages.POST("", handlers.SendMessage)
				messages.POST("/mark-read", handlers.MarkMessagesAsRead)
				messages.DELETE("/:id", handlers.DeleteMessage)
			}

			notifications := protected.Group("/notifications")
			{
				notifications.GET("", handlers.GetNotifications)
				notifications.GET("/unread-count", handlers.GetUnreadNotificationCount)
				notifications.PUT("/:id/read", handlers.MarkNotificationAsRead)
				notifications.PUT("/read-all", handlers.MarkAllNotificationsAsRead)
				notifications.DELETE("/:id", handlers.DeleteNotification)
			}

			admin := protected.Group("/admin")
			admin.Use(middleware.RoleMiddleware(models.RoleAdmin))
			{
				admin.GET("/stats", handlers.GetAdminStats)
				admin.GET("/logs", handlers.GetSystemLogs)
				admin.GET("/admin-actions", handlers.GetAdminActions)
				admin.GET("/pending-approvals", handlers.GetPendingTeacherApprovals)
				admin.POST("/teachers/:id/approve", handlers.ApproveTeacherProfile)
				admin.POST("/teachers/:id/reject", handlers.RejectTeacherProfile)
				admin.GET("/withdraw-requests", handlers.GetWithdrawRequests)
				admin.POST("/withdraw-requests/:id/process", handlers.ProcessWithdrawRequest)
				admin.POST("/subjects", handlers.CreateSubject)
				admin.PUT("/subjects/:id", handlers.UpdateSubject)
				admin.GET("/payment-configs", handlers.GetPaymentConfigs)
				admin.PUT("/payment-configs/:id", handlers.UpdatePaymentConfig)
				admin.POST("/reviews/:id/hide", handlers.HideReview)
			}
		}
	}
}
