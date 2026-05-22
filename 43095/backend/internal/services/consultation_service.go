package services

import (
	"errors"
	"fmt"
	"io"
	"medical-platform/internal/config"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"
	"medical-platform/pkg/utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type ConsultationService struct {
	db *gorm.DB
}

func NewConsultationService() *ConsultationService {
	return &ConsultationService{
		db: database.GetDB(),
	}
}

type CreateConsultationRequest struct {
	AppointmentID uint   `json:"appointment_id" binding:"required"`
	Diagnosis     string `json:"diagnosis" binding:"required"`
	TreatmentPlan string `json:"treatment_plan"`
	DoctorNotes   string `json:"doctor_notes"`
	FollowUpDate  *time.Time `json:"follow_up_date"`
}

type UpdateConsultationRequest struct {
	Diagnosis     string `json:"diagnosis"`
	TreatmentPlan string `json:"treatment_plan"`
	DoctorNotes   string `json:"doctor_notes"`
	FollowUpDate  *time.Time `json:"follow_up_date"`
}

type CreatePrescriptionRequest struct {
	ConsultationID uint                      `json:"consultation_id" binding:"required"`
	Notes          string                    `json:"notes"`
	Items          []PrescriptionItemRequest `json:"items" binding:"required,min=1"`
}

type PrescriptionItemRequest struct {
	DrugName      string  `json:"drug_name" binding:"required"`
	Specification  string  `json:"specification"`
	Dosage         string  `json:"dosage" binding:"required"`
	Frequency      string  `json:"frequency"`
	Duration       string  `json:"duration"`
	Quantity       int     `json:"quantity" binding:"min=1"`
	UnitPrice      float64 `json:"unit_price" binding:"min=0"`
	Notes          string  `json:"notes"`
}

type UploadReportRequest struct {
	ConsultationID uint   `json:"consultation_id" binding:"required"`
	ReportType     string `json:"report_type" binding:"required"`
	ReportName     string `json:"report_name" binding:"required"`
	Findings       string `json:"findings"`
	Conclusion     string `json:"conclusion"`
}

func (s *ConsultationService) CreateConsultation(doctorUserID uint, req *CreateConsultationRequest) (*models.Consultation, error) {
	var consultation *models.Consultation

	err := database.WithTransaction(func(tx *gorm.DB) error {
		var doctor models.Doctor
		if err := tx.Where("user_id = ?", doctorUserID).First(&doctor).Error; err != nil {
			return errors.New("医生信息不存在")
		}

		var appointment models.Appointment
		if err := tx.First(&appointment, req.AppointmentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("预约不存在")
			}
			return err
		}

		if appointment.DoctorID != doctor.ID {
			return errors.New("无权为该预约创建问诊记录")
		}

		if appointment.Status != models.AppointmentConfirmed {
			return errors.New("预约未确认，无法创建问诊记录")
		}

		var existingConsultation models.Consultation
		if err := tx.Where("appointment_id = ?", req.AppointmentID).First(&existingConsultation).Error; err == nil {
			return errors.New("该预约已存在问诊记录")
		}

		consultation = &models.Consultation{
			AppointmentID: req.AppointmentID,
			Diagnosis:     req.Diagnosis,
			TreatmentPlan: req.TreatmentPlan,
			DoctorNotes:   req.DoctorNotes,
			FollowUpDate:  req.FollowUpDate,
		}

		if err := tx.Create(consultation).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return consultation, nil
}

func (s *ConsultationService) GetConsultationDetail(id uint, userID uint, role models.UserRole) (*models.Consultation, error) {
	var consultation models.Consultation
	query := s.db.Preload("Appointment.Patient.User").Preload("Appointment.Doctor.User").Preload("Prescription.Items").Preload("Reports")

	if err := query.First(&consultation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("问诊记录不存在")
		}
		return nil, err
	}

	switch role {
	case models.RolePatient:
		var patient models.Patient
		if err := s.db.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			return nil, errors.New("患者信息不存在")
		}
		if consultation.Appointment.PatientID != patient.ID {
			return nil, errors.New("无权查看该问诊记录")
		}
	case models.RoleDoctor:
		var doctor models.Doctor
		if err := s.db.Where("user_id = ?", userID).First(&doctor).Error; err != nil {
			return nil, errors.New("医生信息不存在")
		}
		if consultation.Appointment.DoctorID != doctor.ID {
			return nil, errors.New("无权查看该问诊记录")
		}
	case models.RoleAdmin:
	default:
		return nil, errors.New("无效的用户角色")
	}

	return &consultation, nil
}

func (s *ConsultationService) UpdateConsultation(id uint, doctorUserID uint, req *UpdateConsultationRequest) (*models.Consultation, error) {
	var updatedConsultation *models.Consultation

	err := database.WithTransaction(func(tx *gorm.DB) error {
		var doctor models.Doctor
		if err := tx.Where("user_id = ?", doctorUserID).First(&doctor).Error; err != nil {
			return errors.New("医生信息不存在")
		}

		var consultation models.Consultation
		if err := tx.First(&consultation, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("问诊记录不存在")
			}
			return err
		}

		var appointment models.Appointment
		if err := tx.First(&appointment, consultation.AppointmentID).Error; err != nil {
			return err
		}

		if appointment.DoctorID != doctor.ID {
			return errors.New("无权修改该问诊记录")
		}

		updates := make(map[string]interface{})
		if req.Diagnosis != "" {
			updates["diagnosis"] = req.Diagnosis
		}
		if req.TreatmentPlan != "" {
			updates["treatment_plan"] = req.TreatmentPlan
		}
		if req.DoctorNotes != "" {
			updates["doctor_notes"] = req.DoctorNotes
		}
		if req.FollowUpDate != nil {
			updates["follow_up_date"] = req.FollowUpDate
		}

		if err := tx.Model(&consultation).Updates(updates).Error; err != nil {
			return err
		}

		updatedConsultation = &consultation
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedConsultation, nil
}

func (s *ConsultationService) CreatePrescription(doctorUserID uint, req *CreatePrescriptionRequest) (*models.Prescription, error) {
	var prescription *models.Prescription

	err := database.WithTransaction(func(tx *gorm.DB) error {
		var doctor models.Doctor
		if err := tx.Where("user_id = ?", doctorUserID).First(&doctor).Error; err != nil {
			return errors.New("医生信息不存在")
		}

		var consultation models.Consultation
		if err := tx.First(&consultation, req.ConsultationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("问诊记录不存在")
			}
			return err
		}

		var appointment models.Appointment
		if err := tx.First(&appointment, consultation.AppointmentID).Error; err != nil {
			return err
		}

		if appointment.DoctorID != doctor.ID {
			return errors.New("无权为该问诊记录开具处方")
		}

		var existingPrescription models.Prescription
		if err := tx.Where("consultation_id = ?", req.ConsultationID).First(&existingPrescription).Error; err == nil {
			return errors.New("该问诊记录已存在处方")
		}

		prescriptionNo := utils.GenerateOrderNo("RX")
		prescription = &models.Prescription{
			ConsultationID: req.ConsultationID,
			PrescriptionNo: prescriptionNo,
			Notes:          req.Notes,
			IsFulfilled:    false,
		}

		if err := tx.Create(prescription).Error; err != nil {
			return err
		}

		var items []models.PrescriptionItem
		var totalAmount float64

		for _, itemReq := range req.Items {
			subtotal := float64(itemReq.Quantity) * itemReq.UnitPrice
			totalAmount += subtotal

			item := models.PrescriptionItem{
				PrescriptionID: prescription.ID,
				DrugName:      itemReq.DrugName,
				Specification: itemReq.Specification,
				Dosage:        itemReq.Dosage,
				Frequency:     itemReq.Frequency,
				Duration:      itemReq.Duration,
				Quantity:      itemReq.Quantity,
				UnitPrice:     itemReq.UnitPrice,
				Subtotal:      subtotal,
				Notes:         itemReq.Notes,
			}
			items = append(items, item)
		}

		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return prescription, nil
}

func (s *ConsultationService) UploadReport(doctorUserID uint, req *UploadReportRequest, file *multipart.FileHeader) (*models.ExaminationReport, error) {
	var report *models.ExaminationReport

	err := database.WithTransaction(func(tx *gorm.DB) error {
		var doctor models.Doctor
		if err := tx.Where("user_id = ?", doctorUserID).First(&doctor).Error; err != nil {
			return errors.New("医生信息不存在")
		}

		var consultation models.Consultation
		if err := tx.First(&consultation, req.ConsultationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("问诊记录不存在")
			}
			return err
		}

		var appointment models.Appointment
		if err := tx.First(&appointment, consultation.AppointmentID).Error; err != nil {
			return err
		}

		if appointment.DoctorID != doctor.ID {
			return errors.New("无权为该问诊记录上传检查报告")
		}

		uploadDir := config.AppConfig.UploadDir
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return fmt.Errorf("创建上传目录失败: %w", err)
		}

		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), utils.GenerateRandomString(8), ext)
		filePath := filepath.Join(uploadDir, fileName)

		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("打开文件失败: %w", err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("创建文件失败: %w", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return fmt.Errorf("保存文件失败: %w", err)
		}

		fileURL := fmt.Sprintf("/uploads/" + fileName)

		report = &models.ExaminationReport{
			ConsultationID: req.ConsultationID,
			ReportType:   req.ReportType,
			ReportName:   req.ReportName,
			FileURL:      fileURL,
			FileSize:     file.Size,
			ContentType:  file.Header.Get("Content-Type"),
			UploadedBy:   doctor.UserID,
			Findings:     req.Findings,
			Conclusion:   req.Conclusion,
		}

		if err := tx.Create(report).Error; err != nil {
			os.Remove(filePath)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ConsultationService) GetPatientConsultationHistory(patientUserID uint, page, pageSize int) ([]models.Consultation, int64, error) {
	var consultations []models.Consultation
	var total int64

	var patient models.Patient
	if err := s.db.Where("user_id = ?", patientUserID).First(&patient).Error; err != nil {
		return nil, 0, errors.New("患者信息不存在")
	}

	query := s.db.Model(&models.Consultation{}).
		Preload("Appointment.Doctor.User").
		Preload("Appointment.Department").
		Preload("Prescription.Items").
		Preload("Reports").
		Joins("JOIN appointments ON consultations.appointment_id = appointments.id").
		Where("appointments.patient_id = ?", patient.ID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).
		Order("appointments.appointment_date DESC, appointments.start_time DESC").
		Find(&consultations).Error; err != nil {
		return nil, 0, err
	}

	return consultations, total, nil
}
