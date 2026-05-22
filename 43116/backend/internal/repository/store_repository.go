package repository

import (
	"car-rental/internal/model"
	cachedb "car-rental/internal/config"

	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository() *StoreRepository {
	return &StoreRepository{db: cachedb.DB}
}

func (r *StoreRepository) CreateStore(store *model.Store) error {
	return r.db.Create(store).Error
}

func (r *StoreRepository) FindStoreByID(id uint) (*model.Store, error) {
	var store model.Store
	err := r.db.Preload("City").First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) FindAllStores(page, pageSize int, cityID uint, keyword string) ([]model.Store, int64, error) {
	var stores []model.Store
	var total int64

	query := r.db.Model(&model.Store{})
	if cityID > 0 {
		query = query.Where("city_id = ?", cityID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Preload("City").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&stores).Error
	return stores, total, err
}

func (r *StoreRepository) UpdateStore(store *model.Store) error {
	return r.db.Save(store).Error
}

func (r *StoreRepository) DeleteStore(id uint) error {
	return r.db.Delete(&model.Store{}, id).Error
}

func (r *StoreRepository) GetStoresByCity(cityID uint) ([]model.Store, error) {
	var stores []model.Store
	err := r.db.Where("city_id = ? AND is_active = ?", cityID, true).Order("name ASC").Find(&stores).Error
	return stores, err
}

func (r *StoreRepository) CreateCity(city *model.City) error {
	return r.db.Create(city).Error
}

func (r *StoreRepository) FindCityByID(id uint) (*model.City, error) {
	var city model.City
	err := r.db.First(&city, id).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *StoreRepository) FindAllCities() ([]model.City, error) {
	var cities []model.City
	err := r.db.Order("name ASC").Find(&cities).Error
	return cities, err
}

func (r *StoreRepository) UpdateCity(city *model.City) error {
	return r.db.Save(city).Error
}

func (r *StoreRepository) DeleteCity(id uint) error {
	return r.db.Delete(&model.City{}, id).Error
}

func (r *StoreRepository) ExistsCityByName(name string) bool {
	var count int64
	r.db.Model(&model.City{}).Where("name = ?", name).Count(&count)
	return count > 0
}

func (r *StoreRepository) ExistsStoreByName(name string) bool {
	var count int64
	r.db.Model(&model.Store{}).Where("name = ?", name).Count(&count)
	return count > 0
}