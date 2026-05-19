package repository

import (
	"gym-management/internal/models"
	"gym-management/internal/pkg/database"
	"time"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *models.Booking) error
	GetByID(id uint) (*models.Booking, error)
	GetByMemberAndSchedule(memberID, scheduleID uint) (*models.Booking, error)
	ListByMember(memberID uint, page, pageSize int) ([]models.Booking, int64, error)
	ListBySchedule(scheduleID uint, page, pageSize int) ([]models.Booking, int64, error)
	Update(booking *models.Booking) error
	Delete(id uint) error
	CountBySchedule(scheduleID uint) (int64, error)
	CheckMemberTimeConflict(memberID uint, startTime, endTime time.Time) (bool, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository() BookingRepository {
	return &bookingRepository{db: database.GetDB()}
}

func (r *bookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) GetByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Member").Preload("Schedule.Course").Preload("Schedule.Course.Coach").First(&booking, id).Error
	return &booking, err
}

func (r *bookingRepository) GetByMemberAndSchedule(memberID, scheduleID uint) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Where("member_id = ? AND schedule_id = ?", memberID, scheduleID).First(&booking).Error
	return &booking, err
}

func (r *bookingRepository) ListByMember(memberID uint, page, pageSize int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).Where("member_id = ?", memberID)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Preload("Schedule.Course").Preload("Schedule.Course.Coach").Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *bookingRepository) ListBySchedule(scheduleID uint, page, pageSize int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).Where("schedule_id = ?", scheduleID)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Preload("Member").Order("created_at ASC").Find(&bookings).Error
	return bookings, total, err
}

func (r *bookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Booking{}, id).Error
}

func (r *bookingRepository) CountBySchedule(scheduleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).Where("schedule_id = ? AND status IN (1, 3)", scheduleID).Count(&count).Error
	return count, err
}

func (r *bookingRepository) CheckMemberTimeConflict(memberID uint, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).
		Where("member_id = ? AND status IN (1, 3)", memberID).
		Joins("JOIN course_schedules ON course_schedules.id = bookings.schedule_id").
		Where("course_schedules.start_time < ? AND course_schedules.end_time > ?", endTime, startTime).
		Count(&count).Error
	return count > 0, err
}

type WaitlistRepository interface {
	Create(waitlist *models.Waitlist) error
	GetByID(id uint) (*models.Waitlist, error)
	GetByMemberAndSchedule(memberID, scheduleID uint) (*models.Waitlist, error)
	ListBySchedule(scheduleID uint) ([]models.Waitlist, error)
	Update(waitlist *models.Waitlist) error
	Delete(id uint) error
	GetMaxPosition(scheduleID uint) (int, error)
	CountBySchedule(scheduleID uint) (int64, error)
}

type waitlistRepository struct {
	db *gorm.DB
}

func NewWaitlistRepository() WaitlistRepository {
	return &waitlistRepository{db: database.GetDB()}
}

func (r *waitlistRepository) Create(waitlist *models.Waitlist) error {
	return r.db.Create(waitlist).Error
}

func (r *waitlistRepository) GetByID(id uint) (*models.Waitlist, error) {
	var waitlist models.Waitlist
	err := r.db.Preload("Member").Preload("Schedule.Course").First(&waitlist, id).Error
	return &waitlist, err
}

func (r *waitlistRepository) GetByMemberAndSchedule(memberID, scheduleID uint) (*models.Waitlist, error) {
	var waitlist models.Waitlist
	err := r.db.Where("member_id = ? AND schedule_id = ? AND status = 1", memberID, scheduleID).First(&waitlist).Error
	return &waitlist, err
}

func (r *waitlistRepository) ListBySchedule(scheduleID uint) ([]models.Waitlist, error) {
	var waitlists []models.Waitlist
	err := r.db.Where("schedule_id = ? AND status = 1", scheduleID).
		Preload("Member").
		Order("position ASC").
		Find(&waitlists).Error
	return waitlists, err
}

func (r *waitlistRepository) Update(waitlist *models.Waitlist) error {
	return r.db.Save(waitlist).Error
}

func (r *waitlistRepository) Delete(id uint) error {
	return r.db.Delete(&models.Waitlist{}, id).Error
}

func (r *waitlistRepository) GetMaxPosition(scheduleID uint) (int, error) {
	var maxPos int
	err := r.db.Model(&models.Waitlist{}).
		Where("schedule_id = ? AND status = 1", scheduleID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPos).Error
	return maxPos, err
}

func (r *waitlistRepository) CountBySchedule(scheduleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Waitlist{}).Where("schedule_id = ? AND status = 1", scheduleID).Count(&count).Error
	return count, err
}

type CheckInRepository interface {
	Create(checkIn *models.CheckIn) error
	GetByID(id uint) (*models.CheckIn, error)
	ListByMember(memberID uint, page, pageSize int) ([]models.CheckIn, int64, error)
	ListByDate(date time.Time) ([]models.CheckIn, error)
	GetByMemberAndDate(memberID uint, date time.Time) ([]models.CheckIn, error)
	Update(checkIn *models.CheckIn) error
}

type checkInRepository struct {
	db *gorm.DB
}

func NewCheckInRepository() CheckInRepository {
	return &checkInRepository{db: database.GetDB()}
}

func (r *checkInRepository) Create(checkIn *models.CheckIn) error {
	return r.db.Create(checkIn).Error
}

func (r *checkInRepository) GetByID(id uint) (*models.CheckIn, error) {
	var checkIn models.CheckIn
	err := r.db.Preload("Member").Preload("Schedule.Course").First(&checkIn, id).Error
	return &checkIn, err
}

func (r *checkInRepository) ListByMember(memberID uint, page, pageSize int) ([]models.CheckIn, int64, error) {
	var checkIns []models.CheckIn
	var total int64

	query := r.db.Model(&models.CheckIn{}).Where("member_id = ?", memberID)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Preload("Schedule.Course").Order("check_in_time DESC").Find(&checkIns).Error
	return checkIns, total, err
}

func (r *checkInRepository) ListByDate(date time.Time) ([]models.CheckIn, error) {
	var checkIns []models.CheckIn
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.Where("check_in_time >= ? AND check_in_time < ?", startOfDay, endOfDay).
		Preload("Member").
		Preload("Schedule.Course").
		Find(&checkIns).Error
	return checkIns, err
}

func (r *checkInRepository) GetByMemberAndDate(memberID uint, date time.Time) ([]models.CheckIn, error) {
	var checkIns []models.CheckIn
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.Where("member_id = ? AND check_in_time >= ? AND check_in_time < ?", memberID, startOfDay, endOfDay).
		Preload("Schedule.Course").
		Find(&checkIns).Error
	return checkIns, err
}

func (r *checkInRepository) Update(checkIn *models.CheckIn) error {
	return r.db.Save(checkIn).Error
}
