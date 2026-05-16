package routes

import (
	"finance-api/config"
	"finance-api/controllers"
	"finance-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	authCtrl := controllers.NewAuthController(cfg)
	accountCtrl := controllers.NewAccountController()
	categoryCtrl := controllers.NewCategoryController()
	transactionCtrl := controllers.NewTransactionController()
	statsCtrl := controllers.NewStatisticsController()
	budgetCtrl := controllers.NewBudgetController()
	exportCtrl := controllers.NewExportController()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
			auth.GET("/profile", middleware.AuthMiddleware(cfg), authCtrl.GetProfile)
		}

		accounts := api.Group("/accounts")
		accounts.Use(middleware.AuthMiddleware(cfg))
		{
			accounts.POST("", accountCtrl.Create)
			accounts.GET("", accountCtrl.List)
			accounts.GET("/:id", accountCtrl.Get)
			accounts.PUT("/:id", accountCtrl.Update)
			accounts.DELETE("/:id", accountCtrl.Delete)
		}

		categories := api.Group("/categories")
		categories.Use(middleware.AuthMiddleware(cfg))
		{
			categories.POST("", categoryCtrl.Create)
			categories.GET("", categoryCtrl.List)
			categories.GET("/:id", categoryCtrl.Get)
			categories.PUT("/:id", categoryCtrl.Update)
			categories.DELETE("/:id", categoryCtrl.Delete)
		}

		transactions := api.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware(cfg))
		{
			transactions.POST("", transactionCtrl.Create)
			transactions.GET("", transactionCtrl.List)
			transactions.GET("/:id", transactionCtrl.Get)
			transactions.PUT("/:id", transactionCtrl.Update)
			transactions.DELETE("/:id", transactionCtrl.Delete)
		}

		statistics := api.Group("/statistics")
		statistics.Use(middleware.AuthMiddleware(cfg))
		{
			statistics.GET("/monthly", statsCtrl.GetMonthly)
		}

		budgets := api.Group("/budgets")
		budgets.Use(middleware.AuthMiddleware(cfg))
		{
			budgets.POST("", budgetCtrl.Create)
			budgets.GET("", budgetCtrl.List)
			budgets.GET("/:id", budgetCtrl.Get)
			budgets.PUT("/:id", budgetCtrl.Update)
			budgets.DELETE("/:id", budgetCtrl.Delete)
		}

		export := api.Group("/export")
		export.Use(middleware.AuthMiddleware(cfg))
		{
			export.GET("/transactions", exportCtrl.ExportTransactions)
		}
	}
}
