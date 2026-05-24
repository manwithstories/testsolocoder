package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"business-registration-platform/config"
	"business-registration-platform/database"
	"business-registration-platform/models"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

type ExportParams struct {
	Type       string            `json:"type"`
	Fields     []string          `json:"fields"`
	StartDate  *time.Time        `json:"startDate"`
	EndDate    *time.Time        `json:"endDate"`
	Conditions map[string]string `json:"conditions"`
}

func (s *ExportService) CreateExportTask(userID uint, params *ExportParams) (*models.ExportTask, error) {
	paramsJSON, _ := json.Marshal(params)

	fileName := fmt.Sprintf("%s_%s.xlsx", params.Type, time.Now().Format("20060102150405"))
	exportPath := filepath.Join("exports", fileName)

	task := &models.ExportTask{
		UserID:   userID,
		Type:     params.Type,
		FileName: fileName,
		FilePath: exportPath,
		Status:   "pending",
		Params:   string(paramsJSON),
	}

	if err := database.DB.Create(task).Error; err != nil {
		return nil, err
	}

	go s.processExportTask(task)

	return task, nil
}

func (s *ExportService) processExportTask(task *models.ExportTask) {
	var params ExportParams
	json.Unmarshal([]byte(task.Params), &params)

	exportPath := filepath.Join(config.AppConfig.File.UploadPath, "exports")
	os.MkdirAll(exportPath, 0755)

	filePath := filepath.Join(exportPath, task.FileName)

	var err error
	switch params.Type {
	case "applications":
		err = s.exportApplications(filePath, &params)
	case "fees":
		err = s.exportFees(filePath, &params)
	case "agents":
		err = s.exportAgentRecords(filePath, &params)
	default:
		err = fmt.Errorf("unsupported export type: %s", params.Type)
	}

	now := time.Now()
	if err != nil {
		task.Status = "failed"
		task.ErrorMsg = err.Error()
	} else {
		task.Status = "completed"
		expiresAt := now.Add(7 * 24 * time.Hour)
		task.ExpiresAt = &expiresAt
	}
	database.DB.Save(task)
}

func (s *ExportService) exportApplications(filePath string, params *ExportParams) error {
	f := excelize.NewFile()
	sheetName := "Applications"
	f.SetCellValue(sheetName, "A1", "申请编号")
	f.SetCellValue(sheetName, "B1", "公司名称")
	f.SetCellValue(sheetName, "C1", "公司类型")
	f.SetCellValue(sheetName, "D1", "注册资本")
	f.SetCellValue(sheetName, "E1", "状态")
	f.SetCellValue(sheetName, "F1", "创业者")
	f.SetCellValue(sheetName, "G1", "代办专员")
	f.SetCellValue(sheetName, "H1", "进度")
	f.SetCellValue(sheetName, "I1", "创建时间")
	f.SetCellValue(sheetName, "J1", "完成时间")

	var applications []models.Application
	db := database.DB.Model(&models.Application{}).Preload("Entrepreneur").Preload("Agent")

	if params.StartDate != nil {
		db = db.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		db = db.Where("created_at <= ?", *params.EndDate)
	}
	if params.Conditions != nil {
		if status, ok := params.Conditions["status"]; ok {
			db = db.Where("status = ?", status)
		}
	}

	db.Find(&applications)

	for i, app := range applications {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), app.ApplicationNo)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), app.CompanyName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), app.CompanyType)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), app.RegisteredCapital)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), app.Status)
		if app.Entrepreneur != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), app.Entrepreneur.RealName)
		}
		if app.Agent != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), app.Agent.RealName)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("%d%%", app.ProgressPercent))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), app.CreatedAt.Format("2006-01-02 15:04:05"))
		if app.CompletedAt != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), app.CompletedAt.Format("2006-01-02 15:04:05"))
		}
	}

	return f.SaveAs(filePath)
}

func (s *ExportService) exportFees(filePath string, params *ExportParams) error {
	f := excelize.NewFile()
	sheetName := "Fees"
	f.SetCellValue(sheetName, "A1", "申请编号")
	f.SetCellValue(sheetName, "B1", "公司名称")
	f.SetCellValue(sheetName, "C1", "总金额")
	f.SetCellValue(sheetName, "D1", "优惠金额")
	f.SetCellValue(sheetName, "E1", "实付金额")
	f.SetCellValue(sheetName, "F1", "支付状态")
	f.SetCellValue(sheetName, "G1", "支付方式")
	f.SetCellValue(sheetName, "H1", "支付时间")
	f.SetCellValue(sheetName, "I1", "交易号")

	var fees []models.ApplicationFee
	db := database.DB.Model(&models.ApplicationFee{}).Preload("Application")

	if params.StartDate != nil {
		db = db.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		db = db.Where("created_at <= ?", *params.EndDate)
	}
	if params.Conditions != nil {
		if status, ok := params.Conditions["status"]; ok {
			db = db.Where("status = ?", status)
		}
	}

	db.Find(&fees)

	for i, fee := range fees {
		row := i + 2
		if fee.Application != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), fee.Application.ApplicationNo)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), fee.Application.CompanyName)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), fee.TotalAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), fee.DiscountAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), fee.PaidAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), fee.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), fee.PaymentMethod)
		if fee.PaymentTime != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), fee.PaymentTime.Format("2006-01-02 15:04:05"))
		}
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), fee.TransactionNo)
	}

	return f.SaveAs(filePath)
}

func (s *ExportService) exportAgentRecords(filePath string, params *ExportParams) error {
	f := excelize.NewFile()
	sheetName := "AgentRecords"
	f.SetCellValue(sheetName, "A1", "工号")
	f.SetCellValue(sheetName, "B1", "姓名")
	f.SetCellValue(sheetName, "C1", "联系方式")
	f.SetCellValue(sheetName, "D1", "专业领域")
	f.SetCellValue(sheetName, "E1", "已处理申请数")
	f.SetCellValue(sheetName, "F1", "当前处理数")
	f.SetCellValue(sheetName, "G1", "绩效评分")

	var agents []models.User
	database.DB.Where("role = ?", models.RoleAgent).Preload("AgentProfile").Find(&agents)

	for i, agent := range agents {
		row := i + 2
		if agent.AgentProfile != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), agent.AgentProfile.EmployeeNo)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), agent.AgentProfile.SpecialtyTags)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), agent.AgentProfile.TotalHandled)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), agent.AgentProfile.CurrentApps)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), agent.AgentProfile.PerformanceScore)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), agent.RealName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), agent.Phone)
	}

	return f.SaveAs(filePath)
}

func (s *ExportService) GetExportTasks(userID uint, page, pageSize int, status string) ([]models.ExportTask, int64, error) {
	var tasks []models.ExportTask
	var total int64

	db := database.DB.Model(&models.ExportTask{}).Where("user_id = ?", userID)

	if status != "" {
		db = db.Where("status = ?", status)
	}

	db.Count(&total)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks)

	return tasks, total, nil
}

func (s *ExportService) GetExportTask(taskID uint) (*models.ExportTask, error) {
	var task models.ExportTask
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *ExportService) DownloadExport(taskID uint, userID uint) (string, error) {
	var task models.ExportTask
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return "", fmt.Errorf("export task not found")
	}

	if task.UserID != userID {
		return "", fmt.Errorf("access denied")
	}

	if task.Status != "completed" {
		return "", fmt.Errorf("export is not ready")
	}

	if task.ExpiresAt != nil && time.Now().After(*task.ExpiresAt) {
		return "", fmt.Errorf("download link has expired")
	}

	filePath := filepath.Join(config.AppConfig.File.UploadPath, task.FilePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	}

	downloadLog := models.DownloadLog{
		ExportID:     task.ID,
		UserID:       userID,
		DownloadedAt: time.Now(),
	}
	database.DB.Create(&downloadLog)

	if !task.Downloaded {
		task.Downloaded = true
		now := time.Now()
		task.DownloadedAt = &now
		database.DB.Save(&task)
	}

	return filePath, nil
}

func (s *ExportService) GenerateDownloadToken(taskID uint) string {
	return uuid.New().String()
}
