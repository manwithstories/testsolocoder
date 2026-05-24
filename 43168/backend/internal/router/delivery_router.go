package router

import (
	"github.com/gin-gonic/gin"

	"furniture-platform/internal/handler"
	"furniture-platform/internal/model"
	"furniture-platform/pkg/middleware"
)

// DeliveryRouter 配送安装路由依赖
type DeliveryRouter struct {
	deliveryHandler *handler.DeliveryHandler
	secret          string
}

// NewDeliveryRouter 创建配送安装路由
func NewDeliveryRouter(deliveryHandler *handler.DeliveryHandler, secret string) *DeliveryRouter {
	return &DeliveryRouter{
		deliveryHandler: deliveryHandler,
		secret:          secret,
	}
}

// Register 注册配送安装路由到已有的 v1 路由组
func (r *DeliveryRouter) Register(v1 *gin.RouterGroup) {
	authDeliveries := v1.Group("/deliveries", middleware.JWTAuth(r.secret))
	{
		// 所有登录用户可查看相关配送安装
		authDeliveries.GET("/", r.deliveryHandler.ListDeliveries)
		authDeliveries.GET("/:id", r.deliveryHandler.GetDelivery)

		// 业主创建预约
		ownerGroup := authDeliveries.Group("", middleware.RoleRequired(model.RoleOwner))
		{
			ownerGroup.POST("/", r.deliveryHandler.CreateDelivery)
		}

		// 业主/厂商/安装人员可更新预约（含确认、完成、取消）
		authDeliveries.PUT("/:id", r.deliveryHandler.UpdateDelivery)
	}
}
