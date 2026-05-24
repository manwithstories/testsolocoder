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

type JobTemplateHandler struct{}

func NewJobTemplateHandler() *JobTemplateHandler {
	return &JobTemplateHandler{}
}

type CreateTemplateRequest struct {
	Name          string  `json:"name" binding:"required,max=100"`
	ActivityType  string  `json:"activity_type"`
	Position      string  `json:"position" binding:"required"`
	Description   string  `json:"description"`
	SalaryPerHour float64 `json:"salary_per_hour"`
	WorkHours     string  `json:"work_hours"`
	Requirements  string  `json:"requirements"`
	Benefits      string  `json:"benefits"`
}

type UpdateTemplateRequest struct {
	Name          string  `json:"name"`
	ActivityType  string  `json:"activity_type"`
	Position      string  `json:"position"`
	Description   string  `json:"description"`
	SalaryPerHour float64 `json:"salary_per_hour"`
	WorkHours     string  `json:"work_hours"`
	Requirements  string  `json:"requirements"`
	Benefits      string  `json:"benefits"`
}

type BatchImportTemplatesRequest struct {
	Templates []CreateTemplateRequest `json:"templates" binding:"required"`
}

type ApplyTemplateRequest struct {
	TemplateID  uuid.UUID `json:"template_id" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     string    `json:"end_date" binding:"required"`
	Headcount   int       `json:"headcount" binding:"required,gt=0"`
}

func (h *JobTemplateHandler) CreateTemplate(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	template := models.JobTemplate{
		EmployerID:    userID.(uuid.UUID),
		Name:          req.Name,
		ActivityType:  req.ActivityType,
		Position:      req.Position,
		Description:   req.Description,
		SalaryPerHour: req.SalaryPerHour,
		WorkHours:     req.WorkHours,
		Requirements:  req.Requirements,
		Benefits:      req.Benefits,
	}

	if err := database.DB.Create(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to create template: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Template created successfully",
		Data:    template,
	})
}

func (h *JobTemplateHandler) GetTemplates(c *gin.Context) {
	userID, _ := c.Get("user_id")
	pagination := utils.GetPagination(c)

	var templates []models.JobTemplate
	var total int64

	query := database.DB.Model(&models.JobTemplate{}).Where("employer_id = ?", userID.(uuid.UUID))

	if pagination.Keyword != "" {
		query = query.Where("name ILIKE ? OR position ILIKE ?",
			"%"+pagination.Keyword+"%",
			"%"+pagination.Keyword+"%",
		)
	}

	if activityType := c.Query("activity_type"); activityType != "" {
		query = query.Where("activity_type = ?", activityType)
	}

	query.Count(&total)
	query.Order("created_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&templates)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       templates,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}

func (h *JobTemplateHandler) GetTemplate(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	templateID := id.(uuid.UUID)
	userID, _ := c.Get("user_id")

	var template models.JobTemplate
	if err := database.DB.Where("id = ? AND employer_id = ?", templateID, userID.(uuid.UUID)).First(&template).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Template not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data:    template,
	})
}

func (h *JobTemplateHandler) UpdateTemplate(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	templateID := id.(uuid.UUID)
	userID, _ := c.Get("user_id")

	var template models.JobTemplate
	if err := database.DB.Where("id = ? AND employer_id = ?", templateID, userID.(uuid.UUID)).First(&template).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Template not found",
		})
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ActivityType != "" {
		updates["activity_type"] = req.ActivityType
	}
	if req.Position != "" {
		updates["position"] = req.Position
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.SalaryPerHour > 0 {
		updates["salary_per_hour"] = req.SalaryPerHour
	}
	if req.WorkHours != "" {
		updates["work_hours"] = req.WorkHours
	}
	if req.Requirements != "" {
		updates["requirements"] = req.Requirements
	}
	if req.Benefits != "" {
		updates["benefits"] = req.Benefits
	}

	if len(updates) > 0 {
		database.DB.Model(&template).Updates(updates)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Template updated successfully",
		Data:    template,
	})
}

func (h *JobTemplateHandler) DeleteTemplate(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	templateID := id.(uuid.UUID)
	userID, _ := c.Get("user_id")

	var template models.JobTemplate
	if err := database.DB.Where("id = ? AND employer_id = ?", templateID, userID.(uuid.UUID)).First(&template).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Template not found",
		})
		return
	}

	database.DB.Delete(&template)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Template deleted successfully",
	})
}

func (h *JobTemplateHandler) BatchImportTemplates(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req BatchImportTemplatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	if len(req.Templates) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "At least one template is required",
		})
		return
	}

	if len(req.Templates) > 100 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Maximum 100 templates can be imported at once",
		})
		return
	}

	var createdTemplates []models.JobTemplate
	for _, t := range req.Templates {
		template := models.JobTemplate{
			EmployerID:    userID.(uuid.UUID),
			Name:          t.Name,
			ActivityType:  t.ActivityType,
			Position:      t.Position,
			Description:   t.Description,
			SalaryPerHour: t.SalaryPerHour,
			WorkHours:     t.WorkHours,
			Requirements:  t.Requirements,
			Benefits:      t.Benefits,
		}
		if err := database.DB.Create(&template).Error; err == nil {
			createdTemplates = append(createdTemplates, template)
		}
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Successfully imported " + string(rune(len(createdTemplates))) + " templates",
		Data: gin.H{
			"imported_count": len(createdTemplates),
			"templates":      createdTemplates,
		},
	})
}

func (h *JobTemplateHandler) ApplyTemplate(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ApplyTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	var template models.JobTemplate
	if err := database.DB.Where("id = ? AND employer_id = ?", req.TemplateID, userID.(uuid.UUID)).First(&template).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Template not found",
		})
		return
	}

	job := models.JobPosting{
		EmployerID:    userID.(uuid.UUID),
		Title:         template.Name,
		Description:   template.Description,
		ActivityType:  template.ActivityType,
		Position:      template.Position,
		Location:      req.Location,
		StartDate:     parseDate(req.StartDate),
		EndDate:       parseDate(req.EndDate),
		SalaryPerHour: template.SalaryPerHour,
		WorkHours:     template.WorkHours,
		Headcount:     req.Headcount,
		Requirements:  template.Requirements,
		Benefits:      template.Benefits,
		Status:        "recruiting",
	}

	if err := database.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to create job from template: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Job created from template successfully",
		Data:    job,
	})
}

func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
