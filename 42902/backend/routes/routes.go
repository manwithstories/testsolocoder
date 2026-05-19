package routes

import (
	"activity-management/controllers"
	"activity-management/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/verify", controllers.VerifyEmail)
		}

		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/profile", controllers.GetProfile)
			user.GET("/registrations", controllers.GetMyRegistrations)
		}

		events := api.Group("/events")
		{
			events.GET("", controllers.GetEvents)
			events.GET("/:id", controllers.GetEvent)
			events.GET("/:id/registrations", middleware.AuthMiddleware(), controllers.GetEventRegistrations)
			events.GET("/:id/export", middleware.AuthMiddleware(), controllers.ExportRegistrations)

			events.POST("", middleware.AuthMiddleware(), controllers.CreateEvent)
			events.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateEvent)
			events.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteEvent)

			events.POST("/:id/register", middleware.AuthMiddleware(), controllers.RegisterEvent)
			events.POST("/:id/cancel", middleware.AuthMiddleware(), controllers.CancelRegistration)
		}
	}
}
