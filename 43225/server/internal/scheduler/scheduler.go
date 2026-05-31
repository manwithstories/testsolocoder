package scheduler

import (
	"time"

	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/model"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var Scheduler *cron.Cron

func InitScheduler() {
	Scheduler = cron.New()

	Scheduler.AddFunc("0 0 8 * * *", checkMaintenanceReminders)
	Scheduler.AddFunc("0 0 0 * * *", updateRentalStatuses)
	Scheduler.AddFunc("0 0 0 1 * *", generateMonthlySettlements)

	Scheduler.Start()
	logrus.Info("Scheduler started")
}

func checkMaintenanceReminders() {
	var schedules []model.MaintenanceSchedule
	now := time.Now()
	reminderDate := now.AddDate(0, 0, 7)

	database.DB.Where("next_due <= ? AND is_active = ? AND last_completed < ?", reminderDate, true, now).
		Find(&schedules)

	for _, schedule := range schedules {
		logrus.Infof("Maintenance reminder: Ship %s - %s (Due: %s)",
			schedule.ShipID, schedule.Title, schedule.NextDue.Format("2006-01-02"))
	}
}

func updateRentalStatuses() {
	now := time.Now()

	database.DB.Model(&model.Rental{}).
		Where("status = ? AND start_date <= ?", model.RentalStatusConfirmed, now).
		Update("status", model.RentalStatusActive)

	database.DB.Model(&model.Rental{}).
		Where("status = ? AND end_date <= ?", model.RentalStatusActive, now).
		Update("status", model.RentalStatusCompleted)

	var completedRentals []model.Rental
	database.DB.Where("status = ?", model.RentalStatusCompleted).
		Where("completed_at IS NULL").
		Find(&completedRentals)

	for _, rental := range completedRentals {
		now := now
		rental.CompletedAt = &now
		database.DB.Save(&rental)

		var activeRentals int64
		database.DB.Model(&model.Rental{}).Where(
			"ship_id = ? AND status IN ?",
			rental.ShipID,
			[]model.RentalStatus{model.RentalStatusConfirmed, model.RentalStatusActive},
		).Count(&activeRentals)

		if activeRentals == 0 {
			database.DB.Model(&model.Ship{}).
				Where("id = ?", rental.ShipID).
				Update("status", model.ShipStatusAvailable)
		}
	}
}

func generateMonthlySettlements() {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	periodStart := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)

	var users []model.User
	database.DB.Where("role IN ?", []model.Role{model.RoleOwner, model.RoleAdmin}).Find(&users)

	for _, user := range users {
		var totalIncome, totalExpense float64

		database.DB.Model(&model.Transaction{}).
			Where("payee_id = ? AND transaction_type = ? AND status = ? AND created_at BETWEEN ? AND ?",
				user.ID, model.TransactionTypeIncome, model.TransactionStatusCompleted, periodStart, periodEnd).
			Select("COALESCE(SUM(net_amount), 0)").
			Scan(&totalIncome)

		database.DB.Model(&model.Transaction{}).
			Where("payer_id = ? AND transaction_type = ? AND status = ? AND created_at BETWEEN ? AND ?",
				user.ID, model.TransactionTypeExpense, model.TransactionStatusCompleted, periodStart, periodEnd).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&totalExpense)

		if totalIncome > 0 || totalExpense > 0 {
			settlement := model.Settlement{
				UserID:       user.ID,
				PeriodStart:  periodStart,
				PeriodEnd:    periodEnd,
				TotalIncome:  totalIncome,
				TotalExpense: totalExpense,
				Currency:     "USD",
				Status:       "pending",
			}
			database.DB.Create(&settlement)
		}
	}

	logrus.Infof("Generated monthly settlements for period %s to %s",
		periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02"))
}

func StopScheduler() {
	if Scheduler != nil {
		Scheduler.Stop()
		logrus.Info("Scheduler stopped")
	}
}
