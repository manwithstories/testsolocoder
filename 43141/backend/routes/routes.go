package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sports-league/handlers"
	"sports-league/pkg/middleware"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authH := handlers.NewAuthHandler(db)
	leagueH := handlers.NewLeagueHandler(db)
	teamH := handlers.NewTeamHandler(db)
	matchH := handlers.NewMatchHandler(db)
	refereeH := handlers.NewRefereeHandler(db)
	feeH := handlers.NewFeeHandler(db)
	notifH := handlers.NewNotificationHandler(db)
	statsH := handlers.NewStatsHandler(db)
	exportH := handlers.NewExportHandler(db)

	api := r.Group("/api")
	{
		api.POST("/auth/register", authH.Register)
		api.POST("/auth/login", authH.Login)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/auth/me", authH.Me)

			auth.GET("/leagues", leagueH.List)
			auth.GET("/leagues/:id", leagueH.Get)
			auth.GET("/leagues/:id/seasons/:season_id", leagueH.GetSeason)

			auth.GET("/teams", teamH.List)
			auth.GET("/teams/:id", teamH.Get)

			auth.GET("/matches", matchH.ListMatches)
			auth.GET("/matches/:id", matchH.GetMatch)
			auth.GET("/venues", matchH.ListVenues)
			auth.GET("/venues/check-conflict", matchH.CheckVenueConflict)
			auth.GET("/seasons/:season_id/standings", matchH.GetStandings)

			auth.GET("/referees", refereeH.ListReferees)
			auth.GET("/referees/assignments", refereeH.ListAssignments)
			auth.PUT("/referees/assignments/:id", refereeH.Respond)

			auth.GET("/fees", feeH.List)
			auth.GET("/fees/:id", feeH.Get)
			auth.GET("/fees/seasons/:season_id/report", feeH.Report)

			auth.GET("/notifications", notifH.List)
			auth.PUT("/notifications/:id/read", notifH.MarkRead)
			auth.PUT("/notifications/read-all", notifH.MarkAllRead)
			auth.GET("/notifications/unread-count", notifH.UnreadCount)

			auth.GET("/stats/rankings", statsH.GetRankings)
			auth.GET("/stats/players/:id", statsH.GetPlayerStats)

			auth.GET("/export/schedule/:season_id", exportH.ExportSchedule)
			auth.GET("/export/standings/:season_id", exportH.ExportStandings)
			auth.GET("/export/stats", exportH.ExportStats)
			auth.GET("/export/pdf/:season_id", exportH.ExportPDF)
		}

		admin := api.Group("")
		admin.Use(middleware.JWTAuth(), middleware.RequireRoles("admin"))
		{
			admin.POST("/leagues", leagueH.Create)
			admin.PUT("/leagues/:id", leagueH.Update)
			admin.DELETE("/leagues/:id", leagueH.Delete)
			admin.POST("/leagues/:id/seasons", leagueH.CreateSeason)
			admin.PUT("/leagues/:id/seasons/:season_id", leagueH.UpdateSeason)

			admin.POST("/teams", teamH.Create)
			admin.PUT("/teams/:id", teamH.Update)
			admin.DELETE("/teams/:id", teamH.Delete)
			admin.POST("/teams/:id/players", teamH.AddPlayer)
			admin.PUT("/teams/:id/players/:player_id", teamH.UpdatePlayer)
			admin.DELETE("/teams/:id/players/:player_id", teamH.DeletePlayer)
			admin.POST("/teams/:id/register", teamH.RegisterSeason)
			admin.GET("/registrations", teamH.ListRegistrations)
			admin.PUT("/registrations/:reg_id/approve", teamH.ApproveRegistration)

			admin.POST("/seasons/:season_id/generate-schedule", matchH.GenerateSchedule)
			admin.PUT("/matches/:id", matchH.UpdateMatch)
			admin.DELETE("/matches/:id", matchH.DeleteMatch)
			admin.POST("/matches/:id/report-score", matchH.ReportScore)
			admin.POST("/seasons/:season_id/generate-knockout", matchH.GenerateKnockout)

			admin.POST("/venues", matchH.CreateVenue)
			admin.PUT("/venues/:id", matchH.UpdateVenue)

			admin.POST("/referees/assign", refereeH.Assign)

			admin.POST("/fees", feeH.Create)
			admin.PUT("/fees/:id/paid", feeH.MarkPaid)
			admin.DELETE("/fees/:id", feeH.Delete)

			admin.POST("/stats", statsH.AddStat)
			admin.DELETE("/stats/:id", statsH.DeleteStat)
		}

		captain := api.Group("")
		captain.Use(middleware.JWTAuth(), middleware.RequireRoles("admin", "captain"))
		{
			captain.POST("/teams/:id/register", teamH.RegisterSeason)
		}
	}

	return r
}
