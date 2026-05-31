package routes

import (
	"consultation-platform/config"
	"consultation-platform/controllers"
	"consultation-platform/middlewares"
	"consultation-platform/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	userCtrl := controllers.NewUserController(cfg)
	serviceCtrl := controllers.NewServiceController(cfg)
	appointmentCtrl := controllers.NewAppointmentController(cfg)
	recordCtrl := controllers.NewRecordController(cfg.SensitiveWords)
	notificationCtrl := controllers.NewNotificationController(cfg)
	statisticsCtrl := controllers.NewStatisticsController()

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", userCtrl.Register)
		auth.POST("/login", userCtrl.Login)
		auth.POST("/refresh", middlewares.JWTAuthMiddleware(cfg.JWT.Secret), userCtrl.RefreshToken)
	}

	users := api.Group("/users")
	users.Use(middlewares.JWTAuthMiddleware(cfg.JWT.Secret))
	{
		users.GET("/profile", userCtrl.GetProfile)
		users.PUT("/profile", userCtrl.UpdateProfile)
		users.GET("/:id", userCtrl.GetUserByID)

		admin := users.Group("")
		admin.Use(middlewares.RoleMiddleware(string(models.RoleAdmin)))
		{
			admin.GET("", userCtrl.GetUsers)
			admin.GET("/verifications/pending", userCtrl.GetPendingVerifications)
			admin.PUT("/:id/verify", userCtrl.VerifyProfessional)
		}
	}

	services := api.Group("/services")
	{
		services.GET("", serviceCtrl.GetAllServices)
		services.GET("/:id", serviceCtrl.GetServiceByID)
		services.GET("/:id/schedules", serviceCtrl.GetSchedules)
		services.GET("/:service_id/reviews", recordCtrl.GetServiceReviews)

		authServices := services.Group("")
		authServices.Use(middlewares.JWTAuthMiddleware(cfg.JWT.Secret))
		{
			professional := authServices.Group("")
			professional.Use(middlewares.RoleMiddleware(string(models.RoleProfessional)))
			{
				professional.POST("", serviceCtrl.CreateService)
				professional.PUT("/:id", serviceCtrl.UpdateService)
				professional.DELETE("/:id", serviceCtrl.DeleteService)
				professional.GET("/professional/list", serviceCtrl.GetProfessionalServices)
				professional.POST("/schedules", serviceCtrl.CreateSchedule)
				professional.POST("/schedules/batch", serviceCtrl.BatchCreateSchedules)
				professional.DELETE("/:service_id/schedules", serviceCtrl.DeleteSchedules)
			}
		}
	}

	appointments := api.Group("/appointments")
	appointments.Use(middlewares.JWTAuthMiddleware(cfg.JWT.Secret))
	{
		appointments.POST("", appointmentCtrl.CreateAppointment)
		appointments.GET("/:id", appointmentCtrl.GetAppointmentByID)

		client := appointments.Group("")
		client.Use(middlewares.RoleMiddleware(string(models.RoleClient)))
		{
			client.GET("/client/list", appointmentCtrl.GetClientAppointments)
			client.PUT("/:id/cancel", appointmentCtrl.CancelAppointment)
		}

		professional := appointments.Group("")
		professional.Use(middlewares.RoleMiddleware(string(models.RoleProfessional)))
		{
			professional.GET("/professional/list", appointmentCtrl.GetProfessionalAppointments)
			professional.PUT("/:id/confirm", appointmentCtrl.ConfirmAppointment)
			professional.PUT("/:id/cancel", appointmentCtrl.CancelAppointment)
			professional.PUT("/:id/complete", appointmentCtrl.CompleteAppointment)
		}

		payment := appointments.Group("")
		{
			payment.PUT("/:id/pay", appointmentCtrl.ProcessPayment)
			payment.PUT("/:id/refund", appointmentCtrl.RefundPayment)
		}
	}

	records := api.Group("/records")
	records.Use(middlewares.JWTAuthMiddleware(cfg.JWT.Secret))
	{
		professional := records.Group("")
		professional.Use(middlewares.RoleMiddleware(string(models.RoleProfessional)))
		{
			professional.POST("/consult", recordCtrl.CreateConsultRecord)
			professional.GET("/professional/list", recordCtrl.GetProfessionalConsultRecords)
		}

		client := records.Group("")
		client.Use(middlewares.RoleMiddleware(string(models.RoleClient)))
		{
			client.GET("/client/list", recordCtrl.GetClientConsultRecords)
			client.POST("/review", recordCtrl.CreateReview)
		}

		records.GET("/:id", recordCtrl.GetConsultRecordByID)
	}

	reviews := api.Group("/reviews")
	reviews.Use(middlewares.JWTAuthMiddleware(cfg.JWT.Secret))
	{
		professional := reviews.Group("")
		professional.Use(middlewares.RoleMiddleware(string(models.RoleProfessional)))
		{
			professional.GET("/professional/list", recordCtrl.GetProfessionalReviews)
			professional.GET("/professional/stats/:professional_id", recordCtrl.GetProfessionalReviewStats)
		}

		admin := reviews.Group("")
		admin.Use(middlewares.RoleMiddleware(string(models.RoleAdmin)))
		{
			admin.GET("/pending", recordCtrl.GetPendingReviews)
			admin.PUT("/status", recordCtrl.UpdateReviewStatus)
		}
	}

	notifications := api.Group("/notifications")
	notifications.Use(middlewares.JWTAuthMiddleware(cfg.JWT.Secret))
	{
		notifications.GET("", notificationCtrl.GetNotifications)
		notifications.PUT("/:id/read", notificationCtrl.MarkAsRead)
		notifications.PUT("/read-all", notificationCtrl.MarkAllAsRead)
		notifications.GET("/unread-count", notificationCtrl.GetUnreadCount)

		admin := notifications.Group("")
		admin.Use(middlewares.RoleMiddleware(string(models.RoleAdmin)))
		{
			admin.GET("/templates", notificationCtrl.GetTemplates)
			admin.PUT("/templates/:id", notificationCtrl.UpdateTemplate)
		}
	}

	statistics := api.Group("/statistics")
	statistics.Use(middlewares.JWTAuthMiddleware(cfg.JWT.Secret))
	{
		professional := statistics.Group("")
		professional.Use(middlewares.RoleMiddleware(string(models.RoleProfessional)))
		{
			professional.GET("/professional", statisticsCtrl.GetProfessionalStats)
			professional.GET("/professional/export/appointments", statisticsCtrl.ExportAppointments)
			professional.GET("/professional/export/revenue", statisticsCtrl.ExportRevenue)
		}

		admin := statistics.Group("")
		admin.Use(middlewares.RoleMiddleware(string(models.RoleAdmin)))
		{
			admin.GET("/admin", statisticsCtrl.GetAdminStats)
			admin.GET("/admin/export/appointments", statisticsCtrl.ExportAppointments)
			admin.GET("/admin/export/revenue", statisticsCtrl.ExportRevenue)
		}
	}
}
