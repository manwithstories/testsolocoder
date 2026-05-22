package services

import (
	"errors"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"

	"gorm.io/gorm"
)

type PatientService struct {
	db *gorm.DB
}

func NewPatientService() *PatientService {
	return &PatientService{
		db: database.GetDB(),
	}
}

type PatientListQuery struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}

type UpdatePatientRequest struct {
	IDCardNo       string `json:"id_card_no"`
	Address         string `json:"address"`
	EmergencyName   string `json:"emergency_contact_name"`
	EmergencyPhone  string `json:"emergency_contact_phone"`
}

func (s *PatientService) GetPatientList(query PatientListQuery) ([]models.Patient, int64, error) {
	var patients []models.Patient
	var total int64

	db := s.db.Model(&models.Patient{}).Preload("User")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("id DESC")

	if err := db.Scopes(database.Paginate(query.Page, query.PageSize)).Find(&patients).Error; err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

func (s *PatientService) GetPatientByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	err := s.db.Preload("User").
		Preload("HealthRecord").
		Preload("Appointments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Doctor.User").Preload("Doctor.Department").Order("appointment_date DESC")
		}).
		First(&patient, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("患者不存在")
		}
		return nil, err
	}
	return &patient, nil
}

func (s *PatientService) GetPatientByUserID(userID uint) (*models.Patient, error) {
	var patient models.Patient
	err := s.db.Where("user_id = ?", userID).First(&patient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("患者不存在")
		}
		return nil, err
	}
	return &patient, nil
}

func (s *PatientService) UpdatePatient(id uint, req UpdatePatientRequest) (*models.Patient, error) {
	patient, err := s.GetPatientByID(id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.IDCardNo != "" {
		updates["id_card_no"] = req.IDCardNo
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.EmergencyName != "" {
		updates["emergency_contact_name"] = req.EmergencyName
	}
	if req.EmergencyPhone != "" {
		updates["emergency_contact_phone"] = req.EmergencyPhone
	}

	if len(updates) > 0 {
		if err := s.db.Model(patient).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return s.GetPatientByID(id)
}

func (s *PatientService) DeletePatient(id uint) error {
	patient, err := s.GetPatientByID(id)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("patient_id = ?", id).Delete(&models.Appointment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("patient_id = ?", id).Delete(&models.HealthRecord{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(patient).Error; err != nil {
			return err
		}
		return nil
	})
}
