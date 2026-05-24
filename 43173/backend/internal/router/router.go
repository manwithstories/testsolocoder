package router

import (
	"music-platform/internal/handler"
	"music-platform/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	userHandler := handler.NewUserHandler()
	workHandler := handler.NewWorkHandler()
	communityHandler := handler.NewCommunityHandler()
	rankingHandler := handler.NewRankingHandler()
	eventHandler := handler.NewEventHandler()
	revenueHandler := handler.NewRevenueHandler()

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)

		api.GET("/ranking", rankingHandler.GetRanking)
		api.GET("/ranking/daily", rankingHandler.GetDailyRanking)
		api.GET("/ranking/weekly", rankingHandler.GetWeeklyRanking)
		api.GET("/ranking/monthly", rankingHandler.GetMonthlyRanking)
		api.GET("/ranking/hot", rankingHandler.GetHotRanking)
		api.GET("/ranking/work/:id", rankingHandler.GetWorkRanking)

		api.GET("/works", workHandler.ListWorks)
		api.GET("/works/:id", workHandler.GetWorkByID)
		api.GET("/works/artist/:artist_id", workHandler.GetArtistWorks)
		api.POST("/works/:id/play", middleware.JWTAuthOptional(), workHandler.RecordPlay)

		api.GET("/albums", workHandler.ListAlbums)
		api.GET("/albums/:id", workHandler.GetAlbumByID)

		api.GET("/tags", workHandler.ListTags)

		api.GET("/events", eventHandler.ListEvents)
		api.GET("/events/:id", eventHandler.GetEventByID)
		api.GET("/events/:id/stats", eventHandler.GetEventStats)
		api.GET("/events/:id/seats", eventHandler.GetSeatAvailability)

		api.GET("/playlists", communityHandler.ListPlaylists)
		api.GET("/playlists/:id", communityHandler.GetPlaylistByID)

		api.GET("/comments", communityHandler.GetComments)

		api.GET("/users/:id", userHandler.GetUserByID)
		api.GET("/users/:id/followers", communityHandler.GetFollowers)
		api.GET("/users/:id/followings", communityHandler.GetFollowings)
		api.GET("/users/:id/follower-count", communityHandler.GetFollowerCount)
		api.GET("/users/:id/following-count", communityHandler.GetFollowingCount)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/auth/profile", userHandler.GetProfile)
			auth.PUT("/auth/profile", userHandler.UpdateProfile)

			auth.GET("/auth/artist-info", userHandler.GetArtistInfo)
			auth.PUT("/auth/artist-info", userHandler.UpdateArtistInfo)
			auth.GET("/auth/balance", userHandler.GetBalance)

			auth.POST("/works/upload", workHandler.UploadWork)
			auth.PUT("/works/:id", workHandler.UpdateWork)
			auth.DELETE("/works/:id", workHandler.DeleteWork)
			auth.POST("/works/batch-publish", workHandler.BatchPublish)

			auth.POST("/albums", workHandler.CreateAlbum)
			auth.PUT("/albums/:id", workHandler.UpdateAlbum)
			auth.DELETE("/albums/:id", workHandler.DeleteAlbum)
			auth.POST("/albums/:album_id/works", workHandler.AddWorkToAlbum)

			auth.POST("/follow", communityHandler.Follow)
			auth.POST("/unfollow", communityHandler.Unfollow)
			auth.GET("/follow/is-following/:id", communityHandler.IsFollowing)

			auth.POST("/comments", communityHandler.CreateComment)
			auth.DELETE("/comments/:id", communityHandler.DeleteComment)

			auth.POST("/playlists", communityHandler.CreatePlaylist)
			auth.PUT("/playlists/:id", communityHandler.UpdatePlaylist)
			auth.DELETE("/playlists/:id", communityHandler.DeletePlaylist)
			auth.POST("/playlists/works", communityHandler.AddWorkToPlaylist)
			auth.DELETE("/playlists/works", communityHandler.RemoveWorkFromPlaylist)

			auth.GET("/notifications", communityHandler.GetNotifications)
			auth.PUT("/notifications/:id/read", communityHandler.MarkNotificationAsRead)
			auth.PUT("/notifications/read-all", communityHandler.MarkAllNotificationsAsRead)
			auth.GET("/notifications/unread-count", communityHandler.GetUnreadNotificationCount)

			auth.GET("/play-records", communityHandler.GetPlayRecords)

			auth.GET("/my-playlists", communityHandler.GetMyPlaylists)
			auth.GET("/follow/artists", communityHandler.GetMyFollowingArtists)
			auth.GET("/follow/users", communityHandler.GetMyFollowingUsers)
			auth.GET("/follow/followers", communityHandler.GetMyFollowers)

			auth.POST("/events", eventHandler.CreateEvent)
			auth.PUT("/events/:id", eventHandler.UpdateEvent)
			auth.DELETE("/events/:id", eventHandler.DeleteEvent)
			auth.POST("/events/:id/publish", eventHandler.PublishEvent)
			auth.POST("/events/purchase", eventHandler.PurchaseTicket)

			auth.GET("/orders", eventHandler.GetOrdersByUser)
			auth.GET("/orders/:id", eventHandler.GetOrderByID)
			auth.GET("/orders/artist/:artist_id", eventHandler.GetOrdersByArtist)

			auth.GET("/tickets", eventHandler.GetTicketsByUser)
			auth.POST("/tickets/:id/use", eventHandler.UseTicket)

			auth.GET("/revenue/records", revenueHandler.GetRevenueRecords)
			auth.GET("/revenue/records/artist/:artist_id", revenueHandler.GetArtistRevenueRecords)
			auth.GET("/revenue/total", revenueHandler.GetTotalRevenue)
			auth.GET("/revenue/summary", revenueHandler.GetRevenueSummary)

			auth.POST("/withdraw", revenueHandler.RequestWithdraw)
			auth.GET("/withdraw", revenueHandler.GetWithdrawRequests)
			auth.GET("/withdraw/status-list", revenueHandler.GetWithdrawStatusList)

			auth.GET("/subscriptions", revenueHandler.GetSubscriptions)
			auth.GET("/subscriptions/artist/:artist_id", revenueHandler.GetArtistSubscribers)

			auth.GET("/stats/daily", revenueHandler.GetDailyStats)
			auth.GET("/stats/daily/artist/:artist_id", revenueHandler.GetArtistDailyStats)
			auth.GET("/stats/artist/:artist_id", revenueHandler.GetArtistStats)

			auth.GET("/export/revenue", revenueHandler.ExportRevenueExcel)
			auth.GET("/export/withdraw", revenueHandler.ExportWithdrawExcel)

			admin := auth.Group("")
			admin.Use(middleware.RoleAuth("admin"))
			{
				admin.GET("/users", userHandler.ListUsers)
				admin.PUT("/users/:id/status", userHandler.UpdateUserStatus)
				admin.PUT("/users/:id/role", userHandler.UpdateUserRole)
				admin.POST("/users/:id/verify", userHandler.VerifyArtist)

				admin.PUT("/works/:id/approve", workHandler.ApproveWork)
				admin.PUT("/works/:id/reject", workHandler.RejectWork)
				admin.PUT("/works/:id/status", workHandler.UpdateWorkStatus)

				admin.GET("/admin/withdraw", revenueHandler.GetAllWithdrawRequests)
				admin.PUT("/admin/withdraw/:id/approve", revenueHandler.ApproveWithdraw)
				admin.PUT("/admin/withdraw/:id/reject", revenueHandler.RejectWithdraw)
				admin.PUT("/admin/withdraw/:id/paid", revenueHandler.MarkWithdrawPaid)
				admin.POST("/admin/revenue/settle", revenueHandler.SettleRevenue)

				admin.GET("/admin/operation-logs", revenueHandler.GetOperationLogs)
			}
		}
	}
}
