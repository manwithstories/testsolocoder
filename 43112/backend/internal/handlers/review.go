package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/database"
	"e-learning-platform/internal/models"
	"e-learning-platform/internal/utils"
)

type ReviewHandler struct {
	cfg *config.Config
}

func NewReviewHandler(cfg *config.Config) *ReviewHandler {
	return &ReviewHandler{cfg: cfg}
}

type CreateReviewRequest struct {
	CourseID    uuid.UUID `json:"course_id" binding:"required"`
	Rating      int       `json:"rating" binding:"required,min=1,max=5"`
	Content     string    `json:"content"`
	IsAnonymous bool      `json:"is_anonymous"`
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var course models.Course
	if err := database.DB.First(&course, req.CourseID).Error; err != nil {
		utils.NotFound(c, "Course not found")
		return
	}

	var order models.Order
	if err := database.DB.Where("user_id = ? AND course_id = ? AND status = ?",
		userID, req.CourseID, models.OrderPaid).First(&order).Error; err != nil {
		if !course.IsFree {
			utils.Forbidden(c, "You need to purchase the course first")
			return
		}
	}

	var existing models.Review
	if database.DB.Where("user_id = ? AND course_id = ?", userID, req.CourseID).First(&existing).Error == nil {
		utils.BadRequest(c, "You have already reviewed this course")
		return
	}

	review := models.Review{
		CourseID:    req.CourseID,
		UserID:      userID.(uuid.UUID),
		Rating:      req.Rating,
		Content:     req.Content,
		IsAnonymous: req.IsAnonymous,
	}

	tx := database.DB.Begin()
	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to create review")
		return
	}

	var avgRating float64
	var reviewCount int64
	tx.Model(&models.Review{}).Where("course_id = ?", req.CourseID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)
	tx.Model(&models.Review{}).Where("course_id = ?", req.CourseID).Count(&reviewCount)

	tx.Model(&models.Course{}).Where("id = ?", req.CourseID).
		Updates(map[string]interface{}{
			"avg_rating":   avgRating,
			"review_count": reviewCount,
		})

	tx.Commit()
	utils.Created(c, review)
}

func (h *ReviewHandler) ListReviews(c *gin.Context) {
	courseID := c.Query("course_id")
	userID := c.Query("user_id")
	minRating := c.Query("min_rating")

	query := database.DB.Model(&models.Review{})
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if minRating != "" {
		query = query.Where("rating >= ?", minRating)
	}

	var reviews []models.Review
	var total int64

	query.Count(&total)

	page, pageSize := getPagination(c)
	query.Preload("User").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reviews)

	utils.Paginated(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) GetMyReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	courseID := c.Query("course_id")

	cID, err := uuid.Parse(courseID)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	var review models.Review
	if err := database.DB.Where("user_id = ? AND course_id = ?", userID, cID).
		First(&review).Error; err != nil {
		utils.Success(c, nil)
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	id := c.Param("id")
	reviewID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid review ID")
		return
	}

	var review models.Review
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		utils.NotFound(c, "Review not found")
		return
	}

	if role != "admin" && review.UserID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not authorized")
		return
	}

	tx := database.DB.Begin()
	if err := tx.Delete(&review).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to delete review")
		return
	}

	var avgRating float64
	var reviewCount int64
	tx.Model(&models.Review{}).Where("course_id = ?", review.CourseID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)
	tx.Model(&models.Review{}).Where("course_id = ?", review.CourseID).Count(&reviewCount)

	tx.Model(&models.Course{}).Where("id = ?", review.CourseID).
		Updates(map[string]interface{}{
			"avg_rating":   avgRating,
			"review_count": reviewCount,
		})

	tx.Commit()
	utils.Success(c, gin.H{"message": "Review deleted"})
}

func (h *ReviewHandler) GetCourseSummary(c *gin.Context) {
	courseID := c.Param("course_id")
	cID, err := uuid.Parse(courseID)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	var avgRating float64
	var reviewCount int64
	var ratingDist []struct {
		Rating int   `json:"rating"`
		Count  int64 `json:"count"`
	}

	database.DB.Model(&models.Review{}).
		Where("course_id = ?", cID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)

	database.DB.Model(&models.Review{}).
		Where("course_id = ?", cID).
		Count(&reviewCount)

	database.DB.Model(&models.Review{}).
		Select("rating, COUNT(*) as count").
		Where("course_id = ?", cID).
		Group("rating").
		Scan(&ratingDist)

	utils.Success(c, gin.H{
		"avg_rating":    avgRating,
		"review_count":  reviewCount,
		"rating_distribution": ratingDist,
	})
}
