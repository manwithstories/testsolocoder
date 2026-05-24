package router

import (
	"github.com/gin-gonic/gin"

	"furniture-platform/internal/handler"
	"furniture-platform/pkg/middleware"
)

// Router 依赖
type Router struct {
	userHandler    *handler.UserHandler
	productHandler *handler.ProductHandler
	secret         string
}

// NewRouter 创建路由
func NewRouter(userHandler *handler.UserHandler, productHandler *handler.ProductHandler, secret string) *Router {
	return &Router{
		userHandler:    userHandler,
		productHandler: productHandler,
		secret:         secret,
	}
}

// Register 注册所有路由
func (r *Router) Register(e *gin.Engine) {
	v1 := e.Group("/api/v1")
	{
		// 认证相关（无需登录）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userHandler.Register)
			auth.POST("/login", r.userHandler.Login)
		}

		// 需要登录的用户相关
		users := v1.Group("/users", middleware.JWTAuth(r.secret))
		{
			users.GET("/profile", r.userHandler.Profile)
			users.PUT("/profile", r.userHandler.UpdateProfile)
			users.PUT("/password", r.userHandler.ChangePassword)
			users.GET("/", r.userHandler.ListUsers)
			users.GET("/:id", r.userHandler.GetUser)
		}

		// 产品相关
		productRouter := NewProductRouter(r.productHandler, r.secret)
		productRouter.Register(v1)
	}
}
