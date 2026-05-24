package handlers

import (
	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApplicationHandler struct {
	DB *gorm.DB
}

func NewApplicationHandler(db *gorm.DB) *ApplicationHandler {
	return &ApplicationHandler{DB: db}
}

type ApplicationRequest struct {
	JobID       uint   `json:"job_id" binding:"required"`
	ResumeID    uint   `json:"resume_id" binding:"required"`
	CoverLetter string `json:"cover_letter"`
}

func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var req ApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var job models.Job
	if err := h.DB.First(&job, req.JobID).Error; err != nil {
		utils.NotFound(c, "Job not found")
		return
	}

	if job.Status != models.JobStatusOpen {
		utils.BadRequest(c, "Job is not open for applications")
		return
	}

	var resume models.Resume
	if err := h.DB.Where("id = ? AND jobseeker_id = ?", req.ResumeID, jobSeeker.ID).First(&resume).Error; err != nil {
		utils.NotFound(c, "Resume not found")
		return
	}

	var existingApp models.Application
	if h.DB.Where("job_id = ? AND jobseeker_id = ?", req.JobID, jobSeeker.ID).First(&existingApp).Error == nil {
		utils.BadRequest(c, "Already applied for this job")
		return
	}

	application := models.Application{
		JobID:       req.JobID,
		JobSeekerID: jobSeeker.ID,
		ResumeID:    req.ResumeID,
		Status:      models.ApplicationStatusPending,
		CoverLetter: req.CoverLetter,
	}

	if err := h.DB.Create(&application).Error; err != nil {
		utils.InternalError(c, "Failed to create application")
		return
	}

	var companyUser models.User
	h.DB.Joins("JOIN companies ON companies.user_id = users.id").
		Where("companies.id = ?", job.CompanyID).
		First(&companyUser)

	h.createNotification(h.DB, companyUser.ID, "New application received", "You have a new application for "+job.Title, models.NotificationTypeApplication, &application.ID)

	h.DB.Preload("Job").Preload("Resume").First(&application, application.ID)
	utils.Success(c, application)
}

func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	id := c.Param("id")

	var application models.Application
	if err := h.DB.Preload("Job").Preload("Job.Company").
		Preload("JobSeeker").Preload("JobSeeker.User").
		Preload("Resume").Preload("Resume.EducationList").
		Preload("Resume.WorkExperiences").Preload("Resume.Skills").
		Preload("Interviews").First(&application, id).Error; err != nil {
		utils.NotFound(c, "Application not found")
		return
	}

	utils.Success(c, application)
}

func (h *ApplicationHandler) ListJobSeekerApplications(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var applications []models.Application
	query := h.DB.Preload("Job").Preload("Job.Company").
		Preload("Resume").
		Where("jobseeker_id = ?", jobSeeker.ID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("created_at desc").Find(&applications)

	utils.Success(c, applications)
}

func (h *ApplicationHandler) ListCompanyApplications(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var applications []models.Application
	query := h.DB.Preload("Job").
		Preload("JobSeeker").Preload("JobSeeker.User").
		Preload("Resume").
		Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("jobs.company_id = ?", company.ID)

	if status := c.Query("status"); status != "" {
		query = query.Where("applications.status = ?", status)
	}

	if jobID := c.Query("job_id"); jobID != "" {
		query = query.Where("applications.job_id = ?", jobID)
	}

	query.Order("applications.created_at desc").Find(&applications)

	utils.Success(c, applications)
}

func (h *ApplicationHandler) UpdateApplicationStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var req struct {
		Status models.ApplicationStatus `json:"status" binding:"required,oneof=pending reviewed interview accepted rejected hold"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var application models.Application
	if err := h.DB.Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("applications.id = ? AND jobs.company_id = ?", id, company.ID).
		First(&application).Error; err != nil {
		utils.NotFound(c, "Application not found")
		return
	}

	h.DB.Model(&application).Update("status", req.Status)

	h.DB.Preload("Job").First(&application, application.ID)
	var jobSeekerUser models.User
	h.DB.Joins("JOIN job_seekers ON users.id = job_seekers.user_id").
		Where("job_seekers.id = ?", application.JobSeekerID).
		First(&jobSeekerUser)

	h.createNotification(h.DB, jobSeekerUser.ID,
		"Application status updated",
		"Your application for "+application.Job.Title+" has been updated to "+string(req.Status),
		models.NotificationTypeApplication, &application.ID)

	utils.Success(c, nil)
}

func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	result := h.DB.Where("id = ? AND jobseeker_id = ?", id, jobSeeker.ID).Delete(&models.Application{})
	if result.RowsAffected == 0 {
		utils.NotFound(c, "Application not found")
		return
	}

	utils.Success(c, nil)
}

func (h *ApplicationHandler) createNotification(db *gorm.DB, userID uint, title, message string, notifType models.NotificationType, relatedID *uint) {
	notification := models.Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		RelatedID: relatedID,
	}
	db.Create(&notification)
}

func CreateNotification(db *gorm.DB, userID uint, title, message string, notifType models.NotificationType, relatedID *uint) {
	notification := models.Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		RelatedID: relatedID,
	}
	db.Create(&notification)
}
