package main

import (
	"booklibrary/internal/config"
	"booklibrary/internal/database"
	"booklibrary/internal/handlers"
	"booklibrary/internal/logger"
	"booklibrary/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.Load(); err != nil {
		panic("Failed to load config: " + err.Error())
	}

	if err := logger.Init(); err != nil {
		panic("Failed to init logger: " + err.Error())
	}

	if err := database.Init(); err != nil {
		logger.Fatalf("Failed to init database: %v", err)
	}

	gin.SetMode(config.AppConfig.Server.Mode)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	r.Static(config.AppConfig.Upload.AccessURL, config.AppConfig.Upload.SavePath)

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Book Library API is running",
			})
		})

		bookHandler := handlers.NewBookHandler()
		books := api.Group("/books")
		{
			books.GET("", bookHandler.GetBooks)
			books.GET("/reading", bookHandler.GetCurrentlyReading)
			books.GET("/:id", bookHandler.GetBook)
			books.POST("", bookHandler.CreateBook)
			books.PUT("/:id", bookHandler.UpdateBook)
			books.DELETE("/:id", bookHandler.DeleteBook)
			books.POST("/:id/cover", bookHandler.UploadCover)
			books.GET("/isbn/:isbn", bookHandler.FetchByISBN)
			books.PATCH("/:id/progress", bookHandler.UpdateReadingProgress)
			books.PATCH("/:id/status", bookHandler.UpdateStatus)
		}

		tagHandler := handlers.NewTagHandler()
		tags := api.Group("/tags")
		{
			tags.GET("", tagHandler.GetTags)
			tags.GET("/:id", tagHandler.GetTag)
			tags.POST("", tagHandler.CreateTag)
			tags.PUT("/:id", tagHandler.UpdateTag)
			tags.DELETE("/:id", tagHandler.DeleteTag)
		}

		categoryHandler := handlers.NewCategoryHandler()
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.POST("", categoryHandler.CreateCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		noteHandler := handlers.NewNoteHandler()
		notes := api.Group("/notes")
		{
			notes.GET("/book/:bookId", noteHandler.GetNotesByBook)
			notes.GET("/:id", noteHandler.GetNote)
			notes.POST("", noteHandler.CreateNote)
			notes.PUT("/:id", noteHandler.UpdateNote)
			notes.DELETE("/:id", noteHandler.DeleteNote)
		}

		borrowHandler := handlers.NewBorrowHandler()
		borrows := api.Group("/borrows")
		{
			borrows.GET("", borrowHandler.GetBorrows)
			borrows.GET("/overdue", borrowHandler.GetOverdue)
			borrows.GET("/book/:bookId", borrowHandler.GetBorrowByBook)
			borrows.GET("/:id", borrowHandler.GetBorrow)
			borrows.POST("", borrowHandler.CreateBorrow)
			borrows.POST("/:id/return", borrowHandler.ReturnBook)
			borrows.DELETE("/:id", borrowHandler.DeleteBorrow)
		}

		statsHandler := handlers.NewStatsHandler()
		stats := api.Group("/stats")
		{
			stats.GET("/overview", statsHandler.GetOverview)
			stats.GET("/yearly-trend", statsHandler.GetYearlyTrend)
			stats.GET("/heatmap", statsHandler.GetReadingHeatmap)
			stats.GET("/duration", statsHandler.GetDurationDistribution)
			stats.GET("/categories", statsHandler.GetCategoryStats)
			stats.GET("/tags", statsHandler.GetTagStats)
		}

		goalHandler := handlers.NewGoalHandler()
		goals := api.Group("/goals")
		{
			goals.GET("", goalHandler.GetGoals)
			goals.GET("/yearly-progress", goalHandler.GetYearlyGoalProgress)
			goals.GET("/:id", goalHandler.GetGoal)
			goals.POST("", goalHandler.CreateGoal)
			goals.PUT("/:id", goalHandler.UpdateGoal)
			goals.DELETE("/:id", goalHandler.DeleteGoal)
		}
	}

	logger.Infof("Server starting on %s", config.AppConfig.Server.Address())
	if err := r.Run(config.AppConfig.Server.Address()); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
