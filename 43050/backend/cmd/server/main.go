package main

import (
	"log"
	"os"

	"splitwise-clone/internal/database"
	"splitwise-clone/internal/handlers"
	"splitwise-clone/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using defaults")
	}

	database.InitDB()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.GET("/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
		}

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/search", handlers.SearchUsers)
			users.GET("/:id", handlers.GetUserByID)
		}

		groups := api.Group("/groups")
		groups.Use(middleware.AuthMiddleware())
		{
			groups.POST("", handlers.CreateGroup)
			groups.GET("", handlers.GetUserGroups)
			groups.GET("/:id", handlers.GetGroupByID)
			groups.GET("/:id/members", handlers.GetGroupMembers)
			groups.POST("/join", handlers.JoinGroup)
			groups.POST("/:id/leave", handlers.LeaveGroup)

			expenses := groups.Group("/:groupId/expenses")
			{
				expenses.POST("", handlers.CreateExpense)
				expenses.GET("", handlers.GetGroupExpenses)
			}

			balances := groups.Group("/:groupId")
			{
				balances.GET("/balances", handlers.GetGroupBalances)
				balances.GET("/transfers", handlers.GetOptimalTransfers)
				balances.GET("/stats", handlers.GetGroupExpenseStats)
				balances.POST("/settlements", handlers.CreateSettlement)
				balances.GET("/settlements", handlers.GetGroupSettlements)
			}
		}

		expenses := api.Group("/expenses")
		expenses.Use(middleware.AuthMiddleware())
		{
			expenses.GET("/:id", handlers.GetExpenseByID)
			expenses.PUT("/:id", handlers.UpdateExpense)
			expenses.DELETE("/:id", handlers.DeleteExpense)
		}

		settlements := api.Group("/settlements")
		settlements.Use(middleware.AuthMiddleware())
		{
			settlements.PATCH("/:id/paid", handlers.MarkSettlementPaid)
		}

		stats := api.Group("/stats")
		stats.Use(middleware.AuthMiddleware())
		{
			stats.GET("/summary", handlers.GetUserStatistics)
			stats.GET("/monthly", handlers.GetUserMonthlyStats)
			stats.GET("/history", handlers.GetUserExpenseHistory)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	r.Run(":" + port)
}
