package services

import (
	"time"
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"

	"gorm.io/gorm"
)

type SchedulerService struct{}

func NewSchedulerService() *SchedulerService {
	return &SchedulerService{}
}

func (s *SchedulerService) CheckPaymentReminders() {
	db := database.GetDB()

	var budgetItems []models.BudgetItem
	now := time.Now()
	reminderDate := now.AddDate(0, 0, 7)

	db.Where("due_date BETWEEN ? AND ? AND status NOT IN ('paid', 'cancelled')", now, reminderDate).
		Find(&budgetItems)

	for _, item := range budgetItems {
		var wedding models.Wedding
		db.First(&wedding, item.WeddingID)

		notification := models.Notification{
			UserID:    wedding.UserID,
			Type:      "payment_reminder",
			Title:     "付款提醒",
			Content:   "您有一笔付款即将到期: " + item.Category + " - " + item.Description,
			RelatedID: item.ID,
		}
		db.Create(&notification)
	}
}

func (s *SchedulerService) CheckTaskDeadlines() {
	db := database.GetDB()

	var tasks []models.Task
	now := time.Now()
	deadlineDate := now.AddDate(0, 0, 3)

	db.Where("due_date BETWEEN ? AND ? AND status != ?", now, deadlineDate, "completed").
		Find(&tasks)

	for _, task := range tasks {
		var wedding models.Wedding
		db.First(&wedding, task.WeddingID)

		notification := models.Notification{
			UserID:    wedding.UserID,
			Type:      "task_reminder",
			Title:     "任务截止提醒",
			Content:   "任务即将截止: " + task.Title,
			RelatedID: task.ID,
		}
		db.Create(&notification)
	}
}

func (s *SchedulerService) CheckOverduePayments() {
	db := database.GetDB()

	var budgetItems []models.BudgetItem
	now := time.Now()

	db.Where("due_date < ? AND status NOT IN ('paid', 'cancelled') AND actual_cost > paid_amount", now).
		Find(&budgetItems)

	for _, item := range budgetItems {
		var wedding models.Wedding
		db.First(&wedding, item.WeddingID)

		notification := models.Notification{
			UserID:    wedding.UserID,
			Type:      "payment_overdue",
			Title:     "付款逾期提醒",
			Content:   "您有一笔付款已逾期: " + item.Category + " - " + item.Description,
			RelatedID: item.ID,
		}
		db.Create(&notification)
	}
}

func (s *SchedulerService) RunDailyTasks() {
	s.CheckPaymentReminders()
	s.CheckTaskDeadlines()
	s.CheckOverduePayments()
}

func (s *SchedulerService) StartScheduler() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			s.RunDailyTasks()
		}
	}()
}

func CreateOperationLog(db *gorm.DB, userID uint, module, action string, targetID uint, detail, ipAddress string) {
	log := models.OperationLog{
		UserID:    userID,
		Module:    module,
		Action:    action,
		TargetID:  targetID,
		Detail:    detail,
		IPAddress: ipAddress,
	}
	db.Create(&log)
}
