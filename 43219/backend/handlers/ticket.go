package handlers

import (
	"time"

	"housekeeping/config"
	"housekeeping/database"
	"housekeeping/models"
	"housekeeping/utils"

	"github.com/gin-gonic/gin"
)

type TicketInput struct {
	OrderID  *uint            `json:"order_id,omitempty"`
	StaffID  *uint            `json:"staff_id,omitempty"`
	Type     models.TicketType  `json:"type" binding:"required"`
	Title    string           `json:"title" binding:"required"`
	Content  string           `json:"content" binding:"required"`
}

func CreateTicket(c *gin.Context) {
	uid, _ := c.Get("uid")
	var in TicketInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	tk := models.Ticket{
		OrderID:      in.OrderID,
		CustomerID: uid.(uint),
		StaffID:    in.StaffID,
		Type:       in.Type,
		Title:      in.Title,
		Content:    in.Content,
		Status:     models.TicketOpen,
		LastActionAt: timeNow(),
	}
	if err := database.DB.Create(&tk).Error; err != nil {
		utils.ServerError(c, "create failed")
		return
	}
	utils.OK(c, tk)
}

func ListTickets(c *gin.Context) {
	uid, _ := c.Get("uid")
	role, _ := c.Get("role")
	var list []models.Ticket
	q := database.DB
	switch role.(string) {
	case string(models.RoleCustomer):
		q = q.Where("customer_id = ?", uid)
	case string(models.RoleStaff):
		q = q.Where("staff_id = ?", uid)
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		utils.ServerError(c, "query failed")
		return
	}
	utils.OK(c, list)
}

func AssignTicket(c *gin.Context) {
	id := c.Param("id")
	var tk models.Ticket
	if err := database.DB.First(&tk, id).Error; err != nil {
		utils.NotFound(c, "ticket not found")
		return
	}
	var body struct {
		AgentID uint `json:"agent_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	updates := map[string]interface{}{
		"agent_id":       body.AgentID,
		"status":         models.TicketAssigned,
		"last_action_at": timeNow(),
	}
	if err := database.DB.Model(&tk).Updates(updates).Error; err != nil {
		utils.ServerError(c, "assign failed")
		return
	}
	utils.OK(c, "ok")
}

func ResolveTicket(c *gin.Context) {
	id := c.Param("id")
	var tk models.Ticket
	if err := database.DB.First(&tk, id).Error; err != nil {
		utils.NotFound(c, "ticket not found")
		return
	}
	var body struct {
		Result string `json:"result"`
		Refund bool   `json:"refund"`
	}
	c.ShouldBindJSON(&body)
	updates := map[string]interface{}{
		"status":         models.TicketResolved,
		"last_action_at": timeNow(),
	}
	if body.Result != "" {
		updates["content"] = tk.Content + "\n[处理结果] " + body.Result
	}
	if err := database.DB.Model(&tk).Updates(updates).Error; err != nil {
		utils.ServerError(c, "resolve failed")
		return
	}
	if body.Refund && tk.OrderID != nil {
		var order models.Order
		database.DB.First(&order, tk.OrderID)
		order.Status = models.OrderRefunded
		database.DB.Save(&order)
	}
	utils.OK(c, "ok")
}

func CloseTicket(c *gin.Context) {
	id := c.Param("id")
	database.DB.Model(&models.Ticket{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         models.TicketClosed,
		"last_action_at": timeNow(),
	})
	utils.OK(c, "ok")
}

func EscalateTickets() {
	hours := config.C.Ticket.EscalateHours
	if hours <= 0 {
		return
	}
	cutoff := timeNow().Add(-time.Duration(hours) * time.Hour)
	database.DB.Model(&models.Ticket{}).
		Where("status IN ? AND last_action_at < ? AND escalated = ?",
			[]models.TicketStatus{models.TicketOpen, models.TicketAssigned, models.TicketPending},
			cutoff, false,
		).Updates(map[string]interface{}{
		"status":    models.TicketEscalate,
		"escalated": true,
	})
}
