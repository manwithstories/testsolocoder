package api

import (
	"gym-management/internal/middleware"
	"gym-management/internal/models"
	"gym-management/internal/pkg/database"
	"gym-management/internal/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CoachHandler struct{}

func NewCoachHandler() *CoachHandler {
	return &CoachHandler{}
}

func (h *CoachHandler) RegisterRoutes(r *gin.RouterGroup) {
	coach := r.Group("/coaches")
	coach.Use(middleware.JWTAuth())
	{
		coach.POST("/", h.Create)
		coach.GET("/", h.List)
		coach.GET("/:id", h.GetByID)
		coach.PUT("/:id", h.Update)
		coach.DELETE("/:id", h.Delete)
		coach.PATCH("/:id/status", h.UpdateStatus)
	}
}

func (h *CoachHandler) Create(c *gin.Context) {
	var coach models.Coach
	if err := c.ShouldBindJSON(&coach); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	coach.Status = 1
	if err := database.GetDB().Create(&coach).Error; err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, coach)
}

func (h *CoachHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPageInfo(c)
	keyword := c.Query("keyword")

	var coaches []models.Coach
	var total int64

	query := database.GetDB().Model(&models.Coach{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR specialty LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Find(&coaches)

	utils.SuccessWithPagination(c, coaches, page, pageSize, total)
}

func (h *CoachHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var coach models.Coach
	if err := database.GetDB().First(&coach, id).Error; err != nil {
		utils.NotFound(c, "教练不存在")
		return
	}

	utils.Success(c, coach)
}

func (h *CoachHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "updated_at")

	if err := database.GetDB().Model(&models.Coach{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *CoachHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	if err := database.GetDB().Delete(&models.Coach{}, id).Error; err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *CoachHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=1 2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := database.GetDB().Model(&models.Coach{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}
