package repositories

import (
	"consultation-platform/models"
	"consultation-platform/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConsultRecordRepository struct {
	db *gorm.DB
}

func NewConsultRecordRepository() *ConsultRecordRepository {
	return &ConsultRecordRepository{db: utils.GetDB()}
}

func (r *ConsultRecordRepository) Create(record *models.ConsultRecord) error {
	return r.db.Create(record).Error
}

func (r *ConsultRecordRepository) FindByID(id uuid.UUID) (*models.ConsultRecord, error) {
	var record models.ConsultRecord
	err := r.db.Preload("Client").Preload("Professional").Preload("Appointment").First(&record, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *ConsultRecordRepository) FindByAppointmentID(appointmentID uuid.UUID) (*models.ConsultRecord, error) {
	var record models.ConsultRecord
	err := r.db.Where("appointment_id = ?", appointmentID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *ConsultRecordRepository) Update(record *models.ConsultRecord) error {
	return r.db.Save(record).Error
}

func (r *ConsultRecordRepository) FindByClientID(clientID uuid.UUID, page, pageSize int) ([]models.ConsultRecord, int64, error) {
	var records []models.ConsultRecord
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.ConsultRecord{}).Where("client_id = ?", clientID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").
		Preload("Professional").
		Preload("Appointment").
		Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *ConsultRecordRepository) FindByProfessionalID(professionalID uuid.UUID, page, pageSize int) ([]models.ConsultRecord, int64, error) {
	var records []models.ConsultRecord
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.ConsultRecord{}).Where("professional_id = ?", professionalID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").
		Preload("Client").
		Preload("Appointment").
		Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{db: utils.GetDB()}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) FindByID(id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("Client").Preload("Professional").Preload("Service").First(&review, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *ReviewRepository) FindByProfessionalID(professionalID uuid.UUID, page, pageSize int, status string) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.Review{}).Where("professional_id = ?", professionalID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").
		Preload("Client").
		Preload("Service").
		Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *ReviewRepository) FindByServiceID(serviceID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.Review{}).Where("service_id = ? AND status = ?", serviceID, models.ReviewApproved)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").
		Preload("Client").
		Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *ReviewRepository) FindPendingReviews(page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.Review{}).Where("status = ?", models.ReviewPending)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").
		Preload("Client").
		Preload("Professional").
		Preload("Service").
		Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *ReviewRepository) UpdateStatus(id uuid.UUID, status models.ReviewStatus, rejectReason string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if rejectReason != "" {
		updates["reject_reason"] = rejectReason
	}
	return r.db.Model(&models.Review{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ReviewRepository) GetAverageRatingByProfessionalID(professionalID uuid.UUID) (float64, int, error) {
	var result struct {
		AvgRating float64
		Count     int
	}

	err := r.db.Model(&models.Review{}).
		Where("professional_id = ? AND status = ?", professionalID, models.ReviewApproved).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as count").
		Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}

	return result.AvgRating, result.Count, nil
}

func (r *ReviewRepository) GetAverageRatingByServiceID(serviceID uuid.UUID) (float64, int, error) {
	var result struct {
		AvgRating float64
		Count     int
	}

	err := r.db.Model(&models.Review{}).
		Where("service_id = ? AND status = ?", serviceID, models.ReviewApproved).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as count").
		Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}

	return result.AvgRating, result.Count, nil
}

func (r *ReviewRepository) FindByAppointmentID(appointmentID uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("appointment_id = ?", appointmentID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}
