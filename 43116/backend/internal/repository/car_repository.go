package repository

import (
	"car-rental/internal/model"
	cachedb "car-rental/internal/config"

	"gorm.io/gorm"
)

type CarRepository struct {
	db *gorm.DB
}

func NewCarRepository() *CarRepository {
	return &CarRepository{db: cachedb.DB}
}

func (r *CarRepository) Create(car *model.Car) error {
	return r.db.Create(car).Error
}

func (r *CarRepository) FindByID(id uint) (*model.Car, error) {
	var car model.Car
	err := r.db.Preload("Store.City").Preload("Images").First(&car, id).Error
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (r *CarRepository) FindAll(page, pageSize int, keyword string, status string, brand string, storeID uint) ([]model.Car, int64, error) {
	var cars []model.Car
	var total int64

	query := r.db.Model(&model.Car{})
	if keyword != "" {
		query = query.Where("brand LIKE ? OR model LIKE ? OR license_plate LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if brand != "" {
		query = query.Where("brand = ?", brand)
	}
	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}

	query.Count(&total)
	err := query.Preload("Store.City").Preload("Images").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&cars).Error
	return cars, total, err
}

func (r *CarRepository) Update(car *model.Car) error {
	return r.db.Save(car).Error
}

func (r *CarRepository) UpdateStatus(id uint, status model.CarStatus) error {
	return r.db.Model(&model.Car{}).Where("id = ?", id).Update("status", status).Error
}

func (r *CarRepository) Delete(id uint) error {
	return r.db.Delete(&model.Car{}, id).Error
}

func (r *CarRepository) AddImage(image *model.CarImage) error {
	return r.db.Create(image).Error
}

func (r *CarRepository) DeleteImage(id uint) error {
	return r.db.Delete(&model.CarImage{}, id).Error
}

func (r *CarRepository) GetImages(carID uint) ([]model.CarImage, error) {
	var images []model.CarImage
	err := r.db.Where("car_id = ?", carID).Order("sort_order ASC").Find(&images).Error
	return images, err
}

func (r *CarRepository) ExistsByLicensePlate(licensePlate string) bool {
	var count int64
	r.db.Model(&model.Car{}).Where("license_plate = ?", licensePlate).Count(&count)
	return count > 0
}

func (r *CarRepository) UpdateRating(carID uint) error {
	var reviews []model.Review
	r.db.Where("car_id = ? AND is_hidden = ?", carID, false).Find(&reviews)

	if len(reviews) == 0 {
		return r.db.Model(&model.Car{}).Where("id = ?", carID).Updates(map[string]interface{}{
			"rating":       0,
			"review_count": 0,
		}).Error
	}

	var totalRating int
	for _, review := range reviews {
		totalRating += review.Rating
	}
	avgRating := float64(totalRating) / float64(len(reviews))

	return r.db.Model(&model.Car{}).Where("id = ?", carID).Updates(map[string]interface{}{
		"rating":       avgRating,
		"review_count": len(reviews),
	}).Error
}

func (r *CarRepository) GetAvailableCars(storeID uint, page, pageSize int) ([]model.Car, int64, error) {
	var cars []model.Car
	var total int64

	query := r.db.Model(&model.Car{}).Where("status = ?", model.CarStatusAvailable)
	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}

	query.Count(&total)
	err := query.Preload("Store.City").Preload("Images").Offset((page - 1) * pageSize).Limit(pageSize).Order("rating DESC").Find(&cars).Error
	return cars, total, err
}