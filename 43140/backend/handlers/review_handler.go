package handlers

import (
	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	DB *gorm.DB
}

func NewReviewHandler(db *gorm.DB) *ReviewHandler {
	return &ReviewHandler{DB: db}
}

type ReviewRequest struct {
	InterviewID uint              `json:"interview_id" binding:"required"`
	Rating      int               `json:"rating" binding:"required,min=1,max=5"`
	Feedback    string            `json:"feedback" binding:"required"`
	Strengths   string            `json:"strengths"`
	Weaknesses  string            `json:"weaknesses"`
	Status      models.ReviewStatus `json:"status" binding:"required,oneof=offer pass reject pending"`
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var interview models.Interview
	if err := h.DB.Preload("Application").Preload("Application.Job").
		First(&interview, req.InterviewID).Error; err != nil {
		utils.NotFound(c, "Interview not found")
		return
	}

	if interview.Application.Job.CompanyID != company.ID {
		utils.Forbidden(c, "You can only create reviews for your company's interviews")
		return
	}

	var existingReview models.Review
	if h.DB.Where("interview_id = ?", req.InterviewID).First(&existingReview).Error == nil {
		utils.BadRequest(c, "Review already exists for this interview")
		return
	}

	review := models.Review{
		InterviewID: req.InterviewID,
		Rating:      req.Rating,
		Feedback:    req.Feedback,
		Strengths:   req.Strengths,
		Weaknesses:  req.Weaknesses,
		Status:      req.Status,
	}

	tx := h.DB.Begin()

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to create review")
		return
	}

	tx.Model(&interview).Update("status", models.InterviewStatusCompleted)

	if req.Status == models.ReviewStatusOffer {
		tx.Model(&interview.Application).Update("status", models.ApplicationStatusAccepted)
	} else if req.Status == models.ReviewStatusReject {
		tx.Model(&interview.Application).Update("status", models.ApplicationStatusRejected)
	}

	tx.Commit()

	var jobSeekerUser models.User
	h.DB.Joins("JOIN job_seekers ON users.id = job_seekers.user_id").
		Where("job_seekers.id = ?", interview.Application.JobSeekerID).
		First(&jobSeekerUser)

	statusText := map[models.ReviewStatus]string{
		models.ReviewStatusOffer:  "offer extended",
		models.ReviewStatusPass:   "passed",
		models.ReviewStatusReject: "rejected",
		models.ReviewStatusPending: "pending",
	}

	CreateNotification(h.DB, jobSeekerUser.ID,
		"Interview feedback available",
		"Your interview for "+interview.Application.Job.Title+" feedback is available. Status: "+statusText[req.Status],
		models.NotificationTypeReview, &review.ID)

	utils.Success(c, review)
}

func (h *ReviewHandler) GetReview(c *gin.Context) {
	id := c.Param("id")

	var review models.Review
	if err := h.DB.Preload("Interview").Preload("Interview.Application").
		Preload("Interview.Application.Job").Preload("Interview.Application.Job.Company").
		First(&review, id).Error; err != nil {
		utils.NotFound(c, "Review not found")
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) GetReviewByInterview(c *gin.Context) {
	interviewID := c.Param("interview_id")

	var review models.Review
	if err := h.DB.Preload("Interview").Preload("Interview.Application").
		Preload("Interview.Application.Job").Preload("Interview.Application.Job.Company").
		Where("interview_id = ?", interviewID).First(&review).Error; err != nil {
		utils.NotFound(c, "Review not found")
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var review models.Review
	if err := h.DB.Preload("Interview").Preload("Interview.Application").
		Preload("Interview.Application.Job").
		First(&review, id).Error; err != nil {
		utils.NotFound(c, "Review not found")
		return
	}

	if review.Interview.Application.Job.CompanyID != company.ID {
		utils.Forbidden(c, "You can only update reviews for your company's interviews")
		return
	}

	var req struct {
		Rating     *int               `json:"rating"`
		Feedback   *string            `json:"feedback"`
		Strengths  *string            `json:"strengths"`
		Weaknesses *string            `json:"weaknesses"`
		Status     *models.ReviewStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Feedback != nil {
		updates["feedback"] = *req.Feedback
	}
	if req.Strengths != nil {
		updates["strengths"] = *req.Strengths
	}
	if req.Weaknesses != nil {
		updates["weaknesses"] = *req.Weaknesses
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		if *req.Status == models.ReviewStatusOffer {
			h.DB.Model(&models.Application{}).Where("id = ?", review.Interview.ApplicationID).Update("status", models.ApplicationStatusAccepted)
		} else if *req.Status == models.ReviewStatusReject {
			h.DB.Model(&models.Application{}).Where("id = ?", review.Interview.ApplicationID).Update("status", models.ApplicationStatusRejected)
		}
	}

	h.DB.Model(&review).Updates(updates)

	utils.Success(c, nil)
}

func (h *ReviewHandler) ListCompanyReviews(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var reviews []models.Review
	query := h.DB.Preload("Interview").Preload("Interview.Application").
		Preload("Interview.Application.Job").
		Preload("Interview.Application.JobSeeker").Preload("Interview.Application.JobSeeker.User").
		Joins("JOIN interviews ON reviews.interview_id = interviews.id").
		Joins("JOIN applications ON interviews.application_id = applications.id").
		Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("jobs.company_id = ?", company.ID)

	if status := c.Query("status"); status != "" {
		query = query.Where("reviews.status = ?", status)
	}

	query.Order("reviews.created_at desc").Find(&reviews)

	utils.Success(c, reviews)
}

func (h *ReviewHandler) ListJobSeekerReviews(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var reviews []models.Review
	h.DB.Preload("Interview").Preload("Interview.Application").
		Preload("Interview.Application.Job").Preload("Interview.Application.Job.Company").
		Joins("JOIN interviews ON reviews.interview_id = interviews.id").
		Joins("JOIN applications ON interviews.application_id = applications.id").
		Where("applications.jobseeker_id = ?", jobSeeker.ID).
		Order("reviews.created_at desc").Find(&reviews)

	utils.Success(c, reviews)
}
