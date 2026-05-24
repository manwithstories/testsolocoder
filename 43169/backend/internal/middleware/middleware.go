package middleware

import (
	"strings"
	"time"

	"matchmaking-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "未登录，请先登录")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			utils.Unauthorized(c, "无效的认证格式")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.Unauthorized(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

func MatchmakerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || (role != "matchmaker" && role != "admin") {
			utils.Forbidden(c, "需要红娘或管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

func VerifiedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		var user struct {
			VerifyStatus string `json:"verify_status"`
		}
		if err := utils.DB.Table("users").Select("verify_status").Where("id = ?", userID).Scan(&user).Error; err != nil {
			utils.ServerError(c, "服务器错误")
			c.Abort()
			return
		}

		if user.VerifyStatus != "verified" {
			utils.Forbidden(c, "请先完成实名认证后再进行互动")
			c.Abort()
			return
		}

		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		latency := time.Since(start)
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		utils.DB.Create(&map[string]interface{}{
			"module":    "api",
			"action":    method + " " + path,
			"ip":        clientIP,
			"detail":    "status: " + string(rune(statusCode)),
			"created_at": time.Now(),
		})
		_ = latency
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				utils.ServerError(c, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimitMiddleware(maxRequests int, duration time.Duration) gin.HandlerFunc {
	var requestCounts = make(map[string]int)
	var lastCleanup = time.Now()

	return func(c *gin.Context) {
		now := time.Now()
		if now.Sub(lastCleanup) > duration {
			requestCounts = make(map[string]int)
			lastCleanup = now
		}

		ip := c.ClientIP()
		requestCounts[ip]++

		if requestCounts[ip] > maxRequests {
			utils.Fail(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
