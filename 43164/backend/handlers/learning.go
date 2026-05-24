package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tutoring-platform/database"
	"tutoring-platform/models"
)

type LessonNoteRequest struct {
	BookingID uuid.UUID `json:"bookingId" binding:"required"`
	SubjectID uuid.UUID `json:"subjectId" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	Tags      string    `json:"tags"`
	IsPrivate bool      `json:"isPrivate"`
}

type HomeworkRequest struct {
	BookingID   uuid.UUID `json:"bookingId" binding:"required"`
	SubjectID   uuid.UUID `json:"subjectId" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	DueDate     string    `json:"dueDate" binding:"required"`
	MaxScore    float64   `json:"maxScore"`
	Attachments string    `json:"attachments"`
}

type HomeworkSubmissionRequest struct {
	HomeworkID  uuid.UUID `json:"homeworkId" binding:"required"`
	Content     string    `json:"content"`
	Attachments string    `json:"attachments"`
}

type GradeHomeworkRequest struct {
	Score    float64 `json:"score" binding:"required"`
	Feedback string  `json:"feedback"`
}

type FeedbackRequest struct {
	BookingID   uuid.UUID `json:"bookingId" binding:"required"`
	SubjectID   uuid.UUID `json:"subjectId" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	Type        string    `json:"type"`
	Progress    string    `json:"progress"`
	Suggestions string    `json:"suggestions"`
	NextSteps   string    `json:"nextSteps"`
}

type MilestoneRequest struct {
	Title       string               `json:"title" binding:"required"`
	Description string               `json:"description"`
	Type        models.MilestoneType `json:"type"`
	Icon        string               `json:"icon"`
	Color       string               `json:"color"`
	SubjectID   *uuid.UUID           `json:"subjectId"`
}

func CreateLessonNote(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")

	if userRole != models.RoleTeacher {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers can create lesson notes"})
		return
	}

	var req LessonNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var booking models.Booking
	if err := database.DB.Where("id = ? AND teacher_id = ?", req.BookingID, userID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	note := models.LessonNote{
		BookingID: req.BookingID,
		TeacherID: userID.(uuid.UUID),
		StudentID: booking.StudentID,
		SubjectID: req.SubjectID,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		IsPrivate: req.IsPrivate,
	}

	if err := database.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lesson note"})
		return
	}

	createNotification(c, booking.StudentID, models.NotificationTypeSystem, "New Lesson Note", "Your teacher has shared a new lesson note")

	c.JSON(http.StatusCreated, gin.H{"message": "Lesson note created", "note": note})
}

func GetLessonNotes(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")
	bookingID := c.Query("bookingId")

	var notes []models.LessonNote
	query := database.DB.Preload("Teacher").Preload("Student").Preload("Subject")

	if userRole == models.RoleTeacher {
		query = query.Where("teacher_id = ?", userID)
	} else if userRole == models.RoleStudent {
		query = query.Where("student_id = ? AND is_private = ?", userID, false)
	}

	if bookingID != "" {
		query = query.Where("booking_id = ?", bookingID)
	}

	if err := query.Order("created_at DESC").Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func GetLessonNoteByID(c *gin.Context) {
	id := c.Param("id")

	var note models.LessonNote
	if err := database.DB.Preload("Teacher").Preload("Student").Preload("Subject").Where("id = ?", id).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func UpdateLessonNote(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userId")

	var note models.LessonNote
	if err := database.DB.Where("id = ? AND teacher_id = ?", id, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var req LessonNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"title":     req.Title,
		"content":   req.Content,
		"tags":      req.Tags,
		"is_private": req.IsPrivate,
	}

	if err := database.DB.Model(&note).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

func DeleteLessonNote(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userId")

	if err := database.DB.Where("id = ? AND teacher_id = ?", id, userID).Delete(&models.LessonNote{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func CreateHomework(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")

	if userRole != models.RoleTeacher {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers can create homework"})
		return
	}

	var req HomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := userID.(uuid.UUID)

	var booking models.Booking
	if err := database.DB.Where("id = ? AND teacher_id = ?", req.BookingID, uid).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format"})
		return
	}

	homework := models.Homework{
		BookingID:   req.BookingID,
		TeacherID:   uid,
		StudentID:   booking.StudentID,
		SubjectID:   req.SubjectID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
		MaxScore:    req.MaxScore,
		Attachments: req.Attachments,
		Status:      models.HomeworkStatusPending,
	}

	if err := database.DB.Create(&homework).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create homework"})
		return
	}

	createNotification(c, booking.StudentID, models.NotificationTypeHomeworkAssigned, "New Homework", "You have a new homework assignment")

	c.JSON(http.StatusCreated, gin.H{"message": "Homework created", "homework": homework})
}

func GetHomeworks(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")

	var homeworks []models.Homework
	query := database.DB.Preload("Teacher").Preload("Student").Preload("Subject").Preload("Submission")

	if userRole == models.RoleTeacher {
		query = query.Where("teacher_id = ?", userID)
	} else if userRole == models.RoleStudent {
		query = query.Where("student_id = ?", userID)
	}

	if err := query.Order("due_date ASC").Find(&homeworks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch homeworks"})
		return
	}

	c.JSON(http.StatusOK, homeworks)
}

func SubmitHomework(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req HomeworkSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var homework models.Homework
	if err := database.DB.Where("id = ? AND student_id = ?", req.HomeworkID, userID).First(&homework).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Homework not found"})
		return
	}

	if homework.Status != models.HomeworkStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Homework already submitted"})
		return
	}

	submission := models.HomeworkSubmission{
		HomeworkID:  req.HomeworkID,
		StudentID:   userID.(uuid.UUID),
		Content:     req.Content,
		Attachments: req.Attachments,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&submission).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit homework"})
		return
	}

	homework.Status = models.HomeworkStatusSubmitted
	if err := tx.Save(&homework).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update homework"})
		return
	}

	tx.Commit()

	createNotification(c, homework.TeacherID, models.NotificationTypeSystem, "Homework Submitted", "A student has submitted homework")

	c.JSON(http.StatusCreated, gin.H{"message": "Homework submitted", "submission": submission})
}

func GradeHomework(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userId")

	var req GradeHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var submission models.HomeworkSubmission
	if err := database.DB.Preload("Homework").Where("id = ?", id).First(&submission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	uid := userID.(uuid.UUID)

	if submission.Homework.TeacherID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to grade this homework"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"score":      req.Score,
		"feedback":   req.Feedback,
		"graded_at":  &now,
		"graded_by":  &uid,
	}

	if err := database.DB.Model(&submission).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to grade homework"})
		return
	}

	submission.Homework.Status = models.HomeworkStatusGraded
	database.DB.Save(&submission.Homework)

	createNotification(c, submission.StudentID, models.NotificationTypeHomeworkGraded, "Homework Graded", "Your homework has been graded")

	c.JSON(http.StatusOK, gin.H{"message": "Homework graded successfully"})
}

func CreateFeedback(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")

	if userRole != models.RoleTeacher {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers can create feedback"})
		return
	}

	var req FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var booking models.Booking
	if err := database.DB.Where("id = ? AND teacher_id = ?", req.BookingID, userID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	feedback := models.Feedback{
		BookingID:   req.BookingID,
		TeacherID:   userID.(uuid.UUID),
		StudentID:   booking.StudentID,
		SubjectID:   req.SubjectID,
		Content:     req.Content,
		Type:        req.Type,
		Progress:    req.Progress,
		Suggestions: req.Suggestions,
		NextSteps:   req.NextSteps,
	}

	if err := database.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feedback"})
		return
	}

	createNotification(c, booking.StudentID, models.NotificationTypeSystem, "New Feedback", "Your teacher has shared feedback")

	c.JSON(http.StatusCreated, gin.H{"message": "Feedback created", "feedback": feedback})
}

func GetFeedbacks(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")

	var feedbacks []models.Feedback
	query := database.DB.Preload("Teacher").Preload("Student").Preload("Subject")

	if userRole == models.RoleTeacher {
		query = query.Where("teacher_id = ?", userID)
	} else if userRole == models.RoleStudent {
		query = query.Where("student_id = ?", userID)
	}

	if err := query.Order("created_at DESC").Find(&feedbacks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feedbacks"})
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

func GetMilestones(c *gin.Context) {
	userID, _ := c.Get("userId")

	var profile models.StudentProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	var milestones []models.Milestone
	if err := database.DB.Preload("Subject").Where("student_id = ?", profile.ID).Order("created_at DESC").Find(&milestones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch milestones"})
		return
	}

	c.JSON(http.StatusOK, milestones)
}

func CreateMilestone(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req MilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.StudentProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	milestone := models.Milestone{
		StudentID:   profile.ID,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Icon:        req.Icon,
		Color:       req.Color,
		SubjectID:   req.SubjectID,
	}

	if err := database.DB.Create(&milestone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create milestone"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Milestone created", "milestone": milestone})
}

func UpdateMilestone(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userId")

	var profile models.StudentProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	var milestone models.Milestone
	if err := database.DB.Where("id = ? AND student_id = ?", id, profile.ID).First(&milestone).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	var req struct {
		IsAchieved bool `json:"isAchieved"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"is_achieved": req.IsAchieved,
	}
	if req.IsAchieved {
		updates["achieved_at"] = &now
	}

	if err := database.DB.Model(&milestone).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update milestone"})
		return
	}

	if req.IsAchieved {
		createNotification(c, userID.(uuid.UUID), models.NotificationTypeMilestone, "Milestone Achieved!", "Congratulations! You have achieved a new milestone")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone updated successfully"})
}
