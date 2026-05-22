package main

import (
	"wedding-planner/internal/handlers"
	"wedding-planner/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	userHandler := handlers.NewUserHandler()
	weddingHandler := handlers.NewWeddingHandler()
	vendorHandler := handlers.NewVendorHandler()
	guestHandler := handlers.NewGuestHandler()
	budgetHandler := handlers.NewBudgetHandler()
	taskHandler := handlers.NewTaskHandler()
	documentHandler := handlers.NewDocumentHandler()
	dashboardHandler := handlers.NewDashboardHandler()
	logHandler := handlers.NewLogHandler()
	notificationHandler := handlers.NewNotificationHandler()

	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		users := api.Group("/users")
		users.Use(middleware.JWTAuth())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.PUT("/password", userHandler.ChangePassword)
			users.GET("", middleware.RequireRole("admin"), middleware.Pagination(), userHandler.GetUsers)
			users.PUT("/:id/status", middleware.RequireRole("admin"), middleware.ParseUintID("id"), userHandler.UpdateUserStatus)
			users.DELETE("/:id", middleware.RequireRole("admin"), middleware.ParseUintID("id"), userHandler.DeleteUser)
		}

		weddings := api.Group("/weddings")
		weddings.Use(middleware.JWTAuth(), middleware.Pagination())
		{
			weddings.POST("", weddingHandler.Create)
			weddings.GET("", weddingHandler.GetList)
			weddings.GET("/:id", middleware.ParseUintID("id"), weddingHandler.GetByID)
			weddings.PUT("/:id", middleware.ParseUintID("id"), weddingHandler.Update)
			weddings.DELETE("/:id", middleware.ParseUintID("id"), weddingHandler.Delete)
			weddings.PUT("/:id/status", middleware.ParseUintID("id"), weddingHandler.UpdateStatus)
		}

		vendors := api.Group("/vendors")
		vendors.Use(middleware.JWTAuth(), middleware.Pagination())
		{
			vendors.POST("", vendorHandler.Create)
			vendors.GET("", vendorHandler.GetList)
			vendors.GET("/categories", vendorHandler.GetCategories)
			vendors.GET("/:id", middleware.ParseUintID("id"), vendorHandler.GetByID)
			vendors.PUT("/:id", middleware.ParseUintID("id"), vendorHandler.Update)
			vendors.DELETE("/:id", middleware.ParseUintID("id"), vendorHandler.Delete)
			vendors.POST("/:id/reviews", middleware.ParseUintID("id"), vendorHandler.AddReview)
		}

		guests := api.Group("/weddings/:wedding_id/guests")
		guests.Use(middleware.JWTAuth(), middleware.ParseUintID("wedding_id"), middleware.Pagination())
		{
			guests.POST("", guestHandler.Create)
			guests.GET("", guestHandler.GetList)
			guests.GET("/export", guestHandler.ExportGuests)
			guests.POST("/import", guestHandler.ImportGuests)
			guests.GET("/:id", middleware.ParseUintID("id"), guestHandler.GetByID)
			guests.PUT("/:id", middleware.ParseUintID("id"), guestHandler.Update)
			guests.DELETE("/:id", middleware.ParseUintID("id"), guestHandler.Delete)
			guests.PUT("/:id/rsvp", middleware.ParseUintID("id"), guestHandler.UpdateRSVPStatus)

			tables := guests.Group("/tables")
			{
				tables.POST("", guestHandler.CreateTable)
				tables.GET("", guestHandler.GetTables)
				tables.PUT("/:id", middleware.ParseUintID("id"), guestHandler.UpdateTable)
				tables.DELETE("/:id", middleware.ParseUintID("id"), guestHandler.DeleteTable)
			}

			guests.POST("/seat-assign", guestHandler.AssignSeat)
		}

		budgets := api.Group("/weddings/:wedding_id/budget")
		budgets.Use(middleware.JWTAuth(), middleware.ParseUintID("wedding_id"))
		{
			budgets.GET("/summary", budgetHandler.GetBudgetSummary)
			budgets.GET("/categories", budgetHandler.GetBudgetCategories)
			budgets.GET("/alerts", budgetHandler.CheckBudgetAlerts)
			budgets.GET("/items", budgetHandler.GetBudgetItems)
			budgets.POST("/items", budgetHandler.CreateBudgetItem)
			budgets.PUT("/items/:id", middleware.ParseUintID("id"), budgetHandler.UpdateBudgetItem)
			budgets.DELETE("/items/:id", middleware.ParseUintID("id"), budgetHandler.DeleteBudgetItem)
			budgets.GET("/payments", budgetHandler.GetPayments)
			budgets.POST("/payments", budgetHandler.RecordPayment)
		}

		tasks := api.Group("/weddings/:wedding_id/tasks")
		tasks.Use(middleware.JWTAuth(), middleware.ParseUintID("wedding_id"))
		{
			tasks.GET("", taskHandler.GetList)
			tasks.POST("", taskHandler.Create)
			tasks.GET("/categories", taskHandler.GetTaskCategories)
			tasks.GET("/:id", middleware.ParseUintID("id"), taskHandler.GetByID)
			tasks.PUT("/:id", middleware.ParseUintID("id"), taskHandler.Update)
			tasks.DELETE("/:id", middleware.ParseUintID("id"), taskHandler.Delete)
			tasks.PUT("/:id/status", middleware.ParseUintID("id"), taskHandler.UpdateStatus)
		}

		taskTemplates := api.Group("/task-templates")
		taskTemplates.Use(middleware.JWTAuth())
		{
			taskTemplates.GET("", taskHandler.GetTemplates)
			taskTemplates.POST("", taskHandler.CreateTemplate)
			taskTemplates.DELETE("/:id", middleware.ParseUintID("id"), taskHandler.DeleteTemplate)
			taskTemplates.POST("/:template_id/apply", middleware.ParseUintID("template_id"), taskHandler.ApplyTemplate)
		}

		documents := api.Group("/weddings/:wedding_id/documents")
		documents.Use(middleware.JWTAuth(), middleware.ParseUintID("wedding_id"))
		{
			documents.GET("", documentHandler.GetList)
			documents.GET("/categories", documentHandler.GetCategories)
			documents.POST("/upload", documentHandler.Upload)
			documents.GET("/:id", middleware.ParseUintID("id"), documentHandler.GetByID)
			documents.GET("/:id/download", middleware.ParseUintID("id"), documentHandler.Download)
			documents.PUT("/:id", middleware.ParseUintID("id"), documentHandler.Update)
			documents.DELETE("/:id", middleware.ParseUintID("id"), documentHandler.Delete)
			documents.POST("/:id/version", middleware.ParseUintID("id"), documentHandler.UploadNewVersion)
		}

		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.JWTAuth())
		{
			dashboard.GET("/stats", dashboardHandler.GetStats)
			dashboard.GET("/budget-chart", dashboardHandler.GetBudgetChart)
			dashboard.GET("/task-progress", dashboardHandler.GetTaskProgress)
			dashboard.GET("/upcoming-tasks", dashboardHandler.GetUpcomingTasks)
			dashboard.GET("/vendor-stats", dashboardHandler.GetVendorStats)
			dashboard.GET("/export", dashboardHandler.ExportReport)
		}

		logs := api.Group("/logs")
		logs.Use(middleware.JWTAuth(), middleware.RequireRole("admin"), middleware.Pagination())
		{
			logs.GET("", logHandler.GetLogs)
		}

		notifications := api.Group("/notifications")
		notifications.Use(middleware.JWTAuth())
		{
			notifications.GET("", notificationHandler.GetList)
			notifications.PUT("/:id/read", middleware.ParseUintID("id"), notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
		}
	}
}
