package repository

import (
	"time"

	"gorm.io/gorm"

	"tea-platform/internal/models"
	"tea-platform/pkg/database"
)

type TastingFilters struct {
	UserID          uint
	TeaID           uint
	MinOverallScore float64
	MaxOverallScore float64
	StartDate       string
	EndDate         string
	BrewMethod      string
	Keyword         string
}

type TastingRepository struct{}

func NewTastingRepository() *TastingRepository {
	return &TastingRepository{}
}

func (r *TastingRepository) CreateTastingRecord(record *models.TastingRecord) error {
	return database.GetDB().Create(record).Error
}

func (r *TastingRepository) GetTastingRecordByID(id uint) (*models.TastingRecord, error) {
	var record models.TastingRecord
	err := database.GetDB().
		Preload("User").
		Preload("Tea").
		Preload("Images").
		First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *TastingRepository) UpdateTastingRecord(record *models.TastingRecord) error {
	return database.GetDB().Save(record).Error
}

func (r *TastingRepository) DeleteTastingRecord(id uint) error {
	return database.GetDB().Delete(&models.TastingRecord{}, id).Error
}

func (r *TastingRepository) GetTastingRecordList(page, pageSize int, filters TastingFilters) ([]models.TastingRecord, int64, error) {
	var records []models.TastingRecord
	var total int64

	db := database.GetDB().Model(&models.TastingRecord{})

	if filters.UserID != 0 {
		db = db.Where("user_id = ?", filters.UserID)
	}
	if filters.TeaID != 0 {
		db = db.Where("tea_id = ?", filters.TeaID)
	}
	if filters.MinOverallScore > 0 {
		db = db.Where("overall_score >= ?", filters.MinOverallScore)
	}
	if filters.MaxOverallScore > 0 {
		db = db.Where("overall_score <= ?", filters.MaxOverallScore)
	}
	if filters.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", filters.StartDate); err == nil {
			db = db.Where("created_at >= ?", startDate)
		}
	}
	if filters.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", filters.EndDate); err == nil {
			endDate = endDate.Add(24 * time.Hour).Add(-time.Second)
			db = db.Where("created_at <= ?", endDate)
		}
	}
	if filters.BrewMethod != "" {
		db = db.Where("brew_method = ?", filters.BrewMethod)
	}
	if filters.Keyword != "" {
		db = db.Where("notes LIKE ?", "%"+filters.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("User").Preload("Tea").Preload("Images").Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *TastingRepository) CreateTastingImage(image *models.TastingImage) error {
	return database.GetDB().Create(image).Error
}

func (r *TastingRepository) GetTastingImages(recordID uint) ([]models.TastingImage, error) {
	var images []models.TastingImage
	err := database.GetDB().Where("tasting_record_id = ?", recordID).Order("id ASC").Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (r *TastingRepository) DeleteTastingImage(id uint) error {
	return database.GetDB().Delete(&models.TastingImage{}, id).Error
}

func (r *TastingRepository) GetTastingStats(userID uint) (totalCount int, avgScore float64, err error) {
	db := database.GetDB().Model(&models.TastingRecord{}).Where("user_id = ?", userID)

	if err = db.Count(&totalCount).Error; err != nil {
		return 0, 0, err
	}

	if totalCount == 0 {
		return 0, 0, nil
	}

	var result struct {
		AvgScore float64
	}
	err = db.Select("AVG(overall_score) as avg_score").Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}

	return totalCount, result.AvgScore, nil
}
