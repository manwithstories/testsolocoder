package repository

import (
	"time"

	"gorm.io/gorm"

	"recruitment-platform/internal/models"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(job *models.Job) error {
	return r.db.Create(job).Error
}

func (r *JobRepository) FindByID(id uint) (*models.Job, error) {
	var job models.Job
	err := r.db.Preload("Company").First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) Update(job *models.Job) error {
	return r.db.Save(job).Error
}

func (r *JobRepository) Delete(id uint) error {
	return r.db.Delete(&models.Job{}, id).Error
}

func (r *JobRepository) List(page, pageSize int, companyID *uint, status models.JobStatus, keyword string) ([]models.Job, int64, error) {
	var jobs []models.Job
	var total int64

	query := r.db.Model(&models.Job{})
	if companyID != nil {
		query = query.Where("company_id = ?", *companyID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ? OR skills LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Preload("Company").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&jobs).Error

	return jobs, total, err
}

func (r *JobRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Job{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *JobRepository) IncrementApplyCount(id uint) error {
	return r.db.Model(&models.Job{}).Where("id = ?", id).
		UpdateColumn("apply_count", gorm.Expr("apply_count + 1")).Error
}

func (r *JobRepository) UpdateStatus(id uint, status models.JobStatus) error {
	return r.db.Model(&models.Job{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *JobRepository) BulkCreate(jobs []models.Job) error {
	return r.db.Create(&jobs).Error
}

func (r *JobRepository) BulkDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&models.Job{}).Error
}

func (r *JobRepository) BulkUpdateStatus(ids []uint, status models.JobStatus) error {
	return r.db.Model(&models.Job{}).Where("id IN ?", ids).
		Update("status", status).Error
}

func (r *JobRepository) LogView(jobID uint, userID *uint, ipAddress string) error {
	log := &models.JobViewLog{
		JobID:     jobID,
		UserID:    userID,
		IPAddress: ipAddress,
		ViewedAt:  time.Now(),
	}
	return r.db.Create(log).Error
}

func (r *JobRepository) GetViewStats(jobID uint, startDate, endDate time.Time) (int64, error) {
	var count int64
	query := r.db.Model(&models.JobViewLog{}).Where("job_id = ?", jobID)
	if !startDate.IsZero() {
		query = query.Where("viewed_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("viewed_at <= ?", endDate)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *JobRepository) GetJobsByCompanyID(companyID uint) ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.Where("company_id = ?", companyID).Find(&jobs).Error
	return jobs, err
}
