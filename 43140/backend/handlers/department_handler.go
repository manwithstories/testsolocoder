package handlers

import (
	"strconv"

	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DepartmentHandler struct {
	DB *gorm.DB
}

func NewDepartmentHandler(db *gorm.DB) *DepartmentHandler {
	return &DepartmentHandler{DB: db}
}

func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	dept := models.Department{
		CompanyID: company.ID,
		Name:      req.Name,
	}

	if err := h.DB.Create(&dept).Error; err != nil {
		utils.InternalError(c, "Failed to create department")
		return
	}

	utils.Success(c, dept)
}

func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var dept models.Department
	if err := h.DB.Where("id = ? AND company_id = ?", id, company.ID).First(&dept).Error; err != nil {
		utils.NotFound(c, "Department not found")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	h.DB.Model(&dept).Update("name", req.Name)
	utils.Success(c, dept)
}

func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	result := h.DB.Where("id = ? AND company_id = ?", id, company.ID).Delete(&models.Department{})
	if result.RowsAffected == 0 {
		utils.NotFound(c, "Department not found")
		return
	}

	utils.Success(c, nil)
}

func (h *DepartmentHandler) ListDepartments(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var departments []models.Department
	h.DB.Preload("Jobs").Where("company_id = ?", company.ID).Find(&departments)

	utils.Success(c, departments)
}

func (h *DepartmentHandler) GetPublicDepartments(c *gin.Context) {
	companyID := c.Param("company_id")

	var departments []models.Department
	h.DB.Where("company_id = ?", companyID).Find(&departments)

	utils.Success(c, departments)
}

func parseInt(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil && v > 0 {
		return v
	}
	return def
}
