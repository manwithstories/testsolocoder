package router

import (
	"github.com/gin-gonic/gin"

	"tea-platform/internal/handlers"
	"tea-platform/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	userHandler := handlers.NewUserHandler()
	teaHandler := handlers.NewTeaHandler()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		teas := api.Group("/teas")
		{
			teas.GET("", teaHandler.GetTeaList)
			teas.GET("/:id", teaHandler.GetTeaDetail)
		}

		authorized := api.Group("/")
		authorized.Use(middleware.JWTAuth())
		{
			user := authorized.Group("/user")
			{
				user.GET("/info", userHandler.GetUserInfo)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.PUT("/password", userHandler.ChangePassword)
				user.POST("/avatar", userHandler.UploadAvatar)
			}

			admin := authorized.Group("/admin")
			admin.Use(middleware.RoleAuth("admin"))
			{
				admin.GET("/users", userHandler.GetUserList)
				admin.PUT("/users/:id/status", userHandler.UpdateUserStatus)
			}

			sellerOrAdmin := authorized.Group("/")
			sellerOrAdmin.Use(middleware.RoleAuth("seller", "admin"))
			{
				teasAuth := sellerOrAdmin.Group("/teas")
				{
					teasAuth.POST("", teaHandler.CreateTea)
					teasAuth.PUT("/:id", teaHandler.UpdateTea)
					teasAuth.DELETE("/:id", teaHandler.DeleteTea)
					teasAuth.POST("/:id/images", teaHandler.UploadTeaImage)
					teasAuth.DELETE("/images/:image_id", teaHandler.DeleteTeaImage)
					teasAuth.POST("/:id/traceability", teaHandler.AddTraceability)
					teasAuth.GET("/:id/traceability", teaHandler.GetTraceability)
				}
			}
		}
	}
}
