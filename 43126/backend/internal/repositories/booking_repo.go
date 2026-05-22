package repositories

import (
	"time"

	"meeting-room/internal/models"
	"meeting-room/internal/utils"

	"gorm.io/gorm"
)

type BookingRepository struct{}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{}
}

func (r *BookingRepository) Create(booking *models.Booking) error {
	return utils.DB.Create(booking).Error
}

func (r *BookingRepository) CreateInTx(tx *gorm.DB, booking *models.Booking) error {
	return tx.Create(booking).Error
}

func (r *BookingRepository) FindByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := utils.DB.Preload("Room").Preload("User").First(&booking, id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) Update(booking *models.Booking) error {
	return utils.DB.Save(booking).Error
}

func (r *BookingRepository) Delete(id uint) error {
	return utils.DB.Delete(&models.Booking{}, id).Error
}

func (r *BookingRepository) List(page, pageSize int, userID uint, roomID uint, status int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	db := utils.DB.Model(&models.Booking{}).Preload("Room").Preload("User")
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if roomID > 0 {
		db = db.Where("room_id = ?", roomID)
	}
	if status >= 0 {
		db = db.Where("status = ?", status)
	}

	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("start_time DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) CheckConflict(roomID uint, startTime, endTime time.Time, excludeID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	db := utils.DB.Where("room_id = ? AND status NOT IN ?", roomID, []int{int(models.BookingStatusCancelled)}).
		Where("start_time < ? AND end_time > ?", endTime, startTime)
	if excludeID > 0 {
		db = db.Where("id != ?", excludeID)
	}
	err := db.Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) GetByDateRange(startTime, endTime time.Time, roomID uint, floor string) ([]models.Booking, error) {
	var bookings []models.Booking
	db := utils.DB.Preload("Room").Preload("User").
		Where("start_time >= ? AND start_time < ?", startTime, endTime).
		Where("status NOT IN ?", []int{int(models.BookingStatusCancelled)})

	if roomID > 0 {
		db = db.Where("room_id = ?", roomID)
	}
	if floor != "" {
		db = db.Joins("JOIN rooms ON rooms.id = bookings.room_id").Where("rooms.floor = ?", floor)
	}

	err := db.Order("start_time ASC").Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) UpdateStatus(id uint, status models.BookingStatus) error {
	return utils.DB.Model(&models.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *BookingRepository) GetUpcomingReminders() ([]models.Booking, error) {
	now := time.Now()
	oneHourLater := now.Add(time.Hour)
	var bookings []models.Booking
	err := utils.DB.Preload("Room").Preload("User").
		Where("start_time > ? AND start_time <= ?", now, oneHourLater).
		Where("status = ? AND reminded = ?", models.BookingStatusConfirmed, false).
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) MarkReminded(id uint) error {
	return utils.DB.Model(&models.Booking{}).Where("id = ?", id).Update("reminded", true).Error
}

func (r *BookingRepository) GetStats(startTime, endTime time.Time, department string) ([]models.Booking, error) {
	var bookings []models.Booking
	db := utils.DB.Preload("Room").Preload("User").
		Where("start_time >= ? AND start_time < ?", startTime, endTime).
		Where("status = ?", models.BookingStatusCompleted)
	if department != "" {
		db = db.Joins("JOIN users ON users.id = bookings.user_id").Where("users.department = ?", department)
	}
	err := db.Find(&bookings).Error
	return bookings, err
}
