package repository

import (
	"gym-management/internal/models"
	"gym-management/internal/pkg/database"
	"time"

	"gorm.io/gorm"
)

type MembershipRepository interface {
	Create(membership *models.Membership) error
	GetByID(id uint) (*models.Membership, error)
	GetByMemberID(memberID uint) (*models.Membership, error)
	List(page, pageSize int, memberID uint) ([]models.Membership, int64, error)
	Update(membership *models.Membership) error
	Delete(id uint) error
	GetExpiringSoon(days int) ([]models.Membership, error)
}

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository() MembershipRepository {
	return &membershipRepository{db: database.GetDB()}
}

func (r *membershipRepository) Create(membership *models.Membership) error {
	return r.db.Create(membership).Error
}

func (r *membershipRepository) GetByID(id uint) (*models.Membership, error) {
	var membership models.Membership
	err := r.db.Preload("Member").First(&membership, id).Error
	return &membership, err
}

func (r *membershipRepository) GetByMemberID(memberID uint) (*models.Membership, error) {
	var membership models.Membership
	err := r.db.Where("member_id = ? AND status = 1", memberID).Preload("Member").First(&membership).Error
	return &membership, err
}

func (r *membershipRepository) List(page, pageSize int, memberID uint) ([]models.Membership, int64, error) {
	var memberships []models.Membership
	var total int64

	query := r.db.Model(&models.Membership{})
	if memberID > 0 {
		query = query.Where("member_id = ?", memberID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Preload("Member").Find(&memberships).Error
	return memberships, total, err
}

func (r *membershipRepository) Update(membership *models.Membership) error {
	return r.db.Save(membership).Error
}

func (r *membershipRepository) Delete(id uint) error {
	return r.db.Delete(&models.Membership{}, id).Error
}

func (r *membershipRepository) GetExpiringSoon(days int) ([]models.Membership, error) {
	var memberships []models.Membership
	cutoffDate := time.Now().AddDate(0, 0, days)
	err := r.db.Where("status = 1 AND end_date <= ? AND end_date >= ?", cutoffDate, time.Now()).
		Preload("Member").
		Find(&memberships).Error
	return memberships, err
}
