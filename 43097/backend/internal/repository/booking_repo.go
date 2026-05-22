package repository

import (
	"hotel-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *model.Booking) error
	GetByID(id uint) (*model.Booking, error)
	GetByBookingNo(bookingNo string) (*model.Booking, error)
	Update(booking *model.Booking) error
	Delete(id uint) error
	List(page, pageSize int, bookingNo string, roomID, memberID uint, guestName, guestPhone string, status model.BookingStatus, checkInDate, checkOutDate time.Time) ([]model.Booking, int64, error)
	CheckRoomAvailability(roomID uint, checkIn, checkOut time.Time, excludeBookingID *uint) (bool, error)
	GetBookingsByRoomAndDateRange(roomID uint, startDate, endDate time.Time) ([]model.Booking, error)
	GetActiveBookingsByRoom(roomID uint) ([]model.Booking, error)
	UpdateStatus(id uint, status model.BookingStatus) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) GetByID(id uint) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("Room").Preload("Room.RoomType").Preload("Member").Preload("Member.Level").First(&booking, id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetByBookingNo(bookingNo string) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("Room").Preload("Room.RoomType").Preload("Member").Preload("Member.Level").Where("booking_no = ?", bookingNo).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) Delete(id uint) error {
	return r.db.Delete(&model.Booking{}, id).Error
}

func (r *bookingRepository) List(page, pageSize int, bookingNo string, roomID, memberID uint, guestName, guestPhone string, status model.BookingStatus, checkInDate, checkOutDate time.Time) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var total int64

	query := r.db.Model(&model.Booking{}).Preload("Room").Preload("Room.RoomType").Preload("Member").Preload("Member.Level")

	if bookingNo != "" {
		query = query.Where("booking_no LIKE ?", "%"+bookingNo+"%")
	}
	if roomID > 0 {
		query = query.Where("room_id = ?", roomID)
	}
	if memberID > 0 {
		query = query.Where("member_id = ?", memberID)
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
	if !checkInDate.IsZero() {
		query = query.Where("check_in_date >= ?", checkInDate)
	}
	if !checkOutDate.IsZero() {
		query = query.Where("check_out_date <= ?", checkOutDate)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&bookings).Error
	if err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

func (r *bookingRepository) CheckRoomAvailability(roomID uint, checkIn, checkOut time.Time, excludeBookingID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.Booking{}).Where(
		"room_id = ? AND status IN ? AND check_in_date < ? AND check_out_date > ?",
		roomID,
		[]model.BookingStatus{model.BookingStatusPending, model.BookingStatusConfirmed},
		checkOut,
		checkIn,
	)

	if excludeBookingID != nil {
		query = query.Where("id != ?", *excludeBookingID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *bookingRepository) GetBookingsByRoomAndDateRange(roomID uint, startDate, endDate time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Where(
		"room_id = ? AND status IN ? AND check_in_date < ? AND check_out_date > ?",
		roomID,
		[]model.BookingStatus{model.BookingStatusPending, model.BookingStatusConfirmed},
		endDate,
		startDate,
	).Order("check_in_date ASC").Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) GetActiveBookingsByRoom(roomID uint) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Where(
		"room_id = ? AND status IN ?",
		roomID,
		[]model.BookingStatus{model.BookingStatusPending, model.BookingStatusConfirmed},
	).Order("check_in_date ASC").Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) UpdateStatus(id uint, status model.BookingStatus) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", id).Update("status", status).Error
}
