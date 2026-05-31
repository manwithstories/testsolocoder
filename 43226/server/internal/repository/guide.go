package repository

import (
	"time"

	"gorm.io/gorm"

	"museum-server/internal/models"
)

type GuideRepository struct {
	db *gorm.DB
}

func NewGuideRepository(db *gorm.DB) *GuideRepository {
	return &GuideRepository{db: db}
}

func (r *GuideRepository) CreateSchedule(schedule *models.GuideSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *GuideRepository) FindScheduleByID(id uint) (*models.GuideSchedule, error) {
	var schedule models.GuideSchedule
	err := r.db.Preload("Guide").First(&schedule, id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *GuideRepository) ListSchedules(guideID uint, startDate, endDate time.Time) ([]models.GuideSchedule, error) {
	var schedules []models.GuideSchedule
	db := r.db.Where("guide_id = ?", guideID)
	if !startDate.IsZero() {
		db = db.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		db = db.Where("date <= ?", endDate)
	}
	err := db.Order("date ASC, start_time ASC").Find(&schedules).Error
	return schedules, err
}

func (r *GuideRepository) UpdateSchedule(schedule *models.GuideSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *GuideRepository) DeleteSchedule(id uint) error {
	return r.db.Delete(&models.GuideSchedule{}, id).Error
}

func (r *GuideRepository) CreateContent(content *models.GuideContent) error {
	return r.db.Create(content).Error
}

func (r *GuideRepository) FindContentByID(id uint) (*models.GuideContent, error) {
	var content models.GuideContent
	err := r.db.Preload("Collection").First(&content, id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *GuideRepository) ListContents(collectionID, exhibitionID uint, language string) ([]models.GuideContent, error) {
	var contents []models.GuideContent
	db := r.db.Preload("Collection")

	if collectionID > 0 {
		db = db.Where("collection_id = ?", collectionID)
	}
	if exhibitionID > 0 {
		db = db.Where("exhibition_id = ?", exhibitionID)
	}
	if language != "" {
		db = db.Where("language = ?", language)
	}

	err := db.Order("sort_order ASC").Find(&contents).Error
	return contents, err
}

func (r *GuideRepository) UpdateContent(content *models.GuideContent) error {
	return r.db.Save(content).Error
}

func (r *GuideRepository) DeleteContent(id uint) error {
	return r.db.Delete(&models.GuideContent{}, id).Error
}
