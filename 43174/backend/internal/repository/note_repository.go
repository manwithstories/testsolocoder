package repository

import (
	"gorm.io/gorm"

	"campus-trade-platform/internal/models"
)

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Create(note *models.Note) error {
	return r.db.Create(note).Error
}

func (r *NoteRepository) Update(note *models.Note) error {
	return r.db.Save(note).Error
}

func (r *NoteRepository) Delete(id string) error {
	return r.db.Delete(&models.Note{}, "id = ?", id).Error
}

func (r *NoteRepository) FindByID(id string) (*models.Note, error) {
	var note models.Note
	err := r.db.Preload("Uploader").Preload("Category").First(&note, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepository) FindAll(page, pageSize int, keyword, subject, categoryID string, isFeatured bool) ([]models.Note, int64, error) {
	var notes []models.Note
	var total int64

	query := r.db.Model(&models.Note{}).Preload("Uploader").Preload("Category")

	if keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if subject != "" {
		query = query.Where("subject = ?", subject)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if isFeatured {
		query = query.Where("is_featured = ?", true)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notes).Error
	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *NoteRepository) FindByUploaderID(uploaderID string, page, pageSize int) ([]models.Note, int64, error) {
	var notes []models.Note
	var total int64

	query := r.db.Model(&models.Note{}).Where("uploader_id = ?", uploaderID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notes).Error
	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *NoteRepository) IncrementDownloadCount(id string) error {
	return r.db.Model(&models.Note{}).Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}

func (r *NoteRepository) IncrementViewCount(id string) error {
	return r.db.Model(&models.Note{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *NoteRepository) UpdateRating(id string, rating float64, ratingCount int) error {
	return r.db.Model(&models.Note{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"rating":       rating,
			"rating_count": ratingCount,
		}).Error
}

func (r *NoteRepository) SetFeatured(id string, isFeatured bool) error {
	return r.db.Model(&models.Note{}).Where("id = ?", id).
		Update("is_featured", isFeatured).Error
}

func (r *NoteRepository) GetFeatured(limit int) ([]models.Note, error) {
	var notes []models.Note
	err := r.db.Preload("Uploader").Preload("Category").
		Where("is_featured = ?", true).
		Order("rating DESC").
		Limit(limit).
		Find(&notes).Error
	return notes, err
}
