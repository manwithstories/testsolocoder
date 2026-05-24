package router

import (
	"github.com/gin-gonic/gin"

	"furniture-platform/internal/handler"
	"furniture-platform/pkg/middleware"
)

// StatisticsRouter 数据统计路由依赖
type StatisticsRouter struct {
	statisticsHandler *handler.StatisticsHandler
	secret            string
}

// NewStatisticsRouter 创建数据统计路由
func NewStatisticsRouter(statisticsHandler *handler.StatisticsHandler, secret string) *StatisticsRouter {
	return &StatisticsRouter{
		statisticsHandler: statisticsHandler,
		secret:            secret,
	}
}

// Register 注册数据统计路由到已有的 v1 路由组
func (r *StatisticsRouter) Register(v1 *gin.RouterGroup) {
	authStats := v1.Group("/statistics", middleware.JWTAuth(r.secret))
	{
		// 所有登录用户可查看统计数据
		authStats.GET("/sales/trend", r.statisticsHandler.GetSalesTrend)
		authStats.GET("/customer/profile", r.statisticsHandler.GetCustomerProfile)
		authStats.GET("/export", r.statisticsHandler.ExportExcel)
	}
}
