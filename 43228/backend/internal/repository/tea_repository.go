package repository

import (
	"gorm.io/gorm"

	"tea-platform/internal/models"
	"tea-platform/pkg/database"
)

type TeaFilters struct {
	Type     string
	Origin   string
	Year     int
	Grade    string
	Keyword  string
	SellerID uint
}

type TeaRepository struct{}

func NewTeaRepository() *TeaRepository {
	return &TeaRepository{}
}

func (r *TeaRepository) CreateTea(tea *models.Tea) error {
	return database.GetDB().Create(tea).Error
}

func (r *TeaRepository) GetTeaByID(id uint) (*models.Tea, error) {
	var tea models.Tea
	err := database.GetDB().
		Preload("Images").
		Preload("Traceability").
		First(&tea, id).Error
	if err != nil {
		return nil, err
	}
	return &tea, nil
}

func (r *TeaRepository) UpdateTea(tea *models.Tea) error {
	return database.GetDB().Save(tea).Error
}

func (r *TeaRepository) DeleteTea(id uint) error {
	return database.GetDB().Delete(&models.Tea{}, id).Error
}

func (r *TeaRepository) GetTeaList(page, pageSize int, filters TeaFilters) ([]models.Tea, int64, error) {
	var teas []models.Tea
	var total int64

	db := database.GetDB().Model(&models.Tea{})

	if filters.Type != "" {
		db = db.Where("type = ?", filters.Type)
	}
	if filters.Origin != "" {
		db = db.Where("origin = ?", filters.Origin)
	}
	if filters.Year != 0 {
		db = db.Where("year = ?", filters.Year)
	}
	if filters.Grade != "" {
		db = db.Where("grade = ?", filters.Grade)
	}
	if filters.Keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+filters.Keyword+"%", "%"+filters.Keyword+"%")
	}
	if filters.SellerID != 0 {
		db = db.Where("seller_id = ?", filters.SellerID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Images").Offset(offset).Limit(pageSize).Order("id DESC").Find(&teas).Error; err != nil {
		return nil, 0, err
	}

	return teas, total, nil
}

func (r *TeaRepository) CreateTeaImage(image *models.TeaImage) error {
	return database.GetDB().Create(image).Error
}

func (r *TeaRepository) GetTeaImages(teaID uint) ([]models.TeaImage, error) {
	var images []models.TeaImage
	err := database.GetDB().Where("tea_id = ?", teaID).Order("sort ASC, id ASC").Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (r *TeaRepository) DeleteTeaImage(id uint) error {
	return database.GetDB().Delete(&models.TeaImage{}, id).Error
}

func (r *TeaRepository) CreateTraceability(trace *models.Traceability) error {
	return database.GetDB().Create(trace).Error
}

func (r *TeaRepository) GetTraceabilityByTeaID(teaID uint) ([]models.Traceability, error) {
	var traces []models.Traceability
	err := database.GetDB().Where("tea_id = ?", teaID).Order("record_time ASC").Find(&traces).Error
	if err != nil {
		return nil, err
	}
	return traces, nil
}

func (r *TeaRepository) UpdateTeaStock(id uint, quantity int) error {
	return database.GetDB().Model(&models.Tea{}).Where("id = ?", id).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}
