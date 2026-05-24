package services

import (
	"errors"
	"fmt"
	"time"

	"pet-adoption-platform/database"
	"pet-adoption-platform/models"

	"gorm.io/gorm"
)

func CreateAdoptionApplication(req *models.CreateAdoptionRequest, adopterID uint) (*models.AdoptionApplication, error) {
	pet, err := GetPetByID(req.PetID)
	if err != nil {
		return nil, errors.New("pet not found")
	}

	if pet.Status != models.PetStatusAdoptable {
		return nil, errors.New("this pet is not available for adoption")
	}

	var existing models.AdoptionApplication
	if err := database.DB.Where("pet_id = ? AND adopter_id = ? AND status IN ?",
		req.PetID, adopterID,
		[]models.AdoptionStatus{models.AdoptionStatusPending, models.AdoptionStatusApproved},
	).First(&existing).Error; err == nil {
		return nil, errors.New("you already have an active application for this pet")
	}

	application := &models.AdoptionApplication{
		PetID:           req.PetID,
		AdopterID:       adopterID,
		RescueID:        pet.RescueID,
		Status:          models.AdoptionStatusPending,
		Reason:          req.Reason,
		LivingSituation: req.LivingSituation,
		PetExperience:   req.PetExperience,
		FamilyMembers:   req.FamilyMembers,
		HasChildren:     req.HasChildren,
		HasOtherPets:    req.HasOtherPets,
		OtherPetsDesc:   req.OtherPetsDesc,
		HousingType:     req.HousingType,
		IncomeLevel:     req.IncomeLevel,
		CanAffordVet:    req.CanAffordVet,
		AgreeToVisit:    req.AgreeToVisit,
	}

	if err := database.DB.Create(application).Error; err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return application, nil
}

func GetAdoptionApplicationByID(id uint) (*models.AdoptionApplication, error) {
	var app models.AdoptionApplication
	if err := database.DB.Preload("Pet").Preload("Adopter").First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func ListAdoptionApplications(query *models.AdoptionListQuery, userID uint, role string, rescueID *uint) ([]models.AdoptionApplication, int64, error) {
	var apps []models.AdoptionApplication
	var total int64

	db := database.DB.Model(&models.AdoptionApplication{})

	switch role {
	case "adopter":
		db = db.Where("adopter_id = ?", userID)
	case "rescue":
		if rescueID != nil {
			db = db.Where("rescue_id = ?", *rescueID)
		}
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.PetID > 0 {
		db = db.Where("pet_id = ?", query.PetID)
	}
	if query.RescueID > 0 {
		db = db.Where("rescue_id = ?", query.RescueID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Pet").Preload("Adopter").
		Order("created_at DESC").
		Offset(offset).Limit(query.PageSize).
		Find(&apps).Error; err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

func ReviewAdoptionApplication(id uint, reviewerID uint, req *models.ReviewAdoptionRequest) (*models.AdoptionApplication, error) {
	app, err := GetAdoptionApplicationByID(id)
	if err != nil {
		return nil, errors.New("application not found")
	}

	if app.Status != models.AdoptionStatusPending {
		return nil, errors.New("application is not in pending status")
	}

	now := time.Now()
	if req.Action == "approve" {
		app.Status = models.AdoptionStatusApproved
		app.ReviewedBy = &reviewerID
		app.ReviewedAt = &now

		var existing models.AdoptionAgreement
		if err := database.DB.Where("application_id = ?", app.ID).First(&existing).Error; err != nil {
			agreement := &models.AdoptionAgreement{
				ApplicationID: app.ID,
				AgreementTerms: "本协议由救助站与领养人共同签署，确保领养宠物得到妥善照顾。救助站承诺提供真实健康的宠物信息，领养人承诺给予宠物良好生活环境、定期医疗、必要训练，不得遗弃或虐待宠物。如有违反，救助站有权收回宠物并追究相关责任。",
			}
			database.DB.Create(agreement)
		}
	} else if req.Action == "reject" {
		app.Status = models.AdoptionStatusRejected
		app.ReviewedBy = &reviewerID
		app.ReviewedAt = &now
		app.RejectReason = req.RejectReason
	}

	if err := database.DB.Save(app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

func SignAdoptionAgreement(applicationID uint, role string, userID uint) (*models.AdoptionAgreement, error) {
	var agreement models.AdoptionAgreement
	if err := database.DB.Where("application_id = ?", applicationID).First(&agreement).Error; err != nil {
		return nil, errors.New("agreement not found")
	}

	now := time.Now()
	if role == "adopter" {
		if agreement.AdopterSign {
			return nil, errors.New("adopter has already signed")
		}
		agreement.AdopterSign = true
		agreement.AdopterSignedAt = &now
	} else if role == "rescue" {
		if agreement.RescueSign {
			return nil, errors.New("rescue has already signed")
		}
		agreement.RescueSign = true
		agreement.RescueSignedAt = &now
	}

	if agreement.AdopterSign && agreement.RescueSign {
		agreement.IsActive = true
	}

	if err := database.DB.Save(&agreement).Error; err != nil {
		return nil, err
	}

	if agreement.AdopterSign && agreement.RescueSign {
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			var app models.AdoptionApplication
			if err := tx.First(&app, applicationID).Error; err != nil {
				return err
			}

			app.Status = models.AdoptionStatusSigned
			app.SignedAt = &now
			if err := tx.Save(&app).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.Pet{}).Where("id = ?", app.PetID).
				Updates(map[string]interface{}{
					"status":       models.PetStatusAdopted,
					"adopter_id":   app.AdopterID,
					"adopted_date": now,
				}).Error; err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to complete adoption: %w", err)
		}
	}

	return &agreement, nil
}

func GetAdoptionAgreement(applicationID uint) (*models.AdoptionAgreement, error) {
	var agreement models.AdoptionAgreement
	if err := database.DB.Preload("Application").Where("application_id = ?", applicationID).First(&agreement).Error; err != nil {
		return nil, err
	}
	return &agreement, nil
}

func CompleteAdoption(applicationID uint) (*models.AdoptionApplication, error) {
	app, err := GetAdoptionApplicationByID(applicationID)
	if err != nil {
		return nil, errors.New("application not found")
	}

	if app.Status != models.AdoptionStatusSigned {
		return nil, errors.New("application must be signed first")
	}

	now := time.Now()
	app.Status = models.AdoptionStatusCompleted
	app.CompletedAt = &now

	if err := database.DB.Save(app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

func CreateFollowUpRecord(req *models.CreateFollowUpRequest, recordedBy uint, rescueID uint) (*models.FollowUpRecord, error) {
	app, err := GetAdoptionApplicationByID(req.ApplicationID)
	if err != nil {
		return nil, errors.New("application not found")
	}

	if app.Status != models.AdoptionStatusSigned && app.Status != models.AdoptionStatusCompleted {
		return nil, errors.New("can only create follow-up for signed or completed applications")
	}

	followUpDate, _ := time.Parse("2006-01-02", req.FollowUpDate)

	record := &models.FollowUpRecord{
		ApplicationID:   req.ApplicationID,
		PetID:           app.PetID,
		AdopterID:       app.AdopterID,
		RescueID:        rescueID,
		FollowUpDate:    &followUpDate,
		HealthStatus:    req.HealthStatus,
		LivingCondition: req.LivingCondition,
		Notes:           req.Notes,
		RecordedBy:      recordedBy,
	}

	if err := database.DB.Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to create follow-up record: %w", err)
	}

	return record, nil
}

func ListFollowUpRecords(petID uint) ([]models.FollowUpRecord, error) {
	var records []models.FollowUpRecord
	err := database.DB.Preload("Application").
		Where("pet_id = ?", petID).
		Order("follow_up_date DESC").
		Find(&records).Error
	return records, err
}
