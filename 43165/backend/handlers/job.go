package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"temp-staff-platform/database"
	"temp-staff-platform/models"
	"temp-staff-platform/utils"
)

type JobHandler struct{}

func NewJobHandler() *JobHandler {
	return &JobHandler{}
}

type CreateJobRequest struct {
	Title         string    `json:"title" binding:"required,max=100"`
	Description   string    `json:"description" binding:"required"`
	ActivityType  string    `json:"activity_type"`
	Position      string    `json:"position" binding:"required"`
	Location      string    `json:"location" binding:"required"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	SalaryPerHour float64   `json:"salary_per_hour" binding:"required,gt=0"`
	SalaryType    string    `json:"salary_type"`
	WorkHours     string    `json:"work_hours"`
	Headcount     int       `json:"headcount" binding:"required,gt=0"`
	Requirements  string    `json:"requirements"`
	Benefits      string    `json:"benefits"`
	ContactPerson string    `json:"contact_person"`
	ContactPhone  string    `json:"contact_phone"`
	Tags          string    `json:"tags"`
	IsUrgent      bool      `json:"is_urgent"`
}

type UpdateJobRequest struct {
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	ActivityType  string  `json:"activity_type"`
	Position      string  `json:"position"`
	Location      string  `json:"location"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	SalaryPerHour float64 `json:"salary_per_hour"`
	SalaryType    string  `json:"salary_type"`
	WorkHours     string  `json:"work_hours"`
	Headcount     int     `json:"headcount"`
	Requirements  string  `json:"requirements"`
	Benefits      string  `json:"benefits"`
	ContactPerson string  `json:"contact_person"`
	ContactPhone  string  `json:"contact_phone"`
	Tags          string  `json:"tags"`
	IsUrgent      bool    `json:"is_urgent"`
	Status        string  `json:"status"`
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	if req.EndDate.Before(req.StartDate) {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "End date must be after start date",
		})
		return
	}

	job := models.JobPosting{
		EmployerID:    userID.(uuid.UUID),
		Title:         req.Title,
		Description:   req.Description,
		ActivityType:  req.ActivityType,
		Position:      req.Position,
		Location:      req.Location,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		SalaryPerHour: req.SalaryPerHour,
		SalaryType:    req.SalaryType,
		WorkHours:     req.WorkHours,
		Headcount:     req.Headcount,
		Requirements:  req.Requirements,
		Benefits:      req.Benefits,
		ContactPerson: req.ContactPerson,
		ContactPhone:  req.ContactPhone,
		Tags:          req.Tags,
		IsUrgent:      req.IsUrgent,
		Status:        "recruiting",
	}

	if err := database.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to create job posting: " + err.Error(),
		})
		return
	}

	database.DB.Preload("Employer").First(&job, job.ID)

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Job posting created successfully",
		Data:    job,
	})
}

func (h *JobHandler) GetJobs(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var jobs []models.JobPosting
	var total int64

	query := database.DB.Model(&models.JobPosting{})

	if pagination.Keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR position ILIKE ? OR location ILIKE ? OR tags ILIKE ?",
			"%"+pagination.Keyword+"%",
			"%"+pagination.Keyword+"%",
			"%"+pagination.Keyword+"%",
			"%"+pagination.Keyword+"%",
			"%"+pagination.Keyword+"%",
		)
	}

	if activityType := c.Query("activity_type"); activityType != "" {
		query = query.Where("activity_type = ?", activityType)
	}
	if location := c.Query("location"); location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if minSalary := c.Query("min_salary"); minSalary != "" {
		query = query.Where("salary_per_hour >= ?", minSalary)
	}
	if isUrgent := c.Query("is_urgent"); isUrgent == "true" {
		query = query.Where("is_urgent = ?", true)
	}

	query.Count(&total)

	sortField := pagination.SortBy
	if sortField == "" {
		sortField = "created_at"
	}
	validSortFields := map[string]bool{
		"created_at":      true,
		"start_date":      true,
		"salary_per_hour": true,
		"is_urgent":       true,
	}
	if !validSortFields[sortField] {
		sortField = "created_at"
	}

	query.Preload("Employer").
		Order(sortField + " " + pagination.SortOrder).
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&jobs)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       jobs,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	jobID := id.(uuid.UUID)

	var job models.JobPosting
	if err := database.DB.Preload("Employer").First(&job, jobID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Job posting not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data:    job,
	})
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	jobID := id.(uuid.UUID)
	userID, _ := c.Get("user_id")

	var job models.JobPosting
	if err := database.DB.First(&job, jobID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Job posting not found",
		})
		return
	}

	if job.EmployerID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    403,
			Message: "You don't have permission to update this job posting",
		})
		return
	}

	var req UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ActivityType != "" {
		updates["activity_type"] = req.ActivityType
	}
	if req.Position != "" {
		updates["position"] = req.Position
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.SalaryPerHour > 0 {
		updates["salary_per_hour"] = req.SalaryPerHour
	}
	if req.SalaryType != "" {
		updates["salary_type"] = req.SalaryType
	}
	if req.WorkHours != "" {
		updates["work_hours"] = req.WorkHours
	}
	if req.Headcount > 0 {
		updates["headcount"] = req.Headcount
	}
	if req.Requirements != "" {
		updates["requirements"] = req.Requirements
	}
	if req.Benefits != "" {
		updates["benefits"] = req.Benefits
	}
	if req.ContactPerson != "" {
		updates["contact_person"] = req.ContactPerson
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}
	updates["is_urgent"] = req.IsUrgent
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		database.DB.Model(&job).Updates(updates)
	}

	database.DB.Preload("Employer").First(&job, jobID)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Job posting updated successfully",
		Data:    job,
	})
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	jobID := id.(uuid.UUID)
	userID, _ := c.Get("user_id")

	var job models.JobPosting
	if err := database.DB.First(&job, jobID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Job posting not found",
		})
		return
	}

	if job.EmployerID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    403,
			Message: "You don't have permission to delete this job posting",
		})
		return
	}

	database.DB.Delete(&job)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Job posting deleted successfully",
	})
}

func (h *JobHandler) GetMyJobs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	pagination := utils.GetPagination(c)

	var jobs []models.JobPosting
	var total int64

	query := database.DB.Model(&models.JobPosting{}).Where("employer_id = ?", userID.(uuid.UUID))

	if pagination.Keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?",
			"%"+pagination.Keyword+"%",
			"%"+pagination.Keyword+"%",
		)
	}

	query.Count(&total)
	query.Preload("Employer").
		Order("created_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&jobs)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       jobs,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}

type ApplyJobRequest struct {
	Message string `json:"message"`
}

func (h *JobHandler) ApplyJob(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	jobID := id.(uuid.UUID)
	userID, _ := c.Get("user_id")

	var req ApplyJobRequest
	c.ShouldBindJSON(&req)

	var job models.JobPosting
	if err := database.DB.First(&job, jobID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Job posting not found",
		})
		return
	}

	if job.Status != "recruiting" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "This job is no longer recruiting",
		})
		return
	}

	var existingApp models.JobApplication
	if err := database.DB.Where("job_id = ? AND temporary_id = ?", jobID, userID.(uuid.UUID)).First(&existingApp).Error; err == nil {
		c.JSON(http.StatusConflict, models.Response{
			Code:    409,
			Message: "You have already applied for this job",
		})
		return
	}

	if job.HiredCount >= job.Headcount {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "This job has reached its headcount limit",
		})
		return
	}

	application := models.JobApplication{
		JobID:       jobID,
		TemporaryID: userID.(uuid.UUID),
		Message:     req.Message,
		Status:      "pending",
		AppliedAt:   time.Now(),
	}

	if err := database.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to apply: " + err.Error(),
		})
		return
	}

	database.DB.Model(&job).Update("applicants", job.Applicants+1)

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Application submitted successfully",
		Data:    application,
	})
}

func (h *JobHandler) GetJobApplications(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	jobID := id.(uuid.UUID)

	pagination := utils.GetPagination(c)

	var applications []models.JobApplication
	var total int64

	query := database.DB.Model(&models.JobApplication{}).Where("job_id = ?", jobID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("Temporary").
		Preload("Agent").
		Order("applied_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&applications)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       applications,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}

type ReviewApplicationRequest struct {
	Status    string `json:"status" binding:"required,oneof=approved rejected"`
	AgentID   *uuid.UUID `json:"agent_id"`
	ReviewNote string `json:"review_note"`
}

func (h *JobHandler) ReviewApplication(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	appID := id.(uuid.UUID)

	var req ReviewApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	var application models.JobApplication
	if err := database.DB.First(&application, appID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Application not found",
		})
		return
	}

	if application.Status != "pending" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Application has already been reviewed",
		})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":      req.Status,
		"reviewed_at": now,
		"review_note": req.ReviewNote,
	}
	if req.AgentID != nil {
		updates["agent_id"] = *req.AgentID
	}

	database.DB.Model(&application).Updates(updates)

	if req.Status == "approved" {
		var job models.JobPosting
		database.DB.First(&job, application.JobID)
		database.DB.Model(&job).Update("hired_count", job.HiredCount+1)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Application reviewed successfully",
		Data:    application,
	})
}

func (h *JobHandler) GetMyApplications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	pagination := utils.GetPagination(c)

	var applications []models.JobApplication
	var total int64

	query := database.DB.Model(&models.JobApplication{}).Where("temporary_id = ?", userID.(uuid.UUID))

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("JobPosting").
		Preload("JobPosting.Employer").
		Preload("Agent").
		Order("applied_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&applications)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       applications,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}
