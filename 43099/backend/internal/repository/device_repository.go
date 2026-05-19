package repository

import (
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/pkg/database"
)

type DeviceRepository struct{}

func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{}
}

func (r *DeviceRepository) CreateCategory(category *model.DeviceCategory) error {
	return database.DB.Create(category).Error
}

func (r *DeviceRepository) GetCategoryByID(id uint) (*model.DeviceCategory, error) {
	var category model.DeviceCategory
	err := database.DB.First(&category, id).Error
	return &category, err
}

func (r *DeviceRepository) GetCategoryByName(name string) (*model.DeviceCategory, error) {
	var category model.DeviceCategory
	err := database.DB.Where("name = ?", name).First(&category).Error
	return &category, err
}

func (r *DeviceRepository) ListCategories() ([]model.DeviceCategory, error) {
	var categories []model.DeviceCategory
	err := database.DB.Order("sort_order ASC, id ASC").Find(&categories).Error
	return categories, err
}

func (r *DeviceRepository) Create(device *model.Device) error {
	return database.DB.Create(device).Error
}

func (r *DeviceRepository) GetByID(id uint) (*model.Device, error) {
	var device model.Device
	err := database.DB.Preload("Category").First(&device, id).Error
	return &device, err
}

func (r *DeviceRepository) List(req *dto.DeviceListRequest) ([]model.Device, int64, error) {
	var devices []model.Device
	var total int64

	query := database.DB.Model(&model.Device{}).Preload("Category")

	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&devices).Error
	return devices, total, err
}

func (r *DeviceRepository) Update(device *model.Device) error {
	return database.DB.Save(device).Error
}

func (r *DeviceRepository) UpdateStatus(id uint, status model.DeviceStatus) error {
	return database.DB.Model(&model.Device{}).Where("id = ?", id).Update("status", status).Error
}

func (r *DeviceRepository) UpdateAvailableQuantity(id uint, quantity int) error {
	return database.DB.Model(&model.Device{}).Where("id = ?", id).
		Update("available_quantity", database.DB.Raw("available_quantity + ?", quantity)).Error
}
