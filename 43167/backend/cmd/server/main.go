package main

import (
	"log"
	"os"

	"watchplatform/internal/config"
	"watchplatform/internal/database"
	"watchplatform/internal/handler"
	"watchplatform/internal/logger"
	"watchplatform/internal/middleware"
	"watchplatform/internal/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	_ = os.MkdirAll(config.Cfg.UploadDir, 0o755)
	logger.Init()
	defer logger.Close()
	database.Init()

	r := gin.New()
	r.Use(logger.GinLogger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Origin"},
		ExposeHeaders:    []string{"Content-Disposition"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", config.Cfg.UploadDir)

	api := r.Group("/api")
	{
		api.POST("/auth/register", handler.Register)
		api.POST("/auth/login", handler.Login)
		api.GET("/watches", handler.ListWatches)
		api.GET("/watches/:id", handler.GetWatch)
		api.GET("/stats/overview", handler.StatsOverview)
		api.GET("/stats/brands", handler.StatsBrands)
		api.GET("/stats/export", handler.StatsExport)
	}

	auth := api.Group("")
	auth.Use(middleware.JWTAuth(database.DB))
	{
		auth.GET("/me", handler.Me)
		auth.PUT("/me", handler.UpdateProfile)
		auth.POST("/me/avatar", handler.UploadAvatar)
		auth.PUT("/me/password", handler.ChangePassword)

		auth.POST("/watches", middleware.RequireRoles(model.RoleSeller), handler.CreateWatch)
		auth.PUT("/watches/:id", middleware.RequireRoles(model.RoleSeller), handler.UpdateWatch)
		auth.DELETE("/watches/:id", middleware.RequireRoles(model.RoleSeller), handler.DeleteWatch)
		auth.POST("/watches/:id/photos", middleware.RequireRoles(model.RoleSeller), handler.UploadWatchPhotos)

		auth.POST("/auth-orders", handler.CreateAuthOrder)
		auth.GET("/auth-orders", handler.ListAuthOrders)
		auth.GET("/auth-orders/:id", handler.GetAuthOrder)
		auth.POST("/auth-orders/:id/assign", middleware.RequireRoles(model.RoleAppraiser), handler.AssignAuthOrder)
		auth.POST("/auth-orders/:id/report", middleware.RequireRoles(model.RoleAppraiser), handler.SubmitAuthReport)
		auth.GET("/auth-orders/:id/report/pdf", handler.DownloadReportPDF)

		auth.POST("/trades", middleware.RequireRoles(model.RoleSeller), handler.CreateTrade)
		auth.GET("/trades", handler.ListTrades)
		auth.GET("/trades/:id", handler.GetTrade)
		auth.POST("/trades/:id/bids", middleware.RequireRoles(model.RoleBuyer), handler.PlaceBid)
		auth.POST("/trades/:id/accept", middleware.RequireRoles(model.RoleSeller), handler.AcceptBid)
		auth.POST("/trades/:id/ship", middleware.RequireRoles(model.RoleSeller), handler.ShipTrade)
		auth.POST("/trades/:id/complete", handler.CompleteTrade)
		auth.POST("/trades/:id/cancel", handler.CancelTrade)

		auth.GET("/favorites/groups", handler.ListFavoriteGroups)
		auth.POST("/favorites/groups", handler.CreateFavoriteGroup)
		auth.DELETE("/favorites/groups/:id", handler.DeleteFavoriteGroup)
		auth.GET("/favorites", handler.ListFavorites)
		auth.POST("/favorites", handler.AddFavorite)
		auth.DELETE("/favorites/:watch_id", handler.RemoveFavorite)

		auth.POST("/reviews", handler.CreateReview)
		auth.GET("/reviews", handler.ListReviews)

		auth.GET("/messages", handler.ListMessages)
		auth.POST("/messages/:id/read", handler.MarkMessageRead)
		auth.POST("/messages/read-all", handler.MarkAllRead)
	}

	if err := r.Run(":" + config.Cfg.Port); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
