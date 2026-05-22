package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/database"
	"e-learning-platform/internal/models"
	"e-learning-platform/internal/utils"
)

type ProgressHandler struct {
	cfg *config.Config
}

func NewProgressHandler(cfg *config.Config) *ProgressHandler {
	return &ProgressHandler{cfg: cfg}
}

type UpdateProgressRequest struct {
	LastPosition  float64 `json:"last_position"`
	TotalDuration float64 `json:"total_duration"`
	IsCompleted   bool    `json:"is_completed"`
}

type CreateNoteRequest struct {
	LessonID  uuid.UUID `json:"lesson_id" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	Timestamp float64   `json:"timestamp"`
}

func (h *ProgressHandler) UpdateProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	lessonID := c.Param("lesson_id")
	lID, err := uuid.Parse(lessonID)
	if err != nil {
		utils.BadRequest(c, "Invalid lesson ID")
		return
	}

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var lesson models.Lesson
	if err := database.DB.First(&lesson, lID).Error; err != nil {
		utils.NotFound(c, "Lesson not found")
		return
	}

	var chapter models.Chapter
	database.DB.First(&chapter, lesson.ChapterID)

	var order models.Order
	if err := database.DB.Where("user_id = ? AND course_id = ? AND status = ?",
		userID, chapter.CourseID, models.OrderPaid).First(&order).Error; err != nil {
		if !lesson.IsFree && !chapter.IsFree {
			var course models.Course
			database.DB.First(&course, chapter.CourseID)
			if !course.IsFree {
				utils.Forbidden(c, "Please purchase the course first")
				return
			}
		}
	}

	var progress models.Progress
	result := database.DB.Where("user_id = ? AND lesson_id = ?", userID, lID).First(&progress)

	if result.Error != nil {
		progress = models.Progress{
			UserID:        userID.(uuid.UUID),
			CourseID:      chapter.CourseID,
			LessonID:      lID,
			LastPosition:  req.LastPosition,
			TotalDuration: req.TotalDuration,
			IsCompleted:   req.IsCompleted,
		}
		if err := database.DB.Create(&progress).Error; err != nil {
			utils.InternalError(c, "Failed to create progress")
			return
		}
	} else {
		updates := map[string]interface{}{
			"last_position":  req.LastPosition,
			"total_duration": req.TotalDuration,
		}
		if req.IsCompleted {
			updates["is_completed"] = true
		}
		database.DB.Model(&progress).Updates(updates)
	}

	utils.Success(c, gin.H{"message": "Progress updated"})
}

func (h *ProgressHandler) GetCourseProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	courseID := c.Param("course_id")
	cID, err := uuid.Parse(courseID)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	var progresses []models.Progress
	database.DB.Where("user_id = ? AND course_id = ?", userID, cID).Find(&progresses)

	var lessons []models.Lesson
	database.DB.Joins("JOIN chapters ON chapters.id = lessons.chapter_id").
		Where("chapters.course_id = ?", cID).
		Find(&lessons)

	totalLessons := len(lessons)
	completedLessons := 0
	for _, p := range progresses {
		if p.IsCompleted {
			completedLessons++
		}
	}

	completionRate := 0.0
	if totalLessons > 0 {
		completionRate = float64(completedLessons) / float64(totalLessons) * 100
	}

	utils.Success(c, gin.H{
		"progresses":      progresses,
		"total_lessons":   totalLessons,
		"completed_lessons": completedLessons,
		"completion_rate": completionRate,
	})
}

func (h *ProgressHandler) GetLessonProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	lessonID := c.Param("lesson_id")
	lID, err := uuid.Parse(lessonID)
	if err != nil {
		utils.BadRequest(c, "Invalid lesson ID")
		return
	}

	var progress models.Progress
	if err := database.DB.Where("user_id = ? AND lesson_id = ?", userID, lID).First(&progress).Error; err != nil {
		utils.Success(c, nil)
		return
	}

	utils.Success(c, progress)
}

func (h *ProgressHandler) CreateNote(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	note := models.Note{
		UserID:    userID.(uuid.UUID),
		LessonID:  req.LessonID,
		Content:   req.Content,
		Timestamp: req.Timestamp,
	}

	if err := database.DB.Create(&note).Error; err != nil {
		utils.InternalError(c, "Failed to create note")
		return
	}

	utils.Created(c, note)
}

func (h *ProgressHandler) ListNotes(c *gin.Context) {
	userID, _ := c.Get("user_id")
	lessonID := c.Query("lesson_id")

	query := database.DB.Where("user_id = ?", userID)
	if lessonID != "" {
		lID, err := uuid.Parse(lessonID)
		if err == nil {
			query = query.Where("lesson_id = ?", lID)
		}
	}

	var notes []models.Note
	page, pageSize := getPagination(c)
	var total int64

	query.Model(&models.Note{}).Count(&total)
	query.Preload("User").Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&notes)

	utils.Paginated(c, notes, total, page, pageSize)
}

func (h *ProgressHandler) UpdateNote(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	noteID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid note ID")
		return
	}

	var note models.Note
	if err := database.DB.First(&note, noteID).Error; err != nil {
		utils.NotFound(c, "Note not found")
		return
	}

	if note.UserID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var req struct {
		Content   string  `json:"content"`
		Timestamp float64 `json:"timestamp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	updates["timestamp"] = req.Timestamp

	database.DB.Model(&note).Updates(updates)
	utils.Success(c, gin.H{"message": "Note updated"})
}

func (h *ProgressHandler) DeleteNote(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	noteID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid note ID")
		return
	}

	var note models.Note
	if err := database.DB.First(&note, noteID).Error; err != nil {
		utils.NotFound(c, "Note not found")
		return
	}

	if note.UserID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	database.DB.Delete(&note)
	utils.Success(c, gin.H{"message": "Note deleted"})
}
