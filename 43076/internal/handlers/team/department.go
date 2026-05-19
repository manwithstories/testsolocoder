package team

import (
	"strconv"

	"ticket-system/internal/database"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
}

type UpdateDepartmentRequest struct {
	Name        string `json:"name" binding:"max=100"`
	Description string `json:"description" binding:"max=500"`
}

func CreateDepartment(c *gin.Context) {
	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	dept := &models.Department{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := database.DB.Create(dept).Error; err != nil {
		utils.InternalServerError(c, "Failed to create department")
		return
	}

	utils.Success(c, dept)
}

func GetDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid department ID")
		return
	}

	var dept models.Department
	if err := database.DB.Preload("Users").First(&dept, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Department not found")
			return
		}
		utils.InternalServerError(c, "Failed to get department")
		return
	}

	utils.Success(c, dept)
}

func ListDepartments(c *gin.Context) {
	var depts []models.Department
	if err := database.DB.Find(&depts).Error; err != nil {
		utils.InternalServerError(c, "Failed to list departments")
		return
	}

	utils.Success(c, depts)
}

func UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid department ID")
		return
	}

	var req UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := database.DB.Model(&models.Department{}).Where("id = ?", uint(id)).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "Failed to update department")
		return
	}

	var dept models.Department
	database.DB.First(&dept, uint(id))
	utils.Success(c, dept)
}

func DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid department ID")
		return
	}

	var count int64
	database.DB.Model(&models.User{}).Where("department_id = ?", uint(id)).Count(&count)
	if count > 0 {
		utils.BadRequest(c, "Cannot delete department with users")
		return
	}

	if err := database.DB.Delete(&models.Department{}, uint(id)).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete department")
		return
	}

	utils.Success(c, nil)
}
