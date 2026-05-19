package attachment

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"ticket-system/internal/config"
	"ticket-system/internal/database"
	"ticket-system/internal/middleware"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UploadAttachment(c *gin.Context) {
	ticketIDStr := c.PostForm("ticket_id")
	commentIDStr := c.PostForm("comment_id")

	if ticketIDStr == "" && commentIDStr == "" {
		utils.BadRequest(c, "Ticket ID or Comment ID is required")
		return
	}

	uploaderID, _ := middleware.GetCurrentUserID(c)
	uploadCfg := config.AppConfig.Upload

	if err := os.MkdirAll(uploadCfg.StoragePath, 0755); err != nil {
		utils.InternalServerError(c, "Failed to create upload directory")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "Failed to get file: "+err.Error())
		return
	}
	defer file.Close()

	maxSize := int64(uploadCfg.MaxSizeMB) * 1024 * 1024
	if header.Size > maxSize {
		utils.BadRequest(c, fmt.Sprintf("File size exceeds maximum allowed size of %d MB", uploadCfg.MaxSizeMB))
		return
	}

	contentType := header.Header.Get("Content-Type")
	allowed := false
	for _, t := range uploadCfg.AllowedTypes {
		if contentType == t {
			allowed = true
			break
		}
	}
	if !allowed {
		utils.BadRequest(c, "File type not allowed: "+contentType)
		return
	}

	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uploadCfg.StoragePath, dateDir)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		utils.InternalServerError(c, "Failed to create date directory")
		return
	}

	filePath := filepath.Join(fullDir, newFileName)
	out, err := os.Create(filePath)
	if err != nil {
		utils.InternalServerError(c, "Failed to create file: "+err.Error())
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		utils.InternalServerError(c, "Failed to save file: "+err.Error())
		return
	}

	tx := database.DB.Begin()

	var ticketID *uint
	var commentID *uint

	if ticketIDStr != "" {
		tid, err := strconv.ParseUint(ticketIDStr, 10, 32)
		if err != nil {
			tx.Rollback()
			utils.BadRequest(c, "Invalid ticket ID")
			return
		}
		tidUint := uint(tid)
		ticketID = &tidUint
	}

	if commentIDStr != "" {
		cid, err := strconv.ParseUint(commentIDStr, 10, 32)
		if err != nil {
			tx.Rollback()
			utils.BadRequest(c, "Invalid comment ID")
			return
		}
		cidUint := uint(cid)
		commentID = &cidUint
	}

	attachment := &models.Attachment{
		TicketID:    ticketID,
		CommentID:   commentID,
		FileName:    header.Filename,
		FilePath:    filePath,
		FileSize:    header.Size,
		ContentType: contentType,
		UploaderID:  uploaderID,
	}

	if err := tx.Create(attachment).Error; err != nil {
		tx.Rollback()
		os.Remove(filePath)
		utils.InternalServerError(c, "Failed to save attachment record")
		return
	}

	tx.Commit()

	database.DB.Preload("Uploader").First(attachment, attachment.ID)
	utils.Success(c, attachment)
}

func GetAttachment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid attachment ID")
		return
	}

	var attachment models.Attachment
	if err := database.DB.First(&attachment, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Attachment not found")
			return
		}
		utils.InternalServerError(c, "Failed to get attachment")
		return
	}

	if _, err := os.Stat(attachment.FilePath); os.IsNotExist(err) {
		utils.NotFound(c, "File not found on disk")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachment.FileName))
	c.Header("Content-Type", attachment.ContentType)
	c.File(attachment.FilePath)
}

func ListAttachments(c *gin.Context) {
	var attachments []models.Attachment
	query := database.DB.Preload("Uploader")

	if ticketIDStr := c.Query("ticket_id"); ticketIDStr != "" {
		ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
		if err == nil {
			query = query.Where("ticket_id = ?", uint(ticketID))
		}
	}

	if commentIDStr := c.Query("comment_id"); commentIDStr != "" {
		commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
		if err == nil {
			query = query.Where("comment_id = ?", uint(commentID))
		}
	}

	if uploaderIDStr := c.Query("uploader_id"); uploaderIDStr != "" {
		uploaderID, err := strconv.ParseUint(uploaderIDStr, 10, 32)
		if err == nil {
			query = query.Where("uploader_id = ?", uint(uploaderID))
		}
	}

	if err := query.Order("created_at DESC").Find(&attachments).Error; err != nil {
		utils.InternalServerError(c, "Failed to list attachments")
		return
	}

	utils.Success(c, attachments)
}

func DeleteAttachment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid attachment ID")
		return
	}

	currentUserID, _ := middleware.GetCurrentUserID(c)
	userRole, _ := middleware.GetCurrentUserRole(c)

	var attachment models.Attachment
	if err := database.DB.First(&attachment, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Attachment not found")
			return
		}
		utils.InternalServerError(c, "Failed to get attachment")
		return
	}

	if attachment.UploaderID != currentUserID && userRole != models.RoleAdmin {
		utils.Forbidden(c, "You can only delete your own attachments")
		return
	}

	tx := database.DB.Begin()

	filePath := attachment.FilePath

	if err := tx.Delete(&attachment).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to delete attachment record")
		return
	}

	tx.Commit()

	if filePath != "" {
		os.Remove(filePath)
	}

	utils.Success(c, nil)
}
