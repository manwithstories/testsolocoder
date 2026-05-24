package router

import (
	"github.com/gin-gonic/gin"

	"furniture-platform/internal/handler"
	"furniture-platform/internal/model"
	"furniture-platform/pkg/middleware"
)

// ReviewRouter 售后评价路由依赖
type ReviewRouter struct {
	reviewHandler *handler.ReviewHandler
	secret        string
}

// NewReviewRouter 创建售后评价路由
func NewReviewRouter(reviewHandler *handler.ReviewHandler, secret string) *ReviewRouter {
	return &ReviewRouter{
		reviewHandler: reviewHandler,
		secret:        secret,
	}
}

// Register 注册售后评价路由到已有的 v1 路由组
func (r *ReviewRouter) Register(v1 *gin.RouterGroup) {
	authReviews := v1.Group("/reviews", middleware.JWTAuth(r.secret))
	{
		// 所有登录用户可查看评价
		authReviews.GET("/", r.reviewHandler.ListReviews)
		authReviews.GET("/:id", r.reviewHandler.GetReview)

		// 业主创建评价
		ownerGroup := authReviews.Group("", middleware.RoleRequired(model.RoleOwner))
		{
			ownerGroup.POST("/", r.reviewHandler.CreateReview)
		}
	}
}
