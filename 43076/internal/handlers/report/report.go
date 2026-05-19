package report

import (
	"strconv"
	"time"

	"ticket-system/internal/database"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
)

type DailyStats struct {
	Date          string `json:"date"`
	NewTickets    int64  `json:"new_tickets"`
	ResolvedTickets int64 `json:"resolved_tickets"`
	ClosedTickets int64  `json:"closed_tickets"`
	AvgResponseTime float64 `json:"avg_response_time_minutes"`
	AvgResolutionTime float64 `json:"avg_resolution_time_minutes"`
}

type WeeklyStats struct {
	Week          string `json:"week"`
	NewTickets    int64  `json:"new_tickets"`
	ResolvedTickets int64 `json:"resolved_tickets"`
	ClosedTickets int64  `json:"closed_tickets"`
	AvgResponseTime float64 `json:"avg_response_time_minutes"`
	AvgResolutionTime float64 `json:"avg_resolution_time_minutes"`
}

type MonthlyStats struct {
	Month         string `json:"month"`
	NewTickets    int64  `json:"new_tickets"`
	ResolvedTickets int64 `json:"resolved_tickets"`
	ClosedTickets int64  `json:"closed_tickets"`
	AvgResponseTime float64 `json:"avg_response_time_minutes"`
	AvgResolutionTime float64 `json:"avg_resolution_time_minutes"`
}

type PriorityDistribution struct {
	Priority string `json:"priority"`
	Count    int64  `json:"count"`
}

type StatusDistribution struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type TypeDistribution struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type AgentPerformance struct {
	UserID        uint   `json:"user_id"`
	Username      string `json:"username"`
	RealName      string `json:"real_name"`
	TicketsHandled int64 `json:"tickets_handled"`
	TicketsResolved int64 `json:"tickets_resolved"`
	AvgResponseTime float64 `json:"avg_response_time_minutes"`
	AvgResolutionTime float64 `json:"avg_resolution_time_minutes"`
}

func GetDailyStats(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 7
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days+1)

	stats := make([]DailyStats, 0)

	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		nextDate := date.AddDate(0, 0, 1)

		var stat DailyStats
		stat.Date = dateStr

		database.DB.Model(&models.Ticket{}).
			Where("DATE(created_at) = ?", dateStr).
			Count(&stat.NewTickets)

		database.DB.Model(&models.Ticket{}).
			Where("status = ? AND DATE(resolved_at) = ?", models.TicketStatusResolved, dateStr).
			Count(&stat.ResolvedTickets)

		database.DB.Model(&models.Ticket{}).
			Where("status = ? AND DATE(closed_at) = ?", models.TicketStatusClosed, dateStr).
			Count(&stat.ClosedTickets)

		database.DB.Model(&models.SLARecord{}).
			Joins("JOIN tickets ON tickets.id = sla_records.ticket_id").
			Where("DATE(tickets.created_at) BETWEEN ? AND ?", dateStr, nextDate.Format("2006-01-02")).
			Where("response_time > ?", 0).
			Select("IFNULL(AVG(response_time), 0)").
			Scan(&stat.AvgResponseTime)

		database.DB.Model(&models.SLARecord{}).
			Joins("JOIN tickets ON tickets.id = sla_records.ticket_id").
			Where("DATE(tickets.resolved_at) = ?", dateStr).
			Where("resolution_time > ?", 0).
			Select("IFNULL(AVG(resolution_time), 0)").
			Scan(&stat.AvgResolutionTime)

		stats = append(stats, stat)
	}

	utils.Success(c, stats)
}

func GetWeeklyStats(c *gin.Context) {
	weeksStr := c.DefaultQuery("weeks", "8")
	weeks, err := strconv.Atoi(weeksStr)
	if err != nil || weeks < 1 || weeks > 52 {
		weeks = 8
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -weeks*7+1)

	stats := make([]WeeklyStats, 0)

	for i := 0; i < weeks; i++ {
		weekStart := startDate.AddDate(0, 0, i*7)
		weekEnd := weekStart.AddDate(0, 0, 6)
		weekStr := weekStart.Format("2006-01-02") + " to " + weekEnd.Format("2006-01-02")

		var stat WeeklyStats
		stat.Week = weekStr

		database.DB.Model(&models.Ticket{}).
			Where("created_at BETWEEN ? AND ?", weekStart, weekEnd.AddDate(0, 0, 1)).
			Count(&stat.NewTickets)

		database.DB.Model(&models.Ticket{}).
			Where("status = ? AND resolved_at BETWEEN ? AND ?", models.TicketStatusResolved, weekStart, weekEnd.AddDate(0, 0, 1)).
			Count(&stat.ResolvedTickets)

		database.DB.Model(&models.Ticket{}).
			Where("status = ? AND closed_at BETWEEN ? AND ?", models.TicketStatusClosed, weekStart, weekEnd.AddDate(0, 0, 1)).
			Count(&stat.ClosedTickets)

		database.DB.Model(&models.SLARecord{}).
			Joins("JOIN tickets ON tickets.id = sla_records.ticket_id").
			Where("tickets.created_at BETWEEN ? AND ?", weekStart, weekEnd.AddDate(0, 0, 1)).
			Where("response_time > ?", 0).
			Select("IFNULL(AVG(response_time), 0)").
			Scan(&stat.AvgResponseTime)

		database.DB.Model(&models.SLARecord{}).
			Joins("JOIN tickets ON tickets.id = sla_records.ticket_id").
			Where("tickets.resolved_at BETWEEN ? AND ?", weekStart, weekEnd.AddDate(0, 0, 1)).
			Where("resolution_time > ?", 0).
			Select("IFNULL(AVG(resolution_time), 0)").
			Scan(&stat.AvgResolutionTime)

		stats = append(stats, stat)
	}

	utils.Success(c, stats)
}

func GetMonthlyStats(c *gin.Context) {
	monthsStr := c.DefaultQuery("months", "6")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months < 1 || months > 24 {
		months = 6
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, -months+1, 1)

	stats := make([]MonthlyStats, 0)

	for i := 0; i < months; i++ {
		monthStart := startDate.AddDate(0, i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0).AddDate(0, 0, -1)
		monthStr := monthStart.Format("2006-01")

		var stat MonthlyStats
		stat.Month = monthStr

		database.DB.Model(&models.Ticket{}).
			Where("created_at BETWEEN ? AND ?", monthStart, monthEnd.AddDate(0, 0, 1)).
			Count(&stat.NewTickets)

		database.DB.Model(&models.Ticket{}).
			Where("status = ? AND resolved_at BETWEEN ? AND ?", models.TicketStatusResolved, monthStart, monthEnd.AddDate(0, 0, 1)).
			Count(&stat.ResolvedTickets)

		database.DB.Model(&models.Ticket{}).
			Where("status = ? AND closed_at BETWEEN ? AND ?", models.TicketStatusClosed, monthStart, monthEnd.AddDate(0, 0, 1)).
			Count(&stat.ClosedTickets)

		database.DB.Model(&models.SLARecord{}).
			Joins("JOIN tickets ON tickets.id = sla_records.ticket_id").
			Where("tickets.created_at BETWEEN ? AND ?", monthStart, monthEnd.AddDate(0, 0, 1)).
			Where("response_time > ?", 0).
			Select("IFNULL(AVG(response_time), 0)").
			Scan(&stat.AvgResponseTime)

		database.DB.Model(&models.SLARecord{}).
			Joins("JOIN tickets ON tickets.id = sla_records.ticket_id").
			Where("tickets.resolved_at BETWEEN ? AND ?", monthStart, monthEnd.AddDate(0, 0, 1)).
			Where("resolution_time > ?", 0).
			Select("IFNULL(AVG(resolution_time), 0)").
			Scan(&stat.AvgResolutionTime)

		stats = append(stats, stat)
	}

	utils.Success(c, stats)
}

func GetPriorityDistribution(c *gin.Context) {
	var distribution []PriorityDistribution

	rows, err := database.DB.Model(&models.Ticket{}).
		Select("priority, COUNT(*) as count").
		Group("priority").
		Rows()

	if err != nil {
		utils.InternalServerError(c, "Failed to get priority distribution")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d PriorityDistribution
		rows.Scan(&d.Priority, &d.Count)
		distribution = append(distribution, d)
	}

	utils.Success(c, distribution)
}

func GetStatusDistribution(c *gin.Context) {
	var distribution []StatusDistribution

	rows, err := database.DB.Model(&models.Ticket{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Rows()

	if err != nil {
		utils.InternalServerError(c, "Failed to get status distribution")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d StatusDistribution
		rows.Scan(&d.Status, &d.Count)
		distribution = append(distribution, d)
	}

	utils.Success(c, distribution)
}

func GetTypeDistribution(c *gin.Context) {
	var distribution []TypeDistribution

	rows, err := database.DB.Model(&models.Ticket{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Rows()

	if err != nil {
		utils.InternalServerError(c, "Failed to get type distribution")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d TypeDistribution
		rows.Scan(&d.Type, &d.Count)
		distribution = append(distribution, d)
	}

	utils.Success(c, distribution)
}

func GetAgentPerformance(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" {
		startDateStr = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	var performances []AgentPerformance

	rows, err := database.DB.Raw(`
		SELECT
			u.id as user_id,
			u.username,
			u.real_name,
			COUNT(DISTINCT t.id) as tickets_handled,
			COUNT(DISTINCT CASE WHEN t.status IN ('resolved', 'closed') THEN t.id END) as tickets_resolved,
			IFNULL(AVG(s.response_time), 0) as avg_response_time_minutes,
			IFNULL(AVG(s.resolution_time), 0) as avg_resolution_time_minutes
		FROM users u
		LEFT JOIN tickets t ON t.assignee_id = u.id AND t.created_at BETWEEN ? AND ?
		LEFT JOIN sla_records s ON s.ticket_id = t.id
		WHERE u.role IN ('admin', 'manager', 'agent')
		GROUP BY u.id, u.username, u.real_name
		ORDER BY tickets_handled DESC
	`, startDateStr, endDateStr+" 23:59:59").Rows()

	if err != nil {
		utils.InternalServerError(c, "Failed to get agent performance")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p AgentPerformance
		rows.Scan(&p.UserID, &p.Username, &p.RealName, &p.TicketsHandled, &p.TicketsResolved, &p.AvgResponseTime, &p.AvgResolutionTime)
		performances = append(performances, p)
	}

	utils.Success(c, performances)
}

func GetOverallStats(c *gin.Context) {
	var totalTickets int64
	var openTickets int64
	var inProgressTickets int64
	var resolvedTickets int64
	var closedTickets int64
	var escalatedTickets int64

	database.DB.Model(&models.Ticket{}).Count(&totalTickets)
	database.DB.Model(&models.Ticket{}).Where("status = ?", models.TicketStatusOpen).Count(&openTickets)
	database.DB.Model(&models.Ticket{}).Where("status = ?", models.TicketStatusInProgress).Count(&inProgressTickets)
	database.DB.Model(&models.Ticket{}).Where("status = ?", models.TicketStatusResolved).Count(&resolvedTickets)
	database.DB.Model(&models.Ticket{}).Where("status = ?", models.TicketStatusClosed).Count(&closedTickets)
	database.DB.Model(&models.Ticket{}).Where("is_escalated = ?", true).Count(&escalatedTickets)

	var avgResponseTime float64
	var avgResolutionTime float64
	var slaComplianceRate float64

	database.DB.Model(&models.SLARecord{}).Where("response_time > ?", 0).Select("IFNULL(AVG(response_time), 0)").Scan(&avgResponseTime)
	database.DB.Model(&models.SLARecord{}).Where("resolution_time > ?", 0).Select("IFNULL(AVG(resolution_time), 0)").Scan(&avgResolutionTime)

	var totalSLA int64
	var breachedSLA int64
	database.DB.Model(&models.SLARecord{}).Count(&totalSLA)
	database.DB.Model(&models.SLARecord{}).Where("response_breached = ? OR resolve_breached = ?", true, true).Count(&breachedSLA)

	if totalSLA > 0 {
		slaComplianceRate = float64(totalSLA-breachedSLA) / float64(totalSLA) * 100
	}

	stats := gin.H{
		"total_tickets":          totalTickets,
		"open_tickets":           openTickets,
		"in_progress_tickets":    inProgressTickets,
		"resolved_tickets":       resolvedTickets,
		"closed_tickets":         closedTickets,
		"escalated_tickets":      escalatedTickets,
		"avg_response_time_minutes":  avgResponseTime,
		"avg_resolution_time_minutes": avgResolutionTime,
		"sla_compliance_rate":    slaComplianceRate,
	}

	utils.Success(c, stats)
}
