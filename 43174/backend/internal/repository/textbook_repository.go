package repository

import (
	"gorm.io/gorm"

	"campus-trade-platform/internal/models"
)

type TextbookRepository struct {
	db *gorm.DB
}

func NewTextbookRepository(db *gorm.DB) *TextbookRepository {
	return &TextbookRepository{db: db}
}

func (r *TextbookRepository) Create(textbook *models.Textbook) error {
	return r.db.Create(textbook).Error
}

func (r *TextbookRepository) Update(textbook *models.Textbook) error {
	return r.db.Save(textbook).Error
}

func (r *TextbookRepository) Delete(id string) error {
	return r.db.Delete(&models.Textbook{}, "id = ?", id).Error
}

func (r *TextbookRepository) FindByID(id string) (*models.Textbook, error) {
	var textbook models.Textbook
	err := r.db.Preload("Seller").Preload("Category").First(&textbook, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &textbook, nil
}

func (r *TextbookRepository) FindAll(page, pageSize int, keyword, categoryID, status string) ([]models.Textbook, int64, error) {
	var textbooks []models.Textbook
	var total int64

	query := r.db.Model(&models.Textbook{}).Preload("Seller").Preload("Category")

	if keyword != "" {
		query = query.Where("title ILIKE ? OR isbn ILIKE ? OR author ILIKE ? OR course_name ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&textbooks).Error
	if err != nil {
		return nil, 0, err
	}

	return textbooks, total, nil
}

func (r *TextbookRepository) FindByISBN(isbn string) (*models.Textbook, error) {
	var textbook models.Textbook
	err := r.db.Preload("Seller").Preload("Category").
		Where("isbn = ? AND status = ?", isbn, models.TextbookStatusAvailable).
		First(&textbook).Error
	if err != nil {
		return nil, err
	}
	return &textbook, nil
}

func (r *TextbookRepository) FindBySellerID(sellerID string, page, pageSize int) ([]models.Textbook, int64, error) {
	var textbooks []models.Textbook
	var total int64

	query := r.db.Model(&models.Textbook{}).Where("seller_id = ?", sellerID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&textbooks).Error
	if err != nil {
		return nil, 0, err
	}

	return textbooks, total, nil
}

func (r *TextbookRepository) UpdateStatus(id string, status models.TextbookStatus) error {
	return r.db.Model(&models.Textbook{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TextbookRepository) IncrementViewCount(id string) error {
	return r.db.Model(&models.Textbook{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *TextbookRepository) GetPopular(limit int) ([]models.Textbook, error) {
	var textbooks []models.Textbook
	err := r.db.Preload("Seller").Preload("Category").
		Where("status = ?", models.TextbookStatusAvailable).
		Order("view_count DESC").
		Limit(limit).
		Find(&textbooks).Error
	return textbooks, err
}
