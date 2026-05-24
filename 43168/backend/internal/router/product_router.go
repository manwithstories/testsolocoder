package router

import (
	"github.com/gin-gonic/gin"

	"furniture-platform/internal/handler"
	"furniture-platform/internal/model"
	"furniture-platform/pkg/middleware"
)

// ProductRouter 产品路由依赖
type ProductRouter struct {
	productHandler *handler.ProductHandler
	secret         string
}

// NewProductRouter 创建产品路由
func NewProductRouter(productHandler *handler.ProductHandler, secret string) *ProductRouter {
	return &ProductRouter{
		productHandler: productHandler,
		secret:         secret,
	}
}

// Register 注册产品路由到已有的 v1 路由组
func (r *ProductRouter) Register(v1 *gin.RouterGroup) {
	// 可匿名访问的产品接口
	publicProducts := v1.Group("/products")
	{
		publicProducts.GET("/hot", r.productHandler.GetHotProducts)
		publicProducts.GET("/", r.productHandler.ListProducts)
		publicProducts.GET("/:id", r.productHandler.GetProduct)
	}

	// 需要登录的产品接口
	authProducts := v1.Group("/products", middleware.JWTAuth(r.secret))
	{
		// 仅厂商角色可操作
		manufacturerGroup := authProducts.Group("", middleware.RoleRequired(model.RoleManufacturer))
		{
			manufacturerGroup.POST("/", r.productHandler.CreateProduct)
			manufacturerGroup.PUT("/:id", r.productHandler.UpdateProduct)
			manufacturerGroup.DELETE("/:id", r.productHandler.DeleteProduct)
		}
	}
}
