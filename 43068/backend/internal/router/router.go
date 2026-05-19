package router

import (
	"freelancer-management/internal/handlers"
	"freelancer-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AccessLogMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	authHandler := handlers.NewAuthHandler()
	clientHandler := handlers.NewClientHandler()
	projectHandler := handlers.NewProjectHandler()
	timeEntryHandler := handlers.NewTimeEntryHandler()
	invoiceHandler := handlers.NewInvoiceHandler()
	dashboardHandler := handlers.NewDashboardHandler()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.Me)
		}

		clients := api.Group("/clients")
		clients.Use(middleware.AuthMiddleware())
		{
			clients.POST("", clientHandler.Create)
			clients.GET("", clientHandler.List)
			clients.GET("/:id", clientHandler.Get)
			clients.PUT("/:id", clientHandler.Update)
			clients.DELETE("/:id", clientHandler.Delete)
		}

		projects := api.Group("/projects")
		projects.Use(middleware.AuthMiddleware())
		{
			projects.POST("", projectHandler.Create)
			projects.GET("", projectHandler.List)
			projects.GET("/:id", projectHandler.Get)
			projects.PUT("/:id", projectHandler.Update)
			projects.DELETE("/:id", projectHandler.Delete)
			projects.POST("/:id/milestones", projectHandler.AddMilestone)
			projects.PUT("/milestones/:milestone_id", projectHandler.UpdateMilestone)
			projects.DELETE("/milestones/:milestone_id", projectHandler.DeleteMilestone)
		}

		timeEntries := api.Group("/time-entries")
		timeEntries.Use(middleware.AuthMiddleware())
		{
			timeEntries.POST("", timeEntryHandler.Create)
			timeEntries.GET("", timeEntryHandler.List)
			timeEntries.GET("/active-timer", timeEntryHandler.GetActiveTimer)
			timeEntries.GET("/:id", timeEntryHandler.Get)
			timeEntries.PUT("/:id", timeEntryHandler.Update)
			timeEntries.DELETE("/:id", timeEntryHandler.Delete)
			timeEntries.POST("/timer/start", timeEntryHandler.StartTimer)
			timeEntries.POST("/timer/:id/stop", timeEntryHandler.StopTimer)
		}

		invoices := api.Group("/invoices")
		invoices.Use(middleware.AuthMiddleware())
		{
			invoices.POST("", invoiceHandler.Create)
			invoices.GET("", invoiceHandler.List)
			invoices.GET("/:id", invoiceHandler.Get)
			invoices.PUT("/:id/status", invoiceHandler.UpdateStatus)
			invoices.GET("/:id/download", invoiceHandler.DownloadPDF)
			invoices.DELETE("/:id", invoiceHandler.Delete)
		}

		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.AuthMiddleware())
		{
			dashboard.GET("", dashboardHandler.GetStats)
		}
	}

	return r
}
