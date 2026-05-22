package services

import (
	"errors"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"

	"gorm.io/gorm"
)

type ReviewService struct {
	db *gorm.DB
}

func NewReviewService() *ReviewService {
	return &ReviewService{
		db: database.GetDB(),
	}
}

type CreateReviewRequest struct {
	AppointmentID uint   `json:"appointment_id" binding:"required"`
	Rating        int    `json:"rating" binding:"required,min=1,max=5"`
	Content       string `json:"content"`
	IsAnonymous   bool   `json:"is_anonymous"`
}

type ReviewListQuery struct {
	DoctorID uint `form:"doctor_id"`
	Page     int  `form:"page,default=1"`
	PageSize int  `form:"page_size,default=10"`
}

func (s *ReviewService) CreateReview(patientID uint, req CreateReviewRequest) (*models.Review, error) {
	var existingReview models.Review
	err := s.db.Where("appointment_id = ?", req.AppointmentID).First(&existingReview).Error
	if err == nil {
		return nil, errors.New("该预约已存在评价")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var appointment models.Appointment
	if err := s.db.First(&appointment, req.AppointmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预约不存在")
		}
		return nil, err
	}

	if appointment.Status != models.AppointmentCompleted {
		return nil, errors.New("只能对已完成的预约进行评价")
	}

	if appointment.PatientID != patientID {
		return nil, errors.New("只能评价自己的预约")
	}

	review := &models.Review{
		AppointmentID: req.AppointmentID,
		PatientID:     patientID,
		DoctorID:      appointment.DoctorID,
		Rating:        req.Rating,
		Content:       req.Content,
		IsAnonymous:   req.IsAnonymous,
		IsVerified:    true,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(review).Error; err != nil {
			return err
		}

		if err := s.updateDoctorRating(tx, appointment.DoctorID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetReviewByID(review.ID)
}

func (s *ReviewService) GetReviewByID(id uint) (*models.Review, error) {
	var review models.Review
	err := s.db.Preload("Appointment").
		Preload("Patient.User").
		Preload("Doctor.User").
		First(&review, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评价不存在")
		}
		return nil, err
	}
	return &review, nil
}

func (s *ReviewService) GetReviewList(query ReviewListQuery) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	db := s.db.Model(&models.Review{}).
		Preload("Appointment").
		Preload("Patient.User").
		Preload("Doctor.User")

	if query.DoctorID > 0 {
		db = db.Where("doctor_id = ?", query.DoctorID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("created_at DESC").
		Scopes(database.Paginate(query.Page, query.PageSize)).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (s *ReviewService) DeleteReview(id uint) error {
	review, err := s.GetReviewByID(id)
	if err != nil {
		return err
	}

	doctorID := review.DoctorID

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(review).Error; err != nil {
			return err
		}

		if err := s.updateDoctorRating(tx, doctorID); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *ReviewService) updateDoctorRating(tx *gorm.DB, doctorID uint) error {
	var avgRating float64
	var reviewCount int64

	if err := tx.Model(&models.Review{}).
		Where("doctor_id = ?", doctorID).
		Count(&reviewCount).Error; err != nil {
		return err
	}

	if reviewCount > 0 {
		row := tx.Model(&models.Review{}).
			Where("doctor_id = ?", doctorID).
			Select("COALESCE(AVG(rating), 0)").
			Row()
		if err := row.Scan(&avgRating); err != nil {
			return err
		}
	}

	updates := make(map[string]interface{})
	updates["average_rating"] = avgRating
	updates["review_count"] = reviewCount

	if err := tx.Model(&models.Doctor{}).
		Where("id = ?", doctorID).
		Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
