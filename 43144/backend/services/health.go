package services

import (
	"fmt"
	"time"

	"pet-adoption-platform/database"
	"pet-adoption-platform/models"
)

func CreateHealthRecord(req *models.CreateHealthRecordRequest, recordedBy uint, rescueID uint) (*models.HealthRecord, error) {
	pet, err := GetPetByID(req.PetID)
	if err != nil {
		return nil, fmt.Errorf("pet not found: %w", err)
	}

	recordDate, _ := time.Parse("2006-01-02", req.RecordDate)

	record := &models.HealthRecord{
		PetID:       req.PetID,
		RecordType:  models.HealthRecordType(req.RecordType),
		Title:       req.Title,
		Description: req.Description,
		VaccineName: req.VaccineName,
		RecordDate:  recordDate,
		Weight:      req.Weight,
		Temperature: req.Temperature,
		VetName:     req.VetName,
		Hospital:    req.Hospital,
		Notes:       req.Notes,
		RecordedBy:  recordedBy,
		RescueID:    pet.RescueID,
	}

	var nextDate *time.Time
	if req.NextDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.NextDate)
		if err == nil {
			record.NextDate = &parsedDate
			nextDate = &parsedDate
		}
	}

	if err := database.DB.Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to create health record: %w", err)
	}

	if nextDate != nil {
		reminderTitle := req.Title
		if req.VaccineName != "" {
			reminderTitle = "下次接种: " + req.VaccineName
		} else if req.RecordType == "vaccine" {
			reminderTitle = "下次疫苗接种"
		} else if req.RecordType == "deworm" {
			reminderTitle = "下次驱虫"
		} else if req.RecordType == "checkup" {
			reminderTitle = "下次体检"
		}

		reminder := &models.HealthReminder{
			PetID:        req.PetID,
			RecordID:     &record.ID,
			Title:        reminderTitle,
			ReminderDate: *nextDate,
			Notes:        req.Notes,
		}
		if err := database.DB.Create(reminder).Error; err != nil {
			return record, nil
		}
	}

	return record, nil
}

func GetHealthRecordByID(id uint) (*models.HealthRecord, error) {
	var record models.HealthRecord
	if err := database.DB.Preload("Pet").First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func ListHealthRecords(query *models.HealthRecordListQuery) ([]models.HealthRecord, int64, error) {
	var records []models.HealthRecord
	var total int64

	db := database.DB.Model(&models.HealthRecord{})

	if query.PetID > 0 {
		db = db.Where("pet_id = ?", query.PetID)
	}
	if query.RecordType != "" {
		db = db.Where("record_type = ?", query.RecordType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Pet").
		Order("record_date DESC").
		Offset(offset).Limit(query.PageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func UpdateHealthRecord(id uint, updates map[string]interface{}) (*models.HealthRecord, error) {
	result := database.DB.Model(&models.HealthRecord{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	return GetHealthRecordByID(id)
}

func DeleteHealthRecord(id uint) error {
	return database.DB.Delete(&models.HealthRecord{}, id).Error
}

func UpdateReportFile(id uint, filePath string) error {
	return database.DB.Model(&models.HealthRecord{}).Where("id = ?", id).Update("report_file", filePath).Error
}

func GetHealthReminders(petID uint) ([]models.HealthReminder, error) {
	var reminders []models.HealthReminder
	err := database.DB.Where("pet_id = ? AND is_completed = ?", petID, false).
		Order("reminder_date ASC").
		Find(&reminders).Error
	return reminders, err
}

func CreateHealthReminder(petID uint, title string, reminderDate time.Time, recordID *uint) (*models.HealthReminder, error) {
	reminder := &models.HealthReminder{
		PetID:        petID,
		RecordID:     recordID,
		Title:        title,
		ReminderDate: reminderDate,
	}
	if err := database.DB.Create(reminder).Error; err != nil {
		return nil, err
	}
	return reminder, nil
}

func CompleteHealthReminder(id uint) error {
	return database.DB.Model(&models.HealthReminder{}).Where("id = ?", id).Update("is_completed", true).Error
}

func GetPetHealthSummary(petID uint) (map[string]interface{}, error) {
	var records []models.HealthRecord
	if err := database.DB.Where("pet_id = ?", petID).Find(&records).Error; err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"total_records": len(records),
		"vaccines":      0,
		"checkups":      0,
		"diseases":      0,
		"surgeries":     0,
		"dewormings":    0,
	}

	for _, r := range records {
		switch r.RecordType {
		case models.HealthRecordVaccine:
			summary["vaccines"] = summary["vaccines"].(int) + 1
		case models.HealthRecordCheckup:
			summary["checkups"] = summary["checkups"].(int) + 1
		case models.HealthRecordDisease:
			summary["diseases"] = summary["diseases"].(int) + 1
		case models.HealthRecordSurgery:
			summary["surgeries"] = summary["surgeries"].(int) + 1
		case models.HealthRecordDeworm:
			summary["dewormings"] = summary["dewormings"].(int) + 1
		}
	}

	return summary, nil
}
