package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"
	"medical-platform/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type PaymentService struct {
	db *gorm.DB
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		db: database.GetDB(),
	}
}

type CreatePaymentRequest struct {
	AppointmentID   uint                `json:"appointment_id" binding:"required"`
	RegistrationFee float64             `json:"registration_fee"`
	ConsultationFee float64             `json:"consultation_fee"`
	DrugFee         float64             `json:"drug_fee"`
	ExaminationFee  float64             `json:"examination_fee"`
	OtherFee        float64             `json:"other_fee"`
	Method          models.PaymentMethod `json:"method" binding:"required,oneof=wechat alipay card cash"`
	Notes           string              `json:"notes"`
}

type UpdatePaymentStatusRequest struct {
	Status models.PaymentStatus `json:"status" binding:"required,oneof=pending paid failed refunded"`
	Notes  string               `json:"notes"`
}

type PaymentListQuery struct {
	PatientID  uint      `form:"patient_id"`
	DoctorID   uint      `form:"doctor_id"`
	Status     string    `form:"status"`
	StartDate  string    `form:"start_date"`
	EndDate    string    `form:"end_date"`
	Page       int       `form:"page,default=1"`
	PageSize   int       `form:"page_size,default=10"`
}

type FeeReportQuery struct {
	StartDate    string `form:"start_date"`
	EndDate      string `form:"end_date"`
	DepartmentID uint   `form:"department_id"`
	DoctorID     uint   `form:"doctor_id"`
	GroupBy      string `form:"group_by,default=date"`
}

type FeeReportItem struct {
	Date           string  `json:"date,omitempty"`
	DepartmentID   uint    `json:"department_id,omitempty"`
	DepartmentName string  `json:"department_name,omitempty"`
	DoctorID       uint    `json:"doctor_id,omitempty"`
	DoctorName     string  `json:"doctor_name,omitempty"`
	TotalAmount    float64 `json:"total_amount"`
	PaymentCount   int64   `json:"payment_count"`
	RegistrationFee float64 `json:"registration_fee"`
	ConsultationFee float64 `json:"consultation_fee"`
	DrugFee        float64 `json:"drug_fee"`
	ExaminationFee float64 `json:"examination_fee"`
	OtherFee       float64 `json:"other_fee"`
}

func (s *PaymentService) CreatePayment(req CreatePaymentRequest) (*models.Payment, error) {
	var existingPayment models.Payment
	err := s.db.Where("appointment_id = ?", req.AppointmentID).First(&existingPayment).Error
	if err == nil {
		return nil, errors.New("该预约已存在支付记录")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var appointment models.Appointment
	if err := s.db.First(&appointment, req.AppointmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预约不存在")
		}
		return nil, err
	}

	totalAmount := req.RegistrationFee + req.ConsultationFee + req.DrugFee + req.ExaminationFee + req.OtherFee
	if totalAmount <= 0 {
		return nil, errors.New("支付金额必须大于0")
	}

	transactionNo := utils.GenerateOrderNo("PAY")

	payment := &models.Payment{
		AppointmentID:   req.AppointmentID,
		TransactionNo:   transactionNo,
		RegistrationFee: req.RegistrationFee,
		ConsultationFee: req.ConsultationFee,
		DrugFee:         req.DrugFee,
		ExaminationFee:  req.ExaminationFee,
		OtherFee:        req.OtherFee,
		TotalAmount:     totalAmount,
		Status:          models.PaymentPending,
		Method:          req.Method,
		Notes:           req.Notes,
	}

	if err := s.db.Create(payment).Error; err != nil {
		return nil, err
	}

	return s.GetPaymentByID(payment.ID)
}

func (s *PaymentService) GetPaymentByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := s.db.Preload("Appointment.Patient.User").
		Preload("Appointment.Doctor.User").
		Preload("Appointment.Doctor.Department").
		First(&payment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("支付记录不存在")
		}
		return nil, err
	}
	return &payment, nil
}

func (s *PaymentService) UpdatePaymentStatus(id uint, req UpdatePaymentStatusRequest) (*models.Payment, error) {
	payment, err := s.GetPaymentByID(id)
	if err != nil {
		return nil, err
	}

	if payment.Status == req.Status {
		return payment, nil
	}

	updates := make(map[string]interface{})
	updates["status"] = req.Status
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}

	now := time.Now()
	if req.Status == models.PaymentPaid {
		updates["paid_at"] = &now
	} else if req.Status == models.PaymentRefunded {
		updates["refunded_at"] = &now
	}

	if err := s.db.Model(payment).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetPaymentByID(id)
}

func (s *PaymentService) GetPaymentList(query PaymentListQuery) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	db := s.db.Model(&models.Payment{}).
		Preload("Appointment.Patient.User").
		Preload("Appointment.Doctor.User").
		Preload("Appointment.Doctor.Department")

	if query.PatientID > 0 {
		db = db.Joins("JOIN appointments ON payments.appointment_id = appointments.id").
			Where("appointments.patient_id = ?", query.PatientID)
	}

	if query.DoctorID > 0 {
		db = db.Joins("JOIN appointments ON payments.appointment_id = appointments.id").
			Where("appointments.doctor_id = ?", query.DoctorID)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", query.StartDate)
		if err == nil {
			db = db.Where("DATE(created_at) >= ?", startDate.Format("2006-01-02"))
		}
	}

	if query.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", query.EndDate)
		if err == nil {
			db = db.Where("DATE(created_at) <= ?", endDate.Format("2006-01-02"))
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("created_at DESC").
		Scopes(database.Paginate(query.Page, query.PageSize)).
		Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (s *PaymentService) GenerateFeeReport(query FeeReportQuery) ([]FeeReportItem, error) {
	var results []FeeReportItem

	db := s.db.Model(&models.Payment{}).
		Where("status = ?", models.PaymentPaid)

	if query.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", query.StartDate)
		if err == nil {
			db = db.Where("DATE(paid_at) >= ?", startDate.Format("2006-01-02"))
		}
	}

	if query.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", query.EndDate)
		if err == nil {
			db = db.Where("DATE(paid_at) <= ?", endDate.Format("2006-01-02"))
		}
	}

	if query.DepartmentID > 0 {
		db = db.Joins("JOIN appointments ON payments.appointment_id = appointments.id").
			Joins("JOIN doctors ON appointments.doctor_id = doctors.id").
			Where("doctors.department_id = ?", query.DepartmentID)
	}

	if query.DoctorID > 0 {
		db = db.Joins("JOIN appointments ON payments.appointment_id = appointments.id").
			Where("appointments.doctor_id = ?", query.DoctorID)
	}

	switch query.GroupBy {
	case "date":
		db = db.Select(`
			DATE(paid_at) as date,
			SUM(total_amount) as total_amount,
			COUNT(*) as payment_count,
			SUM(registration_fee) as registration_fee,
			SUM(consultation_fee) as consultation_fee,
			SUM(drug_fee) as drug_fee,
			SUM(examination_fee) as examination_fee,
			SUM(other_fee) as other_fee
		`).Group("DATE(paid_at)").Order("date DESC")
	case "department":
		db = db.Joins("JOIN appointments ON payments.appointment_id = appointments.id").
			Joins("JOIN doctors ON appointments.doctor_id = doctors.id").
			Joins("JOIN departments ON doctors.department_id = departments.id").
			Select(`
				departments.id as department_id,
				departments.name as department_name,
				SUM(payments.total_amount) as total_amount,
				COUNT(*) as payment_count,
				SUM(payments.registration_fee) as registration_fee,
				SUM(payments.consultation_fee) as consultation_fee,
				SUM(payments.drug_fee) as drug_fee,
				SUM(payments.examination_fee) as examination_fee,
				SUM(payments.other_fee) as other_fee
			`).Group("departments.id, departments.name").Order("total_amount DESC")
	case "doctor":
		db = db.Joins("JOIN appointments ON payments.appointment_id = appointments.id").
			Joins("JOIN doctors ON appointments.doctor_id = doctors.id").
			Joins("JOIN users ON doctors.user_id = users.id").
			Select(`
				doctors.id as doctor_id,
				users.full_name as doctor_name,
				SUM(payments.total_amount) as total_amount,
				COUNT(*) as payment_count,
				SUM(payments.registration_fee) as registration_fee,
				SUM(payments.consultation_fee) as consultation_fee,
				SUM(payments.drug_fee) as drug_fee,
				SUM(payments.examination_fee) as examination_fee,
				SUM(payments.other_fee) as other_fee
			`).Group("doctors.id, users.full_name").Order("total_amount DESC")
	default:
		return nil, errors.New("无效的分组方式")
	}

	if err := db.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (s *PaymentService) ExportFeeReportCSV(query FeeReportQuery) (string, error) {
	report, err := s.GenerateFeeReport(query)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll("./exports", 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("fee_report_%s.csv", time.Now().Format("20060102150405"))
	filepath := filepath.Join("./exports", filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var headers []string
	switch query.GroupBy {
	case "date":
		headers = []string{"日期", "总金额", "支付笔数", "挂号费", "诊疗费", "药品费", "检查费", "其他费用"}
	case "department":
		headers = []string{"科室ID", "科室名称", "总金额", "支付笔数", "挂号费", "诊疗费", "药品费", "检查费", "其他费用"}
	case "doctor":
		headers = []string{"医生ID", "医生姓名", "总金额", "支付笔数", "挂号费", "诊疗费", "药品费", "检查费", "其他费用"}
	}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	for _, item := range report {
		var row []string
		switch query.GroupBy {
		case "date":
			row = []string{
				item.Date,
				fmt.Sprintf("%.2f", item.TotalAmount),
				fmt.Sprintf("%d", item.PaymentCount),
				fmt.Sprintf("%.2f", item.RegistrationFee),
				fmt.Sprintf("%.2f", item.ConsultationFee),
				fmt.Sprintf("%.2f", item.DrugFee),
				fmt.Sprintf("%.2f", item.ExaminationFee),
				fmt.Sprintf("%.2f", item.OtherFee),
			}
		case "department":
			row = []string{
				strconv.FormatUint(uint64(item.DepartmentID), 10),
				item.DepartmentName,
				fmt.Sprintf("%.2f", item.TotalAmount),
				fmt.Sprintf("%d", item.PaymentCount),
				fmt.Sprintf("%.2f", item.RegistrationFee),
				fmt.Sprintf("%.2f", item.ConsultationFee),
				fmt.Sprintf("%.2f", item.DrugFee),
				fmt.Sprintf("%.2f", item.ExaminationFee),
				fmt.Sprintf("%.2f", item.OtherFee),
			}
		case "doctor":
			row = []string{
				strconv.FormatUint(uint64(item.DoctorID), 10),
				item.DoctorName,
				fmt.Sprintf("%.2f", item.TotalAmount),
				fmt.Sprintf("%d", item.PaymentCount),
				fmt.Sprintf("%.2f", item.RegistrationFee),
				fmt.Sprintf("%.2f", item.ConsultationFee),
				fmt.Sprintf("%.2f", item.DrugFee),
				fmt.Sprintf("%.2f", item.ExaminationFee),
				fmt.Sprintf("%.2f", item.OtherFee),
			}
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	return filepath, nil
}
