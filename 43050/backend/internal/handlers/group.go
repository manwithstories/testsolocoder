package handlers

import (
	"net/http"
	"splitwise-clone/internal/database"
	"splitwise-clone/internal/models"
	"splitwise-clone/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type JoinGroupRequest struct {
	InviteCode string `json:"inviteCode" binding:"required"`
}

func CreateGroup(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inviteCode := utils.GenerateInviteCode(8)

	group := models.Group{
		Name:        req.Name,
		Description: req.Description,
		CreatorID:   userID,
		InviteCode:  inviteCode,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&group).Error; err != nil {
			return err
		}

		member := models.GroupMember{
			GroupID:  group.ID,
			UserID:   userID,
			JoinedAt: time.Now(),
			IsActive: true,
		}
		return tx.Create(&member).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	database.DB.Preload("Creator").First(&group, group.ID)
	c.JSON(http.StatusCreated, group)
}

func GetUserGroups(c *gin.Context) {
	userID := c.GetUint("userID")

	var groups []models.Group
	database.DB.Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ? AND group_members.is_active = ?", userID, true).
		Preload("Creator").
		Preload("Members").
		Find(&groups)

	c.JSON(http.StatusOK, groups)
}

func GetGroupByID(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	var group models.Group
	if err := database.DB.Preload("Creator").Preload("Members").First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	c.JSON(http.StatusOK, group)
}

func GetGroupMembers(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	var members []models.GroupMember
	database.DB.Where("group_id = ?", groupID).
		Preload("User").
		Find(&members)

	c.JSON(http.StatusOK, members)
}

func JoinGroup(c *gin.Context) {
	userID := c.GetUint("userID")

	var req JoinGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var group models.Group
	if err := database.DB.Where("invite_code = ?", req.InviteCode).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid invite code"})
		return
	}

	var existingMember models.GroupMember
	err := database.DB.Where("group_id = ? AND user_id = ?", group.ID, userID).First(&existingMember).Error
	if err == nil {
		if existingMember.IsActive {
			c.JSON(http.StatusConflict, gin.H{"error": "Already a member of this group"})
			return
		}
		existingMember.IsActive = true
		existingMember.JoinedAt = time.Now()
		database.DB.Save(&existingMember)
		c.JSON(http.StatusOK, group)
		return
	}

	member := models.GroupMember{
		GroupID:  group.ID,
		UserID:   userID,
		JoinedAt: time.Now(),
		IsActive: true,
	}

	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join group"})
		return
	}

	c.JSON(http.StatusOK, group)
}

func LeaveGroup(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	var group models.Group
	if err := database.DB.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	if group.CreatorID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Creator cannot leave the group"})
		return
	}

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not a member of this group"})
		return
	}

	member.IsActive = false
	database.DB.Save(&member)

	c.JSON(http.StatusOK, gin.H{"message": "Left group successfully"})
}
