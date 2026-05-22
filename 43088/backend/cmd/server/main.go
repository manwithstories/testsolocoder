package main

import (
	"fmt"
	"podcast-manager/internal/config"
	"podcast-manager/internal/database"
	"podcast-manager/internal/handlers"
	"podcast-manager/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig()

	setupLogger()

	database.InitDB()

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	setupRoutes(r)

	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	logrus.Infof("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func setupLogger() {
	level, err := logrus.ParseLevel(config.AppConfig.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		podcastHandler := handlers.NewPodcastHandler()
		episodeHandler := handlers.NewEpisodeHandler()
		noteHandler := handlers.NewNoteHandler()
		playlistHandler := handlers.NewPlaylistHandler()
		statsHandler := handlers.NewStatsHandler()
		importExportHandler := handlers.NewImportExportHandler()

		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		podcasts := api.Group("/podcasts")
		{
			podcasts.POST("", podcastHandler.AddPodcast)
			podcasts.GET("", podcastHandler.GetPodcasts)
			podcasts.GET("/:id", podcastHandler.GetPodcast)
			podcasts.PUT("/:id", podcastHandler.UpdatePodcast)
			podcasts.DELETE("/:id", podcastHandler.DeletePodcast)
			podcasts.POST("/:id/refresh", podcastHandler.RefreshPodcast)
		}

		api.GET("/episodes/new-count", podcastHandler.GetNewEpisodesCount)

		episodes := api.Group("/episodes")
		{
			episodes.GET("", episodeHandler.GetEpisodes)
			episodes.GET("/:id", episodeHandler.GetEpisode)
			episodes.POST("/:id/progress", episodeHandler.UpdatePlaybackProgress)
			episodes.GET("/:id/progress", episodeHandler.GetPlaybackProgress)
			episodes.POST("/:id/complete", episodeHandler.MarkAsCompleted)
			episodes.POST("/:id/history", episodeHandler.AddListeningHistory)
		}

		api.GET("/history", episodeHandler.GetListeningHistory)

		notes := api.Group("/notes")
		{
			notes.GET("/search", noteHandler.SearchNotes)
			notes.GET("/episode/:episode_id", noteHandler.GetNotes)
			notes.POST("/episode/:episode_id", noteHandler.AddNote)
			notes.PUT("/:id", noteHandler.UpdateNote)
			notes.DELETE("/:id", noteHandler.DeleteNote)
		}

		playlists := api.Group("/playlists")
		{
			playlists.POST("", playlistHandler.CreatePlaylist)
			playlists.GET("", playlistHandler.GetPlaylists)
			playlists.GET("/:id", playlistHandler.GetPlaylist)
			playlists.PUT("/:id", playlistHandler.UpdatePlaylist)
			playlists.DELETE("/:id", playlistHandler.DeletePlaylist)
			playlists.POST("/:id/episodes", playlistHandler.AddEpisodeToPlaylist)
			playlists.DELETE("/:id/episodes/:item_id", playlistHandler.RemoveEpisodeFromPlaylist)
			playlists.POST("/:id/reorder", playlistHandler.ReorderPlaylistItems)
		}

		stats := api.Group("/stats")
		{
			stats.GET("/listening", statsHandler.GetListeningStats)
			stats.GET("/distribution", statsHandler.GetPodcastDistribution)
			stats.GET("/completion", statsHandler.GetCompletionStats)
			stats.GET("/habits", statsHandler.GetListeningHabits)
		}

		export := api.Group("/export")
		{
			export.GET("/opml", importExportHandler.ExportOPML)
			export.GET("/history/csv", importExportHandler.ExportHistoryCSV)
			export.GET("/notes/csv", importExportHandler.ExportNotesCSV)
			export.GET("/notes/markdown", importExportHandler.ExportNotesMarkdown)
		}

		importGroup := api.Group("/import")
		{
			importGroup.POST("/opml", importExportHandler.ImportOPML)
		}
	}
}
