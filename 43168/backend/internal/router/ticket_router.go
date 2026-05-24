package router

import (
	"github.com/gin-gonic/gin"

	"furniture-platform/internal/handler"
	"furniture-platform/internal/model"
	"furniture-platform/pkg/middleware"
)

// TicketRouter 售后工单路由依赖
type TicketRouter struct {
	ticketHandler *handler.TicketHandler
	secret        string
}

// NewTicketRouter 创建售后工单路由
func NewTicketRouter(ticketHandler *handler.TicketHandler, secret string) *TicketRouter {
	return &TicketRouter{
		ticketHandler: ticketHandler,
		secret:        secret,
	}
}

// Register 注册售后工单路由到已有的 v1 路由组
func (r *TicketRouter) Register(v1 *gin.RouterGroup) {
	authTickets := v1.Group("/tickets", middleware.JWTAuth(r.secret))
	{
		// 所有登录用户可查看工单
		authTickets.GET("/", r.ticketHandler.ListTickets)
		authTickets.GET("/:id", r.ticketHandler.GetTicket)

		// 业主创建工单
		ownerGroup := authTickets.Group("", middleware.RoleRequired(model.RoleOwner))
		{
			ownerGroup.POST("/", r.ticketHandler.CreateTicket)
		}

		// 业主/厂商可更新工单（状态流转、内容修改）
		authTickets.PUT("/:id", r.ticketHandler.UpdateTicket)
	}
}
