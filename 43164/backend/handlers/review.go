package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tutoring-platform/database"
	"tutoring-platform/models"
)

type ReviewRequest struct {
	BookingID   uuid.UUID `json:"bookingId" binding:"required"`
	Rating      int       `json:"rating" binding:"required,min=1,max=5"`
	Content     string    `json:"content" binding:"required"`
	Tags        string    `json:"tags"`
	IsAnonymous bool      `json:"isAnonymous"`
}

type TeacherReplyRequest struct {
	Reply string `json:"reply" binding:"required"`
}

func CreateReview(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var booking models.Booking
	if err := database.DB.Where("id = ?", req.BookingID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.Status != models.BookingStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only review completed bookings"})
		return
	}

	var existingReview models.Review
	if database.DB.Where("booking_id = ? AND reviewer_id = ?", req.BookingID, userID).First(&existingReview).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already reviewed this booking"})
		return
	}

	var revieweeID uuid.UUID
	if booking.StudentID == userID {
		revieweeID = booking.TeacherID
	} else {
		revieweeID = booking.StudentID
	}

	review := models.Review{
		BookingID:   req.BookingID,
		TeacherID:   booking.TeacherID,
		StudentID:   booking.StudentID,
		ReviewerID:  userID.(uuid.UUID),
		RevieweeID:  revieweeID,
		Rating:      req.Rating,
		Content:     req.Content,
		Tags:        req.Tags,
		IsAnonymous: req.IsAnonymous,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	var teacherProfile models.TeacherProfile
	if tx.Where("user_id = ?", booking.TeacherID).First(&teacherProfile).Error == nil {
		var reviews []models.Review
		tx.Where("reviewee_id = ?", booking.TeacherID).Find(&reviews)

		totalRating := 0
		for _, r := range reviews {
			totalRating += r.Rating
		}
		avgRating := float64(totalRating) / float64(len(reviews))

		tx.Model(&teacherProfile).Updates(map[string]interface{}{
			"rating":       avgRating,
			"review_count": len(reviews),
		})
	}

	tx.Commit()

	createNotification(c, revieweeID, models.NotificationTypeNewReview, "New Review", "You have received a new review")

	c.JSON(http.StatusCreated, gin.H{"message": "Review created", "review": review})
}

func GetReviews(c *gin.Context) {
	teacherID := c.Query("teacherId")
	studentID := c.Query("studentId")
	minRating := c.Query("minRating")

	var reviews []models.Review
	query := database.DB.Preload("Reviewer").Preload("Reviewee").Preload("Booking")

	if teacherID != "" {
		query = query.Where("reviewee_id = ?", teacherID)
	}
	if studentID != "" {
		query = query.Where("reviewer_id = ?", studentID)
	}
	if minRating != "" {
		query = query.Where("rating >= ?", minRating)
	}

	if err := query.Order("created_at DESC").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func GetReviewByID(c *gin.Context) {
	id := c.Param("id")

	var review models.Review
	if err := database.DB.Preload("Reviewer").Preload("Reviewee").Preload("Booking").Where("id = ? AND is_hidden = ?", id, false).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func ReplyToReview(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userId")

	var review models.Review
	if err := database.DB.Where("id = ?", id).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if review.TeacherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the teacher can reply to this review"})
		return
	}

	if review.TeacherReply != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already replied to this review"})
		return
	}

	var req TeacherReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"teacher_reply":     req.Reply,
		"teacher_replied_at": &now,
	}

	if err := database.DB.Model(&review).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reply to review"})
		return
	}

	createNotification(c, review.ReviewerID, models.NotificationTypeSystem, "Review Reply", "The teacher has replied to your review")

	c.JSON(http.StatusOK, gin.H{"message": "Reply posted successfully"})
}

func GetTeacherReviews(c *gin.Context) {
	id := c.Param("teacherId")

	var reviews []models.Review
	if err := database.DB.Preload("Reviewer").Preload("Booking").Where("reviewee_id = ? AND is_hidden = ?", id, false).Order("created_at DESC").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	var summary struct {
		TotalReviews int     `json:"totalReviews"`
		AverageRating float64 `json:"averageRating"`
		FiveStar     int     `json:"fiveStar"`
		FourStar     int     `json:"fourStar"`
		ThreeStar    int     `json:"threeStar"`
		TwoStar      int     `json:"twoStar"`
		OneStar      int     `json:"oneStar"`
		Reviews      []models.Review `json:"reviews"`
	}

	summary.Reviews = reviews
	summary.TotalReviews = len(reviews)

	for _, r := range reviews {
		switch r.Rating {
		case 5:
			summary.FiveStar++
		case 4:
			summary.FourStar++
		case 3:
			summary.ThreeStar++
		case 2:
			summary.TwoStar++
		case 1:
			summary.OneStar++
		}
	}

	if summary.TotalReviews > 0 {
		totalRating := 0
		for _, r := range reviews {
			totalRating += r.Rating
		}
		summary.AverageRating = float64(totalRating) / float64(summary.TotalReviews)
	}

	c.JSON(http.StatusOK, summary)
}

func HideReview(c *gin.Context) {
	id := c.Param("id")

	var review models.Review
	if err := database.DB.Where("id = ?", id).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	review.IsHidden = true
	if err := database.DB.Save(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hide review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review hidden successfully"})
}
