package repository

import (
	"gorm.io/gorm"

	"recruitment-platform/internal/models"
)

type ResumeRepository struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) *ResumeRepository {
	return &ResumeRepository{db: db}
}

func (r *ResumeRepository) Create(resume *models.Resume) error {
	return r.db.Create(resume).Error
}

func (r *ResumeRepository) FindByID(id uint) (*models.Resume, error) {
	var resume models.Resume
	err := r.db.Preload("User").First(&resume, id).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumeRepository) Update(resume *models.Resume) error {
	return r.db.Save(resume).Error
}

func (r *ResumeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Resume{}, id).Error
}

func (r *ResumeRepository) ListByUserID(userID uint) ([]models.Resume, error) {
	var resumes []models.Resume
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&resumes).Error
	return resumes, err
}

func (r *ResumeRepository) GetDefaultResume(userID uint) (*models.Resume, error) {
	var resume models.Resume
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&resume).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumeRepository) SetDefault(userID, resumeID uint) error {
	r.db.Model(&models.Resume{}).Where("user_id = ?", userID).Update("is_default", false)
	return r.db.Model(&models.Resume{}).Where("id = ?", resumeID).Update("is_default", true).Error
}

func (r *ResumeRepository) Search(keyword string, skills []string, page, pageSize int) ([]models.Resume, int64, error) {
	var resumes []models.Resume
	var total int64

	query := r.db.Model(&models.Resume{})
	if keyword != "" {
		query = query.Where("full_name LIKE ? OR summary LIKE ? OR experience LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Preload("User").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&resumes).Error

	return resumes, total, err
}
