package services

import (
	"errors"
	"time"

	"business-registration-platform/database"
	"business-registration-platform/models"
	"business-registration-platform/utils"
)

type AgentService struct{}

func NewAgentService() *AgentService {
	return &AgentService{}
}

type CreateAgentRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=100"`
	Password     string `json:"password" binding:"required,min=6,max=100"`
	RealName     string `json:"realName" binding:"required,max=100"`
	Email        string `json:"email" binding:"omitempty,email,max=100"`
	Phone        string `json:"phone" binding:"required,max=20"`
	EmployeeNo   string `json:"employeeNo" binding:"required,max=50"`
	SpecialtyTags string `json:"specialtyTags" binding:"required"`
	MaxApps      int    `json:"maxApps"`
}

type UpdateAgentProfileRequest struct {
	SpecialtyTags string `json:"specialtyTags"`
	MaxApps       int    `json:"maxApps"`
	WorkStartTime string `json:"workStartTime"`
	WorkEndTime   string `json:"workEndTime"`
	Status        string `json:"status"`
}

func (s *AgentService) GetAgentList(page, pageSize int, status string, keyword string) ([]models.User, int64, error) {
	var agents []models.User
	var total int64

	db := database.DB.Model(&models.User{}).Where("role = ?", models.RoleAgent)

	if status != "" {
		db = db.Joins("JOIN agent_profiles ON agent_profiles.user_id = users.id").
			Where("agent_profiles.status = ?", status)
	}

	if keyword != "" {
		db = db.Where("username LIKE ? OR real_name LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	db.Count(&total)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	db.Preload("AgentProfile").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&agents)

	return agents, total, nil
}

func (s *AgentService) GetAgentByID(id uint) (*models.User, error) {
	var agent models.User
	if err := database.DB.Preload("AgentProfile").First(&agent, id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (s *AgentService) CreateAgent(req *CreateAgentRequest) (*models.User, error) {
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	var existingProfile models.AgentProfile
	if err := database.DB.Where("employee_no = ?", req.EmployeeNo).First(&existingProfile).Error; err == nil {
		return nil, errors.New("employee number already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	tx := database.DB.Begin()

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		RealName: req.RealName,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     models.RoleAgent,
		Status:   models.UserStatusActive,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	profile := &models.AgentProfile{
		UserID:          user.ID,
		EmployeeNo:      req.EmployeeNo,
		SpecialtyTags:   req.SpecialtyTags,
		MaxApplications: req.MaxApps,
		CurrentApps:     0,
		WorkStartTime:   "09:00",
		WorkEndTime:     "18:00",
		Status:          "available",
	}

	if req.MaxApps <= 0 {
		profile.MaxApplications = 5
	}

	if err := tx.Create(profile).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return user, nil
}

func (s *AgentService) UpdateAgentProfile(agentID uint, req *UpdateAgentProfileRequest) error {
	updates := map[string]interface{}{}

	if req.SpecialtyTags != "" {
		updates["specialty_tags"] = req.SpecialtyTags
	}
	if req.MaxApps > 0 {
		updates["max_applications"] = req.MaxApps
	}
	if req.WorkStartTime != "" {
		updates["work_start_time"] = req.WorkStartTime
	}
	if req.WorkEndTime != "" {
		updates["work_end_time"] = req.WorkEndTime
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	return database.DB.Model(&models.AgentProfile{}).Where("user_id = ?", agentID).Updates(updates).Error
}

func (s *AgentService) DeleteAgent(agentID uint) error {
	tx := database.DB.Begin()

	if err := tx.Model(&models.AgentProfile{}).Where("user_id = ?", agentID).Delete(nil).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.User{}, agentID).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *AgentService) AutoAssignAgent(applicationID uint) (*models.User, error) {
	var agents []models.User
	database.DB.Joins("JOIN agent_profiles ON agent_profiles.user_id = users.id").
		Where("agent_profiles.status = ? AND agent_profiles.current_apps < agent_profiles.max_applications", "available").
		Preload("AgentProfile").
		Order("agent_profiles.current_apps ASC").
		Find(&agents)

	if len(agents) == 0 {
		return nil, errors.New("no available agent")
	}

	selectedAgent := agents[0]

	appService := NewApplicationService()
	if err := appService.AssignAgent(applicationID, selectedAgent.ID); err != nil {
		return nil, err
	}

	return &selectedAgent, nil
}

func (s *AgentService) GetAgentApplications(agentID uint, page, pageSize int, status string) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	db := database.DB.Model(&models.Application{}).Where("agent_id = ?", agentID)

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
	db.Preload("Entrepreneur").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&applications)

	return applications, total, nil
}

func (s *AgentService) GetAgentStats(agentID uint) (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	var totalHandled int64
	database.DB.Model(&models.Application{}).Where("agent_id = ?", agentID).Count(&totalHandled)
	stats["totalHandled"] = totalHandled

	var completedCount int64
	database.DB.Model(&models.Application{}).Where("agent_id = ? AND status = ?", agentID, models.AppStatusCompleted).Count(&completedCount)
	stats["completedCount"] = completedCount

	var inProgressCount int64
	database.DB.Model(&models.Application{}).Where("agent_id = ? AND status = ?", agentID, models.AppStatusProcessing).Count(&inProgressCount)
	stats["inProgressCount"] = inProgressCount

	var totalRevenue float64
	database.DB.Table("applications a").
		Select("COALESCE(SUM(f.paid_amount), 0)").
		Joins("JOIN application_fees f ON f.application_id = a.id").
		Where("a.agent_id = ? AND f.status = ?", agentID, models.FeeStatusPaid).
		Scan(&totalRevenue)
	stats["totalRevenue"] = totalRevenue

	return stats, nil
}

func (s *AgentService) UpdateAgentWorkSchedule(agentID uint, startTime, endTime string) error {
	return database.DB.Model(&models.AgentProfile{}).Where("user_id = ?", agentID).Updates(map[string]interface{}{
		"work_start_time": startTime,
		"work_end_time":   endTime,
	}).Error
}

func (s *AgentService) UpdateAgentMaxApps(agentID uint, maxApps int) error {
	return database.DB.Model(&models.AgentProfile{}).Where("user_id = ?", agentID).Update("max_applications", maxApps).Error
}

func (s *AgentService) GetAvailableAgents() ([]models.User, error) {
	var agents []models.User
	err := database.DB.Joins("JOIN agent_profiles ON agent_profiles.user_id = users.id").
		Where("agent_profiles.status = ? AND agent_profiles.current_apps < agent_profiles.max_applications", "available").
		Preload("AgentProfile").
		Find(&agents).Error
	return agents, err
}

func (s *AgentService) DecrementAgentCurrentApps(agentID uint) error {
	return database.DB.Model(&models.AgentProfile{}).
		Where("user_id = ? AND current_apps > 0", agentID).
		UpdateColumn("current_apps", database.DB.Raw("current_apps - 1")).Error
}

func (s *AgentService) UpdatePerformanceScore(agentID uint) error {
	var agentProfile models.AgentProfile
	if err := database.DB.Where("user_id = ?", agentID).First(&agentProfile).Error; err != nil {
		return err
	}

	if agentProfile.TotalHandled == 0 {
		return nil
	}

	var completedCount int64
	database.DB.Model(&models.Application{}).
		Where("agent_id = ? AND status = ?", agentID, models.AppStatusCompleted).
		Count(&completedCount)

	score := float64(completedCount) / float64(agentProfile.TotalHandled) * 100

	return database.DB.Model(&agentProfile).Update("performance_score", score).Error
}

func (s *AgentService) UpdateAgentStatus(agentID uint, status string) error {
	return database.DB.Model(&models.AgentProfile{}).Where("user_id = ?", agentID).Update("status", status).Error
}

func (s *AgentService) GetAgentPerformanceReport(agentID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	report := map[string]interface{}{}

	var totalApps int64
	db := database.DB.Model(&models.Application{}).Where("agent_id = ?", agentID)
	if !startDate.IsZero() {
		db = db.Where("created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		db = db.Where("created_at <= ?", endDate)
	}
	db.Count(&totalApps)
	report["totalApps"] = totalApps

	var completedApps int64
	db2 := database.DB.Model(&models.Application{}).Where("agent_id = ? AND status = ?", agentID, models.AppStatusCompleted)
	if !startDate.IsZero() {
		db2 = db2.Where("created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		db2 = db2.Where("created_at <= ?", endDate)
	}
	db2.Count(&completedApps)
	report["completedApps"] = completedApps

	var totalRevenue float64
	feeDB := database.DB.Table("applications a").
		Select("COALESCE(SUM(f.paid_amount), 0)").
		Joins("JOIN application_fees f ON f.application_id = a.id").
		Where("a.agent_id = ? AND f.status = ?", agentID, models.FeeStatusPaid)
	if !startDate.IsZero() {
		feeDB = feeDB.Where("a.created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		feeDB = feeDB.Where("a.created_at <= ?", endDate)
	}
	feeDB.Scan(&totalRevenue)
	report["totalRevenue"] = totalRevenue

	return report, nil
}
