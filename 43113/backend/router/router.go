package router

import (
	"qa-platform/handlers"
	"qa-platform/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	userHandler := handlers.NewUserHandler()
	questionHandler := handlers.NewQuestionHandler()
	answerHandler := handlers.NewAnswerHandler()
	commentHandler := handlers.NewCommentHandler()
	auditHandler := handlers.NewAuditHandler()
	favoriteHandler := handlers.NewFavoriteHandler()
	followHandler := handlers.NewFollowHandler()
	notificationHandler := handlers.NewNotificationHandler()
	rewardHandler := handlers.NewRewardHandler()
	searchHandler := handlers.NewSearchHandler()
	statsHandler := handlers.NewStatsHandler()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		public := api.Group("/public")
		{
			public.GET("/questions", questionHandler.GetQuestionList)
			public.GET("/questions/:id", questionHandler.GetQuestion)
			public.GET("/categories", questionHandler.GetCategories)
			public.GET("/tags", questionHandler.GetTags)
			public.GET("/users/:id", userHandler.GetUserByID)
			public.GET("/search/questions", searchHandler.SearchQuestions)
		}

		authRequired := api.Group("")
		authRequired.Use(middleware.AuthMiddleware())
		{
			user := authRequired.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.GET("/points/logs", userHandler.GetPointLogs)
				user.POST("/expert/apply", userHandler.ApplyExpert)
			}

			questions := authRequired.Group("/questions")
			{
				questions.POST("", questionHandler.CreateQuestion)
				questions.PUT("/:id", questionHandler.UpdateQuestion)
				questions.DELETE("/:id", questionHandler.DeleteQuestion)
				questions.POST("/:id/like", questionHandler.LikeQuestion)
				questions.POST("/:id/accept", questionHandler.AcceptAnswer)
			}

			answers := authRequired.Group("/answers")
			{
				answers.POST("", answerHandler.CreateAnswer)
				answers.GET("/:id", answerHandler.GetAnswer)
				answers.PUT("/:id", answerHandler.UpdateAnswer)
				answers.DELETE("/:id", answerHandler.DeleteAnswer)
				answers.POST("/:id/like", answerHandler.LikeAnswer)
				answers.POST("/:id/dislike", answerHandler.DislikeAnswer)
			}

			comments := authRequired.Group("/comments")
			{
				comments.POST("", commentHandler.CreateComment)
				comments.DELETE("/:id", commentHandler.DeleteComment)
				comments.POST("/:id/like", commentHandler.LikeComment)
			}

			favorites := authRequired.Group("/favorites")
			{
				favorites.POST("", favoriteHandler.AddFavorite)
				favorites.DELETE("", favoriteHandler.RemoveFavorite)
				favorites.GET("", favoriteHandler.GetUserFavorites)
			}

			follows := authRequired.Group("/follows")
			{
				follows.POST("", followHandler.Follow)
				follows.DELETE("", followHandler.Unfollow)
				follows.GET("", followHandler.GetUserFollows)
				follows.GET("/followers", followHandler.GetUserFollowers)
				follows.GET("/check", followHandler.IsFollowing)
			}

			notifications := authRequired.Group("/notifications")
			{
				notifications.GET("", notificationHandler.GetNotifications)
				notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
				notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
				notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			}

			rewards := authRequired.Group("/rewards")
			{
				rewards.GET("", rewardHandler.GetRewardList)
				rewards.POST("/:id/exchange", rewardHandler.ExchangeReward)
				rewards.GET("/exchanges", rewardHandler.GetExchangeList)
				rewards.GET("/points/logs", rewardHandler.GetPointLogs)
			}

			reports := authRequired.Group("/reports")
			{
				reports.POST("", auditHandler.CreateReport)
			}

			search := authRequired.Group("/search")
			{
				search.GET("/recommendations", searchHandler.GetRecommendations)
			}

			admin := authRequired.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				adminUsers := admin.Group("/users")
				{
					adminUsers.GET("", userHandler.GetUserList)
					adminUsers.PUT("/:id/status", userHandler.UpdateUserStatus)
				}

				adminCategories := admin.Group("/categories")
				{
					adminCategories.POST("", questionHandler.CreateCategory)
				}

				adminTags := admin.Group("/tags")
				{
					adminTags.POST("", questionHandler.CreateTag)
				}

				adminAudit := admin.Group("/audit")
				{
					adminAudit.GET("", auditHandler.GetAuditList)
					adminAudit.POST("", auditHandler.AuditContent)
					adminAudit.GET("/pending-count", auditHandler.GetPendingAuditCount)
				}

				adminReports := admin.Group("/reports")
				{
					adminReports.GET("", auditHandler.GetReportList)
					adminReports.PUT("/:id", auditHandler.HandleReport)
				}

				adminSensitiveWords := admin.Group("/sensitive-words")
				{
					adminSensitiveWords.GET("", auditHandler.GetSensitiveWords)
					adminSensitiveWords.POST("", auditHandler.CreateSensitiveWord)
					adminSensitiveWords.DELETE("/:id", auditHandler.DeleteSensitiveWord)
				}

				adminRewards := admin.Group("/rewards")
				{
					adminRewards.POST("", rewardHandler.CreateReward)
					adminRewards.PUT("/:id", rewardHandler.UpdateReward)
					adminRewards.DELETE("/:id", rewardHandler.DeleteReward)
				}

				adminExpert := admin.Group("/expert")
				{
					adminExpert.GET("/applications", userHandler.GetExpertApplications)
					adminExpert.PUT("/applications/:id", userHandler.ReviewExpertApplication)
				}

				adminStats := admin.Group("/stats")
				{
					adminStats.GET("/dashboard", statsHandler.GetDashboardStats)
					adminStats.GET("/activity", statsHandler.GetActivityReport)
					adminStats.GET("/audit", statsHandler.GetAuditReport)
					adminStats.GET("/activity/export", statsHandler.ExportActivityReport)
					adminStats.GET("/audit/export", statsHandler.ExportAuditReport)
				}
			}

			expert := authRequired.Group("/expert")
			expert.Use(middleware.ExpertMiddleware())
			{
				expert.GET("/privileges", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"code":    0,
						"message": "success",
						"data": gin.H{
							"canReview":       true,
							"priorityReview":  true,
							"specialBadge":    true,
							"rewardMultiplier": 1.5,
						},
					})
				})
			}
		}
	}

	return r
}
