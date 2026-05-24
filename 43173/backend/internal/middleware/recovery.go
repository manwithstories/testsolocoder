package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"music-platform/pkg/logger"
	"music-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error("Broken pipe",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)

				response.InternalError(c, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Token")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Disposition")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func RateLimit(maxRequests int, duration time.Duration) gin.HandlerFunc {
	type client struct {
		count    int
		lastSeen time.Time
	}

	clients := make(map[string]*client)

	go func() {
		for {
			time.Sleep(time.Minute)
			for ip, c := range clients {
				if time.Since(c.lastSeen) > duration {
					delete(clients, ip)
				}
			}
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{count: 0, lastSeen: time.Now()}
		}

		clients[ip].count++
		clients[ip].lastSeen = time.Now()

		if clients[ip].count > maxRequests {
			response.Error(c, http.StatusTooManyRequests, 429, "请求过于频繁")
			c.Abort()
			return
		}

		c.Next()
	}
}
