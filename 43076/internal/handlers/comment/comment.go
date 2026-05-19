package comment

import (
	"strconv"
	"time"

	"ticket-system/internal/database"
	"ticket-system/internal/middleware"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateCommentRequest struct {
	TicketID uint   `json:"ticket_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Type     string `json:"type" binding:"required,comment_type"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"omitempty,comment_type"`
}

func CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	authorID, _ := middleware.GetCurrentUserID(c)

	tx := database.DB.Begin()

	var ticket models.Ticket
	if err := tx.First(&ticket, req.TicketID).Error; err != nil {
		tx.Rollback()
		utils.NotFound(c, "Ticket not found")
		return
	}

	comment := &models.Comment{
		TicketID: req.TicketID,
		AuthorID: authorID,
		Content:  req.Content,
		Type:     req.Type,
	}

	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create comment")
		return
	}

	if ticket.FirstResponseAt == nil && req.Type == models.CommentTypePublic {
		now := time.Now()
		tx.Model(&ticket).Update("first_response_at", &now)

		var slaRecord models.SLARecord
		tx.Where("ticket_id = ?", ticket.ID).First(&slaRecord)
		responseMinutes := int(now.Sub(ticket.CreatedAt).Minutes())
		responseBreached := utils.IsResponseSLABreached(ticket.CreatedAt, now, ticket.Priority)
		tx.Model(&slaRecord).Updates(map[string]interface{}{
			"response_time":    responseMinutes,
			"response_breached": responseBreached,
		})
	}

	tx.Commit()

	database.DB.Preload("Author").Preload("Attachments").First(comment, comment.ID)
	utils.Success(c, comment)
}

func GetComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid comment ID")
		return
	}

	var comment models.Comment
	if err := database.DB.Preload("Author").Preload("Attachments").First(&comment, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Comment not found")
			return
		}
		utils.InternalServerError(c, "Failed to get comment")
		return
	}

	utils.Success(c, comment)
}

func ListComments(c *gin.Context) {
	ticketIDStr := c.Query("ticket_id")
	if ticketIDStr == "" {
		utils.BadRequest(c, "Ticket ID is required")
		return
	}

	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ticket ID")
		return
	}

	var comments []models.Comment
	query := database.DB.Preload("Author").Preload("Attachments").Where("ticket_id = ?", uint(ticketID))

	if commentType := c.Query("type"); commentType != "" {
		query = query.Where("type = ?", commentType)
	}

	if err := query.Order("created_at DESC").Find(&comments).Error; err != nil {
		utils.InternalServerError(c, "Failed to list comments")
		return
	}

	utils.Success(c, comments)
}

func UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid comment ID")
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	currentUserID, _ := middleware.GetCurrentUserID(c)
	userRole, _ := middleware.GetCurrentUserRole(c)

	var comment models.Comment
	if err := database.DB.First(&comment, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Comment not found")
			return
		}
		utils.InternalServerError(c, "Failed to get comment")
		return
	}

	if comment.AuthorID != currentUserID && userRole != models.RoleAdmin {
		utils.Forbidden(c, "You can only update your own comments")
		return
	}

	updates := make(map[string]interface{})
	updates["content"] = req.Content
	if req.Type != "" {
		updates["type"] = req.Type
	}

	if err := database.DB.Model(&comment).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "Failed to update comment")
		return
	}

	database.DB.Preload("Author").Preload("Attachments").First(&comment, uint(id))
	utils.Success(c, comment)
}

func DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid comment ID")
		return
	}

	currentUserID, _ := middleware.GetCurrentUserID(c)
	userRole, _ := middleware.GetCurrentUserRole(c)

	var comment models.Comment
	if err := database.DB.First(&comment, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Comment not found")
			return
		}
		utils.InternalServerError(c, "Failed to get comment")
		return
	}

	if comment.AuthorID != currentUserID && userRole != models.RoleAdmin {
		utils.Forbidden(c, "You can only delete your own comments")
		return
	}

	tx := database.DB.Begin()

	if err := tx.Model(&comment).Association("Attachments").Clear(); err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to delete comment attachments")
		return
	}

	if err := tx.Delete(&comment).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to delete comment")
		return
	}

	tx.Commit()
	utils.Success(c, nil)
}
