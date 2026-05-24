package router

import (
	"github.com/gin-gonic/gin"

	"furniture-platform/internal/handler"
	"furniture-platform/internal/model"
	"furniture-platform/pkg/middleware"
)

type DesignRouter struct {
	designHandler *handler.DesignHandler
	secret        string
}

func NewDesignRouter(designHandler *handler.DesignHandler, secret string) *DesignRouter {
	return &DesignRouter{
		designHandler: designHandler,
		secret:        secret,
	}
}

func (r *DesignRouter) Register(v1 *gin.RouterGroup) {
	auth := v1.Group("/designs", middleware.JWTAuth(r.secret))
	{
		auth.GET("/", r.designHandler.ListProjects)
		auth.GET("/:id", r.designHandler.GetProject)
		auth.GET("/:id/comments", r.designHandler.GetComments)

		designer := auth.Group("", middleware.RoleRequired(model.RoleDesigner))
		{
			designer.POST("/", r.designHandler.CreateProject)
			designer.DELETE("/:id", r.designHandler.DeleteProject)
			designer.POST("/:id/images", r.designHandler.UploadImage)
		}

		designerOrOwner := auth.Group("", middleware.RoleRequired(model.RoleDesigner, model.RoleOwner))
		{
			designerOrOwner.PUT("/:id", r.designHandler.UpdateProject)
			designerOrOwner.POST("/:id/comments", r.designHandler.AddComment)
		}
	}
}
