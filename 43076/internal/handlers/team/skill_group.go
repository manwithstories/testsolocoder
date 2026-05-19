package team

import (
	"strconv"

	"ticket-system/internal/database"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateSkillGroupRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	UserIDs     []uint `json:"user_ids"`
}

type UpdateSkillGroupRequest struct {
	Name        string `json:"name" binding:"max=100"`
	Description string `json:"description" binding:"max=500"`
	UserIDs     []uint `json:"user_ids"`
}

func CreateSkillGroup(c *gin.Context) {
	var req CreateSkillGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	tx := database.DB.Begin()

	group := &models.SkillGroup{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := tx.Create(group).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create skill group")
		return
	}

	if len(req.UserIDs) > 0 {
		var users []models.User
		tx.Find(&users, req.UserIDs)
		if err := tx.Model(group).Association("Users").Append(&users); err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to add users to skill group")
			return
		}
	}

	tx.Commit()

	database.DB.Preload("Users").First(group, group.ID)
	utils.Success(c, group)
}

func GetSkillGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid skill group ID")
		return
	}

	var group models.SkillGroup
	if err := database.DB.Preload("Users").First(&group, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Skill group not found")
			return
		}
		utils.InternalServerError(c, "Failed to get skill group")
		return
	}

	utils.Success(c, group)
}

func ListSkillGroups(c *gin.Context) {
	var groups []models.SkillGroup
	if err := database.DB.Find(&groups).Error; err != nil {
		utils.InternalServerError(c, "Failed to list skill groups")
		return
	}

	utils.Success(c, groups)
}

func UpdateSkillGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid skill group ID")
		return
	}

	var req UpdateSkillGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	tx := database.DB.Begin()

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) > 0 {
		if err := tx.Model(&models.SkillGroup{}).Where("id = ?", uint(id)).Updates(updates).Error; err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to update skill group")
			return
		}
	}

	if req.UserIDs != nil {
		var group models.SkillGroup
		tx.First(&group, uint(id))
		tx.Model(&group).Association("Users").Clear()
		var users []models.User
		tx.Find(&users, req.UserIDs)
		if err := tx.Model(&group).Association("Users").Append(&users); err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to update skill group users")
			return
		}
	}

	tx.Commit()

	var group models.SkillGroup
	database.DB.Preload("Users").First(&group, uint(id))
	utils.Success(c, group)
}

func DeleteSkillGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid skill group ID")
		return
	}

	tx := database.DB.Begin()

	var group models.SkillGroup
	if err := tx.First(&group, uint(id)).Error; err != nil {
		tx.Rollback()
		utils.NotFound(c, "Skill group not found")
		return
	}

	tx.Model(&group).Association("Users").Clear()

	if err := tx.Delete(&group).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to delete skill group")
		return
	}

	tx.Commit()
	utils.Success(c, nil)
}
