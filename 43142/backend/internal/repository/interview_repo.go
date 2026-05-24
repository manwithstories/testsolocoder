package repository

import (
	"time"

	"gorm.io/gorm"

	"recruitment-platform/internal/models"
)

type InterviewRepository struct {
	db *gorm.DB
}

func NewInterviewRepository(db *gorm.DB) *InterviewRepository {
	return &InterviewRepository{db: db}
}

func (r *InterviewRepository) Create(interview *models.Interview) error {
	return r.db.Create(interview).Error
}

func (r *InterviewRepository) FindByID(id uint) (*models.Interview, error) {
	var interview models.Interview
	err := r.db.Preload("Application.Job.Company").Preload("Application.Applicant.Profile").
		First(&interview, id).Error
	if err != nil {
		return nil, err
	}
	return &interview, nil
}

func (r *InterviewRepository) Update(interview *models.Interview) error {
	return r.db.Save(interview).Error
}

func (r *InterviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Interview{}, id).Error
}

func (r *InterviewRepository) ListByCompany(companyID uint, page, pageSize int, status models.InterviewStatus) ([]models.Interview, int64, error) {
	var interviews []models.Interview
	var total int64

	query := r.db.Model(&models.Interview{}).
		Joins("JOIN applications ON applications.id = interviews.application_id").
		Joins("JOIN jobs ON jobs.id = applications.job_id").
		Where("jobs.company_id = ?", companyID)

	if status != "" {
		query = query.Where("interviews.status = ?", status)
	}

	query.Count(&total)
	err := query.Preload("Application.Job.Company").Preload("Application.Applicant.Profile").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("scheduled_at DESC").
		Find(&interviews).Error

	return interviews, total, err
}

func (r *InterviewRepository) ListByApplicant(applicantID uint, page, pageSize int, status models.InterviewStatus) ([]models.Interview, int64, error) {
	var interviews []models.Interview
	var total int64

	query := r.db.Model(&models.Interview{}).Where("applicant_id = ?", applicantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Preload("Application.Job.Company").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("scheduled_at DESC").
		Find(&interviews).Error

	return interviews, total, err
}

func (r *InterviewRepository) CheckTimeConflict(applicantID uint, scheduledAt time.Time, duration int) (bool, error) {
	var count int64
	endTime := scheduledAt.Add(time.Duration(duration) * time.Minute)

	err := r.db.Model(&models.Interview{}).
		Where("applicant_id = ? AND status IN ?", applicantID,
			[]models.InterviewStatus{models.InterviewStatusPending, models.InterviewStatusAccepted}).
		Where("scheduled_at < ? AND datetime(scheduled_at, '+' || duration || ' minutes') > ?",
			endTime, scheduledAt).
		Count(&count).Error

	return count > 0, err
}

func (r *InterviewRepository) UpdateStatus(id uint, status models.InterviewStatus) error {
	return r.db.Model(&models.Interview{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *InterviewRepository) GetUpcomingInterviews(hours int) ([]models.Interview, error) {
	var interviews []models.Interview
	threshold := time.Now().Add(time.Duration(hours) * time.Hour)

	err := r.db.Preload("Application.Job.Company").Preload("Application.Applicant.Profile").
		Where("status IN ? AND scheduled_at <= ? AND scheduled_at >= ?",
			[]models.InterviewStatus{models.InterviewStatusPending, models.InterviewStatusAccepted},
			threshold, time.Now()).
		Find(&interviews).Error

	return interviews, err
}
