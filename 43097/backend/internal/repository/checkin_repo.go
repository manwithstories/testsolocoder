package repository

import (
	"hotel-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type CheckInRepository interface {
	Create(checkIn *model.CheckIn) error
	GetByID(id uint) (*model.CheckIn, error)
	GetByCheckInNo(checkInNo string) (*model.CheckIn, error)
	Update(checkIn *model.CheckIn) error
	List(page, pageSize int, checkInNo string, bookingID, roomID uint, guestName, guestPhone string, status model.CheckInStatus, checkInTime, checkOutTime time.Time) ([]model.CheckIn, int64, error)
	GetActiveCheckInByRoom(roomID uint) (*model.CheckIn, error)
	UpdateStatus(id uint, status model.CheckInStatus) error
	AddExtraCharge(checkinID uint, amount float64, description string) error
}

type checkInRepository struct {
	db *gorm.DB
}

func NewCheckInRepository(db *gorm.DB) CheckInRepository {
	return &checkInRepository{db: db}
}

func (r *checkInRepository) Create(checkIn *model.CheckIn) error {
	return r.db.Create(checkIn).Error
}

func (r *checkInRepository) GetByID(id uint) (*model.CheckIn, error) {
	var checkIn model.CheckIn
	err := r.db.Preload("Room").Preload("Room.RoomType").Preload("Booking").First(&checkIn, id).Error
	if err != nil {
		return nil, err
	}
	return &checkIn, nil
}

func (r *checkInRepository) GetByCheckInNo(checkInNo string) (*model.CheckIn, error) {
	var checkIn model.CheckIn
	err := r.db.Preload("Room").Preload("Room.RoomType").Preload("Booking").Where("check_in_no = ?", checkInNo).First(&checkIn).Error
	if err != nil {
		return nil, err
	}
	return &checkIn, nil
}

func (r *checkInRepository) Update(checkIn *model.CheckIn) error {
	return r.db.Save(checkIn).Error
}

func (r *checkInRepository) List(page, pageSize int, checkInNo string, bookingID, roomID uint, guestName, guestPhone string, status model.CheckInStatus, checkInTime, checkOutTime time.Time) ([]model.CheckIn, int64, error) {
	var checkIns []model.CheckIn
	var total int64

	query := r.db.Model(&model.CheckIn{}).Preload("Room").Preload("Room.RoomType").Preload("Booking")

	if checkInNo != "" {
		query = query.Where("check_in_no LIKE ?", "%"+checkInNo+"%")
	}
	if bookingID > 0 {
		query = query.Where("booking_id = ?", bookingID)
	}
	if roomID > 0 {
		query = query.Where("room_id = ?", roomID)
	}
	if guestName != "" {
		query = query.Where("guest_name LIKE ?", "%"+guestName+"%")
	}
	if guestPhone != "" {
		query = query.Where("guest_phone LIKE ?", "%"+guestPhone+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if !checkInTime.IsZero() {
		query = query.Where("check_in_time >= ?", checkInTime)
	}
	if !checkOutTime.IsZero() {
		query = query.Where("expected_check_out <= ?", checkOutTime)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&checkIns).Error
	if err != nil {
		return nil, 0, err
	}

	return checkIns, total, nil
}

func (r *checkInRepository) GetActiveCheckInByRoom(roomID uint) (*model.CheckIn, error) {
	var checkIn model.CheckIn
	err := r.db.Where("room_id = ? AND status = ?", roomID, model.CheckInStatusActive).First(&checkIn).Error
	if err != nil {
		return nil, err
	}
	return &checkIn, nil
}

func (r *checkInRepository) UpdateStatus(id uint, status model.CheckInStatus) error {
	return r.db.Model(&model.CheckIn{}).Where("id = ?", id).Update("status", status).Error
}

func (r *checkInRepository) AddExtraCharge(checkinID uint, amount float64, description string) error {
	return r.db.Model(&model.CheckIn{}).Where("id = ?", checkinID).Updates(map[string]interface{}{
		"extra_charges": gorm.Expr("extra_charges + ?", amount),
		"total_amount":  gorm.Expr("total_amount + ?", amount),
	}).Error
}
