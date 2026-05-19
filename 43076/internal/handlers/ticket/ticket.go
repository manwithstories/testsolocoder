package ticket

import (
	"encoding/json"
	"strconv"
	"time"

	"ticket-system/internal/database"
	"ticket-system/internal/middleware"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTicketRequest struct {
	Title        string `json:"title" binding:"required,max=200"`
	Description  string `json:"description"`
	Priority     string `json:"priority" binding:"required,priority"`
	Type         string `json:"type" binding:"required,ticket_type"`
	Tags         string `json:"tags" binding:"max=500"`
	CustomerID   uint   `json:"customer_id" binding:"required"`
	AssigneeID   *uint  `json:"assignee_id"`
	SkillGroupID *uint  `json:"skill_group_id"`
}

type UpdateTicketRequest struct {
	Title        string `json:"title" binding:"max=200"`
	Description  string `json:"description"`
	Priority     string `json:"priority" binding:"omitempty,priority"`
	Type         string `json:"type" binding:"omitempty,ticket_type"`
	Tags         string `json:"tags" binding:"max=500"`
	AssigneeID   *uint  `json:"assignee_id"`
	SkillGroupID *uint  `json:"skill_group_id"`
	Remark       string `json:"remark" binding:"max=500"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,ticket_status"`
	Remark string `json:"remark" binding:"max=500"`
}

type AssignTicketRequest struct {
	AssigneeID   uint  `json:"assignee_id" binding:"required"`
	SkillGroupID *uint `json:"skill_group_id"`
	Remark       string `json:"remark" binding:"max=500"`
}

func CreateTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, "Customer not found")
			return
		}
		utils.InternalServerError(c, "Failed to create ticket")
		return
	}

	creatorID, _ := middleware.GetCurrentUserID(c)

	ticketNo, err := utils.GenerateTicketNo()
	if err != nil {
		utils.InternalServerError(c, "Failed to generate ticket number")
		return
	}

	responseDeadline, resolveDeadline := utils.CalculateSLADeadlines(req.Priority, time.Now())

	tx := database.DB.Begin()

	ticket := &models.Ticket{
		TicketNo:         ticketNo,
		Title:            req.Title,
		Description:      req.Description,
		Priority:         req.Priority,
		Status:           models.TicketStatusOpen,
		Type:             req.Type,
		Tags:             req.Tags,
		CustomerID:       req.CustomerID,
		CreatorID:        creatorID,
		AssigneeID:       req.AssigneeID,
		SkillGroupID:     req.SkillGroupID,
		ResponseDeadline: &responseDeadline,
		ResolveDeadline:  &resolveDeadline,
	}

	if req.AssigneeID != nil {
		ticket.Status = models.TicketStatusAssigned
	}

	if err := tx.Create(ticket).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create ticket")
		return
	}

	slaRecord := &models.SLARecord{
		TicketID: ticket.ID,
	}
	if err := tx.Create(slaRecord).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create SLA record")
		return
	}

	logEntry := &models.TicketLog{
		TicketID:   ticket.ID,
		OperatorID: creatorID,
		Action:     "create",
		NewValue:   ticket.Title,
		Remark:     "工单创建",
	}
	if err := tx.Create(logEntry).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create ticket log")
		return
	}

	tx.Commit()

	c.Set("ticket_no", ticketNo)

	database.DB.Preload("Customer").Preload("Creator").Preload("Assignee").Preload("SkillGroup").First(ticket, ticket.ID)
	utils.Success(c, ticket)
}

func GetTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ticket ID")
		return
	}

	var ticket models.Ticket
	if err := database.DB.Preload("Customer").Preload("Creator").Preload("Assignee").
		Preload("SkillGroup").Preload("Logs", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Preload("Operator")
		}).Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Preload("Author").Preload("Attachments")
		}).Preload("Attachments").First(&ticket, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Ticket not found")
			return
		}
		utils.InternalServerError(c, "Failed to get ticket")
		return
	}

	utils.Success(c, ticket)
}

func GetTicketByNo(c *gin.Context) {
	ticketNo := c.Param("ticket_no")

	var ticket models.Ticket
	if err := database.DB.Preload("Customer").Preload("Creator").Preload("Assignee").
		Preload("SkillGroup").Preload("Logs", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Preload("Operator")
		}).Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Preload("Author").Preload("Attachments")
		}).Preload("Attachments").Where("ticket_no = ?", ticketNo).First(&ticket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Ticket not found")
			return
		}
		utils.InternalServerError(c, "Failed to get ticket")
		return
	}

	utils.Success(c, ticket)
}

func ListTickets(c *gin.Context) {
	var tickets []models.Ticket
	query := database.DB.Preload("Customer").Preload("Creator").Preload("Assignee").Preload("SkillGroup")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if ticketType := c.Query("type"); ticketType != "" {
		query = query.Where("type = ?", ticketType)
	}
	if customerID := c.Query("customer_id"); customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	if creatorID := c.Query("creator_id"); creatorID != "" {
		query = query.Where("creator_id = ?", creatorID)
	}
	if assigneeID := c.Query("assignee_id"); assigneeID != "" {
		query = query.Where("assignee_id = ?", assigneeID)
	}
	if skillGroupID := c.Query("skill_group_id"); skillGroupID != "" {
		query = query.Where("skill_group_id = ?", skillGroupID)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ? OR ticket_no LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if pn, err := strconv.Atoi(p); err == nil && pn > 0 {
			page = pn
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if psn, err := strconv.Atoi(ps); err == nil && psn > 0 {
			pageSize = psn
		}
	}

	var total int64
	query.Model(&models.Ticket{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tickets).Error; err != nil {
		utils.InternalServerError(c, "Failed to list tickets")
		return
	}

	utils.Success(c, gin.H{
		"items":     tickets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ticket ID")
		return
	}

	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	operatorID, _ := middleware.GetCurrentUserID(c)

	tx := database.DB.Begin()

	var ticket models.Ticket
	if err := tx.First(&ticket, uint(id)).Error; err != nil {
		tx.Rollback()
		utils.NotFound(c, "Ticket not found")
		return
	}

	if ticket.Status == models.TicketStatusClosed {
		tx.Rollback()
		utils.BadRequest(c, "Cannot update closed ticket")
		return
	}

	oldValues := make(map[string]interface{})
	newValues := make(map[string]interface{})

	updates := make(map[string]interface{})
	if req.Title != "" {
		oldValues["title"] = ticket.Title
		newValues["title"] = req.Title
		updates["title"] = req.Title
	}
	if req.Description != "" {
		oldValues["description"] = ticket.Description
		newValues["description"] = req.Description
		updates["description"] = req.Description
	}
	if req.Priority != "" {
		oldValues["priority"] = ticket.Priority
		newValues["priority"] = req.Priority
		updates["priority"] = req.Priority
	}
	if req.Type != "" {
		oldValues["type"] = ticket.Type
		newValues["type"] = req.Type
		updates["type"] = req.Type
	}
	if req.Tags != "" {
		oldValues["tags"] = ticket.Tags
		newValues["tags"] = req.Tags
		updates["tags"] = req.Tags
	}
	if req.AssigneeID != nil {
		oldValues["assignee_id"] = ticket.AssigneeID
		newValues["assignee_id"] = *req.AssigneeID
		updates["assignee_id"] = *req.AssigneeID
		if ticket.Status == models.TicketStatusOpen {
			updates["status"] = models.TicketStatusAssigned
			oldValues["status"] = models.TicketStatusOpen
			newValues["status"] = models.TicketStatusAssigned
		}
	}
	if req.SkillGroupID != nil {
		oldValues["skill_group_id"] = ticket.SkillGroupID
		newValues["skill_group_id"] = *req.SkillGroupID
		updates["skill_group_id"] = *req.SkillGroupID
	}

	if len(updates) > 0 {
		if err := tx.Model(&ticket).Updates(updates).Error; err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to update ticket")
			return
		}

		oldJSON, _ := json.Marshal(oldValues)
		newJSON, _ := json.Marshal(newValues)

		logEntry := &models.TicketLog{
			TicketID:   ticket.ID,
			OperatorID: operatorID,
			Action:     "update",
			OldValue:   string(oldJSON),
			NewValue:   string(newJSON),
			Remark:     req.Remark,
		}
		if err := tx.Create(logEntry).Error; err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to create ticket log")
			return
		}
	}

	tx.Commit()

	database.DB.Preload("Customer").Preload("Creator").Preload("Assignee").Preload("SkillGroup").First(&ticket, uint(id))
	utils.Success(c, ticket)
}

func UpdateTicketStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ticket ID")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	operatorID, _ := middleware.GetCurrentUserID(c)

	tx := database.DB.Begin()

	var ticket models.Ticket
	if err := tx.First(&ticket, uint(id)).Error; err != nil {
		tx.Rollback()
		utils.NotFound(c, "Ticket not found")
		return
	}

	if ticket.Status == models.TicketStatusClosed && req.Status != models.TicketStatusClosed {
		tx.Rollback()
		utils.BadRequest(c, "Cannot update closed ticket")
		return
	}

	if !models.IsValidStatusTransition(ticket.Status, req.Status) {
		tx.Rollback()
		utils.BadRequest(c, "Invalid status transition from "+ticket.Status+" to "+req.Status)
		return
	}

	now := time.Now()
	updates := make(map[string]interface{})
	updates["status"] = req.Status

	if req.Status == models.TicketStatusResolved && ticket.ResolvedAt == nil {
		updates["resolved_at"] = &now
	}
	if req.Status == models.TicketStatusClosed {
		updates["closed_at"] = &now
	}

	if err := tx.Model(&ticket).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to update ticket status")
		return
	}

	if req.Status == models.TicketStatusResolved && ticket.ResolvedAt == nil {
		var slaRecord models.SLARecord
		tx.Where("ticket_id = ?", ticket.ID).First(&slaRecord)
		resolutionMinutes := int(now.Sub(ticket.CreatedAt).Minutes())
		resolveBreached := utils.IsResolveSLABreached(ticket.CreatedAt, now, ticket.Priority)
		tx.Model(&slaRecord).Updates(map[string]interface{}{
			"resolution_time":  resolutionMinutes,
			"resolve_breached": resolveBreached,
		})
	}

	logEntry := &models.TicketLog{
		TicketID:   ticket.ID,
		OperatorID: operatorID,
		Action:     "status_change",
		OldValue:   ticket.Status,
		NewValue:   req.Status,
		Remark:     req.Remark,
	}
	if err := tx.Create(logEntry).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create ticket log")
		return
	}

	tx.Commit()

	database.DB.Preload("Customer").Preload("Creator").Preload("Assignee").Preload("SkillGroup").First(&ticket, uint(id))
	utils.Success(c, ticket)
}

func AssignTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ticket ID")
		return
	}

	var req AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	operatorID, _ := middleware.GetCurrentUserID(c)

	tx := database.DB.Begin()

	var ticket models.Ticket
	if err := tx.First(&ticket, uint(id)).Error; err != nil {
		tx.Rollback()
		utils.NotFound(c, "Ticket not found")
		return
	}

	if ticket.Status == models.TicketStatusClosed {
		tx.Rollback()
		utils.BadRequest(c, "Cannot assign closed ticket")
		return
	}

	var assignee models.User
	if err := tx.First(&assignee, req.AssigneeID).Error; err != nil {
		tx.Rollback()
		utils.BadRequest(c, "Assignee not found")
		return
	}

	updates := make(map[string]interface{})
	updates["assignee_id"] = req.AssigneeID
	if req.SkillGroupID != nil {
		updates["skill_group_id"] = *req.SkillGroupID
	}
	if ticket.Status == models.TicketStatusOpen {
		updates["status"] = models.TicketStatusAssigned
	}

	if err := tx.Model(&ticket).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to assign ticket")
		return
	}

	oldAssignee := "未分配"
	if ticket.AssigneeID != nil {
		oldAssignee = strconv.FormatUint(uint64(*ticket.AssigneeID), 10)
	}

	logEntry := &models.TicketLog{
		TicketID:   ticket.ID,
		OperatorID: operatorID,
		Action:     "assign",
		OldValue:   oldAssignee,
		NewValue:   strconv.FormatUint(uint64(req.AssigneeID), 10),
		Remark:     req.Remark,
	}
	if err := tx.Create(logEntry).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create ticket log")
		return
	}

	tx.Commit()

	database.DB.Preload("Customer").Preload("Creator").Preload("Assignee").Preload("SkillGroup").First(&ticket, uint(id))
	utils.Success(c, ticket)
}

func GetTicketLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ticket ID")
		return
	}

	var logs []models.TicketLog
	if err := database.DB.Where("ticket_id = ?", uint(id)).Preload("Operator").
		Order("created_at DESC").Find(&logs).Error; err != nil {
		utils.InternalServerError(c, "Failed to get ticket logs")
		return
	}

	utils.Success(c, logs)
}
