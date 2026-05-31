package services

import (
	"errors"

	"secondhand-platform/database"
	"secondhand-platform/models"

	"github.com/sirupsen/logrus"
)

type ReviewService struct{}

func NewReviewService() *ReviewService {
	return &ReviewService{}
}

func (s *ReviewService) CreateReview(reviewerID, revieweeID uint, orderID, repairOrderID *uint, reviewType string, rating int, content, images string, qualityScore, serviceScore *int) (*models.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, errors.New("评分必须在1-5之间")
	}

	review := &models.Review{
		ReviewerID:    reviewerID,
		RevieweeID:    revieweeID,
		OrderID:       orderID,
		RepairOrderID: repairOrderID,
		ReviewType:    reviewType,
		Rating:        rating,
		Content:       content,
		Images:        images,
		QualityScore:  qualityScore,
		ServiceScore:  serviceScore,
		Status:        1,
	}

	result := database.DB.Create(review)
	if result.Error != nil {
		return nil, result.Error
	}

	creditDelta := 0
	switch rating {
	case 5:
		creditDelta = 2
	case 4:
		creditDelta = 1
	case 3:
		creditDelta = 0
	case 2:
		creditDelta = -1
	case 1:
		creditDelta = -2
	}

	if creditDelta != 0 {
		userService := NewUserService()
		userService.UpdateCreditScore(revieweeID, creditDelta)
	}

	return review, nil
}

func (s *ReviewService) GetReviewByID(id uint) (*models.Review, error) {
	var review models.Review
	if err := database.DB.Preload("Reviewer").Preload("Reviewee").First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (s *ReviewService) ListReviews(page, pageSize int, revieweeID uint, reviewType string, minRating int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	db := database.DB.Model(&models.Review{}).Where("status = ?", 1)

	if revieweeID > 0 {
		db = db.Where("reviewee_id = ?", revieweeID)
	}
	if reviewType != "" {
		db = db.Where("review_type = ?", reviewType)
	}
	if minRating > 0 {
		db = db.Where("rating >= ?", minRating)
	}

	db.Count(&total)
	if err := db.Preload("Reviewer").Preload("Reviewee").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (s *ReviewService) GetAverageRating(userID uint, reviewType string) (float64, error) {
	var result struct {
		AvgRating float64
	}

	db := database.DB.Model(&models.Review{}).Where("reviewee_id = ? AND status = ?", userID, 1)
	if reviewType != "" {
		db = db.Where("review_type = ?", reviewType)
	}

	if err := db.Select("COALESCE(AVG(rating), 0) as avg_rating").Scan(&result).Error; err != nil {
		return 0, err
	}

	return result.AvgRating, nil
}

func (s *ReviewService) DeleteReview(reviewerID, reviewID uint) error {
	var review models.Review
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		return errors.New("评价不存在")
	}

	if review.ReviewerID != reviewerID {
		return errors.New("无权删除此评价")
	}

	return database.DB.Delete(&review).Error
}

func (s *ReviewService) CheckReviewed(orderID, repairOrderID *uint, reviewerID uint) bool {
	db := database.DB.Model(&models.Review{}).Where("reviewer_id = ?", reviewerID)
	if orderID != nil {
		db = db.Where("order_id = ?", *orderID)
	}
	if repairOrderID != nil {
		db = db.Where("repair_order_id = ?", *repairOrderID)
	}

	var count int64
	db.Count(&count)
	return count > 0
}

type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

func (s *ReportService) CreateReport(reporterID uint, targetType string, targetID uint, reason, description, images string) (*models.Report, error) {
	report := &models.Report{
		ReporterID:  reporterID,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
		Description: description,
		Images:      images,
		Status:      models.ReportStatusPending,
	}

	result := database.DB.Create(report)
	if result.Error != nil {
		return nil, result.Error
	}

	return report, nil
}

func (s *ReportService) GetReportByID(id uint) (*models.Report, error) {
	var report models.Report
	if err := database.DB.Preload("Reporter").First(&report, id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *ReportService) ListReports(page, pageSize int, status int, targetType string) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	db := database.DB.Model(&models.Report{})
	if status > 0 {
		db = db.Where("status = ?", status)
	}
	if targetType != "" {
		db = db.Where("target_type = ?", targetType)
	}

	db.Count(&total)
	if err := db.Preload("Reporter").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (s *ReportService) HandleReport(reportID uint, approved bool, handleResult string, handledBy uint) error {
	var report models.Report
	if err := database.DB.First(&report, reportID).Error; err != nil {
		return errors.New("举报不存在")
	}

	report.HandledBy = &handledBy
	report.HandleResult = handleResult

	if approved {
		report.Status = models.ReportStatusApproved
		switch report.TargetType {
		case models.ReportTargetProduct:
			database.DB.Model(&models.Product{}).Where("id = ?", report.TargetID).Update("status", models.ProductStatusOffShelf)
		case models.ReportTargetUser:
			database.DB.Model(&models.User{}).Where("id = ?", report.TargetID).Update("status", models.UserStatusFrozen)
		}
	} else {
		report.Status = models.ReportStatusRejected
	}

	return database.DB.Save(&report).Error
}

type WarrantyService struct{}

func NewWarrantyService() *WarrantyService {
	return &WarrantyService{}
}

func (s *WarrantyService) CreateWarrantyClaim(userID uint, orderID, repairOrderID *uint, warrantyType, description, images string) (*models.Warranty, error) {
	warranty := &models.Warranty{
		UserID:        userID,
		OrderID:       orderID,
		RepairOrderID: repairOrderID,
		Type:          warrantyType,
		Description:   description,
		Images:        images,
		Status:        models.WarrantyStatusPending,
	}

	result := database.DB.Create(warranty)
	if result.Error != nil {
		return nil, result.Error
	}

	return warranty, nil
}

func (s *WarrantyService) GetWarrantyByID(id uint) (*models.Warranty, error) {
	var warranty models.Warranty
	if err := database.DB.First(&warranty, id).Error; err != nil {
		return nil, err
	}
	return &warranty, nil
}

func (s *WarrantyService) ListWarranties(page, pageSize int, userID uint, status int) ([]models.Warranty, int64, error) {
	var warranties []models.Warranty
	var total int64

	db := database.DB.Model(&models.Warranty{})
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	db.Count(&total)
	if err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&warranties).Error; err != nil {
		return nil, 0, err
	}

	return warranties, total, nil
}

func (s *WarrantyService) HandleWarranty(warrantyID uint, status int, handleResult string, handledBy uint) error {
	var warranty models.Warranty
	if err := database.DB.First(&warranty, warrantyID).Error; err != nil {
		return errors.New("质保申请不存在")
	}

	warranty.Status = status
	warranty.HandleResult = handleResult
	warranty.HandledBy = &handledBy

	return database.DB.Save(&warranty).Error
}

func init() {
	logrus.Info("Review service initialized")
}
