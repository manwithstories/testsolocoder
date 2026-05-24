package handlers

import (
	"time"

	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatisticsHandler struct {
	DB *gorm.DB
}

func NewStatisticsHandler(db *gorm.DB) *StatisticsHandler {
	return &StatisticsHandler{DB: db}
}

func (h *StatisticsHandler) GetCompanyStatistics(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var stats struct {
		TotalJobs          int64   `json:"total_jobs"`
		OpenJobs           int64   `json:"open_jobs"`
		TotalApplications  int64   `json:"total_applications"`
		TotalInterviews    int64   `json:"total_interviews"`
		TotalOffers        int64   `json:"total_offers"`
		InterviewRate      float64 `json:"interview_rate"`
		OfferRate          float64 `json:"offer_rate"`
		TotalViews         int64   `json:"total_views"`
		ConversionRate     float64 `json:"conversion_rate"`
	}

	h.DB.Model(&models.Job{}).Where("company_id = ?", company.ID).Count(&stats.TotalJobs)
	h.DB.Model(&models.Job{}).Where("company_id = ? AND status = ?", company.ID, models.JobStatusOpen).Count(&stats.OpenJobs)

	h.DB.Table("applications").Joins("JOIN jobs ON applications.job_id = jobs.id").Where("jobs.company_id = ?", company.ID).Count(&stats.TotalApplications)
	h.DB.Table("interviews").Joins("JOIN applications ON interviews.application_id = applications.id").Joins("JOIN jobs ON applications.job_id = jobs.id").Where("jobs.company_id = ?", company.ID).Count(&stats.TotalInterviews)
	h.DB.Table("reviews").Joins("JOIN interviews ON reviews.interview_id = interviews.id").Joins("JOIN applications ON interviews.application_id = applications.id").Joins("JOIN jobs ON applications.job_id = jobs.id").Where("jobs.company_id = ? AND reviews.status = ?", company.ID, models.ReviewStatusOffer).Count(&stats.TotalOffers)

	if stats.TotalApplications > 0 {
		stats.InterviewRate = float64(stats.TotalInterviews) / float64(stats.TotalApplications) * 100
		stats.OfferRate = float64(stats.TotalOffers) / float64(stats.TotalApplications) * 100
	}

	var totalViews int64
	h.DB.Model(&models.Job{}).Where("company_id = ?", company.ID).Select("COALESCE(SUM(views), 0)").Scan(&totalViews)
	stats.TotalViews = totalViews

	if stats.TotalViews > 0 {
		stats.ConversionRate = float64(stats.TotalApplications) / float64(stats.TotalViews) * 100
	}

	utils.Success(c, stats)
}

func (h *StatisticsHandler) GetJobStatistics(c *gin.Context) {
	userID := middleware.GetUserID(c)
	jobID := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var job models.Job
	if err := h.DB.Where("id = ? AND company_id = ?", jobID, company.ID).First(&job).Error; err != nil {
		utils.NotFound(c, "Job not found")
		return
	}

	var stats struct {
		JobID             uint    `json:"job_id"`
		JobTitle          string  `json:"job_title"`
		Views             int     `json:"views"`
		Applications      int64   `json:"applications"`
		Interviews        int64   `json:"interviews"`
		Offers            int64   `json:"offers"`
		ConversionRate    float64 `json:"conversion_rate"`
		InterviewRate     float64 `json:"interview_rate"`
	}

	stats.JobID = job.ID
	stats.JobTitle = job.Title
	stats.Views = job.Views

	h.DB.Model(&models.Application{}).Where("job_id = ?", job.ID).Count(&stats.Applications)
	h.DB.Table("interviews").Joins("JOIN applications ON interviews.application_id = applications.id").Where("applications.job_id = ?", job.ID).Count(&stats.Interviews)
	h.DB.Table("reviews").Joins("JOIN interviews ON reviews.interview_id = interviews.id").Joins("JOIN applications ON interviews.application_id = applications.id").Where("applications.job_id = ? AND reviews.status = ?", job.ID, models.ReviewStatusOffer).Count(&stats.Offers)

	if stats.Views > 0 {
		stats.ConversionRate = float64(stats.Applications) / float64(stats.Views) * 100
	}
	if stats.Applications > 0 {
		stats.InterviewRate = float64(stats.Interviews) / float64(stats.Applications) * 100
	}

	utils.Success(c, stats)
}

func (h *StatisticsHandler) GetApplicationTrend(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	days := 30
	if d := c.Query("days"); d != "" {
		days = parseInt(d, 30)
	}

	type TrendData struct {
		Date         string `json:"date"`
		Applications int64  `json:"applications"`
		Interviews   int64  `json:"interviews"`
	}

	var data []TrendData
	startDate := time.Now().AddDate(0, 0, -days)

	h.DB.Table("applications").
		Select("DATE(applications.created_at) as date, COUNT(*) as applications").
		Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("jobs.company_id = ? AND applications.created_at >= ?", company.ID, startDate).
		Group("DATE(applications.created_at)").
		Scan(&data)

	utils.Success(c, data)
}

func (h *StatisticsHandler) GetJobSeekerStatistics(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var stats struct {
		TotalApplications int64   `json:"total_applications"`
		InterviewCount    int64   `json:"interview_count"`
		OfferCount        int64   `json:"offer_count"`
		RejectionCount    int64   `json:"rejection_count"`
		PendingCount      int64   `json:"pending_count"`
		SuccessRate       float64 `json:"success_rate"`
		ResumeCount       int64   `json:"resume_count"`
	}

	h.DB.Model(&models.Application{}).Where("jobseeker_id = ?", jobSeeker.ID).Count(&stats.TotalApplications)
	h.DB.Table("interviews").Joins("JOIN applications ON interviews.application_id = applications.id").Where("applications.jobseeker_id = ?", jobSeeker.ID).Count(&stats.InterviewCount)
	h.DB.Table("reviews").Joins("JOIN interviews ON reviews.interview_id = interviews.id").Joins("JOIN applications ON interviews.application_id = applications.id").Where("applications.jobseeker_id = ? AND reviews.status = ?", jobSeeker.ID, models.ReviewStatusOffer).Count(&stats.OfferCount)
	h.DB.Model(&models.Application{}).Where("jobseeker_id = ? AND status = ?", jobSeeker.ID, models.ApplicationStatusRejected).Count(&stats.RejectionCount)
	h.DB.Model(&models.Application{}).Where("jobseeker_id = ? AND status = ?", jobSeeker.ID, models.ApplicationStatusPending).Count(&stats.PendingCount)
	h.DB.Model(&models.Resume{}).Where("jobseeker_id = ?", jobSeeker.ID).Count(&stats.ResumeCount)

	if stats.TotalApplications > 0 {
		stats.SuccessRate = float64(stats.OfferCount) / float64(stats.TotalApplications) * 100
	}

	utils.Success(c, stats)
}

func (h *StatisticsHandler) GetNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var notifications []models.Notification
	query := h.DB.Where("user_id = ?", userID)

	if read := c.Query("read"); read != "" {
		if read == "true" {
			query = query.Where("read = ?", true)
		} else if read == "false" {
			query = query.Where("read = ?", false)
		}
	}

	query.Order("created_at desc").Limit(50).Find(&notifications)

	utils.Success(c, notifications)
}

func (h *StatisticsHandler) MarkNotificationRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := h.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read", true)

	if result.RowsAffected == 0 {
		utils.NotFound(c, "Notification not found")
		return
	}

	utils.Success(c, nil)
}

func (h *StatisticsHandler) MarkAllNotificationsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)

	h.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true)

	utils.Success(c, nil)
}
