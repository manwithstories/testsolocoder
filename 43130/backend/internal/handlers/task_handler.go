package handlers

import (
	"encoding/json"
	"time"
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"
	"wedding-planner/pkg/response"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct{}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

type TaskRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Assignee    string `json:"assignee"`
	DueDate     string `json:"due_date"`
	Priority    string `json:"priority"`
	ParentID    *uint  `json:"parent_id"`
	Order       int    `json:"order"`
}

type TaskTemplateRequest struct {
	Name      string `json:"name" binding:"required,max=200"`
	Category  string `json:"category"`
	TasksJSON string `json:"tasks_json"`
	IsDefault bool   `json:"is_default"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	task := models.Task{
		WeddingID:   weddingID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Assignee:    req.Assignee,
		Priority:    req.Priority,
		ParentID:    req.ParentID,
		Order:       req.Order,
		Status:      "pending",
	}

	if task.Priority == "" {
		task.Priority = "medium"
	}

	if req.DueDate != "" {
		dueDate, err := parseDate(req.DueDate)
		if err == nil {
			task.DueDate = &dueDate
		}
	}

	if err := db.Create(&task).Error; err != nil {
		response.InternalError(c, "Failed to create task")
		return
	}

	response.Created(c, task)
}

func (h *TaskHandler) GetList(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	db := database.GetDB()

	var tasks []models.Task

	status := c.Query("status")
	category := c.Query("category")

	query := db.Where("wedding_id = ?", weddingID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Order("`order` ASC, due_date ASC, created_at ASC").Find(&tasks)

	response.Success(c, tasks)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&task).Error; err != nil {
		response.NotFound(c, "Task not found")
		return
	}

	var subtasks []models.Task
	db.Where("wedding_id = ? AND parent_id = ?", weddingID, id).Order("`order` ASC").Find(&subtasks)

	response.Success(c, gin.H{
		"task":     task,
		"subtasks": subtasks,
	})
}

func (h *TaskHandler) Update(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&task).Error; err != nil {
		response.NotFound(c, "Task not found")
		return
	}

	updates := map[string]interface{}{
		"title":       req.Title,
		"description": req.Description,
		"category":    req.Category,
		"assignee":    req.Assignee,
		"priority":    req.Priority,
		"parent_id":   req.ParentID,
		"`order`":     req.Order,
	}

	if req.DueDate != "" {
		dueDate, err := parseDate(req.DueDate)
		if err == nil {
			updates["due_date"] = dueDate
		}
	} else {
		updates["due_date"] = nil
	}

	db.Model(&task).Updates(updates)

	response.Success(c, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&task).Error; err != nil {
		response.NotFound(c, "Task not found")
		return
	}

	db.Where("parent_id = ?", id).Delete(&models.Task{})

	db.Delete(&task)

	response.Success(c, gin.H{"message": "Task deleted successfully"})
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	type StatusRequest struct {
		Status string `json:"status" binding:"required,oneof=pending in_progress completed cancelled"`
	}

	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid status")
		return
	}

	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&task).Error; err != nil {
		response.NotFound(c, "Task not found")
		return
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}

	if req.Status == "completed" {
		now := time.Now()
		updates["completed_at"] = &now
	} else {
		updates["completed_at"] = nil
	}

	db.Model(&task).Updates(updates)

	response.Success(c, task)
}

func (h *TaskHandler) CreateTemplate(c *gin.Context) {
	var req TaskTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	if req.IsDefault {
		db.Model(&models.TaskTemplate{}).Where("is_default = ?", true).Update("is_default", false)
	}

	template := models.TaskTemplate{
		Name:      req.Name,
		Category:  req.Category,
		TasksJSON: req.TasksJSON,
		IsDefault: req.IsDefault,
	}

	if err := db.Create(&template).Error; err != nil {
		response.InternalError(c, "Failed to create template")
		return
	}

	response.Created(c, template)
}

func (h *TaskHandler) GetTemplates(c *gin.Context) {
	db := database.GetDB()

	var templates []models.TaskTemplate
	db.Order("is_default DESC, created_at DESC").Find(&templates)

	response.Success(c, templates)
}

func (h *TaskHandler) ApplyTemplate(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	templateID := c.GetUint("template_id")

	db := database.GetDB()

	var template models.TaskTemplate
	if err := db.First(&template, templateID).Error; err != nil {
		response.NotFound(c, "Template not found")
		return
	}

	var templateTasks []map[string]interface{}
	if err := json.Unmarshal([]byte(template.TasksJSON), &templateTasks); err != nil {
		response.BadRequest(c, "Invalid template format")
		return
	}

	var createdTasks []models.Task
	for _, t := range templateTasks {
		task := models.Task{
			WeddingID:  weddingID,
			Title:      getStringValue(t, "title"),
			Category:   getStringValue(t, "category"),
			Assignee:   getStringValue(t, "assignee"),
			Priority:   getStringValue(t, "priority"),
			Status:     "pending",
			TemplateID: &templateID,
		}

		if desc, ok := t["description"].(string); ok {
			task.Description = desc
		}
		if priority, ok := t["priority"].(string); ok && priority != "" {
			task.Priority = priority
		} else {
			task.Priority = "medium"
		}

		if order, ok := t["order"].(float64); ok {
			task.Order = int(order)
		}

		if dueDate, ok := t["due_date"].(string); ok && dueDate != "" {
			parsed, err := parseDate(dueDate)
			if err == nil {
				task.DueDate = &parsed
			}
		}

		createdTasks = append(createdTasks, task)
	}

	if len(createdTasks) > 0 {
		db.Create(&createdTasks)
	}

	response.Created(c, gin.H{
		"message": "Template applied successfully",
		"count":   len(createdTasks),
	})
}

func (h *TaskHandler) DeleteTemplate(c *gin.Context) {
	id := c.GetUint("id")

	db := database.GetDB()

	var template models.TaskTemplate
	if err := db.First(&template, id).Error; err != nil {
		response.NotFound(c, "Template not found")
		return
	}

	db.Delete(&template)

	response.Success(c, gin.H{"message": "Template deleted successfully"})
}

func (h *TaskHandler) GetTaskCategories(c *gin.Context) {
	categories := []string{
		"前期准备", "场地预定", "供应商选择", "婚礼筹备", "婚礼当天", "婚后事项",
	}

	response.Success(c, categories)
}

func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
