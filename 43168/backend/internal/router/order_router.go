package router

import (
	"github.com/gin-gonic/gin"

	"furniture-platform/internal/handler"
	"furniture-platform/internal/model"
	"furniture-platform/pkg/middleware"
)

// OrderRouter 订单路由依赖
type OrderRouter struct {
	orderHandler *handler.OrderHandler
	secret       string
}

// NewOrderRouter 创建订单路由
func NewOrderRouter(orderHandler *handler.OrderHandler, secret string) *OrderRouter {
	return &OrderRouter{
		orderHandler: orderHandler,
		secret:       secret,
	}
}

// Register 注册订单路由到已有的 v1 路由组
func (r *OrderRouter) Register(v1 *gin.RouterGroup) {
	authOrders := v1.Group("/orders", middleware.JWTAuth(r.secret))
	{
		// 所有登录用户可查看相关订单
		authOrders.GET("/", r.orderHandler.ListOrders)
		authOrders.GET("/:id", r.orderHandler.GetOrder)

		// 业主发起询价
		ownerGroup := authOrders.Group("", middleware.RoleRequired(model.RoleOwner))
		{
			ownerGroup.POST("/", r.orderHandler.CreateInquiry)
			ownerGroup.POST("/:id/confirm", r.orderHandler.ConfirmOrder)
			ownerGroup.POST("/:id/cancel", r.orderHandler.CancelOrder)
			// 业主可推进已确认订单到已付款
			ownerGroup.POST("/:id/status", r.orderHandler.UpdateStatus)
		}

		// 厂商报价与生产状态流转
		manufacturerGroup := authOrders.Group("", middleware.RoleRequired(model.RoleManufacturer))
		{
			manufacturerGroup.POST("/:id/quote", r.orderHandler.Quote)
			manufacturerGroup.POST("/:id/cancel", r.orderHandler.CancelOrder)
			manufacturerGroup.POST("/:id/status", r.orderHandler.UpdateStatus)
		}
	}
}
