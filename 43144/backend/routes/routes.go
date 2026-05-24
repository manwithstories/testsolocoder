package routes

import (
	"pet-adoption-platform/handlers"
	"pet-adoption-platform/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/health", handlers.GetHealthCheck)

		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		api.GET("/uploads/pets/:filename", handlers.GetUploadedFile)

		api.Static("/uploads", "./uploads")

		api.Use(middleware.AuthMiddleware())
		{
			profile := api.Group("/profile")
			{
				profile.GET("", handlers.GetProfile)
				profile.PUT("", handlers.UpdateProfile)
			}

			pets := api.Group("/pets")
			{
				pets.GET("", handlers.ListPets)
				pets.GET("/my", handlers.GetMyPets)
				pets.GET("/adopted", handlers.GetMyAdoptedPets)
				pets.GET("/:id", handlers.GetPet)
				pets.GET("/:id/history", handlers.GetPetAdoptionHistory)
			}

			petsAdmin := pets.Group("")
			petsAdmin.Use(middleware.RoleMiddleware("rescue", "admin"))
			{
				petsAdmin.POST("", handlers.CreatePet)
				petsAdmin.PUT("/:id", handlers.UpdatePet)
				petsAdmin.DELETE("/:id", handlers.DeletePet)
				petsAdmin.PUT("/:id/status", handlers.UpdatePetStatus)
				petsAdmin.POST("/:id/photos", handlers.UploadPetPhotos)
				petsAdmin.POST("/:id/videos", handlers.UploadPetVideos)
			}

			adoption := api.Group("/adoption")
			{
				adoption.GET("/applications", handlers.ListAdoptionApplications)
				adoption.GET("/applications/:id", handlers.GetAdoptionApplication)
				adoption.GET("/applications/:id/agreement", handlers.GetAdoptionAgreement)
				adoption.POST("/applications", handlers.CreateAdoptionApplication)
				adoption.PUT("/applications/:id/sign", handlers.SignAdoptionAgreement)
				adoption.PUT("/applications/:id/complete", handlers.CompleteAdoption)

				adoption.POST("/follow-ups", handlers.CreateFollowUpRecord)
				adoption.GET("/pets/:pet_id/follow-ups", handlers.ListFollowUpRecords)
			}

			adoptionRescue := adoption.Group("")
			adoptionRescue.Use(middleware.RoleMiddleware("rescue", "admin"))
			{
				adoptionRescue.PUT("/applications/:id/review", handlers.ReviewAdoptionApplication)
			}

			health := api.Group("/health")
			{
				health.GET("/records", handlers.ListHealthRecords)
				health.GET("/records/:id", handlers.GetHealthRecord)
				health.GET("/pets/:pet_id/summary", handlers.GetPetHealthSummary)
				health.GET("/pets/:pet_id/reminders", handlers.GetHealthReminders)
				health.PUT("/reminders/:id/complete", handlers.CompleteHealthReminder)
			}

			healthRescue := health.Group("")
			healthRescue.Use(middleware.RoleMiddleware("rescue", "admin"))
			{
				healthRescue.POST("/records", handlers.CreateHealthRecord)
				healthRescue.PUT("/records/:id", handlers.UpdateHealthRecord)
				healthRescue.DELETE("/records/:id", handlers.DeleteHealthRecord)
				healthRescue.POST("/records/:id/report", handlers.UploadHealthReport)
			}

			appointments := api.Group("/appointments")
			{
				appointments.GET("", handlers.ListAppointments)
				appointments.GET("/:id", handlers.GetAppointment)
				appointments.POST("", handlers.CreateAppointment)
				appointments.PUT("/:id", handlers.UpdateAppointment)
				appointments.PUT("/:id/cancel", handlers.CancelAppointment)
				appointments.PUT("/:id/reschedule", handlers.RescheduleAppointment)
			}

			appointmentsRescue := appointments.Group("")
			appointmentsRescue.Use(middleware.RoleMiddleware("rescue", "admin"))
			{
				appointmentsRescue.PUT("/:id/confirm", handlers.ConfirmAppointment)
				appointmentsRescue.PUT("/:id/complete", handlers.CompleteAppointment)
			}

			rescue := api.Group("/rescue")
			{
				rescue.GET("", handlers.ListRescueStations)
				rescue.GET("/:id", handlers.GetRescueStation)
				rescue.GET("/:id/stats", handlers.GetRescueStatsByID)
			}

			rescueAdmin := rescue.Group("")
			rescueAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				rescueAdmin.PUT("/:id/review", handlers.ReviewRescueStation)
				rescueAdmin.GET("/stats/all", handlers.GetAllRescuesStats)
			}

			rescueSelf := rescue.Group("/me")
			rescueSelf.Use(middleware.RoleMiddleware("rescue"))
			{
				rescueSelf.GET("/stats", handlers.GetRescueStats)
			}

			users := api.Group("/users")
			users.Use(middleware.RoleMiddleware("admin"))
			{
				users.GET("", handlers.ListUsers)
				users.GET("/:id", handlers.GetUserByID)
				users.PUT("/:id/verify", handlers.VerifyUser)
			}

			export := api.Group("/export")
			export.Use(middleware.RoleMiddleware("rescue", "admin"))
			{
				export.GET("/adoption", handlers.ExportAdoptionReport)
				export.GET("/health/:pet_id", handlers.ExportHealthReport)
			}
		}
	}
}
