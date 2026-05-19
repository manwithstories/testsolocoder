package repository

import (
	"gym-management/internal/models"
	"gym-management/internal/pkg/database"

	"gorm.io/gorm"
)

type MemberRepository interface {
	Create(member *models.Member) error
	GetByID(id uint) (*models.Member, error)
	GetByPhone(phone string) (*models.Member, error)
	GetByEmail(email string) (*models.Member, error)
	List(page, pageSize int, keyword string) ([]models.Member, int64, error)
	Update(member *models.Member) error
	Delete(id uint) error
	UpdateStatus(id uint, status int) error
	DB() *gorm.DB
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository() MemberRepository {
	return &memberRepository{db: database.GetDB()}
}

func (r *memberRepository) Create(member *models.Member) error {
	return r.db.Create(member).Error
}

func (r *memberRepository) GetByID(id uint) (*models.Member, error) {
	var member models.Member
	err := r.db.Preload("Membership").First(&member, id).Error
	return &member, err
}

func (r *memberRepository) GetByPhone(phone string) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("phone = ?", phone).Preload("Membership").First(&member).Error
	return &member, err
}

func (r *memberRepository) GetByEmail(email string) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("email = ?", email).Preload("Membership").First(&member).Error
	return &member, err
}

func (r *memberRepository) List(page, pageSize int, keyword string) ([]models.Member, int64, error) {
	var members []models.Member
	var total int64

	query := r.db.Model(&models.Member{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Preload("Membership").Find(&members).Error
	return members, total, err
}

func (r *memberRepository) Update(member *models.Member) error {
	return r.db.Save(member).Error
}

func (r *memberRepository) Delete(id uint) error {
	return r.db.Delete(&models.Member{}, id).Error
}

func (r *memberRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&models.Member{}).Where("id = ?", id).Update("status", status).Error
}

func (r *memberRepository) DB() *gorm.DB {
	return r.db
}
