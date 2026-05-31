package repository

import (
	"health-platform/models"
	"time"

	"gorm.io/gorm"
)

type HealthRecordRepository struct {
	*BaseRepository
}

func NewHealthRecordRepository() *HealthRecordRepository {
	return &HealthRecordRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *HealthRecordRepository) FindByEmployeeID(employeeID uint, page, pageSize int) ([]models.HealthRecord, int64, error) {
	var records []models.HealthRecord
	var total int64

	query := r.DB.Model(&models.HealthRecord{}).Where("employee_id = ?", employeeID)
	query.Count(&total)

	err := query.Preload("Report").
		Order("record_year DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records).Error
	return records, total, err
}

func (r *HealthRecordRepository) GetByEmployeeAndYear(employeeID uint, year int) (*models.HealthRecord, error) {
	var record models.HealthRecord
	err := r.DB.Where("employee_id = ? AND record_year = ?", employeeID, year).
		Preload("Report").Preload("Report.Items").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *HealthRecordRepository) GetAllByEmployee(employeeID uint) ([]models.HealthRecord, error) {
	var records []models.HealthRecord
	err := r.DB.Where("employee_id = ?", employeeID).
		Preload("Report").Preload("Report.Items").
		Order("record_year ASC").
		Find(&records).Error
	return records, err
}

type AbnormalItemRepository struct {
	*BaseRepository
}

func NewAbnormalItemRepository() *AbnormalItemRepository {
	return &AbnormalItemRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *AbnormalItemRepository) FindByEmployeeID(employeeID uint, page, pageSize int) ([]models.AbnormalItem, int64, error) {
	var items []models.AbnormalItem
	var total int64

	query := r.DB.Model(&models.AbnormalItem{}).Where("employee_id = ?", employeeID)
	query.Count(&total)

	err := query.Preload("HealthRecord").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error
	return items, total, err
}

func (r *AbnormalItemRepository) GetAllByEmployee(employeeID uint) ([]models.AbnormalItem, error) {
	var items []models.AbnormalItem
	err := r.DB.Where("employee_id = ?", employeeID).
		Preload("HealthRecord").
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *AbnormalItemRepository) GetNeedRecheckItems(employeeID uint) ([]models.AbnormalItem, error) {
	var items []models.AbnormalItem
	err := r.DB.Where("employee_id = ? AND recheck_date IS NOT NULL AND recheck_status = 0",
		employeeID).
		Order("recheck_date ASC").
		Find(&items).Error
	return items, err
}

func (r *AbnormalItemRepository) UpdateRecheckStatus(id uint, status int) error {
	return r.DB.Model(&models.AbnormalItem{}).Where("id = ?", id).
		Update("recheck_status", status).Error
}

type RecheckReminderRepository struct {
	*BaseRepository
}

func NewRecheckReminderRepository() *RecheckReminderRepository {
	return &RecheckReminderRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *RecheckReminderRepository) FindByEmployeeID(employeeID uint, page, pageSize int) ([]models.RecheckReminder, int64, error) {
	var reminders []models.RecheckReminder
	var total int64

	query := r.DB.Model(&models.RecheckReminder{}).Where("employee_id = ?", employeeID)
	query.Count(&total)

	err := query.Order("remind_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reminders).Error
	return reminders, total, err
}

func (r *RecheckReminderRepository) GetUnreadByEmployee(employeeID uint) ([]models.RecheckReminder, error) {
	var reminders []models.RecheckReminder
	err := r.DB.Where("employee_id = ? AND is_read = ?", employeeID, false).
		Order("remind_date DESC").
		Find(&reminders).Error
	return reminders, err
}

func (r *RecheckReminderRepository) MarkAsRead(id uint) error {
	return r.DB.Model(&models.RecheckReminder{}).Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *RecheckReminderRepository) GetNeedSendReminders() ([]models.RecheckReminder, error) {
	var reminders []models.RecheckReminder
	err := r.DB.Where("is_sent = ? AND remind_date <= ?", false, time.Now()).
		Find(&reminders).Error
	return reminders, err
}

func (r *RecheckReminderRepository) MarkAsSent(id uint) error {
	return r.DB.Model(&models.RecheckReminder{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_sent": true,
			"sent_at": time.Now(),
		}).Error
}
