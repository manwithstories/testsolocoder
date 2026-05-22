package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/database"
	"e-learning-platform/internal/models"
	"e-learning-platform/internal/utils"
)

type QAHandler struct {
	cfg *config.Config
}

func NewQAHandler(cfg *config.Config) *QAHandler {
	return &QAHandler{cfg: cfg}
}

type CreateQuestionRequest struct {
	CourseID uuid.UUID  `json:"course_id" binding:"required"`
	LessonID *uuid.UUID `json:"lesson_id"`
	Title    string     `json:"title" binding:"required,max=500"`
	Content  string     `json:"content" binding:"required"`
}

type CreateAnswerRequest struct {
	QuestionID uuid.UUID `json:"question_id" binding:"required"`
	Content    string    `json:"content" binding:"required"`
}

func (h *QAHandler) CreateQuestion(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var course models.Course
	if err := database.DB.First(&course, req.CourseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	question := models.Question{
		CourseID: req.CourseID,
		UserID:   userID.(uuid.UUID),
		LessonID: req.LessonID,
		Title:    req.Title,
		Content:  req.Content,
	}

	if err := database.DB.Create(&question).Error; err != nil {
		utils.InternalError(c, "Failed to create question")
		return
	}

	utils.Created(c, question)
}

func (h *QAHandler) ListQuestions(c *gin.Context) {
	courseID := c.Query("course_id")
	userID := c.Query("user_id")
	resolved := c.Query("is_resolved")
	search := c.Query("search")

	query := database.DB.Model(&models.Question{})
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if resolved != "" {
		query = query.Where("is_resolved = ?", resolved == "true")
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var questions []models.Question
	var total int64

	query.Count(&total)

	page, pageSize := getPagination(c)
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	query.Preload("User").Preload("Answers.User").
		Order(sort+" "+order).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&questions)

	utils.Paginated(c, questions, total, page, pageSize)
}

func (h *QAHandler) GetQuestion(c *gin.Context) {
	id := c.Param("id")
	questionID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid question ID")
		return
	}

	var question models.Question
	if err := database.DB.Preload("User").
		Preload("Answers.User").
		Preload("Course").
		First(&question, questionID).Error; err != nil {
		utils.NotFound(c, "Question not found")
		return
	}

	database.DB.Model(&question).UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	utils.Success(c, question)
}

func (h *QAHandler) CreateAnswer(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var question models.Question
	if err := database.DB.First(&question, req.QuestionID).Error; err != nil {
		utils.NotFound(c, "Question not found")
		return
	}

	answer := models.Answer{
		QuestionID: req.QuestionID,
		UserID:     userID.(uuid.UUID),
		Content:    req.Content,
	}

	tx := database.DB.Begin()
	if err := tx.Create(&answer).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to create answer")
		return
	}

	if err := tx.Model(&question).UpdateColumn("reply_count", gorm.Expr("reply_count + 1")).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to update reply count")
		return
	}

	tx.Commit()
	utils.Created(c, answer)
}

func (h *QAHandler) MarkBestAnswer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	answerID := c.Param("answer_id")
	aID, err := uuid.Parse(answerID)
	if err != nil {
		utils.BadRequest(c, "Invalid answer ID")
		return
	}

	var answer models.Answer
	if err := database.DB.First(&answer, aID).Error; err != nil {
		utils.NotFound(c, "Answer not found")
		return
	}

	var question models.Question
	if err := database.DB.First(&question, answer.QuestionID).Error; err != nil {
		utils.NotFound(c, "Question not found")
		return
	}

	if question.UserID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Only the question owner can mark best answer")
		return
	}

	tx := database.DB.Begin()

	if err := tx.Model(&models.Answer{}).
		Where("question_id = ?", question.ID).
		Update("is_best", false).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to reset best answers")
		return
	}

	if err := tx.Model(&answer).Update("is_best", true).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to mark best answer")
		return
	}

	if err := tx.Model(&question).Update("is_resolved", true).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to resolve question")
		return
	}

	tx.Commit()
	utils.Success(c, gin.H{"message": "Best answer marked"})
}

func (h *QAHandler) LikeAnswer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	answerID := c.Param("answer_id")
	aID, err := uuid.Parse(answerID)
	if err != nil {
		utils.BadRequest(c, "Invalid answer ID")
		return
	}

	var existing models.AnswerLike
	if database.DB.Where("answer_id = ? AND user_id = ?", aID, userID).First(&existing).Error == nil {
		tx := database.DB.Begin()
		tx.Delete(&existing)
		tx.Model(&models.Answer{}).Where("id = ?", aID).UpdateColumn("like_count", gorm.Expr("like_count - 1"))
		tx.Commit()
		utils.Success(c, gin.H{"liked": false})
		return
	}

	tx := database.DB.Begin()
	like := models.AnswerLike{
		AnswerID: aID,
		UserID:   userID.(uuid.UUID),
	}
	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to like answer")
		return
	}
	tx.Model(&models.Answer{}).Where("id = ?", aID).UpdateColumn("like_count", gorm.Expr("like_count + 1"))
	tx.Commit()

	utils.Success(c, gin.H{"liked": true})
}

func (h *QAHandler) LikeQuestion(c *gin.Context) {
	userID, _ := c.Get("user_id")
	questionID := c.Param("id")
	qID, err := uuid.Parse(questionID)
	if err != nil {
		utils.BadRequest(c, "Invalid question ID")
		return
	}

	var existing models.QuestionLike
	if database.DB.Where("question_id = ? AND user_id = ?", qID, userID).First(&existing).Error == nil {
		tx := database.DB.Begin()
		tx.Delete(&existing)
		tx.Model(&models.Question{}).Where("id = ?", qID).UpdateColumn("like_count", gorm.Expr("like_count - 1"))
		tx.Commit()
		utils.Success(c, gin.H{"liked": false})
		return
	}

	tx := database.DB.Begin()
	like := models.QuestionLike{
		QuestionID: qID,
		UserID:     userID.(uuid.UUID),
	}
	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to like question")
		return
	}
	tx.Model(&models.Question{}).Where("id = ?", qID).UpdateColumn("like_count", gorm.Expr("like_count + 1"))
	tx.Commit()

	utils.Success(c, gin.H{"liked": true})
}

func (h *QAHandler) DeleteQuestion(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	id := c.Param("id")
	questionID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid question ID")
		return
	}

	var question models.Question
	if err := database.DB.First(&question, questionID).Error; err != nil {
		utils.NotFound(c, "Question not found")
		return
	}

	if role != "admin" && question.UserID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	database.DB.Delete(&question)
	utils.Success(c, gin.H{"message": "Question deleted"})
}

func (h *QAHandler) DeleteAnswer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	answerID := c.Param("answer_id")
	aID, err := uuid.Parse(answerID)
	if err != nil {
		utils.BadRequest(c, "Invalid answer ID")
		return
	}

	var answer models.Answer
	if err := database.DB.First(&answer, aID).Error; err != nil {
		utils.NotFound(c, "Answer not found")
		return
	}

	if role != "admin" && answer.UserID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	tx := database.DB.Begin()
	tx.Delete(&answer)
	tx.Model(&models.Question{}).Where("id = ?", answer.QuestionID).UpdateColumn("reply_count", gorm.Expr("reply_count - 1"))
	tx.Commit()

	utils.Success(c, gin.H{"message": "Answer deleted"})
}
