package repository

import (
	"gorm.io/gorm"

	"museum-server/internal/dto"
	"museum-server/internal/models"
)

type CollectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) Create(collection *models.Collection) error {
	return r.db.Create(collection).Error
}

func (r *CollectionRepository) FindByID(id uint) (*models.Collection, error) {
	var collection models.Collection
	err := r.db.Preload("Category").Preload("Museum").First(&collection, id).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *CollectionRepository) FindByCode(code string) (*models.Collection, error) {
	var collection models.Collection
	err := r.db.Where("code = ?", code).First(&collection).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *CollectionRepository) Update(collection *models.Collection) error {
	return r.db.Save(collection).Error
}

func (r *CollectionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Collection{}, id).Error
}

func (r *CollectionRepository) List(query *dto.CollectionListQuery) ([]models.Collection, int64, error) {
	var collections []models.Collection
	var total int64

	db := r.db.Model(&models.Collection{})

	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ? OR description LIKE ? OR tags LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Era != "" {
		db = db.Where("era LIKE ?", "%"+query.Era+"%")
	}
	if query.Material != "" {
		db = db.Where("material LIKE ?", "%"+query.Material+"%")
	}

	db.Count(&total)

	sortField := "created_at"
	if query.SortBy != "" {
		sortField = query.SortBy
	}
	sortOrder := "DESC"
	if query.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	db.Preload("Category").Preload("Museum").
		Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).
		Order(sortField + " " + sortOrder).
		Find(&collections)

	return collections, total, nil
}

func (r *CollectionRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Collection{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *CollectionRepository) BatchCreate(collections []models.Collection) error {
	return r.db.CreateInBatches(collections, 50).Error
}

func (r *CollectionRepository) CreateCategory(category *models.CollectionCategory) error {
	return r.db.Create(category).Error
}

func (r *CollectionRepository) UpdateCategory(category *models.CollectionCategory) error {
	return r.db.Save(category).Error
}

func (r *CollectionRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&models.CollectionCategory{}, id).Error
}

func (r *CollectionRepository) ListCategories() ([]models.CollectionCategory, error) {
	var categories []models.CollectionCategory
	err := r.db.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (r *CollectionRepository) FindCategoryByID(id uint) (*models.CollectionCategory, error) {
	var category models.CollectionCategory
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CollectionRepository) ListTags() ([]models.CollectionTag, error) {
	var tags []models.CollectionTag
	err := r.db.Find(&tags).Error
	return tags, err
}

func (r *CollectionRepository) CreateTag(tag *models.CollectionTag) error {
	return r.db.Create(tag).Error
}
