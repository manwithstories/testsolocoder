package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tutoring-platform/database"
	"tutoring-platform/models"
)

type StudentProfileRequest struct {
	GradeLevel    string `json:"gradeLevel"`
	School        string `json:"school"`
	LearningStyle string `json:"learningStyle"`
	PreferredTime string `json:"preferredTime"`
	Notes         string `json:"notes"`
	ParentName    string `json:"parentName"`
	ParentPhone   string `json:"parentPhone"`
	ParentEmail   string `json:"parentEmail"`
}

type LearningGoalRequest struct {
	SubjectID    uuid.UUID  `json:"subjectId" binding:"required"`
	Title        string     `json:"title" binding:"required"`
	Description  string     `json:"description"`
	TargetScore  float64    `json:"targetScore"`
	CurrentScore float64    `json:"currentScore"`
	Deadline     *time.Time `json:"deadline"`
}

type AssessmentAnswerRequest struct {
	QuestionID uuid.UUID `json:"questionId" binding:"required"`
	Answer     string    `json:"answer" binding:"required"`
	Score      float64   `json:"score"`
}

func UpdateStudentProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req StudentProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.StudentProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	updates := map[string]interface{}{
		"grade_level":    req.GradeLevel,
		"school":         req.School,
		"learning_style": req.LearningStyle,
		"preferred_time": req.PreferredTime,
		"notes":          req.Notes,
		"parent_name":    req.ParentName,
		"parent_phone":   req.ParentPhone,
		"parent_email":   req.ParentEmail,
	}

	if err := database.DB.Model(&profile).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetStudentProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var profile models.StudentProfile
	if err := database.DB.Preload("LearningGoals.Subject").Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func GetStudentByID(c *gin.Context) {
	id := c.Param("id")
	studentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var profile models.StudentProfile
	if err := database.DB.Preload("LearningGoals.Subject").Preload("User").Where("user_id = ?", studentID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func AddLearningGoal(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req LearningGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.StudentProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	goal := models.LearningGoal{
		StudentID:    profile.ID,
		SubjectID:    req.SubjectID,
		Title:        req.Title,
		Description:  req.Description,
		TargetScore:  req.TargetScore,
		CurrentScore: req.CurrentScore,
		Deadline:     req.Deadline,
	}

	if err := database.DB.Create(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add learning goal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Learning goal added successfully", "goal": goal})
}

func UpdateLearningGoal(c *gin.Context) {
	goalID := c.Param("goalId")
	userID, _ := c.Get("userId")

	var req LearningGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.StudentProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	var goal models.LearningGoal
	if err := database.DB.Where("id = ? AND student_id = ?", goalID, profile.ID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Learning goal not found"})
		return
	}

	updates := map[string]interface{}{
		"title":         req.Title,
		"description":   req.Description,
		"subject_id":    req.SubjectID,
		"target_score":  req.TargetScore,
		"current_score": req.CurrentScore,
		"deadline":      req.Deadline,
	}

	if err := database.DB.Model(&goal).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update learning goal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Learning goal updated successfully"})
}

func DeleteLearningGoal(c *gin.Context) {
	goalID := c.Param("goalId")
	userID, _ := c.Get("userId")

	var profile models.StudentProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	if err := database.DB.Where("id = ? AND student_id = ?", goalID, profile.ID).Delete(&models.LearningGoal{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete learning goal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Learning goal deleted successfully"})
}

func GetAssessmentQuestions(c *gin.Context) {
	subjectID := c.Query("subjectId")

	var questions []models.AssessmentQuestion
	query := database.DB.Preload("Subject").Where("is_active = ?", true)

	if subjectID != "" {
		query = query.Where("subject_id = ?", subjectID)
	}

	if err := query.Order("sort_order").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func SubmitAssessment(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req struct {
		Answers []AssessmentAnswerRequest `json:"answers" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.StudentProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	tx := database.DB.Begin()

	now := time.Now()
	for _, ans := range req.Answers {
		answer := models.AssessmentAnswer{
			StudentID:   profile.ID,
			QuestionID:  ans.QuestionID,
			Answer:      ans.Answer,
			Score:       ans.Score,
			CompletedAt: &now,
		}
		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save answer"})
			return
		}
	}

	tx.Model(&profile).Update("assessment_status", "completed")
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Assessment submitted successfully"})
}

func GetMyAssessment(c *gin.Context) {
	userID, _ := c.Get("userId")

	var profile models.StudentProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	var answers []models.AssessmentAnswer
	if err := database.DB.Preload("Question.Subject").Where("student_id = ?", profile.ID).Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assessment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
		"answers": answers,
	})
}

func MatchTeachers(c *gin.Context) {
	userID, _ := c.Get("userId")

	var profile models.StudentProfile
	if err := database.DB.Preload("Assessments").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	var subjectIDs []uuid.UUID
	for _, answer := range profile.Assessments {
		subjectIDs = append(subjectIDs, answer.Question.SubjectID)
	}

	var teachers []models.TeacherProfile
	query := database.DB.Preload("Subjects.Subject").Preload("User").
		Where("approval_status = ?", "approved")

	if len(subjectIDs) > 0 {
		query = query.Joins("JOIN teacher_subjects ON teacher_subjects.teacher_id = teacher_profiles.id").
			Where("teacher_subjects.subject_id IN ?", subjectIDs)
	}

	if err := query.Order("rating DESC").Limit(10).Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to match teachers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teachers": teachers, "matchCount": len(teachers)})
}
