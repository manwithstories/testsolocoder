package routes

import (
	"housekeeping-platform/config"
	"housekeeping-platform/handlers"
	"housekeeping-platform/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.RecoveryMiddleware())

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		authRequired := api.Group("")
		authRequired.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			authRequired.POST("/auth/logout", handlers.Logout)
			authRequired.GET("/auth/me", handlers.GetCurrentUser)
			authRequired.PUT("/auth/profile", handlers.UpdateProfile)
			authRequired.PUT("/auth/password", handlers.ChangePassword)

			addresses := authRequired.Group("/addresses")
			{
				addresses.GET("", handlers.GetAddressList)
				addresses.POST("", handlers.CreateAddress)
				addresses.PUT("/:id", handlers.UpdateAddress)
				addresses.DELETE("/:id", handlers.DeleteAddress)
				addresses.PUT("/:id/default", handlers.SetDefaultAddress)
			}

			certification := authRequired.Group("/certification")
			{
				certification.POST("", handlers.SubmitCertification)
			}

			services := authRequired.Group("/services")
			{
				services.GET("/categories", handlers.GetServiceCategories)
				services.GET("/areas", handlers.GetServiceAreas)
				services.GET("", handlers.GetServiceList)
				services.GET("/:id", handlers.GetServiceDetail)
				services.POST("", middleware.RequireRole("service_provider"), handlers.CreateServiceItem)
				services.PUT("/:id", middleware.RequireRole("service_provider"), handlers.UpdateServiceItem)
				services.DELETE("/:id", middleware.RequireRole("service_provider"), handlers.DeleteServiceItem)
				services.GET("/mine/list", middleware.RequireRole("service_provider"), handlers.GetMyServices)
			}

			orders := authRequired.Group("/orders")
			{
				orders.GET("", handlers.GetOrderList)
				orders.POST("", middleware.RequireRole("customer"), handlers.CreateOrder)
				orders.GET("/:id", handlers.GetOrderDetail)
				orders.PUT("/:id/cancel", handlers.CancelOrder)
				orders.PUT("/:id/start", middleware.RequireRole("service_provider"), handlers.StartService)
				orders.PUT("/:id/complete", middleware.RequireRole("service_provider"), handlers.CompleteService)
			}

			invitations := authRequired.Group("/invitations")
			{
				invitations.GET("", middleware.RequireRole("service_provider"), handlers.GetMyInvitations)
				invitations.PUT("/:id/respond", middleware.RequireRole("service_provider"), handlers.RespondInvitation)
			}

			reviews := authRequired.Group("/reviews")
			{
				reviews.GET("", handlers.GetReviewList)
				reviews.GET("/:id", handlers.GetReviewDetail)
				reviews.POST("", handlers.CreateReview)
				reviews.PUT("/:id/reply", middleware.RequireRole("service_provider"), handlers.ReplyReview)
			}

			complaints := authRequired.Group("/complaints")
			{
				complaints.GET("", handlers.GetComplaintList)
				complaints.POST("", handlers.CreateComplaint)
				complaints.PUT("/:id/handle", middleware.RequireRole("admin"), handlers.HandleComplaint)
			}

			bills := authRequired.Group("/bills")
			{
				bills.GET("", handlers.GetBillList)
				bills.GET("/balance", middleware.RequireRole("service_provider"), handlers.GetBalance)
				bills.GET("/income-summary", middleware.RequireRole("service_provider"), handlers.GetIncomeSummary)
			}

			withdraws := authRequired.Group("/withdraws")
			{
				withdraws.GET("", handlers.GetWithdrawList)
				withdraws.POST("", middleware.RequireRole("service_provider"), handlers.CreateWithdrawRequest)
				withdraws.PUT("/:id/handle", middleware.RequireRole("admin"), handlers.HandleWithdrawRequest)
			}

			messages := authRequired.Group("/messages")
			{
				messages.GET("", handlers.GetMessageList)
				messages.GET("/unread-count", handlers.GetUnreadCount)
				messages.PUT("/:id/read", handlers.ReadMessage)
				messages.PUT("/read-all", handlers.ReadAllMessages)
				messages.DELETE("/:id", handlers.DeleteMessage)
			}

			admin := authRequired.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/dashboard", handlers.GetDashboardStats)
				admin.GET("/statistics/orders", handlers.GetOrderStatistics)
				admin.GET("/statistics/users", handlers.GetUserStatistics)
				admin.GET("/statistics/export", handlers.ExportStatistics)

				admin.GET("/users", handlers.GetUserList)
				admin.GET("/users/:id", handlers.GetUserDetail)
				admin.PUT("/users/:id/toggle-status", handlers.ToggleUserStatus)
				admin.PUT("/users/:id/review-certification", handlers.ReviewCertification)

				admin.POST("/service-categories", handlers.CreateServiceCategory)
				admin.PUT("/service-categories/:id", handlers.UpdateServiceCategory)
				admin.DELETE("/service-categories/:id", handlers.DeleteServiceCategory)
				admin.POST("/service-areas", handlers.CreateServiceArea)
			}
		}
	}
}
