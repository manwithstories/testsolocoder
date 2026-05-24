package repository

import (
	"time"

	"gorm.io/gorm"

	"recruitment-platform/internal/models"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(application *models.Application) error {
	return r.db.Create(application).Error
}

func (r *ApplicationRepository) FindByID(id uint) (*models.Application, error) {
	var application models.Application
	err := r.db.Preload("Job.Company").Preload("Applicant.Profile").Preload("Resume").
		First(&application, id).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *ApplicationRepository) Update(application *models.Application) error {
	return r.db.Save(application).Error
}

func (r *ApplicationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Application{}, id).Error
}

func (r *ApplicationRepository) FindByJobAndApplicant(jobID, applicantID uint) (*models.Application, error) {
	var application models.Application
	err := r.db.Where("job_id = ? AND applicant_id = ?", jobID, applicantID).First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *ApplicationRepository) ListByApplicant(applicantID uint, page, pageSize int, status models.ApplicationStatus) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	query := r.db.Model(&models.Application{}).Where("applicant_id = ?", applicantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Preload("Job.Company").Preload("Resume").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("applied_at DESC").
		Find(&applications).Error

	return applications, total, err
}

func (r *ApplicationRepository) ListByJob(jobID uint, page, pageSize int, status models.ApplicationStatus) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	query := r.db.Model(&models.Application{}).Where("job_id = ?", jobID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Preload("Applicant.Profile").Preload("Resume").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("applied_at DESC").
		Find(&applications).Error

	return applications, total, err
}

func (r *ApplicationRepository) ListByCompany(companyID uint, page, pageSize int, status models.ApplicationStatus, keyword string) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	query := r.db.Model(&models.Application{}).
		Joins("JOIN jobs ON jobs.id = applications.job_id").
		Where("jobs.company_id = ?", companyID)

	if status != "" {
		query = query.Where("applications.status = ?", status)
	}
	if keyword != "" {
		query = query.Joins("JOIN users ON users.id = applications.applicant_id").
			Where("users.email LIKE ? OR applications.cover_letter LIKE ?",
				"%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Preload("Job.Company").Preload("Applicant.Profile").Preload("Resume").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("applied_at DESC").
		Find(&applications).Error

	return applications, total, err
}

func (r *ApplicationRepository) UpdateStatus(id uint, status models.ApplicationStatus) error {
	return r.db.Model(&models.Application{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        status,
			"last_update_at": time.Now(),
		}).Error
}

func (r *ApplicationRepository) BulkUpdateStatus(ids []uint, status models.ApplicationStatus) error {
	return r.db.Model(&models.Application{}).Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":        status,
			"last_update_at": time.Now(),
		}).Error
}

func (r *ApplicationRepository) AddHistory(history *models.ApplicationHistory) error {
	return r.db.Create(history).Error
}

func (r *ApplicationRepository) GetHistory(applicationID uint) ([]models.ApplicationHistory, error) {
	var histories []models.ApplicationHistory
	err := r.db.Where("application_id = ?", applicationID).
		Order("created_at DESC").Find(&histories).Error
	return histories, err
}

func (r *ApplicationRepository) GetStatusCountByJob(jobID uint) (map[string]int64, error) {
	var results []struct {
		Status models.ApplicationStatus
		Count  int64
	}

	err := r.db.Model(&models.Application{}).
		Where("job_id = ?", jobID).
		Select("status, count(*) as count").
		Group("status").
		Scan(&results).Error

	counts := make(map[string]int64)
	for _, r := range results {
		counts[string(r.Status)] = r.Count
	}

	return counts, err
}

func (r *ApplicationRepository) GetApplicationStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	var total int64
	var byStatus []struct {
		Status models.ApplicationStatus
		Count  int64
	}

	query := r.db.Model(&models.Application{})
	if !startDate.IsZero() {
		query = query.Where("applied_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("applied_at <= ?", endDate)
	}

	query.Count(&total)
	query.Select("status, count(*) as count").Group("status").Scan(&byStatus)

	statusCount := make(map[string]int64)
	for _, s := range byStatus {
		statusCount[string(s.Status)] = s.Count
	}

	return map[string]interface{}{
		"total":     total,
		"by_status": statusCount,
	}, nil
}
