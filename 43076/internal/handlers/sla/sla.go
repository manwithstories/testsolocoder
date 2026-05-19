package sla

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

type CheckSLAResponse struct {
	TicketID         uint   `json:"ticket_id"`
	TicketNo       string `json:"ticket_no"`
	Title          string `json:"title"`
	Priority       string `json:"priority"`
	Status         string `json:"status"`
	ResponseTime   int    `json:"response_time_minutes"`
	ResolutionTime int    `json:"resolution_time_minutes"`
	ResponseDeadline string `json:"response_deadline"`
	ResolveDeadline  string `json:"resolve_deadline"`
	ResponseBreached bool `json:"response_breached"`
	ResolveBreached  bool `json:"resolve_breached"`
	IsEscalated    bool   `json:"is_escalated"`
	EscalationCount int  `json:"escalation_count"`
}

func GetSLARecord(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("ticket_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ticket ID")
		return
	}

	var slaRecord models.SLARecord
	if err := database.DB.Where("ticket_id = ?", uint(ticketID)).First(&slaRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "SLA record not found")
			return
		}
		utils.InternalServerError(c, "Failed to get SLA record")
		return
	}

	var ticket models.Ticket
	database.DB.Where("id = ?", uint(ticketID)).Select("ticket_no", "title", "priority", "status", "is_escalated", "response_deadline", "resolve_deadline", "first_response_at", "resolved_at").First(&ticket)

	responseTime := 0
	if ticket.FirstResponseAt != nil {
		responseTime = int(ticket.FirstResponseAt.Sub(ticket.CreatedAt).Minutes())
	}

	resolutionTime := 0
	if ticket.ResolvedAt != nil {
		resolutionTime = int(ticket.ResolvedAt.Sub(ticket.CreatedAt).Minutes())
	}

	responseBreached := false
	if ticket.ResponseDeadline != nil {
		responseBreached = time.Now().After(*ticket.ResponseDeadline) && ticket.FirstResponseAt == nil
	} else if ticket.FirstResponseAt != nil && ticket.ResponseDeadline != nil {
		responseBreached = ticket.FirstResponseAt.After(*ticket.ResponseDeadline)
	}

	resolveBreached := false
	if ticket.ResolveDeadline != nil {
		resolveBreached = time.Now().After(*ticket.ResolveDeadline) && ticket.ResolvedAt == nil
	} else if ticket.ResolvedAt != nil && ticket.ResolveDeadline != nil {
		resolveBreached = ticket.ResolvedAt.After(*ticket.ResolveDeadline)
	}

	response := CheckSLAResponse{
		TicketID:       slaRecord.TicketID,
		TicketNo:       ticket.TicketNo,
		Title:          ticket.Title,
		Priority:       ticket.Priority,
		Status:         ticket.Status,
		ResponseTime:   responseTime,
		ResolutionTime: resolutionTime,
		ResponseDeadline: ticket.ResponseDeadline.Format(time.RFC3339),
		ResolveDeadline:  ticket.ResolveDeadline.Format(time.RFC3339),
		ResponseBreached: responseBreached,
		ResolveBreached:  resolveBreached,
		IsEscalated:    ticket.IsEscalated,
		EscalationCount: slaRecord.EscalationCount,
	}

	utils.Success(c, response)
}

func ListBreachedSLA(c *gin.Context) {
	var tickets []models.Ticket
	now := time.Now()

	query := database.DB.Preload("Assignee").Preload("Customer").Preload("SkillGroup").Where("status NOT IN ?", []string{models.TicketStatusResolved, models.TicketStatusClosed})

	breachType := c.Query("type")
	if breachType == "response" {
		query = query.Where("response_deadline < ? AND first_response_at IS NULL", now)
	} else if breachType == "resolve" {
		query = query.Where("resolve_deadline < ?", now)
	} else {
		query = query.Where("(response_deadline < ? AND first_response_at IS NULL) OR (resolve_deadline < ?)", now, now)
	}

	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
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
	if err := query.Offset(offset).Limit(pageSize).Order("priority DESC, created_at ASC").Find(&tickets).Error; err != nil {
		utils.InternalServerError(c, "Failed to list breached SLA")
		return
	}

	utils.Success(c, gin.H{
		"items":     tickets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func EscalateTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("ticket_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ticket ID")
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

	if ticket.Status == models.TicketStatusClosed || ticket.Status == models.TicketStatusResolved {
		tx.Rollback()
		utils.BadRequest(c, "Cannot escalate resolved or closed ticket")
		return
	}

	now := time.Now()
	updates := make(map[string]interface{})
	updates["is_escalated"] = true
	updates["escalated_at"] = &now
	updates["status"] = models.TicketStatusEscalated

	if err := tx.Model(&ticket).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to escalate ticket")
		return
	}

	var slaRecord models.SLARecord
	tx.Where("ticket_id = ?", ticket.ID).First(&slaRecord)
	tx.Model(&slaRecord).Update("escalation_count", gorm.Expr("escalation_count + ?", 1))

	logEntry := &models.TicketLog{
		TicketID:   ticket.ID,
		OperatorID: operatorID,
		Action:     "escalate",
		OldValue:   ticket.Status,
		NewValue:   models.TicketStatusEscalated,
		Remark:     "SLA超时升级",
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

func GetSLAStats(c *gin.Context) {
	var totalTickets int64
	var breachedTickets int64
	var respondedOnTime int64
	var resolvedOnTime int64

	database.DB.Model(&models.Ticket{}).Count(&totalTickets)

	database.DB.Model(&models.SLARecord{}).Where("response_breached = ? OR resolve_breached = ?", true, true).Count(&breachedTickets)

	database.DB.Model(&models.SLARecord{}).Where("response_breached = ?", false).Where("response_time > ?", 0).Count(&respondedOnTime)

	database.DB.Model(&models.SLARecord{}).Where("resolve_breached = ?", false).Where("resolution_time > ?", 0).Count(&resolvedOnTime)

	var avgResponseTime float64
	var avgResolutionTime float64

	database.DB.Model(&models.SLARecord{}).Select("AVG(response_time)").Where("response_time > ?", 0).Scan(&avgResponseTime)
	database.DB.Model(&models.SLARecord{}).Select("AVG(resolution_time)").Where("resolution_time > ?", 0).Scan(&avgResolutionTime)

	stats := gin.H{
		"total_tickets":        totalTickets,
		"breached_tickets":     breachedTickets,
		"responded_on_time":   respondedOnTime,
		"resolved_on_time":    resolvedOnTime,
		"avg_response_time_minutes":  avgResponseTime,
		"avg_resolution_time_minutes": avgResolutionTime,
		"sla_compliance_rate": 0.0,
	}

	if totalTickets > 0 {
		compliant := totalTickets - breachedTickets
		stats["sla_compliance_rate"] = float64(compliant) / float64(totalTickets) * 100
	}

	utils.Success(c, stats)
}
