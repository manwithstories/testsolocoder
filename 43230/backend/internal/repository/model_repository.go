package repository

import (
	"context"
	"errors"
	"time"

	"print3d-platform/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) Create(ctx context.Context, model *models.Model3D) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ModelRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Model3D, error) {
	var model models.Model3D
	err := r.db.WithContext(ctx).
		Preload("Designer").
		Preload("PreviousVersions").
		First(&model, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *ModelRepository) Update(ctx context.Context, model *models.Model3D) error {
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *ModelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Model3D{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("model not found")
	}
	return nil
}

func (r *ModelRepository) List(ctx context.Context, filter ModelFilter, page, pageSize int) ([]models.Model3D, int64, error) {
	var modelsList []models.Model3D
	var total int64

	query := r.db.WithContext(ctx).Preload("Designer").Where("status = ?", models.ModelStatusPublished)

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	if filter.DesignerID != uuid.Nil {
		query = query.Where("designer_id = ?", filter.DesignerID)
	}
	if len(filter.Tags) > 0 {
		query = query.Where("tags && ?", filter.Tags)
	}
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	query.Count(&total)

	switch filter.SortBy {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "rating_desc":
		query = query.Order("rating DESC")
	case "downloads_desc":
		query = query.Order("download_count DESC")
	default:
		query = query.Order("created_at DESC")
	}

	if filter.IsFeatured {
		query = query.Where("is_featured = ?", true)
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&modelsList).Error
	return modelsList, total, err
}

type ModelFilter struct {
	Category    string
	Keyword     string
	DesignerID  uuid.UUID
	Tags        []string
	MinPrice    float64
	MaxPrice    float64
	SortBy      string
	IsFeatured  bool
}

func (r *ModelRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Model3D{}).
		Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *ModelRepository) IncrementDownloadCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Model3D{}).
		Where("id = ?", id).
		Update("download_count", gorm.Expr("download_count + 1")).Error
}

func (r *ModelRepository) IncrementPurchaseCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Model3D{}).
		Where("id = ?", id).
		Update("purchase_count", gorm.Expr("purchase_count + 1")).Error
}

func (r *ModelRepository) UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error {
	return r.db.WithContext(ctx).Model(&models.Model3D{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"rating":       rating,
			"rating_count": gorm.Expr("rating_count + 1"),
		}).Error
}

func (r *ModelRepository) CreateVersion(ctx context.Context, version *models.ModelVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *ModelRepository) GetVersions(ctx context.Context, modelID uuid.UUID) ([]models.ModelVersion, error) {
	var versions []models.ModelVersion
	err := r.db.WithContext(ctx).
		Where("model_id = ?", modelID).
		Order("created_at DESC").
		Find(&versions).Error
	return versions, err
}

func (r *ModelRepository) CreatePurchase(ctx context.Context, purchase *models.ModelPurchase) error {
	return r.db.WithContext(ctx).Create(purchase).Error
}

func (r *ModelRepository) GetPurchase(ctx context.Context, modelID, userID uuid.UUID) (*models.ModelPurchase, error) {
	var purchase models.ModelPurchase
	err := r.db.WithContext(ctx).
		Where("model_id = ? AND user_id = ?", modelID, userID).
		Order("created_at DESC").
		First(&purchase).Error
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}

func (r *ModelRepository) GetUserPurchases(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.ModelPurchase, int64, error) {
	var purchases []models.ModelPurchase
	var total int64

	query := r.db.WithContext(ctx).Preload("Model").Preload("Model.Designer").Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&purchases).Error
	return purchases, total, err
}

func (r *ModelRepository) IsPurchased(ctx context.Context, modelID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ModelPurchase{}).
		Where("model_id = ? AND user_id = ?", modelID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *ModelRepository) AddFavorite(ctx context.Context, favorite *models.ModelFavorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *ModelRepository) RemoveFavorite(ctx context.Context, modelID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("model_id = ? AND user_id = ?", modelID, userID).
		Delete(&models.ModelFavorite{}).Error
}

func (r *ModelRepository) IsFavorite(ctx context.Context, modelID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ModelFavorite{}).
		Where("model_id = ? AND user_id = ?", modelID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *ModelRepository) GetUserFavorites(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Model3D, int64, error) {
	var modelsList []models.Model3D
	var total int64

	subQuery := r.db.WithContext(ctx).
		Model(&models.ModelFavorite{}).
		Where("user_id = ?", userID).
		Select("model_id")

	query := r.db.WithContext(ctx).Preload("Designer").Where("id IN (?)", subQuery)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&modelsList).Error
	return modelsList, total, err
}

func (r *ModelRepository) CreateDownloadRecord(ctx context.Context, record *models.DownloadRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *ModelRepository) GetHotModels(ctx context.Context, limit int) ([]models.Model3D, error) {
	var modelsList []models.Model3D
	err := r.db.WithContext(ctx).
		Preload("Designer").
		Where("status = ?", models.ModelStatusPublished).
		Order("download_count DESC, view_count DESC").
		Limit(limit).
		Find(&modelsList).Error
	return modelsList, err
}

func (r *ModelRepository) GetDesignerModels(ctx context.Context, designerID uuid.UUID, page, pageSize int) ([]models.Model3D, int64, error) {
	var modelsList []models.Model3D
	var total int64

	query := r.db.WithContext(ctx).Where("designer_id = ?", designerID)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&modelsList).Error
	return modelsList, total, err
}

func (r *ModelRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ModelStatus) error {
	return r.db.WithContext(ctx).Model(&models.Model3D{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *ModelRepository) UpdateFeatured(ctx context.Context, id uuid.UUID, featured bool) error {
	return r.db.WithContext(ctx).Model(&models.Model3D{}).
		Where("id = ?", id).
		Update("is_featured", featured).Error
}

func (r *ModelRepository) GetStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalModels int64
	r.db.WithContext(ctx).Model(&models.Model3D{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&totalModels)
	stats["total_models"] = totalModels

	var totalDownloads int64
	r.db.WithContext(ctx).Model(&models.DownloadRecord{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&totalDownloads)
	stats["total_downloads"] = totalDownloads

	var totalPurchases int64
	r.db.WithContext(ctx).Model(&models.ModelPurchase{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&totalPurchases)
	stats["total_purchases"] = totalPurchases

	var totalRevenue float64
	r.db.WithContext(ctx).Model(&models.ModelPurchase{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)
	stats["total_revenue"] = totalRevenue

	return stats, nil
}
