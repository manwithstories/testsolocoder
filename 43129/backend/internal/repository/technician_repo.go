package repository

import (
	"beauty-salon-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type TechnicianRepository struct {
	db *gorm.DB
}

func NewTechnicianRepository(db *gorm.DB) *TechnicianRepository {
	return &TechnicianRepository{db: db}
}

func (r *TechnicianRepository) Create(tech *model.Technician) error {
	return r.db.Create(tech).Error
}

func (r *TechnicianRepository) GetByID(id uint) (*model.Technician, error) {
	var tech model.Technician
	err := r.db.Preload("User").First(&tech, id).Error
	if err != nil {
		return nil, err
	}
	return &tech, nil
}

func (r *TechnicianRepository) GetByUserID(userID uint) (*model.Technician, error) {
	var tech model.Technician
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&tech).Error
	if err != nil {
		return nil, err
	}
	return &tech, nil
}

func (r *TechnicianRepository) Update(tech *model.Technician) error {
	return r.db.Save(tech).Error
}

func (r *TechnicianRepository) UpdateRating(id uint, rating float64, reviewCount int) error {
	return r.db.Model(&model.Technician{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"rating":       rating,
			"review_count": reviewCount,
		}).Error
}

func (r *TechnicianRepository) List(page, pageSize int, status int) ([]model.Technician, int64, error) {
	var technicians []model.Technician
	var total int64

	query := r.db.Model(&model.Technician{}).Preload("User")
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("rating DESC").Find(&technicians).Error
	return technicians, total, err
}

func (r *TechnicianRepository) ListAll() ([]model.Technician, error) {
	var technicians []model.Technician
	err := r.db.Where("status = 1").Preload("User").Find(&technicians).Error
	return technicians, err
}

func (r *TechnicianRepository) IsOnLeave(technicianID uint, date time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&model.TechnicianLeave{}).
		Where("technician_id = ? AND DATE(leave_date) = ? AND approved = ?",
			technicianID, date.Format("2006-01-02"), true).
		Count(&count).Error
	return count > 0, err
}

func (r *TechnicianRepository) AddLeave(leave *model.TechnicianLeave) error {
	return r.db.Create(leave).Error
}

func (r *TechnicianRepository) GetLeaves(technicianID uint, month int) ([]model.TechnicianLeave, error) {
	var leaves []model.TechnicianLeave
	query := r.db.Where("technician_id = ?", technicianID)
	if month > 0 {
		query = query.Where("MONTH(leave_date) = ?", month)
	}
	err := query.Order("leave_date DESC").Find(&leaves).Error
	return leaves, err
}
