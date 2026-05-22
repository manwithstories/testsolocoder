package repository

import (
	"hotel-system/internal/model"

	"gorm.io/gorm"
)

type RoomTypeRepository interface {
	Create(roomType *model.RoomType) error
	GetByID(id uint) (*model.RoomType, error)
	GetByName(name string) (*model.RoomType, error)
	Update(roomType *model.RoomType) error
	Delete(id uint) error
	List() ([]model.RoomType, error)
	HasRooms(roomTypeID uint) (bool, error)
}

type RoomRepository interface {
	Create(room *model.Room) error
	GetByID(id uint) (*model.Room, error)
	GetByRoomNo(roomNo string) (*model.Room, error)
	Update(room *model.Room) error
	Delete(id uint) error
	List(page, pageSize int, roomNo string, floor int, roomTypeID uint, status model.RoomStatus) ([]model.Room, int64, error)
	GetAvailableRooms(roomTypeID uint) ([]model.Room, error)
	UpdateStatus(id uint, status model.RoomStatus) error
	BatchCreate(rooms []model.Room) error
	GetRoomsByStatus(status model.RoomStatus) ([]model.Room, error)
	CountByStatus() (map[string]int64, error)
	CountAll() (int64, error)
}

type roomTypeRepository struct {
	db *gorm.DB
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomTypeRepository(db *gorm.DB) RoomTypeRepository {
	return &roomTypeRepository{db: db}
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomTypeRepository) Create(roomType *model.RoomType) error {
	return r.db.Create(roomType).Error
}

func (r *roomTypeRepository) GetByID(id uint) (*model.RoomType, error) {
	var roomType model.RoomType
	err := r.db.First(&roomType, id).Error
	if err != nil {
		return nil, err
	}
	return &roomType, nil
}

func (r *roomTypeRepository) GetByName(name string) (*model.RoomType, error) {
	var roomType model.RoomType
	err := r.db.Where("name = ?", name).First(&roomType).Error
	if err != nil {
		return nil, err
	}
	return &roomType, nil
}

func (r *roomTypeRepository) Update(roomType *model.RoomType) error {
	return r.db.Save(roomType).Error
}

func (r *roomTypeRepository) Delete(id uint) error {
	return r.db.Delete(&model.RoomType{}, id).Error
}

func (r *roomTypeRepository) List() ([]model.RoomType, error) {
	var roomTypes []model.RoomType
	err := r.db.Find(&roomTypes).Error
	if err != nil {
		return nil, err
	}
	return roomTypes, nil
}

func (r *roomTypeRepository) HasRooms(roomTypeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Room{}).Where("room_type_id = ?", roomTypeID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *roomRepository) Create(room *model.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) GetByID(id uint) (*model.Room, error) {
	var room model.Room
	err := r.db.Preload("RoomType").First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) GetByRoomNo(roomNo string) (*model.Room, error) {
	var room model.Room
	err := r.db.Where("room_no = ?", roomNo).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) Update(room *model.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepository) Delete(id uint) error {
	return r.db.Delete(&model.Room{}, id).Error
}

func (r *roomRepository) List(page, pageSize int, roomNo string, floor int, roomTypeID uint, status model.RoomStatus) ([]model.Room, int64, error) {
	var rooms []model.Room
	var total int64

	query := r.db.Model(&model.Room{}).Preload("RoomType")

	if roomNo != "" {
		query = query.Where("room_no LIKE ?", "%"+roomNo+"%")
	}
	if floor > 0 {
		query = query.Where("floor = ?", floor)
	}
	if roomTypeID > 0 {
		query = query.Where("room_type_id = ?", roomTypeID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("floor ASC, room_no ASC").Find(&rooms).Error
	if err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *roomRepository) GetAvailableRooms(roomTypeID uint) ([]model.Room, error) {
	var rooms []model.Room
	query := r.db.Preload("RoomType").Where("status = ?", model.RoomStatusAvailable)

	if roomTypeID > 0 {
		query = query.Where("room_type_id = ?", roomTypeID)
	}

	err := query.Order("floor ASC, room_no ASC").Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) UpdateStatus(id uint, status model.RoomStatus) error {
	return r.db.Model(&model.Room{}).Where("id = ?", id).Update("status", status).Error
}

func (r *roomRepository) BatchCreate(rooms []model.Room) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, room := range rooms {
			if err := tx.Create(&room).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roomRepository) GetRoomsByStatus(status model.RoomStatus) ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.Preload("RoomType").Where("status = ?", status).Order("floor ASC, room_no ASC").Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) CountByStatus() (map[string]int64, error) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result

	err := r.db.Model(&model.Room{}).Select("status, count(*) as count").Group("status").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

func (r *roomRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&model.Room{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
