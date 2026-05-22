package repositories

import (
	"meeting-room/internal/models"
	"meeting-room/internal/utils"

	"gorm.io/gorm"
)

type RoomRepository struct{}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{}
}

func (r *RoomRepository) Create(room *models.Room) error {
	return utils.DB.Create(room).Error
}

func (r *RoomRepository) FindByID(id uint) (*models.Room, error) {
	var room models.Room
	err := utils.DB.Preload("Photos").First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) Update(room *models.Room) error {
	return utils.DB.Save(room).Error
}

func (r *RoomRepository) Delete(id uint) error {
	return utils.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("room_id = ?", id).Delete(&models.RoomPhoto{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Room{}, id).Error
	})
}

func (r *RoomRepository) List(page, pageSize int, floor string, equipment string) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	db := utils.DB.Model(&models.Room{}).Preload("Photos")
	if floor != "" {
		db = db.Where("floor = ?", floor)
	}
	if equipment != "" {
		db = db.Where("equipment LIKE ?", "%"+equipment+"%")
	}

	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&rooms).Error
	return rooms, total, err
}

func (r *RoomRepository) ListAll() ([]models.Room, error) {
	var rooms []models.Room
	err := utils.DB.Preload("Photos").Where("status = ?", models.RoomStatusActive).Find(&rooms).Error
	return rooms, err
}

func (r *RoomRepository) AddPhoto(photo *models.RoomPhoto) error {
	return utils.DB.Create(photo).Error
}

func (r *RoomRepository) DeletePhoto(id uint) error {
	return utils.DB.Delete(&models.RoomPhoto{}, id).Error
}

func (r *RoomRepository) GetAllFloors() ([]string, error) {
	var floors []string
	err := utils.DB.Model(&models.Room{}).Distinct("floor").Pluck("floor", &floors).Error
	return floors, err
}
