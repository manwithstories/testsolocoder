package handlers

import (
	"strconv"
	"time"

	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"housekeeping-platform/utils"

	"github.com/gin-gonic/gin"
)

func GetMessageList(c *gin.Context) {
	userID := c.GetUint("user_id")
	msgType := c.Query("type")
	isRead := c.Query("is_read")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := config.DB.Model(&models.Message{}).Where("user_id = ?", userID)

	if msgType != "" {
		query = query.Where("type = ?", msgType)
	}
	if isRead != "" {
		read := isRead == "true"
		query = query.Where("is_read = ?", read)
	}

	var total int64
	query.Count(&total)

	var messages []models.Message
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&messages)

	var unreadCount int64
	config.DB.Model(&models.Message{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	utils.Success(c, gin.H{
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"unread_count": unreadCount,
		"list":        messages,
	})
}

func ReadMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var message models.Message
	if result := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&message); result.Error != nil {
		utils.NotFound(c, "消息不存在")
		return
	}

	now := time.Now()
	message.IsRead = true
	message.ReadAt = &now

	config.DB.Save(&message)

	utils.Success(c, nil)
}

func ReadAllMessages(c *gin.Context) {
	userID := c.GetUint("user_id")

	now := time.Now()
	config.DB.Model(&models.Message{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})

	utils.Success(c, nil)
}

func DeleteMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if result := config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Message{}); result.Error != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.Success(c, nil)
}

func GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var count int64
	config.DB.Model(&models.Message{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count)

	utils.Success(c, gin.H{"unread_count": count})
}
