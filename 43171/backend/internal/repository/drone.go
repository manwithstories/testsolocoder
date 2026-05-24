package repository

import (
	"drone-rental/internal/config"
	"drone-rental/internal/model"
	"time"
)

type DroneRepo struct{}

func NewDroneRepo() *DroneRepo {
	return &DroneRepo{}
}

func (r *DroneRepo) Create(drone *model.Drone) error {
	return config.DB.Create(drone).Error
}

func (r *DroneRepo) GetByID(id uint) (*model.Drone, error) {
	var drone model.Drone
	err := config.DB.Preload("Owner").First(&drone, id).Error
	return &drone, err
}

func (r *DroneRepo) GetBySerialNo(serialNo string) (*model.Drone, error) {
	var drone model.Drone
	err := config.DB.Where("serial_no = ?", serialNo).First(&drone).Error
	return &drone, err
}

func (r *DroneRepo) Update(drone *model.Drone) error {
	return config.DB.Save(drone).Error
}

func (r *DroneRepo) Delete(id uint) error {
	return config.DB.Delete(&model.Drone{}, id).Error
}

func (r *DroneRepo) UpdateStatus(id uint, status model.DroneStatus) error {
	return config.DB.Model(&model.Drone{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *DroneRepo) UpdateRating(id uint, rating float64) error {
	return config.DB.Model(&model.Drone{}).Where("id = ?", id).
		Update("rating", rating).Error
}

func (r *DroneRepo) List(page, pageSize int, ownerID uint, keyword, region, brand string, status model.DroneStatus) ([]model.Drone, int64, error) {
	var drones []model.Drone
	var total int64
	db := config.DB.Model(&model.Drone{}).Preload("Owner")
	if ownerID != 0 {
		db = db.Where("owner_id = ?", ownerID)
	}
	if keyword != "" {
		db = db.Where("name LIKE ? OR brand LIKE ? OR model LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if region != "" {
		db = db.Where("region = ?", region)
	}
	if brand != "" {
		db = db.Where("brand = ?", brand)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)
	err := db.Order("rating DESC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&drones).Error
	return drones, total, err
}

func (r *DroneRepo) SearchAvailable(startDate, endDate time.Time, region, keyword string, minPrice, maxPrice float64, page, pageSize int) ([]model.Drone, int64, error) {
	var drones []model.Drone
	var total int64
	db := config.DB.Model(&model.Drone{}).Preload("Owner").
		Where("status = ?", model.DroneStatusOnline)

	db = db.Where(`
		NOT EXISTS (
			SELECT 1 FROM orders 
			WHERE orders.drone_id = drones.id 
			AND orders.status IN ('pending', 'paid', 'picked_up')
			AND orders.start_date <= ? 
			AND orders.end_date >= ?
		)
	`, endDate, startDate)

	if !startDate.IsZero() {
		db = db.Where("(available_from IS NULL OR available_from <= ?)", startDate)
	}
	if !endDate.IsZero() {
		db = db.Where("(available_to IS NULL OR available_to >= ?)", endDate)
	}

	if region != "" {
		db = db.Where("region = ?", region)
	}
	if keyword != "" {
		db = db.Where("name LIKE ? OR brand LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if minPrice > 0 {
		db = db.Where("price_per_day >= ?", minPrice)
	}
	if maxPrice > 0 {
		db = db.Where("price_per_day <= ?", maxPrice)
	}

	db.Count(&total)
	err := db.Order("rating DESC, price_per_day ASC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&drones).Error
	return drones, total, err
}

func (r *DroneRepo) BatchCreate(drones []model.Drone) error {
	return config.DB.Create(&drones).Error
}
