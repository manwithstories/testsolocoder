package router

import (
	"furniture-platform/internal/handler"
	"furniture-platform/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type ProductionRouter struct {
	handler *handler.ProductionHandler
	secret  string
}

func NewProductionRouter(h *handler.ProductionHandler, secret string) *ProductionRouter {
	return &ProductionRouter{handler: h, secret: secret}
}

func (r *ProductionRouter) Register(v1 *gin.RouterGroup) {
	productions := v1.Group("/productions", middleware.JWTAuth(r.secret))
	{
		productions.GET("/", r.handler.ListProductions)
		productions.GET("/:id", r.handler.GetProduction)
		productions.POST("/", r.handler.CreateProduction)
		productions.PUT("/:id/status", middleware.RoleRequired("manufacturer"), r.handler.UpdateStatus)
	}
}
