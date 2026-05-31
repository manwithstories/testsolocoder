package repository

import (
	"health-platform/models"
	"time"

	"gorm.io/gorm"
)

type PackageRepository struct {
	*BaseRepository
}

func NewPackageRepository() *PackageRepository {
	return &PackageRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *PackageRepository) FindByAgencyID(agencyID uint, page, pageSize int) ([]models.Package, int64, error) {
	var packages []models.Package
	var total int64

	query := r.DB.Model(&models.Package{}).Where("agency_id = ?", agencyID)
	query.Count(&total)

	err := query.Preload("Items").Preload("Agency").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&packages).Error
	return packages, total, err
}

func (r *PackageRepository) GetWithItems(packageID uint) (*models.Package, error) {
	var pkg models.Package
	err := r.DB.Preload("Items").Preload("Agency").First(&pkg, packageID).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *PackageRepository) ListOnlinePackages(page, pageSize int, keyword string) ([]models.Package, int64, error) {
	var packages []models.Package
	var total int64

	query := r.DB.Model(&models.Package{}).Where("status = ?", models.PackageStatusOnline)
	
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	query.Count(&total)

	err := query.Preload("Items").Preload("Agency").
		Order("view_count DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&packages).Error
	return packages, total, err
}

func (r *PackageRepository) UpdateStatus(packageID uint, status models.PackageStatus) error {
	return r.DB.Model(&models.Package{}).Where("id = ?", packageID).
		Update("status", status).Error
}

func (r *PackageRepository) UpdatePrice(packageID uint, price float64) error {
	return r.DB.Model(&models.Package{}).Where("id = ?", packageID).
		Updates(map[string]interface{}{
			"price":      price,
			"updated_at": time.Now(),
		}).Error
}

func (r *PackageRepository) IncrementViewCount(packageID uint) error {
	return r.DB.Model(&models.Package{}).Where("id = ?", packageID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *PackageRepository) IncrementSaleCount(packageID uint) error {
	return r.DB.Model(&models.Package{}).Where("id = ?", packageID).
		UpdateColumn("sale_count", gorm.Expr("sale_count + 1")).Error
}

func (r *PackageRepository) GetHotPackages(limit int) ([]models.Package, error) {
	var packages []models.Package
	err := r.DB.Where("status = ?", models.PackageStatusOnline).
		Preload("Agency").
		Order("view_count DESC").
		Limit(limit).
		Find(&packages).Error
	return packages, err
}
