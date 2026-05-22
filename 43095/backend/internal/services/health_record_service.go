package services

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"
	"medical-platform/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type HealthRecordService struct {
	db *gorm.DB
}

func NewHealthRecordService() *HealthRecordService {
	return &HealthRecordService{
		db: database.GetDB(),
	}
}

type UpdateHealthRecordRequest struct {
	MedicalHistory map[string]interface{} `json:"medical_history"`
	Allergies      map[string]interface{} `json:"allergies"`
	Medications    map[string]interface{} `json:"medications"`
	Vaccinations   map[string]interface{} `json:"vaccinations"`
	FamilyHistory  string                 `json:"family_history"`
	LifeHabits     string                 `json:"life_habits"`
	Remarks        string                 `json:"remarks"`
}

type MedicalVisitHistory struct {
	Appointment   *models.Appointment       `json:"appointment"`
	Consultation  *models.Consultation      `json:"consultation,omitempty"`
	Prescription  *models.Prescription      `json:"prescription,omitempty"`
	Reports       []models.ExaminationReport `json:"reports,omitempty"`
}

func (s *HealthRecordService) GetHealthRecord(patientUserID uint, currentUserID uint, role models.UserRole) (*models.HealthRecord, error) {
	var patient models.Patient
	if err := s.db.Where("user_id = ?", patientUserID).First(&patient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("患者信息不存在")
		}
		return nil, err
	}

	if role == models.RolePatient {
		var currentPatient models.Patient
		if err := s.db.Where("user_id = ?", currentUserID).First(&currentPatient).Error; err != nil {
			return nil, errors.New("患者信息不存在")
		}
		if currentPatient.ID != patient.ID {
			return nil, errors.New("无权查看该患者的健康档案")
		}
	} else if role == models.RoleDoctor {
		var doctor models.Doctor
		if err := s.db.Where("user_id = ?", currentUserID).First(&doctor).Error; err != nil {
			return nil, errors.New("医生信息不存在")
		}
		var count int64
		if err := s.db.Model(&models.Appointment{}).
			Where("patient_id = ? AND doctor_id = ? AND status IN (?)",
				patient.ID, doctor.ID,
				[]models.AppointmentStatus{models.AppointmentConfirmed, models.AppointmentCompleted}).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, errors.New("无权查看该患者的健康档案")
		}
	} else if role != models.RoleAdmin {
		return nil, errors.New("无效的用户角色")
	}

	var healthRecord models.HealthRecord
	if err := s.db.Where("patient_id = ?", patient.ID).
		Preload("Patient.User").
		First(&healthRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			healthRecord = models.HealthRecord{
				PatientID: patient.ID,
			}
			if err := s.db.Create(&healthRecord).Error; err != nil {
				return nil, err
			}
			if err := s.db.Preload("Patient.User").First(&healthRecord, healthRecord.ID).Error; err != nil {
				return nil, err
			}
			return &healthRecord, nil
		}
		return nil, err
	}

	return &healthRecord, nil
}

func (s *HealthRecordService) UpdateHealthRecord(patientUserID uint, currentUserID uint, req *UpdateHealthRecordRequest) (*models.HealthRecord, error) {
	var patient models.Patient
	if err := s.db.Where("user_id = ?", patientUserID).First(&patient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("患者信息不存在")
		}
		return nil, err
	}

	var currentPatient models.Patient
	if err := s.db.Where("user_id = ?", currentUserID).First(&currentPatient).Error; err != nil {
		return nil, errors.New("患者信息不存在")
	}
	if currentPatient.ID != patient.ID {
		return nil, errors.New("无权修改该患者的健康档案")
	}

	var healthRecord models.HealthRecord
	if err := s.db.Where("patient_id = ?", patient.ID).First(&healthRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			healthRecord = models.HealthRecord{
				PatientID: patient.ID,
			}
			if err := s.db.Create(&healthRecord).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	updates := make(map[string]interface{})
	if req.MedicalHistory != nil {
		updates["medical_history"] = req.MedicalHistory
	}
	if req.Allergies != nil {
		updates["allergies"] = req.Allergies
	}
	if req.Medications != nil {
		updates["medications"] = req.Medications
	}
	if req.Vaccinations != nil {
		updates["vaccinations"] = req.Vaccinations
	}
	if req.FamilyHistory != "" {
		updates["family_history"] = req.FamilyHistory
	}
	if req.LifeHabits != "" {
		updates["life_habits"] = req.LifeHabits
	}
	if req.Remarks != "" {
		updates["remarks"] = req.Remarks
	}

	if err := s.db.Model(&healthRecord).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.db.Preload("Patient.User").First(&healthRecord, healthRecord.ID).Error; err != nil {
		return nil, err
	}

	return &healthRecord, nil
}

func (s *HealthRecordService) GetMedicalVisitHistory(patientUserID uint, currentUserID uint, role models.UserRole, page, pageSize int) ([]MedicalVisitHistory, int64, error) {
	var patient models.Patient
	if err := s.db.Where("user_id = ?", patientUserID).First(&patient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("患者信息不存在")
		}
		return nil, 0, err
	}

	if role == models.RolePatient {
		var currentPatient models.Patient
		if err := s.db.Where("user_id = ?", currentUserID).First(&currentPatient).Error; err != nil {
			return nil, 0, errors.New("患者信息不存在")
		}
		if currentPatient.ID != patient.ID {
			return nil, 0, errors.New("无权查看该患者的就诊历史")
		}
	} else if role == models.RoleDoctor {
		var doctor models.Doctor
		if err := s.db.Where("user_id = ?", currentUserID).First(&doctor).Error; err != nil {
			return nil, 0, errors.New("医生信息不存在")
		}
	} else if role != models.RoleAdmin {
		return nil, 0, errors.New("无效的用户角色")
	}

	var appointments []models.Appointment
	var total int64

	query := s.db.Model(&models.Appointment{}).
		Where("patient_id = ? AND status IN (?)",
			patient.ID,
			[]models.AppointmentStatus{models.AppointmentConfirmed, models.AppointmentCompleted})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).
		Preload("Doctor.User").
		Preload("Doctor.Department").
		Preload("Consultation").
		Preload("Consultation.Prescription.Items").
		Preload("Consultation.Reports").
		Order("appointment_date DESC, start_time DESC").
		Find(&appointments).Error; err != nil {
		return nil, 0, err
	}

	var history []MedicalVisitHistory
	for _, apt := range appointments {
		item := MedicalVisitHistory{
			Appointment: &apt,
		}
		if apt.Consultation != nil {
			item.Consultation = apt.Consultation
			item.Prescription = apt.Consultation.Prescription
			item.Reports = apt.Consultation.Reports
		}
		history = append(history, item)
	}

	return history, total, nil
}

func (s *HealthRecordService) ExportHealthRecordHTML(patientUserID uint, currentUserID uint, role models.UserRole) (string, error) {
	healthRecord, err := s.GetHealthRecord(patientUserID, currentUserID, role)
	if err != nil {
		return "", err
	}

	history, _, err := s.GetMedicalVisitHistory(patientUserID, currentUserID, role, 1, 100)
	if err != nil {
		return "", err
	}

	data := struct {
		HealthRecord *models.HealthRecord
		History      []MedicalVisitHistory
		ExportTime   string
	}{
		HealthRecord: healthRecord,
		History:      history,
		ExportTime:   utils.FormatDateTime(time.Now()),
	}

	const htmlTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>健康档案</title>
    <style>
        body { font-family: "Microsoft YaHei", sans-serif; margin: 40px; }
        .header { text-align: center; border-bottom: 2px solid #333; padding-bottom: 20px; margin-bottom: 30px; }
        .header h1 { color: #2c3e50; margin: 0; }
        .section { margin-bottom: 30px; }
        .section h2 { color: #3498db; border-bottom: 1px solid #ddd; padding-bottom: 10px; }
        .info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
        .info-item { display: flex; }
        .info-label { font-weight: bold; color: #555; min-width: 100px; }
        .info-value { color: #333; }
        .visit-card { border: 1px solid #ddd; border-radius: 8px; padding: 15px; margin-bottom: 15px; }
        .visit-header { background: #f8f9fa; padding: 10px; border-radius: 6px; margin-bottom: 10px; }
        .visit-header .date { font-weight: bold; color: #2c3e50; }
        .visit-header .doctor { color: #3498db; margin-left: 20px; }
        table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background: #f8f9fa; }
        .footer { text-align: center; margin-top: 40px; color: #999; font-size: 12px; }
        .json-content { background: #f8f9fa; padding: 10px; border-radius: 4px; white-space: pre-wrap; }
    </style>
</head>
<body>
    <div class="header">
        <h1>患者健康档案</h1>
        <p>导出时间：{{.ExportTime}}</p>
    </div>

    <div class="section">
        <h2>基本信息</h2>
        <div class="info-grid">
            <div class="info-item">
                <span class="info-label">姓名：</span>
                <span class="info-value">{{.HealthRecord.Patient.User.FullName}}</span>
            </div>
            <div class="info-item">
                <span class="info-label">性别：</span>
                <span class="info-value">{{.HealthRecord.Patient.User.Gender}}</span>
            </div>
            <div class="info-item">
                <span class="info-label">手机号：</span>
                <span class="info-value">{{.HealthRecord.Patient.User.Phone}}</span>
            </div>
            <div class="info-item">
                <span class="info-label">身份证号：</span>
                <span class="info-value">{{.HealthRecord.Patient.IDCardNo}}</span>
            </div>
        </div>
    </div>

    <div class="section">
        <h2>健康信息</h2>
        <div class="info-grid">
            <div class="info-item">
                <span class="info-label">家族病史：</span>
                <span class="info-value">{{if .HealthRecord.FamilyHistory}}{{.HealthRecord.FamilyHistory}}{{else}}无{{end}}</span>
            </div>
            <div class="info-item">
                <span class="info-label">生活习惯：</span>
                <span class="info-value">{{if .HealthRecord.LifeHabits}}{{.HealthRecord.LifeHabits}}{{else}}无{{end}}</span>
            </div>
            <div class="info-item">
                <span class="info-label">备注：</span>
                <span class="info-value">{{if .HealthRecord.Remarks}}{{.HealthRecord.Remarks}}{{else}}无{{end}}</span>
            </div>
        </div>
    </div>

    <div class="section">
        <h2>就诊历史</h2>
        {{if .History}}
            {{range .History}}
            <div class="visit-card">
                <div class="visit-header">
                    <span class="date">{{.Appointment.AppointmentDate.Format "2006-01-02"}} {{.Appointment.StartTime}}-{{.Appointment.EndTime}}</span>
                    <span class="doctor">{{.Appointment.Doctor.Department.Name}} - {{.Appointment.Doctor.User.FullName}} {{.Appointment.Doctor.Title}}</span>
                </div>
                {{if .Appointment.Symptoms}}
                <p><strong>主诉：</strong>{{.Appointment.Symptoms}}</p>
                {{end}}
                {{if .Consultation}}
                <p><strong>诊断：</strong>{{.Consultation.Diagnosis}}</p>
                {{if .Consultation.TreatmentPlan}}
                <p><strong>治疗方案：</strong>{{.Consultation.TreatmentPlan}}</p>
                {{end}}
                {{end}}
                {{if .Prescription}}
                <h4>处方</h4>
                <table>
                    <tr>
                        <th>药品名称</th>
                        <th>规格</th>
                        <th>用法用量</th>
                        <th>数量</th>
                    </tr>
                    {{range .Prescription.Items}}
                    <tr>
                        <td>{{.DrugName}}</td>
                        <td>{{.Specification}}</td>
                        <td>{{.Dosage}} {{.Frequency}}</td>
                        <td>{{.Quantity}}</td>
                    </tr>
                    {{end}}
                </table>
                {{end}}
                {{if .Reports}}
                <h4>检查报告</h4>
                <ul>
                    {{range .Reports}}
                    <li>{{.ReportType}} - {{.ReportName}}</li>
                    {{end}}
                </ul>
                {{end}}
            </div>
            {{end}}
        {{else}}
        <p>暂无就诊记录</p>
        {{end}}
    </div>

    <div class="footer">
        <p>本健康档案由系统自动生成，仅供参考</p>
    </div>
</body>
</html>`

	tmpl, err := template.New("health_record").Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("模板解析失败: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("模板渲染失败: %w", err)
	}

	return buf.String(), nil
}
