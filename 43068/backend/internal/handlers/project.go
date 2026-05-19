package handlers

import (
	"net/http"
	"strconv"
	"time"

	"freelancer-management/internal/database"
	"freelancer-management/internal/logger"
	"freelancer-management/internal/middleware"
	"freelancer-management/internal/models"
	"freelancer-management/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	db *gorm.DB
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{db: database.GetDB()}
}

type CreateProjectRequest struct {
	ClientID    uint           `json:"client_id" binding:"required"`
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	HourlyRate  float64        `json:"hourly_rate"`
	Deadline    *time.Time     `json:"deadline"`
	Budget      float64        `json:"budget"`
	Milestones  []MilestoneReq `json:"milestones"`
}

type MilestoneReq struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateProjectRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	HourlyRate  float64    `json:"hourly_rate"`
	Deadline    *time.Time `json:"deadline"`
	Budget      float64    `json:"budget"`
}

type UpdateMilestoneRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Completed   *bool      `json:"completed"`
}

var validTransitions = map[models.ProjectStatus][]models.ProjectStatus{
	models.ProjectStatusDraft:     {models.ProjectStatusActive},
	models.ProjectStatusActive:    {models.ProjectStatusCompleted, models.ProjectStatusDraft},
	models.ProjectStatusCompleted: {models.ProjectStatusArchived, models.ProjectStatusActive},
	models.ProjectStatusArchived:  {models.ProjectStatusCompleted},
}

func isValidStatusTransition(oldStatus, newStatus models.ProjectStatus) bool {
	validNext, exists := validTransitions[oldStatus]
	if !exists {
		return false
	}
	for _, s := range validNext {
		if s == newStatus {
			return true
		}
	}
	return false
}

func (h *ProjectHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var client models.Client
	if err := h.db.Where("id = ? AND user_id = ?", req.ClientID, userID).First(&client).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Client not found")
		return
	}

	tx := h.db.Begin()

	project := models.Project{
		UserID:      userID,
		ClientID:    req.ClientID,
		Name:        req.Name,
		Description: req.Description,
		Status:      models.ProjectStatusDraft,
		HourlyRate:  req.HourlyRate,
		Deadline:    req.Deadline,
		Budget:      req.Budget,
	}

	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		logger.LogError("Failed to create project: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create project")
		return
	}

	for _, m := range req.Milestones {
		milestone := models.Milestone{
			ProjectID:   project.ID,
			Title:       m.Title,
			Description: m.Description,
			DueDate:     m.DueDate,
			Completed:   false,
		}
		if err := tx.Create(&milestone).Error; err != nil {
			tx.Rollback()
			logger.LogError("Failed to create milestone: %v", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create milestone")
			return
		}
	}

	tx.Commit()

	logger.LogOperation(userID, "create_project", "Project created: "+project.Name)
	utils.SuccessResponse(c, project)
}

func (h *ProjectHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	status := c.Query("status")
	clientID := c.Query("client_id")
	offset := (page - 1) * perPage

	var projects []models.Project
	var total int64

	query := h.db.Model(&models.Project{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}

	query.Count(&total)
	query.Preload("Client").Preload("Milestones").Offset(offset).Limit(perPage).Find(&projects)

	utils.PaginatedSuccessResponse(c, projects, utils.PaginationMeta{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	})
}

func (h *ProjectHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var project models.Project
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Client").
		Preload("Milestones").
		Preload("TimeEntries").
		First(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}

	utils.SuccessResponse(c, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var project models.Project
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Status != "" {
		newStatus := models.ProjectStatus(req.Status)
		if !isValidStatusTransition(project.Status, newStatus) {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid status transition from "+string(project.Status)+" to "+string(newStatus))
			return
		}
		project.Status = newStatus
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	project.Description = req.Description
	if req.HourlyRate != 0 {
		project.HourlyRate = req.HourlyRate
	}
	if req.Deadline != nil {
		project.Deadline = req.Deadline
	}
	if req.Budget != 0 {
		project.Budget = req.Budget
	}

	if err := h.db.Save(&project).Error; err != nil {
		logger.LogError("Failed to update project: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update project")
		return
	}

	logger.LogOperation(userID, "update_project", "Project updated: "+project.Name)
	utils.SuccessResponse(c, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var project models.Project
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}

	tx := h.db.Begin()

	if err := tx.Where("project_id = ?", id).Delete(&models.Milestone{}).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete milestones")
		return
	}

	if err := tx.Delete(&project).Error; err != nil {
		tx.Rollback()
		logger.LogError("Failed to delete project: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	tx.Commit()

	logger.LogOperation(userID, "delete_project", "Project deleted: "+project.Name)
	utils.SuccessResponseWithMessage(c, "Project deleted successfully", nil)
}

func (h *ProjectHandler) AddMilestone(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var project models.Project
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}

	var req MilestoneReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	milestone := models.Milestone{
		ProjectID:   project.ID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Completed:   false,
	}

	if err := h.db.Create(&milestone).Error; err != nil {
		logger.LogError("Failed to create milestone: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create milestone")
		return
	}

	logger.LogOperation(userID, "add_milestone", "Milestone added: "+milestone.Title)
	utils.SuccessResponse(c, milestone)
}

func (h *ProjectHandler) UpdateMilestone(c *gin.Context) {
	userID := middleware.GetUserID(c)
	milestoneID, _ := strconv.Atoi(c.Param("milestone_id"))

	var milestone models.Milestone
	if err := h.db.First(&milestone, milestoneID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Milestone not found")
		return
	}

	var project models.Project
	if err := h.db.Where("id = ? AND user_id = ?", milestone.ProjectID, userID).First(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}

	var req UpdateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title != "" {
		milestone.Title = req.Title
	}
	milestone.Description = req.Description
	if req.DueDate != nil {
		milestone.DueDate = req.DueDate
	}
	if req.Completed != nil {
		milestone.Completed = *req.Completed
	}

	if err := h.db.Save(&milestone).Error; err != nil {
		logger.LogError("Failed to update milestone: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update milestone")
		return
	}

	logger.LogOperation(userID, "update_milestone", "Milestone updated: "+milestone.Title)
	utils.SuccessResponse(c, milestone)
}

func (h *ProjectHandler) DeleteMilestone(c *gin.Context) {
	userID := middleware.GetUserID(c)
	milestoneID, _ := strconv.Atoi(c.Param("milestone_id"))

	var milestone models.Milestone
	if err := h.db.First(&milestone, milestoneID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Milestone not found")
		return
	}

	var project models.Project
	if err := h.db.Where("id = ? AND user_id = ?", milestone.ProjectID, userID).First(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}

	if err := h.db.Delete(&milestone).Error; err != nil {
		logger.LogError("Failed to delete milestone: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete milestone")
		return
	}

	logger.LogOperation(userID, "delete_milestone", "Milestone deleted: "+milestone.Title)
	utils.SuccessResponseWithMessage(c, "Milestone deleted successfully", nil)
}
