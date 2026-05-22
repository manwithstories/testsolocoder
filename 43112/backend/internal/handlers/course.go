package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/database"
	"e-learning-platform/internal/models"
	"e-learning-platform/internal/utils"
)

type CourseHandler struct {
	cfg *config.Config
}

func NewCourseHandler(cfg *config.Config) *CourseHandler {
	return &CourseHandler{cfg: cfg}
}

type CreateCourseRequest struct {
	Title         string  `json:"title" binding:"required,max=200"`
	Subtitle      string  `json:"subtitle" binding:"max=500"`
	Description   string  `json:"description" binding:"required"`
	Cover         string  `json:"cover"`
	Category      string  `json:"category" binding:"required"`
	Level         string  `json:"level" binding:"required,oneof=beginner intermediate advanced"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	IsFree        bool    `json:"is_free"`
	Tags          string  `json:"tags"`
}

type UpdateCourseRequest struct {
	Title         string  `json:"title" binding:"max=200"`
	Subtitle      string  `json:"subtitle" binding:"max=500"`
	Description   string  `json:"description"`
	Cover         string  `json:"cover"`
	Category      string  `json:"category"`
	Level         string  `json:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	IsFree        bool    `json:"is_free"`
	Tags          string  `json:"tags"`
}

type UpdateCourseStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft published offline rejected"`
}

type CreateChapterRequest struct {
	Title    string `json:"title" binding:"required,max=200"`
	Position int    `json:"position"`
	IsFree   bool   `json:"is_free"`
}

type CreateLessonRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Type        string `json:"type" binding:"required,oneof=video document quiz"`
	Content     string `json:"content"`
	VideoURL    string `json:"video_url"`
	VideoLength int    `json:"video_length"`
	DocURL      string `json:"doc_url"`
	DocName     string `json:"doc_name"`
	Position    int    `json:"position"`
	IsFree      bool   `json:"is_free"`
}

type CreateQuizRequest struct {
	Title     string                 `json:"title" binding:"required"`
	PassScore int                    `json:"pass_score"`
	TimeLimit int                    `json:"time_limit"`
	Questions []QuizQuestionRequest  `json:"questions" binding:"required"`
}

type QuizQuestionRequest struct {
	Content  string           `json:"content" binding:"required"`
	Type     string           `json:"type" binding:"required,oneof=single multiple"`
	Score    int              `json:"score"`
	Position int              `json:"position"`
	Options  []QuizOptionReq  `json:"options" binding:"required"`
}

type QuizOptionReq struct {
	Content   string `json:"content" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role == "instructor" {
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			utils.Forbidden(c, "Instructor not found")
			return
		}
		if user.InstructorStatus != models.InstructorApproved {
			utils.Forbidden(c, "Instructor qualification not approved")
			return
		}
	}

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	course := models.Course{
		InstructorID:  userID.(uuid.UUID),
		Title:         req.Title,
		Subtitle:      req.Subtitle,
		Description:   req.Description,
		Cover:         req.Cover,
		Category:      req.Category,
		Level:         models.CourseLevel(req.Level),
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		IsFree:        req.IsFree,
		Tags:          req.Tags,
		Status:        models.CourseDraft,
	}

	if req.IsFree {
		course.Price = 0
	}

	if err := database.DB.Create(&course).Error; err != nil {
		utils.InternalError(c, "Failed to create course")
		return
	}

	utils.Created(c, course)
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")
	courseID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	var course models.Course
	if err := database.DB.Preload("Instructor").
		Preload("Chapters.Lessons.Quiz.Questions.Options").
		Preload("Chapters.Lessons", "is_published = ?", true).
		First(&course, courseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	utils.Success(c, course)
}

func (h *CourseHandler) ListCourses(c *gin.Context) {
	category := c.Query("category")
	level := c.Query("level")
	status := c.Query("status")
	search := c.Query("search")
	isFree := c.Query("is_free")
	instructorID := c.Query("instructor_id")

	query := database.DB.Model(&models.Course{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", models.CoursePublished)
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if isFree != "" {
		query = query.Where("is_free = ?", isFree == "true")
	}
	if instructorID != "" {
		query = query.Where("instructor_id = ?", instructorID)
	}

	var courses []models.Course
	var total int64

	query.Count(&total)

	page, pageSize := getPagination(c)
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	query.Preload("Instructor").
		Order(sort+" "+order).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&courses)

	utils.Paginated(c, courses, total, page, pageSize)
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	courseID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized to edit this course")
		return
	}

	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Subtitle != "" {
		updates["subtitle"] = req.Subtitle
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Level != "" {
		updates["level"] = req.Level
	}
	updates["price"] = req.Price
	updates["original_price"] = req.OriginalPrice
	updates["is_free"] = req.IsFree
	updates["tags"] = req.Tags

	if err := database.DB.Model(&course).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update course")
		return
	}

	utils.Success(c, gin.H{"message": "Course updated successfully"})
}

func (h *CourseHandler) UpdateCourseStatus(c *gin.Context) {
	id := c.Param("id")
	courseID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var req UpdateCourseStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}
	if req.Status == "published" {
		now := time.Now()
		updates["published_at"] = &now
	}

	if err := database.DB.Model(&course).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update course status")
		return
	}

	utils.Success(c, gin.H{"message": "Course status updated"})
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	courseID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	if err := database.DB.Delete(&course).Error; err != nil {
		utils.InternalError(c, "Failed to delete course")
		return
	}

	utils.Success(c, gin.H{"message": "Course deleted successfully"})
}

func (h *CourseHandler) ListMyCourses(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page, pageSize := getPagination(c)
	var courses []models.Course
	var total int64

	database.DB.Model(&models.Course{}).Where("instructor_id = ?", userID).Count(&total)
	database.DB.Where("instructor_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&courses)

	utils.Paginated(c, courses, total, page, pageSize)
}

func (h *CourseHandler) CreateChapter(c *gin.Context) {
	courseID := c.Param("course_id")
	cID, err := uuid.Parse(courseID)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var course models.Course
	if err := database.DB.First(&course, cID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var req CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	chapter := models.Chapter{
		CourseID: cID,
		Title:    req.Title,
		Position: req.Position,
		IsFree:   req.IsFree,
	}

	if err := database.DB.Create(&chapter).Error; err != nil {
		utils.InternalError(c, "Failed to create chapter")
		return
	}

	utils.Created(c, chapter)
}

func (h *CourseHandler) UpdateChapter(c *gin.Context) {
	id := c.Param("id")
	chapterID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid chapter ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var chapter models.Chapter
	if err := database.DB.First(&chapter, chapterID).Error; err != nil {
		utils.NotFound(c, "Chapter not found")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, chapter.CourseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var req CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	database.DB.Model(&chapter).Updates(map[string]interface{}{
		"title":    req.Title,
		"position": req.Position,
		"is_free":  req.IsFree,
	})

	utils.Success(c, gin.H{"message": "Chapter updated"})
}

func (h *CourseHandler) DeleteChapter(c *gin.Context) {
	id := c.Param("id")
	chapterID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid chapter ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var chapter models.Chapter
	if err := database.DB.First(&chapter, chapterID).Error; err != nil {
		utils.NotFound(c, "Chapter not found")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, chapter.CourseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	database.DB.Delete(&chapter)
	utils.Success(c, gin.H{"message": "Chapter deleted"})
}

func (h *CourseHandler) CreateLesson(c *gin.Context) {
	chapterID := c.Param("chapter_id")
	cID, err := uuid.Parse(chapterID)
	if err != nil {
		utils.BadRequest(c, "Invalid chapter ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var chapter models.Chapter
	if err := database.DB.First(&chapter, cID).Error; err != nil {
		utils.NotFound(c, "Chapter not found")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, chapter.CourseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var req CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	lesson := models.Lesson{
		ChapterID:   cID,
		Title:       req.Title,
		Type:        models.LessonType(req.Type),
		Content:     req.Content,
		VideoURL:    req.VideoURL,
		VideoLength: req.VideoLength,
		DocURL:      req.DocURL,
		DocName:     req.DocName,
		Position:    req.Position,
		IsFree:      req.IsFree,
		IsPublished: true,
	}

	if err := database.DB.Create(&lesson).Error; err != nil {
		utils.InternalError(c, "Failed to create lesson")
		return
	}

	utils.Created(c, lesson)
}

func (h *CourseHandler) UpdateLesson(c *gin.Context) {
	id := c.Param("id")
	lessonID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid lesson ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		utils.NotFound(c, "Lesson not found")
		return
	}

	var chapter models.Chapter
	if err := database.DB.First(&chapter, lesson.ChapterID).Error; err != nil {
		utils.NotFound(c, "Chapter not found")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, chapter.CourseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var req CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	database.DB.Model(&lesson).Updates(map[string]interface{}{
		"title":       req.Title,
		"type":        req.Type,
		"content":     req.Content,
		"video_url":   req.VideoURL,
		"video_length": req.VideoLength,
		"doc_url":     req.DocURL,
		"doc_name":    req.DocName,
		"position":    req.Position,
		"is_free":     req.IsFree,
	})

	utils.Success(c, gin.H{"message": "Lesson updated"})
}

func (h *CourseHandler) DeleteLesson(c *gin.Context) {
	id := c.Param("id")
	lessonID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid lesson ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		utils.NotFound(c, "Lesson not found")
		return
	}

	var chapter models.Chapter
	if err := database.DB.First(&chapter, lesson.ChapterID).Error; err != nil {
		utils.NotFound(c, "Chapter not found")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, chapter.CourseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	database.DB.Delete(&lesson)
	utils.Success(c, gin.H{"message": "Lesson deleted"})
}

func (h *CourseHandler) CreateQuiz(c *gin.Context) {
	lessonID := c.Param("lesson_id")
	lID, err := uuid.Parse(lessonID)
	if err != nil {
		utils.BadRequest(c, "Invalid lesson ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var lesson models.Lesson
	if err := database.DB.First(&lesson, lID).Error; err != nil {
		utils.NotFound(c, "Lesson not found")
		return
	}

	var chapter models.Chapter
	if err := database.DB.First(&chapter, lesson.ChapterID).Error; err != nil {
		utils.NotFound(c, "Chapter not found")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, chapter.CourseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	if role != "admin" && course.InstructorID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var req CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tx := database.DB.Begin()

	quiz := models.Quiz{
		LessonID:  lID,
		Title:     req.Title,
		PassScore: req.PassScore,
		TimeLimit: req.TimeLimit,
	}

	if err := tx.Create(&quiz).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to create quiz")
		return
	}

	for _, q := range req.Questions {
		question := models.QuizQuestion{
			QuizID:   quiz.ID,
			Content:  q.Content,
			Type:     q.Type,
			Score:    q.Score,
			Position: q.Position,
		}
		if err := tx.Create(&question).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to create question")
			return
		}

		for _, o := range q.Options {
			option := models.QuizOption{
				QuestionID: question.ID,
				Content:    o.Content,
				IsCorrect:  o.IsCorrect,
			}
			if err := tx.Create(&option).Error; err != nil {
				tx.Rollback()
				utils.InternalError(c, "Failed to create option")
				return
			}
		}
	}

	tx.Commit()
	utils.Created(c, quiz)
}

func (h *CourseHandler) SubmitQuiz(c *gin.Context) {
	userID, _ := c.Get("user_id")
	quizIDStr := c.Param("quiz_id")
	quizID, err := uuid.Parse(quizIDStr)
	if err != nil {
		utils.BadRequest(c, "Invalid quiz ID")
		return
	}

	var req struct {
		Answers []struct {
			QuestionID uuid.UUID   `json:"question_id"`
			OptionIDs  []uuid.UUID `json:"option_ids"`
		} `json:"answers"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var quiz models.Quiz
	if err := database.DB.Preload("Questions.Options").First(&quiz, quizID).Error; err != nil {
		utils.NotFound(c, "Quiz not found")
		return
	}

	totalScore := 0
	userScore := 0

	for _, question := range quiz.Questions {
		totalScore += question.Score

		for _, ua := range req.Answers {
			if ua.QuestionID == question.ID {
				correct := true
				for _, opt := range question.Options {
					if opt.IsCorrect {
						found := false
						for _, uid := range ua.OptionIDs {
							if uid == opt.ID {
								found = true
								break
							}
						}
						if !found {
							correct = false
							break
						}
					} else {
						for _, uid := range ua.OptionIDs {
							if uid == opt.ID {
								correct = false
								break
							}
						}
					}
				}
				if correct {
					userScore += question.Score
				}
				break
			}
		}
	}

	now := time.Now()
	userQuiz := models.UserQuiz{
		UserID:      userID.(uuid.UUID),
		QuizID:      quizID,
		Score:       userScore,
		IsPassed:    userScore >= quiz.PassScore,
		CompletedAt: &now,
	}

	database.DB.Create(&userQuiz)

	utils.Success(c, gin.H{
		"score":       userScore,
		"total_score": totalScore,
		"is_passed":   userScore >= quiz.PassScore,
		"pass_score":  quiz.PassScore,
	})
}

func (h *CourseHandler) ListCategories(c *gin.Context) {
	var categories []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	database.DB.Model(&models.Course{}).
		Select("category, COUNT(*) as count").
		Where("status = ?", models.CoursePublished).
		Group("category").
		Find(&categories)

	utils.Success(c, categories)
}

func (h *CourseHandler) UpdateCourseHours(c *gin.Context) {
	id := c.Param("id")
	courseID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	var totalHours float64
	database.DB.Model(&models.Lesson{}).
		Joins("JOIN chapters ON chapters.id = lessons.chapter_id").
		Where("chapters.course_id = ? AND lessons.type = ?", courseID, models.LessonTypeVideo).
		Select("COALESCE(SUM(lessons.video_length) / 60.0, 0)").
		Scan(&totalHours)

	database.DB.Model(&course).Update("total_hours", totalHours)
	utils.Success(c, gin.H{"total_hours": totalHours})
}

func parseIntSafe(s string, def int) int {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return def
	}
	return n
}
