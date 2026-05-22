package routes

import (
	"github.com/gin-gonic/gin"

	"travel-planner/internal/handlers"
	"travel-planner/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.GET("/me", middleware.JWTAuth(), handlers.GetCurrentUser)
			auth.PUT("/profile", middleware.JWTAuth(), handlers.UpdateProfile)
		}

		plans := api.Group("/plans")
		plans.Use(middleware.JWTAuth())
		{
			plans.GET("", handlers.GetPlans)
			plans.POST("", handlers.CreatePlan)
			plans.GET("/:id", middleware.CheckPlanAccess(false), handlers.GetPlan)
			plans.PUT("/:id", middleware.CheckPlanAccess(true), handlers.UpdatePlan)
			plans.DELETE("/:id", middleware.CheckPlanAccess(true), handlers.DeletePlan)

			plans.POST("/:id/participants", middleware.CheckPlanAccess(true), handlers.AddParticipant)
			plans.DELETE("/:id/participants/:participant_id", middleware.CheckPlanAccess(true), handlers.RemoveParticipant)

			activities := plans.Group("/:plan_id/activities")
			activities.Use(middleware.CheckPlanAccess(false))
			{
				activities.GET("", handlers.GetActivities)
				activities.POST("", middleware.CheckPlanAccess(true), handlers.CreateActivity)
				activities.GET("/:id", handlers.GetActivity)
				activities.PUT("/:id", middleware.CheckPlanAccess(true), handlers.UpdateActivity)
				activities.DELETE("/:id", middleware.CheckPlanAccess(true), handlers.DeleteActivity)
			}

			budget := plans.Group("/:plan_id/budget")
			budget.Use(middleware.CheckPlanAccess(false))
			{
				budget.GET("", handlers.GetBudgetSummary)
			}

			files := plans.Group("/:plan_id/files")
			files.Use(middleware.CheckPlanAccess(false))
			{
				files.GET("", handlers.GetFiles)
				files.POST("", middleware.CheckPlanAccess(true), handlers.UploadFile)
				files.DELETE("/:id", middleware.CheckPlanAccess(true), handlers.DeleteFile)
			}

			checklists := plans.Group("/:plan_id/checklists")
			checklists.Use(middleware.CheckPlanAccess(false))
			{
				checklists.GET("", handlers.GetChecklists)
				checklists.POST("", middleware.CheckPlanAccess(true), handlers.CreateChecklist)
				checklists.GET("/:id", handlers.GetChecklist)
				checklists.PUT("/:id", middleware.CheckPlanAccess(true), handlers.UpdateChecklist)
				checklists.DELETE("/:id", middleware.CheckPlanAccess(true), handlers.DeleteChecklist)

				items := checklists.Group("/:id/items")
				{
					items.POST("", middleware.CheckPlanAccess(true), handlers.AddChecklistItem)
					items.PUT("/:item_id", middleware.CheckPlanAccess(true), handlers.UpdateChecklistItem)
					items.DELETE("/:item_id", middleware.CheckPlanAccess(true), handlers.DeleteChecklistItem)
				}
			}

			reminders := plans.Group("/:plan_id/reminders")
			reminders.Use(middleware.CheckPlanAccess(false))
			{
				reminders.POST("", middleware.CheckPlanAccess(true), handlers.CreateReminder)
			}

			api.GET("/reminders", middleware.JWTAuth(), handlers.GetReminders)
			api.PUT("/reminders/:id", middleware.JWTAuth(), handlers.UpdateReminder)
			api.DELETE("/reminders/:id", middleware.JWTAuth(), handlers.DeleteReminder)

			maps := plans.Group("/:plan_id/map")
			maps.Use(middleware.CheckPlanAccess(false))
			{
				maps.GET("", handlers.GetMapData)
			}

			export := plans.Group("/:id/export")
			export.Use(middleware.CheckPlanAccess(false))
			{
				export.GET("/json", handlers.ExportJSON)
				export.GET("/pdf", handlers.ExportPDF)
			}
		}

		files := api.Group("/files")
		{
			files.GET("/:filename", handlers.GetFile)
			files.GET("/download/:id", middleware.JWTAuth(), handlers.DownloadFile)
		}
	}
}
