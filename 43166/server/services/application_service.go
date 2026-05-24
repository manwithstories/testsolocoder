package services

import (
	"errors"
	"fmt"
	"time"

	"business-registration-platform/database"
	"business-registration-platform/models"

	"github.com/google/uuid"
)

type ApplicationService struct{}

func NewApplicationService() *ApplicationService {
	return &ApplicationService{}
}

type CreateApplicationRequest struct {
	EntrepreneurID    uint              `json:"entrepreneurId"`
	CompanyName       string            `json:"companyName" binding:"required,max=200"`
	CompanyType       models.CompanyType `json:"companyType" binding:"required,oneof=llc joint_stock sole partnership"`
	RegisteredCapital float64           `json:"registeredCapital" binding:"required,gt=0"`
	BusinessScope     string            `json:"businessScope" binding:"required"`
	RegisteredAddress string            `json:"registeredAddress" binding:"required,max=500"`
	ShareholderInfo   string            `json:"shareholderInfo" binding:"required"`
}

type UpdateApplicationRequest struct {
	CompanyName       string  `json:"companyName"`
	CompanyType       string  `json:"companyType"`
	RegisteredCapital float64 `json:"registeredCapital"`
	BusinessScope     string  `json:"businessScope"`
	RegisteredAddress string  `json:"registeredAddress"`
	ShareholderInfo   string  `json:"shareholderInfo"`
}

type ApplicationListQuery struct {
	Page           int                    `json:"page" form:"page"`
	PageSize       int                    `json:"pageSize" form:"pageSize"`
	Status         models.ApplicationStatus `json:"status" form:"status"`
	Keyword        string                 `json:"keyword" form:"keyword"`
	EntrepreneurID *uint                  `json:"entrepreneurId" form:"entrepreneurId"`
	AgentID        *uint                  `json:"agentId" form:"agentId"`
}

func (s *ApplicationService) CreateApplication(req *CreateApplicationRequest) (*models.Application, error) {
	applicationNo := fmt.Sprintf("APP%s%s", time.Now().Format("20060102"), uuid.New().String()[:8])

	application := &models.Application{
		ApplicationNo:    applicationNo,
		EntrepreneurID:   req.EntrepreneurID,
		CompanyName:      req.CompanyName,
		CompanyType:      req.CompanyType,
		RegisteredCapital: req.RegisteredCapital,
		BusinessScope:    req.BusinessScope,
		RegisteredAddress: req.RegisteredAddress,
		ShareholderInfo:  req.ShareholderInfo,
		Status:           models.AppStatusDraft,
		ProgressPercent:  0,
	}

	if err := database.DB.Create(application).Error; err != nil {
		return nil, err
	}

	return application, nil
}

func (s *ApplicationService) SubmitApplication(applicationID uint) error {
	var application models.Application
	if err := database.DB.First(&application, applicationID).Error; err != nil {
		return errors.New("application not found")
	}

	if application.Status != models.AppStatusDraft {
		return errors.New("application is not in draft status")
	}

	now := time.Now()
	application.Status = models.AppStatusPaymentPending
	application.SubmittedAt = &now

	return database.DB.Save(&application).Error
}

func (s *ApplicationService) AssignAgent(applicationID, agentID uint) error {
	var application models.Application
	if err := database.DB.First(&application, applicationID).Error; err != nil {
		return errors.New("application not found")
	}

	var agentProfile models.AgentProfile
	if err := database.DB.Where("user_id = ?", agentID).First(&agentProfile).Error; err != nil {
		return errors.New("agent not found")
	}

	if agentProfile.CurrentApps >= agentProfile.MaxApplications {
		return errors.New("agent has reached maximum application limit")
	}

	tx := database.DB.Begin()

	application.AgentID = &agentID
	application.Status = models.AppStatusProcessing
	application.CurrentStep = string(models.StepTypeNaming)
	if err := tx.Save(&application).Error; err != nil {
		tx.Rollback()
		return err
	}

	agentProfile.CurrentApps++
	if err := tx.Save(&agentProfile).Error; err != nil {
		tx.Rollback()
		return err
	}

	steps := []models.ProcessStep{
		{ApplicationID: applicationID, StepType: models.StepTypeNaming, StepName: "核名", StepOrder: 1, Status: models.StepStatusInProgress},
		{ApplicationID: applicationID, StepType: models.StepTypeRegistration, StepName: "工商登记", StepOrder: 2, Status: models.StepStatusPending},
		{ApplicationID: applicationID, StepType: models.StepTypeTax, StepName: "税务登记", StepOrder: 3, Status: models.StepStatusPending},
		{ApplicationID: applicationID, StepType: models.StepTypeBank, StepName: "银行开户", StepOrder: 4, Status: models.StepStatusPending},
		{ApplicationID: applicationID, StepType: models.StepTypeSeal, StepName: "刻章备案", StepOrder: 5, Status: models.StepStatusPending},
		{ApplicationID: applicationID, StepType: models.StepTypeCompletion, StepName: "完成交付", StepOrder: 6, Status: models.StepStatusPending},
	}
	if err := tx.Create(&steps).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *ApplicationService) GetApplicationByID(id uint) (*models.Application, error) {
	var application models.Application
	if err := database.DB.
		Preload("Entrepreneur").
		Preload("Agent").
		Preload("ProcessSteps").
		Preload("Fee").
		Preload("Fee.FeeItems").
		First(&application, id).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (s *ApplicationService) GetApplicationList(query *ApplicationListQuery) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	db := database.DB.Model(&models.Application{})

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.Keyword != "" {
		db = db.Where("company_name LIKE ? OR application_no LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.EntrepreneurID != nil {
		db = db.Where("entrepreneur_id = ?", *query.EntrepreneurID)
	}

	if query.AgentID != nil {
		db = db.Where("agent_id = ?", *query.AgentID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Entrepreneur").Preload("Agent").
		Order("created_at DESC").
		Offset(offset).Limit(query.PageSize).
		Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

func (s *ApplicationService) UpdateApplication(id uint, req *UpdateApplicationRequest) error {
	updates := map[string]interface{}{}
	if req.CompanyName != "" {
		updates["company_name"] = req.CompanyName
	}
	if req.CompanyType != "" {
		updates["company_type"] = req.CompanyType
	}
	if req.RegisteredCapital > 0 {
		updates["registered_capital"] = req.RegisteredCapital
	}
	if req.BusinessScope != "" {
		updates["business_scope"] = req.BusinessScope
	}
	if req.RegisteredAddress != "" {
		updates["registered_address"] = req.RegisteredAddress
	}
	if req.ShareholderInfo != "" {
		updates["shareholder_info"] = req.ShareholderInfo
	}

	return database.DB.Model(&models.Application{}).Where("id = ?", id).Updates(updates).Error
}

func (s *ApplicationService) UploadMaterials(applicationID uint, field, filePath string) error {
	updates := map[string]interface{}{}
	updates[field] = filePath
	return database.DB.Model(&models.Application{}).Where("id = ?", applicationID).Updates(updates).Error
}

func (s *ApplicationService) ReviewApplication(applicationID uint, approved bool, comments string) error {
	var application models.Application
	if err := database.DB.First(&application, applicationID).Error; err != nil {
		return errors.New("application not found")
	}

	if application.Status != models.AppStatusPendingReview {
		return errors.New("application is not in review status")
	}

	if approved {
		application.Status = models.AppStatusProcessing
		application.CurrentStep = string(models.StepTypeNaming)
	} else {
		application.Status = models.AppStatusRejected
		now := time.Now()
		application.RejectedAt = &now
	}
	application.ReviewComments = comments

	return database.DB.Save(&application).Error
}

func (s *ApplicationService) CancelApplication(applicationID uint) error {
	var application models.Application
	if err := database.DB.First(&application, applicationID).Error; err != nil {
		return errors.New("application not found")
	}

	if application.Status == models.AppStatusCompleted || application.Status == models.AppStatusCancelled {
		return errors.New("cannot cancel this application")
	}

	application.Status = models.AppStatusCancelled
	return database.DB.Save(&application).Error
}
