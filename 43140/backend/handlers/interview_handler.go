package handlers

import (
	"time"

	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InterviewHandler struct {
	DB *gorm.DB
}

func NewInterviewHandler(db *gorm.DB) *InterviewHandler {
	return &InterviewHandler{DB: db}
}

type InterviewRequest struct {
	ApplicationID   uint      `json:"application_id" binding:"required"`
	ScheduledAt     time.Time `json:"scheduled_at" binding:"required"`
	Duration        int       `json:"duration"`
	Location        string    `json:"location" binding:"required"`
	Interviewer     string    `json:"interviewer" binding:"required"`
	InterviewerEmail string   `json:"interviewer_email"`
	Notes           string    `json:"notes"`
}

func (h *InterviewHandler) CreateInterview(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var req InterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if req.ScheduledAt.Before(time.Now()) {
		utils.BadRequest(c, "Interview time must be in the future")
		return
	}

	var application models.Application
	if err := h.DB.Preload("Job").First(&application, req.ApplicationID).Error; err != nil {
		utils.NotFound(c, "Application not found")
		return
	}

	if application.Job.CompanyID != company.ID {
		utils.Forbidden(c, "You can only schedule interviews for your company's applications")
		return
	}

	interview := models.Interview{
		ApplicationID:   req.ApplicationID,
		ScheduledAt:     req.ScheduledAt,
		Duration:        req.Duration,
		Location:        req.Location,
		Interviewer:     req.Interviewer,
		InterviewerEmail: req.InterviewerEmail,
		Notes:           req.Notes,
		Status:          models.InterviewStatusScheduled,
	}

	if err := h.DB.Create(&interview).Error; err != nil {
		utils.InternalError(c, "Failed to create interview")
		return
	}

	h.DB.Model(&application).Update("status", models.ApplicationStatusInterview)

	var jobSeekerUser models.User
	h.DB.Joins("JOIN job_seekers ON users.id = job_seekers.user_id").
		Where("job_seekers.id = ?", application.JobSeekerID).
		First(&jobSeekerUser)

	CreateNotification(h.DB, jobSeekerUser.ID,
		"Interview scheduled",
		"You have an interview scheduled for "+application.Job.Title+" on "+req.ScheduledAt.Format("2006-01-02 15:04"),
		models.NotificationTypeInterview, &interview.ID)

	h.DB.Preload("Application").Preload("Application.Job").First(&interview, interview.ID)
	utils.Success(c, interview)
}

func (h *InterviewHandler) GetInterview(c *gin.Context) {
	id := c.Param("id")

	var interview models.Interview
	if err := h.DB.Preload("Application").Preload("Application.Job").
		Preload("Application.JobSeeker").Preload("Application.JobSeeker.User").
		Preload("Application.Resume").Preload("Review").
		First(&interview, id).Error; err != nil {
		utils.NotFound(c, "Interview not found")
		return
	}

	utils.Success(c, interview)
}

func (h *InterviewHandler) ListCompanyInterviews(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var interviews []models.Interview
	query := h.DB.Preload("Application").Preload("Application.Job").
		Preload("Application.JobSeeker").Preload("Application.JobSeeker.User").
		Joins("JOIN applications ON interviews.application_id = applications.id").
		Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("jobs.company_id = ?", company.ID)

	if status := c.Query("status"); status != "" {
		query = query.Where("interviews.status = ?", status)
	}

	if applicationID := c.Query("application_id"); applicationID != "" {
		query = query.Where("interviews.application_id = ?", applicationID)
	}

	query.Order("interviews.scheduled_at desc").Find(&interviews)

	utils.Success(c, interviews)
}

func (h *InterviewHandler) ListJobSeekerInterviews(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var interviews []models.Interview
	query := h.DB.Preload("Application").Preload("Application.Job").
		Preload("Application.Job.Company").Preload("Review").
		Joins("JOIN applications ON interviews.application_id = applications.id").
		Where("applications.jobseeker_id = ?", jobSeeker.ID)

	if status := c.Query("status"); status != "" {
		query = query.Where("interviews.status = ?", status)
	}

	query.Order("interviews.scheduled_at desc").Find(&interviews)

	utils.Success(c, interviews)
}

func (h *InterviewHandler) ConfirmInterview(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var interview models.Interview
	if err := h.DB.Joins("JOIN applications ON interviews.application_id = applications.id").
		Where("interviews.id = ? AND applications.jobseeker_id = ?", id, jobSeeker.ID).
		First(&interview).Error; err != nil {
		utils.NotFound(c, "Interview not found")
		return
	}

	if interview.Status != models.InterviewStatusScheduled {
		utils.BadRequest(c, "Interview is not in scheduled status")
		return
	}

	h.DB.Model(&interview).Update("status", models.InterviewStatusConfirmed)

	h.DB.Preload("Application").First(&interview, interview.ID)
	var companyUser models.User
	h.DB.Joins("JOIN companies ON users.id = companies.user_id").
		Where("companies.id = ?", interview.Application.Job.CompanyID).
		First(&companyUser)

	CreateNotification(h.DB, companyUser.ID,
		"Interview confirmed",
		"The candidate has confirmed the interview scheduled for "+interview.ScheduledAt.Format("2006-01-02 15:04"),
		models.NotificationTypeInterview, &interview.ID)

	utils.Success(c, nil)
}

func (h *InterviewHandler) DeclineInterview(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var interview models.Interview
	if err := h.DB.Joins("JOIN applications ON interviews.application_id = applications.id").
		Where("interviews.id = ? AND applications.jobseeker_id = ?", id, jobSeeker.ID).
		First(&interview).Error; err != nil {
		utils.NotFound(c, "Interview not found")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	h.DB.Model(&interview).Updates(map[string]interface{}{
		"status": models.InterviewStatusDeclined,
		"notes":  interview.Notes + "\nDecline reason: " + req.Reason,
	})

	h.DB.Preload("Application").First(&interview, interview.ID)
	var companyUser models.User
	h.DB.Joins("JOIN companies ON users.id = companies.user_id").
		Where("companies.id = ?", interview.Application.Job.CompanyID).
		First(&companyUser)

	CreateNotification(h.DB, companyUser.ID,
		"Interview declined",
		"The candidate has declined the interview scheduled for "+interview.ScheduledAt.Format("2006-01-02 15:04"),
		models.NotificationTypeInterview, &interview.ID)

	utils.Success(c, nil)
}

func (h *InterviewHandler) UpdateInterview(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var interview models.Interview
	if err := h.DB.Joins("JOIN applications ON interviews.application_id = applications.id").
		Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("interviews.id = ? AND jobs.company_id = ?", id, company.ID).
		First(&interview).Error; err != nil {
		utils.NotFound(c, "Interview not found")
		return
	}

	var req struct {
		ScheduledAt      *time.Time `json:"scheduled_at"`
		Duration         *int       `json:"duration"`
		Location         *string    `json:"location"`
		Interviewer      *string    `json:"interviewer"`
		InterviewerEmail *string    `json:"interviewer_email"`
		Notes            *string    `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.ScheduledAt != nil {
		updates["scheduled_at"] = *req.ScheduledAt
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Interviewer != nil {
		updates["interviewer"] = *req.Interviewer
	}
	if req.InterviewerEmail != nil {
		updates["interviewer_email"] = *req.InterviewerEmail
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}

	h.DB.Model(&interview).Updates(updates)

	utils.Success(c, nil)
}

func (h *InterviewHandler) CancelInterview(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var interview models.Interview
	if err := h.DB.Joins("JOIN applications ON interviews.application_id = applications.id").
		Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("interviews.id = ? AND jobs.company_id = ?", id, company.ID).
		First(&interview).Error; err != nil {
		utils.NotFound(c, "Interview not found")
		return
	}

	h.DB.Model(&interview).Update("status", models.InterviewStatusCancelled)
	utils.Success(c, nil)
}
