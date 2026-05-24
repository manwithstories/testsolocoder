package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tutoring-platform/database"
	"tutoring-platform/models"
)

type SendMessageRequest struct {
	ReceiverID uuid.UUID           `json:"receiverId" binding:"required"`
	Type       models.MessageType  `json:"type"`
	Content    string              `json:"content" binding:"required"`
	BookingID  *uuid.UUID          `json:"bookingId"`
	Files      []MessageFileRequest `json:"files"`
}

type MessageFileRequest struct {
	FileName string `json:"fileName" binding:"required"`
	FileURL  string `json:"fileUrl" binding:"required"`
	FileSize int64  `json:"fileSize"`
	FileType string `json:"fileType"`
}

func SendMessage(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type == "" {
		req.Type = models.MessageTypeText
	}

	message := models.Message{
		SenderID:   userID.(uuid.UUID),
		ReceiverID: req.ReceiverID,
		BookingID:  req.BookingID,
		Type:       req.Type,
		Content:    req.Content,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	for _, file := range req.Files {
		messageFile := models.MessageFile{
			MessageID: message.ID,
			FileName:  file.FileName,
			FileURL:   file.FileURL,
			FileSize:  file.FileSize,
			FileType:  file.FileType,
		}
		if err := tx.Create(&messageFile).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to attach file"})
			return
		}
	}

	tx.Commit()

	createNotification(c, req.ReceiverID, models.NotificationTypeNewMessage, "New Message", "You have received a new message")

	c.JSON(http.StatusCreated, gin.H{"message": "Message sent", "data": message})
}

func GetMessages(c *gin.Context) {
	userID, _ := c.Get("userId")
	otherUserID := c.Query("userId")
	bookingID := c.Query("bookingId")

	var messages []models.Message
	query := database.DB.Preload("Sender").Preload("Receiver").Preload("Files").Preload("Booking")

	query = query.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, otherUserID, otherUserID, userID)

	if bookingID != "" {
		query = query.Where("booking_id = ?", bookingID)
	}

	if err := query.Order("created_at ASC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func GetConversationList(c *gin.Context) {
	userID, _ := c.Get("userId")

	type ConversationResult struct {
		UserID      uuid.UUID `json:"userId"`
		FirstName   string    `json:"firstName"`
		LastName    string    `json:"lastName"`
		AvatarURL   string    `json:"avatarUrl"`
		LastMessage string    `json:"lastMessage"`
		LastTime    string    `json:"lastTime"`
		UnreadCount int64     `json:"unreadCount"`
	}

	var conversations []ConversationResult

	rows, err := database.DB.Raw(`
		SELECT 
			CASE WHEN m.sender_id = ? THEN m.receiver_id ELSE m.sender_id END as user_id,
			u.first_name,
			u.last_name,
			u.avatar_url,
			m.content as last_message,
			m.created_at as last_time,
			(SELECT COUNT(*) FROM messages WHERE receiver_id = ? AND sender_id = CASE WHEN m.sender_id = ? THEN m.receiver_id ELSE m.sender_id END AND is_read = false) as unread_count
		FROM messages m
		JOIN users u ON u.id = CASE WHEN m.sender_id = ? THEN m.receiver_id ELSE m.sender_id END
		WHERE m.sender_id = ? OR m.receiver_id = ?
		ORDER BY m.created_at DESC
	`, userID, userID, userID, userID, userID, userID).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var conv ConversationResult
		rows.Scan(&conv.UserID, &conv.FirstName, &conv.LastName, &conv.AvatarURL, &conv.LastMessage, &conv.LastTime, &conv.UnreadCount)
		conversations = append(conversations, conv)
	}

	c.JSON(http.StatusOK, conversations)
}

func MarkMessagesAsRead(c *gin.Context) {
	userID, _ := c.Get("userId")
	senderID := c.Query("senderId")

	if err := database.DB.Model(&models.Message{}).
		Where("receiver_id = ? AND sender_id = ? AND is_read = ?", userID, senderID, false).
		Updates(map[string]interface{}{"is_read": true}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark messages as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Messages marked as read"})
}

func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userId")

	var message models.Message
	if err := database.DB.Where("id = ?", id).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	if message.SenderID != userID && message.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this message"})
		return
	}

	uid := userID.(uuid.UUID)
	message.IsDeleted = true
	message.DeletedBy = &uid

	if err := database.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

func GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("userId")

	var count int64
	if err := database.DB.Model(&models.Message{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch unread count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unreadCount": count})
}
