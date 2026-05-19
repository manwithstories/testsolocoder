package router

import (
	"github.com/gin-gonic/gin"
	"github.com/notification-center/internal/handlers"
	"github.com/notification-center/internal/middleware"
)

func SetupRouter(
	channelHandler *handlers.ChannelHandler,
	templateHandler *handlers.TemplateHandler,
	recipientHandler *handlers.RecipientHandler,
	messageHandler *handlers.MessageHandler,
	webhookHandler *handlers.WebhookHandler,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api/v1")
	{
		channels := api.Group("/channels")
		{
			channels.POST("", channelHandler.Create)
			channels.GET("", channelHandler.List)
			channels.GET("/:id", channelHandler.GetByID)
			channels.PUT("/:id", channelHandler.Update)
			channels.DELETE("/:id", channelHandler.Delete)
			channels.POST("/:id/enable", channelHandler.Enable)
			channels.POST("/:id/disable", channelHandler.Disable)
			channels.PUT("/:id/priority", channelHandler.UpdatePriority)
			channels.POST("/:id/test", channelHandler.TestConnection)
		}

		templates := api.Group("/templates")
		{
			templates.POST("", templateHandler.Create)
			templates.GET("", templateHandler.List)
			templates.GET("/:id", templateHandler.GetByID)
			templates.PUT("/:id", templateHandler.Update)
			templates.DELETE("/:id", templateHandler.Delete)
			templates.POST("/:id/render", templateHandler.Render)
			templates.POST("/validate", templateHandler.Validate)
		}

		recipients := api.Group("/recipients")
		{
			recipients.POST("", recipientHandler.Create)
			recipients.GET("", recipientHandler.List)
			recipients.GET("/:id", recipientHandler.GetByID)
			recipients.PUT("/:id", recipientHandler.Update)
			recipients.DELETE("/:id", recipientHandler.Delete)
			recipients.POST("/:id/tags", recipientHandler.AddTags)
			recipients.DELETE("/:id/tags", recipientHandler.RemoveTags)
			recipients.POST("/:id/groups", recipientHandler.AddToGroups)
			recipients.DELETE("/:id/groups", recipientHandler.RemoveFromGroups)
			recipients.POST("/import", recipientHandler.BatchImport)
			recipients.GET("/export", recipientHandler.Export)
		}

		tags := api.Group("/tags")
		{
			tags.POST("", recipientHandler.CreateTag)
			tags.GET("", recipientHandler.ListTags)
			tags.DELETE("/:id", recipientHandler.DeleteTag)
		}

		groups := api.Group("/groups")
		{
			groups.POST("", recipientHandler.CreateGroup)
			groups.GET("", recipientHandler.ListGroups)
			groups.GET("/:id", recipientHandler.GetGroup)
			groups.DELETE("/:id", recipientHandler.DeleteGroup)
		}

		messages := api.Group("/messages")
		{
			messages.POST("/send", messageHandler.Send)
			messages.POST("/batch-send", messageHandler.BatchSend)
			messages.GET("/:message_id/status", messageHandler.GetStatus)
			messages.GET("", messageHandler.List)
			messages.GET("/stats", messageHandler.GetStats)
			messages.GET("/stats/channel", messageHandler.GetChannelStats)
			messages.GET("/stats/daily", messageHandler.GetDailyStats)
			messages.GET("/queue/stats", messageHandler.GetQueueStats)
		}

		webhooks := api.Group("/webhooks")
		{
			webhooks.POST("", webhookHandler.Create)
			webhooks.GET("", webhookHandler.List)
			webhooks.GET("/:id", webhookHandler.GetByID)
			webhooks.PUT("/:id", webhookHandler.Update)
			webhooks.DELETE("/:id", webhookHandler.Delete)
		}

		health := api.Group("/health")
		{
			health.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"code":    0,
					"message": "ok",
					"data": gin.H{
						"status": "running",
					},
				})
			})
		}
	}

	return r
}
