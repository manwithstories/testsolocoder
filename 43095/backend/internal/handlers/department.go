package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	departmentService *services.DepartmentService
}

func NewDepartmentHandler() *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: services.NewDepartmentService(),
	}
}

func RegisterDepartmentRoutes(r *gin.RouterGroup) {
	handler := NewDepartmentHandler()

	departments := r.Group("/departments")
	{
		departments.GET("", handler.List)
		departments.GET("/:id", handler.GetByID)

		admin := departments.Group("")
		admin.Use(middleware.Auth(), middleware.AdminRequired())
		{
			admin.POST("", handler.Create)
			admin.PUT("/:id", handler.Update)
			admin.DELETE("/:id", handler.Delete)
		}
	}
}

func (h *DepartmentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	departments, total, err := h.departmentService.GetList(page, pageSize, keyword)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.SuccessWithPagination(c, departments, total, page, pageSize)
}

func (h *DepartmentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的科室ID")
		return
	}

	department, err := h.departmentService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "科室不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, department)
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	var req services.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	department, err := h.departmentService.Create(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, department)
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的科室ID")
		return
	}

	var req services.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	department, err := h.departmentService.Update(uint(id), &req)
	if err != nil {
		if err.Error() == "科室不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, department)
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的科室ID")
		return
	}

	if err := h.departmentService.Delete(uint(id)); err != nil {
		if err.Error() == "科室不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}
