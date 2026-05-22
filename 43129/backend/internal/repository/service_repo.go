package repository

import (
	"beauty-salon-system/internal/model"

	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(service *model.Service) error {
	return r.db.Create(service).Error
}

func (r *ServiceRepository) GetByID(id uint) (*model.Service, error) {
	var service model.Service
	err := r.db.First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepository) Update(service *model.Service) error {
	return r.db.Save(service).Error
}

func (r *ServiceRepository) Delete(id uint) error {
	return r.db.Delete(&model.Service{}, id).Error
}

func (r *ServiceRepository) List(page, pageSize int, category string, isPackage bool) ([]model.Service, int64, error) {
	var services []model.Service
	var total int64

	query := r.db.Model(&model.Service{}).Where("status = 1")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if isPackage {
		query = query.Where("is_package = ?", true)
	}
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&services).Error
	return services, total, err
}

func (r *ServiceRepository) ListAll() ([]model.Service, error) {
	var services []model.Service
	err := r.db.Where("status = 1").Find(&services).Error
	return services, err
}

func (r *ServiceRepository) AddPackageService(pkg *model.PackageService) error {
	return r.db.Create(pkg).Error
}

func (r *ServiceRepository) GetPackageServices(serviceID uint) ([]model.PackageService, error) {
	var pkgServices []model.PackageService
	err := r.db.Where("service_id = ?", serviceID).Preload("ChildService").Find(&pkgServices).Error
	return pkgServices, err
}

func (r *ServiceRepository) DeletePackageServices(serviceID uint) error {
	return r.db.Where("service_id = ?", serviceID).Delete(&model.PackageService{}).Error
}
