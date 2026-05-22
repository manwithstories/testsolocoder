package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-Refresh-Token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Token-Will-Expire"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
