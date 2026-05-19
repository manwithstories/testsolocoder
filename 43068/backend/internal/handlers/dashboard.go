package handlers

import (
	"time"

	"freelancer-management/internal/database"
	"freelancer-management/internal/middleware"
	"freelancer-management/internal/models"
	"freelancer-management/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{db: database.GetDB()}
}

type DashboardStats struct {
	MonthlyRevenue   float64                    `json:"monthly_revenue"`
	TotalClients     int64                      `json:"total_clients"`
	ActiveProjects   int64                      `json:"active_projects"`
	TotalHours       float64                    `json:"total_hours"`
	OverdueInvoices  int64                      `json:"overdue_invoices"`
	PendingInvoices  int64                      `json:"pending_invoices"`
	ProjectProgress  []ProjectProgressStats     `json:"project_progress"`
	MonthlyEarnings  []MonthlyEarningsStats     `json:"monthly_earnings"`
	OverdueReminders []OverdueReminder          `json:"overdue_reminders"`
}

type ProjectProgressStats struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	TotalHours   float64 `json:"total_hours"`
	Budget       float64 `json:"budget"`
	EarnedAmount float64 `json:"earned_amount"`
	Progress     float64 `json:"progress"`
	DaysLeft     int     `json:"days_left"`
}

type MonthlyEarningsStats struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
	Hours   float64 `json:"hours"`
}

type OverdueReminder struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	DueDate     time.Time `json:"due_date"`
	DaysOverdue int       `json:"days_overdue"`
	Amount      float64   `json:"amount,omitempty"`
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := currentMonth.AddDate(0, 1, 0)

	var stats DashboardStats

	h.db.Model(&models.Invoice{}).
		Where("user_id = ? AND status IN ? AND issue_date >= ? AND issue_date < ?",
			userID, []models.InvoiceStatus{models.InvoiceStatusPaid, models.InvoiceStatusSent},
			currentMonth, nextMonth).
		Select("COALESCE(SUM(total), 0)").
		Scan(&stats.MonthlyRevenue)

	h.db.Model(&models.Client{}).Where("user_id = ?", userID).Count(&stats.TotalClients)

	h.db.Model(&models.Project{}).
		Where("user_id = ? AND status = ?", userID, models.ProjectStatusActive).
		Count(&stats.ActiveProjects)

	h.db.Model(&models.TimeEntry{}).
		Where("user_id = ? AND date >= ? AND date < ?", userID, currentMonth, nextMonth).
		Select("COALESCE(SUM(hours), 0)").
		Scan(&stats.TotalHours)

	h.db.Model(&models.Invoice{}).
		Where("user_id = ? AND status = ? AND due_date < ?", userID, models.InvoiceStatusSent, now).
		Count(&stats.OverdueInvoices)

	h.db.Model(&models.Invoice{}).
		Where("user_id = ? AND status = ?", userID, models.InvoiceStatusSent).
		Count(&stats.PendingInvoices)

	stats.ProjectProgress = h.getProjectProgress(userID)
	stats.MonthlyEarnings = h.getMonthlyEarnings(userID)
	stats.OverdueReminders = h.getOverdueReminders(userID, now)

	utils.SuccessResponse(c, stats)
}

func (h *DashboardHandler) getProjectProgress(userID uint) []ProjectProgressStats {
	var projects []models.Project
	h.db.Where("user_id = ? AND status IN ?", userID,
		[]models.ProjectStatus{models.ProjectStatusActive, models.ProjectStatusDraft}).
		Preload("TimeEntries").
		Find(&projects)

	var progress []ProjectProgressStats
	now := time.Now()

	for _, p := range projects {
		totalHours := 0.0
		for _, t := range p.TimeEntries {
			totalHours += t.Hours
		}

		earnedAmount := totalHours * p.HourlyRate
		progressPercent := 0.0
		if p.Budget > 0 {
			progressPercent = (earnedAmount / p.Budget) * 100
		}

		daysLeft := 0
		if p.Deadline != nil {
			daysLeft = int(p.Deadline.Sub(now).Hours() / 24)
		}

		progress = append(progress, ProjectProgressStats{
			ID:           p.ID,
			Name:         p.Name,
			Status:       string(p.Status),
			TotalHours:   totalHours,
			Budget:       p.Budget,
			EarnedAmount: earnedAmount,
			Progress:     progressPercent,
			DaysLeft:     daysLeft,
		})
	}

	return progress
}

func (h *DashboardHandler) getMonthlyEarnings(userID uint) []MonthlyEarningsStats {
	var earnings []MonthlyEarningsStats
	now := time.Now()

	for i := 5; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		monthStart := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, 0)

		var revenue float64
		h.db.Model(&models.Invoice{}).
			Where("user_id = ? AND status = ? AND issue_date >= ? AND issue_date < ?",
				userID, models.InvoiceStatusPaid, monthStart, monthEnd).
			Select("COALESCE(SUM(total), 0)").
			Scan(&revenue)

		var hours float64
		h.db.Model(&models.TimeEntry{}).
			Where("user_id = ? AND date >= ? AND date < ?", userID, monthStart, monthEnd).
			Select("COALESCE(SUM(hours), 0)").
			Scan(&hours)

		earnings = append(earnings, MonthlyEarningsStats{
			Month:   month.Format("Jan 2006"),
			Revenue: revenue,
			Hours:   hours,
		})
	}

	return earnings
}

func (h *DashboardHandler) getOverdueReminders(userID uint, now time.Time) []OverdueReminder {
	var reminders []OverdueReminder

	var overdueInvoices []models.Invoice
	h.db.Where("user_id = ? AND status = ? AND due_date < ?", userID, models.InvoiceStatusSent, now).
		Preload("Client").
		Find(&overdueInvoices)

	for _, inv := range overdueInvoices {
		daysOverdue := int(now.Sub(inv.DueDate).Hours() / 24)
		reminders = append(reminders, OverdueReminder{
			ID:          inv.ID,
			Type:        "invoice",
			Title:       "Overdue Invoice: " + inv.InvoiceNumber,
			DueDate:     inv.DueDate,
			DaysOverdue: daysOverdue,
			Amount:      inv.Total,
		})
	}

	var overdueProjects []models.Project
	h.db.Where("user_id = ? AND status = ? AND deadline < ?", userID, models.ProjectStatusActive, now).
		Find(&overdueProjects)

	for _, p := range overdueProjects {
		daysOverdue := int(now.Sub(*p.Deadline).Hours() / 24)
		reminders = append(reminders, OverdueReminder{
			ID:          p.ID,
			Type:        "project",
			Title:       "Overdue Project: " + p.Name,
			DueDate:     *p.Deadline,
			DaysOverdue: daysOverdue,
		})
	}

	return reminders
}
