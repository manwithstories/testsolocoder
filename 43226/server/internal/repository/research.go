package repository

import (
	"gorm.io/gorm"

	"museum-server/internal/models"
)

type ResearchRepository struct {
	db *gorm.DB
}

func NewResearchRepository(db *gorm.DB) *ResearchRepository {
	return &ResearchRepository{db: db}
}

func (r *ResearchRepository) CreateApplication(app *models.ResearchApplication) error {
	return r.db.Create(app).Error
}

func (r *ResearchRepository) FindApplicationByID(id uint) (*models.ResearchApplication, error) {
	var app models.ResearchApplication
	err := r.db.Preload("User").Preload("Collection").Preload("Reviewer").First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ResearchRepository) UpdateApplication(app *models.ResearchApplication) error {
	return r.db.Save(app).Error
}

func (r *ResearchRepository) ListApplications(status string, page, pageSize int) ([]models.ResearchApplication, int64, error) {
	var apps []models.ResearchApplication
	var total int64

	db := r.db.Model(&models.ResearchApplication{})
	if status != "" {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)
	db.Preload("User").Preload("Collection").Preload("Reviewer").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&apps)

	return apps, total, nil
}

func (r *ResearchRepository) ListApplicationsByUser(userID uint, page, pageSize int) ([]models.ResearchApplication, int64, error) {
	var apps []models.ResearchApplication
	var total int64

	r.db.Model(&models.ResearchApplication{}).Where("user_id = ?", userID).Count(&total)
	r.db.Preload("User").Preload("Collection").Preload("Reviewer").
		Where("user_id = ?", userID).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&apps)

	return apps, total, nil
}

type MuseumRepository struct {
	db *gorm.DB
}

func NewMuseumRepository(db *gorm.DB) *MuseumRepository {
	return &MuseumRepository{db: db}
}

func (r *MuseumRepository) Create(museum *models.Museum) error {
	return r.db.Create(museum).Error
}

func (r *MuseumRepository) FindByID(id uint) (*models.Museum, error) {
	var museum models.Museum
	err := r.db.First(&museum, id).Error
	if err != nil {
		return nil, err
	}
	return &museum, nil
}

func (r *MuseumRepository) Update(museum *models.Museum) error {
	return r.db.Save(museum).Error
}

func (r *MuseumRepository) List() ([]models.Museum, error) {
	var museums []models.Museum
	err := r.db.Order("id ASC").Find(&museums).Error
	return museums, err
}

func (r *MuseumRepository) Delete(id uint) error {
	return r.db.Delete(&models.Museum{}, id).Error
}
