package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"skillshare/internal/models"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) FindByID(id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Posting").Preload("Student").Preload("Teacher").
		First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *BookingRepository) ListByUser(userID uuid.UUID, role string, page, pageSize int, status *string) ([]*models.Booking, int64, error) {
	var bookings []*models.Booking
	var total int64

	query := r.db.Model(&models.Booking{})
	if role == "student" {
		query = query.Where("student_id = ?", userID)
	} else if role == "teacher" {
		query = query.Where("teacher_id = ?", userID)
	}

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)
	query.Preload("Posting").Preload("Student").Preload("Teacher").
		Order("scheduled_start DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&bookings)

	return bookings, total, nil
}

func (r *BookingRepository) GetConflictingBookings(teacherID uuid.UUID, start, end time.Time) ([]*models.Booking, error) {
	var bookings []*models.Booking
	err := r.db.Where("teacher_id = ? AND status IN ? AND scheduled_start < ? AND scheduled_end > ?",
		teacherID,
		[]string{string(models.BookingStatusPending), string(models.BookingStatusConfirmed)},
		end,
		start,
	).Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) CompleteBooking(id uuid.UUID, actualStart, actualEnd time.Time) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       models.BookingStatusCompleted,
			"actual_start": actualStart,
			"actual_end":   actualEnd,
		}).Error
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) FindByID(id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("Reviewer").Preload("Reviewee").First(&review, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *BookingRepository) GetReviewsByPosting(postingID uuid.UUID, page, pageSize int) ([]*models.Review, int64, error) {
	var reviews []*models.Review
	var total int64

	query := r.db.Model(&models.Review{}).Where("posting_id = ? AND is_public = ? AND status = ?", postingID, true, "approved")

	query.Count(&total)
	query.Preload("Reviewer").Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&reviews)

	return reviews, total, nil
}

func (r *BookingRepository) GetReviewsByUser(userID uuid.UUID, page, pageSize int) ([]*models.Review, int64, error) {
	var reviews []*models.Review
	var total int64

	query := r.db.Model(&models.Review{}).Where("reviewee_id = ? AND is_public = ? AND status = ?", userID, true, "approved")

	query.Count(&total)
	query.Preload("Reviewer").Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&reviews)

	return reviews, total, nil
}

func (r *BookingRepository) CreateComplaint(complaint *models.Complaint) error {
	return r.db.Create(complaint).Error
}

func (r *BookingRepository) GetComplaints(page, pageSize int, status *string) ([]*models.Complaint, int64, error) {
	var complaints []*models.Complaint
	var total int64

	query := r.db.Model(&models.Complaint{})
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)
	query.Preload("Reporter").Preload("Respondent").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&complaints)

	return complaints, total, nil
}

func (r *BookingRepository) HandleComplaint(id uuid.UUID, handlerID uuid.UUID, result string) error {
	return r.db.Model(&models.Complaint{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     models.ComplaintStatusResolved,
			"handler_id": handlerID,
			"result":     result,
			"handled_at": gorm.Expr("NOW()"),
		}).Error
}
