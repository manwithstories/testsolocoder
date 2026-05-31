package routes

import (
	"ship-rental-platform/internal/handler"
	"ship-rental-platform/internal/middleware"
	"ship-rental-platform/internal/model"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	userHandler := handler.NewUserHandler()
	shipHandler := handler.NewShipHandler()
	berthHandler := handler.NewBerthHandler()
	rentalHandler := handler.NewRentalHandler()
	voyageHandler := handler.NewVoyageHandler()
	maintenanceHandler := handler.NewMaintenanceHandler()
	financeHandler := handler.NewFinanceHandler()
	reviewHandler := handler.NewReviewHandler()

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)

		api.GET("/ships", shipHandler.GetShips)
		api.GET("/ships/:id", shipHandler.GetShip)

		api.GET("/docks", berthHandler.GetDocks)
		api.GET("/docks/:id", berthHandler.GetDock)
		api.GET("/berths", berthHandler.GetBerths)
		api.GET("/berths/:id", berthHandler.GetBerth)
		api.GET("/berths/availability", berthHandler.CheckAvailability)
		api.GET("/docks/:id/water-levels", berthHandler.GetWaterLevels)

		api.GET("/reviews", reviewHandler.GetReviews)
	}

	auth := api.Group("")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/auth/profile", userHandler.GetProfile)
		auth.PUT("/auth/profile", userHandler.UpdateProfile)
		auth.PUT("/auth/password", userHandler.ChangePassword)
		auth.GET("/auth/me", userHandler.GetProfile)

		auth.GET("/my-ships", shipHandler.GetMyShips)
		auth.POST("/ships", shipHandler.CreateShip)
		auth.PUT("/ships/:id", shipHandler.UpdateShip)
		auth.DELETE("/ships/:id", shipHandler.DeleteShip)
		auth.POST("/ships/:id/images", shipHandler.UploadImage)
		auth.DELETE("/ships/:id/images/:imageId", shipHandler.DeleteImage)

		auth.GET("/berth-reservations", berthHandler.GetReservations)
		auth.POST("/berth-reservations", berthHandler.CreateReservation)
		auth.PUT("/berth-reservations/:id/cancel", berthHandler.CancelReservation)

		auth.GET("/rentals", rentalHandler.GetRentals)
		auth.GET("/rentals/:id", rentalHandler.GetRental)
		auth.POST("/rentals", rentalHandler.CreateRental)
		auth.PUT("/rentals/:id/status", rentalHandler.UpdateRentalStatus)
		auth.DELETE("/rentals/:id", rentalHandler.CancelRental)
		auth.GET("/my-rentals", rentalHandler.GetMyRentals)

		auth.GET("/voyage-logs", voyageHandler.GetVoyageLogs)
		auth.GET("/voyage-logs/:id", voyageHandler.GetVoyageLog)
		auth.POST("/voyage-logs", voyageHandler.CreateVoyageLog)
		auth.PUT("/voyage-logs/:id", voyageHandler.UpdateVoyageLog)
		auth.DELETE("/voyage-logs/:id", voyageHandler.DeleteVoyageLog)
		auth.GET("/voyage-logs/export", voyageHandler.ExportVoyageLogs)

		auth.GET("/maintenance", maintenanceHandler.GetMaintenances)
		auth.GET("/maintenance/:id", maintenanceHandler.GetMaintenance)
		auth.POST("/maintenance", maintenanceHandler.CreateMaintenance)
		auth.PUT("/maintenance/:id", maintenanceHandler.UpdateMaintenance)
		auth.DELETE("/maintenance/:id", maintenanceHandler.DeleteMaintenance)

		auth.GET("/maintenance-schedules", maintenanceHandler.GetSchedules)
		auth.POST("/maintenance-schedules", maintenanceHandler.CreateSchedule)
		auth.PUT("/maintenance-schedules/:id", maintenanceHandler.UpdateSchedule)
		auth.DELETE("/maintenance-schedules/:id", maintenanceHandler.DeleteSchedule)

		auth.GET("/transactions", financeHandler.GetTransactions)
		auth.GET("/transactions/:id", financeHandler.GetTransaction)
		auth.POST("/transactions", financeHandler.CreateTransaction)
		auth.PUT("/transactions/:id/status", financeHandler.UpdateTransactionStatus)

		auth.GET("/settlements", financeHandler.GetSettlements)
		auth.POST("/settlements", financeHandler.CreateSettlement)

		auth.GET("/financial-summary", financeHandler.GetFinancialSummary)
		auth.GET("/financial-report/export", financeHandler.ExportFinancialReport)
		auth.GET("/monthly-report/export", financeHandler.ExportMonthlyReport)

		auth.GET("/my-reviews", reviewHandler.GetMyReviews)
		auth.POST("/reviews", reviewHandler.CreateReview)
		auth.POST("/reviews/:id/respond", reviewHandler.RespondToReview)
		auth.POST("/reviews/:id/helpful", reviewHandler.MarkHelpful)
		auth.DELETE("/reviews/:id", reviewHandler.DeleteReview)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.RoleAuth(model.RoleAdmin))
	{
		admin.GET("/users", userHandler.GetUsers)
		admin.GET("/users/:id", userHandler.GetUserByID)
		admin.PUT("/users/:id/toggle-status", userHandler.ToggleUserStatus)

		admin.POST("/docks", berthHandler.CreateDock)
		admin.PUT("/docks/:id", berthHandler.UpdateDock)
		admin.DELETE("/docks/:id", berthHandler.DeleteDock)

		admin.POST("/berths", berthHandler.CreateBerth)
		admin.PUT("/berths/:id", berthHandler.UpdateBerth)
		admin.DELETE("/berths/:id", berthHandler.DeleteBerth)

		admin.POST("/docks/:id/water-levels", berthHandler.RecordWaterLevel)
	}
}
